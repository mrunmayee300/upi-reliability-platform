package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/kafka"
	"github.com/upi-platform/shared/kafkax"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/metrics"
	"github.com/upi-platform/shared/models"
	"github.com/upi-platform/shared/httpx"
	"github.com/upi-platform/shared/store"
)

const serviceName = "retry-orchestrator"

func main() {
	cfg := config.LoadBase(serviceName, "8085")
	log := logging.New(serviceName)
	ctx := context.Background()

	maxAttempts := config.GetInt("MAX_RETRY_ATTEMPTS", 5)
	initialDelay := config.GetDuration("RETRY_INITIAL_INTERVAL", time.Second)
	maxDelay := config.GetDuration("RETRY_MAX_INTERVAL", 60*time.Second)

	pg, err := store.NewPostgres(ctx, cfg.PostgresDSN)
	if err != nil {
		panic(err)
	}
	defer pg.Close()

	producer := kafkax.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	consumer := kafkax.NewConsumer(cfg.KafkaBrokers, kafka.TopicFailedTransactions, kafka.GroupRetryOrchestrator)
	defer consumer.Close()

	r := httpx.NewRouter()
	httpx.RegisterHealth(r, nil, func() error { return pg.Ping(ctx) })
	go func() { _ = httpx.Run(r, cfg.HTTPPort) }()

	log.Info("retry orchestrator running", slog.Int("max_attempts", maxAttempts))

	for {
		msg, err := consumer.Read(ctx)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		env, err := kafkax.ParseEnvelope(msg)
		if err != nil {
			continue
		}
		var failure models.TransactionFailure
		body, _ := json.Marshal(env.Payload)
		if err := json.Unmarshal(body, &failure); err != nil {
			continue
		}
		if !failure.Enriched {
			continue
		}
		if !failure.Retryable {
			sendDLQ(ctx, producer, pg, env.CorrelationID, failure)
			continue
		}

		attempt := failure.AttemptNumber + 1
		if attempt > maxAttempts {
			sendDLQ(ctx, producer, pg, env.CorrelationID, failure)
			continue
		}

		delay := backoff(initialDelay, maxDelay, attempt)
		scheduled := time.Now().UTC().Add(delay)
		_ = pg.InsertRetryAttempt(ctx, failure.TransactionID, attempt, scheduled)

		go func(f models.TransactionFailure, cid string, att int, d time.Duration) {
			time.Sleep(d)
			if f.OriginalTransaction == nil {
				return
			}
			txn := *f.OriginalTransaction
			f.AttemptNumber = att
			_ = producer.PublishEvent(context.Background(), kafka.TopicRetryTransactions, txn.TransactionID, cid, "",
				models.EventTransactionRetry, serviceName, txn)
			metrics.TxProcessed.WithLabelValues(serviceName, "retry_scheduled").Inc()
		}(failure, env.CorrelationID, attempt, delay)

		metrics.KafkaMessages.WithLabelValues(serviceName, kafka.TopicRetryTransactions, "schedule").Inc()
	}
}

func backoff(initial, max time.Duration, attempt int) time.Duration {
	delay := time.Duration(float64(initial) * math.Pow(2, float64(attempt-1)))
	if delay > max {
		delay = max
	}
	jitter := time.Duration(rand.Float64() * 0.2 * float64(delay))
	return delay + jitter
}

func sendDLQ(ctx context.Context, producer *kafkax.Producer, pg *store.Postgres, correlationID string, failure models.TransactionFailure) {
	history, _ := json.Marshal([]models.TransactionFailure{failure})
	_ = pg.InsertDLQ(ctx, failure.TransactionID, string(failure.FailureCode), correlationID, history)
	_ = producer.PublishEvent(ctx, kafka.TopicDeadLetterEvents, failure.TransactionID, correlationID, "",
		models.EventDeadLetter, serviceName, failure)
	metrics.TxProcessed.WithLabelValues(serviceName, "dlq").Inc()
}

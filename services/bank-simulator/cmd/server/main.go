package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/kafka"
	"github.com/upi-platform/shared/kafkax"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/metrics"
	"github.com/upi-platform/shared/models"
	"github.com/upi-platform/shared/httpx"
)

const serviceName = "bank-simulator"

type bankConfig struct {
	BaseLatencyMs int
	FailureRate   float64
}

var (
	bankConfigs = map[string]bankConfig{
		"HDFC":  {BaseLatencyMs: 80, FailureRate: 0.04},
		"ICICI": {BaseLatencyMs: 100, FailureRate: 0.05},
		"SBI":   {BaseLatencyMs: 150, FailureRate: 0.07},
		"AXIS":  {BaseLatencyMs: 90, FailureRate: 0.04},
		"KOTAK": {BaseLatencyMs: 110, FailureRate: 0.05},
		"PNB":   {BaseLatencyMs: 130, FailureRate: 0.06},
		"BOB":   {BaseLatencyMs: 120, FailureRate: 0.06},
		"YES":   {BaseLatencyMs: 95, FailureRate: 0.05},
	}
	bankMu sync.RWMutex
)

func main() {
	cfg := config.LoadBase(serviceName, "8083")
	log := logging.New(serviceName)
	ctx := context.Background()

	producer := kafkax.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	go consumeTopic(ctx, log, cfg, producer, kafka.TopicValidatedTransactions, kafka.GroupBankSimulator)
	go consumeTopic(ctx, log, cfg, producer, kafka.TopicRetryTransactions, kafka.GroupBankSimulator)
	go publishBankHealth(ctx, log, cfg, producer)

	r := httpx.NewRouter()
	httpx.RegisterHealth(r, nil, nil)
	go func() {
		_ = httpx.Run(r, cfg.HTTPPort)
	}()

	log.Info("bank simulator running", slog.String("port", cfg.HTTPPort))
	select {}
}

func consumeTopic(ctx context.Context, log *slog.Logger, cfg config.Base, producer *kafkax.Producer, topic, group string) {
	consumer := kafkax.NewConsumer(cfg.KafkaBrokers, topic, group)
	defer consumer.Close()

	for {
		msg, err := consumer.Read(ctx)
		if err != nil {
			log.Error("kafka read", slog.Any("error", err))
			time.Sleep(time.Second)
			continue
		}
		env, err := kafkax.ParseEnvelope(msg)
		if err != nil {
			continue
		}
		var txn models.UpiTransaction
		body, _ := json.Marshal(env.Payload)
		if err := json.Unmarshal(body, &txn); err != nil {
			continue
		}
		processPayment(ctx, log, producer, env.CorrelationID, txn)
		metrics.KafkaMessages.WithLabelValues(serviceName, topic, "consume").Inc()
	}
}

func processPayment(ctx context.Context, log *slog.Logger, producer *kafkax.Producer, correlationID string, txn models.UpiTransaction) {
	bankMu.RLock()
	bc, ok := bankConfigs[txn.BankCode]
	bankMu.RUnlock()
	if !ok {
		bc = bankConfig{BaseLatencyMs: 120, FailureRate: 0.05}
	}

	latency := bc.BaseLatencyMs + rand.Intn(50)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	_ = producer.PublishEvent(ctx, kafka.TopicLatencyEvents, txn.TransactionID, correlationID, "",
		models.EventLatencyRecorded, serviceName, models.LatencyEvent{
			TransactionID: txn.TransactionID,
			Hop:           "bank.process",
			LatencyMs:     int64(latency),
			BankCode:      txn.BankCode,
		})

	fail := rand.Float64() < bc.FailureRate
	if fail {
		code := models.FailureBankTimeout
		if rand.Float64() < 0.3 {
			code = models.FailureBankOverload
		}
		failure := models.TransactionFailure{
			TransactionID:       txn.TransactionID,
			FailureCode:         code,
			FailureCause:        string(code),
			Retryable:           true,
			BankCode:            txn.BankCode,
			FailedAt:            time.Now().UTC(),
			AttemptNumber:       0,
			LatencyMs:           int64(latency),
			Enriched:            false,
			OriginalTransaction: &txn,
		}
		_ = producer.PublishEvent(ctx, kafka.TopicFailedTransactions, txn.TransactionID, correlationID, "",
			models.EventTransactionFailed, serviceName, failure)
		metrics.TxProcessed.WithLabelValues(serviceName, "failed").Inc()
		return
	}

	_ = producer.PublishEvent(ctx, kafka.TopicAnalyticsEvents, txn.TransactionID, correlationID, "",
		models.EventAnalyticsMetric, serviceName, models.AnalyticsMetric{
			MetricName: "transaction.success",
			Value:      1,
			BankCode:   txn.BankCode,
		})
	metrics.TxProcessed.WithLabelValues(serviceName, "success").Inc()
}

func publishBankHealth(ctx context.Context, log *slog.Logger, cfg config.Base, producer *kafkax.Producer) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		bankMu.RLock()
		for code, bc := range bankConfigs {
			health := map[string]any{
				"bank_code":       code,
				"status":          "HEALTHY",
				"success_rate":    1 - bc.FailureRate,
				"p95_latency_ms":  bc.BaseLatencyMs + 40,
				"error_rate":      bc.FailureRate,
				"circuit_state":   "CLOSED",
				"recorded_at":     time.Now().UTC(),
			}
			_ = producer.PublishEvent(ctx, kafka.TopicBankHealth, code, uuid.New().String(), "",
				models.EventBankHealth, serviceName, health)
		}
		bankMu.RUnlock()
	}
}

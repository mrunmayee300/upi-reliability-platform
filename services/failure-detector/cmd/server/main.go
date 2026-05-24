package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/kafka"
	"github.com/upi-platform/shared/kafkax"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/metrics"
	"github.com/upi-platform/shared/models"
	"github.com/upi-platform/shared/httpx"
)

const serviceName = "failure-detector"

func main() {
	cfg := config.LoadBase(serviceName, "8084")
	log := logging.New(serviceName)
	ctx := context.Background()

	producer := kafkax.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	consumer := kafkax.NewConsumer(cfg.KafkaBrokers, kafka.TopicFailedTransactions, kafka.GroupFailureDetector)
	defer consumer.Close()

	r := httpx.NewRouter()
	httpx.RegisterHealth(r, nil, nil)
	go func() { _ = httpx.Run(r, cfg.HTTPPort) }()

	log.Info("failure detector running")
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
		if failure.Enriched {
			continue
		}

		failure.Enriched = true
		failure.FailureCause = classifyFailure(failure.FailureCode)

		_ = producer.PublishEvent(ctx, kafka.TopicFailedTransactions, failure.TransactionID, env.CorrelationID, "",
			models.EventTransactionFailed, serviceName, failure)
		_ = producer.PublishEvent(ctx, kafka.TopicAnalyticsEvents, failure.TransactionID, env.CorrelationID, "",
			models.EventAnalyticsMetric, serviceName, models.AnalyticsMetric{
				MetricName: "transaction.failed",
				Value:      1,
				BankCode:   failure.BankCode,
				Labels:     map[string]string{"failure_code": string(failure.FailureCode)},
			})

		metrics.KafkaMessages.WithLabelValues(serviceName, kafka.TopicFailedTransactions, "produce").Inc()
		metrics.TxProcessed.WithLabelValues(serviceName, "enriched").Inc()
	}
}

func classifyFailure(code models.FailureCode) string {
	switch code {
	case models.FailureBankTimeout:
		return "Bank PSP timed out while processing debit request"
	case models.FailureBankOverload:
		return "Bank returned overload/congestion response"
	case models.FailureInsufficientFunds:
		return "Issuer declined due to insufficient balance"
	case models.FailureInvalidVPA:
		return "Payee VPA validation failed at issuer"
	case models.FailureNPCISwitchDown:
		return "NPCI switch unavailable for routing"
	default:
		return "Unclassified processing failure"
	}
}

package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/kafka"
	"github.com/upi-platform/shared/kafkax"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/models"
	"github.com/upi-platform/shared/httpx"
)

const serviceName = "fraud-detector"

func main() {
	cfg := config.LoadBase(serviceName, "8087")
	log := logging.New(serviceName)
	ctx := context.Background()
	producer := kafkax.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()
	consumer := kafkax.NewConsumer(cfg.KafkaBrokers, kafka.TopicValidatedTransactions, kafka.GroupFraudDetector)
	defer consumer.Close()

	go func() {
		r := httpx.NewRouter()
		httpx.RegisterHealth(r, nil, nil)
		_ = httpx.Run(r, cfg.HTTPPort)
	}()

	for {
		msg, err := consumer.Read(ctx)
		if err != nil {
			log.Warn("read", "err", err)
			time.Sleep(time.Second)
			continue
		}
		env, _ := kafkax.ParseEnvelope(msg)
		var txn models.UpiTransaction
		body, _ := json.Marshal(env.Payload)
		_ = json.Unmarshal(body, &txn)
		if txn.AmountPaise > 5000000 && rand.Float64() < 0.1 {
			_ = producer.PublishEvent(ctx, kafka.TopicFraudAlerts, txn.TransactionID, env.CorrelationID, "",
				models.EventFraudAlert, serviceName, map[string]any{
					"alert_id": env.EventID, "transaction_id": txn.TransactionID,
					"anomaly_score": 0.92, "severity": "HIGH", "signals": []string{"amount_outlier"},
					"created_at": time.Now().UTC(),
				})
		}
	}
}

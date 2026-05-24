package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/kafka"
	"github.com/upi-platform/shared/kafkax"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/models"
	"github.com/upi-platform/shared/httpx"
)

const serviceName = "intelligent-routing"

func main() {
	cfg := config.LoadBase(serviceName, "8086")
	_ = logging.New(serviceName)
	ctx := context.Background()
	producer := kafkax.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	go consume(ctx, cfg, producer, kafka.TopicLatencyEvents)
	go consume(ctx, cfg, producer, kafka.TopicBankHealth)

	r := httpx.NewRouter()
	httpx.RegisterHealth(r, nil, nil)
	_ = httpx.Run(r, cfg.HTTPPort)
}

func consume(ctx context.Context, cfg config.Base, producer *kafkax.Producer, topic string) {
	consumer := kafkax.NewConsumer(cfg.KafkaBrokers, topic, kafka.GroupIntelligentRouting)
	defer consumer.Close()
	for {
		msg, err := consumer.Read(ctx)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		env, _ := kafkax.ParseEnvelope(msg)
		body, _ := json.Marshal(env.Payload)
		var raw map[string]any
		if json.Unmarshal(body, &raw) != nil {
			continue
		}
		bank, _ := raw["bank_code"].(string)
		if bank == "" {
			continue
		}
		_ = producer.PublishEvent(ctx, kafka.TopicCongestionEvents, bank, env.CorrelationID, "",
			models.EventCongestionDetected, serviceName, map[string]any{
				"bank_code": bank, "congestion_score": 0.5, "recommended_action": "NONE",
				"source": "ROUTING", "recorded_at": time.Now().UTC(),
			})
	}
}

package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/kafka"
	"github.com/upi-platform/shared/kafkax"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/httpx"
)

const serviceName = "notification"

func main() {
	cfg := config.LoadBase(serviceName, "8090")
	log := logging.New(serviceName)
	ctx := context.Background()
	consumer := kafkax.NewConsumer(cfg.KafkaBrokers, kafka.TopicFraudAlerts, kafka.GroupNotification)
	defer consumer.Close()

	go func() {
		r := httpx.NewRouter()
		httpx.RegisterHealth(r, nil, nil)
		_ = httpx.Run(r, cfg.HTTPPort)
	}()

	for {
		msg, err := consumer.Read(ctx)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		env, _ := kafkax.ParseEnvelope(msg)
		log.Info("alert dispatched", slog.String("type", env.EventType), slog.String("id", env.EventID))
	}
}

package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/upi-platform/analytics/internal/aggregator"
	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/httpx"
	"github.com/upi-platform/shared/kafka"
	"github.com/upi-platform/shared/kafkax"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/metrics"
	"github.com/upi-platform/shared/models"
)

const serviceName = "analytics"

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	cfg := config.LoadBase(serviceName, "8089")
	log := logging.New(serviceName)
	ctx := context.Background()

	agg := aggregator.New()
	var wsMu sync.Mutex
	clients := map[*websocket.Conn]struct{}{}

	broadcast := func() {
		payload := agg.Dashboard()
		wsMu.Lock()
		defer wsMu.Unlock()
		for c := range clients {
			_ = c.WriteJSON(payload)
		}
	}

	go consumeMetrics(ctx, log, cfg, agg, broadcast)
	go consumeBankHealth(ctx, log, cfg, agg, broadcast)

	r := httpx.NewRouter()
	r.Use(corsMiddleware())
	httpx.RegisterHealth(r, nil, nil)

	r.GET("/v1/metrics/summary", func(c *gin.Context) {
		c.JSON(http.StatusOK, agg.Snapshot())
	})
	r.GET("/v1/metrics/banks", func(c *gin.Context) {
		c.JSON(http.StatusOK, agg.Banks())
	})
	r.GET("/v1/metrics/retries", func(c *gin.Context) {
		c.JSON(http.StatusOK, agg.RetryMetrics())
	})
	r.GET("/v1/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, agg.Dashboard())
	})

	r.GET("/v1/ws/live", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		wsMu.Lock()
		clients[conn] = struct{}{}
		wsMu.Unlock()
		_ = conn.WriteJSON(agg.Dashboard())
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					wsMu.Lock()
					delete(clients, conn)
					wsMu.Unlock()
					conn.Close()
					return
				}
			}
		}()
	})

	log.Info("analytics running", slog.String("port", cfg.HTTPPort))
	_ = httpx.Run(r, cfg.HTTPPort)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func consumeMetrics(ctx context.Context, log *slog.Logger, cfg config.Base, agg *aggregator.Aggregator, broadcast func()) {
	topics := []string{kafka.TopicAnalyticsEvents, kafka.TopicLatencyEvents}
	for _, topic := range topics {
		go func(t string) {
			consumer := kafkax.NewConsumer(cfg.KafkaBrokers, t, kafka.GroupAnalytics)
			defer consumer.Close()
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
				switch env.EventType {
				case models.EventAnalyticsMetric:
					var m models.AnalyticsMetric
					body, _ := json.Marshal(env.Payload)
					_ = json.Unmarshal(body, &m)
					switch m.MetricName {
					case "transaction.success":
						agg.RecordSuccess(m.BankCode, env.CorrelationID, 0)
					case "transaction.failed":
						agg.RecordFailure(m.BankCode, env.CorrelationID, 0)
					case "transaction.retry":
						agg.RecordRetry()
					}
				case models.EventLatencyRecorded:
					var l models.LatencyEvent
					body, _ := json.Marshal(env.Payload)
					_ = json.Unmarshal(body, &l)
					agg.RecordLatency(l.LatencyMs)
				case models.EventDeadLetter:
					agg.RecordDLQ()
				}
				metrics.KafkaMessages.WithLabelValues(serviceName, t, "consume").Inc()
				broadcast()
			}
		}(topic)
	}
}

func consumeBankHealth(ctx context.Context, log *slog.Logger, cfg config.Base, agg *aggregator.Aggregator, broadcast func()) {
	consumer := kafkax.NewConsumer(cfg.KafkaBrokers, kafka.TopicBankHealth, kafka.GroupAnalytics+"-banks")
	defer consumer.Close()
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
		body, _ := json.Marshal(env.Payload)
		var raw map[string]any
		if json.Unmarshal(body, &raw) == nil {
			b := aggregator.BankHealth{
				BankCode:     str(raw["bank_code"]),
				Status:       str(raw["status"]),
				SuccessRate:  num(raw["success_rate"]),
				P95LatencyMs: int(num(raw["p95_latency_ms"])),
				ErrorRate:    num(raw["error_rate"]),
				CircuitState: str(raw["circuit_state"]),
			}
			agg.UpdateBankHealth(b)
			broadcast()
		}
	}
}

func str(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func num(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	default:
		return 0
	}
}

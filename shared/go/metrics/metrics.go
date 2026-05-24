package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "upi_platform_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"service", "method", "path", "status"},
	)
	KafkaMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "upi_platform_kafka_messages_total",
			Help: "Kafka messages produced/consumed",
		},
		[]string{"service", "topic", "direction"},
	)
	TxProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "upi_platform_transactions_total",
			Help: "Transactions processed by outcome",
		},
		[]string{"service", "status"},
	)
)

func init() {
	prometheus.MustRegister(HTTPRequests, KafkaMessages, TxProcessed)
}

package kafka

// Topic names — must match docs/kafka/topics.md
const (
	TopicUpiTransactions       = "upi-transactions"
	TopicValidatedTransactions = "validated-transactions"
	TopicFailedTransactions    = "failed-transactions"
	TopicRetryTransactions     = "retry-transactions"
	TopicFraudAlerts           = "fraud-alerts"
	TopicLatencyEvents         = "latency-events"
	TopicBankHealth            = "bank-health"
	TopicCongestionEvents      = "congestion-events"
	TopicAnalyticsEvents       = "analytics-events"
	TopicDeadLetterEvents      = "dead-letter-events"
)

// Consumer group IDs.
const (
	GroupBankSimulator      = "bank-simulator-v1"
	GroupFailureDetector    = "failure-detector-v1"
	GroupRetryOrchestrator  = "retry-orchestrator-v1"
	GroupRetryWorker        = "retry-worker-v1"
	GroupIntelligentRouting = "intelligent-routing-v1"
	GroupFraudDetector      = "fraud-detector-v1"
	GroupAnalytics          = "analytics-aggregator-v1"
	GroupNotification       = "notification-v1"
	GroupAIPrediction       = "ai-prediction-v1"
)

// Standard Kafka headers.
const (
	HeaderCorrelationID = "correlation_id"
	HeaderTraceParent   = "traceparent"
	HeaderIdempotencyKey = "idempotency_key"
	HeaderEventVersion  = "event_version"
	HeaderContentType   = "content-type"
)

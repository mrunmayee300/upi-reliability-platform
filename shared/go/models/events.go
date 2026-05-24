package models

import "time"

// EventEnvelope wraps all Kafka payloads (schema: event-envelope.json).
type EventEnvelope struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	EventVersion   string    `json:"event_version"`
	OccurredAt     time.Time `json:"occurred_at"`
	CorrelationID  string    `json:"correlation_id"`
	TraceID        string    `json:"trace_id,omitempty"`
	SourceService  string    `json:"source_service"`
	Payload        any       `json:"payload"`
}

// Event type constants.
const (
	EventTransactionReceived  = "transaction.received"
	EventTransactionValidated = "transaction.validated"
	EventTransactionFailed    = "transaction.failed"
	EventTransactionRetry     = "transaction.retry"
	EventTransactionSuccess   = "transaction.success"
	EventFraudAlert           = "fraud.alert"
	EventLatencyRecorded      = "latency.recorded"
	EventBankHealth           = "bank.health"
	EventCongestionDetected   = "congestion.detected"
	EventAnalyticsMetric      = "analytics.metric"
	EventDeadLetter           = "transaction.dead_letter"
)

// FailureCode taxonomy.
type FailureCode string

const (
	FailureBankTimeout      FailureCode = "BANK_TIMEOUT"
	FailureBankOverload     FailureCode = "BANK_OVERLOAD"
	FailureInsufficientFunds FailureCode = "INSUFFICIENT_FUNDS"
	FailureInvalidVPA       FailureCode = "INVALID_VPA"
	FailureDuplicateTxn     FailureCode = "DUPLICATE_TXN"
	FailureFraudSuspect     FailureCode = "FRAUD_SUSPECT"
	FailureNPCISwitchDown   FailureCode = "NPCI_SWITCH_DOWN"
	FailureUnknown          FailureCode = "UNKNOWN"
)

// TransactionFailure payload for failed-transactions topic.
type TransactionFailure struct {
	TransactionID       string          `json:"transaction_id"`
	FailureCode         FailureCode     `json:"failure_code"`
	FailureCause        string          `json:"failure_cause"`
	Retryable           bool            `json:"retryable"`
	BankCode            string          `json:"bank_code"`
	FailedAt            time.Time       `json:"failed_at"`
	AttemptNumber       int             `json:"attempt_number"`
	LatencyMs           int64           `json:"latency_ms,omitempty"`
	Enriched            bool            `json:"enriched"`
	OriginalTransaction *UpiTransaction `json:"original_transaction,omitempty"`
}

// LatencyEvent records hop timing.
type LatencyEvent struct {
	TransactionID string `json:"transaction_id"`
	Hop           string `json:"hop"`
	LatencyMs     int64  `json:"latency_ms"`
	BankCode      string `json:"bank_code,omitempty"`
}

// AnalyticsMetric is a unified metric event.
type AnalyticsMetric struct {
	MetricName string  `json:"metric_name"`
	Value      float64 `json:"value"`
	BankCode   string  `json:"bank_code,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

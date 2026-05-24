package aggregator

import (
	"sync"
	"time"
)

type Snapshot struct {
	TPS           float64   `json:"tps"`
	SuccessRate   float64   `json:"success_rate"`
	FailureRate   float64   `json:"failure_rate"`
	RetryRate     float64   `json:"retry_rate"`
	P95LatencyMs  int64     `json:"p95_latency_ms"`
	KafkaLagTotal int64     `json:"kafka_lag_total"`
	Timestamp     time.Time `json:"timestamp"`
}

type BankHealth struct {
	BankCode      string  `json:"bank_code"`
	Status        string  `json:"status"`
	SuccessRate   float64 `json:"success_rate"`
	P95LatencyMs  int     `json:"p95_latency_ms"`
	ErrorRate     float64 `json:"error_rate"`
	CircuitState  string  `json:"circuit_state"`
}

type RetryMetrics struct {
	PendingRetries  int64   `json:"pending_retries"`
	DLQCount24h     int64   `json:"dlq_count_24h"`
	AvgRetryAttempts float64 `json:"avg_retry_attempts"`
}

type LiveEvent struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	Type          string    `json:"type"`
	BankCode      string    `json:"bank_code,omitempty"`
	AmountPaise   int64     `json:"amount_paise,omitempty"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"timestamp"`
}

type DashboardPayload struct {
	Summary      Snapshot       `json:"summary"`
	Banks        []BankHealth   `json:"banks"`
	Retries      RetryMetrics   `json:"retries"`
	RecentEvents []LiveEvent    `json:"recent_events"`
}

type Aggregator struct {
	mu          sync.RWMutex
	success     int64
	failed      int64
	retries     int64
	dlqCount    int64
	latencies   []int64
	banks       map[string]BankHealth
	recent      []LiveEvent
	windowStart time.Time
}

func New() *Aggregator {
	return &Aggregator{
		windowStart: time.Now(),
		banks:       defaultBanks(),
	}
}

func defaultBanks() map[string]BankHealth {
	codes := []string{"HDFC", "ICICI", "SBI", "AXIS", "KOTAK", "PNB", "BOB", "YES"}
	m := make(map[string]BankHealth, len(codes))
	for _, code := range codes {
		m[code] = BankHealth{
			BankCode:     code,
			Status:       "HEALTHY",
			SuccessRate:  0.95,
			P95LatencyMs: 120,
			ErrorRate:    0.05,
			CircuitState: "CLOSED",
		}
	}
	return m
}

func (a *Aggregator) RecordSuccess(bankCode, txnID string, amountPaise int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.success++
	a.pushEvent(LiveEvent{
		ID: txnID, TransactionID: txnID, Type: "transaction.success",
		BankCode: bankCode, AmountPaise: amountPaise, Status: "SUCCESS", Timestamp: time.Now().UTC(),
	})
	a.bumpBank(bankCode, true)
}

func (a *Aggregator) RecordFailure(bankCode, txnID string, amountPaise int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.failed++
	a.pushEvent(LiveEvent{
		ID: txnID, TransactionID: txnID, Type: "transaction.failed",
		BankCode: bankCode, AmountPaise: amountPaise, Status: "FAILED", Timestamp: time.Now().UTC(),
	})
	a.bumpBank(bankCode, false)
}

func (a *Aggregator) RecordRetry() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.retries++
}

func (a *Aggregator) RecordDLQ() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.dlqCount++
}

func (a *Aggregator) RecordLatency(ms int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.latencies = append(a.latencies, ms)
	if len(a.latencies) > 5000 {
		a.latencies = a.latencies[len(a.latencies)-5000:]
	}
}

func (a *Aggregator) UpdateBankHealth(b BankHealth) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.banks[b.BankCode] = b
}

func (a *Aggregator) pushEvent(e LiveEvent) {
	a.recent = append([]LiveEvent{e}, a.recent...)
	if len(a.recent) > 100 {
		a.recent = a.recent[:100]
	}
}

func (a *Aggregator) bumpBank(code string, success bool) {
	b, ok := a.banks[code]
	if !ok {
		b = BankHealth{BankCode: code, Status: "HEALTHY", CircuitState: "CLOSED"}
	}
	if success {
		b.SuccessRate = min(1.0, b.SuccessRate+0.001)
		b.ErrorRate = max(0, b.ErrorRate-0.001)
	} else {
		b.SuccessRate = max(0, b.SuccessRate-0.005)
		b.ErrorRate = min(1.0, b.ErrorRate+0.005)
	}
	if b.ErrorRate > 0.15 {
		b.Status = "CONGESTED"
	} else if b.ErrorRate > 0.08 {
		b.Status = "DEGRADED"
	} else {
		b.Status = "HEALTHY"
	}
	a.banks[code] = b
}

func (a *Aggregator) Snapshot() Snapshot {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.snapshotLocked()
}

func (a *Aggregator) Dashboard() DashboardPayload {
	a.mu.RLock()
	defer a.mu.RUnlock()

	banks := make([]BankHealth, 0, len(a.banks))
	for _, b := range a.banks {
		banks = append(banks, b)
	}

	recent := make([]LiveEvent, len(a.recent))
	copy(recent, a.recent)

	total := a.success + a.failed
	avgRetry := 0.0
	if a.retries > 0 && total > 0 {
		avgRetry = float64(a.retries) / float64(total)
	}

	return DashboardPayload{
		Summary:      a.snapshotLocked(),
		Banks:        banks,
		Retries: RetryMetrics{
			PendingRetries:   a.retries,
			DLQCount24h:      a.dlqCount,
			AvgRetryAttempts: avgRetry,
		},
		RecentEvents: recent,
	}
}

func (a *Aggregator) snapshotLocked() Snapshot {
	total := a.success + a.failed
	elapsed := time.Since(a.windowStart).Seconds()
	if elapsed < 1 {
		elapsed = 1
	}
	tps := float64(total) / elapsed

	var failureRate, successRate, retryRate float64
	if total > 0 {
		failureRate = float64(a.failed) / float64(total)
		successRate = float64(a.success) / float64(total)
		retryRate = float64(a.retries) / float64(total)
	}

	return Snapshot{
		TPS:           tps,
		SuccessRate:   successRate,
		FailureRate:   failureRate,
		RetryRate:     retryRate,
		P95LatencyMs:  p95(a.latencies),
		KafkaLagTotal: 0,
		Timestamp:     time.Now().UTC(),
	}
}

func (a *Aggregator) RetryMetrics() RetryMetrics {
	a.mu.RLock()
	defer a.mu.RUnlock()
	total := a.success + a.failed
	avgRetry := 0.0
	if a.retries > 0 && total > 0 {
		avgRetry = float64(a.retries) / float64(total)
	}
	return RetryMetrics{
		PendingRetries:   a.retries,
		DLQCount24h:      a.dlqCount,
		AvgRetryAttempts: avgRetry,
	}
}

func (a *Aggregator) Banks() []BankHealth {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]BankHealth, 0, len(a.banks))
	for _, b := range a.banks {
		out = append(out, b)
	}
	return out
}

func p95(vals []int64) int64 {
	if len(vals) == 0 {
		return 0
	}
	sorted := append([]int64(nil), vals...)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] < sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	idx := int(float64(len(sorted)-1) * 0.95)
	return sorted[idx]
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

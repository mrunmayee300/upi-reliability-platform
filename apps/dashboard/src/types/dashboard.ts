export interface PlatformSummary {
  tps: number;
  success_rate: number;
  failure_rate: number;
  retry_rate: number;
  p95_latency_ms: number;
  kafka_lag_total: number;
  timestamp: string;
}

export interface BankHealth {
  bank_code: string;
  status: "HEALTHY" | "DEGRADED" | "CONGESTED" | "DOWN";
  success_rate: number;
  p95_latency_ms: number;
  error_rate: number;
  circuit_state: "CLOSED" | "OPEN" | "HALF_OPEN";
}

export interface RetryMetrics {
  pending_retries: number;
  dlq_count_24h: number;
  avg_retry_attempts: number;
}

export interface LiveEvent {
  id: string;
  transaction_id: string;
  type: string;
  bank_code?: string;
  amount_paise?: number;
  status: string;
  timestamp: string;
}

export interface DashboardPayload {
  summary: PlatformSummary;
  banks: BankHealth[];
  retries: RetryMetrics;
  recent_events: LiveEvent[];
}

export const emptyDashboard = (): DashboardPayload => ({
  summary: {
    tps: 0,
    success_rate: 0,
    failure_rate: 0,
    retry_rate: 0,
    p95_latency_ms: 0,
    kafka_lag_total: 0,
    timestamp: new Date().toISOString(),
  },
  banks: [],
  retries: { pending_retries: 0, dlq_count_24h: 0, avg_retry_attempts: 0 },
  recent_events: [],
});

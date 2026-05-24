# Scalability Strategy

## Target SLOs (simulation)

| Metric | Target | Measurement |
|--------|--------|-------------|
| Ingestion throughput | 10K+ TPS (k6) | `tx_ingestion_requests_total` |
| P95 end-to-end latency | < 500ms (no congestion) | Trace span `payment.e2e` |
| Failure detection lag | < 2s | Kafka consumer lag `failure-detector` |
| Retry scheduling accuracy | ±5s of backoff | `retry_scheduled_delay_seconds` |
| Dashboard refresh | < 1s WS push | `dashboard_ws_push_latency` |

## Horizontal scaling levers

| Component | Scale trigger | Mechanism |
|-----------|---------------|-----------|
| Tx Ingestion | CPU > 70%, lag > 5K | K8s HPA on Deployment |
| Bank Simulator | `validated` lag | Consumer group + partition count |
| Failure Detector | `failed` lag | Independent consumer group |
| Retry Orchestrator | Retry queue depth | HPA + bounded worker pool |
| Analytics | `analytics-events` lag | Multiple aggregators (idempotent windows) |
| Kafka | Throughput ceiling | Partition increase (planned rebalance) |

## Partition sizing (initial)

| Topic | Partitions | Replication |
|-------|------------|-------------|
| `upi-transactions` | 24 | 3 (prod) / 1 (local) |
| `validated-transactions` | 24 | 3 / 1 |
| `failed-transactions` | 12 | 3 / 1 |
| `retry-transactions` | 12 | 3 / 1 |
| `analytics-events` | 12 | 3 / 1 |
| `bank-health` | 6 | 3 / 1 |

**Rule of thumb:** partitions ≥ max consumer instances for parallel consumers; key skew monitored via `kafka_partition_lag_skew`.

## Backpressure

1. **API Gateway** — token bucket rate limit (Redis sliding window).
2. **Ingestion** — reject with `429` when internal queue > threshold.
3. **Bank Simulator** — shed load via circuit breaker → `BANK_OVERLOAD` events.
4. **Retry Orchestrator** — cap in-flight retries per bank (`max_retries_inflight_per_bank`).

## Caching

| Cache | TTL | Purpose |
|-------|-----|---------|
| Idempotency keys | 24h | Dedupe |
| Bank routing table | 30s | Hot path routing |
| Congestion scores | 10s | Routing decisions |
| Analytics snapshots | 5s | Dashboard API |

## AI prediction scaling

- Python service stateless; scale on CPU.
- Model artifacts in object storage (Phase 5); in-memory for local.
- Batch inference every 60s; push `congestion-events` forecasts.

## Capacity planning formula

```
Required ingestion pods ≈ ceil(target_TPS / per_pod_TPS)
per_pod_TPS ≈ 2000 (baseline, profile with k6)

Kafka disk/day ≈ avg_event_bytes × TPS × 86400 × retention_days
```

Phase 6 k6 scripts validate these assumptions with benchmark reports.

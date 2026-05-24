# Kafka Topic Architecture

## Cluster configuration (target)

| Setting | Local (docker-compose) | Production (K8s/Strimzi) |
|---------|------------------------|---------------------------|
| Brokers | 1 | 3 |
| `min.insync.replicas` | 1 | 2 |
| `replication.factor` | 1 | 3 |
| Retention | 24h | 72h (audit: 7d for DLQ) |
| Compression | lz4 | lz4 |

## Topic catalog

| Topic | Purpose | Partitions | Retention | Compaction |
|-------|---------|------------|-----------|------------|
| `upi-transactions` | Raw ingested transactions (audit) | 24 | 72h | false |
| `validated-transactions` | Post-validation, ready for bank | 24 | 48h | false |
| `failed-transactions` | Bank/processing failures | 12 | 72h | false |
| `retry-transactions` | Scheduled retry attempts | 12 | 48h | false |
| `fraud-alerts` | Anomaly / fraud signals | 6 | 7d | false |
| `latency-events` | Per-hop latency measurements | 12 | 24h | false |
| `bank-health` | Periodic bank health snapshots | 6 | 24h | true (key) |
| `congestion-events` | Congestion scores & routing decisions | 6 | 48h | false |
| `analytics-events` | Unified metrics stream for Analytics | 12 | 48h | false |
| `dead-letter-events` | Exhausted retries, manual review | 6 | 30d | false |

## Consumer groups

| Group ID | Service | Subscribes | Notes |
|----------|---------|------------|-------|
| `bank-simulator-v1` | Bank Simulator | `validated-transactions` | Horizontally scaled |
| `failure-detector-v1` | Failure Detector | `failed-transactions` | Enriches failure_cause |
| `retry-orchestrator-v1` | Retry Orchestrator | `failed-transactions` | Competing consumer OK with detector* |
| `retry-worker-v1` | Retry Orchestrator | `retry-transactions` | Executes retries |
| `intelligent-routing-v1` | Intelligent Routing | `latency-events`, `bank-health`, `congestion-events` | |
| `fraud-detector-v1` | Fraud Detector | `validated-transactions` | Parallel path |
| `analytics-aggregator-v1` | Analytics | `analytics-events`, `latency-events` | |
| `notification-v1` | Notification | `fraud-alerts`, `dead-letter-events`, `congestion-events` | |
| `ai-prediction-v1` | AI Prediction | `analytics-events`, `bank-health` | Feature ingestion |

\*Failure Detector enriches and republishes; Retry Orchestrator consumes enriched failures. Alternative: single `payment-failures` internal topic (Phase 3 implementation choice documented in ADR-004).

## Retry topic pattern

We use **application-level retry scheduling** (not Kafka retry topics alone):

1. `failed-transactions` → Retry Orchestrator computes `retry_at`.
2. Publishes to `retry-transactions` with header `scheduled_at`.
3. Consumer delays via worker queue / Redis ZSET (not blocking Kafka poll).

Optional future: `retry-transactions-dlq` as alias consumer on same topic with `x-dead-letter` semantics.

## Dead letter queue

`dead-letter-events` receives:

- Max retries exceeded
- Non-retryable failure codes
- Manual poison message quarantine

Headers:

```
x-original-topic: failed-transactions
x-original-partition: 3
x-original-offset: 12890
x-failure-count: 5
```

## Message headers (standard)

| Header | Required | Description |
|--------|----------|-------------|
| `correlation_id` | yes | End-to-end trace correlation |
| `traceparent` | yes | W3C trace context |
| `idempotency_key` | yes* | *Required on write commands |
| `content-type` | yes | `application/json` |
| `event_version` | yes | Schema version e.g. `1.0` |

## Partition assignment examples

```
Key: payer_vpa "alice@paytm" → hash % 24 → partition 7
Key: bank_code "HDFC" → hash % 6 → partition 2
Key: transaction_id "TXN-uuid" → hash % 24 → partition 15
```

## Topic creation script

Phase 2: `scripts/kafka/create-topics.sh` — idempotent topic creation for local and CI.

## Monitoring

| Metric | Alert threshold |
|--------|-----------------|
| `kafka_consumer_group_lag` | > 10000 for 5m |
| `kafka_topic_partition_under_replicated` | > 0 |
| `kafka_log_size_bytes` | > 80% disk |

Grafana dashboard: `dashboards/grafana/kafka-overview.json` (Phase 2).

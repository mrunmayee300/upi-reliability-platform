# Resilience Engineering

## Circuit breaker (Bank Simulator client & Routing)

States: `CLOSED` → `OPEN` → `HALF_OPEN`

| Parameter | Default |
|-----------|---------|
| Failure threshold | 50% over 30s window |
| Open duration | 30s |
| Half-open probes | 3 |

When open: Routing publishes `congestion-events`, Intelligent Routing shifts traffic to alternate `bank_code`.

## Retry policy

```yaml
max_attempts: 5
initial_interval: 1s
multiplier: 2.0
max_interval: 60s
jitter: 0.2  # full jitter
retryable_codes:
  - BANK_TIMEOUT
  - BANK_OVERLOAD
  - NPCI_SWITCH_DOWN
```

Non-retryable → immediate `dead-letter-events` + notification.

## Dead letter queue

Topic: `dead-letter-events`

Payload includes: original event, failure history[], final `failure_cause`, `correlation_id`.

Ops replay via admin API (Phase 3): `POST /v1/admin/dlq/{event_id}/replay`

## Chaos engineering (Phase 6)

| Experiment | Tool | Expected behavior |
|------------|------|-------------------|
| Kafka broker kill | Chaos Mesh / script | Producer buffer, lag alert, recovery |
| Bank 100% timeout | Config toggle | Circuit opens, routing shifts |
| Ingestion pod kill | `kubectl delete pod` | HPA replaces, no duplicate with idempotency |
| Network delay | toxiproxy / tc | Latency metrics spike, AI forecast |

## Disaster recovery (documentation target)

- Kafka: RF=3, min ISR=2
- Postgres: PITR enabled in Terraform module (prod profile)
- Redis: AOF persistence for idempotency (optional local: RDB only)

## Health checks

Every service exposes:

- `GET /health/live` — process up
- `GET /health/ready` — Kafka + DB + Redis connectivity

K8s probes use these endpoints (Phase 5 manifests).

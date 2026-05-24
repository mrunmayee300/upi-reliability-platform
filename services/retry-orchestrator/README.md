# Retry Orchestrator Service

**Port:** 8085 | **Phase 3**

## Responsibilities

- Exponential backoff with full jitter
- Idempotent retry scheduling (Redis ZSET + worker pool)
- Publish `retry-transactions`
- Dead letter to `dead-letter-events` after max attempts
- Retry metrics for Prometheus/Grafana

## Kafka

- **Consumes:** `failed-transactions`, `retry-transactions`
- **Produces:** `retry-transactions`, `dead-letter-events`, `analytics-events`

## Config

```
MAX_RETRY_ATTEMPTS=5
RETRY_INITIAL_INTERVAL=1s
RETRY_MAX_INTERVAL=60s
```

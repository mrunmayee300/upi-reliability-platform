# Analytics Service

**Port:** 8089 | **Phase 3**

## Responsibilities

- Consume `analytics-events`, `latency-events`
- Compute TPS, failure rate, bank health, retry metrics
- Expose REST + WebSocket for dashboard
- Persist rollups to PostgreSQL

## Kafka

- **Consumes:** `analytics-events`, `latency-events`

## Contract

`shared/contracts/openapi/analytics.yaml`

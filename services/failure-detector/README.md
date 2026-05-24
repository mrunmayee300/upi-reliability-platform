# Failure Detection Service

**Port:** 8084 (metrics/health only) | **Phase 3**

## Responsibilities

- Consume `failed-transactions`
- Classify and enrich `failure_cause`
- Map raw errors to `FailureCode` taxonomy
- Emit enriched failures and `analytics-events`

## Kafka

- **Consumes:** `failed-transactions`
- **Produces:** `failed-transactions` (enriched), `analytics-events`

See ADR-004 for detector vs orchestrator split.

# Bank Simulation Service

**Port:** 8083 | **Phase 3**

## Responsibilities

- Consume `validated-transactions`
- Simulate per-bank latency, jitter, failure rates
- Overload detection → `BANK_OVERLOAD`
- Circuit breaker per bank
- Emit `latency-events`, `bank-health`, `failed-transactions`, `analytics-events`

## Kafka

- **Consumes:** `validated-transactions`, `retry-transactions`
- **Produces:** `failed-transactions`, `latency-events`, `bank-health`, `analytics-events`

## Contract

`shared/contracts/openapi/bank-simulator.yaml`

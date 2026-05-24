# Intelligent Routing Service

**Port:** 8086 | **Phase 3**

## Responsibilities

- Aggregate `latency-events` and `bank-health` into congestion scores
- Consume AI forecasts from `congestion-events`
- Reroute traffic to alternate `bank_code`
- Publish `congestion-events` with `recommended_action`

## Kafka

- **Consumes:** `latency-events`, `bank-health`, `congestion-events`
- **Produces:** `congestion-events`, `analytics-events`

## Algorithms

- Weighted score: `0.4*p95_norm + 0.3*error_rate + 0.3*circuit_penalty`
- Reroute when score > 0.75 and alternate bank healthy

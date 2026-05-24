# Fraud / Anomaly Detection Service

**Port:** 8087 | **Phase 3**

## Responsibilities

- Consume `validated-transactions` (parallel path)
- Rule-based + statistical anomaly scoring (ML-ready hooks)
- Publish `fraud-alerts` for high-severity cases
- Block path: mark `FRAUD_SUSPECT` (non-retryable)

## Kafka

- **Consumes:** `validated-transactions`
- **Produces:** `fraud-alerts`, `analytics-events`

## Signals

- Velocity per payer VPA (sliding 5m window)
- Amount z-score vs merchant baseline
- Geo region mismatch

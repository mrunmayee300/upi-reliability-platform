# Transaction Ingestion Service

**Port:** 8081 | **Phase 3**

## Responsibilities

- Validate against `transaction.json` schema
- Idempotency enforcement (Redis + Postgres audit)
- Publish `upi-transactions` and `validated-transactions`
- Transaction status API

## Kafka

- **Produces:** `upi-transactions`, `validated-transactions`

## Contract

`shared/contracts/openapi/tx-ingestion.yaml`

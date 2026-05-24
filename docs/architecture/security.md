# Security (Simulation Layer)

## Authentication (simulated)

API Gateway validates `Authorization: Bearer <api_key>` against ConfigMap/Secret `API_KEYS` (comma-separated). Production would integrate OAuth2/mTLS with PSP partners.

## Rate limiting

- Per API key: 1000 req/s default (configurable).
- Per IP: 100 req/s for unauthenticated health endpoints.
- Implemented via Redis `INCR` + sliding window in Gateway middleware.

## Network (Kubernetes)

- `NetworkPolicy`: only Gateway ingress from Ingress controller; services talk via cluster DNS.
- Kafka: TLS in production Terraform module; PLAINTEXT local only.

## Secrets management

| Secret | Storage |
|--------|---------|
| Postgres credentials | K8s Secret / Terraform output |
| Redis password | K8s Secret |
| API keys | K8s Secret |
| Grafana admin | Helm values (override in prod) |

Never commit `.env` with real credentials. Use `.env.example` templates (Phase 2).

## PII handling (simulation)

- VPA addresses are synthetic (`user@paytm`, `user@ybl`).
- No real PAN/Aadhaar; `masked_account` is fake data.
- Logs redact `device_fingerprint` in production mode via structured log filter.

## Audit trail

All state transitions written to `transaction_audit` table with `correlation_id`, `actor_service`, `payload_hash`.

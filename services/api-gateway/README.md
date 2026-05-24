# API Gateway Service

**Port:** 8080 | **Phase 3**

## Responsibilities

- Request routing to Tx Ingestion and Analytics
- Bearer token auth simulation
- Redis-backed rate limiting (sliding window)
- Request validation and correlation ID injection
- Prometheus metrics: `gateway_requests_total`, `gateway_rate_limited_total`

## Dependencies

- Tx Ingestion (HTTP)
- Redis
- OpenTelemetry Collector

## Contract

`shared/contracts/openapi/api-gateway.yaml`

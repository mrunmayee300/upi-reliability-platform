# Phase 3 Summary — Backend Microservices

## Implemented services (7)

| Service | Port | Responsibility |
|---------|------|----------------|
| **tx-ingestion** | 8081 | Validate, idempotency (Redis), Postgres persist, Kafka publish |
| **bank-simulator** | 8083 | Consume `validated-transactions` + `retry-transactions`, simulate bank latency/failures |
| **failure-detector** | 8084 | Enrich `failed-transactions` with classified `failure_cause` |
| **retry-orchestrator** | 8085 | Exponential backoff retries, DLQ to Postgres + Kafka |
| **analytics** | 8089 | Aggregate metrics, REST summary, WebSocket live stream |
| **api-gateway** | 8080 | Auth, rate limit, proxy to ingestion/analytics |
| **tx-generator** | 8082 | Synthetic load generator → ingestion HTTP |

## Event flow (implemented)

```
Generator/Gateway → Ingestion → validated-transactions
       → Bank Simulator → success: analytics-events
                         → fail: failed-transactions (enriched=false)
       → Failure Detector → failed-transactions (enriched=true)
       → Retry Orchestrator → retry-transactions (after backoff)
       → Bank Simulator (retry path)
       → DLQ after max attempts
```

## Shared libraries (`shared/go`)

- `config` — environment configuration
- `kafkax` — Kafka producer/consumer wrappers
- `idempotency` — Redis SET NX store
- `store` — PostgreSQL transaction/retry/DLQ persistence
- `events` — envelope builder
- `httpx` — Gin router, health, correlation middleware
- `metrics` — Prometheus counters

## Run locally

### 1. Infra (Phase 2)

```powershell
.\scripts\dev\bootstrap.ps1
```

### 2. Build services (requires Go 1.22+)

Go is often installed at `C:\Program Files\Go\bin` but missing from PATH. Scripts auto-detect it.

```powershell
# From repo root (not shared/go):
cd C:\Users\Mrunmayee\OneDrive\Desktop\upi
.\scripts\dev\build-all.ps1
```

If Go is not installed:

```powershell
winget install GoLang.Go
# Close and reopen PowerShell, then:
go version
```

### 3. Start services

```powershell
.\scripts\dev\run-services.ps1
```

### 4. Generate load

```powershell
Invoke-RestMethod -Method POST http://localhost:8082/v1/generator/start `
  -ContentType "application/json" -Body '{"tps": 50}'
```

### 5. Observe

- Metrics: http://localhost:8089/v1/metrics/summary
- Gateway: http://localhost:8080/v1/metrics/summary (with `Authorization: Bearer dev-api-key-001`)
- Grafana: http://localhost:3001
- Kafka UI: http://localhost:8088

## Docker (optional)

```powershell
docker compose -f docker-compose.yml -f docker-compose.apps.yml up -d --build
```

> Note: containerized apps require Kafka dual-listener config for `kafka:9092` internal access (host dev uses `run-services.ps1`).

## Deferred to Phase 3b / Phase 4

- intelligent-routing
- fraud-detector
- notification
- ai-prediction (Python)
- dashboard (Next.js)

## Prerequisites

- Go **1.22+** installed and on `PATH`
- Docker Desktop running (for Kafka/Postgres/Redis)

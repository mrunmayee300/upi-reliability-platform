# Phase 4 Summary — Enterprise Dashboard

![Dashboard preview](../images/dashboard-preview.svg)

## Delivered

### Next.js operations UI (`apps/dashboard`)

- **Stack:** Next.js 14, TypeScript, TailwindCSS, Recharts, Lucide icons
- **Theme:** Dark enterprise layout with live status indicator
- **Real-time:** WebSocket to Analytics `ws://localhost:8089/v1/ws/live`
- **Panels:**
  - KPI cards (TPS, success/failure/retry rates, P95 latency)
  - Throughput & latency chart (30-sample window)
  - Bank health grid with congestion status
  - Live transaction feed
  - Retry & DLQ metrics

### Analytics service enhancements

- `GET /v1/dashboard` — full payload for UI
- WebSocket broadcasts `DashboardPayload` (summary + banks + retries + events)
- Consumes `bank-health` Kafka topic
- CORS enabled for browser access from `:3000`

## Run locally

```powershell
# 1. Infra + backend (from repo root)
.\scripts\dev\bootstrap.ps1
.\scripts\dev\build-all.ps1
.\scripts\dev\run-services.ps1

# 2. Rebuild analytics after Phase 4 changes
.\bin\analytics.exe   # or rebuild: .\scripts\dev\build-all.ps1

# 3. Dashboard
.\scripts\dev\run-dashboard.ps1
```

Open **http://localhost:3000**

Generate traffic:

```powershell
Invoke-RestMethod -Uri http://localhost:8082/v1/generator/start -Method POST -ContentType application/json -Body '{"tps":50}'
```

## Environment

Copy `apps/dashboard/.env.local.example` → `.env.local`:

```
NEXT_PUBLIC_ANALYTICS_URL=http://localhost:8089
NEXT_PUBLIC_WS_URL=ws://localhost:8089/v1/ws/live
```

## Phase 5 preview

Kubernetes manifests, Helm charts, Terraform modules for cloud deployment.

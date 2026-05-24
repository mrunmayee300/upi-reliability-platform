# Load Tests (k6)

## Prerequisites

```powershell
# Install k6: https://k6.io/docs/get-started/installation/
choco install k6
# or: winget install Grafana.k6
```

Backend must be running (`.\scripts\dev\run-services.ps1`).

## Run

```powershell
cd load-tests
.\run-benchmark.ps1 -Test baseline
.\run-benchmark.ps1 -Test surge
.\run-benchmark.ps1 -Test 10k-tps
```

Results: `load-tests/results/`

## Env

| Var | Default |
|-----|---------|
| INGESTION_URL | http://localhost:8081 |
| GATEWAY_URL | http://localhost:8080 |
| API_KEY | dev-api-key-001 |

# Phase 6 — Load Testing & Chaos Engineering

## k6 tests (`load-tests/k6/`)

| Script | Scenario |
|--------|----------|
| `baseline.js` | 50 VUs, 2m steady |
| `surge-traffic.js` | Spike to 800 VUs |
| `10k-tps.js` | Constant arrival rate 10K/s |
| `gateway-load.js` | Via API gateway + auth |
| `bank-failure.js` | High volume bank path |
| `retry-storm.js` | Retry pressure |
| `congestion-spike.js` | Gateway spike |

Run: `.\load-tests\run-benchmark.ps1 -Test all`

## Chaos (`scripts/chaos/`)

- `kafka-outage.ps1` — broker stop/start
- `kill-service.ps1` — kill one microservice

## Results

Exported to `load-tests/results/*.json` and `*-summary.txt`.

# Local Infra Runbook (Phase 2)

## Prerequisites

- Docker Desktop (Compose v2)
- Bash (Git Bash / WSL) for `.sh` scripts
- PowerShell for Windows bootstrap convenience

## Start infrastructure

```powershell
docker compose up -d
```

Or bootstrap with topic creation:

```powershell
.\scripts\dev\bootstrap.ps1
```

## Infra endpoints

- Kafka broker: `localhost:9092`
- Kafka UI: [http://localhost:8088](http://localhost:8088)
- PostgreSQL: `localhost:5433` (`upi` / `upi_dev_password`) — port **5433** avoids conflict with a local PostgreSQL on 5432
- Redis: `localhost:6379`
- Prometheus: [http://localhost:9090](http://localhost:9090)
- Grafana: [http://localhost:3001](http://localhost:3001) (`admin/admin`)
- Jaeger: [http://localhost:16686](http://localhost:16686)
- OTel collector metrics: [http://localhost:8889/metrics](http://localhost:8889/metrics)

## Verify stack

```powershell
docker ps
docker compose logs kafka --tail 50
docker compose logs prometheus --tail 50
```

Topic checks:

```bash
docker exec upi-kafka kafka-topics.sh --bootstrap-server localhost:9092 --list
```

## Stop stack

```powershell
docker compose down
```

Delete volumes:

```powershell
docker compose down -v
```

## Common issues

### `dockerDesktopLinuxEngine: The system cannot find the file specified`

**Cause:** Docker Desktop is not running (or the Linux engine is not started).

**Fix:**

1. Open **Docker Desktop** from the Start menu.
2. Wait until the tray icon shows **Engine running**.
3. Re-run:

```powershell
.\scripts\dev\bootstrap.ps1
```

Verify Docker works:

```powershell
docker info
```

### `bitnami/kafka:3.7: not found`

**Cause:** Invalid or removed Docker image tag.

**Fix:** This repo uses `apache/kafka:3.8.1`. Pull latest compose and re-run bootstrap:

```powershell
docker compose pull kafka
.\scripts\dev\bootstrap.ps1
```

### `go : The term 'go' is not recognized`

**Cause:** Go is not on your PowerShell `PATH` (common on Windows after install).

**Fix:**

1. Install Go if needed: `winget install GoLang.Go`
2. **Close and reopen** PowerShell (installer updates PATH).
3. Or run scripts from repo root — they auto-add `C:\Program Files\Go\bin`:

```powershell
cd C:\Users\Mrunmayee\OneDrive\Desktop\upi
.\scripts\dev\build-all.ps1
```

Manual PATH for current session:

```powershell
$env:PATH = "C:\Program Files\Go\bin;$env:PATH"
go version
```

### Other issues

- `port already allocated`: free the port or map a different host port in `docker-compose.yml`.
- Kafka not healthy: wait for KRaft metadata init, then run bootstrap again.
- Grafana empty: ensure `dashboards/grafana/*.json` exists and `dashboard.yaml` path matches.
- No traces in Jaeger: verify services export OTLP to `http://localhost:4318`.

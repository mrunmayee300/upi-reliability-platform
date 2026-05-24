# UPI Platform — Setup & Run

<p align="center">
  <img src="docs/images/banner.svg" alt="UPI Platform" width="100%"/>
</p>

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Docker Desktop | latest | https://www.docker.com/products/docker-desktop/ |
| Go | 1.22+ | https://go.dev/dl/ |
| Node.js | 20+ | https://nodejs.org/ |
| Python | 3.11+ | https://www.python.org/ (AI service, optional) |
| k6 | latest | `winget install Grafana.k6` (load tests, optional) |
| kind + Helm | optional | K8s deploy only |

After installing Go, **restart PowerShell**. If `go` not found:

```powershell
.\scripts\dev\refresh-path.ps1
```

---

## 1. Local development (recommended)

```powershell
cd C:\Users\Mrunmayee\OneDrive\Desktop\upi

# Infra: Kafka, Postgres :5433, Redis, Grafana, Jaeger
.\scripts\dev\bootstrap.ps1

# Build all Go services
.\scripts\dev\build-all.ps1

# Start all services (background, logs in .\logs\)
.\scripts\dev\run-services.ps1

# Dashboard (new terminal)
.\scripts\dev\run-dashboard.ps1

# Generate traffic
.\scripts\dev\start-generator.ps1
```

### URLs

| Service | URL |
|---------|-----|
| Dashboard | http://localhost:3000 |
| API Gateway | http://localhost:8080 |
| Analytics / WS | http://localhost:8089 |
| Grafana | http://localhost:3001 (admin/admin) |
| Kafka UI | http://localhost:8088 |
| Jaeger | http://localhost:16686 |
| Prometheus | http://localhost:9090 |

### Stop

```powershell
.\scripts\dev\stop-services.ps1
docker compose down
```

---

## 2. Load tests (k6)

```powershell
# Services must be running
.\load-tests\run-benchmark.ps1 -Test baseline
.\load-tests\run-benchmark.ps1 -Test surge
.\load-tests\run-benchmark.ps1 -Test all
```

Results: `load-tests/results/`

---

## 3. Chaos tests

```powershell
.\scripts\chaos\kafka-outage.ps1
.\scripts\chaos\kill-service.ps1 -ServiceName bank-simulator
```

---

## 4. Kubernetes (kind)

```powershell
.\scripts\k8s\create-kind.ps1
.\scripts\k8s\build-images.ps1
.\scripts\k8s\deploy-helm.ps1

kubectl port-forward -n upi-platform svc/dashboard 3000:3000
kubectl port-forward -n upi-platform svc/api-gateway 8080:8080
```

### Terraform

```powershell
cd terraform\environments\local
copy terraform.tfvars.example terraform.tfvars
terraform init
terraform apply
```

---

## 5. Optional services

```powershell
# AI prediction (Python) — port 8091
.\scripts\dev\run-ai-prediction.ps1
```

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| `go` not found | `.\scripts\dev\refresh-path.ps1` or restart terminal |
| Postgres auth failed | Use port **5433** not 5432; `docker compose up -d postgres` |
| Port 8082 connection refused | `.\scripts\dev\run-services.ps1` |
| Grafana empty | Start generator traffic first |
| Kafka topics missing | `.\scripts\kafka\create-topics.ps1` |

Logs: `Get-Content .\logs\<service>.err.log -Tail 30`

---

## Deployment

Full guide: [docs/runbooks/deployment.md](docs/runbooks/deployment.md)

**Docker Compose (easiest):**

```powershell
.\scripts\deploy\docker-compose-deploy.ps1
```

**Kubernetes (kind + Helm):**

```powershell
.\scripts\k8s\install-prereqs.ps1   # once, restart shell
.\scripts\k8s\deploy-local.ps1
```

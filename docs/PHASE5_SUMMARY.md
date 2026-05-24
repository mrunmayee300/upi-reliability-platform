# Phase 5 Summary — Kubernetes, Helm & Terraform

## Delivered

### Helm chart (`helm/upi-platform/`)

| Component | Resources |
|-----------|-----------|
| **Infra** | Kafka (StatefulSet), Postgres, Redis |
| **Apps** | 7 Go microservices + Next.js dashboard |
| **Networking** | Services, Ingress (nginx) |
| **Scaling** | HPA on ingestion, bank, gateway, analytics |
| **Observability** | OTel Collector, Jaeger |
| **Jobs** | Kafka topic bootstrap |

**Values files:**
- `values.yaml` — production-like defaults (replicas, HPA)
- `values-local.yaml` — kind/local overrides

### Terraform

- Module `terraform/modules/helm-upi` — wraps `helm_release`
- Environment `terraform/environments/local` — deploys to kind context

### Scripts

| Script | Purpose |
|--------|---------|
| `scripts/k8s/create-kind.ps1` | Create `kind-upi-local` cluster |
| `scripts/k8s/build-images.ps1` | Docker build + `kind load` |
| `scripts/k8s/deploy-helm.ps1` | `helm upgrade --install` |

### CI (`.github/workflows/ci.yml`)

- Go service compile check
- Helm lint
- Dashboard Next.js build

### Docker

- `infra/docker/Dockerfile.service` — Go microservices (existing)
- `infra/docker/Dockerfile.dashboard` — Next.js standalone

## Architecture decisions

| Decision | Rationale |
|----------|-----------|
| Helm as deploy unit | Parameterized replicas/HPA/env per environment |
| In-cluster Kafka (local) | Parity with docker-compose; prod swaps to MSK |
| Terraform → Helm only | Avoid duplicating K8s resources in TF and Helm |
| kind for local K8s | Lightweight recruiter/demo cluster on laptop |

## Deploy locally

```powershell
# Prerequisites: Docker, kind, kubectl, helm

.\scripts\k8s\create-kind.ps1
.\scripts\dev\build-all.ps1          # optional: verify Go builds
.\scripts\k8s\build-images.ps1
.\scripts\k8s\deploy-helm.ps1

kubectl port-forward -n upi-platform svc/dashboard 3000:3000
kubectl port-forward -n upi-platform svc/api-gateway 8080:8080
```

## HPA targets (default `values.yaml`)

| Service | Min | Max | CPU target |
|---------|-----|-----|--------------|
| tx-ingestion | 3 | 20 | 70% |
| bank-simulator | 2 | 15 | 75% |
| api-gateway | 2 | 8 | 70% |
| analytics | 2 | 10 | 70% |

## Phase 6 preview

- k6 load tests (10K TPS)
- Chaos Mesh experiments
- Benchmark reports

# Deployment Runbook

![Deployment topology](../images/deployment-topology.svg)

## Your machine (checked)

| Tool | Status |
|------|--------|
| Docker | OK |
| kubectl | OK (v1.32) |
| kind | Install with `.\scripts\k8s\install-prereqs.ps1` |
| Helm | Install with `.\scripts\k8s\install-prereqs.ps1` |

---

## Option A — Docker Compose (fastest, no K8s)

Everything runs in Docker. Good for demos and staging-like local deploy.

```powershell
cd C:\Users\Mrunmayee\OneDrive\Desktop\upi
.\scripts\deploy\docker-compose-deploy.ps1
```

**URLs:** Gateway `http://localhost:8080`, Grafana `http://localhost:3001`, Kafka UI `http://localhost:8088`

**Dashboard** (still runs on host for now):

```powershell
.\scripts\dev\run-dashboard.ps1
```

**Stop:**

```powershell
docker compose -f docker-compose.yml -f docker-compose.apps.yml down
```

---

## Option B — Kubernetes on kind (full Helm chart)

### 1. Install tools (once)

```powershell
.\scripts\k8s\install-prereqs.ps1
```

Restart PowerShell, then:

```powershell
kind version
helm version
```

### 2. Deploy

```powershell
cd C:\Users\Mrunmayee\OneDrive\Desktop\upi
.\scripts\k8s\deploy-local.ps1
```

This creates cluster `upi-local`, installs ingress, builds/loads images, runs Helm.

### 3. Access

```powershell
kubectl port-forward -n upi-platform svc/dashboard 3000:3000
kubectl port-forward -n upi-platform svc/api-gateway 8080:8080
```

**Verify:**

```powershell
kubectl get pods -n upi-platform
kubectl logs -n upi-platform deploy/tx-ingestion --tail=30
```

**Redeploy after code changes:**

```powershell
.\scripts\k8s\deploy-local.ps1 -SkipCluster
```

---

## Option C — Terraform (wraps Helm)

```powershell
cd terraform\environments\local
copy terraform.tfvars.example terraform.tfvars
terraform init
terraform apply
```

Requires kind cluster + images already built (`deploy-local.ps1 -SkipCluster` after images).

---

## Option D — Cloud (AKS / EKS / GKE)

1. **Registry** — push images from `build-images.ps1` to ACR/ECR/GCR; set `global.imageRegistry` in Helm values.
2. **Managed Kafka** — MSK / Confluent / Aiven; set `global.kafka.brokers`.
3. **Managed Postgres** — RDS / Cloud SQL; set `global.postgres.*` and use K8s secrets.
4. **Secrets** — replace `values.yaml` passwords with `kubectl create secret` or External Secrets.
5. **Ingress** — cloud LB + cert-manager TLS; disable bundled Kafka/Postgres in values if using managed services:

```yaml
infra:
  kafka:
    enabled: false
  postgres:
    enabled: false
```

6. **Deploy:**

```powershell
helm upgrade --install upi-platform .\helm\upi-platform `
  -f .\helm\upi-platform\values.yaml `
  -f .\helm\upi-platform\values-prod.yaml `
  --namespace upi-platform --create-namespace
```

---

## Pre-production checklist

- [ ] Change `global.apiKeys` and DB passwords (not `dev-api-key-001`)
- [ ] Kafka topics created (Helm job or `create-topics.ps1`)
- [ ] HPA enabled for gateway/ingestion in prod values
- [ ] Prometheus scraping all `/metrics` endpoints
- [ ] Backups for Postgres
- [ ] Network policies / private subnets for data stores

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| `ImagePullBackOff` on kind | Re-run `build-images.ps1` (uses `imagePullPolicy: Never` locally) |
| Pods `CrashLoopBackOff` | `kubectl logs -n upi-platform <pod>` — usually Kafka not ready; wait for `upi-kafka` |
| Ingress 404 | Install ingress: `install-ingress.ps1`; port-forward ingress controller :80 |
| Helm timeout | `kubectl get events -n upi-platform`; increase `--timeout 15m` |
| Compose Postgres port | Host uses **5433** in dev scripts; inside compose network use `postgres:5432` |

# Kubernetes Deployment

Production deployments use **Helm** as the primary packaging format.

## Quick start (kind + Helm)

```powershell
# 1. Create cluster
.\scripts\k8s\create-kind.ps1

# 2. Build & load images
.\scripts\k8s\build-images.ps1

# 3. Deploy chart
.\scripts\k8s\deploy-helm.ps1
```

## Helm chart

- Chart: [`../helm/upi-platform/`](../helm/upi-platform/)
- Values: `values.yaml` (prod-like), `values-local.yaml` (kind/local)

### Includes

| Resource | Description |
|----------|-------------|
| Namespace | `upi-platform` |
| Kafka, Postgres, Redis | Infra StatefulSet/Deployment |
| 7 Go microservices | Deployments + Services |
| Dashboard | Next.js Deployment |
| Ingress | nginx routes |
| HPA | tx-ingestion, bank-simulator, api-gateway, analytics |
| OTel + Jaeger | Observability stack |
| Job | Kafka topic bootstrap |

## Terraform

```powershell
cd terraform/environments/local
terraform init
terraform apply
```

Requires existing kubeconfig (e.g. `kind-upi-local` context).

## Raw manifests

Helm is canonical; `k8s/base/` is reserved for kubectl-kustomize overlays if needed later.

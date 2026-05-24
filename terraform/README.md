# Terraform

Infrastructure-as-code for deploying the UPI platform Helm chart.

## Layout

```
terraform/
├── modules/helm-upi/          # Helm release module
└── environments/
    └── local/                 # kind / local kubeconfig
```

## Local (kind)

```powershell
.\scripts\k8s\create-kind.ps1
.\scripts\k8s\build-images.ps1

cd terraform/environments/local
copy terraform.tfvars.example terraform.tfvars
terraform init
terraform apply
```

## Production (stub)

For EKS/GKE/AKS:

1. Provision cluster (Terraform cloud module or managed K8s)
2. Install ingress-nginx + cert-manager
3. Point `chart_path` to `helm/upi-platform`
4. Use `values.yaml` with managed Kafka (MSK/Confluent) and RDS/ElastiCache endpoints

See `docs/PHASE5_SUMMARY.md` for scaling guidance.

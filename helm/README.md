# Helm Charts

## upi-platform

Primary umbrella chart deploying the full stack.

```powershell
helm upgrade --install upi-platform ./helm/upi-platform `
  -f ./helm/upi-platform/values.yaml `
  -f ./helm/upi-platform/values-local.yaml `
  --namespace upi-platform --create-namespace
```

Or use `.\scripts\k8s\deploy-helm.ps1`.

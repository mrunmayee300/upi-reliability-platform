# Deploy UPI platform via Helm
$ErrorActionPreference = "Stop"
$root = Resolve-Path "$PSScriptRoot\..\.."
$chart = Join-Path $root "helm\upi-platform"

if (-not (Get-Command helm -ErrorAction SilentlyContinue)) {
    Write-Host "Install Helm 3: https://helm.sh/docs/intro/install/" -ForegroundColor Red
    exit 1
}

helm upgrade --install upi-platform $chart `
    -f "$chart\values.yaml" `
    -f "$chart\values-local.yaml" `
    --namespace upi-platform `
    --create-namespace `
    --wait `
    --timeout 10m

Write-Host ""
Write-Host "Deployment complete." -ForegroundColor Green
kubectl get pods -n upi-platform
Write-Host ""
Write-Host "Port-forward dashboard: kubectl port-forward -n upi-platform svc/dashboard 3000:3000"
Write-Host "Port-forward gateway:   kubectl port-forward -n upi-platform svc/api-gateway 8080:8080"

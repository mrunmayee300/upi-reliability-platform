# NGINX Ingress for kind (required for Helm ingress resource)
$ErrorActionPreference = "Stop"

$manifest = "https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.11.3/deploy/static/provider/kind/deploy.yaml"
Write-Host "Applying ingress-nginx for kind..." -ForegroundColor Cyan
kubectl apply -f $manifest

Write-Host "Waiting for ingress-nginx controller..." -ForegroundColor Cyan
kubectl wait --namespace ingress-nginx `
    --for=condition=ready pod `
    --selector=app.kubernetes.io/component=controller `
    --timeout=180s

Write-Host "Ingress controller ready." -ForegroundColor Green

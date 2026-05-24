# Full local K8s deploy: kind cluster -> ingress -> images -> Helm
param(
    [switch]$SkipCluster,
    [switch]$SkipBuild
)

$ErrorActionPreference = "Stop"
$root = Resolve-Path "$PSScriptRoot\..\.."
Set-Location $root

foreach ($cmd in @("docker", "kubectl", "kind", "helm")) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        Write-Host "Missing '$cmd'. Run .\scripts\k8s\install-prereqs.ps1 (kind + helm) and ensure Docker Desktop is running." -ForegroundColor Red
        exit 1
    }
}

if (-not $SkipCluster) {
    & "$PSScriptRoot\create-kind.ps1"
    & "$PSScriptRoot\install-ingress.ps1"
}

if (-not $SkipBuild) {
    & "$PSScriptRoot\build-images.ps1"
}

& "$PSScriptRoot\deploy-helm.ps1"

Write-Host ""
Write-Host "=== Access (pick one) ===" -ForegroundColor Green
Write-Host "Port-forward (simplest):"
Write-Host "  kubectl port-forward -n upi-platform svc/dashboard 3000:3000"
Write-Host "  kubectl port-forward -n upi-platform svc/api-gateway 8080:8080"
Write-Host ""
Write-Host "Ingress (add to hosts: 127.0.0.1 localhost):"
Write-Host "  kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 80:80"
Write-Host "  Dashboard: http://localhost/"
Write-Host "  API:       http://localhost/api"
Write-Host ""
Write-Host "Status: kubectl get pods -n upi-platform"

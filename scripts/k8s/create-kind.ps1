# Create a local kind cluster for UPI platform
$ErrorActionPreference = "Stop"

$clusterName = "upi-local"

if (-not (Get-Command kind -ErrorAction SilentlyContinue)) {
    Write-Host "Install kind: https://kind.sigs.k8s.io/docs/user/quick-start/" -ForegroundColor Red
    exit 1
}

$existing = kind get clusters 2>$null | Select-String $clusterName
if ($existing) {
    Write-Host "Cluster '$clusterName' already exists." -ForegroundColor Yellow
} else {
    kind create cluster --name $clusterName
}

Write-Host "Cluster ready. Context: kind-$clusterName" -ForegroundColor Green
kubectl cluster-info --context "kind-$clusterName"

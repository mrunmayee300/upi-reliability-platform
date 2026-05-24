# Build and load Docker images into kind cluster
$ErrorActionPreference = "Stop"
$root = Resolve-Path "$PSScriptRoot\..\.."
$clusterName = "upi-local"

$services = @(
    "api-gateway", "tx-ingestion", "tx-generator", "bank-simulator",
    "failure-detector", "retry-orchestrator", "analytics",
    "fraud-detector", "notification", "intelligent-routing"
)

foreach ($svc in $services) {
    $path = "services/$svc"
    Write-Host "Building upi/${svc}:latest ..." -ForegroundColor Cyan
    docker build -f "$root/infra/docker/Dockerfile.service" `
        --build-arg "SERVICE_PATH=$path" `
        -t "upi/${svc}:latest" $root
}

Write-Host "Building upi/dashboard:latest ..." -ForegroundColor Cyan
docker build -f "$root/infra/docker/Dockerfile.dashboard" -t upi/dashboard:latest $root

if (Get-Command kind -ErrorAction SilentlyContinue) {
    foreach ($svc in $services) {
        kind load docker-image "upi/${svc}:latest" --name $clusterName
    }
    kind load docker-image upi/dashboard:latest --name $clusterName
    Write-Host "Images loaded into kind cluster '$clusterName'" -ForegroundColor Green
} else {
    Write-Host "kind not found - images built locally only" -ForegroundColor Yellow
}

# Deploy full stack with Docker Compose (infra + app services)
$ErrorActionPreference = "Stop"
$root = Resolve-Path "$PSScriptRoot\..\.."
Set-Location $root

Write-Host "Starting infra + apps (build may take several minutes)..." -ForegroundColor Cyan
docker compose -f docker-compose.yml -f docker-compose.apps.yml up -d --build

Write-Host "Waiting for Kafka..." -ForegroundColor Cyan
$deadline = (Get-Date).AddMinutes(3)
do {
    Start-Sleep -Seconds 5
    $ok = docker exec upi-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list 2>$null
} while (-not $ok -and (Get-Date) -lt $deadline)

if (-not $ok) {
    Write-Host "Kafka not ready — check: docker compose logs kafka" -ForegroundColor Red
    exit 1
}

& "$root\scripts\kafka\create-topics.ps1"

Write-Host ""
Write-Host "Deployed." -ForegroundColor Green
Write-Host "  API Gateway:  http://localhost:8080"
Write-Host "  Dashboard:    run .\scripts\dev\run-dashboard.ps1  (or build dashboard image separately)"
Write-Host "  Grafana:      http://localhost:3001  (admin/admin)"
Write-Host "  Kafka UI:     http://localhost:8088"
Write-Host ""
Write-Host "Start traffic: curl -X POST http://localhost:8082/v1/generator/start -H Content-Type: application/json -d '{\"rate_per_second\":50}'"
Write-Host "Logs:          docker compose -f docker-compose.yml -f docker-compose.apps.yml logs -f tx-ingestion"

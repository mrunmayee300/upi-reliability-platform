# Bootstrap local infra (Kafka, Postgres, Redis, observability)
$ErrorActionPreference = "Stop"

function Test-DockerRunning {
    $prev = $ErrorActionPreference
    $ErrorActionPreference = "SilentlyContinue"
    docker info *> $null
    $ok = ($LASTEXITCODE -eq 0)
    $ErrorActionPreference = $prev
    return $ok
}

function Wait-KafkaHealthy {
    param([int]$MaxAttempts = 60)

    Write-Host "Waiting for Kafka health..." -ForegroundColor Cyan
    $prev = $ErrorActionPreference
    $ErrorActionPreference = "SilentlyContinue"
    for ($i = 0; $i -lt $MaxAttempts; $i++) {
        $health = docker inspect --format '{{.State.Health.Status}}' upi-kafka 2>$null
        if ($health -eq "healthy") {
            Write-Host "Kafka is healthy." -ForegroundColor Green
            return $true
        }
        Start-Sleep -Seconds 2
    }
    $ErrorActionPreference = $prev
    return $false
}

if (-not (Test-DockerRunning)) {
    Write-Host ""
    Write-Host "ERROR: Docker is not running." -ForegroundColor Red
    Write-Host ""
    Write-Host "The bootstrap failed because Docker Desktop is not started."
    Write-Host "Fix:"
    Write-Host "  1. Open Docker Desktop"
    Write-Host "  2. Wait until status shows 'Engine running'"
    Write-Host "  3. Run: .\scripts\dev\bootstrap.ps1"
    Write-Host ""
    Write-Host "If Docker Desktop is installed but still fails, restart Docker Desktop"
    Write-Host "or enable WSL2 backend in Docker Desktop settings."
    Write-Host ""
    exit 1
}

Write-Host "Starting infra stack..." -ForegroundColor Cyan
docker compose up -d kafka postgres redis jaeger otel-collector prometheus grafana kafka-ui
if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "ERROR: docker compose up failed." -ForegroundColor Red
    Write-Host "Check logs: docker compose logs"
    Write-Host ""
    exit 1
}

if (-not (Wait-KafkaHealthy)) {
    Write-Host ""
    Write-Host "ERROR: Kafka did not become healthy in time." -ForegroundColor Red
    Write-Host "Check logs: docker compose logs kafka"
    Write-Host ""
    exit 1
}

Write-Host "Creating Kafka topics..." -ForegroundColor Cyan
& "$PSScriptRoot\..\kafka\create-topics.ps1"
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}

Write-Host ""
Write-Host "Bootstrap complete." -ForegroundColor Green
Write-Host "Kafka UI:   http://localhost:8088"
Write-Host "Prometheus: http://localhost:9090"
Write-Host "Grafana:    http://localhost:3001 (admin/admin)"
Write-Host "Jaeger:     http://localhost:16686"

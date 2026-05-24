# Start synthetic UPI traffic (checks tx-generator is up first)
param(
    [int]$Tps = 50,
    [int]$DurationSeconds = 0
)

$ErrorActionPreference = "Stop"
$uri = "http://localhost:8082/v1/generator/start"

try {
    $null = Invoke-WebRequest -Uri "http://localhost:8082/health/live" -UseBasicParsing -TimeoutSec 3
} catch {
    Write-Host ""
    Write-Host "ERROR: tx-generator is not running on port 8082." -ForegroundColor Red
    Write-Host ""
    Write-Host "Start backend services first:"
    Write-Host "  cd `"$((Resolve-Path "$PSScriptRoot\..\..").Path)`""
    Write-Host "  .\scripts\dev\run-services.ps1"
    Write-Host ""
    exit 1
}

$body = @{ tps = $Tps }
if ($DurationSeconds -gt 0) { $body.duration_seconds = $DurationSeconds }

$response = Invoke-RestMethod -Uri $uri -Method POST -ContentType "application/json" -Body ($body | ConvertTo-Json)
Write-Host "Generator started:" -ForegroundColor Green
$response | ConvertTo-Json

Write-Host ""
Write-Host "Status:  Invoke-RestMethod http://localhost:8082/v1/generator/status"
Write-Host "Metrics: Invoke-RestMethod http://localhost:8089/v1/metrics/summary"
Write-Host "Dashboard: http://localhost:3000"

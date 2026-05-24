# Simulate Kafka outage — stop broker, wait, restart
param([int]$DurationSeconds = 30)
$ErrorActionPreference = "Stop"
Write-Host "Stopping Kafka..." -ForegroundColor Yellow
docker stop upi-kafka
Start-Sleep -Seconds $DurationSeconds
Write-Host "Starting Kafka..." -ForegroundColor Green
docker start upi-kafka
Write-Host "Wait for health, then run: .\scripts\dev\bootstrap.ps1 (topics)" -ForegroundColor Cyan

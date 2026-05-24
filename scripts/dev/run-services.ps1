$ErrorActionPreference = "Stop"

$root = (Resolve-Path "$PSScriptRoot\..\..").Path
$binDir = Join-Path $root "bin"
$logDir = Join-Path $root "logs"
$pidDir = Join-Path $logDir "pids"
New-Item -ItemType Directory -Force -Path $logDir, $pidDir | Out-Null

$env:KAFKA_BROKERS = "localhost:9092"
$env:POSTGRES_DSN = "postgres://upi:upi_dev_password@localhost:5433/upi_platform?sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
$env:API_KEYS = "dev-api-key-001"
$env:INGESTION_URL = "http://localhost:8081"
$env:ANALYTICS_URL = "http://localhost:8089"

$services = @(
    @{ Name = "tx-ingestion";       Port = "8081" },
    @{ Name = "bank-simulator";     Port = "8083" },
    @{ Name = "failure-detector";   Port = "8084" },
    @{ Name = "retry-orchestrator"; Port = "8085" },
    @{ Name = "intelligent-routing"; Port = "8086" },
    @{ Name = "fraud-detector";     Port = "8087" },
    @{ Name = "analytics";          Port = "8089" },
    @{ Name = "notification";       Port = "8090" },
    @{ Name = "api-gateway";        Port = "8080" },
    @{ Name = "tx-generator";       Port = "8082" }
)

foreach ($svc in $services) {
    if (-not (Test-Path (Join-Path $binDir "$($svc.Name).exe"))) {
        Write-Host "Missing bin\$($svc.Name).exe - run .\scripts\dev\build-all.ps1" -ForegroundColor Red
        exit 1
    }
}

& "$PSScriptRoot\stop-services.ps1" | Out-Null
Start-Sleep -Seconds 1

Write-Host "Starting services..." -ForegroundColor Cyan
foreach ($svc in $services) {
    $exe = Join-Path $binDir "$($svc.Name).exe"
    $env:HTTP_PORT = $svc.Port
    $proc = Start-Process -FilePath $exe -WorkingDirectory $root -WindowStyle Hidden `
        -RedirectStandardOutput (Join-Path $logDir "$($svc.Name).out.log") `
        -RedirectStandardError (Join-Path $logDir "$($svc.Name).err.log") -PassThru
    Set-Content (Join-Path $pidDir "$($svc.Name).pid") -Value $proc.Id
    Write-Host "  $($svc.Name) :$($svc.Port)" -ForegroundColor DarkGray
    Start-Sleep -Milliseconds 500
}

Start-Sleep -Seconds 8
& "$PSScriptRoot\verify-services.ps1"
Write-Host "Dashboard: .\scripts\dev\run-dashboard.ps1" -ForegroundColor Green
Write-Host "Traffic:   .\scripts\dev\start-generator.ps1" -ForegroundColor Green

$ErrorActionPreference = "Stop"
. "$PSScriptRoot\ensure-go.ps1"

$root = Resolve-Path "$PSScriptRoot\..\.."
Set-Location $root

$binDir = Join-Path $root "bin"
New-Item -ItemType Directory -Force -Path $binDir | Out-Null

Push-Location shared/go
go mod tidy
Pop-Location

$services = @(
    "services/tx-ingestion",
    "services/bank-simulator",
    "services/failure-detector",
    "services/retry-orchestrator",
    "services/analytics",
    "services/api-gateway",
    "services/tx-generator",
    "services/fraud-detector",
    "services/notification",
    "services/intelligent-routing"
)

foreach ($svc in $services) {
    $name = Split-Path $svc -Leaf
    $out = Join-Path $binDir "$name.exe"
    Write-Host "Building $svc ..." -ForegroundColor Cyan
    Push-Location $svc
    go mod tidy
    if ($LASTEXITCODE -ne 0) { Pop-Location; exit 1 }
    go build -o $out ./cmd/server
    if ($LASTEXITCODE -ne 0) { Pop-Location; exit 1 }
    Pop-Location
}

Write-Host "Built into .\bin\" -ForegroundColor Green

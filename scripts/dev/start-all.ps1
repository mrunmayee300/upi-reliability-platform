# One command: infra + build + services (dashboard separate)
$ErrorActionPreference = "Stop"
$root = Resolve-Path "$PSScriptRoot\..\.."
Set-Location $root
.\scripts\dev\bootstrap.ps1
.\scripts\dev\build-all.ps1
.\scripts\dev\run-services.ps1
Write-Host ""
Write-Host "Next: .\scripts\dev\run-dashboard.ps1" -ForegroundColor Green
Write-Host "      .\scripts\dev\start-generator.ps1" -ForegroundColor Green

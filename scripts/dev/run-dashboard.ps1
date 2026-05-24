$ErrorActionPreference = "Stop"
$root = Resolve-Path "$PSScriptRoot\..\.."
$dash = Join-Path $root "apps\dashboard"

if (-not (Test-Path (Join-Path $dash ".env.local"))) {
    Copy-Item (Join-Path $dash ".env.local.example") (Join-Path $dash ".env.local")
    Write-Host "Created apps/dashboard/.env.local" -ForegroundColor DarkGray
}

Set-Location $dash
if (-not (Test-Path "node_modules")) {
    Write-Host "Installing npm dependencies..." -ForegroundColor Cyan
    npm install
}

Write-Host "Starting dashboard at http://localhost:3000" -ForegroundColor Green
Write-Host "Requires analytics on http://localhost:8089" -ForegroundColor Yellow
npm run dev

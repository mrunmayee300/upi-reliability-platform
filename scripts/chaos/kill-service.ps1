param([Parameter(Mandatory=$true)][string]$ServiceName)
# ServiceName: tx-ingestion, bank-simulator, etc.
$ErrorActionPreference = "Stop"
$pidFile = "logs\pids\$ServiceName.pid"
if (Test-Path $pidFile) {
    $pid = Get-Content $pidFile
    Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
    Write-Host "Killed $ServiceName (PID $pid)" -ForegroundColor Yellow
} else {
    Write-Host "No PID file. Use: Get-NetTCPConnection -LocalPort <port>" -ForegroundColor Red
}

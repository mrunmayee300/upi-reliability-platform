$names = @(
    "tx-ingestion", "bank-simulator", "failure-detector", "retry-orchestrator",
    "intelligent-routing", "fraud-detector", "analytics", "notification",
    "api-gateway", "tx-generator"
)
$root = Resolve-Path "$PSScriptRoot\..\.."
$pidDir = Join-Path $root "logs\pids"
foreach ($n in $names) {
    $f = Join-Path $pidDir "$n.pid"
    if (Test-Path $f) {
        Stop-Process -Id (Get-Content $f) -Force -ErrorAction SilentlyContinue
        Remove-Item $f -Force
    }
}
@(8080,8081,8082,8083,8084,8085,8086,8087,8089,8090) | ForEach-Object {
    $c = Get-NetTCPConnection -LocalPort $_ -State Listen -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($c) { Stop-Process -Id $c.OwningProcess -Force -ErrorAction SilentlyContinue }
}
Write-Host "Stopped."

$checks = @(
    @{ Name = "api-gateway";          Port = 8080 },
    @{ Name = "tx-ingestion";         Port = 8081 },
    @{ Name = "tx-generator";         Port = 8082 },
    @{ Name = "bank-simulator";       Port = 8083 },
    @{ Name = "failure-detector";     Port = 8084 },
    @{ Name = "retry-orchestrator";   Port = 8085 },
    @{ Name = "intelligent-routing";  Port = 8086 },
    @{ Name = "fraud-detector";       Port = 8087 },
    @{ Name = "analytics";           Port = 8089 },
    @{ Name = "notification";        Port = 8090 }
)

Write-Host "Health check:" -ForegroundColor Cyan
$ok = 0
foreach ($c in $checks) {
    try {
        $r = Invoke-WebRequest -Uri "http://localhost:$($c.Port)/health/live" -UseBasicParsing -TimeoutSec 2
        if ($r.StatusCode -eq 200) { Write-Host "  OK   $($c.Name)" -ForegroundColor Green; $ok++ }
    } catch { Write-Host "  DOWN $($c.Name)" -ForegroundColor Red }
}
Write-Host "$ok / $($checks.Count) up"

param(
    [ValidateSet("baseline", "surge", "10k-tps", "gateway", "bank-failure", "retry-storm", "congestion", "all")]
    [string]$Test = "baseline"
)

$ErrorActionPreference = "Stop"
$root = $PSScriptRoot
$results = Join-Path $root "results"
New-Item -ItemType Directory -Force -Path $results | Out-Null

$tests = @{
    baseline  = "k6/baseline.js"
    surge     = "k6/surge-traffic.js"
    "10k-tps" = "k6/10k-tps.js"
    gateway   = "k6/gateway-load.js"
    "bank-failure" = "k6/bank-failure.js"
    "retry-storm"  = "k6/retry-storm.js"
    congestion     = "k6/congestion-spike.js"
}

function Run-K6($name, $script) {
    $ts = Get-Date -Format "yyyyMMdd-HHmmss"
    $out = Join-Path $results "$name-$ts.json"
    $sum = Join-Path $results "$name-$ts-summary.txt"
    Write-Host "Running $name ..." -ForegroundColor Cyan
    k6 run --summary-export $out $script 2>&1 | Tee-Object -FilePath $sum
}

Set-Location $root
if ($Test -eq "all") {
    foreach ($k in $tests.Keys) { Run-K6 $k $tests[$k] }
} else {
    Run-K6 $Test $tests[$Test]
}

Write-Host "Results in load-tests/results/" -ForegroundColor Green

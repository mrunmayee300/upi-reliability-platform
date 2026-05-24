# Install kind + Helm on Windows (winget). Requires admin for first install.
$ErrorActionPreference = "Stop"

$packages = @(
    @{ Id = "Kubernetes.kind"; Name = "kind" },
    @{ Id = "Helm.Helm"; Name = "helm" }
)

foreach ($p in $packages) {
    if (Get-Command $p.Name -ErrorAction SilentlyContinue) {
        Write-Host "$($p.Name) already installed" -ForegroundColor Green
        continue
    }
    Write-Host "Installing $($p.Name) via winget..." -ForegroundColor Cyan
    winget install --id $p.Id -e --accept-package-agreements --accept-source-agreements
}

Write-Host ""
Write-Host "Close and reopen PowerShell, then verify:" -ForegroundColor Yellow
Write-Host "  kind version"
Write-Host "  helm version"
Write-Host "  kubectl cluster-info"

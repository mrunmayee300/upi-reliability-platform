# Reload PATH from Windows registry (use after installing Go without restarting terminal)
$machine = [Environment]::GetEnvironmentVariable("Path", "Machine")
$user = [Environment]::GetEnvironmentVariable("Path", "User")
$env:Path = "$machine;$user"

if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "OK: $(go version)" -ForegroundColor Green
} else {
    Write-Host "Go still not found. Install with: winget install GoLang.Go" -ForegroundColor Red
    Write-Host "Then close ALL terminal windows and open a new one." -ForegroundColor Yellow
    exit 1
}

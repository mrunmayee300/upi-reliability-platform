# Ensures Go is available in the current PowerShell session.
# Dot-source from other scripts: . "$PSScriptRoot\ensure-go.ps1"

function Initialize-GoToolchain {
    if (Get-Command go -ErrorAction SilentlyContinue) {
        return (Get-Command go).Source
    }

    $candidates = @(
        "$env:ProgramFiles\Go\bin\go.exe",
        "${env:ProgramFiles(x86)}\Go\bin\go.exe",
        "$env:LOCALAPPDATA\Programs\Go\bin\go.exe"
    )

    foreach ($candidate in $candidates) {
        if (Test-Path $candidate) {
            $goBin = Split-Path $candidate -Parent
            $goRoot = Split-Path $goBin -Parent
            $env:GOROOT = $goRoot
            $env:PATH = "$goBin;$env:PATH"
            return $candidate
        }
    }

    return $null
}

$goExe = Initialize-GoToolchain
if (-not $goExe) {
    Write-Host ""
    Write-Host "ERROR: Go is not installed or not on PATH." -ForegroundColor Red
    Write-Host ""
    Write-Host "Install Go 1.22+ then reopen PowerShell:"
    Write-Host "  winget install GoLang.Go"
    Write-Host "  # or download: https://go.dev/dl/"
    Write-Host ""
    Write-Host "After install, verify:"
    Write-Host "  go version"
    Write-Host ""
    exit 1
}

if (-not $env:GOPATH) {
    $env:GOPATH = Join-Path $env:USERPROFILE "go"
}
if (-not $env:PATH.Contains("$env:GOPATH\bin")) {
    $env:PATH = "$env:GOPATH\bin;$env:PATH"
}

Write-Host "Using Go: $goExe ($(& $goExe version))" -ForegroundColor DarkGray

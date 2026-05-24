$ErrorActionPreference = "Stop"
$root = Resolve-Path "$PSScriptRoot\..\.."
$svc = Join-Path $root "services\ai-prediction"
if (-not (Test-Path (Join-Path $svc ".venv"))) {
    python -m venv (Join-Path $svc ".venv")
}
& (Join-Path $svc ".venv\Scripts\Activate.ps1")
pip install -q -r (Join-Path $svc "requirements.txt")
$env:HTTP_PORT = "8091"
python (Join-Path $svc "app\main.py")

# Run any go command with auto-detected toolchain (works when go is not on PATH)
. "$PSScriptRoot\ensure-go.ps1"
& go @args

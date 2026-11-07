# Loads PocketBase\.env into the current shell session.
# Usage: . .\loadenv.ps1; .\pb.exe serve
$envFile = Join-Path $PSScriptRoot ".env"
if (-not (Test-Path $envFile)) {
    Write-Error "Missing $envFile - copy .env.example to .env and fill in your IGDB credentials."
    exit 1
}
Get-Content $envFile | ForEach-Object {
    if ($_ -match '^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)$') {
        $value = $Matches[2].Trim().Trim('"').Trim("'")
        [Environment]::SetEnvironmentVariable($Matches[1], $value, "Process")
    }
}
Write-Host "Loaded environment from $envFile"

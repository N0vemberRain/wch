# env.ps1
# Sets up environment variables for local testing

Write-Host "Setting up local environment variables..."

# --- Clean up problematic global variables ---
if (Test-Path Env:PGLOCALEDIR) {
    Remove-Item Env:PGLOCALEDIR
    Write-Host "Removed PGLOCALEDIR to prevent lib/pq panic."
}

$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "683951"
$env:DB_NAME = "wch"
$env:DB_SSLMODE = "disable"

Write-Host "Environment variables loaded."
Write-Host "Run your service with: go run ./services/users/cmd/main.go"
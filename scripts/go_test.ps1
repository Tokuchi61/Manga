$ErrorActionPreference = 'Stop'
Push-Location apps/api
try {
    go test ./...
} finally {
    Pop-Location
}

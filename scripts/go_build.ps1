$ErrorActionPreference = 'Stop'
Push-Location apps/api
try {
    go build ./...
} finally {
    Pop-Location
}

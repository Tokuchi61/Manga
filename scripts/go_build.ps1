$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrWhiteSpace($env:APP_VERSION)) {
    $env:APP_VERSION = (Get-Content VERSION -Raw).Trim()
}

$cacheRoot = '.cache'
New-Item -ItemType Directory -Path "$cacheRoot/go-build", "$cacheRoot/go-mod", "$cacheRoot/gopath" -Force | Out-Null

if ([string]::IsNullOrWhiteSpace($env:GOCACHE)) {
    $env:GOCACHE = (Resolve-Path "$cacheRoot/go-build").Path
}
if ([string]::IsNullOrWhiteSpace($env:GOMODCACHE)) {
    $env:GOMODCACHE = (Resolve-Path "$cacheRoot/go-mod").Path
}
if ([string]::IsNullOrWhiteSpace($env:GOPATH)) {
    $env:GOPATH = (Resolve-Path "$cacheRoot/gopath").Path
}

Push-Location apps/api
try {
    go build ./...
} finally {
    Pop-Location
}

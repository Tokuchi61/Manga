$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrWhiteSpace($env:APP_VERSION)) {
    $env:APP_VERSION = (Get-Content VERSION -Raw).Trim()
}

docker compose -f deploy/docker-compose.yml up -d --build

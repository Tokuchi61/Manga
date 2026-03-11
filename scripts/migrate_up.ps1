param(
    [Parameter(Mandatory = $true)]
    [string]$DatabaseUrl
)

$ErrorActionPreference = 'Stop'
migrate -path apps/api/migrations -database $DatabaseUrl up

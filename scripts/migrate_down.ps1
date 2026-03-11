param(
    [Parameter(Mandatory = $true)]
    [string]$DatabaseUrl,
    [int]$Steps = 1
)

$ErrorActionPreference = 'Stop'
migrate -path apps/api/migrations -database $DatabaseUrl down $Steps

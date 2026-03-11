param(
    [Parameter(Mandatory = $true)]
    [string]$ModuleName,
    [string]$DomainGroup = ""
)

$ErrorActionPreference = 'Stop'

if ($ModuleName -notmatch '^[a-z][a-z0-9_]*$') {
    throw "ModuleName yalnizca kucuk harf, rakam ve underscore icerebilir; harf ile baslamalidir."
}

if ($DomainGroup -ne "" -and $DomainGroup -notmatch '^[a-z][a-z0-9_]*$') {
    throw "DomainGroup yalnizca kucuk harf, rakam ve underscore icerebilir; harf ile baslamalidir."
}

$moduleRoot = if ($DomainGroup -eq "") {
    Join-Path "apps/api/internal/modules" $ModuleName
} else {
    Join-Path (Join-Path "apps/api/internal/modules" $DomainGroup) $ModuleName
}

if (Test-Path $moduleRoot) {
    throw "Hedef modul yolu zaten mevcut: $moduleRoot"
}

$dirs = @(
    "entity",
    "dto",
    "service",
    "repository",
    "handler",
    "middleware",
    "validator",
    "mapper",
    "contract",
    "events",
    "consumer",
    "producer",
    "jobs",
    "readmodel"
)

New-Item -ItemType Directory -Path $moduleRoot -Force | Out-Null
foreach ($dir in $dirs) {
    New-Item -ItemType Directory -Path (Join-Path $moduleRoot $dir) -Force | Out-Null
}

$moduleGo = @"
package $ModuleName

import "github.com/go-chi/chi/v5"

// Module bu leaf modulu merkezi app bootstrap katmanina baglar.
type Module struct{}

func New() Module {
    return Module{}
}

func (m Module) Name() string {
    return "$ModuleName"
}

func (m Module) RegisterRoutes(router chi.Router) {
    registerRoutes(router)
}
"@

$routesGo = @"
package $ModuleName

import "github.com/go-chi/chi/v5"

func registerRoutes(router chi.Router) {
    // Stage 1 iskeleti: route kayitlari bu modulde acilir.
}
"@

$errorsGo = @"
package $ModuleName

// Stage 1 iskeleti: modul bazli hata tipleri burada tanimlanabilir.
"@

Set-Content -Path (Join-Path $moduleRoot "module.go") -Value $moduleGo -Encoding utf8
Set-Content -Path (Join-Path $moduleRoot "routes.go") -Value $routesGo -Encoding utf8
Set-Content -Path (Join-Path $moduleRoot "errors.go") -Value $errorsGo -Encoding utf8

Write-Output "Modul iskeleti olusturuldu: $moduleRoot"

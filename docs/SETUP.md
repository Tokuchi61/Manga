# Setup

Bu dokuman proje kurulumunu Asama 0-5 omurgasina gore tanimlar.

## Gereksinimler

- Go 1.26
- Docker ve Docker Compose
- (Opsiyonel) `migrate` CLI

## Versiyon Hazirligi

- Canonical versiyon dosyasi: `VERSION`
- Runtime icin `APP_VERSION` env degeri gerekir.
- PowerShell icin ornek:

```powershell
$env:APP_VERSION = (Get-Content VERSION -Raw).Trim()
```

## Ortam Degiskenleri

1. `cp .env.example .env`
2. Gerekirse DSN ve port degerlerini guncelle.
3. Auth guvenlik esikleri icin opsiyonel env degerleri:
   - `AUTH_LOGIN_FAILED_ATTEMPT_LIMIT_PER_MINUTE`
   - `AUTH_LOGIN_COOLDOWN_SECONDS`
   - `AUTH_EMAIL_VERIFICATION_RESEND_COOLDOWN_SECONDS`

## Docker ile Calistirma

1. `docker compose -f deploy/docker-compose.yml up -d --build`
2. `http://localhost:8080/health` endpointini kontrol et.
3. Kapatmak icin: `docker compose -f deploy/docker-compose.yml down`

## Local Calistirma

1. `cd apps/api`
2. `go mod tidy`
3. `go run ./cmd/api`

## Migration

- Up: `migrate -path apps/api/migrations -database "$DB_MAIN_DSN" up`
- Down: `migrate -path apps/api/migrations -database "$DB_MAIN_DSN" down 1`

## Modul Iskeleti Uretimi (Asama 1)

Yeni bir module klasor omurgasi acmak icin:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/scaffold_module.ps1 -ModuleName auth
powershell -ExecutionPolicy Bypass -File scripts/scaffold_module.ps1 -ModuleName chapter -DomainGroup content
```



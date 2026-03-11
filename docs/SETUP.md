# Setup

Bu dokuman proje kurulumunu Asama 0 omurgasina gore tanimlar.

## Gereksinimler

- Go 1.26
- Docker ve Docker Compose
- (Opsiyonel) `migrate` CLI

## Ortam Degiskenleri

1. `cp .env.example .env`
2. Gerekirse DSN ve port degerlerini guncelle.

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

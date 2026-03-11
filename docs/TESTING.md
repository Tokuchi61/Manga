# TESTING

Bu dokuman test katmanlarini ve Asama 0-1 dogrulama adimlarini listeler.

## Katmanlar

- Unit: config, servis ve saf is kurallari
- Integration: repository ve kritik HTTP akislarinin dogrulamasi
- Contract: moduller arasi acik kontratlarin dogrulamasi
- E2E: kritik uc-tan-uca akislar

## Asama 0 Temel Kontroller

- `cd apps/api && go test ./...`
- `cd apps/api && go build ./...`
- `docker build -f apps/api/Dockerfile -t novascans-api:local .`
- `docker compose -f deploy/docker-compose.yml up -d --build`
- `GET /health` ve `GET /version`

## Asama 1 Mimari Kontroller

- `cd apps/api && go test ./internal/modules/...`
- `cd apps/api && go test ./tests/contract/...`
- `cd apps/api && go test ./tests/integration/...`
- Module registry duplicate-name ve canonical-name kurallarinin test edilmesi
- App bootstrap katmaninda module route mount davranisinin test edilmesi

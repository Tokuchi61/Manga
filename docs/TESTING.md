# TESTING

Bu dokuman test katmanlarini ve Asama 0 dogrulama adimlarini listeler.

## Katmanlar

- Unit: config, servis ve saf is kurallari
- Integration: repository ve db akislarinin dogrulamasi
- Contract: moduller arasi acik kontratlarin dogrulamasi
- E2E: kritik HTTP akislari

## Asama 0 Temel Kontroller

- `cd apps/api && go test ./...`
- `cd apps/api && go build ./...`
- `docker build -f apps/api/Dockerfile -t novascans-api:local .`
- `docker compose -f deploy/docker-compose.yml up -d --build`
- `GET /health` ve `GET /version`

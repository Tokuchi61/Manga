# TESTING

Bu dokuman test katmanlarini ve Asama 0-7 dogrulama adimlarini listeler.

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

## Asama 4 Auth Kontrolleri

- `cd apps/api && go test ./internal/modules/auth/...`
- `cd apps/api && go test ./tests/contract -run Auth`
- `cd apps/api && go test ./tests/integration -run Auth`
- Register/login/logout/session revoke ve token rotation akislarinin dogrulanmasi
- Forgot/reset/change password ile verification/resend cooldown senaryolarinin dogrulanmasi
- Auth migration smoke kontrolu (`202603120002_auth_create_core_tables`)

## Asama 5 User Kontrolleri

- `cd apps/api && go test ./internal/modules/user/...`
- `cd apps/api && go test ./tests/contract -run User`
- `cd apps/api && go test ./tests/integration -run User`
- Profil okuma-guncelleme ve public/private response ayriminin dogrulanmasi
- History visibility preference ust sinir (global deny) davranisinin dogrulanmasi
- Account state (deactivated/banned) ve VIP lifecycle (activate/freeze/resume/deactivate) senaryolarinin dogrulanmasi
- User migration smoke kontrolu (`202603120003_user_create_core_tables`)

## Asama 6 Access Kontrolleri

- `cd apps/api && go test ./internal/modules/access/...`
- `cd apps/api && go test ./tests/contract -run Access`
- `cd apps/api && go test ./tests/integration -run Access`
- Role/permission tanimlari, role-permission baglama ve user-role assignment akislarinin dogrulanmasi
- Policy create/update davranislarinda conflict kontrolu ve precedence yorumunun dogrulanmasi
- Guest/authenticated/vip/blocked ayrimlari ile own/any decision matrix senaryolarinin dogrulanmasi
- Access migration smoke kontrolu (`202603120004_access_create_core_tables`)

## Asama 7 Manga Kontrolleri

- `cd apps/api && go test ./internal/modules/manga/...`
- `cd apps/api && go test ./tests/contract -run Manga`
- `cd apps/api && go test ./tests/integration -run Manga`
- Manga create/update/list/detail akislarinin ve search/filter/sort davranisinin dogrulanmasi
- Publish/archive/unpublish/schedule lifecycle akislari ile visibility daraltma davranisinin dogrulanmasi
- Discovery/editorial (featured/recommended/collection) ve counter sync (chapter/comment/view) dogrulamalarinin yapilmasi
- Soft delete/restore davranisi ile public detail/list gorunurlugunun dogrulanmasi
- Manga migration smoke kontrolu (`202603120005_manga_create_core_tables`)

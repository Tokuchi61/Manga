# Changelog

Bu dosya yalnizca projede gercekte yapilan islemleri kaydeder.
Bu proje SemVer (`MAJOR.MINOR.PATCH`) standardini takip eder.

## [Unreleased]

## [0.4.0-alpha.1] - 2026-03-12

### Added
- `Asama 4` kapsaminda canonical `auth` modulu eklendi: `apps/api/internal/modules/auth`.
- Auth module owner akis omurgasi eklendi:
  - register/login/logout
  - session list + revoke (`current`, `others`, `all`)
  - refresh token rotation
  - forgot/reset/change password
  - email verification send/confirm + resend cooldown
  - failed login limit + cooldown davranisi
- `auth -> access` kontrat yuzeyi eklendi: `apps/api/internal/modules/auth/contract/identity_contract.go`.
- Auth event sabitleri eklendi: `apps/api/internal/modules/auth/events/events.go`.
- In-memory auth repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/auth/repository/memory_store.go`
  - `apps/api/internal/modules/auth/repository/auth_credential_repository.go`
  - `apps/api/internal/modules/auth/repository/auth_session_repository.go`
  - `apps/api/internal/modules/auth/repository/auth_token_repository.go`
  - `apps/api/internal/modules/auth/repository/auth_security_event_repository.go`
  - `apps/api/internal/modules/auth/repository/memory_store_test.go`
- Auth service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/auth/service/auth_service.go`
  - `apps/api/internal/modules/auth/service/auth_register_service.go`
  - `apps/api/internal/modules/auth/service/auth_login_service.go`
  - `apps/api/internal/modules/auth/service/auth_refresh_token_service.go`
  - `apps/api/internal/modules/auth/service/auth_session_service.go`
  - `apps/api/internal/modules/auth/service/auth_password_service.go`
  - `apps/api/internal/modules/auth/service/auth_verification_service.go`
  - `apps/api/internal/modules/auth/service/auth_token_issue_service.go`
  - `apps/api/internal/modules/auth/service/auth_security_event_service.go`
  - `apps/api/internal/modules/auth/service/auth_request_meta_service.go`
  - `apps/api/internal/modules/auth/service/auth_session_parse_service.go`
  - `apps/api/internal/modules/auth/service/service_test.go`
- Auth HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/auth/handler/auth_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_register_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_login_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_token_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_session_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_password_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_verification_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_request_meta_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_response_handler.go`
- Auth route registration yapisi islem ailelerine gore parcalandi:
  - `apps/api/internal/modules/auth/routes.go`
  - `apps/api/internal/modules/auth/routes/register_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_identity_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_session_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_password_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_verification_routes.go`
- Auth migration cifti eklendi:
  - `apps/api/migrations/202603120002_auth_create_core_tables.up.sql`
  - `apps/api/migrations/202603120002_auth_create_core_tables.down.sql`
- Auth stage testleri eklendi:
  - contract: `apps/api/tests/contract/auth_access_contract_test.go`
  - integration: `apps/api/tests/integration/auth_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/auth_migration_integration_test.go`
- Canonical validation ve password altyapisi eklendi:
  - `apps/api/internal/platform/validation/validator.go`
  - `apps/api/internal/shared/crypto/password/argon2id.go`

### Changed
- API bootstrap'ta auth module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- Config'e auth runtime esikleri eklendi:
  - `AUTH_LOGIN_FAILED_ATTEMPT_LIMIT_PER_MINUTE`
  - `AUTH_LOGIN_COOLDOWN_SECONDS`
  - `AUTH_EMAIL_VERIFICATION_RESEND_COOLDOWN_SECONDS`
- `.env.example`, `README.md` ve `docs/upgrade.md` Asama 4 ile hizalandi.
- Auth modulu icinde coklu akis tasiyan tek dosyalar parcalama kuralina gore operasyon bazli ayrildi (`service.go` ve `http_handler.go` kaldirildi, route kayitlari `routes/` altina ayrildi).

### Fixed
- Yok.

### Docs
- Asama 4 auth omurgasi ve versiyonlama guncellemeleri changelog ile izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 4 auth owner omurgasi (credential/session/token/verification/recovery) kod seviyesine tasindi.
- Etkilenen Moduller: `auth`, `platform/validation`, `shared/crypto/password`, `config`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120002_auth_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.3.0-alpha.1] - 2026-03-12

### Added
- `Asama 3` kapsaminda `apps/api/internal/shared/policy` altina olceklenme hazirlik policy seti eklendi:
  - `scaling.go`: domain-group kullanim rehberi, module inventory kayit modeli, module status sozlugu
  - `readmodel.go`: projection/read model uygulama kurallari ve eventual consistency penceresi
  - `reporting.go`: operasyon summary, analytics aggregate ve export query layer sozlugu
  - `reconcile.go`: reconcile guardrail kurallari ve kritik reconcile akis referanslari
  - `maintenance.go`: bakim/refactor disiplini icin canonical guardrail seti
- Asama 3 testleri eklendi:
  - `apps/api/internal/shared/policy/scaling_test.go`
  - `apps/api/tests/contract/shared_scaling_contract_test.go`
  - `apps/api/tests/integration/shared_reconcile_integration_test.go`

### Changed
- `apps/api/internal/shared/README.md` Asama 3 policy kapsamiyla guncellendi.
- `README.md` proje kapsam metni Asama 3 seviyesine cekildi.
- `docs/upgrade.md` durum tablosunda Asama 3 tamamlandi olarak isaretlendi.

### Docs
- Olceklenme, reconcile ve bakim/refactor guardrail setinin kod seviyesindeki karsiliklari changelog ve README uzerinden izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 3 olceklenme hazirlik policy omurgasi kod seviyesine tasindi.
- Etkilenen Moduller: `shared/policy`, `tests/contract`, `tests/integration`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: Yok.

## [0.2.0-alpha.1] - 2026-03-12

### Added
- `Asama 2` kapsaminda `apps/api/internal/shared/catalog` paketi eklendi ve canonical sozlukler kod seviyesinde sabitlendi:
  - audit event types
  - moderation case/assignment/action durumlari
  - notification categories
  - policy effects
  - purchase source ve reward source tipleri
  - support status ve reply visibility
  - target types
  - visibility states
- `apps/api/internal/shared/policy` paketi eklendi:
  - request/correlation izleme alanlari
  - rate-limit surface kayitlari
  - outbox zorunlu bilesenleri ve mesaj alanlari
  - canonical projection kayitlari
  - technical stack, cache/queue, media/reporting/search karar kayitlari
  - transaction boundary referans akislari
- `apps/api/internal/shared/settings` paketi eklendi:
  - audience/scope/disabled behavior/error/entitlement sozlukleri
  - settings kayit semasi modeli
  - access yorumlama sirasi ve kill-switch seviyeleri
  - key grammar validator yardimcilari
- Asama 2 icin yeni test katmani eklendi:
  - unit: `internal/shared/*`
  - contract: `apps/api/tests/contract/shared_catalog_contract_test.go`
  - integration: `apps/api/tests/integration/shared_policy_integration_test.go`

### Changed
- `apps/api/internal/shared/README.md` Asama 2 paketlerini dokumante edecek sekilde guncellendi.
- `README.md` proje kapsam metni Asama 2 seviyesine cekildi.
- `docs/upgrade.md` durum tablosunda Asama 2 tamamlandi olarak isaretlendi.

### Docs
- Shared kararlarin kod karsiliklari changelog ve README uzerinden izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 2 canonical shared sozluk/policy/settings omurgasi kod seviyesine tasindi.
- Etkilenen Moduller: `shared`, `tests/contract`, `tests/integration`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: Yok.

## [0.1.0-alpha.1] - 2026-03-12

### Added
- `Asama 0` kapsaminda repo iskeleti olusturuldu: `apps/`, `docs/`, `scripts/`, `deploy/`, `.github/`.
- `apps/api` altinda temel Go uygulama omurgasi kuruldu (`chi`, `caarlos0/env`, `zap`, `pgxpool`).
- Health (`/health`), readiness (`/ready`) ve version (`/version`) endpointleri eklendi.
- Baslangic migration ciftleri eklendi (`up/down`): `202603120001_core_bootstrap`.
- Docker-first calisma icin `apps/api/Dockerfile` ve `deploy/docker-compose.yml` eklendi.
- Temel script seti eklendi: build, test, docker up/down, migration up/down.
- CI workflow ve PR template omurgasi eklendi.
- `README.md`, `Makefile`, `.env.example`, `docs/SETUP.md` ve `docs/TESTING.md` eklendi.
- `Asama 1` kapsaminda `internal/modules` icin merkezi `Module` kontrati ve registry omurgasi eklendi.
- App bootstrap katmaninda merkezi module route mount akisi eklendi (`NewHTTPHandler` + registry mount).
- `apps/api/tests/contract` ve `apps/api/tests/integration` altinda Asama 1 test katmanlari eklendi.
- Yeni module iskeleti acmak icin `scripts/scaffold_module.ps1` eklendi.
- `internal/app`, `internal/platform`, `internal/shared`, `internal/modules` klasorleri icin katman siniri README dosyalari eklendi.
- Canonical versiyon kaynagi olarak `VERSION` dosyasi eklendi.

### Changed
- `APP_VERSION` runtime zorunlu hale getirildi; versiyon degeri kod icinden default/hardcode olarak alinmiyor.
- Docker compose APP_VERSION fallback'i kaldirilarak zorunlu env kontrolu eklendi.
- CI pipeline test/build adimlarinda `APP_VERSION` `VERSION` dosyasindan export edilecek sekilde guncellendi.

### Fixed
- Konfig yuklenirken `APP_VERSION` bos gecilmesi acik hata ile engelleniyor.

### Docs
- Versiyonlama kullanimi `README.md` ve `docs/SETUP.md` icinde `VERSION` + `APP_VERSION` modeliyle netlestirildi.

### Release Notes
- Degisiklik Ozeti: Asama 0 ve Asama 1 omurgasi ile birlikte versiyonlama modeli kural 16 ile hizalandi.
- Etkilenen Moduller: `app`, `platform/config`, `modules`, `deploy`, `scripts`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120001_core_bootstrap` migration cifti eklendi (uyumlu bootstrap kurulumu).




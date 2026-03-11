# Changelog

Bu dosya yalnizca projede gercekte yapilan islemleri kaydeder.
Bu proje SemVer (`MAJOR.MINOR.PATCH`) standardini takip eder.

## [Unreleased]

### Added
- `Asama 0` kapsaminda repo iskeleti olusturuldu: `apps/`, `docs/`, `scripts/`, `deploy/`, `.github/`.
- `apps/api` altinda temel Go uygulama omurgasi kuruldu (`chi`, `caarlos0/env`, `zap`, `pgxpool`).
- `APP_VERSION` env tabanli versiyonlama, `DB_MAIN_DSN` ve `DB_TEST_DSN` ayrimi tanimlandi.
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

### Changed
- Ana dokumantasyon dosyalari canonical yapiya alinip `docs/` altina tasindi:
  - `docs/rules.md`
  - `docs/roadmap.md`
  - `docs/changelog.md`
  - `docs/modules.md`
  - `docs/shared.md`
- `README.md`, `docs/SETUP.md` ve `docs/TESTING.md` Asama 1 mimari omurga kapsamiyla hizalandi.

### Docs
- `docs/issues.md` ve `docs/upgrade.md` olusturuldu.

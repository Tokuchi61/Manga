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

### Changed
- Ana dokumantasyon dosyalari canonical yapiya alinip `docs/` altina tasindi:
  - `docs/rules.md`
  - `docs/roadmap.md`
  - `docs/changelog.md`
  - `docs/modules.md`
  - `docs/shared.md`

### Docs
- `docs/issues.md` ve `docs/upgrade.md` olusturuldu.

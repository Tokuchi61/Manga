# Changelog

Bu dosya yalnizca projede gercekte yapilan islemleri kaydeder.
Bu proje SemVer (`MAJOR.MINOR.PATCH`) standardini takip eder.

## [Unreleased]

### Changed
- Versiyonlama akisi `VERSION` + `APP_VERSION` tek kaynak prensibi ile hizalandi.

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

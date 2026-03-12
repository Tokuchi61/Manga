# Changelog

Bu dosya yalnizca projede gercekte yapilan islemleri kaydeder.
Bu proje SemVer (`MAJOR.MINOR.PATCH`) standardini takip eder.

## [Unreleased]

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

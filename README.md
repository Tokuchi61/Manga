# NovaScans

NovaScans, oyunlastirilmis manga/manhwa/manhua okuma platformudur.

Bu repo su anda `Asama 3 - Genisleme ve Olceklenme Hazirligi` kapsaminda shared olceklenme policy omurgasini icerir.

## Canonical Versiyon

- Canonical versiyon kaynagi: `VERSION`
- Runtime versiyon kaynagi: `APP_VERSION` environment variable
- Su anki surum: `0.3.0-alpha.1`

## Dizin Yapisi

```text
/apps/
  /api/
    /cmd/
    /internal/
      /app/
      /platform/
      /shared/
      /modules/
    /migrations/
    /tests/
  /web/
/docs/
/scripts/
/deploy/
/.github/
README.md
Makefile
.env.example
VERSION
```

## API Mimari Katmanlar

- `apps/api/internal/app`: bootstrap, composition root ve merkezi route mount
- `apps/api/internal/platform`: config, db, logger ve teknik altyapi kodlari
- `apps/api/internal/shared`: domain-agnostic ortak yapilar
- `apps/api/internal/modules`: leaf moduller ve module registry kontrati

## Asama 2-3 Shared Paketleri

- `apps/api/internal/shared/catalog`: canonical enum ve sozluk kayitlari
- `apps/api/internal/shared/settings`: runtime settings sozlugu, key grameri ve yorumlama modeli
- `apps/api/internal/shared/policy`: transaction/outbox/projection policy'lerine ek olarak AÅŸama 3 olceklenme guardrail'lari:
  - domain-group ve module envanter policy kayitlari
  - projection/read model uygulama kurallari
  - reporting katmanlari ve reconcile akis referanslari
  - bakim/refactor disiplini checklisti

## Modul Iskeleti

Yeni bir leaf modul iskeleti olusturmak icin:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/scaffold_module.ps1 -ModuleName auth
powershell -ExecutionPolicy Bypass -File scripts/scaffold_module.ps1 -ModuleName manga -DomainGroup content
```

## Dokumantasyon

- Ana kurallar: `docs/rules.md`
- Yol haritasi: `docs/roadmap.md`
- Moduller: `docs/modules.md`
- Shared kararlar: `docs/shared.md`
- Changelog: `docs/changelog.md`
- Kurulum: `docs/SETUP.md`
- Test stratejisi: `docs/TESTING.md`

## Hizli Baslangic (Docker)

1. `.env.example` dosyasini referans alarak `.env` olustur.
2. `docker compose -f deploy/docker-compose.yml up -d --build`
3. API kontrolu:
   - `GET http://localhost:8080/health`
   - `GET http://localhost:8080/version`

## Local Gelistirme

```bash
cd apps/api
go test ./...
go build ./...
```

## Versiyonlama

- SemVer kullanilir.
- Runtime versiyonu `APP_VERSION` env degiskeninden okunur.
- `APP_VERSION` degeri kod icinde hardcode edilmez.

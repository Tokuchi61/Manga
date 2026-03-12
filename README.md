# NovaScans

NovaScans, oyunlastirilmis manga/manhwa/manhua okuma platformudur.

Bu repo su anda `Asama 8 - Chapter` kapsaminda kimlik, kullanici, merkezi erisim/policy, manga icerik owner ve chapter okuma/release omurgasini icerir.

## Canonical Versiyon

- Canonical versiyon kaynagi: `VERSION`
- Runtime versiyon kaynagi: `APP_VERSION` environment variable
- Su anki surum: `0.8.0-alpha.1`

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
- `apps/api/internal/platform`: config, db, logger, validation ve teknik altyapi kodlari
- `apps/api/internal/shared`: domain-agnostic ortak yapilar
- `apps/api/internal/modules`: leaf moduller ve module registry kontrati

## Asama 4-8 Omurga

- `apps/api/internal/modules/auth`: register/login/logout, session list/revoke, token refresh rotation, verification ve password reset/change akislari
- `apps/api/internal/modules/user`: profil okuma/guncelleme, public-private profil ayrimi, account state gecisleri, history visibility preference ve VIP lifecycle akislari
- `apps/api/internal/modules/access`: merkezi role/permission/policy yonetimi, temporary grant, own/any authorization ve evaluate karar katmani
- `apps/api/internal/modules/manga`: manga CRUD/listing/detail, search/filter/sort, publish/archive lifecycle, discovery/editorial, counter sync ve soft delete/restore akislarinin owner modulu
- `apps/api/internal/modules/chapter`: chapter CRUD/list/detail/read/navigation, preview/early access state, media/integrity sinyalleri, soft delete/restore ve history resume/read kontrat omurgasi
- `apps/api/internal/shared/crypto/password`: canonical argon2id sifre hash/verify yardimcilari
- `apps/api/internal/platform/validation`: canonical validator wrapper (`go-playground/validator/v10`)
- `apps/api/migrations/202603120002_auth_create_core_tables.*`: auth migration omurgasi
- `apps/api/migrations/202603120003_user_create_core_tables.*`: user migration omurgasi
- `apps/api/migrations/202603120004_access_create_core_tables.*`: access migration omurgasi
- `apps/api/migrations/202603120005_manga_create_core_tables.*`: manga migration omurgasi
- `apps/api/migrations/202603120006_chapter_create_core_tables.*`: chapter migration omurgasi

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

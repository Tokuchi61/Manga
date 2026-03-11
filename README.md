# NovaScans

NovaScans, oyunlastirilmis manga/manhwa/manhua okuma platformudur.

Bu repo su anda `Asama 0 - Temel Standartlar Baslangic` kapsaminda kurulum omurgasini icerir.

## Dizin Yapisi

```text
/apps/
  /api/
  /web/
/docs/
/scripts/
/deploy/
/.github/
README.md
Makefile
.env.example
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

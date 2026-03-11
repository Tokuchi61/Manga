# Kurallar

> Bu dokÃ¼man, proje geneli kurallarÄ±, Ã¶lÃ§eklenebilir mimari ilkeleri ve ortak Ã§alÄ±ÅŸma standartlarÄ±nÄ± toplar. ModÃ¼l ve feature detaylarÄ± ayrÄ± dokÃ¼manlarda geniÅŸletilir; ancak burada tanÄ±mlanan genel kurallar tÃ¼m modÃ¼ller iÃ§in baÄŸlayÄ±cÄ±dÄ±r.

## 1) Proje KimliÄŸi ve Kapsam
- Proje adÄ±: NovaScans.
- Proje alanÄ±: oyunlaÅŸtÄ±rÄ±lmÄ±ÅŸ manga, manhwa ve manhua okuma platformu.
- Bu dokÃ¼man proje geneli kurallarÄ±, mimari ilkeleri ve Ã§alÄ±ÅŸma standartlarÄ±nÄ± tanÄ±mlar.
- Bu dokÃ¼man sprint planÄ± deÄŸildir.
- Bu dokÃ¼man gÃ¶rev listesi deÄŸildir.
- Bu dokÃ¼man tek tek modÃ¼l implementasyon dokÃ¼manÄ± deÄŸildir.
- DetaylÄ± implementasyon kararlarÄ± ve modÃ¼l iÃ§i feature ayrÄ±ntÄ±larÄ± sonraki iterasyonlarda ayrÄ± dokÃ¼manlarda tanÄ±mlanmalÄ±dÄ±r.
- Genel mimari kararlar korunmadan yapÄ±lan geliÅŸtirme tamamlanmÄ±ÅŸ sayÄ±lmaz.

## 2) DokÃ¼manÄ±n AmacÄ± ve KullanÄ±mÄ±
- Bu dokÃ¼man; proje sahibi, geliÅŸtiriciler, AI destekli araÃ§lar ve ajanlar iÃ§in baÄŸlayÄ±cÄ±dÄ±r.
- Ana amaÃ§, proje bÃ¼yÃ¼rken mimariyi, sÄ±nÄ±rlarÄ± ve sÃ¼rdÃ¼rÃ¼lebilirliÄŸi korumaktÄ±r.
- Yeni geliÅŸtirme, refactor, veri modeli deÄŸiÅŸikliÄŸi, API deÄŸiÅŸikliÄŸi, altyapÄ± deÄŸiÅŸikliÄŸi ve sÃ¼reÃ§ tasarÄ±mÄ±nda bu dokÃ¼man esas alÄ±nmalÄ±dÄ±r.
- Projede modÃ¼l, Ã¶zellik, hotfix veya herhangi bir dÃ¼zenleme yapÄ±lacaÄŸÄ± zaman Ã¶nce aktif kurallar dokÃ¼manÄ± baz alÄ±nmalÄ±dÄ±r.
- Genel kurallar ile daha detaylÄ± modÃ¼l dokÃ¼manlarÄ± Ã§eliÅŸirse Ã¶nce bu dokÃ¼mandaki genel mimari ilkeler korunmalÄ±dÄ±r.
- Yeni ana kararlar dokÃ¼mana yansÄ±tÄ±lmadan iÅŸ tamamlanmÄ±ÅŸ kabul edilmemelidir.
- HÄ±zlÄ± Ã§Ã¶zÃ¼m, geÃ§ici Ã§Ã¶zÃ¼m veya acil dÃ¼zeltme gerekÃ§esi bu dokÃ¼mandaki genel mimari kurallarÄ± kalÄ±cÄ± olarak ihlal etme nedeni olamaz.

## 3) Sabit Teknik Kararlar
- Backend dili: Go 1.26.
- Canonical env/config loader olarak `caarlos0/env` kullanÄ±lmalÄ±dÄ±r.
- SQL eriÅŸimi ve connection pooling iÃ§in `pgx/v5` ve `pgxpool` kullanÄ±lmalÄ±dÄ±r.
- Structured logging iÃ§in canonical seÃ§im `zap` olmalÄ±dÄ±r.
- Input validation iÃ§in canonical seÃ§im `go-playground/validator/v10` olmalÄ±dÄ±r.
- UUID Ã¼retimi iÃ§in canonical seÃ§im `google/uuid` olmalÄ±dÄ±r.
- Password hashing iÃ§in canonical seÃ§im `argon2id` olmalÄ±dÄ±r.
- Test assertion ve helper standardÄ± iÃ§in `testify` kullanÄ±lmalÄ±dÄ±r.
- BaÅŸlangÄ±Ã§ async iÅŸleme standardÄ± PostgreSQL-backed jobs + transactional outbox olmalÄ±dÄ±r.
- Cache ihtiyacÄ± oluÅŸtuÄŸunda canonical backend `Redis` olmalÄ±dÄ±r; cache source-of-truth kabul edilmemelidir.
- VeritabanÄ±: PostgreSQL 18.
- HTTP router olarak Chi kullanÄ±lmalÄ±dÄ±r.
- SÃ¼rÃ¼m kontrol sistemi: Git.
- Commit mesajlarÄ± Conventional Commits standardÄ±na uymalÄ±dÄ±r.
- Branch modeli: `main + feature/* + hotfix/*`.
- Migration yÃ¶netiminde `golang-migrate` kullanÄ±lmalÄ±dÄ±r.
- YapÄ±landÄ±rma env tabanlÄ± olmalÄ±, config eriÅŸimi merkezi katmandan yapÄ±lmalÄ±dÄ±r.
- Proje Docker iÃ§inde build alabilmelidir.
- Proje Docker iÃ§inde ayaÄŸa kalkabilmelidir.
- Main DB ve test DB kesin olarak ayrÄ±lmalÄ±dÄ±r.

## 4) GeliÅŸtirme Prensipleri
- Her deÄŸiÅŸiklik doÄŸru sorumluluk alanÄ±nda yapÄ±lmalÄ±dÄ±r.
- Gereksiz refactor yapÄ±lmamalÄ±dÄ±r.
- Gereksiz dosya taÅŸÄ±ma veya isim deÄŸiÅŸikliÄŸi yapÄ±lmamalÄ±dÄ±r.
- Gereksiz paket, feature veya soyutlama eklenmemelidir.
- Yeni yapÄ± eklenmeden Ã¶nce mevcut yapÄ± incelenmelidir.
- AynÄ± sorumluluk birden fazla alana daÄŸÄ±lmamalÄ±dÄ±r.
- OrtaklaÅŸtÄ±rma yalnÄ±zca gerÃ§ek ihtiyaÃ§ oluÅŸtuÄŸunda yapÄ±lmalÄ±dÄ±r.
- GeÃ§ici Ã§Ã¶zÃ¼mler aÃ§Ä±kÃ§a iÅŸaretlenmeli, kalÄ±cÄ± mimari karar gibi bÄ±rakÄ±lmamalÄ±dÄ±r.
- Kod okunabilir, sÃ¼rdÃ¼rÃ¼lebilir, test edilebilir ve izlenebilir olmalÄ±dÄ±r.

## 5) Proje YapÄ±sÄ± ve Dizin Organizasyonu
- Proje Ã§ok modÃ¼llÃ¼ bÃ¼yÃ¼meye ve ileride frontend eklenmesine uygun repo kÃ¶k standardÄ± ile baÅŸlamalÄ±dÄ±r.
- Repo kÃ¶kÃ¼ yalnÄ±zca Ã¼st seviye uygulama dizinleri, ortak dokÃ¼mantasyon, ortak scriptler, deploy dosyalarÄ± ve repo seviyesi yapÄ±landÄ±rmalarÄ± iÃ§ermelidir.
- Backend ve frontend aynÄ± repo iÃ§inde yer alacaksa uygulamalar `apps/` altÄ±nda ayrÄ±ÅŸtÄ±rÄ±lmalÄ±dÄ±r.
- Ã–nerilen kÃ¶k dizin yapÄ±sÄ± aÅŸaÄŸÄ±daki omurgayÄ± korumalÄ±dÄ±r:

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
    go.mod
    go.sum
    Dockerfile
  /web/
/docs/
/scripts/
/deploy/
/.github/
README.md
Makefile
.env.example
```

- `apps/api/` Go backend uygulamasÄ±nÄ±n tek canonical kÃ¶kÃ¼ olmalÄ±dÄ±r.
- `apps/web/` frontend uygulamasÄ±nÄ±n tek canonical kÃ¶kÃ¼ olmalÄ±dÄ±r.
- `apps/api/cmd/` yalnÄ±zca backend uygulama giriÅŸ noktalarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `apps/api/internal/app/` uygulama bootstrap, composition root, dependency wiring ve merkezi baÅŸlatma akÄ±ÅŸlarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `apps/api/internal/platform/` config, DB, logger, middleware, mailer, cache, queue, storage ve benzeri teknik altyapÄ± kodlarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `apps/api/internal/shared/` yalnÄ±zca gerÃ§ekten modÃ¼lden baÄŸÄ±msÄ±z, domain-agnostic ve tekrar kullanÄ±labilir yapÄ±lar iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `apps/api/internal/modules/` tÃ¼m backend iÅŸ modÃ¼llerinin ana yerleÅŸim alanÄ± olmalÄ±dÄ±r.
- `apps/api/migrations/` veritabanÄ± migration dosyalarÄ± iÃ§in tek merkez olmalÄ±dÄ±r.
- `apps/api/tests/` entegrasyon, contract veya uÃ§tan uca testlerin Ã¼st seviye yerleÅŸimi iÃ§in kullanÄ±labilir.
- `docs/` tÃ¼m mimari, sÃ¼reÃ§ ve modÃ¼l belgeleri iÃ§in tek merkez olmalÄ±dÄ±r.
- `scripts/` ortak geliÅŸtirme ve bakÄ±m scriptleri iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `deploy/` Docker, Compose, deployment ve operasyonel Ã§alÄ±ÅŸma dosyalarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `.github/` CI/CD workflow, issue template ve PR template dosyalarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- Repo kÃ¶kÃ¼nde backend uygulama dosyalarÄ± daÄŸÄ±nÄ±k ÅŸekilde tutulmamalÄ±; backend kodu `apps/api/` altÄ±nda toplanmalÄ±dÄ±r.
- Repo kÃ¶kÃ¼nde frontend uygulama dosyalarÄ± daÄŸÄ±nÄ±k ÅŸekilde tutulmamalÄ±; frontend kodu `apps/web/` altÄ±nda toplanmalÄ±dÄ±r.

## 6) ModÃ¼l Organizasyonu ve Ã–lÃ§eklenme KurallarÄ±
- VarsayÄ±lan backend modÃ¼l kÃ¶k dizini `apps/api/internal/modules/<module>/` olmalÄ±dÄ±r.
- Bir modÃ¼lÃ¼n gerÃ§ek kÃ¶k dizini, modÃ¼lÃ¼n kendi leaf klasÃ¶rÃ¼dÃ¼r.
- ModÃ¼l sayÄ±sÄ± arttÄ±ÄŸÄ±nda veya okuma/bakÄ±m maliyeti yÃ¼kseldiÄŸinde opsiyonel domain grubu kullanÄ±labilir.
- Domain grubu kullanÄ±lan durumda yapÄ± `apps/api/internal/modules/<domain-group>/<module>/` formatÄ±na geÃ§ebilir.
- Domain group klasÃ¶rÃ¼ yalnÄ±zca gruplayÄ±cÄ±dÄ±r; gerÃ§ek iÅŸ sÄ±nÄ±rÄ± yine leaf modÃ¼l klasÃ¶rÃ¼dÃ¼r.
- BaÅŸlangÄ±Ã§ aÅŸamasÄ±nda gereksiz klasÃ¶r derinliÄŸi oluÅŸturulmamalÄ±dÄ±r; gerÃ§ek ihtiyaÃ§ yoksa `apps/api/internal/modules/<module>/` yeterlidir.
- Domain group kullanÄ±mÄ± bir gereklilik deÄŸil, Ã¶lÃ§eklenme aracÄ±dÄ±r.
- Backend modÃ¼lleri yalnÄ±zca `apps/api/internal/modules/` altÄ±nda yer almalÄ±dÄ±r.
- Frontend tarafÄ±ndaki feature veya page organizasyonu backend modÃ¼l kÃ¶k yapÄ±sÄ± ile karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.
- Ã–rnek domain group alanlarÄ± baÄŸlayÄ±cÄ± olmadan ÅŸu ÅŸekilde dÃ¼ÅŸÃ¼nÃ¼lebilir:
  - `identity`
  - `content`
  - `community`
  - `operations`
  - `engagement`
  - `commerce`
  - `gameplay`
- Yeni modÃ¼l aÃ§mak iÃ§in aÃ§Ä±k veri sahipliÄŸi, aÃ§Ä±k use-case sÄ±nÄ±rÄ± ve net baÄŸÄ±mlÄ±lÄ±k gerekÃ§esi bulunmalÄ±dÄ±r.
- Alt Ã¶zellikler varsayÄ±lan olarak ayrÄ± modÃ¼l yapÄ±lmamalÄ±, Ã¶nce mevcut modÃ¼l iÃ§inde kalmalÄ±dÄ±r.
- Bir Ã¶zellik ancak baÄŸÄ±msÄ±z veri sahipliÄŸi, baÄŸÄ±msÄ±z servis akÄ±ÅŸÄ± ve baÄŸÄ±msÄ±z access kontratÄ± gerektiriyorsa ayrÄ± modÃ¼le ayrÄ±lmalÄ±dÄ±r.
- ModÃ¼l isimleri kÄ±sa, tek anlamlÄ± ve iÅŸ alanÄ±nÄ± yansÄ±tan canonical adlar olmalÄ±dÄ±r.
- AynÄ± anlama gelen birden fazla modÃ¼l adÄ± aÃ§Ä±lmamalÄ±dÄ±r.
- ModÃ¼l adlarÄ± tekil veya Ã§oÄŸul kullanÄ±m aÃ§Ä±sÄ±ndan tutarlÄ± olmalÄ±; aynÄ± alan iÃ§in iki farklÄ± yazÄ±m standardÄ± aÃ§Ä±lmamalÄ±dÄ±r.
- TÃ¼m aktif modÃ¼ller iÃ§in canonical modÃ¼l kaydÄ± tutulmalÄ±dÄ±r.
- Ã–nerilen kayÄ±t dosyasÄ± `docs/modules/index.md` veya benzeri merkezi bir modÃ¼l envanteri olmalÄ±dÄ±r.
- Bu kayÄ±t en az ÅŸu alanlarÄ± iÃ§ermelidir:
  - canonical modÃ¼l adÄ±
  - varsa domain group
  - kÄ±sa aÃ§Ä±klama
  - durum
  - ana dokÃ¼man yolu
- `durum` alanÄ± iÃ§in Ã¶nerilen canonical deÄŸerler: `planned`, `active`, `deprecated`, `archived`.

## 7) Standart ModÃ¼l YapÄ±sÄ± ve Katman Ä°lkeleri
- Her leaf modÃ¼l, yalnÄ±zca ihtiyaÃ§ duyduÄŸu dosya ve klasÃ¶rleri iÃ§ermelidir; ancak kullanÄ±lan yapÄ±lar ortak isimlendirme ve katman standardÄ±nÄ± korumalÄ±dÄ±r.
- AÅŸaÄŸÄ±daki yapÄ± zorunlu tam ÅŸablon deÄŸildir; modÃ¼lde ihtiyaÃ§ doÄŸduÄŸunda hangi klasÃ¶r veya dosyanÄ±n hangi amaÃ§la aÃ§Ä±lacaÄŸÄ±nÄ± gÃ¶steren standart backend omurgadÄ±r:

```text
apps/api/internal/modules/<module>/
  entity/
  dto/
  service/
  repository/
  handler/
  middleware/
  validator/
  mapper/
  contract/
  events/
  consumer/
  producer/
  jobs/
  readmodel/
  errors.go
  module.go
  routes.go
```

- Domain group kullanÄ±lÄ±yorsa aynÄ± mantÄ±k `apps/api/internal/modules/<domain-group>/<module>/` altÄ±nda uygulanmalÄ±dÄ±r.

- `entity/` modÃ¼lÃ¼n kendi veri yapÄ±larÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `dto/` request, response ve modÃ¼l dÄ±ÅŸÄ± contract veri yapÄ±larÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `service/` iÅŸ kurallarÄ± ve use-case akÄ±ÅŸlarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `repository/` yalnÄ±zca veri eriÅŸimi iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `handler/` HTTP handler, request parse ve response Ã¼retimi gibi giriÅŸ katmanÄ± sorumluluklarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `middleware/` modÃ¼le Ã¶zel middleware yapÄ±larÄ± iÃ§in kullanÄ±labilir.
- `validator/` modÃ¼le Ã¶zel doÄŸrulama kurallarÄ± ve input validation yardÄ±mcÄ±larÄ± iÃ§in kullanÄ±labilir.
- `mapper/` entity, dto, response, read model veya contract dÃ¶nÃ¼ÅŸÃ¼mleri karmaÅŸÄ±k hale geldiÄŸinde kullanÄ±labilir.
- `contract/` diÄŸer modÃ¼llerle paylaÅŸÄ±lan resmi modÃ¼l kontratlarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `events/` modÃ¼lÃ¼n yayÄ±nladÄ±ÄŸÄ± veya tÃ¼kettiÄŸi event tanÄ±mlarÄ± iÃ§in kullanÄ±labilir.
- `consumer/` queue consumer, event consumer, webhook consumer veya dÄ±ÅŸ giriÅŸ entegrasyonlarÄ± iÃ§in kullanÄ±labilir.
- `producer/` message producer, event producer veya dÄ±ÅŸ sisteme yayÄ±n yapan entegrasyon akÄ±ÅŸlarÄ± iÃ§in kullanÄ±labilir.
- `jobs/` asenkron iÅŸleyiciler veya zamanlanmÄ±ÅŸ iÅŸler iÃ§in kullanÄ±labilir.
- `readmodel/` yalnÄ±zca okuma odaklÄ± Ã¶zel projection veya denormalize yapÄ±lar gerektiÄŸinde kullanÄ±lmalÄ±dÄ±r.
- `module.go` gerekiyorsa modÃ¼lÃ¼n composition veya registration giriÅŸ noktasÄ± olmalÄ±dÄ±r.
- `routes.go` gerekiyorsa modÃ¼lÃ¼n route kayÄ±t giriÅŸ noktasÄ± olmalÄ±dÄ±r.
- Bu bÃ¶lÃ¼mde listelenen klasÃ¶r ve dosyalarÄ±n hiÃ§biri her modÃ¼lde zorunlu deÄŸildir; zorunlu olan, ihtiyaÃ§ doÄŸduÄŸunda burada tanÄ±mlanan amaÃ§ ve isim standardÄ±na uyulmasÄ±dÄ±r.
- Tekil ve modÃ¼l seviyesinde kalan dosyalar modÃ¼l kÃ¶k dizininde tutulmalÄ±dÄ±r.
- Ã–rnek tekil modÃ¼l kÃ¶k dosyalarÄ±: `module.go`, `routes.go`, `errors.go`, `constants.go`, `types.go`, `service.go`, `repository.go`, `handler.go`, `middleware.go`, `validator.go`, `mapper.go`, `contract.go`, `events.go`, `consumer.go`, `producer.go`, `jobs.go`, `readmodel.go`.
- Tek bir katman veya modÃ¼l bileÅŸeni tek dosya ile temsil edilebiliyorsa ilgili dosya modÃ¼l kÃ¶kÃ¼nde tutulabilir.
- Bu kÃ¶k dosyalar yalnÄ±zca tek bir aÃ§Ä±k sorumluluk veya tek bir akÄ±ÅŸ ailesi taÅŸÄ±dÄ±ÄŸÄ± sÃ¼rece kÃ¶kte kalmalÄ±dÄ±r.
- Bir alan birden fazla iÅŸlem, akÄ±ÅŸ veya dosya gerektiriyorsa modÃ¼l kÃ¶kÃ¼nde bÃ¼yÃ¼tÃ¼lmemeli, aynÄ± amaÃ§la aÃ§Ä±lan klasÃ¶r altÄ±nda parÃ§alanmalÄ±dÄ±r.
- Ã–rnek dÃ¶nÃ¼ÅŸÃ¼m: `service.go -> service/`, `handler.go -> handler/`, `middleware.go -> middleware/`, `routes.go -> routes/`, `events.go -> events/`, `consumer.go -> consumer/`.
- `errors.go`, `constants.go`, `types.go` veya benzeri tekil kÃ¶k dosyalar tekil olmaktan Ã§Ä±karsa daha aÃ§Ä±k isimli dosyalara veya uygun klasÃ¶r yapÄ±sÄ±na ayrÄ±ÅŸtÄ±rÄ±lmalÄ±dÄ±r.
- Bu kural yalnÄ±zca belirli birkaÃ§ klasÃ¶r iÃ§in deÄŸil, modÃ¼l iÃ§indeki tÃ¼m Ã§oklu iÅŸlev alanlarÄ± iÃ§in geÃ§erlidir.
- `service/`, `repository/`, `dto/`, `entity/`, `handler/`, `middleware/`, `validator/`, `mapper/`, `contract/`, `events/`, `consumer/`, `producer/`, `jobs/`, `readmodel/` ve gelecekte aÃ§Ä±lacak benzer Ã§oklu iÅŸlev klasÃ¶rleri aynÄ± parÃ§alama ilkesine uymalÄ±dÄ±r.
- ModÃ¼l kÃ¶kÃ¼ tekil dosyalar iÃ§indir; Ã§oklu iÅŸlev taÅŸÄ±yan bÃ¼yÃ¼k dosyalar kÃ¶kte biriktirilmemelidir.
- Birden fazla use-case veya iÅŸlem iÃ§eren modÃ¼llerde her ana iÅŸlem mÃ¼mkÃ¼n olduÄŸunda ayrÄ± dosyada tutulmalÄ±dÄ±r.
- Service dosyalarÄ± iÅŸlem veya use-case bazlÄ± parÃ§alanmalÄ±dÄ±r.
- Repository dosyalarÄ± mÃ¼mkÃ¼n olduÄŸunda entity, aggregate veya belirgin veri eriÅŸim sorumluluÄŸu bazlÄ± parÃ§alanmalÄ±dÄ±r.
- DTO dosyalarÄ± request ve response olarak ayrÄ±lmalÄ±; birden fazla farklÄ± akÄ±ÅŸ tek DTO dosyasÄ±na doldurulmamalÄ±dÄ±r.
- Entity dosyalarÄ± tek bir dev entity dosyasÄ±na dÃ¶nÃ¼ÅŸmemeli; alan kÃ¼meleri, alt yapÄ±lar veya anlamlÄ± domain ayrÄ±mlarÄ± varsa kontrollÃ¼ ÅŸekilde bÃ¶lÃ¼nmelidir.
- Handler dosyalarÄ± mÃ¼mkÃ¼n olduÄŸunda endpoint veya iÅŸlem ailesi bazÄ±nda parÃ§alanmalÄ±dÄ±r.
- Middleware dosyalarÄ± farklÄ± middleware akÄ±ÅŸlarÄ±nÄ± tek dosyada gereksiz ÅŸekilde biriktirmemelidir.
- Validator dosyalarÄ± farklÄ± iÅŸlem ailelerini tek doÄŸrulama dosyasÄ±nda toplamamalÄ±dÄ±r.
- Mapper dosyalarÄ± farklÄ± dÃ¶nÃ¼ÅŸÃ¼m ailelerini tek dev dosyada toplamamalÄ±dÄ±r.
- Contract dosyalarÄ± farklÄ± entegrasyon akÄ±ÅŸlarÄ±nÄ± tek dosyada biriktirmemeli; entegrasyon yÃ¼zeyleri gerektiÄŸinde iÅŸlem ailesi bazÄ±nda ayrÄ±lmalÄ±dÄ±r.
- Event dosyalarÄ± producer veya consumer tarafÄ±nda farklÄ± domain olaylarÄ±nÄ± anlamsÄ±z ÅŸekilde tek dosyada biriktirmemelidir.
- Consumer dosyalarÄ± farklÄ± dÄ±ÅŸ giriÅŸ akÄ±ÅŸlarÄ±nÄ± tek dosyada biriktirmemelidir.
- Producer dosyalarÄ± farklÄ± dÄ±ÅŸ yayÄ±n akÄ±ÅŸlarÄ±nÄ± tek dosyada biriktirmemelidir.
- Job dosyalarÄ± zamanlanmÄ±ÅŸ iÅŸler veya asenkron akÄ±ÅŸlar bazÄ±nda ayrÄ±ÅŸtÄ±rÄ±lmalÄ±dÄ±r.
- Read model dosyalarÄ± tek bir dev projection dosyasÄ±na dÃ¶nÃ¼ÅŸmemelidir.
- ModÃ¼l kÃ¶kÃ¼ndeki tekil dosyalar yalnÄ±zca kendi tekil modÃ¼l sorumluluÄŸunu taÅŸÄ±malÄ±dÄ±r.
- EÄŸer route, middleware, error mapping veya registration alanÄ± tekil olmaktan Ã§Ä±kÄ±p birden fazla akÄ±ÅŸ taÅŸÄ±maya baÅŸlarsa bu yapÄ± kÃ¶kte bÃ¼yÃ¼tÃ¼lmemeli, uygun klasÃ¶r altÄ±nda ayrÄ±ÅŸtÄ±rÄ±lmalÄ±dÄ±r.
- AynÄ± modÃ¼lde 5-6 farklÄ± ana akÄ±ÅŸÄ± tek `service.go`, tek `repository.go`, tek `dto.go` veya benzeri tek bir dosyada toplamak mimari olarak yanlÄ±ÅŸ kabul edilmelidir.
- Dosya parÃ§alama keyfi deÄŸil, okunabilirlik ve sorumluluk ayrÄ±mÄ± amacÄ±yla yapÄ±lmalÄ±dÄ±r; anlamsÄ±z aÅŸÄ±rÄ± bÃ¶lme de yapÄ±lmamalÄ±dÄ±r.
- Dosya adlarÄ±nda modÃ¼l prefixi ve iÅŸlem adÄ± birlikte kullanÄ±lmalÄ±dÄ±r.
- Ã–rnek yaklaÅŸÄ±m:

```text
service/
  auth_login_service.go
  auth_register_service.go
  auth_session_service.go

repository/
  auth_user_repository.go
  auth_session_repository.go

dto/
  auth_login_request_dto.go
  auth_login_response_dto.go
  auth_register_request_dto.go
```

- Bu Ã¶rnek baÄŸlayÄ±cÄ± isim listesi deÄŸildir; baÄŸlayÄ±cÄ± olan ilke, dosyanÄ±n tek bir ana sorumluluÄŸu ve aÃ§Ä±k bir adÄ± olmasÄ±dÄ±r.
- Yeni iÅŸlem eklendiÄŸinde varsayÄ±lan yaklaÅŸÄ±m mevcut bÃ¼yÃ¼k dosyayÄ± daha da bÃ¼yÃ¼tmek deÄŸil, modÃ¼l iÃ§indeki doÄŸru klasÃ¶rde veya gerekiyorsa modÃ¼l kÃ¶kÃ¼nde yeni iÅŸlem dosyasÄ± aÃ§mak olmalÄ±dÄ±r.
- Refactor sÄ±rasÄ±nda devasa dosyalar gÃ¶rÃ¼lÃ¼rse Ã¶nce iÅŸlem sÄ±nÄ±rlarÄ± ayrÄ±ÅŸtÄ±rÄ±lmalÄ±, sonra dosyalar kontrollÃ¼ ÅŸekilde bÃ¶lÃ¼nmelidir.
- ModÃ¼l dizin adlarÄ± ve katman klasÃ¶rleri kÃ¼Ã§Ã¼k harfli, ASCII uyumlu ve tek bir canonical yazÄ±m ile oluÅŸturulmalÄ±dÄ±r.
- Dizin, modÃ¼l ve klasÃ¶r adlarÄ±nda boÅŸluk, tire veya rastgele kÄ±saltma kullanÄ±lmamalÄ±dÄ±r.
- Go package adlarÄ± kÃ¼Ã§Ã¼k harfli, kÄ±sa, tek anlamlÄ± ve mÃ¼mkÃ¼n olduÄŸunda klasÃ¶r adÄ± ile uyumlu olmalÄ±dÄ±r.
- Go package adlarÄ±nda anlamsÄ±z kÄ±saltmalar, underscore veya birden fazla kelimeyi gereksiz ÅŸekilde birleÅŸtiren karmaÅŸÄ±k adlar kullanÄ±lmamalÄ±dÄ±r.
- Dosya adlarÄ± kÃ¼Ã§Ã¼k harfli `snake_case` biÃ§iminde, aÃ§Ä±k sorumluluk ve aÃ§Ä±k iÅŸlem adÄ± ile oluÅŸturulmalÄ±dÄ±r.
- `util.go`, `utils.go`, `helper.go`, `common.go`, `misc.go`, `temp.go` gibi belirsiz dosya adlarÄ± kullanÄ±lmamalÄ±dÄ±r.
- Export edilen Go identifier adlarÄ± `PascalCase`, export edilmeyen adlar `camelCase` standardÄ±na uymalÄ±dÄ±r.
- KÄ±saltmalar tutarlÄ± kullanÄ±lmalÄ±dÄ±r: `ID`, `API`, `HTTP`, `URL` gibi yaygÄ±n kÄ±saltmalar proje genelinde aynÄ± biÃ§imde yazÄ±lmalÄ±dÄ±r.
- Katmanlar arasÄ± baÄŸÄ±mlÄ±lÄ±k yÃ¶nÃ¼ kontrollÃ¼ olmalÄ±dÄ±r.
- HTTP katmanÄ± iÅŸ kararÄ± vermemelidir.
- Service katmanÄ± HTTP detayÄ±na baÄŸÄ±mlÄ± olmamalÄ±dÄ±r.
- Repository katmanÄ± business karar vermemelidir.
- Entity yapÄ±larÄ± baÅŸka modÃ¼lÃ¼n repository detayÄ±na baÄŸÄ±mlÄ± olmamalÄ±dÄ±r.
- Frontend klasÃ¶r yapÄ±sÄ± bu backend modÃ¼l standardÄ±na gÃ¶re ÅŸekillendirilmemelidir.
- Backend modÃ¼l standardÄ± yalnÄ±zca `apps/api/` iÃ§indeki Go uygulamasÄ± iÃ§in geÃ§erlidir.

## 8) Shared, Platform ve Bootstrap SÄ±nÄ±rlarÄ±
- `platform` teknik altyapÄ±nÄ±n sahibidir; iÅŸ kuralÄ± barÄ±ndÄ±rmamalÄ±dÄ±r.
- `shared` gerÃ§ekten ortak ve domain-agnostic yapÄ±larÄ±n sahibidir; modÃ¼l bazlÄ± iÅŸ mantÄ±ÄŸÄ± iÃ§eremez.
- `shared` klasÃ¶rÃ¼, modÃ¼llerden kod kaÃ§Ä±rma alanÄ±na dÃ¶nÃ¼ÅŸmemelidir.
- Bir kod parÃ§asÄ± yalnÄ±zca iki farklÄ± modÃ¼l tarafÄ±ndan tekrar kullanÄ±ldÄ±ÄŸÄ± iÃ§in otomatik olarak `shared` iÃ§ine taÅŸÄ±nmamalÄ±dÄ±r.
- `app` veya bootstrap katmanÄ± modÃ¼lleri birbirine baÄŸlar; modÃ¼ller bootstrap katmanÄ±nÄ± import etmemelidir.
- ModÃ¼l baÄŸÄ±mlÄ±lÄ±klarÄ± merkezi wiring katmanÄ±nda birleÅŸtirilmelidir.
- Route mount iÅŸlemleri merkezi uygulama baÅŸlangÄ±Ã§ katmanÄ±nda yapÄ±lmalÄ±dÄ±r.

## 9) Ã‡apraz ModÃ¼l BaÄŸÄ±mlÄ±lÄ±k ve Entegrasyon KurallarÄ±
- Bir modÃ¼l baÅŸka modÃ¼lÃ¼n tablolarÄ±na doÄŸrudan yazmamalÄ±dÄ±r.
- Bir modÃ¼l baÅŸka modÃ¼lÃ¼n repository implementasyonunu doÄŸrudan kullanmamalÄ±dÄ±r.
- Bir modÃ¼l baÅŸka modÃ¼lÃ¼n iÃ§ klasÃ¶rlerine doÄŸrudan baÄŸÄ±mlÄ± olmamalÄ±dÄ±r.
- BaÅŸka bir modÃ¼lÃ¼n `handler/`, `repository/`, `entity/`, `jobs/`, `consumer/`, `producer/` veya benzeri iÃ§ yapÄ±larÄ± dÄ±ÅŸarÄ±ya aÃ§Ä±k API kabul edilmemelidir.
- Her modÃ¼l kendi public surface'inin tek sahibidir.
- Bir modÃ¼lÃ¼n public surface'i yalnÄ±zca aÃ§Ä±kÃ§a belgelediÄŸi ve dÄ±ÅŸa aÃ§tÄ±ÄŸÄ± ÅŸu yÃ¼zeylerden oluÅŸmalÄ±dÄ±r:
  - `contract/`
  - dÄ±ÅŸ kullanÄ±ma aÃ§Ä±lmÄ±ÅŸ service interface'leri
  - dÄ±ÅŸ kullanÄ±m iÃ§in tanÄ±mlanmÄ±ÅŸ DTO veya contract modelleri
  - yayÄ±nlanan event sÃ¶zleÅŸmeleri
  - gerekiyorsa aÃ§Ä±k read model veya projection yÃ¼zeyi
- Bu yÃ¼zey dÄ±ÅŸÄ±nda kalan tÃ¼m yapÄ± modÃ¼lÃ¼n iÃ§ implementasyonudur ve dÄ±ÅŸ baÄŸÄ±mlÄ±lÄ±k noktasÄ± olarak kullanÄ±lamaz.
- Senkron service veya contract sÃ¶zleÅŸmesinin sahipliÄŸi provider modÃ¼ldedir.
- Event ÅŸemasÄ±nÄ±n sahipliÄŸi event'i publish eden producer modÃ¼ldedir.
- Read model veya projection yÃ¼zeyinin sahipliÄŸi ilgili veriyi servis eden modÃ¼ldedir.
- Consumer modÃ¼l kendi iÃ§inde local interface veya adapter tanÄ±mlayabilir; ancak bu durum provider modÃ¼lÃ¼n iÃ§ klasÃ¶rlerini resmi public API'ye dÃ¶nÃ¼ÅŸtÃ¼rmez.
- ModÃ¼ller arasÄ± iletiÅŸim yalnÄ±zca aÅŸaÄŸÄ±daki yollardan biriyle kurulmalÄ±dÄ±r:
  - aÃ§Ä±k service/contract arayÃ¼zÃ¼
  - aÃ§Ä±k DTO/contract modeli
  - aÃ§Ä±k event sÃ¶zleÅŸmesi
  - aÃ§Ä±k read model veya projection ihtiyacÄ±
- Ã‡apraz modÃ¼l veri ihtiyacÄ± iÃ§in diÄŸer modÃ¼lÃ¼n entity yapÄ±sÄ± doÄŸrudan dÄ±ÅŸarÄ± sÄ±zdÄ±rÄ±lmamalÄ±dÄ±r.
- Bir modÃ¼lÃ¼n taÅŸÄ±dÄ±ÄŸÄ± denormalize sayaÃ§, Ã¶zet veya projection alanÄ± o modÃ¼lde tutulabilir; ancak bu alanlarÄ±n canonical kaynak verisi ilgili kaynak modÃ¼lde kalmalÄ±dÄ±r.
- Ã‡apraz modÃ¼l sayaÃ§ veya Ã¶zet gÃ¼ncellemeleri kaynak modÃ¼l tarafÄ±ndan owner modÃ¼lÃ¼n tablosuna doÄŸrudan yazÄ±larak yapÄ±lmamalÄ±dÄ±r; aÃ§Ä±k event, projection veya owner modÃ¼lÃ¼n aÃ§Ä±k counter contract yÃ¼zeyi kullanÄ±lmalÄ±dÄ±r.
- SayaÃ§ veya Ã¶zet alanÄ± tanÄ±mlanan her yerde canonical source, gÃ¼ncelleme tetikleyicisi, kabul edilen gecikme modeli ve gerektiÄŸinde reconcile veya yeniden hesaplama yolu dokÃ¼mante edilmelidir.
- `target_type` veya benzeri paylaÅŸÄ±lan hedef tipleri taÅŸÄ±yan modÃ¼ller canonical kayÄ±t dosyasÄ± olarak `docs/shared/target-types.md` kullanmalÄ±dÄ±r.
- Yeni target type yalnÄ±zca ilgili modÃ¼l dokÃ¼manÄ±, consumer modÃ¼l dokÃ¼manlarÄ± ve `docs/shared/target-types.md` aynÄ± deÄŸiÅŸiklik setinde gÃ¼ncellendiÄŸinde kullanÄ±labilir hale gelmelidir.
- DÃ¶ngÃ¼sel modÃ¼l baÄŸÄ±mlÄ±lÄ±ÄŸÄ± kesin olarak yasaktÄ±r.
- Senkron baÄŸÄ±mlÄ±lÄ±k yalnÄ±zca gerÃ§ekten anlÄ±k doÄŸrulama veya kritik iÅŸlem bÃ¼tÃ¼nlÃ¼ÄŸÃ¼ gerektiÄŸinde kullanÄ±lmalÄ±dÄ±r.
- ZayÄ±f baÄŸlÄ± entegrasyonlarda event tabanlÄ± veya asenkron akÄ±ÅŸ tercih edilmelidir.
- Event tabanlÄ± entegrasyon kullanÄ±lan yerlerde event adÄ±, payload, producer, consumer ve idempotency beklentisi dokÃ¼mante edilmelidir.
- Public surface Ã¼zerinde yapÄ±lan her deÄŸiÅŸiklik modÃ¼l dokÃ¼manÄ±na, gerekiyorsa changelog kaydÄ±na ve versiyonlama deÄŸerlendirmesine yansÄ±tÄ±lmalÄ±dÄ±r.

## 10) Yeni ModÃ¼l AÃ§ma Zorunlu Kriterleri
- Yeni modÃ¼l aÃ§Ä±lmadan Ã¶nce en az aÅŸaÄŸÄ±daki baÅŸlÄ±klar tanÄ±mlanmalÄ±dÄ±r:
  - sorumluluk alanÄ±
  - veri sahipliÄŸi
  - dÄ±ÅŸa aÃ§Ä±lan API veya contract sÄ±nÄ±rÄ±
  - access veya authorization entegrasyonu
  - state ve lifecycle yapÄ±sÄ±
  - diÄŸer modÃ¼llerle baÄŸÄ±mlÄ±lÄ±k iliÅŸkisi
  - event ihtiyacÄ± varsa event sÃ¶zleÅŸmesi
  - shared registry etkisi varsa ilgili canonical kayÄ±t dosyasÄ±
  - config ihtiyacÄ±
  - migration ihtiyacÄ±
  - test gereksinimleri
  - log ve audit gereksinimleri
  - dokÃ¼mantasyon dosyasÄ±
- Bu baÅŸlÄ±klar netleÅŸmeden yeni modÃ¼l implementasyonu baÅŸlatÄ±lmamalÄ±dÄ±r.
- Her yeni modÃ¼l iÃ§in en az bir modÃ¼l dokÃ¼manÄ± oluÅŸturulmalÄ±dÄ±r.
- Domain group kullanÄ±lmÄ±yorsa modÃ¼l dokÃ¼manlarÄ± iÃ§in Ã¶nerilen yerleÅŸim `docs/modules/<module>.md` olmalÄ±dÄ±r.
- Domain group kullanÄ±lÄ±yorsa modÃ¼l dokÃ¼manlarÄ± iÃ§in Ã¶nerilen yerleÅŸim `docs/modules/<domain-group>/<module>.md` olmalÄ±dÄ±r.
- Yeni leaf modÃ¼l aÃ§Ä±ldÄ±ÄŸÄ±nda merkezi modÃ¼l envanteri aynÄ± deÄŸiÅŸiklik seti iÃ§inde eklenmeli veya gÃ¼ncellenmelidir.

## 11) Veri, API ve State KurallarÄ±
- Åema iliÅŸkileri temiz, aÃ§Ä±k ve tutarlÄ± kurulmalÄ±dÄ±r.
- Foreign key mantÄ±ÄŸÄ± net olmalÄ±dÄ±r.
- Soft delete kullanÄ±lan alanlar kontrollÃ¼ uygulanmalÄ±dÄ±r.
- Unique alanlar aÃ§Ä±k tanÄ±mlanmalÄ± ve Ã§akÄ±ÅŸma senaryolarÄ± yÃ¶netilmelidir.
- Gerekli alanlarda makul index kullanÄ±lmalÄ±dÄ±r.
- Veri modeli gereksiz karmaÅŸÄ±klÄ±k oluÅŸturmadan bÃ¼yÃ¼meye uygun olmalÄ±dÄ±r.
- API stili REST + JSON olmalÄ±dÄ±r.
- Request ve response modelleri aÃ§Ä±kÃ§a ayrÄ±lmalÄ±dÄ±r.
- DB modeli doÄŸrudan API response olarak kullanÄ±lmamalÄ±dÄ±r.
- Standart hata cevabÄ± formatÄ± kullanÄ±lmalÄ±dÄ±r.
- `401` yalnÄ±zca kimlik doÄŸrulama eksik veya geÃ§ersiz olduÄŸunda kullanÄ±lmalÄ±dÄ±r.
- `403` yetki, policy veya access kararÄ± nedeniyle reddedilen iÅŸlemlerde kullanÄ±lmalÄ±dÄ±r.
- `404` gÃ¶rÃ¼nÃ¼rlÃ¼k veya kaynak gizleme kararÄ± nedeniyle dÄ±ÅŸarÄ±ya kapalÄ± bÄ±rakÄ±lan kaynaklarda kullanÄ±lmalÄ±dÄ±r.
- `429` rate limit, throttling, cooldown veya benzeri geÃ§ici eÅŸik ihlali durumlarÄ±nda kullanÄ±lmalÄ±dÄ±r.
- `503` bakÄ±m modu, kill switch, emergency deny veya sistem kaynaklÄ± geÃ§ici pasiflik durumlarÄ±nda kullanÄ±lmalÄ±dÄ±r.
- Liste endpoint'lerinde pagination, filter ve sort parametreleri tutarlÄ± adlandÄ±rÄ±lmalÄ±dÄ±r.
- State, visibility, moderation ve publish kavramlarÄ± birbirine karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.
- State deÄŸiÅŸiklikleri kontrollÃ¼, izlenebilir ve yetki kontrolÃ¼ altÄ±nda olmalÄ±dÄ±r.
- Admin tarafÄ±ndan yÃ¶netilen runtime ayarlar, modÃ¼l durumu ve Ã¶zellik durumu veri modeli seviyesinde aÃ§Ä±k scope ile temsil edilebilmelidir.
- Sistem en az site geneli, modÃ¼l bazlÄ±, Ã¶zellik bazlÄ± ve gerektiÄŸinde context veya resource bazlÄ± aÃ§ma-kapama davranÄ±ÅŸlarÄ±nÄ± destekleyebilmelidir.
- Bir modÃ¼l veya alt Ã¶zellik pasife alÄ±ndÄ±ÄŸÄ±nda beklenen fallback davranÄ±ÅŸÄ±, gÃ¶rÃ¼nÃ¼rlÃ¼k etkisi ve hata cevabÄ± aÃ§Ä±kÃ§a tanÄ±mlanmalÄ±dÄ±r.
- Ä°ÅŸ kurallarÄ±nÄ± etkileyen eÅŸik, oran, sÃ¼re ve limit deÄŸerleri mÃ¼mkÃ¼n olduÄŸunda sabit koda gÃ¶mÃ¼lmemeli; yÃ¶netilebilir ayar yÃ¼zeyi ile kontrol edilebilmelidir.
- Runtime ayarlarda `scope`, ayarÄ±n nerede uygulandÄ±ÄŸÄ±nÄ±; `audience`, ayarÄ±n kim iÃ§in uygulandÄ±ÄŸÄ±nÄ± ifade eder. Bu iki kavram birbirine karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.
- BaÅŸlangÄ±Ã§ audience kapsamÄ± en az `all`, `guest`, `authenticated`, `authenticated_non_vip` ve `vip` seviyelerini destekleyebilmelidir.
- Ä°leri audience geniÅŸletmeleri `role:<name>`, `group:<name>` veya `user:<id>` gibi hedeflemeler olabilir; bu tÃ¼r kapsamlar aÃ§Ä±kÃ§a dokÃ¼mana eklenmeden kullanÄ±lmamalÄ±dÄ±r.
- Runtime kontrol modeli tek bir global `disabled` bayraÄŸÄ±na indirgenmemelidir; gerektiÄŸinde `read`, `write`, `intake`, `preview`, `visibility` ve `benefit` gibi yÃ¼zeyler ayrÄ± ayrÄ± yÃ¶netilebilmelidir.
- Runtime kapatma davranÄ±ÅŸlarÄ± canonical olarak en az `visibility_off`, `read_only`, `write_off`, `intake_pause`, `preview_off` ve `benefit_pause` tiplerini ifade edebilmelidir.
- `disabled_behavior` iÅŸlevsel davranÄ±ÅŸÄ± ifade eder; tek baÅŸÄ±na HTTP cevap kodunu belirlemez. API davranÄ±ÅŸÄ± iÃ§in ayrÄ±ca `error_response_policy` tanÄ±mlanmalÄ±dÄ±r.
- `error_response_policy` en az `not_found`, `forbidden`, `rate_limited` ve `temporarily_unavailable` deÄŸerlerini desteklemelidir.
- VarsayÄ±lan hizalama olarak gÃ¶rÃ¼nÃ¼rlÃ¼k gizleme kararlarÄ± `not_found`, eÅŸik veya cooldown ihlalleri `rate_limited`, sistem kaynaklÄ± pause veya kill switch kararlarÄ± `temporarily_unavailable`, kullanÄ±cÄ±ya aÃ§Ä±k ama yazma veya eriÅŸim kÄ±sÄ±tÄ± taÅŸÄ±yan yÃ¼zeyler ise `forbidden` ile modellenmelidir.
- Availability veya kill switch tÃ¼rÃ¼ndeki ayarlarda gÃ¼venlik Ã¶nceliklidir; eÅŸleÅŸen bir `emergency_deny` her durumda en yÃ¼ksek Ã¶nceliÄŸe sahip olmalÄ±dÄ±r.
- Availability veya kill switch tÃ¼rÃ¼ndeki ayarlarda eÅŸleÅŸen herhangi bir `deny/off` kuralÄ±, `allow/on` kuralÄ±nÄ± bastÄ±rmalÄ±dÄ±r.
- AynÄ± availability anahtarÄ± iÃ§in aynÄ± `audience_kind + audience_selector` ve aynÄ± `scope_kind + scope_selector` kombinasyonunda birden fazla aktif kural bÄ±rakÄ±lamaz; bu tÃ¼r Ã§akÄ±ÅŸmalar kayÄ±t aÅŸamasÄ±nda reddedilmelidir.
- EÅŸik veya deÄŸer taÅŸÄ±yan runtime ayarlarda en spesifik geÃ§erli kayÄ±t kullanÄ±lmalÄ±dÄ±r.
- EÅŸik veya deÄŸer ayarlarÄ±nda audience Ã¶zgÃ¼llÃ¼k sÄ±rasÄ± `user/group/role` > `vip/authenticated_non_vip` > `authenticated/guest` > `all` olmalÄ±dÄ±r.
- EÅŸik veya deÄŸer ayarlarÄ±nda scope Ã¶zgÃ¼llÃ¼k sÄ±rasÄ± `resource/context` > `feature` > `module` > `site` olmalÄ±dÄ±r.
- Bir modÃ¼l dokÃ¼manÄ± bir yÃ¼zeyin ayrÄ± ayrÄ± runtime kontrol edilebildiÄŸini sÃ¶ylÃ¼yorsa `docs/settings/index.md` iÃ§inde o yÃ¼zey iÃ§in en az bir canonical baseline key veya bu alt yÃ¼zeyleri aÃ§Ä±kÃ§a kapsayan umbrella key kaydÄ± bulunmalÄ±dÄ±r.
- Umbrella key kullanÄ±lan durumda kapsanan alt yÃ¼zeyler `affected_surfaces` ve `notes` alanlarÄ±nda aÃ§Ä±kÃ§a listelenmeli; yÃ¼zey dokÃ¼manda tanÄ±mlÄ± kalÄ±rken settings envanterinde tamamen isimsiz bÄ±rakÄ±lamaz.
- Ãœcretli veya sÃ¼reli haklarÄ± etkileyen runtime kapatmalar aÃ§Ä±k bir entitlement impact policy taÅŸÄ±mak zorundadÄ±r.
- Sistem kaynaklÄ± pasiflikte Ã¼cretli veya sÃ¼reli avantajlarÄ±n kalan sÃ¼resi sessizce tÃ¼ketilemez; varsayÄ±lan gÃ¼venli davranÄ±ÅŸ sÃ¼renin dondurulmasÄ± ve sistem tekrar aÃ§Ä±ldÄ±ÄŸÄ±nda kaldÄ±ÄŸÄ± yerden devam etmesidir.
- Ã‡apraz modÃ¼l precedence kararlarÄ±nda sistem veya admin kaynaklÄ± emergency deny en yÃ¼ksek Ã¶nceliktedir; bunun altÄ±nda `access` deny kararÄ±, modÃ¼l iÃ§i allow veya paylaÅŸÄ±m sinyalini her zaman bastÄ±rmalÄ±dÄ±r.
- `user` modÃ¼lÃ¼ndeki global gÃ¶rÃ¼nÃ¼rlÃ¼k veya paylaÅŸÄ±m preference sinyali ilgili yÃ¼zey iÃ§in Ã¼st sÄ±nÄ±rÄ± tanÄ±mlar; `history` iÃ§indeki entry-level paylaÅŸÄ±m metadata'sÄ± bu Ã¼st sÄ±nÄ±rÄ± daraltabilir veya yalnÄ±zca izin verilen tavan iÃ§inde opt-in paylaÅŸÄ±m saÄŸlayabilir, ancak global deny kararÄ±nÄ± geniÅŸletemez.
- `social` modÃ¼lÃ¼nÃ¼n Ã¼rettiÄŸi block, privacy veya mute sinyalleri ham iliÅŸki verisidir; final allow veya deny kararÄ± `access` tarafÄ±ndan verilir. Block veya aÃ§Ä±k privacy deny sinyali final deny Ã¼retmelidir; mute sinyali ise aksi ayrÄ±ca dokÃ¼mante edilmedikÃ§e tek baÅŸÄ±na genel authorization deny sayÄ±lmamalÄ±dÄ±r.
- `moderation` gÃ¼nlÃ¼k scoped vaka akÄ±ÅŸÄ±nÄ±n sahibidir; ancak aynÄ± case Ã¼zerinde `admin` tarafÄ±ndan verilen override, reopen, freeze, reassignment veya final kararlar moderator aksiyonunun Ã¼zerinde precedence taÅŸÄ±r ve yeni bir handoff kaydÄ± oluÅŸmadan moderator tarafÄ±ndan bastÄ±rÄ±lamaz.
- `support` iÃ§indeki report kaydÄ± ile `moderation` case yaÅŸam dÃ¶ngÃ¼sÃ¼ aynÄ± ÅŸey sayÄ±lmamalÄ±dÄ±r; her report zorunlu olarak case aÃ§maz. Moderation incelemesi gerektiÄŸinde support kaydÄ± kaynak intake olarak kalÄ±r, moderation tarafÄ±nda ise buna baÄŸlÄ± ama ayrÄ± bir case lifecycle baÅŸlatÄ±lÄ±r.

## 12) GÃ¼venlik, Loglama ve Audit KurallarÄ±
- Loglar yapÄ±landÄ±rÄ±lmÄ±ÅŸ formatta Ã¼retilmelidir.
- JSON log formatÄ± tercih edilmelidir.
- Her request iÃ§in `request_id` Ã¼retilmeli veya taÅŸÄ±nmalÄ±dÄ±r.
- Hassas veri loglara yazÄ±lmamalÄ±dÄ±r.
- GÃ¼venlik olaylarÄ± gerektiÄŸinde ayrÄ± izlenebilmelidir.
- Audit log ile operasyonel log birbirine karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.
- YÃ¼ksek riskli iÅŸlemler izlenebilir ve aÃ§Ä±klanabilir olmalÄ±dÄ±r.
- Admin tarafÄ±ndan yapÄ±lan runtime ayar, modÃ¼l aÃ§ma-kapama, Ã¶zellik aÃ§ma-kapama ve eÅŸik gÃ¼ncelleme iÅŸlemleri actor, reason, scope, eski deÄŸer ve yeni deÄŸer ile audit kaydÄ± Ã¼retmelidir.
- Emergency deny veya kill switch iÅŸlemleri ayrÄ± kritik operasyon olayÄ± olarak izlenebilmelidir.
- Ã‡apraz modÃ¼l kritik iÅŸlemler gerektiÄŸinde ayrÄ± audit veya domain event kaydÄ± Ã¼retmelidir.

## 13) Config ve Ortam KurallarÄ±
- Config deÄŸerleri ortam deÄŸiÅŸkenlerinden okunmalÄ±dÄ±r.
- Config eriÅŸimi merkezi config yapÄ±sÄ±ndan yapÄ±lmalÄ±dÄ±r.
- ModÃ¼l bazlÄ± config alanlarÄ± namespace mantÄ±ÄŸÄ± ile adlandÄ±rÄ±lmalÄ±dÄ±r.
- Ã–rnek yaklaÅŸÄ±m: `AUTH_`, `MANGA_`, `PAYMENT_`, `NOTIFICATION_` gibi modÃ¼l prefixleri kullanÄ±labilir.
- YalnÄ±zca `.env.example` repoda tutulmalÄ±dÄ±r.
- GerÃ§ek secret deÄŸerleri repoya commit edilmemelidir.
- Ortam profilleri aÃ§Ä±kÃ§a ayrÄ±lmalÄ±dÄ±r: `local`, `test`, `staging`, `prod`.
- Main DB ve test DB config seviyesinde kesin olarak ayrÄ±lmalÄ±dÄ±r.
- Ortam baÄŸÄ±mlÄ± deÄŸerler kod iÃ§ine gÃ¶mÃ¼lmemelidir.
- Eksik veya geÃ§ersiz config deÄŸerleri kontrollÃ¼ ÅŸekilde doÄŸrulanmalÄ±dÄ±r.
- Ortam deÄŸiÅŸkenleri deploy veya Ã§alÄ±ÅŸma ortamÄ± seviyesindeki teknik config iÃ§indir; admin tarafÄ±ndan deÄŸiÅŸtirilebilen runtime ayarlar env config ile karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.
- Admin tarafÄ±ndan yÃ¶netilen runtime ayarlar merkezi ve kalÄ±cÄ± bir ayar deposunda tutulmalÄ±; uygulama yeniden build edilmeden deÄŸiÅŸtirilebilir olmalÄ±dÄ±r.
- Runtime ayar anahtarlarÄ± canonical namespace ile tanÄ±mlanmalÄ±dÄ±r; Ã¶rnek yaklaÅŸÄ±m `site.maintenance.enabled`, `auth.login.failed_attempt_limit_per_minute`, `comment.write.cooldown_seconds`, `feature.user.vip_benefits.enabled`.
- Boolean availability veya feature toggle anahtarlarÄ± mÃ¼mkÃ¼n olduÄŸunda `feature.<module>.<surface>.enabled` biÃ§imini kullanmalÄ±dÄ±r.
- EÅŸik, limit, cooldown veya davranÄ±ÅŸ deÄŸeri taÅŸÄ±yan anahtarlar mÃ¼mkÃ¼n olduÄŸunda `<module>.<surface>.<metric>` biÃ§imini kullanmalÄ±dÄ±r.
- Site geneli operasyon veya bakÄ±m anahtarlarÄ± mÃ¼mkÃ¼n olduÄŸunda `site.<surface>.<metric_or_flag>` biÃ§imini kullanmalÄ±dÄ±r.
- Audience, role, grup, kullanÄ±cÄ± veya resource bilgisi runtime key iÃ§ine gÃ¶mÃ¼lmemeli; bu bilgiler `scope_selector` ve `audience_selector` alanlarÄ±nda taÅŸÄ±nmalÄ±dÄ±r.
- Runtime key ve selector iÃ§indeki modÃ¼l adÄ± her zaman canonical leaf modÃ¼l adÄ± ile aynÄ± yazÄ±lmalÄ±dÄ±r.
- Runtime ayarlar en az ÅŸu canonical kategorilere geniÅŸleyebilir olmalÄ±dÄ±r: `site`, `communication`, `operations`, `security_auth`, `access_availability`, `content`, `reading`, `engagement`, `support`, `membership`, `social`, `gamification` ve `economy`.
- Runtime ayarlar tip, aralÄ±k, zorunluluk ve scope aÃ§Ä±sÄ±ndan doÄŸrulanmalÄ±; geÃ§ersiz ayar deÄŸiÅŸikliÄŸi sessizce kabul edilmemelidir.
- Teknik altyapÄ± config'i hiÃ§bir koÅŸulda admin runtime ayarÄ± sayÄ±lmamalÄ±dÄ±r; DB host veya pool, SMTP credential, queue DSN, object storage anahtarÄ±, secret key ve servis URL gibi deÄŸerler yalnÄ±zca env veya secret yÃ¶netimi ile taÅŸÄ±nmalÄ±dÄ±r.
- `site` kategorisi yalnÄ±zca kullanÄ±cÄ±ya gÃ¶rÃ¼nen genel Ã¼rÃ¼n davranÄ±ÅŸlarÄ±, genel site yÃ¼zeyleri ve public deneyim ayarlarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `communication` kategorisi iletiÅŸim sayfasÄ±, iletiÅŸim kanalÄ± gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼, destek giriÅŸ yÃ¼zeyi ve benzeri public iletiÅŸim verileri ile sÄ±nÄ±rlÄ± olmalÄ±dÄ±r; provider credential veya gizli anahtarlarÄ± kapsamaz.
- `operations` kategorisi dÃ¼ÅŸÃ¼k seviye altyapÄ± tuning'i iÃ§in deÄŸil, bakÄ±m modu, kayÄ±t aÃ§ma-kapama veya belirli runtime sÃ¼reÃ§leri durdurma gibi kontrollÃ¼ Ã¼rÃ¼n operasyon davranÄ±ÅŸlarÄ± iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `security_auth` kategorisi baÅŸarÄ±sÄ±z giriÅŸ limiti, cooldown, resend verification aralÄ±ÄŸÄ±, MFA zorunluluÄŸu ve benzeri auth gÃ¼venlik eÅŸikleri iÃ§in kullanÄ±lmalÄ±dÄ±r.
- `access_availability` kategorisi audience targeting, entitlement gating, feature availability ve kill switch karar yÃ¼zeyleri iÃ§in kullanÄ±lmalÄ±dÄ±r.
- Authorization, audience targeting, entitlement gating, feature availability ve kill switch kararlarÄ± `access` tarafÄ±ndan yorumlanmalÄ±dÄ±r.
- Auth gÃ¼venlik eÅŸikleri, yorum gÃ¶nderme aralÄ±ÄŸÄ±, attachment sÄ±nÄ±rlarÄ±, support intake davranÄ±ÅŸÄ± ve benzeri eriÅŸim dÄ±ÅŸÄ± runtime davranÄ±ÅŸlarÄ± ilgili modÃ¼lÃ¼n service katmanÄ±nda yorumlanmalÄ±dÄ±r; bunlar yalnÄ±zca eriÅŸim veya entitlement kararÄ± Ã¼rettiÄŸinde `access` ile entegre Ã§alÄ±ÅŸmalÄ±dÄ±r.
- Site iÃ§eriÄŸi, iletiÅŸim iÃ§eriÄŸi ve eriÅŸim dÄ±ÅŸÄ± runtime ayarlar `access` Ã¼zerinden Ã§Ã¶zÃ¼lmemelidir.
- Her runtime ayarÄ± iÃ§in en az ÅŸu metadata alanlarÄ± tanÄ±mlanmalÄ±dÄ±r: `key`, `description`, `category`, `owner_module`, `consumer_layer`, `value_type`, `default_value`, `allowed_range_or_enum`, `scope_kind`, `scope_selector`, `audience_kind`, `audience_selector`, `sensitive`, `apply_mode`, `cache_strategy`, `schedule_support`, `audit_required`, `affected_surfaces` ve gerektiÄŸinde `disabled_behavior` ile `error_response_policy`.
- Selector gerektirmeyen `site` veya `all` gibi kayÄ±tlarda `scope_selector` ve `audience_selector` iÃ§in aÃ§Ä±k boÅŸ deÄŸer standardÄ± kullanÄ±lmalÄ±dÄ±r.
- `scope_selector` iÃ§in canonical yaklaÅŸÄ±m en az `-`, `<module>`, `<module>.<surface>`, `<module>.<surface>.<subsurface>` ve gerektiÄŸinde `resource:<module>:<resource_kind>:<identifier>` biÃ§imlerini desteklemelidir.
- `audience_selector` iÃ§in canonical yaklaÅŸÄ±m en az `-`, `role:<name>`, `group:<name>` ve `user:<id>` biÃ§imlerini desteklemelidir.
- `apply_mode` en az `immediate`, `cache_refresh` ve `scheduled` deÄŸerlerini destekleyecek ÅŸekilde tasarlanmalÄ±dÄ±r.
- `cache_strategy` en az `none`, `ttl` ve `manual_invalidate` gibi aÃ§Ä±k stratejilerle tanÄ±mlanmalÄ±dÄ±r.
- `schedule_support` en az `none`, `start_at` ve `time_window` gibi aÃ§Ä±k planlama modlarÄ± ile tanÄ±mlanmalÄ±dÄ±r.
- Runtime ayar envanterinin canonical kayÄ±t dosyasÄ± `docs/settings/index.md` olmalÄ±; yeni ayar, toggle, kill switch veya limit eklendiÄŸinde aynÄ± deÄŸiÅŸiklik setinde bu dosya gÃ¼ncellenmelidir.
- Ãœcretli veya sÃ¼reli avantajÄ± etkileyen ayarlarda metadata'ya ek olarak `entitlement_impact_policy` zorunlu olmalÄ±dÄ±r.

## 14) Migration ve Ã‡alÄ±ÅŸma OrtamÄ± KurallarÄ±
- Migration yÃ¶netiminde `golang-migrate` kullanÄ±lmalÄ±dÄ±r.
- Her migration iÃ§in `up` ve `down` script zorunludur.
- Backend migration dosyalarÄ± yalnÄ±zca `apps/api/migrations/` altÄ±nda tutulmalÄ±dÄ±r.
- Migration dosyalarÄ± standart isimlendirme ile oluÅŸturulmalÄ±dÄ±r.
- Ã‡ok modÃ¼llÃ¼ yapÄ±da migration isimlerinde modÃ¼l veya alan prefixi kullanÄ±lmalÄ±dÄ±r.
- Ã–rnek yaklaÅŸÄ±m: `YYYYMMDDHHMM_auth_create_sessions.up.sql` veya `YYYYMMDDHHMM_content_create_manga.up.sql`.
- Åema deÄŸiÅŸiklikleri migration olmadan uygulanmamalÄ±dÄ±r.
- Seed ve migration sÃ¼reÃ§leri birbirine karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.
- Backend iÃ§in gerekli uygulama Ã§alÄ±ÅŸma dosyalarÄ± `apps/api/` altÄ±nda, repo seviyesi deploy dosyalarÄ± ise `deploy/` altÄ±nda tutulmalÄ±dÄ±r.
- Dockerfile uygulama kÃ¶kÃ¼nde ilgili app altÄ±nda yer almalÄ±; compose ve benzeri Ã§ok servisli Ã§alÄ±ÅŸma dosyalarÄ± repo seviyesinde merkezi olarak yÃ¶netilmelidir.
- Proje Docker iÃ§inde build alabilmeli ve Ã§alÄ±ÅŸabilmelidir.
- Ã‡alÄ±ÅŸma iÃ§in gereken servisler tekrarlanabilir ÅŸekilde tanÄ±mlanmalÄ±dÄ±r.
- Local, test ve benzeri ortamlar mÃ¼mkÃ¼n olduÄŸunca tutarlÄ± olmalÄ±dÄ±r.

## 15) Git, PR ve Kod Ä°nceleme KurallarÄ±

- Git reposu: `https://github.com/Tokuchi61/Manga`
- VarsayÄ±lan remote `origin` olmalÄ±dÄ±r.
- Push iÅŸlemleri yalnÄ±zca bu repoya yapÄ±lmalÄ±dÄ±r.
- Onay olmadan farklÄ± remote eklenmemeli ve farklÄ± repolara push yapÄ±lmamalÄ±dÄ±r.

- Branch modeli `main + feature/* + hotfix/*` olmalÄ±dÄ±r.
- `main` daima deploy edilebilir durumda kalmalÄ±dÄ±r.
- DoÄŸrudan `main` branch'e push yapÄ±lmamalÄ±dÄ±r.
- TÃ¼m deÄŸiÅŸiklikler PR Ã¼zerinden ilerlemelidir.

- Branch adlarÄ± kÄ±sa, aÃ§Ä±klayÄ±cÄ± ve konu odaklÄ± olmalÄ±dÄ±r.
- Feature branch formatÄ±: `feature/<konu>`
- Hotfix branch formatÄ±: `hotfix/<konu>`

- Commit'ler kÃ¼Ã§Ã¼k, anlamlÄ± ve geri alÄ±nabilir olmalÄ±dÄ±r.
- Commit mesajlarÄ± Conventional Commits standardÄ±na uygun olmalÄ±dÄ±r.
- Tek commit iÃ§inde birden fazla baÄŸÄ±msÄ±z konu birleÅŸtirilmemelidir.

- Her PR tek bir konuya odaklÄ± olmalÄ±dÄ±r.
- BÃ¼yÃ¼k geliÅŸtirmeler kÃ¼Ã§Ã¼k ve incelenebilir PR'lara bÃ¶lÃ¼nmelidir.
- AltyapÄ±, modÃ¼l geliÅŸtirmesi, refactor ve dokÃ¼man gÃ¼ncellemeleri mÃ¼mkÃ¼n olduÄŸunda mantÄ±klÄ± parÃ§alara ayrÄ±lmalÄ±dÄ±r.

- PR aÃ§Ä±klamasÄ± en az ÅŸu bÃ¶lÃ¼mleri iÃ§ermelidir:
  - ne deÄŸiÅŸti
  - neden deÄŸiÅŸti
  - nasÄ±l test edildi

- PR aÃ§Ä±lmadan Ã¶nce en az ÅŸu kontroller yapÄ±lmÄ±ÅŸ olmalÄ±dÄ±r:
  - ilgili testler Ã§alÄ±ÅŸtÄ±rÄ±lmÄ±ÅŸ olmalÄ±
  - lint/format kontrolleri geÃ§miÅŸ olmalÄ±
  - yeni migration varsa kontrol edilmiÅŸ olmalÄ±
  - ilgili dokÃ¼man gÃ¼ncellemeleri eklenmiÅŸ olmalÄ±
  - gereksiz debug/log/yorum satÄ±rlarÄ± temizlenmiÅŸ olmalÄ±

- VeritabanÄ± ÅŸemasÄ±nÄ± etkileyen deÄŸiÅŸikliklerde migration zorunludur.
- Migration iÃ§eren PR'larda ilgili model, repository, servis ve test gÃ¼ncellemeleri birlikte deÄŸerlendirilmelidir.
- Geri alma etkisi yÃ¼ksek migration'lar PR aÃ§Ä±klamasÄ±nda ayrÄ±ca belirtilmelidir.

- Kod deÄŸiÅŸikliÄŸi mimari, modÃ¼l sÄ±nÄ±rÄ±, veri sahipliÄŸi, ayar, event akÄ±ÅŸÄ± veya eriÅŸim davranÄ±ÅŸÄ±nÄ± etkiliyorsa ilgili dokÃ¼manlar aynÄ± PR iÃ§inde gÃ¼ncellenmelidir.
- Ä°lgili dokÃ¼manlar gÃ¼ncel deÄŸilse PR merge edilmemelidir.

- Merge Ã¶ncesi CI sonucu baÅŸarÄ±lÄ± olmalÄ±dÄ±r.
- PaylaÅŸÄ±lan branch'lerde force push kullanÄ±lmamalÄ±dÄ±r.
- Commit geÃ§miÅŸi inceleme sÃ¼recini bozacak ÅŸekilde yeniden yazÄ±lmamalÄ±dÄ±r.
- VarsayÄ±lan merge yÃ¶ntemi ekip standardÄ±na gÃ¶re belirlenmeli; aksi belirtilmedikÃ§e squash merge tercih edilmelidir.

- Ajanlar doÄŸrudan `main` branch Ã¼zerinde Ã§alÄ±ÅŸmamalÄ±dÄ±r.
- Her gÃ¶rev iÃ§in uygun bir `feature/*` veya `hotfix/*` branch aÃ§Ä±lmalÄ±dÄ±r.
- Her aÅŸama sonunda deÄŸiÅŸiklikler ilgili branch'e push edilmeli ve PR aÃ§Ä±lmaya hazÄ±r halde bÄ±rakÄ±lmalÄ±dÄ±r.
- Remote veya branch belirsizse ajan varsayÄ±m yapmamalÄ±, mevcut git yapÄ±landÄ±rmasÄ±nÄ± korumalÄ±dÄ±r.

## 16) Versiyonlama KurallarÄ±
- Proje versiyonlamasÄ±nda SemVer (`MAJOR.MINOR.PATCH`) standardÄ± kullanÄ±lmalÄ±dÄ±r.
- Versiyon formatÄ± Ã¼retim ve kalÄ±cÄ± release'ler iÃ§in yalnÄ±zca `X.Y.Z` biÃ§iminde olmalÄ±dÄ±r.
- GeliÅŸtirme ve release adayÄ± sÃ¼reÃ§lerinde gerekirse pre-release etiketleri kullanÄ±labilir:
  - `X.Y.Z-alpha.N`
  - `X.Y.Z-beta.N`
  - `X.Y.Z-rc.N`
- Build metadata gerekiyorsa SemVer ile uyumlu `+build` eki kullanÄ±labilir; ancak asÄ±l Ã¼rÃ¼n versiyonu bunun Ã¶ncesindeki canonical deÄŸerdir.
- UygulamanÄ±n Ã§alÄ±ÅŸÄ±rken gÃ¶sterdiÄŸi canonical versiyon tek kaynak Ã¼zerinden yÃ¶netilmelidir.
- Runtime tarafÄ±nda versiyon bilgisi `APP_VERSION` environment deÄŸiÅŸkeni Ã¼zerinden okunmalÄ±dÄ±r.
- Versiyon bilgisi kod iÃ§ine sabit string olarak gÃ¶mÃ¼lmemelidir.
- Repo iÃ§indeki dokÃ¼man, release kaydÄ±, tag ve daÄŸÄ±tÄ±m Ã§Ä±ktÄ±larÄ± aynÄ± canonical versiyon ile hizalÄ± olmalÄ±dÄ±r.
- AynÄ± iÃ§erik iÃ§in birden fazla canonical versiyon adÄ± Ã¼retilmemelidir.
- Her release tek bir versiyon numarasÄ±na sahip olmalÄ±dÄ±r.
- YayÄ±nlanmÄ±ÅŸ bir versiyon sonradan sessizce deÄŸiÅŸtirilmemelidir; yeni deÄŸiÅŸiklik gerekiyorsa yeni versiyon Ã§Ä±karÄ±lmalÄ±dÄ±r.
- Geri alma gerekiyorsa eski versiyonu sessizce oynatmak yerine yeni bir dÃ¼zeltme versiyonu Ã¼retilmelidir.
- Versiyon artÄ±rÄ±mÄ± gerektiren deÄŸiÅŸiklikler release Ã¶ncesinde netleÅŸtirilmelidir.
- Versiyon artÄ±rÄ±mÄ± ÅŸu kurallara gÃ¶re yapÄ±lmalÄ±dÄ±r:
  - `MAJOR`: Geriye dÃ¶nÃ¼k uyumsuz API deÄŸiÅŸikliÄŸi, veri modeli kÄ±rÄ±lmasÄ±, davranÄ±ÅŸ deÄŸiÅŸikliÄŸi, kaldÄ±rÄ±lan alan/endpoint, zorunlu migration uyumsuzluÄŸu, mevcut entegrasyonlarÄ± bozan mimari deÄŸiÅŸiklik.
  - `MINOR`: Geriye dÃ¶nÃ¼k uyumlu yeni Ã¶zellik, yeni endpoint, yeni modÃ¼l, yeni opsiyonel alan, mevcut davranÄ±ÅŸÄ± kÄ±rmadan yapÄ±lan anlamlÄ± kapasite artÄ±ÅŸÄ±.
  - `PATCH`: Geriye dÃ¶nÃ¼k uyumlu bugfix, gÃ¼venlik dÃ¼zeltmesi, kÃ¼Ã§Ã¼k performans iyileÅŸtirmesi, davranÄ±ÅŸÄ± kÄ±rmayan iÃ§ dÃ¼zeltme.
- YalnÄ±zca dokÃ¼mantasyon, yorum, metin dÃ¼zeltmesi veya release Ã§Ä±ktÄ±sÄ±nÄ± etkilemeyen iÃ§ temizlikler tek baÅŸÄ±na versiyon artÄ±rmak zorunda deÄŸildir.
- Ancak dokÃ¼man deÄŸiÅŸikliÄŸi mevcut release'in kullanÄ±mÄ±nÄ±, kurulumunu, entegrasyonunu veya gÃ¼venliÄŸini fiilen etkiliyorsa uygun versiyon artÄ±rÄ±mÄ± deÄŸerlendirilmelidir.
- VeritabanÄ± migration iÃ§eren her deÄŸiÅŸiklik iÃ§in versiyon etkisi ayrÄ±ca deÄŸerlendirilmelidir.
- Geriye dÃ¶nÃ¼k uyumsuz migration deÄŸiÅŸiklikleri `MAJOR`, uyumlu ÅŸema geniÅŸletmeleri en az `MINOR` olarak ele alÄ±nmalÄ±dÄ±r.
- GÃ¼venlik aÃ§Ä±ÄŸÄ± kapatan deÄŸiÅŸiklikler varsayÄ±lan olarak en az `PATCH` artÄ±ÅŸÄ± ile yayÄ±nlanmalÄ±dÄ±r.
- Public API sÃ¶zleÅŸmesini etkileyen her deÄŸiÅŸiklikte versiyon etkisi aÃ§Ä±kÃ§a belirtilmelidir.
- PR aÃ§Ä±klamalarÄ±nda gerekiyorsa hedef versiyon veya beklenen bump tipi belirtilmelidir.
- Release hazÄ±rlÄ±ÄŸÄ±nda en az ÅŸu alanlar birlikte gÃ¼ncellenmelidir:
  - `APP_VERSION`
  - `docs/changelog.md`
  - gerekiyorsa `README.md`
  - gerekiyorsa kurulum, migration veya breaking change notlarÄ±
- `docs/changelog.md` iÃ§inde her release iÃ§in en az ÅŸu bilgiler yer almalÄ±dÄ±r:
  - versiyon
  - tarih
  - deÄŸiÅŸiklik Ã¶zeti
  - etkilenen modÃ¼ller
  - breaking change bilgisi varsa aÃ§Ä±k not
  - migration etkisi varsa aÃ§Ä±k not
- `docs/changelog.md` yalnÄ±zca final release baÅŸlÄ±klarÄ±ndan ibaret olmamalÄ±; modÃ¼l, feature, hotfix, fix, refactor, security ve operasyonel dÃ¼zeltmeler uygun release girdisi altÄ±nda izlenebilir ÅŸekilde gruplanmalÄ±dÄ±r.
- Release giriÅŸlerinde mÃ¼mkÃ¼n olduÄŸunda ÅŸu alt baÅŸlÄ±klar kullanÄ±lmalÄ±dÄ±r:
  - `Added`
  - `Changed`
  - `Fixed`
  - `Removed`
  - `Deprecated`
  - `Security`
  - `Docs`
- Hotfix kayÄ±tlarÄ±nda en az etkilenen alan, sorunun kÄ±sa Ã¶zeti ve dÃ¼zeltilen kapsam belirtilmelidir.
- API, migration, config, access veya operasyonel davranÄ±ÅŸÄ± etkileyen deÄŸiÅŸikliklerde gerekli kullanÄ±cÄ± veya geliÅŸtirici aksiyonlarÄ± changelog iÃ§inde aÃ§Ä±kÃ§a yazÄ±lmalÄ±dÄ±r.
- Release Ã§Ä±kmadan Ã¶nce changelog taslak girdileri hazÄ±rlanabilir; ancak yayÄ±n anÄ±nda hepsi canonical versiyon baÅŸlÄ±ÄŸÄ± altÄ±nda birleÅŸtirilmelidir.
- Breaking change iÃ§eren release'lerde upgrade notu zorunlu olmalÄ±dÄ±r.
- Release candidate kullanÄ±lÄ±yorsa production'a Ã§Ä±kmadan Ã¶nce final versiyon ayrÄ±ca Ã¼retilmelidir.
- Git tag standardÄ± canonical versiyon ile uyumlu olmalÄ± ve `vX.Y.Z` biÃ§iminde aÃ§Ä±lmalÄ±dÄ±r.
- Pre-release tag'leri gerekiyorsa `vX.Y.Z-rc.N` benzeri biÃ§imde oluÅŸturulmalÄ±dÄ±r.
- Tag, changelog ve daÄŸÄ±tÄ±lan artifact versiyonu birbiriyle Ã§eliÅŸmemelidir.
- Release alÄ±nmadan Ã¶nce build, test ve kritik doÄŸrulama adÄ±mlarÄ± baÅŸarÄ±lÄ± olmalÄ±dÄ±r.
- Versiyon artÄ±ÅŸÄ± yapÄ±lan deÄŸiÅŸikliklerde rollback, migration ve uyumluluk etkisi en az bir kez gÃ¶zden geÃ§irilmelidir.

## 17) Test ve DoÄŸrulama KurallarÄ±
- Yeni eklenen veya gÃ¼ncellenen yapÄ± test edilebilir olmalÄ±dÄ±r.
- Testler ana veritabanÄ±nÄ± asla kullanmamalÄ±, yalnÄ±zca test veritabanÄ± ile Ã§alÄ±ÅŸmalÄ±dÄ±r.
- Unit test Ã¶nceliÄŸi iÅŸ kurallarÄ± ve veri doÄŸrulama katmanlarÄ±dÄ±r.
- Integration test Ã¶nceliÄŸi veri eriÅŸimi, modÃ¼l kontratlarÄ± ve kritik HTTP akÄ±ÅŸlarÄ±dÄ±r.
- Yeni endpoint eklendiÄŸinde en az bir baÅŸarÄ±lÄ± ve bir baÅŸarÄ±sÄ±z senaryo test edilmelidir.
- Veri sÄ±zÄ±ntÄ±sÄ± ve yetkisiz eriÅŸim senaryolarÄ± test edilmelidir.
- Ã‡apraz modÃ¼l entegrasyonlarÄ± varsa contract veya integration test ile doÄŸrulanmalÄ±dÄ±r.
- `go test ./...` temel doÄŸrulama kontrolÃ¼ olarak kabul edilmelidir.
- BaÅŸlangÄ±Ã§ coverage hedefi minimum `%60` olarak Ã¶nerilir.

## 18) Kapsam ve DokÃ¼mantasyon YÃ¶netimi
- Kapsam dÄ±ÅŸÄ± feature'lar varsayÄ±lÄ±p eklenmemelidir.
- ÃœrÃ¼n veya mimari kapsamÄ±nÄ± deÄŸiÅŸtiren kararlar Ã¶nce dokÃ¼mana yansÄ±tÄ±lmalÄ±dÄ±r.
- AynÄ± anda Ã§ok fazla sorumluluk aÃ§mak yerine en kÃ¼Ã§Ã¼k sÃ¼rdÃ¼rÃ¼lebilir kapsam seÃ§ilmelidir.
- `RULES.md`, `ROADMAP.md` ve gerekli diÄŸer dokÃ¼manlar birlikte gÃ¼ncellenmelidir.
- `docs/settings/index.md`, `docs/shared/target-types.md` ve benzeri canonical kayÄ±t dosyalarÄ± etkilendiklerinde aynÄ± deÄŸiÅŸiklik setinde gÃ¼ncellenmelidir.
- AynÄ± bilgi birden fazla dosyada Ã§eliÅŸkili ÅŸekilde bÄ±rakÄ±lmamalÄ±dÄ±r.
- `README.md` proje giriÅŸi ve ana dokÃ¼man baÄŸlantÄ±larÄ± iÃ§in gÃ¼ncel tutulmalÄ±dÄ±r.
- SÃ¼rÃ¼m bazlÄ± deÄŸiÅŸiklikler `docs/changelog.md` iÃ§inde kayÄ±t altÄ±na alÄ±nmalÄ±dÄ±r.
- Bilinen sorunlar ve teknik borÃ§lar `docs/issues.md` iÃ§inde tutulmalÄ±dÄ±r.
- Repo kÃ¶k dizin yapÄ±sÄ±nÄ±, uygulama yerleÅŸimini veya deployment klasÃ¶r yapÄ±sÄ±nÄ± etkileyen deÄŸiÅŸikliklerde `RULES.md`, `README.md` ve `SETUP.md` aynÄ± deÄŸiÅŸiklik setinde birlikte gÃ¼ncellenmelidir.
- Her yeni leaf modÃ¼l iÃ§in en az bir modÃ¼l dokÃ¼manÄ± aÃ§Ä±lmalÄ±dÄ±r.
- Yeni runtime ayar, feature toggle, kill switch veya oran limiti eklendiÄŸinde ilgili key, scope, scope selector, audience, audience selector, varsayÄ±lan deÄŸer, disabled behavior varsa bunun tipi, error response policy varsa bunun tipi ve etkilediÄŸi modÃ¼ller dokÃ¼mana yansÄ±tÄ±lmalÄ±; `docs/settings/index.md` aynÄ± deÄŸiÅŸiklikte gÃ¼ncellenmelidir.
- ModÃ¼l dokÃ¼manlarÄ±nda en az ÅŸu baÅŸlÄ±klar yer almalÄ±dÄ±r:
  - amaÃ§
  - sorumluluk alanÄ±
  - veri sahipliÄŸi
  - access kontratÄ±
  - API veya event sÄ±nÄ±rÄ±
  - baÄŸÄ±mlÄ±lÄ±klar
  - state yapÄ±sÄ±
  - test notlarÄ±

## 19) Son Mimari Ä°lke
- Kimlik doÄŸrulama, veri sahipliÄŸi, eriÅŸim kararÄ± ve yÃ¶netim akÄ±ÅŸlarÄ± birbirine karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.
- Veri sahipliÄŸi ilgili alanda, authorization kararÄ± merkezi yapÄ±da kalmalÄ±dÄ±r.
- Ä°Ã§erik, topluluk, operasyon, ticaret, etkileÅŸim ve altyapÄ± sorumluluklarÄ± ayrÄ±ÅŸmÄ±ÅŸ kalmalÄ±dÄ±r.
- Yeni Ã¶zellik eklemek bu ayrÄ±mÄ± bozma gerekÃ§esi olamaz.

## 20) ModÃ¼l Referans YapÄ±sÄ±
- ModÃ¼l kurallarÄ± dokÃ¼man iÃ§inde tek bir ana bÃ¶lÃ¼m altÄ±nda toplanmalÄ±dÄ±r.
- Bu yapÄ± iÃ§in ana modÃ¼l bÃ¶lÃ¼mÃ¼ `21) ModÃ¼l Genel KurallarÄ±` baÅŸlÄ±ÄŸÄ± olmalÄ±dÄ±r.
- Her leaf modÃ¼l bu ana bÃ¶lÃ¼m altÄ±nda `21.1`, `21.2`, `21.3` ÅŸeklinde ayrÄ± alt baÅŸlÄ±k olarak eklenmelidir.
- Yeni modÃ¼l eklendikÃ§e numaralandÄ±rma aynÄ± ana bÃ¶lÃ¼m altÄ±nda devam etmelidir: `21.5`, `21.6`, `21.7` ... `21.20`.
- ModÃ¼l sayÄ±sÄ± artsa bile her yeni modÃ¼l iÃ§in yeni bir Ã¼st seviye bÃ¶lÃ¼m aÃ§Ä±lmamalÄ±; modÃ¼l kurallarÄ± `21.x` formatÄ±nda sÃ¼rdÃ¼rÃ¼lmelidir.
- Yeni modÃ¼ller mevcut modÃ¼l maddelerinin iÃ§ine gÃ¶mÃ¼lmemeli; her yeni leaf modÃ¼l ayrÄ± alt baÅŸlÄ±k olarak eklenmelidir.
- Roadmap tarafÄ±nda da aynÄ± yaklaÅŸÄ±m korunmalÄ±; her yeni leaf modÃ¼l mevcut aÅŸamalarÄ±n iÃ§ine sÄ±kÄ±ÅŸtÄ±rÄ±lmadan ayrÄ± bir aÅŸama olarak eklenmelidir.

## 21) ModÃ¼l Genel KurallarÄ±
- Bu bÃ¶lÃ¼m, proje iÃ§in kabul edilmiÅŸ leaf modÃ¼l referanslarÄ±nÄ± merkezi olarak kaydeder.
- Bu bÃ¶lÃ¼mdeki kayÄ±tlar yÃ¼ksek seviye modÃ¼l sÄ±nÄ±rÄ± ve sahiplik Ã¶zeti taÅŸÄ±r; detaylÄ± modÃ¼l tasarÄ±mÄ± ilgili modÃ¼l dokÃ¼manÄ±nda yer almalÄ±dÄ±r.
- Her modÃ¼l alt baÅŸlÄ±ÄŸÄ± canonical modÃ¼l adÄ±nÄ±, temel sorumluluk alanÄ±nÄ±, veri sahipliÄŸini ve ana referans dokÃ¼manÄ±nÄ± aÃ§Ä±kÃ§a gÃ¶stermelidir.
- Her modÃ¼l, ihtiyaÃ§ duyduÄŸu durumda admin tarafÄ±ndan yÃ¶netilen runtime ayarlar, modÃ¼l aÃ§ma-kapama ve Ã¶zellik aÃ§ma-kapama yÃ¼zeyleri ile uyumlu Ã§alÄ±ÅŸacak ÅŸekilde tasarlanmalÄ±dÄ±r; veri sahipliÄŸi ilgili modÃ¼lde, authorization veya availability yorumlama `access`te, eriÅŸim dÄ±ÅŸÄ± runtime davranÄ±ÅŸ yorumlama ilgili modÃ¼l service katmanÄ±nda, operasyon yÃ¶netimi `admin`de kalmalÄ±dÄ±r.

### 21.1) `auth` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `auth` olmalÄ±dÄ±r.
- `auth` modÃ¼lÃ¼ kimlik doÄŸrulama, oturum gÃ¼venliÄŸi ve hesap giriÅŸ gÃ¼venliÄŸi alanlarÄ±nÄ±n sahibidir.
- `auth` modÃ¼lÃ¼ kayÄ±t, giriÅŸ, Ã§Ä±kÄ±ÅŸ, token, session, email verification, password reset, password change, login gÃ¼venliÄŸi ve auth audit akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `auth` modÃ¼lÃ¼ kullanÄ±cÄ± profil verisi, kullanÄ±cÄ± tercihleri, Ã¼yelik avantajlarÄ±, rol/permission yÃ¶netimi veya authorization kararÄ± Ã¼retmemelidir.
- `auth` veri sahipliÄŸi; credential benzeri kimlik bilgileri, auth session kayÄ±tlarÄ±, token yaÅŸam dÃ¶ngÃ¼sÃ¼ ve auth gÃ¼venlik olaylarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `auth` modÃ¼lÃ¼nÃ¼n public surface'i, kimlik doÄŸrulama akÄ±ÅŸlarÄ± ve gÃ¼venli oturum yÃ¶netimi iÃ§in gerekli contract yÃ¼zeyi ile sÄ±nÄ±rlÄ± olmalÄ±dÄ±r.
- `auth` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/auth.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - refresh token Ã¼retimi ve yenileme akÄ±ÅŸÄ±
  - session listing, session revoke ve logout all akÄ±ÅŸlarÄ±
  - forgot password, reset password ve change password ayrÄ±mÄ±
  - email verification token, resend verification ve verification gÃ¼venliÄŸi
  - login rate limit, failed login limit ve cooldown veya temporary lock davranÄ±ÅŸÄ±
  - device, IP, son giriÅŸ ve ÅŸÃ¼pheli giriÅŸ takibi
  - email doÄŸrulanmamÄ±ÅŸ, suspend edilmiÅŸ veya banlÄ± kullanÄ±cÄ± iÃ§in auth kontrol davranÄ±ÅŸÄ±
  - baÅŸarÄ±sÄ±z giriÅŸ limiti, cooldown sÃ¼resi, resend verification aralÄ±ÄŸÄ± ve gerektiÄŸinde MFA zorunluluÄŸu gibi gÃ¼venlik eÅŸiklerinin admin tarafÄ±ndan yÃ¶netilen runtime ayarlar ile kontrol edilebilmesi
  - Ã§ok faktÃ¶rlÃ¼ doÄŸrulama, trusted device ve risk score temelli giriÅŸ deÄŸerlendirmesi
  - ÅŸÃ¼pheli giriÅŸ durumlarÄ±nda ek challenge veya doÄŸrulama adÄ±mÄ±

### 21.2) `user` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `user` olmalÄ±dÄ±r.
- `user` modÃ¼lÃ¼ kullanÄ±cÄ± hesabÄ±, profil, hesap durumu, tercih, gÃ¶rÃ¼nÃ¼m ve Ã¼yelik verisinin sahibidir.
- `user` modÃ¼lÃ¼ public/private profil ayrÄ±mÄ±, kullanÄ±cÄ± arama veya listeleme, hesap durumu alanlarÄ± ve Ã¼yelik verisi gibi kullanÄ±cÄ± merkezli alanlarÄ± taÅŸÄ±malÄ±dÄ±r.
- `user` modÃ¼lÃ¼ kimlik doÄŸrulama akÄ±ÅŸlarÄ±nÄ±, authorization kararlarÄ±nÄ±, role/permission yÃ¶netimini veya admin operasyon kararlarÄ±nÄ± sahiplenmemelidir.
- `user` veri sahipliÄŸi; kullanÄ±cÄ± kimliÄŸi ile iliÅŸkili profil alanlarÄ±, hesap durum alanlarÄ±, tercih alanlarÄ± ve Ã¼yelik durum verileri ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `user` modÃ¼lÃ¼ veri taÅŸÄ±r ve yayÄ±nlar; eriÅŸim kararÄ± veya feature access kararÄ± Ã¼retmez.
- `user` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/user.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - kullanÄ±cÄ± arama, kullanÄ±cÄ± listeleme ve reserved username kontrolÃ¼
  - email change, preferences update ve privacy alanlarÄ±
  - aktif veya pasif, suspend, ban, soft delete ve restore hesap durumlarÄ±
  - display name, avatar, banner, bio ve gÃ¶rÃ¼nÃ¼m alanlarÄ±
  - visibility preset yapÄ±larÄ± ve profil gÃ¶rÃ¼nÃ¼rlÃ¼k ÅŸablonlarÄ±
  - Ã¼yelik ve VIP veri alanlarÄ±
  - EXP, level, level progress ve kullanÄ±cÄ±ya ait oyunlaÅŸtÄ±rma Ã¶zeti
  - VIP rozet ve envanterden seÃ§ilen profil efekti veya nameplate referanslarÄ± gibi profil gÃ¶rÃ¼nÃ¼m alanlarÄ±
  - profil gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼, VIP gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ve benzeri kullanÄ±cÄ± alt yÃ¼zeylerinin admin tarafÄ±ndan yÃ¶netilen feature toggle veya runtime ayarlar ile kontrollÃ¼ biÃ§imde aÃ§Ä±lÄ±p kapatÄ±labilmesi
  - sistem kaynaklÄ± VIP global pasiflikte mevcut kullanÄ±cÄ±nÄ±n kalan VIP sÃ¼resinin dondurulmasÄ± ve tekrar aÃ§Ä±ldÄ±ÄŸÄ±nda kaldÄ±ÄŸÄ± yerden devam etmesi
  - kullanÄ±cÄ±ya ait genel preference sinyalleri; ancak detaylÄ± bildirim tercihleri `notification`, sosyal iliÅŸki blok veya mute listeleri ise `social` modÃ¼lÃ¼nde sahiplenilmelidir
  - reading history veya library gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼nÃ¼ etkileyen kullanÄ±cÄ± preference sinyalleri; ancak continue reading, bookmark veya okuma kayÄ±tlarÄ±nÄ±n kendisi `history` modÃ¼lÃ¼nde sahiplenilmelidir
  - reading activity veya public library gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼nde `user` modÃ¼lÃ¼ndeki global preference sinyalinin Ã¼st sÄ±nÄ±r, `history` iÃ§indeki entry-level share metadata'sÄ±nÄ±n ise bu Ã¼st sÄ±nÄ±r iÃ§inde Ã§alÄ±ÅŸan alt karar olmasÄ±
  - hesap dÄ±ÅŸa aktarma veya hesap verisi export yÃ¼zeyi
  - profil deÄŸiÅŸim geÃ§miÅŸi veya profile history gÃ¶rÃ¼nÃ¼mÃ¼

### 21.3) `access` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `access` olmalÄ±dÄ±r.
- `access` modÃ¼lÃ¼ sistemdeki tÃ¼m authorization, policy ve eriÅŸim kararlarÄ±nÄ±n merkezi sahibi olmalÄ±dÄ±r.
- `access` modÃ¼lÃ¼ role, permission, policy, ownership, guest/authenticated/vip kararlarÄ±, endpoint guard ve modÃ¼l bazlÄ± authorization contract alanlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `access` modÃ¼lÃ¼ kullanÄ±cÄ± profili, credential verisi, iÃ§erik verisi veya yÃ¶netim use-case verisi taÅŸÄ±mamalÄ±dÄ±r.
- `access` veri sahipliÄŸi; authorization sÃ¶zlÃ¼ÄŸÃ¼, rol-permission iliÅŸkileri, policy yorumlarÄ± ve eriÅŸim kararÄ±na temel olan kurallarla sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `access` modÃ¼lÃ¼ kimlik doÄŸrulama yapmaz; `auth` tarafÄ±ndan doÄŸrulanan kimliÄŸi ve `user` tarafÄ±ndan taÅŸÄ±nan kullanÄ±cÄ± verisini kullanarak karar Ã¼retir.
- `access` modÃ¼lÃ¼ yalnÄ±zca authorization, audience targeting, entitlement gating, feature availability ve kill switch kararlarÄ±nÄ± yorumlamalÄ±dÄ±r; `site` veya `communication` kategorisindeki Ã¼rÃ¼n ayarlarÄ± ile eriÅŸim dÄ±ÅŸÄ± iÅŸ kuralÄ± eÅŸikleri `access` iÃ§inde Ã§Ã¶zÃ¼lmemelidir.
- `access` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/access.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - public, guest, authenticated, vip, early access ve gerektiÄŸinde restricted kararlarÄ±
  - role CRUD, permission CRUD, user-role ve role-permission iliÅŸkileri
  - default role atama, Ã§oklu rol desteÄŸi ve rol Ã¶nceliÄŸi kurallarÄ±
  - own veya any ayrÄ±mÄ±, ownership ve resource access kurallarÄ±
  - endpoint guard, use-case guard ve admin panel gÃ¶rÃ¼nÃ¼rlÃ¼k kararlarÄ±
  - super admin bypass, moderator veya admin override kurallarÄ±
  - chapter iÃ§in minimum `authenticated` okuma kapÄ±sÄ±, misafir kullanÄ±cÄ±nÄ±n site ve manga detayÄ±na eriÅŸebilmesi, ancak chapter okuma ve yorum yazma gibi akÄ±ÅŸlarda kÄ±sÄ±tlÄ± kalmasÄ±
  - VIP Ã¶zel bÃ¶lÃ¼m eriÅŸimi ve belirli bÃ¶lÃ¼mler iÃ§in VIP erken eriÅŸim kararlarÄ±nÄ±n merkezi yÃ¶netimi
  - reklam gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼, VIP reklamsÄ±z deneyim, kozmetik gÃ¶rÃ¼nÃ¼rlÃ¼k ve ileride eklenecek Ã¶zellikler iÃ§in feature access kararlarÄ±
  - yeni modÃ¼ller iÃ§in canonical permission Ã¶rnekleri; Ã¶rnek olarak `history.continue_reading.read.own`, `history.timeline.read.own`, `history.library.read.own`, `history.bookmark.write.own`, `history.library.read.public`, `manga.discovery.view`, `ads.view`, `shop.item.purchase`, `payment.mana.purchase` ve `payment.transaction.read.own`
  - her modÃ¼l iÃ§in zorunlu authorization kontratÄ± ve canonical permission isimlendirmesi
  - feature flag tabanlÄ± policy, rollout ve geÃ§ici davranÄ±ÅŸ kontrolÃ¼
  - admin tarafÄ±ndan tetiklenebilen site geneli, modÃ¼l bazlÄ±, Ã¶zellik bazlÄ±, audience bazlÄ± ve gerektiÄŸinde context bazlÄ± acil kapatma, emergency deny veya kill switch yÃ¼zeyi
  - availability kurallarÄ±nda `emergency_deny` > `deny/off` > `allow/on` > varsayÄ±lan deÄŸer Ã¶nceliÄŸi
  - audience scope yÃ¶netimi; baÅŸlangÄ±Ã§ta `all`, `guest`, `authenticated`, `authenticated_non_vip` ve `vip`, ileride gerekirse daha spesifik hedefler
  - aynÄ± key, aynÄ± `audience_kind + audience_selector` ve aynÄ± `scope_kind + scope_selector` iÃ§in Ã§akÄ±ÅŸan aktif kural bÄ±rakmama ve kayÄ±t aÅŸamasÄ±nda reddetme davranÄ±ÅŸÄ±
  - temporary grants, sÃ¼reli yetki verme ve kontrollÃ¼ yetki geri alma
  - policy versioning ve denial explanation surface
  - VIP kullanÄ±labilirliÄŸi ile VIP entitlement sÃ¼resinin birbirinden ayrÄ± ele alÄ±nmasÄ±; sistem kaynaklÄ± global pasiflikte kalan sÃ¼renin dondurulmasÄ±
  - kiÅŸi bazlÄ± ve alan bazlÄ± moderatÃ¶r yetkilendirmesi; Ã¶rnek olarak yorum moderatÃ¶rÃ¼, bÃ¶lÃ¼m moderatÃ¶rÃ¼ veya manga moderatÃ¶rÃ¼ gibi ayrÄ±k yetki yÃ¼zeyleri
  - `moderation` modÃ¼lÃ¼ ile baÄŸlanacak moderator scope ve delegation altyapÄ±sÄ±

### 21.4) `admin` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `admin` olmalÄ±dÄ±r.
- `admin` modÃ¼lÃ¼ yÃ¶netim, tam yetkili inceleme, merkezi ayar ve operasyon use-case'lerinin sahibidir.
- `admin` modÃ¼lÃ¼ dashboard, yÃ¶netim giriÅŸ noktalarÄ±, kullanÄ±cÄ± yÃ¶netim akÄ±ÅŸlarÄ±, support review, tam yetkili moderasyon gÃ¶zetimi, operasyonel kontrol ve admin audit akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `admin` modÃ¼lÃ¼ kendi iÃ§inde authorization kararÄ±, role/permission kararÄ± veya kullanÄ±cÄ± profil veri sahipliÄŸi Ã¼retmemelidir.
- `admin` veri sahipliÄŸi; yÃ¶netimsel iÅŸlem kayÄ±tlarÄ±, admin notlarÄ±, admin use-case akÄ±ÅŸlarÄ± ve operasyonel gÃ¶rÃ¼nÃ¼m alanlarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `admin` modÃ¼lÃ¼nÃ¼n tÃ¼m kritik akÄ±ÅŸlarÄ± `access` guard veya policy kararlarÄ± ile korunmalÄ±dÄ±r.
- Gerekli role veya permission'a sahip admin kullanÄ±cÄ±larÄ± sistemdeki yÃ¶netim ve inceleme yÃ¼zeylerine tam eriÅŸim taÅŸÄ±yabilir; bu durum gÃ¼nlÃ¼k scoped moderator use-case sahipliÄŸini `moderation` modÃ¼lÃ¼nden almaz.
- `admin` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/admin.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - admin dashboard ve dashboard veri ihtiyaÃ§larÄ±
  - kullanÄ±cÄ± yÃ¶netimi, kullanÄ±cÄ± durum mÃ¼dahaleleri, warning, restriction, suspend ve ban akÄ±ÅŸlarÄ±
  - manga, chapter ve comment iÃ§in yÃ¼ksek riskli moderasyon handoff, escalation ve tam yetkili yÃ¶netimsel inceleme yÃ¼zeyleri
  - support review queue, destek karar yÃ¼rÃ¼tme yÃ¼zeyleri ve iletiÅŸim odaklÄ± yÃ¶netim akÄ±ÅŸlarÄ±
  - yÃ¼ksek riskli admin aksiyonlarÄ±nda aÃ§Ä±k permission, zorunlu reason ve gerektiÄŸinde ek doÄŸrulama
  - cache temizleme, log gÃ¶rÃ¼ntÃ¼leme ve sistem saÄŸlÄ±k durumu gibi operasyonel araÃ§lar
  - ileri aÅŸamalarda eklenecek yeni yÃ¶netim yÃ¼zeyleri iÃ§in hazÄ±rlÄ±k
  - bulk actions, case timeline ve toplu moderasyon araÃ§larÄ±
  - canned moderation actions ve approval workflow yapÄ±larÄ±
  - sistem genelindeki runtime ayarlarÄ±n, modÃ¼l aÃ§ma-kapama ve alt Ã¶zellik aÃ§ma-kapama yÃ¼zeylerinin merkezi yÃ¶netimi
  - `site`, `communication`, `operations`, `security_auth`, `access_availability`, `content`, `reading`, `engagement`, `support`, `membership`, `social`, `gamification` ve `economy` kategorileri iÃ§in geniÅŸleyebilir settings merkezi
  - baÅŸarÄ±sÄ±z giriÅŸ denemesi limiti, yorum gÃ¶nderme aralÄ±ÄŸÄ±, yorumlarÄ±n manga detayÄ±nda aÃ§Ä±k veya kapalÄ± olmasÄ± gibi iÅŸ kuralÄ± eÅŸiklerinin yÃ¶netim yÃ¼zeyi
  - env veya secret yÃ¶netimi gerektiren teknik config ile admin runtime ayarlarÄ±nÄ±n kesin olarak ayrÄ±lmasÄ±
  - ayar metadata kataloÄŸu; key, tip, scope kind, scope selector, audience kind, audience selector, apply mode, cache strategy, error response policy ve entitlement impact policy gibi alanlarÄ±n merkezi yÃ¶netimi
  - `moderation` modÃ¼lÃ¼ eklense bile merkezi ayar ve kill switch yÃ¼zeylerinin yalnÄ±zca admin tarafÄ±nda kalmasÄ±; moderatÃ¶r paneline delegasyon yapÄ±lmamasÄ±
  - operatÃ¶r iÅŸ yÃ¼kÃ¼ gÃ¶rÃ¼nÃ¼mÃ¼ ve moderasyon veya destek yÃ¼k daÄŸÄ±lÄ±mÄ± takibi
  - access iÃ§indeki emergency deny veya kill switch yÃ¼zeyini yÃ¶netebilecek operasyon kontrol noktalarÄ±

### 21.5) `manga` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `manga` olmalÄ±dÄ±r.
- `manga` modÃ¼lÃ¼ ana iÃ§erik varlÄ±ÄŸÄ±nÄ±n, metadata yapÄ±sÄ±nÄ±n, taxonomy iliÅŸkilerinin ve iÃ§erik yaÅŸam dÃ¶ngÃ¼sÃ¼ verisinin sahibidir.
- `manga` modÃ¼lÃ¼ CRUD, listing, detail, search, filtering, sorting, publish akÄ±ÅŸlarÄ±, metadata/taxonomy alanlarÄ± ve iÃ§erik sayaÃ§larÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `manga` modÃ¼lÃ¼ chapter iÃ§in varsayÄ±lan okuma eriÅŸim verisini taÅŸÄ±yabilir; ancak eriÅŸim kararÄ±nÄ± kendi iÃ§inde Ã¼retmemelidir.
- `manga` veri sahipliÄŸi; baÅŸlÄ±k, Ã¶zet, gÃ¶rsel, taxonomy, yayÄ±n durumu, gÃ¶rÃ¼nÃ¼rlÃ¼kle iliÅŸkili state alanlarÄ± ve iÃ§erik sayaÃ§larÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `manga` iÃ§indeki `chapter_count`, `comment_count` ve benzeri sayaÃ§ alanlarÄ± denormalize okuma alanlarÄ± olarak ele alÄ±nmalÄ±; canonical kaynak verisi ilgili kaynak modÃ¼lde kalmalÄ±dÄ±r.
- `manga` sayaÃ§ gÃ¼ncellemeleri `chapter` veya `comment` modÃ¼lÃ¼nÃ¼n manga tablosuna doÄŸrudan yazmasÄ± ile yapÄ±lmamalÄ±; event, projection veya aÃ§Ä±k counter contract yÃ¼zeyi ile senkronize edilmelidir.
- `manga` sayaÃ§larÄ± iÃ§in kabul edilen gecikme modeli ve gerektiÄŸinde reconcile veya yeniden hesaplama yolu dokÃ¼mante edilmelidir.
- `manga` modÃ¼lÃ¼nÃ¼n public surface'i iÃ§erik listeleme, iÃ§erik detay, yÃ¶netimsel iÃ§erik iÅŸlemleri ve taxonomy iliÅŸkileri iÃ§in gerekli yÃ¼zey ile sÄ±nÄ±rlÄ± olmalÄ±dÄ±r.
- `manga` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/manga.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - slug ve benzersizlik kurallarÄ±
  - alternative titles, short summary, cover image, banner image ve SEO alanlarÄ±
  - genre, tag, theme ve content warning taxonomy yapÄ±larÄ±
  - draft, scheduled, published ve archived veya unpublished benzeri publish yaÅŸam dÃ¶ngÃ¼sÃ¼
  - featured veya recommended iÅŸaretleri, editoryal koleksiyonlar ve iÃ§erik sayaÃ§larÄ±
  - soft delete ve restore desteÄŸi
  - chapter iÃ§in varsayÄ±lan read access ve varsayÄ±lan VIP erken eriÅŸim ayarlarÄ±
  - release schedule ve translation group gibi yayÄ±n planÄ± alanlarÄ±
  - toplu planlÄ± yayÄ±n ve editoryal yayÄ±n paketleri
  - recommendation, iÃ§erik koleksiyonu ve editoryal keÅŸif yÃ¼zeyleri; ancak kullanÄ±cÄ±ya ait continue reading, reading history veya bookmark-library kayÄ±tlarÄ± `history` modÃ¼lÃ¼nde kalmalÄ±dÄ±r
  - manga listeleme, manga detay ve editoryal gÃ¶rÃ¼nÃ¼rlÃ¼k gibi yÃ¼zeylerin admin tarafÄ±ndan yÃ¶netilen runtime ayarlar ile kontrollÃ¼ biÃ§imde daraltÄ±labilmesi

### 21.6) `chapter` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `chapter` olmalÄ±dÄ±r.
- `chapter` modÃ¼lÃ¼ manga iÃ§eriÄŸinin okunabilir bÃ¶lÃ¼m yapÄ±sÄ±nÄ±n, bÃ¶lÃ¼m sÄ±ralamasÄ±nÄ±n, bÃ¶lÃ¼m sayfalarÄ±nÄ±n ve bÃ¶lÃ¼m yaÅŸam dÃ¶ngÃ¼sÃ¼ verisinin sahibidir.
- `chapter` modÃ¼lÃ¼ CRUD, manga bazlÄ± chapter listesi, detail, read akÄ±ÅŸÄ±, navigation, page/media iliÅŸkileri, numbering ve publish akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `chapter` modÃ¼lÃ¼ chapter eriÅŸimini etkileyen veri alanlarÄ±nÄ± taÅŸÄ±yabilir; ancak guest/authenticated/vip/early access kararlarÄ±nÄ± kendi iÃ§inde Ã¼retmemelidir.
- `chapter` veri sahipliÄŸi; chapter metadata alanlarÄ±, page yapÄ±sÄ±, publish state, access state verisi ve navigation alanlarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `chapter` kullanÄ±cÄ±ya ait son okuma pozisyonu, reading session progress, continue reading kaydÄ± veya bookmark-library state'i taÅŸÄ±mamalÄ±; bunlar `history` modÃ¼lÃ¼nde tutulmalÄ±dÄ±r.
- `chapter` modÃ¼lÃ¼nÃ¼n public surface'i okuma akÄ±ÅŸÄ±, chapter detail ve yÃ¶netimsel chapter iÅŸlemleri iÃ§in gerekli contract yÃ¼zeyi ile sÄ±nÄ±rlÄ± olmalÄ±dÄ±r.
- `chapter` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/chapter.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - latest chapter list ile previous, next, first ve last navigation akÄ±ÅŸlarÄ±
  - chapter page yapÄ±sÄ±, width veya height bilgisi ve gerekirse long strip desteÄŸi
  - read_access_level, inherit_access_from_manga, early_access_enabled, early_access_level ve fallback alanlarÄ±
  - VIP Ã¶zel bÃ¶lÃ¼m ve VIP iÃ§in erken eriÅŸim bÃ¶lÃ¼m yapÄ±larÄ±nÄ±n birbirinden ayrÄ±lmasÄ±
  - misafir kullanÄ±cÄ±nÄ±n manga detayÄ±na eriÅŸebilmesi ama chapter okuma iÃ§in minimum `authenticated` kapÄ±sÄ±nÄ±n zorunlu olmasÄ±
  - early access zaman penceresi, pencere sonrasÄ± fallback eriÅŸimi ve access policy hizasÄ±
  - soft delete ve restore desteÄŸi
  - media validation ve CDN health kontrol yÃ¼zeyleri
  - bozuk veya eksik medya tespiti iÃ§in checksum veya benzeri bÃ¼tÃ¼nlÃ¼k doÄŸrulamasÄ±
  - `history` modÃ¼lÃ¼ iÃ§in continue reading, reading history, resume anchor ve progress entegrasyon yÃ¼zeyi
  - chapter okuma, preview veya belirli okuma alt yÃ¼zeylerinin admin tarafÄ±ndan yÃ¶netilen runtime ayarlar ile kontrollÃ¼ aÃ§Ä±lÄ±p kapatÄ±labilmesi
  - preview kapalÄ±yken detail veya tam read yÃ¼zeyinin otomatik olarak kapanmamasÄ±; her yÃ¼zeyin ayrÄ± kontrol edilebilmesi

### 21.7) `comment` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `comment` olmalÄ±dÄ±r.
- `comment` modÃ¼lÃ¼ yorum verisinin, thread/reply yapÄ±sÄ±nÄ±n, yorum gÃ¶rÃ¼nÃ¼rlÃ¼k verisinin ve yorum yaÅŸam dÃ¶ngÃ¼sÃ¼nÃ¼n sahibidir.
- `comment` modÃ¼lÃ¼ en az `manga` ve `chapter` hedef tiplerini desteklemeli; yorum create, edit, delete, reply, listing, moderation state ve thread akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `comment` modÃ¼lÃ¼ yorum gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼nÃ¼ etkileyen veri alanlarÄ±nÄ± taÅŸÄ±yabilir; ancak create/edit/delete/pin/lock gibi iÅŸlemlerin yetki kararÄ±nÄ± kendi iÃ§inde Ã¼retmemelidir.
- `comment` modÃ¼lÃ¼ sosyal duvar post'u veya sosyal duvar reply akÄ±ÅŸÄ±nÄ± sahiplenmemeli; bunlar `social` modÃ¼lÃ¼nde kalmalÄ± ve `comment` sistemine Ã¶rtÃ¼k olarak dÃ¶nÃ¼ÅŸtÃ¼rÃ¼lmemelidir.
- `comment` veri sahipliÄŸi; yorum iÃ§eriÄŸi, hedef iliÅŸkisi, reply yapÄ±sÄ±, moderation/spoiler/lock verileri ve sÄ±ralama/listeleme alanlarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `comment` modÃ¼lÃ¼ndeki `target_type` deÄŸerleri canonical olarak `docs/shared/target-types.md` dosyasÄ±ndaki kayÄ±tlarla hizalÄ± olmalÄ±dÄ±r.
- Yeni yorum hedef tipi eklendiÄŸinde `comment` modÃ¼lÃ¼, hedef modÃ¼l dokÃ¼manÄ± ve canonical target type kaydÄ± aynÄ± deÄŸiÅŸiklik setinde gÃ¼ncellenmelidir.
- `comment` modÃ¼lÃ¼nÃ¼n public surface'i yorum listeleme, yorum detay, thread akÄ±ÅŸÄ± ve hedef iliÅŸkisi iÃ§in gerekli yÃ¼zey ile sÄ±nÄ±rlÄ± olmalÄ±dÄ±r.
- `comment` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/comment.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - root comment ve reply thread yapÄ±sÄ±
  - newest, oldest ve popular sÄ±ralama seÃ§enekleri ile pagination
  - sanitize edilmiÅŸ iÃ§erik veya gÃ¼venli Ã§Ä±ktÄ± alanÄ±
  - anonymous gÃ¶rÃ¼ntÃ¼leme ve authenticated yorum yazma ayrÄ±mÄ±
  - reply derinliÄŸi, edit window, lock etkisi ve restore sÄ±nÄ±rlarÄ±
  - report edilebilir hedef olma ve ileri aÅŸamalarda yeni target_type geniÅŸletmeleri
  - geliÅŸmiÅŸ moderasyon ve anti-spam akÄ±ÅŸlarÄ±
  - yorum attachment desteÄŸi
  - anti-spam score veya yorum risk puanÄ± yaklaÅŸÄ±mÄ±
  - site geneli, manga detayÄ± veya chapter altÄ± yorum alanlarÄ±nÄ±n ayrÄ± ayrÄ± aÃ§Ä±lÄ±p kapatÄ±labilmesi ve yorum gÃ¶nderme aralÄ±ÄŸÄ± gibi etkileÅŸim eÅŸiklerinin yÃ¶netilebilmesi
  - yorum modÃ¼lÃ¼ iÃ§in read ve write yÃ¼zeylerinin ayrÄ± ayrÄ± kontrol edilebilmesi; varsayÄ±lan kapatma senaryosunda mevcut yorumlarÄ±n gÃ¶rÃ¼nÃ¼r kalÄ±p yeni yorumlarÄ±n engellenebilmesi
  - kullanÄ±cÄ± bazlÄ± mute ve moderation escalation davranÄ±ÅŸlarÄ±
  - sessiz moderasyon; iÃ§eriÄŸi tamamen silmeden gÃ¶rÃ¼nÃ¼rlÃ¼k kapsamÄ±nÄ± kÄ±sÄ±tlama yaklaÅŸÄ±mÄ±

### 21.8) `support` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `support` olmalÄ±dÄ±r.
- `support` modÃ¼lÃ¼ kullanÄ±cÄ± iletiÅŸim taleplerinin, destek biletlerinin, manga/chapter/comment iÃ§in hedefe baÄŸlÄ± iÃ§erik bildirimlerinin ve destek yaÅŸam dÃ¶ngÃ¼sÃ¼nÃ¼n sahibidir.
- AyrÄ± bir `report` leaf modÃ¼l aÃ§Ä±lmamalÄ±; iÃ§erik bildirimi akÄ±ÅŸlarÄ± `support` modÃ¼lÃ¼nÃ¼n destek ve ticket alanÄ± iÃ§indeki bir feature yÃ¼zeyi olarak ele alÄ±nmalÄ±dÄ±r.
- `support` modÃ¼lÃ¼ communication/create, ticket/create, hedefe baÄŸlÄ± report/create, own support list, support detail, support reply, category, priority, duplicate/spam kontrolÃ¼, review queue verisi, status update ve resolution note akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `support` veri sahipliÄŸi; destek kaydÄ±, `support_kind`, category, isteÄŸe baÄŸlÄ± hedef iliÅŸkisi, destek durumu, mesaj veya reply zinciri, Ã§Ã¶zÃ¼m notlarÄ± ve inceleme yaÅŸam dÃ¶ngÃ¼sÃ¼ verileri ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `support` modÃ¼lÃ¼ review ve karar verisini taÅŸÄ±r; ancak yetki kararÄ±nÄ± kendi iÃ§inde Ã¼retmez, yÃ¶netimsel karar yÃ¼rÃ¼tÃ¼mÃ¼ `admin`, authorization ise `access` ile korunur.
- `support` modÃ¼lÃ¼ndeki report kaydÄ± varsayÄ±lan olarak moderation case ile aynÄ± kayÄ±t sayÄ±lmamalÄ±dÄ±r; moderation ihtiyacÄ± oluÅŸtuÄŸunda linked case aÃ§Ä±labilir, ancak support intake kaydÄ± ve moderation case yaÅŸam dÃ¶ngÃ¼sÃ¼ ayrÄ± owner sÄ±nÄ±rlarÄ±nda kalmalÄ±dÄ±r.
- `support` modÃ¼lÃ¼ndeki `target_type` deÄŸerleri canonical olarak `docs/shared/target-types.md` dosyasÄ±ndaki kayÄ±tlarla hizalÄ± olmalÄ±dÄ±r.
- Genel iletiÅŸim veya hedefsiz destek biletlerinde `target_type` zorunlu olmamalÄ±; bu alan yalnÄ±zca manga/chapter/comment gibi hedefe baÄŸlÄ± kayÄ±tlar iÃ§in kullanÄ±lmalÄ±dÄ±r.
- Yeni support hedef tipi eklendiÄŸinde `support` modÃ¼lÃ¼, hedef modÃ¼l dokÃ¼manÄ± ve canonical target type kaydÄ± aynÄ± deÄŸiÅŸiklik setinde gÃ¼ncellenmelidir.
- `support` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/support.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - communication, ticket ve hedefe baÄŸlÄ± report akÄ±ÅŸlarÄ±nÄ±n aynÄ± modÃ¼lde ama ayrÄ± kayÄ±t mantÄ±ÄŸÄ± ile taÅŸÄ±nmasÄ±
  - manga, chapter ve comment hedef tipleri ile isteÄŸe baÄŸlÄ± target relation yapÄ±sÄ±
  - support kind, category, priority ve reason code kataloÄŸu
  - duplicate veya spam kontrolÃ¼ ve aynÄ± kullanÄ±cÄ± ile aynÄ± hedef iÃ§in tekrar bildirim davranÄ±ÅŸÄ±
  - review queue, status update, support reply, resolution note, assignee, reviewed_by ve resolved_at alanlarÄ±
  - iletiÅŸim kaydÄ± ile hedefe baÄŸlÄ± iÃ§erik bildiriminin aynÄ± yaÅŸam dÃ¶ngÃ¼sÃ¼nde ama ayrÄ± kurallarla iÅŸlenebilmesi
  - support attachment desteÄŸi
  - canned support replies ve cevap ÅŸablonlarÄ±
  - SLA, Ã¶nceliklendirme ve escalation yÃ¼zeyi
  - support modÃ¼lÃ¼ altÄ±ndaki communication, ticket ve report create yÃ¼zeylerinin, attachment kabulÃ¼nÃ¼n ve intake davranÄ±ÅŸlarÄ±nÄ±n admin tarafÄ±ndan yÃ¶netilen runtime ayarlar ile kontrol edilebilmesi
  - support iÃ§in yeni kayÄ±t alÄ±mÄ±nÄ±n durdurulabilmesi, ancak mevcut kayÄ±tlarÄ±n kullanÄ±cÄ± veya admin tarafÄ±ndan okunabilir kalmasÄ± gibi intake pause davranÄ±ÅŸlarÄ±
  - raporlayan kullanÄ±cÄ± gÃ¼ven skoru ve gÃ¼ven puanÄ±nÄ±n inceleme Ã¶nceliÄŸine etkisi
  - bozuk medya veya eksik sayfa bildirimlerinde checksum veya bÃ¼tÃ¼nlÃ¼k doÄŸrulamasÄ± ile desteklenebilen inceleme verisi

### 21.9) `moderation` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `moderation` olmalÄ±dÄ±r.
- `moderation` modÃ¼lÃ¼ role bazlÄ± veya kullanÄ±cÄ± bazlÄ± scoped moderatÃ¶r panelinin, vaka inceleme akÄ±ÅŸlarÄ±nÄ±n ve moderatÃ¶r iÅŸ yÃ¼kÃ¼ sÃ¼reÃ§lerinin sahibi olmalÄ±dÄ±r.
- `moderation` modÃ¼lÃ¼ queue, assignment, case detail, moderator note, sÄ±nÄ±rlÄ± aksiyon yÃ¼rÃ¼tme ve escalation akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `moderation` modÃ¼lÃ¼ authorization, role veya permission sahipliÄŸi Ã¼retmemeli; moderator scope ve yetki kararlarÄ± `access` ile korunmalÄ±dÄ±r.
- `moderation` veri sahipliÄŸi; moderation case, assignment, moderator note, action summary ve escalation lifecycle alanlarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `moderation` gÃ¼nlÃ¼k scoped inceleme sahibidir; ancak `admin` tarafÄ±ndan aynÄ± case Ã¼zerinde verilen override, reopen, freeze, reassignment veya final kararlar daha yÃ¼ksek precedence taÅŸÄ±r.
- `moderation` case hedefleri canonical olarak `docs/shared/target-types.md` dosyasÄ±ndaki kayÄ±tlarla hizalÄ± olmalÄ±; alt yÃ¼zey bilgisi `target_type` iÃ§ine deÄŸil context verisine taÅŸÄ±nmalÄ±dÄ±r.
- `moderation` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/moderation.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - yorum, chapter ve manga yÃ¼zeyleri iÃ§in scoped queue yapÄ±larÄ±
  - comment moderator, chapter moderator veya manga moderator gibi role veya kullanÄ±cÄ± bazlÄ± scope modelleri
  - moderator assignment, handoff ve escalation akÄ±ÅŸlarÄ±
  - vaka timeline, moderator note ve karar Ã¶zeti alanlarÄ±
  - sÄ±nÄ±rlÄ± moderator aksiyonlarÄ±; Ã¶rnek olarak hide, unhide, lock, unlock, escalate veya review complete yÃ¼zeyleri
  - admin tarafÄ±ndan yÃ¶netilen runtime ayarlar ile queue veya aksiyon yÃ¼zeylerinin ayrÄ± ayrÄ± aÃ§Ä±lÄ±p kapatÄ±labilmesi
  - moderasyon panelinin aÃ§Ä±k kalmasÄ±, ancak merkezi settings ve kill switch yÃ¼zeylerinin admin modÃ¼lÃ¼nde kalmasÄ±
  - admin tarafÄ±na workload, escalation ve audit sinyali Ã¼retebilen entegrasyon yÃ¼zeyleri

### 21.10) `notification` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `notification` olmalÄ±dÄ±r.
- `notification` modÃ¼lÃ¼ sistem genelindeki bildirim Ã¼retimi, teslimi, kategori yÃ¶netimi ve detaylÄ± bildirim tercihleri verisinin sahibi olmalÄ±dÄ±r.
- `notification` modÃ¼lÃ¼ in-app inbox, read veya unread akÄ±ÅŸlarÄ±, category, channel, template, delivery attempt ve suppression yÃ¼zeylerini taÅŸÄ±malÄ±dÄ±r.
- `notification` modÃ¼lÃ¼ business event sahipliÄŸi veya authorization kararÄ± Ã¼retmemelidir; bildirim olaylarÄ±nÄ± producer modÃ¼llerden almalÄ± ve own-surface eriÅŸimini `access` ile korumalÄ±dÄ±r.
- `notification` veri sahipliÄŸi; notification kaydÄ±, delivery durumu, template veya category tanÄ±mÄ± ve kullanÄ±cÄ± bildirim tercihleri ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `notification` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/notification.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - in-app inbox ve read veya unread davranÄ±ÅŸlarÄ±
  - category, template ve channel yÃ¶netimi
  - module veya category bazlÄ± delivery pause ve flood control yÃ¼zeyleri
  - quiet-hour, mute, digest veya batch delivery davranÄ±ÅŸlarÄ±
  - detail bildirim tercihleri sahipliÄŸi; `user` modÃ¼lÃ¼nde yalnÄ±zca Ã¶zet preference sinyali kalmasÄ±
  - admin tarafÄ±ndan category, channel veya feature bazlÄ± aÃ§ma-kapama ve eÅŸik yÃ¶netimi
  - `social`, `mission`, `royalpass`, `support`, `moderation` ve diÄŸer producer modÃ¼ller ile event kontratlarÄ±

### 21.11) `social` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `social` olmalÄ±dÄ±r.
- `social` modÃ¼lÃ¼ kullanÄ±cÄ±lar arasÄ± arkadaÅŸlÄ±k, takip, sosyal duvar, duvar altÄ± etkileÅŸim ve mesajlaÅŸma iÅŸ alanlarÄ±nÄ±n sahibi olmalÄ±dÄ±r.
- `social` modÃ¼lÃ¼ friendship request, friendship state, follow relation, wall post veya wall reply ve direct message thread akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `social` modÃ¼lÃ¼ manga veya chapter iÃ§erik yorumlarÄ±nÄ± sahiplenmemeli; bu alanlar `comment` modÃ¼lÃ¼nde kalmalÄ±dÄ±r.
- `social` modÃ¼lÃ¼ndeki wall reply yapÄ±sÄ± social-native kabul edilmeli; `comment` thread sistemi ile Ã¶rtÃ¼k olarak birleÅŸtirilmemelidir.
- `social` veri sahipliÄŸi; sosyal iliÅŸki kayÄ±tlarÄ±, social block veya mute iliÅŸkileri, sosyal iÃ§erik kayÄ±tlarÄ± ve sosyal privacy sinyalleri ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `social` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/social.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - friend request, accept, reject, remove ve friend list akÄ±ÅŸlarÄ±
  - follow veya unfollow davranÄ±ÅŸlarÄ± ve takip listeleri
  - profil sosyal duvarÄ±, duvar post veya duvar reply yÃ¼zeyleri
  - direct message thread, mesaj gÃ¶nderme ve mesaj gÃ¶rÃ¼nÃ¼rlÃ¼k kurallarÄ±
  - social block veya mute sahipliÄŸi; bu alanlarÄ±n `user` modÃ¼lÃ¼nde veri sahipliÄŸi olarak tutulmamasÄ±
  - social block veya privacy deny sinyalinin final authorization kararÄ±nda `access` tarafÄ±ndan deny precedence ile yorumlanmasÄ±; mute sinyalinin ise ayrÄ±ca dokÃ¼mante edilmedikÃ§e teslim veya gÃ¶rÃ¼nÃ¼rlÃ¼k sinyali olarak kalmasÄ±
  - friendship, follow, messaging ve wall yÃ¼zeyleri iÃ§in ayrÄ± runtime control anahtarlarÄ±
  - admin tarafÄ±ndan messaging, wall, follow veya friendship yÃ¼zeylerinin ayrÄ± ayrÄ± aÃ§Ä±lÄ±p kapatÄ±labilmesi
  - bildirim, anti-spam ve social privacy sinyallerinin birbirine karÄ±ÅŸmadan Ã§alÄ±ÅŸmasÄ±

### 21.12) `inventory` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `inventory` olmalÄ±dÄ±r.
- `inventory` modÃ¼lÃ¼ item tanÄ±mÄ±, kullanÄ±cÄ± envanteri, reward sahipliÄŸi ve final grant yÃ¼rÃ¼tÃ¼mÃ¼nÃ¼n sahibi olmalÄ±dÄ±r.
- `inventory` modÃ¼lÃ¼ ownable item definition, user inventory entry, grant, revoke, claim, consume ve equip akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `inventory` modÃ¼lÃ¼ sellable shop product veya offer catalog sahipliÄŸi Ã¼retmemeli; `shop` ile iliÅŸkisi product-to-item mapping ve final grant kontratÄ± Ã¼zerinden kurulmalÄ±dÄ±r.
- `inventory` modÃ¼lÃ¼ Ã¶deme, gÃ¶rev ilerlemesi veya pass season sahipliÄŸi Ã¼retmemeli; yalnÄ±zca item sahipliÄŸi ve item durumunu taÅŸÄ±malÄ±dÄ±r.
- `inventory` veri sahipliÄŸi; ownable item tanÄ±mÄ±, quantity veya stack durumu, expiry, equip state ve source reference alanlarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±; sellable catalog verisini iÃ§ermemelidir.
- `inventory` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/inventory.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - stackable ve non-stackable item ayrÄ±mÄ±
  - grant, reward teslim yÃ¼rÃ¼tÃ¼mÃ¼, revoke, consume ve equip yÃ¼zeyleri
  - source reference ve idempotent grant davranÄ±ÅŸlarÄ±
  - `user` modÃ¼lÃ¼nde seÃ§ilen kozmetik gÃ¶rÃ¼nÃ¼m referanslarÄ±nÄ±n `inventory` sahipliÄŸi ile hizalanmasÄ±
  - admin tarafÄ±ndan inventory gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼, claim, consume veya equip yÃ¼zeylerinin aÃ§Ä±lÄ±p kapatÄ±lmasÄ±
  - `mission`, `royalpass` ve ileride gelecek diÄŸer reward producer modÃ¼lleri ile kontrollÃ¼ grant kontratlarÄ±

### 21.13) `mission` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `mission` olmalÄ±dÄ±r.
- `mission` modÃ¼lÃ¼ gÃ¼nlÃ¼k, haftalÄ±k, aylÄ±k, event ve seviye bazlÄ± gÃ¶rev tanÄ±mlarÄ±nÄ±n, gÃ¶rev ilerlemesinin ve reward iÃ§in claim uygunluÄŸu veya claim request yaÅŸam dÃ¶ngÃ¼sÃ¼nÃ¼n sahibi olmalÄ±dÄ±r.
- `mission` modÃ¼lÃ¼ mission definition, objective yapÄ±sÄ±, progress kaydÄ±, completion, claim eligibility ve reset pencerelerini taÅŸÄ±malÄ±dÄ±r.
- `mission` modÃ¼lÃ¼ global EXP veya level sahipliÄŸini tek baÅŸÄ±na Ã¼retmemeli; `user` modÃ¼lÃ¼ndeki progression sinyallerini tÃ¼keterek gÃ¶rev deÄŸerlendirmesi yapabilmelidir.
- `mission` veri sahipliÄŸi; gÃ¶rev tanÄ±mÄ±, gÃ¶rev kategorisi, progress state, claim eligibility state ve reward reference alanlarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `mission` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/mission.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - daily, weekly, monthly, event ve level-based mission tipleri
  - recurring reset, dÃ¶nemsel yenileme ve gerektiÄŸinde streak davranÄ±ÅŸlarÄ±
  - producer event kontratlarÄ±; Ã¶rnek olarak okuma, yorum, sosyal etkileÅŸim veya diÄŸer gÃ¶rev sinyalleri
  - claim request yÃ¼zeyinin `inventory` iÃ§indeki final grant sahipliÄŸi ile hizalanmasÄ±
  - claim, reset veya mission category yÃ¼zeylerinin admin tarafÄ±ndan ayrÄ± ayrÄ± aÃ§Ä±lÄ±p kapatÄ±labilmesi
  - `notification` ile gÃ¶rev tamamlama ve claim bildirim entegrasyonu
  - `royalpass` iÃ§in gÃ¶revden progress besleme yÃ¼zeyi

### 21.14) `royalpass` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `royalpass` olmalÄ±dÄ±r.
- `royalpass` modÃ¼lÃ¼ aylÄ±k season yapÄ±sÄ±nÄ±n, free veya premium track ilerlemesinin ve season Ã¶dÃ¼lÃ¼ iÃ§in claim uygunluÄŸu veya claim request yaÅŸam dÃ¶ngÃ¼sÃ¼nÃ¼n sahibi olmalÄ±dÄ±r.
- `royalpass` modÃ¼lÃ¼ season, tier, track, user season progress ve reward claim eligibility akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `royalpass` modÃ¼lÃ¼ gÃ¶rev tanÄ±mÄ±, item sahipliÄŸi veya Ã¶deme sahipliÄŸi Ã¼retmemelidir; season iÃ§i progress ve claim eligibility sahipliÄŸi ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `royalpass` veri sahipliÄŸi; season tanÄ±mÄ±, track veya tier yapÄ±sÄ±, user season claim eligibility state ve premium activation referanslarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `royalpass` premium aktivasyonu Ã¼rÃ¼nleÅŸmiÅŸ satÄ±n alma akÄ±ÅŸÄ±nda canonical olarak `shop` Ã¼zerinden baÅŸlamalÄ±, gerÃ§ek para veya mana checkout veya bakiye doÄŸruluÄŸu gerekiyorsa `payment` tarafÄ±ndan tamamlanmalÄ± ve final premium activation referansÄ± `royalpass` tarafÄ±ndan tÃ¼ketilmelidir.
- `royalpass` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/royalpass.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - aylÄ±k season yaÅŸam dÃ¶ngÃ¼sÃ¼ ve season archive yapÄ±larÄ±
  - free track ve premium track ayrÄ±mÄ±
  - mission tabanlÄ± progress veya puan besleme davranÄ±ÅŸlarÄ±
  - claim request yÃ¼zeyinin `inventory` iÃ§indeki final grant sahipliÄŸi ile hizalanmasÄ±
  - season gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼, claim yÃ¼zeyi veya premium track yÃ¼zeyinin admin tarafÄ±ndan ayrÄ± ayrÄ± aÃ§Ä±lÄ±p kapatÄ±labilmesi
  - Ã¼rÃ¼nleÅŸmiÅŸ premium pass satÄ±n alma zincirinde `shop` Ã¼rÃ¼n orkestrasyonu, `payment` checkout veya bakiye doÄŸruluÄŸu ve `royalpass` entitlement sahipliÄŸinin aÃ§Ä±kÃ§a ayrÄ±lmasÄ±
  - season pause veya sistem kaynaklÄ± pasiflikte claim ve reward davranÄ±ÅŸlarÄ±nÄ±n gÃ¼venli yÃ¶netimi
  - `notification` ile season baÅŸlangÄ±cÄ±, reward claim ve kalan tier bildirimi entegrasyonu

### 21.15) `history` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `history` olmalÄ±dÄ±r.
- `history` modÃ¼lÃ¼ kullanÄ±cÄ±ya ait continue reading, reading history, bookmark-library ve okuma devamlÄ±lÄ±ÄŸÄ± kayÄ±tlarÄ±nÄ±n sahibi olmalÄ±dÄ±r.
- `history` modÃ¼lÃ¼ user-manga library iliÅŸkisi, user-chapter son okuma durumu, reading session checkpoint, resume anchor ve own history timeline akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `history` modÃ¼lÃ¼ manga metadata, chapter page yapÄ±sÄ±, kullanÄ±cÄ± profil verisi veya authorization kararÄ± Ã¼retmemelidir.
- `history` veri sahipliÄŸi; kullanÄ±cÄ±ya ait library entry kayÄ±tlarÄ±, bookmark veya favorite iÅŸaretleri, son okunan chapter veya page referanslarÄ±, okuma progress snapshot'larÄ±, history timeline verileri ve entry-level share metadata alanlarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `history` modÃ¼lÃ¼ public veya shared library gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼nÃ¼ doÄŸrudan kendi iÃ§inde karara baÄŸlamamalÄ±; global visibility default'larÄ± `user`, permission kararlarÄ± `access`, entry-level share metadata ise `history` iÃ§inde kalacak ÅŸekilde Ã§alÄ±ÅŸmalÄ±dÄ±r.
- `history` access entegrasyonunda en az `history.continue_reading.read.own`, `history.timeline.read.own`, `history.library.read.own`, `history.bookmark.write.own` ve gerektiÄŸinde `history.library.read.public` gibi canonical permission Ã¶rnekleri dokÃ¼mante edilmelidir.
- `history` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/history.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - continue reading yÃ¼zeyi
  - reading history timeline veya own reading log gÃ¶rÃ¼nÃ¼mÃ¼
  - bookmark, favorite, library status ve gerektiÄŸinde custom shelf benzeri kayÄ±tlar
  - chapter read start, checkpoint, finish ve resume anchor entegrasyonu
  - cihazlar arasÄ± okuma devamlÄ±lÄ±ÄŸÄ± ve duplicate progress yazÄ±mÄ±na karÅŸÄ± idempotent gÃ¼ncelleme davranÄ±ÅŸÄ±
  - `manga`, `chapter`, `mission` ve ileride recommendation tarafÄ±na kontrollÃ¼ okuma sinyali veya Ã¶zet yÃ¼zeyi
  - continue reading, library, timeline veya bookmark write alt yÃ¼zeylerinin admin tarafÄ±ndan ayrÄ± ayrÄ± aÃ§Ä±lÄ±p kapatÄ±labilmesi

### 21.16) `ads` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `ads` olmalÄ±dÄ±r.
- `ads` modÃ¼lÃ¼ reklam yerleÅŸimi, kampanya, kreatif, teslim planÄ± ve gÃ¶sterim Ã¶lÃ§Ã¼mlemesi alanlarÄ±nÄ±n sahibi olmalÄ±dÄ±r.
- `ads` modÃ¼lÃ¼ placement, campaign, creative, active window, priority, frequency cap, impression ve click akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `ads` modÃ¼lÃ¼ VIP reklamsÄ±z deneyim veya audience eriÅŸim kararÄ±nÄ± kendi iÃ§inde Ã¼retmemelidir; bu yorumlar `access` ile yapÄ±lmalÄ±dÄ±r.
- `ads` veri sahipliÄŸi; placement tanÄ±mlarÄ±, campaign yapÄ±landÄ±rmalarÄ±, creative kayÄ±tlarÄ±, delivery metadata, gÃ¶sterim veya tÄ±klama loglarÄ± ve reklam gÃ¶rÃ¼nÃ¼rlÃ¼k state alanlarÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `ads` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/ads.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - ana sayfa, listeleme, manga detay ve chapter Ã§evresi placement yapÄ±larÄ±
  - campaign aktiflik penceresi, Ã¶ncelik ve frequency cap yÃ¶netimi
  - impression, click ve temel performans Ã¶lÃ§Ã¼mleme yÃ¼zeyleri
  - VIP reklamsÄ±z deneyim ile access policy hizasÄ±
  - surface, placement veya campaign bazlÄ± admin runtime control anahtarlarÄ± ve operasyon yÃ¼zeyleri

### 21.17) `shop` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `shop` olmalÄ±dÄ±r.
- `shop` modÃ¼lÃ¼ sellable product veya offer kataloÄŸu, teklif gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼, fiyatlandÄ±rma, satÄ±n alma orkestrasyonu ve Ã¼rÃ¼n kullanÄ±m kurallarÄ±nÄ±n sahibi olmalÄ±dÄ±r.
- `shop` modÃ¼lÃ¼ sellable product, offer, price plan, purchase intent, purchase eligibility ve kullanÄ±m kuralÄ± akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `shop` modÃ¼lÃ¼ final item sahipliÄŸi, equip state veya ledger doÄŸruluÄŸunu kendi iÃ§inde Ã¼retmemelidir; sahiplik `inventory`, finansal bakiye ve iÅŸlem doÄŸruluÄŸu `payment` modÃ¼lÃ¼nde kalmalÄ±dÄ±r. Stage 29 geÃ§iÅŸinde kullanabileceÄŸi allowance bridge verisi canonical bakiye owner'lÄ±ÄŸÄ± sayÄ±lmamalÄ±dÄ±r.
- `shop` veri sahipliÄŸi; sellable product veya offer kataloÄŸu, fiyat tanÄ±mÄ±, indirim veya kampanya metadata'sÄ±, purchase request veya order kayÄ±tlarÄ±, Ã¼rÃ¼n gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ve kullanÄ±m kÄ±sÄ±tÄ± verileri ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `shop` modÃ¼lÃ¼ `payment` Ã¶ncesi aÅŸamada yalnÄ±zca purchase eligibility iÃ§in geÃ§ici `seed_mana_allowance_snapshot` veya operasyonel allowance read modelini kullanabilir; bu kÃ¶prÃ¼ veri `payment` devreye girince kaldÄ±rÄ±lmalÄ±dÄ±r.
- `shop` modÃ¼lÃ¼ Ã¼rÃ¼nleÅŸmiÅŸ RoyalPass veya benzeri entitlement Ã¼rÃ¼nlerinde canonical purchase intent giriÅŸ noktasÄ± olabilir; ancak entitlement sahipliÄŸi ilgili hedef modÃ¼lde, bakiye veya checkout doÄŸruluÄŸu ise `payment` iÃ§inde kalmalÄ±dÄ±r.
- `shop` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/shop.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - kozmetik Ã¼rÃ¼n kataloÄŸu, kategori ve slot uyumluluÄŸu
  - mana bazlÄ± fiyatlandÄ±rma ve kampanya gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼
  - satÄ±n alma isteÄŸi, idempotent purchase davranÄ±ÅŸÄ± ve duplicate purchase korumasÄ±
  - VIP, level, RoyalPass veya gÃ¶rev kaynaklÄ± Ã¼rÃ¼n uygunluÄŸu sinyalleri
  - `inventory` ile final grant veya equip kontratÄ± ve `payment` ile bakiye dÃ¼ÅŸÃ¼m mutabakatÄ±
  - katalog, campaign, purchase veya belirli shop alt yÃ¼zeylerinin admin tarafÄ±ndan ayrÄ± ayrÄ± aÃ§Ä±lÄ±p kapatÄ±labilmesi

### 21.18) `payment` ModÃ¼lÃ¼ KurallarÄ±
- Canonical modÃ¼l adÄ± `payment` olmalÄ±dÄ±r.
- `payment` modÃ¼lÃ¼ mana satÄ±n alma, Ã¶deme saÄŸlayÄ±cÄ±sÄ± entegrasyonu, cÃ¼zdan veya ledger doÄŸruluÄŸu ve finansal iÅŸlem kayÄ±tlarÄ±nÄ±n sahibi olmalÄ±dÄ±r.
- `payment` modÃ¼lÃ¼ mana package, checkout session, provider callback, transaction, ledger entry, balance snapshot ve refund veya reversal akÄ±ÅŸlarÄ±nÄ± taÅŸÄ±malÄ±dÄ±r.
- `payment` modÃ¼lÃ¼ devreye girdiÄŸinde `shop` iÃ§indeki geÃ§ici `seed_mana_allowance_snapshot` veya operasyonel allowance bridge yÃ¼zeyini devralmalÄ± ve canonical bakiye owner'lÄ±ÄŸÄ±nÄ± tek baÅŸÄ±na Ã¼stlenmelidir.
- `payment` modÃ¼lÃ¼ Ã¼rÃ¼n kataloÄŸu, item sahipliÄŸi veya authorization kararÄ± Ã¼retmemelidir; katalog `shop`, sahiplik `inventory`, eriÅŸim kararÄ± `access` modÃ¼lÃ¼nde kalmalÄ±dÄ±r.
- `payment` modÃ¼lÃ¼ checkout, mana purchase ve finansal doÄŸruluÄŸun sahibidir; ancak Ã¼rÃ¼nleÅŸmiÅŸ RoyalPass benzeri entitlement akÄ±ÅŸlarÄ±nda final entitlement owner'lÄ±ÄŸÄ±na doÄŸrudan geÃ§mez, onaylanmÄ±ÅŸ Ã¶deme sonucunu ilgili modÃ¼le kontrollÃ¼ kontrat ile aktarÄ±r.
- `payment` veri sahipliÄŸi; provider session kayÄ±tlarÄ±, purchase order veya transaction kayÄ±tlarÄ±, ledger hareketleri, bakiye snapshot'larÄ±, fraud review state ve finansal audit metadata'sÄ± ile sÄ±nÄ±rlÄ± kalmalÄ±dÄ±r.
- `payment` modÃ¼lÃ¼nÃ¼n ana referans dokÃ¼manÄ± `docs/modules/payment.md` olmalÄ±dÄ±r.
- GeliÅŸtirme ile gelen baÅŸlÄ±ca baÅŸlÄ±klar ÅŸunlardÄ±r:
  - mana satÄ±n alma paketleri ve provider checkout oturumlarÄ±
  - pending, success, failed, cancelled, refunded veya reversed iÅŸlem durumlarÄ±
  - webhook veya callback doÄŸrulamasÄ± ve idempotent iÅŸleme
  - `shop` ile bakiye dÃ¼ÅŸÃ¼m veya mutabakat entegrasyonu
  - fraud review, audit ve finansal gÃ¶rÃ¼nÃ¼rlÃ¼k yÃ¼zeyleri
  - mana purchase, checkout veya iÅŸlem gÃ¶rÃ¼ntÃ¼leme alt yÃ¼zeylerinin admin tarafÄ±ndan kontrollÃ¼ biÃ§imde daraltÄ±labilmesi

## 22) DokÃ¼mantasyon ve Ã‡apraz Kesit StandartlarÄ±
- Proje geneli dokÃ¼mantasyon dili TÃ¼rkÃ§e tutulmalÄ±; dosya baÅŸlÄ±klarÄ±, bÃ¶lÃ¼m adlarÄ± ve tablo kolonlarÄ± tutarlÄ± yazÄ±lmalÄ±dÄ±r.
- ModÃ¼l dokÃ¼manlarÄ± en az `AmaÃ§`, `Sorumluluk AlanÄ±`, `Bu ModÃ¼l Neyi Yapmaz?`, `Veri SahipliÄŸi`, `Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?`, `Access KontratÄ±`, `API veya Event SÄ±nÄ±rÄ±`, `BaÄŸÄ±mlÄ±lÄ±klar`, `Settings Etkileri`, `Event AkÄ±ÅŸlarÄ±`, `Audit ve Ä°zleme`, `Ä°dempotency ve Retry`, `State YapÄ±sÄ±` ve `Test NotlarÄ±` bÃ¶lÃ¼mlerini taÅŸÄ±malÄ±dÄ±r.
- Negatif sÄ±nÄ±r bÃ¶lÃ¼mleri boÅŸ geÃ§ilmemeli; modÃ¼lÃ¼n yapmadÄ±ÄŸÄ± iÅŸler ve sahip olmadÄ±ÄŸÄ± veriler aÃ§Ä±kÃ§a yazÄ±lmalÄ±dÄ±r.
- Ã‡apraz kesit kararlarÄ± iÃ§in canonical shared dokÃ¼manlar `docs/shared/precedence-rules.md`, `docs/shared/projection-strategy.md`, `docs/shared/idempotency-policy.md`, `docs/shared/transaction-boundaries.md`, `docs/shared/audit-policy.md`, `docs/shared/outbox-pattern.md` ve `docs/shared/operational-standards.md` olmalÄ±dÄ±r.
- Teknik paket, cache/queue, media, search ve reporting/analytics kararlarÄ± iÃ§in aktif referanslar `docs/shared/technical-stack.md`, `docs/shared/cache-queue-strategy.md`, `docs/shared/media-asset-strategy.md`, `docs/shared/search-strategy.md` ve `docs/shared/reporting-analytics-strategy.md` olmalÄ±dÄ±r.
- Bir modÃ¼l veya feature dokÃ¼manÄ± bu yardÄ±mcÄ± altyapÄ± alanlarÄ±nda aktif sistem kararÄ± Ã¼retiyorsa ilgili shared dokÃ¼man aynÄ± deÄŸiÅŸiklik setinde gÃ¼ncellenmelidir; karar yalnÄ±zca modÃ¼l iÃ§inde not olarak bÄ±rakÄ±lamaz.
- Enum veya karar sÃ¶zlÃ¼ÄŸÃ¼ niteliÄŸindeki paylaÅŸÄ±lan kayÄ±tlar `docs/shared/target-types.md`, `docs/shared/visibility-states.md`, `docs/shared/moderation-statuses.md`, `docs/shared/support-statuses.md`, `docs/shared/reward-source-types.md`, `docs/shared/purchase-source-types.md`, `docs/shared/audit-event-types.md`, `docs/shared/notification-categories.md` ve `docs/shared/policy-effects.md` iÃ§inde tutulmalÄ±dÄ±r.
- Runtime ayar yorumlama sÄ±rasÄ± `global kill switch -> module/surface availability -> audience selector -> entitlement etkisi -> action policy -> rate limit` biÃ§iminde dokÃ¼mante edilmeli ve `access` tarafÄ±ndan yorumlanan yÃ¼zeyler `docs/settings/index.md` ile hizalÄ± kalmalÄ±dÄ±r.
- `docs/settings/index.md` yaÅŸayan dokÃ¼mandÄ±r; yeni bir surface eklendiÄŸinde availability anahtarÄ± yanÄ±nda rate limit, threshold, disabled behavior veya degrade davranÄ±ÅŸÄ± da aynÄ± deÄŸiÅŸiklikte yazÄ±lmalÄ± ya da neden henÃ¼z `planned` kaldÄ±ÄŸÄ± aÃ§Ä±kÃ§a not edilmelidir.
- `docs/upgrade.md` ham Ã¶neri arÅŸivi olarak deÄŸil, uygulanan, kÄ±smi kalan ve bekleyen iÅŸlerin durumunu izleyen operasyonel takip belgesi olarak kullanÄ±lmalÄ±dÄ±r.
- `support` intake tek baÅŸÄ±na otomatik moderation case sayÄ±lmamalÄ±; `support -> moderation` iliÅŸkisi aÃ§Ä±k handoff politikasÄ± ile tanÄ±mlanmalÄ±dÄ±r.
- `admin` hard override kararÄ± scoped moderator aksiyonunun Ã¼zerinde precedence taÅŸÄ±r; bu kural modÃ¼l dokÃ¼manlarÄ±nda ve precedence matrisi iÃ§inde aynÄ± ÅŸekilde yazÄ±lmalÄ±dÄ±r.
- Projection veya event tabanlÄ± read model kullanan modÃ¼ller canonical write model, rebuild yolu, replay desteÄŸi ve kabul edilen eventual consistency penceresini `docs/shared/projection-strategy.md` ile hizalÄ± yazmalÄ±dÄ±r.
- `payment`, `inventory`, `mission`, `royalpass`, `notification`, `support`, `moderation`, `history` ve benzeri event Ã¼reticisi modÃ¼llerde transactional outbox, retry ve dead-letter yaklaÅŸÄ±mÄ± plan dÄ±ÅŸÄ± bÄ±rakÄ±lmamalÄ±dÄ±r.
- Audit kaydÄ± gereken modÃ¼ller actor, target, action, result, reason, `correlation_id` ve `request_id` alan setini ortak kullanmalÄ±dÄ±r.
- Request ID, correlation ID, rate limit, secret/config ayrÄ±mÄ±, backup/restore/rollback ve PII retention kurallarÄ± `docs/shared/operational-standards.md` ile hizalÄ± kalmalÄ±dÄ±r.
- Test stratejisi, contract test zorunluluklarÄ± ve fixture standardÄ± `docs/TESTING.md` iÃ§inde tutulmalÄ±; modÃ¼l dokÃ¼manlarÄ± buradaki katmanlara referans vermelidir.


## 23) AÅŸama BaÅŸlatma ve Tamamlama Zorunlu AkÄ±ÅŸÄ±
- Bir ajan herhangi bir aÅŸamayÄ± oluÅŸturmaya veya uygulamaya baÅŸladÄ±ÄŸÄ±nda, aÅŸaÄŸÄ±daki adÄ±mlarÄ± sÄ±rasÄ±yla ve eksiksiz uygulamak zorundadÄ±r.
- Bu akÄ±ÅŸ, diÄŸer tÃ¼m genel kurallarÄ±n yanÄ±nda zorunlu operasyonel Ã§alÄ±ÅŸma akÄ±ÅŸÄ± olarak kabul edilmelidir.
- Zorunlu sÄ±ralÄ± akÄ±ÅŸ:
  - Ã–nce kurallarÄ± okur.
  - ArdÄ±ndan proje yapÄ±sÄ±nÄ± inceler.
  - YapÄ±lacak aÅŸamayÄ± docs/ROADMAP.md Ã¼zerinden okur.
  - Ä°lgili aÅŸama iÃ§in uygulanabilir planÄ± hazÄ±rlar.
  - HazÄ±rlanan planÄ± gerÃ§ekleÅŸtirir.
  - AÅŸamaya ait tÃ¼m testleri oluÅŸturur.
  - OluÅŸturulan tÃ¼m testleri baÅŸarÄ±yla tamamlar.
  - Docker build alÄ±r ve uygulamayÄ± Docker iÃ§inde Ã§alÄ±ÅŸtÄ±rÄ±r.
  - Versiyonlama iÅŸlemlerini bu dokÃ¼mandaki sÃ¼rÃ¼m/commit/branch kurallarÄ±na uygun ÅŸekilde uygular.
  - TÃ¼m deÄŸiÅŸiklikleri Git'e yÃ¼kler.
  - Son olarak planÄ±n tamamlandÄ±ÄŸÄ±nÄ± kontrol eder ve kÄ±sa, Ã¶z bir rapor ile sonucu iletir.
## 30) DokÃ¼mantasyon YapÄ±sÄ± ve GÃ¼ncelleme KuralÄ±
- Projede aktif ana dokÃ¼mantasyon seti aÅŸaÄŸÄ±daki beÅŸ dosyadan oluÅŸmalÄ±dÄ±r:
  - `docs/rules.md`
  - `docs/roadmap.md`
  - `docs/changelog.md`
  - `docs/modules.md`
  - `docs/shared.md`
- AyrÄ± modÃ¼l ve shared alt dokÃ¼manlarÄ± yalnÄ±zca Ã§alÄ±ÅŸma taslaÄŸÄ± olarak tutulabilir; aktif referans seti yukarÄ±daki beÅŸ dosyadÄ±r.
- `rules.md` proje geneli baÄŸlayÄ±cÄ± kurallarÄ± taÅŸÄ±r.
- `roadmap.md` sistemin oluÅŸturulma sÄ±rasÄ±nÄ±, fazlarÄ±nÄ± ve teslim sÄ±nÄ±rlarÄ±nÄ± taÅŸÄ±r.
- `changelog.md` yalnÄ±zca projede gerÃ§ekten yapÄ±lan deÄŸiÅŸiklikleri kronolojik olarak kaydeder; plan, niyet veya gelecekte yapÄ±lacak iÅŸler changelog'a yazÄ±lmaz.
- `modules.md` tÃ¼m modÃ¼l ownerlÄ±klarÄ±, sÄ±nÄ±rlar, veri sahipliÄŸi, baÄŸÄ±mlÄ±lÄ±klar ve modÃ¼l bazlÄ± aÃ§Ä±klamalar iÃ§in tek referans dosyadÄ±r.
- `shared.md` shared sÃ¶zlÃ¼kler, ortak teknik kararlar, settings envanteri ve modÃ¼ller Ã¼stÃ¼ politikalar iÃ§in tek referans dosyadÄ±r.
- Yeni modÃ¼l eklendiÄŸinde veya modÃ¼l ownerlÄ±ÄŸÄ± deÄŸiÅŸtiÄŸinde aynÄ± deÄŸiÅŸiklikte `modules.md` gÃ¼ncellenmelidir.
- Ortak enum, policy, precedence, ayar anahtarÄ± veya operasyon standardÄ± eklendiÄŸinde aynÄ± deÄŸiÅŸiklikte `shared.md` gÃ¼ncellenmelidir.
- Faz, sÄ±ra veya kapsam deÄŸiÅŸtiÄŸinde aynÄ± deÄŸiÅŸiklikte `roadmap.md` gÃ¼ncellenmelidir.
- Mimariyi, klasÃ¶r yapÄ±sÄ±nÄ±, katman sÄ±nÄ±rlarÄ±nÄ± veya geliÅŸtirme kurallarÄ±nÄ± etkileyen deÄŸiÅŸikliklerde `rules.md` gÃ¼ncellenmelidir.
- Proje oluÅŸturulurken veya geliÅŸtirme sÃ¼recinde gerÃ§ekten yapÄ±lan iÅŸlemler, eklenen dosyalar, Ã§Ä±karÄ±lan dosyalar ve Ã¶nemli kararlar `changelog.md` iÃ§inde kayÄ±t altÄ±na alÄ±nmalÄ±dÄ±r.
- DokÃ¼man gÃ¼ncellemesi olmadan yapÄ±lan mimari deÄŸiÅŸiklik tamamlanmÄ±ÅŸ sayÄ±lmamalÄ±dÄ±r.

# Kurallar

> Bu doküman, proje geneli kuralları, ölçeklenebilir mimari ilkeleri ve ortak çalışma standartlarını toplar. Modül ve feature detayları ayrı dokümanlarda genişletilir; ancak burada tanımlanan genel kurallar tüm modüller için bağlayıcıdır.

## 1) Proje Kimliği ve Kapsam
- Proje adı: NovaScans.
- Proje alanı: oyunlaştırılmış manga, manhwa ve manhua okuma platformu.
- Bu doküman proje geneli kuralları, mimari ilkeleri ve çalışma standartlarını tanımlar.
- Bu doküman sprint planı değildir.
- Bu doküman görev listesi değildir.
- Bu doküman tek tek modül implementasyon dokümanı değildir.
- Detaylı implementasyon kararları ve modül içi feature ayrıntıları sonraki iterasyonlarda ayrı dokümanlarda tanımlanmalıdır.
- Genel mimari kararlar korunmadan yapılan geliştirme tamamlanmış sayılmaz.

## 2) Dokümanın Amacı ve Kullanımı
- Bu doküman; proje sahibi, geliştiriciler, AI destekli araçlar ve ajanlar için bağlayıcıdır.
- Ana amaç, proje büyürken mimariyi, sınırları ve sürdürülebilirliği korumaktır.
- Yeni geliştirme, refactor, veri modeli değişikliği, API değişikliği, altyapı değişikliği ve süreç tasarımında bu doküman esas alınmalıdır.
- Projede modül, özellik, hotfix veya herhangi bir düzenleme yapılacağı zaman önce aktif kurallar dokümanı baz alınmalıdır.
- Genel kurallar ile daha detaylı modül dokümanları çelişirse önce bu dokümandaki genel mimari ilkeler korunmalıdır.
- Yeni ana kararlar dokümana yansıtılmadan iş tamamlanmış kabul edilmemelidir.
- Hızlı çözüm, geçici çözüm veya acil düzeltme gerekçesi bu dokümandaki genel mimari kuralları kalıcı olarak ihlal etme nedeni olamaz.

## 3) Sabit Teknik Kararlar
- Backend dili: Go 1.26.
- Canonical env/config loader olarak `caarlos0/env` kullanılmalıdır.
- SQL erişimi ve connection pooling için `pgx/v5` ve `pgxpool` kullanılmalıdır.
- Structured logging için canonical seçim `zap` olmalıdır.
- Input validation için canonical seçim `go-playground/validator/v10` olmalıdır.
- UUID üretimi için canonical seçim `google/uuid` olmalıdır.
- Password hashing için canonical seçim `argon2id` olmalıdır.
- Test assertion ve helper standardı için `testify` kullanılmalıdır.
- Başlangıç async işleme standardı PostgreSQL-backed jobs + transactional outbox olmalıdır.
- Cache ihtiyacı oluştuğunda canonical backend `Redis` olmalıdır; cache source-of-truth kabul edilmemelidir.
- Veritabanı: PostgreSQL 18.
- HTTP router olarak Chi kullanılmalıdır.
- Sürüm kontrol sistemi: Git.
- Commit mesajları Conventional Commits standardına uymalıdır.
- Branch modeli: `main + feature/* + hotfix/*` (ajan calismalari icin `codex/**` gecici branch modeli desteklenebilir).
- Migration yönetiminde `golang-migrate` kullanılmalıdır.
- Yapılandırma env tabanlı olmalı, config erişimi merkezi katmandan yapılmalıdır.
- Proje Docker içinde build alabilmelidir.
- Proje Docker içinde ayağa kalkabilmelidir.
- Main DB ve test DB kesin olarak ayrılmalıdır.

## 4) Geliştirme Prensipleri
- Her değişiklik doğru sorumluluk alanında yapılmalıdır.
- Gereksiz refactor yapılmamalıdır.
- Gereksiz dosya taşıma veya isim değişikliği yapılmamalıdır.
- Gereksiz paket, feature veya soyutlama eklenmemelidir.
- Yeni yapı eklenmeden önce mevcut yapı incelenmelidir.
- Aynı sorumluluk birden fazla alana dağılmamalıdır.
- Ortaklaştırma yalnızca gerçek ihtiyaç oluştuğunda yapılmalıdır.
- Geçici çözümler açıkça işaretlenmeli, kalıcı mimari karar gibi bırakılmamalıdır.
- Kod okunabilir, sürdürülebilir, test edilebilir ve izlenebilir olmalıdır.

## 5) Proje Yapısı ve Dizin Organizasyonu
- Proje çok modüllü büyümeye ve ileride frontend eklenmesine uygun repo kök standardı ile başlamalıdır.
- Repo kökü yalnızca üst seviye uygulama dizinleri, ortak dokümantasyon, ortak scriptler, deploy dosyaları ve repo seviyesi yapılandırmaları içermelidir.
- Backend ve frontend aynı repo içinde yer alacaksa uygulamalar `apps/` altında ayrıştırılmalıdır.
- Önerilen kök dizin yapısı aşağıdaki omurgayı korumalıdır:

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

- `apps/api/` Go backend uygulamasının tek canonical kökü olmalıdır.
- `apps/web/` frontend uygulamasının tek canonical kökü olmalıdır.
- `apps/api/cmd/` yalnızca backend uygulama giriş noktaları için kullanılmalıdır.
- `apps/api/internal/app/` uygulama bootstrap, composition root, dependency wiring ve merkezi başlatma akışları için kullanılmalıdır.
- `apps/api/internal/platform/` config, DB, logger, middleware, mailer, cache, queue, storage ve benzeri teknik altyapı kodları için kullanılmalıdır.
- `apps/api/internal/shared/` yalnızca gerçekten modülden bağımsız, domain-agnostic ve tekrar kullanılabilir yapılar için kullanılmalıdır.
- `apps/api/internal/modules/` tüm backend iş modüllerinin ana yerleşim alanı olmalıdır.
- `apps/api/migrations/` veritabanı migration dosyaları için tek merkez olmalıdır.
- `apps/api/tests/` entegrasyon, contract veya uçtan uca testlerin üst seviye yerleşimi için kullanılabilir.
- `docs/` tüm mimari, süreç ve modül belgeleri için tek merkez olmalıdır.
- `scripts/` ortak geliştirme ve bakım scriptleri için kullanılmalıdır.
- `deploy/` Docker, Compose, deployment ve operasyonel çalışma dosyaları için kullanılmalıdır.
- `.github/` CI/CD workflow, issue template ve PR template dosyaları için kullanılmalıdır.
- Repo kökünde backend uygulama dosyaları dağınık şekilde tutulmamalı; backend kodu `apps/api/` altında toplanmalıdır.
- Repo kökünde frontend uygulama dosyaları dağınık şekilde tutulmamalı; frontend kodu `apps/web/` altında toplanmalıdır.

## 6) Modül Organizasyonu ve Ölçeklenme Kuralları
- Varsayılan backend modül kök dizini `apps/api/internal/modules/<module>/` olmalıdır.
- Bir modülün gerçek kök dizini, modülün kendi leaf klasörüdür.
- Modül sayısı arttığında veya okuma/bakım maliyeti yükseldiğinde opsiyonel domain grubu kullanılabilir.
- Domain grubu kullanılan durumda yapı `apps/api/internal/modules/<domain-group>/<module>/` formatına geçebilir.
- Domain group klasörü yalnızca gruplayıcıdır; gerçek iş sınırı yine leaf modül klasörüdür.
- Başlangıç aşamasında gereksiz klasör derinliği oluşturulmamalıdır; gerçek ihtiyaç yoksa `apps/api/internal/modules/<module>/` yeterlidir.
- Domain group kullanımı bir gereklilik değil, ölçeklenme aracıdır.
- Backend modülleri yalnızca `apps/api/internal/modules/` altında yer almalıdır.
- Frontend tarafındaki feature veya page organizasyonu backend modül kök yapısı ile karıştırılmamalıdır.
- Örnek domain group alanları bağlayıcı olmadan şu şekilde düşünülebilir:
  - `identity`
  - `content`
  - `community`
  - `operations`
  - `engagement`
  - `commerce`
  - `gameplay`
- Yeni modül açmak için açık veri sahipliği, açık use-case sınırı ve net bağımlılık gerekçesi bulunmalıdır.
- Alt özellikler varsayılan olarak ayrı modül yapılmamalı, önce mevcut modül içinde kalmalıdır.
- Bir özellik ancak bağımsız veri sahipliği, bağımsız servis akışı ve bağımsız access kontratı gerektiriyorsa ayrı modüle ayrılmalıdır.
- Modül isimleri kısa, tek anlamlı ve iş alanını yansıtan canonical adlar olmalıdır.
- Aynı anlama gelen birden fazla modül adı açılmamalıdır.
- Modül adları tekil veya çoğul kullanım açısından tutarlı olmalı; aynı alan için iki farklı yazım standardı açılmamalıdır.
- Tüm aktif modüller için canonical modül kaydı tutulmalıdır.
- Önerilen kayıt dosyası `docs/modules.md` veya benzeri merkezi bir modül envanteri olmalıdır.
- Bu kayıt en az şu alanları içermelidir:
  - canonical modül adı
  - varsa domain group
  - kısa açıklama
  - durum
  - ana doküman yolu
- `durum` alanı için önerilen canonical değerler: `planned`, `active`, `deprecated`, `archived`.

## 7) Standart Modül Yapısı ve Katman İlkeleri
- Her leaf modül, yalnızca ihtiyaç duyduğu dosya ve klasörleri içermelidir; ancak kullanılan yapılar ortak isimlendirme ve katman standardını korumalıdır.
- Aşağıdaki yapı zorunlu tam şablon değildir; modülde ihtiyaç doğduğunda hangi klasör veya dosyanın hangi amaçla açılacağını gösteren standart backend omurgadır:

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

- Domain group kullanılıyorsa aynı mantık `apps/api/internal/modules/<domain-group>/<module>/` altında uygulanmalıdır.

- `entity/` modülün kendi veri yapıları için kullanılmalıdır.
- `dto/` request, response ve modül dışı contract veri yapıları için kullanılmalıdır.
- `service/` iş kuralları ve use-case akışları için kullanılmalıdır.
- `repository/` yalnızca veri erişimi için kullanılmalıdır.
- `handler/` HTTP handler, request parse ve response üretimi gibi giriş katmanı sorumlulukları için kullanılmalıdır.
- `middleware/` modüle özel middleware yapıları için kullanılabilir.
- `validator/` modüle özel doğrulama kuralları ve input validation yardımcıları için kullanılabilir.
- `mapper/` entity, dto, response, read model veya contract dönüşümleri karmaşık hale geldiğinde kullanılabilir.
- `contract/` diğer modüllerle paylaşılan resmi modül kontratları için kullanılmalıdır.
- `events/` modülün yayınladığı veya tükettiği event tanımları için kullanılabilir.
- `consumer/` queue consumer, event consumer, webhook consumer veya dış giriş entegrasyonları için kullanılabilir.
- `producer/` message producer, event producer veya dış sisteme yayın yapan entegrasyon akışları için kullanılabilir.
- `jobs/` asenkron işleyiciler veya zamanlanmış işler için kullanılabilir.
- `readmodel/` yalnızca okuma odaklı özel projection veya denormalize yapılar gerektiğinde kullanılmalıdır.
- `module.go` gerekiyorsa modülün composition veya registration giriş noktası olmalıdır.
- `routes.go` gerekiyorsa modülün route kayıt giriş noktası olmalıdır.
- Bu bölümde listelenen klasör ve dosyaların hiçbiri her modülde zorunlu değildir; zorunlu olan, ihtiyaç doğduğunda burada tanımlanan amaç ve isim standardına uyulmasıdır.
- Tekil ve modül seviyesinde kalan dosyalar modül kök dizininde tutulmalıdır.
- Örnek tekil modül kök dosyaları: `module.go`, `routes.go`, `errors.go`, `constants.go`, `types.go`, `service.go`, `repository.go`, `handler.go`, `middleware.go`, `validator.go`, `mapper.go`, `contract.go`, `events.go`, `consumer.go`, `producer.go`, `jobs.go`, `readmodel.go`.
- Tek bir katman veya modül bileşeni tek dosya ile temsil edilebiliyorsa ilgili dosya modül kökünde tutulabilir.
- Bu kök dosyalar yalnızca tek bir açık sorumluluk veya tek bir akış ailesi taşıdığı sürece kökte kalmalıdır.
- Bir alan birden fazla işlem, akış veya dosya gerektiriyorsa modül kökünde büyütülmemeli, aynı amaçla açılan klasör altında parçalanmalıdır.
- Örnek dönüşüm: `service.go -> service/`, `handler.go -> handler/`, `middleware.go -> middleware/`, `routes.go -> routes/`, `events.go -> events/`, `consumer.go -> consumer/`.
- `errors.go`, `constants.go`, `types.go` veya benzeri tekil kök dosyalar tekil olmaktan çıkarsa daha açık isimli dosyalara veya uygun klasör yapısına ayrıştırılmalıdır.
- Bu kural yalnızca belirli birkaç klasör için değil, modül içindeki tüm çoklu işlev alanları için geçerlidir.
- `service/`, `repository/`, `dto/`, `entity/`, `handler/`, `middleware/`, `validator/`, `mapper/`, `contract/`, `events/`, `consumer/`, `producer/`, `jobs/`, `readmodel/` ve gelecekte açılacak benzer çoklu işlev klasörleri aynı parçalama ilkesine uymalıdır.
- Modül kökü tekil dosyalar içindir; çoklu işlev taşıyan büyük dosyalar kökte biriktirilmemelidir.
- Birden fazla use-case veya işlem içeren modüllerde her ana işlem mümkün olduğunda ayrı dosyada tutulmalıdır.
- Service dosyaları işlem veya use-case bazlı parçalanmalıdır.
- Repository dosyaları mümkün olduğunda entity, aggregate veya belirgin veri erişim sorumluluğu bazlı parçalanmalıdır.
- DTO dosyaları request ve response olarak ayrılmalı; birden fazla farklı akış tek DTO dosyasına doldurulmamalıdır.
- Entity dosyaları tek bir dev entity dosyasına dönüşmemeli; alan kümeleri, alt yapılar veya anlamlı domain ayrımları varsa kontrollü şekilde bölünmelidir.
- Handler dosyaları mümkün olduğunda endpoint veya işlem ailesi bazında parçalanmalıdır.
- Middleware dosyaları farklı middleware akışlarını tek dosyada gereksiz şekilde biriktirmemelidir.
- Validator dosyaları farklı işlem ailelerini tek doğrulama dosyasında toplamamalıdır.
- Mapper dosyaları farklı dönüşüm ailelerini tek dev dosyada toplamamalıdır.
- Contract dosyaları farklı entegrasyon akışlarını tek dosyada biriktirmemeli; entegrasyon yüzeyleri gerektiğinde işlem ailesi bazında ayrılmalıdır.
- Event dosyaları producer veya consumer tarafında farklı domain olaylarını anlamsız şekilde tek dosyada biriktirmemelidir.
- Consumer dosyaları farklı dış giriş akışlarını tek dosyada biriktirmemelidir.
- Producer dosyaları farklı dış yayın akışlarını tek dosyada biriktirmemelidir.
- Job dosyaları zamanlanmış işler veya asenkron akışlar bazında ayrıştırılmalıdır.
- Read model dosyaları tek bir dev projection dosyasına dönüşmemelidir.
- Modül kökündeki tekil dosyalar yalnızca kendi tekil modül sorumluluğunu taşımalıdır.
- Eğer route, middleware, error mapping veya registration alanı tekil olmaktan çıkıp birden fazla akış taşımaya başlarsa bu yapı kökte büyütülmemeli, uygun klasör altında ayrıştırılmalıdır.
- Aynı modülde 5-6 farklı ana akışı tek `service.go`, tek `repository.go`, tek `dto.go` veya benzeri tek bir dosyada toplamak mimari olarak yanlış kabul edilmelidir.
- Dosya parçalama keyfi değil, okunabilirlik ve sorumluluk ayrımı amacıyla yapılmalıdır; anlamsız aşırı bölme de yapılmamalıdır.
- Dosya adlarında modül prefixi ve işlem adı birlikte kullanılmalıdır.
- Örnek yaklaşım:

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

- Bu örnek bağlayıcı isim listesi değildir; bağlayıcı olan ilke, dosyanın tek bir ana sorumluluğu ve açık bir adı olmasıdır.
- Yeni işlem eklendiğinde varsayılan yaklaşım mevcut büyük dosyayı daha da büyütmek değil, modül içindeki doğru klasörde veya gerekiyorsa modül kökünde yeni işlem dosyası açmak olmalıdır.
- Refactor sırasında devasa dosyalar görülürse önce işlem sınırları ayrıştırılmalı, sonra dosyalar kontrollü şekilde bölünmelidir.
- Modül dizin adları ve katman klasörleri küçük harfli, ASCII uyumlu ve tek bir canonical yazım ile oluşturulmalıdır.
- Dizin, modül ve klasör adlarında boşluk, tire veya rastgele kısaltma kullanılmamalıdır.
- Go package adları küçük harfli, kısa, tek anlamlı ve mümkün olduğunda klasör adı ile uyumlu olmalıdır.
- Go package adlarında anlamsız kısaltmalar, underscore veya birden fazla kelimeyi gereksiz şekilde birleştiren karmaşık adlar kullanılmamalıdır.
- Dosya adları küçük harfli `snake_case` biçiminde, açık sorumluluk ve açık işlem adı ile oluşturulmalıdır.
- `util.go`, `utils.go`, `helper.go`, `common.go`, `misc.go`, `temp.go` gibi belirsiz dosya adları kullanılmamalıdır.
- Export edilen Go identifier adları `PascalCase`, export edilmeyen adlar `camelCase` standardına uymalıdır.
- Kısaltmalar tutarlı kullanılmalıdır: `ID`, `API`, `HTTP`, `URL` gibi yaygın kısaltmalar proje genelinde aynı biçimde yazılmalıdır.
- Katmanlar arası bağımlılık yönü kontrollü olmalıdır.
- HTTP katmanı iş kararı vermemelidir.
- Service katmanı HTTP detayına bağımlı olmamalıdır.
- Repository katmanı business karar vermemelidir.
- Entity yapıları başka modülün repository detayına bağımlı olmamalıdır.
- Frontend klasör yapısı bu backend modül standardına göre şekillendirilmemelidir.
- Backend modül standardı yalnızca `apps/api/` içindeki Go uygulaması için geçerlidir.

## 8) Shared, Platform ve Bootstrap Sınırları
- `platform` teknik altyapının sahibidir; iş kuralı barındırmamalıdır.
- `shared` gerçekten ortak ve domain-agnostic yapıların sahibidir; modül bazlı iş mantığı içeremez.
- `shared` klasörü, modüllerden kod kaçırma alanına dönüşmemelidir.
- Bir kod parçası yalnızca iki farklı modül tarafından tekrar kullanıldığı için otomatik olarak `shared` içine taşınmamalıdır.
- `app` veya bootstrap katmanı modülleri birbirine bağlar; modüller bootstrap katmanını import etmemelidir.
- Modül bağımlılıkları merkezi wiring katmanında birleştirilmelidir.
- Route mount işlemleri merkezi uygulama başlangıç katmanında yapılmalıdır.

## 9) Çapraz Modül Bağımlılık ve Entegrasyon Kuralları
- Bir modül başka modülün tablolarına doğrudan yazmamalıdır.
- Bir modül başka modülün repository implementasyonunu doğrudan kullanmamalıdır.
- Bir modül başka modülün iç klasörlerine doğrudan bağımlı olmamalıdır.
- Başka bir modülün `handler/`, `repository/`, `entity/`, `jobs/`, `consumer/`, `producer/` veya benzeri iç yapıları dışarıya açık API kabul edilmemelidir.
- Her modül kendi public surface'inin tek sahibidir.
- Bir modülün public surface'i yalnızca açıkça belgelediği ve dışa açtığı şu yüzeylerden oluşmalıdır:
  - `contract/`
  - dış kullanıma açılmış service interface'leri
  - dış kullanım için tanımlanmış DTO veya contract modelleri
  - yayınlanan event sözleşmeleri
  - gerekiyorsa açık read model veya projection yüzeyi
- Bu yüzey dışında kalan tüm yapı modülün iç implementasyonudur ve dış bağımlılık noktası olarak kullanılamaz.
- Senkron service veya contract sözleşmesinin sahipliği provider modüldedir.
- Event şemasının sahipliği event'i publish eden producer modüldedir.
- Read model veya projection yüzeyinin sahipliği ilgili veriyi servis eden modüldedir.
- Consumer modül kendi içinde local interface veya adapter tanımlayabilir; ancak bu durum provider modülün iç klasörlerini resmi public API'ye dönüştürmez.
- Modüller arası iletişim yalnızca aşağıdaki yollardan biriyle kurulmalıdır:
  - açık service/contract arayüzü
  - açık DTO/contract modeli
  - açık event sözleşmesi
  - açık read model veya projection ihtiyacı
- Çapraz modül veri ihtiyacı için diğer modülün entity yapısı doğrudan dışarı sızdırılmamalıdır.
- Bir modülün taşıdığı denormalize sayaç, özet veya projection alanı o modülde tutulabilir; ancak bu alanların canonical kaynak verisi ilgili kaynak modülde kalmalıdır.
- Çapraz modül sayaç veya özet güncellemeleri kaynak modül tarafından owner modülün tablosuna doğrudan yazılarak yapılmamalıdır; açık event, projection veya owner modülün açık counter contract yüzeyi kullanılmalıdır.
- Sayaç veya özet alanı tanımlanan her yerde canonical source, güncelleme tetikleyicisi, kabul edilen gecikme modeli ve gerektiğinde reconcile veya yeniden hesaplama yolu dokümante edilmelidir.
- `target_type` veya benzeri paylaşılan hedef tipleri taşıyan modüller canonical kayıt dosyası olarak `docs/shared.md` kullanmalıdır.
- Yeni target type yalnızca ilgili modül dokümanı, consumer modül dokümanları ve `docs/shared.md` aynı değişiklik setinde güncellendiğinde kullanılabilir hale gelmelidir.
- Döngüsel modül bağımlılığı kesin olarak yasaktır.
- Senkron bağımlılık yalnızca gerçekten anlık doğrulama veya kritik işlem bütünlüğü gerektiğinde kullanılmalıdır.
- Zayıf bağlı entegrasyonlarda event tabanlı veya asenkron akış tercih edilmelidir.
- Event tabanlı entegrasyon kullanılan yerlerde event adı, payload, producer, consumer ve idempotency beklentisi dokümante edilmelidir.
- Public surface üzerinde yapılan her değişiklik modül dokümanına, gerekiyorsa changelog kaydına ve versiyonlama değerlendirmesine yansıtılmalıdır.

## 10) Yeni Modül Açma Zorunlu Kriterleri
- Yeni modül açılmadan önce en az aşağıdaki başlıklar tanımlanmalıdır:
  - sorumluluk alanı
  - veri sahipliği
  - dışa açılan API veya contract sınırı
  - access veya authorization entegrasyonu
  - state ve lifecycle yapısı
  - diğer modüllerle bağımlılık ilişkisi
  - event ihtiyacı varsa event sözleşmesi
  - shared registry etkisi varsa ilgili canonical kayıt dosyası
  - config ihtiyacı
  - migration ihtiyacı
  - test gereksinimleri
  - log ve audit gereksinimleri
  - dokümantasyon dosyası
- Bu başlıklar netleşmeden yeni modül implementasyonu başlatılmamalıdır.
- Her yeni modül için en az bir modül dokümanı oluşturulmalıdır.
- Domain group kullanılmıyorsa modül dokümanları için önerilen yerleşim `docs/modules.md` olmalıdır.
- Domain group kullanılıyorsa modül dokümanları için önerilen yerleşim `docs/modules.md` olmalıdır.
- Yeni leaf modül açıldığında merkezi modül envanteri aynı değişiklik seti içinde eklenmeli veya güncellenmelidir.

## 11) Veri, API ve State Kuralları
- Şema ilişkileri temiz, açık ve tutarlı kurulmalıdır.
- Foreign key mantığı net olmalıdır.
- Soft delete kullanılan alanlar kontrollü uygulanmalıdır.
- Unique alanlar açık tanımlanmalı ve çakışma senaryoları yönetilmelidir.
- Gerekli alanlarda makul index kullanılmalıdır.
- Veri modeli gereksiz karmaşıklık oluşturmadan büyümeye uygun olmalıdır.
- API stili REST + JSON olmalıdır.
- Request ve response modelleri açıkça ayrılmalıdır.
- DB modeli doğrudan API response olarak kullanılmamalıdır.
- Standart hata cevabı formatı kullanılmalıdır.
- `401` yalnızca kimlik doğrulama eksik veya geçersiz olduğunda kullanılmalıdır.
- `403` yetki, policy veya access kararı nedeniyle reddedilen işlemlerde kullanılmalıdır.
- `404` görünürlük veya kaynak gizleme kararı nedeniyle dışarıya kapalı bırakılan kaynaklarda kullanılmalıdır.
- `429` rate limit, throttling, cooldown veya benzeri geçici eşik ihlali durumlarında kullanılmalıdır.
- `503` bakım modu, kill switch, emergency deny veya sistem kaynaklı geçici pasiflik durumlarında kullanılmalıdır.
- Liste endpoint'lerinde pagination, filter ve sort parametreleri tutarlı adlandırılmalıdır.
- State, visibility, moderation ve publish kavramları birbirine karıştırılmamalıdır.
- State değişiklikleri kontrollü, izlenebilir ve yetki kontrolü altında olmalıdır.
- Admin tarafından yönetilen runtime ayarlar, modül durumu ve özellik durumu veri modeli seviyesinde açık scope ile temsil edilebilmelidir.
- Sistem en az site geneli, modül bazlı, özellik bazlı ve gerektiğinde context veya resource bazlı açma-kapama davranışlarını destekleyebilmelidir.
- Bir modül veya alt özellik pasife alındığında beklenen fallback davranışı, görünürlük etkisi ve hata cevabı açıkça tanımlanmalıdır.
- İş kurallarını etkileyen eşik, oran, süre ve limit değerleri mümkün olduğunda sabit koda gömülmemeli; yönetilebilir ayar yüzeyi ile kontrol edilebilmelidir.
- Runtime ayarlarda `scope`, ayarın nerede uygulandığını; `audience`, ayarın kim için uygulandığını ifade eder. Bu iki kavram birbirine karıştırılmamalıdır.
- Başlangıç audience kapsamı en az `all`, `guest`, `authenticated`, `authenticated_non_vip` ve `vip` seviyelerini destekleyebilmelidir.
- İleri audience genişletmeleri `role:<name>`, `group:<name>` veya `user:<id>` gibi hedeflemeler olabilir; bu tür kapsamlar açıkça dokümana eklenmeden kullanılmamalıdır.
- Runtime kontrol modeli tek bir global `disabled` bayrağına indirgenmemelidir; gerektiğinde `read`, `write`, `intake`, `preview`, `visibility` ve `benefit` gibi yüzeyler ayrı ayrı yönetilebilmelidir.
- Runtime kapatma davranışları canonical olarak en az `visibility_off`, `read_only`, `write_off`, `intake_pause`, `preview_off` ve `benefit_pause` tiplerini ifade edebilmelidir.
- `disabled_behavior` işlevsel davranışı ifade eder; tek başına HTTP cevap kodunu belirlemez. API davranışı için ayrıca `error_response_policy` tanımlanmalıdır.
- `error_response_policy` en az `not_found`, `forbidden`, `rate_limited` ve `temporarily_unavailable` değerlerini desteklemelidir.
- Varsayılan hizalama olarak görünürlük gizleme kararları `not_found`, eşik veya cooldown ihlalleri `rate_limited`, sistem kaynaklı pause veya kill switch kararları `temporarily_unavailable`, kullanıcıya açık ama yazma veya erişim kısıtı taşıyan yüzeyler ise `forbidden` ile modellenmelidir.
- Availability veya kill switch türündeki ayarlarda güvenlik önceliklidir; eşleşen bir `emergency_deny` her durumda en yüksek önceliğe sahip olmalıdır.
- Availability veya kill switch türündeki ayarlarda eşleşen herhangi bir `deny/off` kuralı, `allow/on` kuralını bastırmalıdır.
- Aynı availability anahtarı için aynı `audience_kind + audience_selector` ve aynı `scope_kind + scope_selector` kombinasyonunda birden fazla aktif kural bırakılamaz; bu tür çakışmalar kayıt aşamasında reddedilmelidir.
- Eşik veya değer taşıyan runtime ayarlarda en spesifik geçerli kayıt kullanılmalıdır.
- Eşik veya değer ayarlarında audience özgüllük sırası `user/group/role` > `vip/authenticated_non_vip` > `authenticated/guest` > `all` olmalıdır.
- Eşik veya değer ayarlarında scope özgüllük sırası `resource/context` > `feature` > `module` > `site` olmalıdır.
- Bir modül dokümanı bir yüzeyin ayrı ayrı runtime kontrol edilebildiğini söylüyorsa `docs/shared.md` içinde o yüzey için en az bir canonical baseline key veya bu alt yüzeyleri açıkça kapsayan umbrella key kaydı bulunmalıdır.
- Umbrella key kullanılan durumda kapsanan alt yüzeyler `affected_surfaces` ve `notes` alanlarında açıkça listelenmeli; yüzey dokümanda tanımlı kalırken settings envanterinde tamamen isimsiz bırakılamaz.
- Ücretli veya süreli hakları etkileyen runtime kapatmalar açık bir entitlement impact policy taşımak zorundadır.
- Sistem kaynaklı pasiflikte ücretli veya süreli avantajların kalan süresi sessizce tüketilemez; varsayılan güvenli davranış sürenin dondurulması ve sistem tekrar açıldığında kaldığı yerden devam etmesidir.
- Çapraz modül precedence kararlarında sistem veya admin kaynaklı emergency deny en yüksek önceliktedir; bunun altında `access` deny kararı, modül içi allow veya paylaşım sinyalini her zaman bastırmalıdır.
- `user` modülündeki global görünürlük veya paylaşım preference sinyali ilgili yüzey için üst sınırı tanımlar; `history` içindeki entry-level paylaşım metadata'sı bu üst sınırı daraltabilir veya yalnızca izin verilen tavan içinde opt-in paylaşım sağlayabilir, ancak global deny kararını genişletemez.
- `social` modülünün ürettiği block, privacy veya mute sinyalleri ham ilişki verisidir; final allow veya deny kararı `access` tarafından verilir. Block veya açık privacy deny sinyali final deny üretmelidir; mute sinyali ise aksi ayrıca dokümante edilmedikçe tek başına genel authorization deny sayılmamalıdır.
- `moderation` günlük scoped vaka akışının sahibidir; ancak aynı case üzerinde `admin` tarafından verilen override, reopen, freeze, reassignment veya final kararlar moderator aksiyonunun üzerinde precedence taşır ve yeni bir handoff kaydı oluşmadan moderator tarafından bastırılamaz.
- `support` içindeki report kaydı ile `moderation` case yaşam döngüsü aynı şey sayılmamalıdır; her report zorunlu olarak case açmaz. Moderation incelemesi gerektiğinde support kaydı kaynak intake olarak kalır, moderation tarafında ise buna bağlı ama ayrı bir case lifecycle başlatılır.

## 12) Güvenlik, Loglama ve Audit Kuralları
- Loglar yapılandırılmış formatta üretilmelidir.
- JSON log formatı tercih edilmelidir.
- Her request için `request_id` üretilmeli veya taşınmalıdır.
- Hassas veri loglara yazılmamalıdır.
- Güvenlik olayları gerektiğinde ayrı izlenebilmelidir.
- Audit log ile operasyonel log birbirine karıştırılmamalıdır.
- Yüksek riskli işlemler izlenebilir ve açıklanabilir olmalıdır.
- Admin tarafından yapılan runtime ayar, modül açma-kapama, özellik açma-kapama ve eşik güncelleme işlemleri actor, reason, scope, eski değer ve yeni değer ile audit kaydı üretmelidir.
- Emergency deny veya kill switch işlemleri ayrı kritik operasyon olayı olarak izlenebilmelidir.
- Çapraz modül kritik işlemler gerektiğinde ayrı audit veya domain event kaydı üretmelidir.

## 13) Config ve Ortam Kuralları
- Config değerleri ortam değişkenlerinden okunmalıdır.
- Config erişimi merkezi config yapısından yapılmalıdır.
- Modül bazlı config alanları namespace mantığı ile adlandırılmalıdır.
- Örnek yaklaşım: `AUTH_`, `MANGA_`, `PAYMENT_`, `NOTIFICATION_` gibi modül prefixleri kullanılabilir.
- Yalnızca `.env.example` repoda tutulmalıdır.
- Gerçek secret değerleri repoya commit edilmemelidir.
- Ortam profilleri açıkça ayrılmalıdır: `local`, `test`, `staging`, `prod`.
- Main DB ve test DB config seviyesinde kesin olarak ayrılmalıdır.
- Ortam bağımlı değerler kod içine gömülmemelidir.
- Eksik veya geçersiz config değerleri kontrollü şekilde doğrulanmalıdır.
- Ortam değişkenleri deploy veya çalışma ortamı seviyesindeki teknik config içindir; admin tarafından değiştirilebilen runtime ayarlar env config ile karıştırılmamalıdır.
- Admin tarafından yönetilen runtime ayarlar merkezi ve kalıcı bir ayar deposunda tutulmalı; uygulama yeniden build edilmeden değiştirilebilir olmalıdır.
- Runtime ayar anahtarları canonical namespace ile tanımlanmalıdır; örnek yaklaşım `site.maintenance.enabled`, `auth.login.failed_attempt_limit_per_minute`, `comment.write.cooldown_seconds`, `feature.user.vip_benefits.enabled`.
- Boolean availability veya feature toggle anahtarları mümkün olduğunda `feature.<module>.<surface>.enabled` biçimini kullanmalıdır.
- Eşik, limit, cooldown veya davranış değeri taşıyan anahtarlar mümkün olduğunda `<module>.<surface>.<metric>` biçimini kullanmalıdır.
- Site geneli operasyon veya bakım anahtarları mümkün olduğunda `site.<surface>.<metric_or_flag>` biçimini kullanmalıdır.
- Audience, role, grup, kullanıcı veya resource bilgisi runtime key içine gömülmemeli; bu bilgiler `scope_selector` ve `audience_selector` alanlarında taşınmalıdır.
- Runtime key ve selector içindeki modül adı her zaman canonical leaf modül adı ile aynı yazılmalıdır.
- Runtime ayarlar en az şu canonical kategorilere genişleyebilir olmalıdır: `site`, `communication`, `operations`, `security_auth`, `access_availability`, `content`, `reading`, `engagement`, `support`, `membership`, `social`, `gamification` ve `economy`.
- Runtime ayarlar tip, aralık, zorunluluk ve scope açısından doğrulanmalı; geçersiz ayar değişikliği sessizce kabul edilmemelidir.
- Teknik altyapı config'i hiçbir koşulda admin runtime ayarı sayılmamalıdır; DB host veya pool, SMTP credential, queue DSN, object storage anahtarı, secret key ve servis URL gibi değerler yalnızca env veya secret yönetimi ile taşınmalıdır.
- `site` kategorisi yalnızca kullanıcıya görünen genel ürün davranışları, genel site yüzeyleri ve public deneyim ayarları için kullanılmalıdır.
- `communication` kategorisi iletişim sayfası, iletişim kanalı görünürlüğü, destek giriş yüzeyi ve benzeri public iletişim verileri ile sınırlı olmalıdır; provider credential veya gizli anahtarları kapsamaz.
- `operations` kategorisi düşük seviye altyapı tuning'i için değil, bakım modu, kayıt açma-kapama veya belirli runtime süreçleri durdurma gibi kontrollü ürün operasyon davranışları için kullanılmalıdır.
- `security_auth` kategorisi başarısız giriş limiti, cooldown, resend verification aralığı, MFA zorunluluğu ve benzeri auth güvenlik eşikleri için kullanılmalıdır.
- `access_availability` kategorisi audience targeting, entitlement gating, feature availability ve kill switch karar yüzeyleri için kullanılmalıdır.
- Authorization, audience targeting, entitlement gating, feature availability ve kill switch kararları `access` tarafından yorumlanmalıdır.
- Auth güvenlik eşikleri, yorum gönderme aralığı, attachment sınırları, support intake davranışı ve benzeri erişim dışı runtime davranışları ilgili modülün service katmanında yorumlanmalıdır; bunlar yalnızca erişim veya entitlement kararı ürettiğinde `access` ile entegre çalışmalıdır.
- Site içeriği, iletişim içeriği ve erişim dışı runtime ayarlar `access` üzerinden çözülmemelidir.
- Her runtime ayarı için en az şu metadata alanları tanımlanmalıdır: `key`, `description`, `category`, `owner_module`, `consumer_layer`, `value_type`, `default_value`, `allowed_range_or_enum`, `scope_kind`, `scope_selector`, `audience_kind`, `audience_selector`, `sensitive`, `apply_mode`, `cache_strategy`, `schedule_support`, `audit_required`, `affected_surfaces` ve gerektiğinde `disabled_behavior` ile `error_response_policy`.
- Selector gerektirmeyen `site` veya `all` gibi kayıtlarda `scope_selector` ve `audience_selector` için açık boş değer standardı kullanılmalıdır.
- `scope_selector` için canonical yaklaşım en az `-`, `<module>`, `<module>.<surface>`, `<module>.<surface>.<subsurface>` ve gerektiğinde `resource:<module>:<resource_kind>:<identifier>` biçimlerini desteklemelidir.
- `audience_selector` için canonical yaklaşım en az `-`, `role:<name>`, `group:<name>` ve `user:<id>` biçimlerini desteklemelidir.
- `apply_mode` en az `immediate`, `cache_refresh` ve `scheduled` değerlerini destekleyecek şekilde tasarlanmalıdır.
- `cache_strategy` en az `none`, `ttl` ve `manual_invalidate` gibi açık stratejilerle tanımlanmalıdır.
- `schedule_support` en az `none`, `start_at` ve `time_window` gibi açık planlama modları ile tanımlanmalıdır.
- Runtime ayar envanterinin canonical kayıt dosyası `docs/shared.md` olmalı; yeni ayar, toggle, kill switch veya limit eklendiğinde aynı değişiklik setinde bu dosya güncellenmelidir.
- Ücretli veya süreli avantajı etkileyen ayarlarda metadata'ya ek olarak `entitlement_impact_policy` zorunlu olmalıdır.

## 14) Migration ve Çalışma Ortamı Kuralları
- Migration yönetiminde `golang-migrate` kullanılmalıdır.
- Her migration için `up` ve `down` script zorunludur.
- Backend migration dosyaları yalnızca `apps/api/migrations/` altında tutulmalıdır.
- Migration dosyaları standart isimlendirme ile oluşturulmalıdır.
- Çok modüllü yapıda migration isimlerinde modül veya alan prefixi kullanılmalıdır.
- Örnek yaklaşım: `YYYYMMDDHHMM_auth_create_sessions.up.sql` veya `YYYYMMDDHHMM_content_create_manga.up.sql`.
- Şema değişiklikleri migration olmadan uygulanmamalıdır.
- Seed ve migration süreçleri birbirine karıştırılmamalıdır.
- Backend için gerekli uygulama çalışma dosyaları `apps/api/` altında, repo seviyesi deploy dosyaları ise `deploy/` altında tutulmalıdır.
- Dockerfile uygulama kökünde ilgili app altında yer almalı; compose ve benzeri çok servisli çalışma dosyaları repo seviyesinde merkezi olarak yönetilmelidir.
- Proje Docker içinde build alabilmeli ve çalışabilmelidir.
- Çalışma için gereken servisler tekrarlanabilir şekilde tanımlanmalıdır.
- Local, test ve benzeri ortamlar mümkün olduğunca tutarlı olmalıdır.

## 15) Git, PR ve Kod İnceleme Kuralları

- Git reposu: `https://github.com/Tokuchi61/Manga`
- Varsayılan remote `origin` olmalıdır.
- Push işlemleri yalnızca bu repoya yapılmalıdır.
- Onay olmadan farklı remote eklenmemeli ve farklı repolara push yapılmamalıdır.

- Branch modeli `main + feature/* + hotfix/*` olmali; ajan tabanli akislarda `codex/**` branchleri gecici calisma dali olarak kabul edilebilir.
- `main` daima deploy edilebilir durumda kalmalıdır.
- Doğrudan `main` branch'e push yapılmamalıdır.
- Tüm değişiklikler PR üzerinden ilerlemelidir.

- Branch adları kısa, açıklayıcı ve konu odaklı olmalıdır.
- Feature branch formatı: `feature/<konu>`
- Hotfix branch formatı: `hotfix/<konu>`

- Commit'ler küçük, anlamlı ve geri alınabilir olmalıdır.
- Commit mesajları Conventional Commits standardına uygun olmalıdır.
- Tek commit içinde birden fazla bağımsız konu birleştirilmemelidir.

- Her PR tek bir konuya odaklı olmalıdır.
- Büyük geliştirmeler küçük ve incelenebilir PR'lara bölünmelidir.
- Altyapı, modül geliştirmesi, refactor ve doküman güncellemeleri mümkün olduğunda mantıklı parçalara ayrılmalıdır.

- PR açıklaması en az şu bölümleri içermelidir:
  - ne değişti
  - neden değişti
  - nasıl test edildi

- PR açılmadan önce en az şu kontroller yapılmış olmalıdır:
  - ilgili testler çalıştırılmış olmalı
  - lint/format kontrolleri geçmiş olmalı
  - yeni migration varsa kontrol edilmiş olmalı
  - ilgili doküman güncellemeleri eklenmiş olmalı
  - gereksiz debug/log/yorum satırları temizlenmiş olmalı

- Veritabanı şemasını etkileyen değişikliklerde migration zorunludur.
- Migration içeren PR'larda ilgili model, repository, servis ve test güncellemeleri birlikte değerlendirilmelidir.
- Geri alma etkisi yüksek migration'lar PR açıklamasında ayrıca belirtilmelidir.

- Kod değişikliği mimari, modül sınırı, veri sahipliği, ayar, event akışı veya erişim davranışını etkiliyorsa ilgili dokümanlar aynı PR içinde güncellenmelidir.
- İlgili dokümanlar güncel değilse PR merge edilmemelidir.

- Merge öncesi CI sonucu başarılı olmalıdır.
- Paylaşılan branch'lerde force push kullanılmamalıdır.
- Commit geçmişi inceleme sürecini bozacak şekilde yeniden yazılmamalıdır.
- Varsayılan merge yöntemi ekip standardına göre belirlenmeli; aksi belirtilmedikçe squash merge tercih edilmelidir.

- Ajanlar doğrudan `main` branch üzerinde çalışmamalıdır.
- Her görev için uygun bir `feature/*` veya `hotfix/*` branch açılmalıdır.
- Her aşama sonunda değişiklikler ilgili branch'e push edilmeli ve PR açılmaya hazır halde bırakılmalıdır.
- Remote veya branch belirsizse ajan varsayım yapmamalı, mevcut git yapılandırmasını korumalıdır.

## 16) Versiyonlama Kuralları
- Proje versiyonlamasında SemVer (`MAJOR.MINOR.PATCH`) standardı kullanılmalıdır.
- Versiyon formatı üretim ve kalıcı release'ler için yalnızca `X.Y.Z` biçiminde olmalıdır.
- Geliştirme ve release adayı süreçlerinde gerekirse pre-release etiketleri kullanılabilir:
  - `X.Y.Z-alpha.N`
  - `X.Y.Z-beta.N`
  - `X.Y.Z-rc.N`
- Build metadata gerekiyorsa SemVer ile uyumlu `+build` eki kullanılabilir; ancak asıl ürün versiyonu bunun öncesindeki canonical değerdir.
- Uygulamanın çalışırken gösterdiği canonical versiyon tek kaynak üzerinden yönetilmelidir.
- Runtime tarafında versiyon bilgisi `APP_VERSION` environment değişkeni üzerinden okunmalıdır.
- Versiyon bilgisi kod içine sabit string olarak gömülmemelidir.
- Repo içindeki doküman, release kaydı, tag ve dağıtım çıktıları aynı canonical versiyon ile hizalı olmalıdır.
- Aynı içerik için birden fazla canonical versiyon adı üretilmemelidir.
- Her release tek bir versiyon numarasına sahip olmalıdır.
- Yayınlanmış bir versiyon sonradan sessizce değiştirilmemelidir; yeni değişiklik gerekiyorsa yeni versiyon çıkarılmalıdır.
- Geri alma gerekiyorsa eski versiyonu sessizce oynatmak yerine yeni bir düzeltme versiyonu üretilmelidir.
- Versiyon artırımı gerektiren değişiklikler release öncesinde netleştirilmelidir.
- Versiyon artırımı şu kurallara göre yapılmalıdır:
  - `MAJOR`: Geriye dönük uyumsuz API değişikliği, veri modeli kırılması, davranış değişikliği, kaldırılan alan/endpoint, zorunlu migration uyumsuzluğu, mevcut entegrasyonları bozan mimari değişiklik.
  - `MINOR`: Geriye dönük uyumlu yeni özellik, yeni endpoint, yeni modül, yeni opsiyonel alan, mevcut davranışı kırmadan yapılan anlamlı kapasite artışı.
  - `PATCH`: Geriye dönük uyumlu bugfix, güvenlik düzeltmesi, küçük performans iyileştirmesi, davranışı kırmayan iç düzeltme.
- Yalnızca dokümantasyon, yorum, metin düzeltmesi veya release çıktısını etkilemeyen iç temizlikler tek başına versiyon artırmak zorunda değildir.
- Ancak doküman değişikliği mevcut release'in kullanımını, kurulumunu, entegrasyonunu veya güvenliğini fiilen etkiliyorsa uygun versiyon artırımı değerlendirilmelidir.
- Veritabanı migration içeren her değişiklik için versiyon etkisi ayrıca değerlendirilmelidir.
- Geriye dönük uyumsuz migration değişiklikleri `MAJOR`, uyumlu şema genişletmeleri en az `MINOR` olarak ele alınmalıdır.
- Güvenlik açığı kapatan değişiklikler varsayılan olarak en az `PATCH` artışı ile yayınlanmalıdır.
- Public API sözleşmesini etkileyen her değişiklikte versiyon etkisi açıkça belirtilmelidir.
- PR açıklamalarında gerekiyorsa hedef versiyon veya beklenen bump tipi belirtilmelidir.
- Release hazırlığında en az şu alanlar birlikte güncellenmelidir:
  - `APP_VERSION`
  - `docs/changelog.md`
  - gerekiyorsa `README.md`
  - gerekiyorsa kurulum, migration veya breaking change notları
- `docs/changelog.md` içinde her release için en az şu bilgiler yer almalıdır:
  - versiyon
  - tarih
  - değişiklik özeti
  - etkilenen modüller
  - breaking change bilgisi varsa açık not
  - migration etkisi varsa açık not
- `docs/changelog.md` yalnızca final release başlıklarından ibaret olmamalı; modül, feature, hotfix, fix, refactor, security ve operasyonel düzeltmeler uygun release girdisi altında izlenebilir şekilde gruplanmalıdır.
- Release girişlerinde mümkün olduğunda şu alt başlıklar kullanılmalıdır:
  - `Added`
  - `Changed`
  - `Fixed`
  - `Removed`
  - `Deprecated`
  - `Security`
  - `Docs`
- Hotfix kayıtlarında en az etkilenen alan, sorunun kısa özeti ve düzeltilen kapsam belirtilmelidir.
- API, migration, config, access veya operasyonel davranışı etkileyen değişikliklerde gerekli kullanıcı veya geliştirici aksiyonları changelog içinde açıkça yazılmalıdır.
- Release çıkmadan önce changelog taslak girdileri hazırlanabilir; ancak yayın anında hepsi canonical versiyon başlığı altında birleştirilmelidir.
- Breaking change içeren release'lerde upgrade notu zorunlu olmalıdır.
- Release candidate kullanılıyorsa production'a çıkmadan önce final versiyon ayrıca üretilmelidir.
- Git tag standardı canonical versiyon ile uyumlu olmalı ve `vX.Y.Z` biçiminde açılmalıdır.
- Pre-release tag'leri gerekiyorsa `vX.Y.Z-rc.N` benzeri biçimde oluşturulmalıdır.
- Tag, changelog ve dağıtılan artifact versiyonu birbiriyle çelişmemelidir.
- Release alınmadan önce build, test ve kritik doğrulama adımları başarılı olmalıdır.
- Versiyon artışı yapılan değişikliklerde rollback, migration ve uyumluluk etkisi en az bir kez gözden geçirilmelidir.

## 17) Test ve Doğrulama Kuralları
- Yeni eklenen veya güncellenen yapı test edilebilir olmalıdır.
- Testler ana veritabanını asla kullanmamalı, yalnızca test veritabanı ile çalışmalıdır.
- Unit test önceliği iş kuralları ve veri doğrulama katmanlarıdır.
- Integration test önceliği veri erişimi, modül kontratları ve kritik HTTP akışlarıdır.
- Yeni endpoint eklendiğinde en az bir başarılı ve bir başarısız senaryo test edilmelidir.
- Veri sızıntısı ve yetkisiz erişim senaryoları test edilmelidir.
- Çapraz modül entegrasyonları varsa contract veya integration test ile doğrulanmalıdır.
- `go test ./...` temel doğrulama kontrolü olarak kabul edilmelidir.
- Başlangıç coverage hedefi minimum `%60` olarak önerilir.

## 18) Kapsam ve Dokümantasyon Yönetimi
- Kapsam dışı feature'lar varsayılıp eklenmemelidir.
- Ürün veya mimari kapsamını değiştiren kararlar önce dokümana yansıtılmalıdır.
- Aynı anda çok fazla sorumluluk açmak yerine en küçük sürdürülebilir kapsam seçilmelidir.
- `RULES.md`, `ROADMAP.md` ve gerekli diğer dokümanlar birlikte güncellenmelidir.
- `docs/shared.md`, `docs/shared.md` ve benzeri canonical kayıt dosyaları etkilendiklerinde aynı değişiklik setinde güncellenmelidir.
- Aynı bilgi birden fazla dosyada çelişkili şekilde bırakılmamalıdır.
- `README.md` proje girişi ve ana doküman bağlantıları için güncel tutulmalıdır.
- Sürüm bazlı değişiklikler `docs/changelog.md` içinde kayıt altına alınmalıdır.
- Bilinen sorunlar ve teknik borçlar `docs/issues.md` içinde tutulmalıdır.
- Repo kök dizin yapısını, uygulama yerleşimini veya deployment klasör yapısını etkileyen değişikliklerde `RULES.md`, `README.md` ve `SETUP.md` aynı değişiklik setinde birlikte güncellenmelidir.
- Her yeni leaf modül için en az bir modül dokümanı açılmalıdır.
- Yeni runtime ayar, feature toggle, kill switch veya oran limiti eklendiğinde ilgili key, scope, scope selector, audience, audience selector, varsayılan değer, disabled behavior varsa bunun tipi, error response policy varsa bunun tipi ve etkilediği modüller dokümana yansıtılmalı; `docs/shared.md` aynı değişiklikte güncellenmelidir.
- Modül dokümanlarında en az şu başlıklar yer almalıdır:
  - amaç
  - sorumluluk alanı
  - veri sahipliği
  - access kontratı
  - API veya event sınırı
  - bağımlılıklar
  - state yapısı
  - test notları

## 19) Son Mimari İlke
- Kimlik doğrulama, veri sahipliği, erişim kararı ve yönetim akışları birbirine karıştırılmamalıdır.
- Veri sahipliği ilgili alanda, authorization kararı merkezi yapıda kalmalıdır.
- İçerik, topluluk, operasyon, ticaret, etkileşim ve altyapı sorumlulukları ayrışmış kalmalıdır.
- Yeni özellik eklemek bu ayrımı bozma gerekçesi olamaz.

## 20) Modül Referans Yapısı
- Modül kuralları doküman içinde tek bir ana bölüm altında toplanmalıdır.
- Bu yapı için ana modül bölümü `21) Modül Genel Kuralları` başlığı olmalıdır.
- Her leaf modül bu ana bölüm altında `21.1`, `21.2`, `21.3` şeklinde ayrı alt başlık olarak eklenmelidir.
- Yeni modül eklendikçe numaralandırma aynı ana bölüm altında devam etmelidir: `21.5`, `21.6`, `21.7` ... `21.20`.
- Modül sayısı artsa bile her yeni modül için yeni bir üst seviye bölüm açılmamalı; modül kuralları `21.x` formatında sürdürülmelidir.
- Yeni modüller mevcut modül maddelerinin içine gömülmemeli; her yeni leaf modül ayrı alt başlık olarak eklenmelidir.
- Roadmap tarafında da aynı yaklaşım korunmalı; her yeni leaf modül mevcut aşamaların içine sıkıştırılmadan ayrı bir aşama olarak eklenmelidir.

## 21) Modül Genel Kuralları
- Bu bölüm, proje için kabul edilmiş leaf modül referanslarını merkezi olarak kaydeder.
- Bu bölümdeki kayıtlar yüksek seviye modül sınırı ve sahiplik özeti taşır; detaylı modül tasarımı ilgili modül dokümanında yer almalıdır.
- Her modül alt başlığı canonical modül adını, temel sorumluluk alanını, veri sahipliğini ve ana referans dokümanını açıkça göstermelidir.
- Her modül, ihtiyaç duyduğu durumda admin tarafından yönetilen runtime ayarlar, modül açma-kapama ve özellik açma-kapama yüzeyleri ile uyumlu çalışacak şekilde tasarlanmalıdır; veri sahipliği ilgili modülde, authorization veya availability yorumlama `access`te, erişim dışı runtime davranış yorumlama ilgili modül service katmanında, operasyon yönetimi `admin`de kalmalıdır.

### 21.1) `auth` Modülü Kuralları
- Canonical modül adı `auth` olmalıdır.
- `auth` modülü kimlik doğrulama, oturum güvenliği ve hesap giriş güvenliği alanlarının sahibidir.
- `auth` modülü kayıt, giriş, çıkış, token, session, email verification, password reset, password change, login güvenliği ve auth audit akışlarını taşımalıdır.
- `auth` modülü kullanıcı profil verisi, kullanıcı tercihleri, üyelik avantajları, rol/permission yönetimi veya authorization kararı üretmemelidir.
- `auth` veri sahipliği; credential benzeri kimlik bilgileri, auth session kayıtları, token yaşam döngüsü ve auth güvenlik olayları ile sınırlı kalmalıdır.
- `auth` modülünün public surface'i, kimlik doğrulama akışları ve güvenli oturum yönetimi için gerekli contract yüzeyi ile sınırlı olmalıdır.
- `auth` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - refresh token üretimi ve yenileme akışı
  - session listing, session revoke ve logout all akışları
  - forgot password, reset password ve change password ayrımı
  - email verification token, resend verification ve verification güvenliği
  - login rate limit, failed login limit ve cooldown veya temporary lock davranışı
  - device, IP, son giriş ve şüpheli giriş takibi
  - email doğrulanmamış, suspend edilmiş veya banlı kullanıcı için auth kontrol davranışı
  - başarısız giriş limiti, cooldown süresi, resend verification aralığı ve gerektiğinde MFA zorunluluğu gibi güvenlik eşiklerinin admin tarafından yönetilen runtime ayarlar ile kontrol edilebilmesi
  - çok faktörlü doğrulama, trusted device ve risk score temelli giriş değerlendirmesi
  - şüpheli giriş durumlarında ek challenge veya doğrulama adımı

### 21.2) `user` Modülü Kuralları
- Canonical modül adı `user` olmalıdır.
- `user` modülü kullanıcı hesabı, profil, hesap durumu, tercih, görünüm ve üyelik verisinin sahibidir.
- `user` modülü public/private profil ayrımı, kullanıcı arama veya listeleme, hesap durumu alanları ve üyelik verisi gibi kullanıcı merkezli alanları taşımalıdır.
- `user` modülü kimlik doğrulama akışlarını, authorization kararlarını, role/permission yönetimini veya admin operasyon kararlarını sahiplenmemelidir.
- `user` veri sahipliği; kullanıcı kimliği ile ilişkili profil alanları, hesap durum alanları, tercih alanları ve üyelik durum verileri ile sınırlı kalmalıdır.
- `user` modülü veri taşır ve yayınlar; erişim kararı veya feature access kararı üretmez.
- `user` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - kullanıcı arama, kullanıcı listeleme ve reserved username kontrolü
  - email change, preferences update ve privacy alanları
  - aktif veya pasif, suspend, ban, soft delete ve restore hesap durumları
  - display name, avatar, banner, bio ve görünüm alanları
  - visibility preset yapıları ve profil görünürlük şablonları
  - üyelik ve VIP veri alanları
  - EXP, level, level progress ve kullanıcıya ait oyunlaştırma özeti
  - VIP rozet ve envanterden seçilen profil efekti veya nameplate referansları gibi profil görünüm alanları
  - profil görünürlüğü, VIP görünürlüğü ve benzeri kullanıcı alt yüzeylerinin admin tarafından yönetilen feature toggle veya runtime ayarlar ile kontrollü biçimde açılıp kapatılabilmesi
  - sistem kaynaklı VIP global pasiflikte mevcut kullanıcının kalan VIP süresinin dondurulması ve tekrar açıldığında kaldığı yerden devam etmesi
  - kullanıcıya ait genel preference sinyalleri; ancak detaylı bildirim tercihleri `notification`, sosyal ilişki blok veya mute listeleri ise `social` modülünde sahiplenilmelidir
  - reading history veya library görünürlüğünü etkileyen kullanıcı preference sinyalleri; ancak continue reading, bookmark veya okuma kayıtlarının kendisi `history` modülünde sahiplenilmelidir
  - reading activity veya public library görünürlüğünde `user` modülündeki global preference sinyalinin üst sınır, `history` içindeki entry-level share metadata'sının ise bu üst sınır içinde çalışan alt karar olması
  - hesap dışa aktarma veya hesap verisi export yüzeyi
  - profil değişim geçmişi veya profile history görünümü

### 21.3) `access` Modülü Kuralları
- Canonical modül adı `access` olmalıdır.
- `access` modülü sistemdeki tüm authorization, policy ve erişim kararlarının merkezi sahibi olmalıdır.
- `access` modülü role, permission, policy, ownership, guest/authenticated/vip kararları, endpoint guard ve modül bazlı authorization contract alanlarını taşımalıdır.
- `access` modülü kullanıcı profili, credential verisi, içerik verisi veya yönetim use-case verisi taşımamalıdır.
- `access` veri sahipliği; authorization sözlüğü, rol-permission ilişkileri, policy yorumları ve erişim kararına temel olan kurallarla sınırlı kalmalıdır.
- `access` modülü kimlik doğrulama yapmaz; `auth` tarafından doğrulanan kimliği ve `user` tarafından taşınan kullanıcı verisini kullanarak karar üretir.
- `access` modülü yalnızca authorization, audience targeting, entitlement gating, feature availability ve kill switch kararlarını yorumlamalıdır; `site` veya `communication` kategorisindeki ürün ayarları ile erişim dışı iş kuralı eşikleri `access` içinde çözülmemelidir.
- `access` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - public, guest, authenticated, vip, early access ve gerektiğinde restricted kararları
  - role CRUD, permission CRUD, user-role ve role-permission ilişkileri
  - default role atama, çoklu rol desteği ve rol önceliği kuralları
  - own veya any ayrımı, ownership ve resource access kuralları
  - endpoint guard, use-case guard ve admin panel görünürlük kararları
  - super admin bypass, moderator veya admin override kuralları
  - chapter için minimum `authenticated` okuma kapısı, misafir kullanıcının site ve manga detayına erişebilmesi, ancak chapter okuma ve yorum yazma gibi akışlarda kısıtlı kalması
  - VIP özel bölüm erişimi ve belirli bölümler için VIP erken erişim kararlarının merkezi yönetimi
  - reklam görünürlüğü, VIP reklamsız deneyim, kozmetik görünürlük ve ileride eklenecek özellikler için feature access kararları
  - yeni modüller için canonical permission örnekleri; örnek olarak `history.continue_reading.read.own`, `history.timeline.read.own`, `history.library.read.own`, `history.bookmark.write.own`, `history.library.read.public`, `manga.discovery.view`, `ads.view`, `shop.item.purchase`, `payment.mana.purchase` ve `payment.transaction.read.own`
  - her modül için zorunlu authorization kontratı ve canonical permission isimlendirmesi
  - feature flag tabanlı policy, rollout ve geçici davranış kontrolü
  - admin tarafından tetiklenebilen site geneli, modül bazlı, özellik bazlı, audience bazlı ve gerektiğinde context bazlı acil kapatma, emergency deny veya kill switch yüzeyi
  - availability kurallarında `emergency_deny` > `deny/off` > `allow/on` > varsayılan değer önceliği
  - audience scope yönetimi; başlangıçta `all`, `guest`, `authenticated`, `authenticated_non_vip` ve `vip`, ileride gerekirse daha spesifik hedefler
  - aynı key, aynı `audience_kind + audience_selector` ve aynı `scope_kind + scope_selector` için çakışan aktif kural bırakmama ve kayıt aşamasında reddetme davranışı
  - temporary grants, süreli yetki verme ve kontrollü yetki geri alma
  - policy versioning ve denial explanation surface
  - VIP kullanılabilirliği ile VIP entitlement süresinin birbirinden ayrı ele alınması; sistem kaynaklı global pasiflikte kalan sürenin dondurulması
  - kişi bazlı ve alan bazlı moderatör yetkilendirmesi; örnek olarak yorum moderatörü, bölüm moderatörü veya manga moderatörü gibi ayrık yetki yüzeyleri
  - `moderation` modülü ile bağlanacak moderator scope ve delegation altyapısı

### 21.4) `admin` Modülü Kuralları
- Canonical modül adı `admin` olmalıdır.
- `admin` modülü yönetim, tam yetkili inceleme, merkezi ayar ve operasyon use-case'lerinin sahibidir.
- `admin` modülü dashboard, yönetim giriş noktaları, kullanıcı yönetim akışları, support review, tam yetkili moderasyon gözetimi, operasyonel kontrol ve admin audit akışlarını taşımalıdır.
- `admin` modülü kendi içinde authorization kararı, role/permission kararı veya kullanıcı profil veri sahipliği üretmemelidir.
- `admin` veri sahipliği; yönetimsel işlem kayıtları, admin notları, admin use-case akışları ve operasyonel görünüm alanları ile sınırlı kalmalıdır.
- `admin` modülünün tüm kritik akışları `access` guard veya policy kararları ile korunmalıdır.
- Gerekli role veya permission'a sahip admin kullanıcıları sistemdeki yönetim ve inceleme yüzeylerine tam erişim taşıyabilir; bu durum günlük scoped moderator use-case sahipliğini `moderation` modülünden almaz.
- `admin` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - admin dashboard ve dashboard veri ihtiyaçları
  - kullanıcı yönetimi, kullanıcı durum müdahaleleri, warning, restriction, suspend ve ban akışları
  - manga, chapter ve comment için yüksek riskli moderasyon handoff, escalation ve tam yetkili yönetimsel inceleme yüzeyleri
  - support review queue, destek karar yürütme yüzeyleri ve iletişim odaklı yönetim akışları
  - yüksek riskli admin aksiyonlarında açık permission, zorunlu reason ve gerektiğinde ek doğrulama
  - cache temizleme, log görüntüleme ve sistem sağlık durumu gibi operasyonel araçlar
  - ileri aşamalarda eklenecek yeni yönetim yüzeyleri için hazırlık
  - bulk actions, case timeline ve toplu moderasyon araçları
  - canned moderation actions ve approval workflow yapıları
  - sistem genelindeki runtime ayarların, modül açma-kapama ve alt özellik açma-kapama yüzeylerinin merkezi yönetimi
  - `site`, `communication`, `operations`, `security_auth`, `access_availability`, `content`, `reading`, `engagement`, `support`, `membership`, `social`, `gamification` ve `economy` kategorileri için genişleyebilir settings merkezi
  - başarısız giriş denemesi limiti, yorum gönderme aralığı, yorumların manga detayında açık veya kapalı olması gibi iş kuralı eşiklerinin yönetim yüzeyi
  - env veya secret yönetimi gerektiren teknik config ile admin runtime ayarlarının kesin olarak ayrılması
  - ayar metadata kataloğu; key, tip, scope kind, scope selector, audience kind, audience selector, apply mode, cache strategy, error response policy ve entitlement impact policy gibi alanların merkezi yönetimi
  - `moderation` modülü eklense bile merkezi ayar ve kill switch yüzeylerinin yalnızca admin tarafında kalması; moderatör paneline delegasyon yapılmaması
  - operatör iş yükü görünümü ve moderasyon veya destek yük dağılımı takibi
  - access içindeki emergency deny veya kill switch yüzeyini yönetebilecek operasyon kontrol noktaları

### 21.5) `manga` Modülü Kuralları
- Canonical modül adı `manga` olmalıdır.
- `manga` modülü ana içerik varlığının, metadata yapısının, taxonomy ilişkilerinin ve içerik yaşam döngüsü verisinin sahibidir.
- `manga` modülü CRUD, listing, detail, search, filtering, sorting, publish akışları, metadata/taxonomy alanları ve içerik sayaçlarını taşımalıdır.
- `manga` modülü chapter için varsayılan okuma erişim verisini taşıyabilir; ancak erişim kararını kendi içinde üretmemelidir.
- `manga` veri sahipliği; başlık, özet, görsel, taxonomy, yayın durumu, görünürlükle ilişkili state alanları ve içerik sayaçları ile sınırlı kalmalıdır.
- `manga` içindeki `chapter_count`, `comment_count` ve benzeri sayaç alanları denormalize okuma alanları olarak ele alınmalı; canonical kaynak verisi ilgili kaynak modülde kalmalıdır.
- `manga` sayaç güncellemeleri `chapter` veya `comment` modülünün manga tablosuna doğrudan yazması ile yapılmamalı; event, projection veya açık counter contract yüzeyi ile senkronize edilmelidir.
- `manga` sayaçları için kabul edilen gecikme modeli ve gerektiğinde reconcile veya yeniden hesaplama yolu dokümante edilmelidir.
- `manga` modülünün public surface'i içerik listeleme, içerik detay, yönetimsel içerik işlemleri ve taxonomy ilişkileri için gerekli yüzey ile sınırlı olmalıdır.
- `manga` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - slug ve benzersizlik kuralları
  - alternative titles, short summary, cover image, banner image ve SEO alanları
  - genre, tag, theme ve content warning taxonomy yapıları
  - draft, scheduled, published ve archived veya unpublished benzeri publish yaşam döngüsü
  - featured veya recommended işaretleri, editoryal koleksiyonlar ve içerik sayaçları
  - soft delete ve restore desteği
  - chapter için varsayılan read access ve varsayılan VIP erken erişim ayarları
  - release schedule ve translation group gibi yayın planı alanları
  - toplu planlı yayın ve editoryal yayın paketleri
  - recommendation, içerik koleksiyonu ve editoryal keşif yüzeyleri; ancak kullanıcıya ait continue reading, reading history veya bookmark-library kayıtları `history` modülünde kalmalıdır
  - manga listeleme, manga detay ve editoryal görünürlük gibi yüzeylerin admin tarafından yönetilen runtime ayarlar ile kontrollü biçimde daraltılabilmesi

### 21.6) `chapter` Modülü Kuralları
- Canonical modül adı `chapter` olmalıdır.
- `chapter` modülü manga içeriğinin okunabilir bölüm yapısının, bölüm sıralamasının, bölüm sayfalarının ve bölüm yaşam döngüsü verisinin sahibidir.
- `chapter` modülü CRUD, manga bazlı chapter listesi, detail, read akışı, navigation, page/media ilişkileri, numbering ve publish akışlarını taşımalıdır.
- `chapter` modülü chapter erişimini etkileyen veri alanlarını taşıyabilir; ancak guest/authenticated/vip/early access kararlarını kendi içinde üretmemelidir.
- `chapter` veri sahipliği; chapter metadata alanları, page yapısı, publish state, access state verisi ve navigation alanları ile sınırlı kalmalıdır.
- `chapter` kullanıcıya ait son okuma pozisyonu, reading session progress, continue reading kaydı veya bookmark-library state'i taşımamalı; bunlar `history` modülünde tutulmalıdır.
- `chapter` modülünün public surface'i okuma akışı, chapter detail ve yönetimsel chapter işlemleri için gerekli contract yüzeyi ile sınırlı olmalıdır.
- `chapter` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - latest chapter list ile previous, next, first ve last navigation akışları
  - chapter page yapısı, width veya height bilgisi ve gerekirse long strip desteği
  - read_access_level, inherit_access_from_manga, early_access_enabled, early_access_level ve fallback alanları
  - VIP özel bölüm ve VIP için erken erişim bölüm yapılarının birbirinden ayrılması
  - misafir kullanıcının manga detayına erişebilmesi ama chapter okuma için minimum `authenticated` kapısının zorunlu olması
  - early access zaman penceresi, pencere sonrası fallback erişimi ve access policy hizası
  - soft delete ve restore desteği
  - media validation ve CDN health kontrol yüzeyleri
  - bozuk veya eksik medya tespiti için checksum veya benzeri bütünlük doğrulaması
  - `history` modülü için continue reading, reading history, resume anchor ve progress entegrasyon yüzeyi
  - chapter okuma, preview veya belirli okuma alt yüzeylerinin admin tarafından yönetilen runtime ayarlar ile kontrollü açılıp kapatılabilmesi
  - preview kapalıyken detail veya tam read yüzeyinin otomatik olarak kapanmaması; her yüzeyin ayrı kontrol edilebilmesi

### 21.7) `comment` Modülü Kuralları
- Canonical modül adı `comment` olmalıdır.
- `comment` modülü yorum verisinin, thread/reply yapısının, yorum görünürlük verisinin ve yorum yaşam döngüsünün sahibidir.
- `comment` modülü en az `manga` ve `chapter` hedef tiplerini desteklemeli; yorum create, edit, delete, reply, listing, moderation state ve thread akışlarını taşımalıdır.
- `comment` modülü yorum görünürlüğünü etkileyen veri alanlarını taşıyabilir; ancak create/edit/delete/pin/lock gibi işlemlerin yetki kararını kendi içinde üretmemelidir.
- `comment` modülü sosyal duvar post'u veya sosyal duvar reply akışını sahiplenmemeli; bunlar `social` modülünde kalmalı ve `comment` sistemine örtük olarak dönüştürülmemelidir.
- `comment` veri sahipliği; yorum içeriği, hedef ilişkisi, reply yapısı, moderation/spoiler/lock verileri ve sıralama/listeleme alanları ile sınırlı kalmalıdır.
- `comment` modülündeki `target_type` değerleri canonical olarak `docs/shared.md` dosyasındaki kayıtlarla hizalı olmalıdır.
- Yeni yorum hedef tipi eklendiğinde `comment` modülü, hedef modül dokümanı ve canonical target type kaydı aynı değişiklik setinde güncellenmelidir.
- `comment` modülünün public surface'i yorum listeleme, yorum detay, thread akışı ve hedef ilişkisi için gerekli yüzey ile sınırlı olmalıdır.
- `comment` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - root comment ve reply thread yapısı
  - newest, oldest ve popular sıralama seçenekleri ile pagination
  - sanitize edilmiş içerik veya güvenli çıktı alanı
  - anonymous görüntüleme ve authenticated yorum yazma ayrımı
  - reply derinliği, edit window, lock etkisi ve restore sınırları
  - report edilebilir hedef olma ve ileri aşamalarda yeni target_type genişletmeleri
  - gelişmiş moderasyon ve anti-spam akışları
  - yorum attachment desteği
  - anti-spam score veya yorum risk puanı yaklaşımı
  - site geneli, manga detayı veya chapter altı yorum alanlarının ayrı ayrı açılıp kapatılabilmesi ve yorum gönderme aralığı gibi etkileşim eşiklerinin yönetilebilmesi
  - yorum modülü için read ve write yüzeylerinin ayrı ayrı kontrol edilebilmesi; varsayılan kapatma senaryosunda mevcut yorumların görünür kalıp yeni yorumların engellenebilmesi
  - kullanıcı bazlı mute ve moderation escalation davranışları
  - sessiz moderasyon; içeriği tamamen silmeden görünürlük kapsamını kısıtlama yaklaşımı

### 21.8) `support` Modülü Kuralları
- Canonical modül adı `support` olmalıdır.
- `support` modülü kullanıcı iletişim taleplerinin, destek biletlerinin, manga/chapter/comment için hedefe bağlı içerik bildirimlerinin ve destek yaşam döngüsünün sahibidir.
- Ayrı bir `report` leaf modül açılmamalı; içerik bildirimi akışları `support` modülünün destek ve ticket alanı içindeki bir feature yüzeyi olarak ele alınmalıdır.
- `support` modülü communication/create, ticket/create, hedefe bağlı report/create, own support list, support detail, support reply, category, priority, duplicate/spam kontrolü, review queue verisi, status update ve resolution note akışlarını taşımalıdır.
- `support` veri sahipliği; destek kaydı, `support_kind`, category, isteğe bağlı hedef ilişkisi, destek durumu, mesaj veya reply zinciri, çözüm notları ve inceleme yaşam döngüsü verileri ile sınırlı kalmalıdır.
- `support` modülü review ve karar verisini taşır; ancak yetki kararını kendi içinde üretmez, yönetimsel karar yürütümü `admin`, authorization ise `access` ile korunur.
- `support` modülündeki report kaydı varsayılan olarak moderation case ile aynı kayıt sayılmamalıdır; moderation ihtiyacı oluştuğunda linked case açılabilir, ancak support intake kaydı ve moderation case yaşam döngüsü ayrı owner sınırlarında kalmalıdır.
- `support` modülündeki `target_type` değerleri canonical olarak `docs/shared.md` dosyasındaki kayıtlarla hizalı olmalıdır.
- Genel iletişim veya hedefsiz destek biletlerinde `target_type` zorunlu olmamalı; bu alan yalnızca manga/chapter/comment gibi hedefe bağlı kayıtlar için kullanılmalıdır.
- Yeni support hedef tipi eklendiğinde `support` modülü, hedef modül dokümanı ve canonical target type kaydı aynı değişiklik setinde güncellenmelidir.
- `support` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - communication, ticket ve hedefe bağlı report akışlarının aynı modülde ama ayrı kayıt mantığı ile taşınması
  - manga, chapter ve comment hedef tipleri ile isteğe bağlı target relation yapısı
  - support kind, category, priority ve reason code kataloğu
  - duplicate veya spam kontrolü ve aynı kullanıcı ile aynı hedef için tekrar bildirim davranışı
  - review queue, status update, support reply, resolution note, assignee, reviewed_by ve resolved_at alanları
  - iletişim kaydı ile hedefe bağlı içerik bildiriminin aynı yaşam döngüsünde ama ayrı kurallarla işlenebilmesi
  - support attachment desteği
  - canned support replies ve cevap şablonları
  - SLA, önceliklendirme ve escalation yüzeyi
  - support modülü altındaki communication, ticket ve report create yüzeylerinin, attachment kabulünün ve intake davranışlarının admin tarafından yönetilen runtime ayarlar ile kontrol edilebilmesi
  - support için yeni kayıt alımının durdurulabilmesi, ancak mevcut kayıtların kullanıcı veya admin tarafından okunabilir kalması gibi intake pause davranışları
  - raporlayan kullanıcı güven skoru ve güven puanının inceleme önceliğine etkisi
  - bozuk medya veya eksik sayfa bildirimlerinde checksum veya bütünlük doğrulaması ile desteklenebilen inceleme verisi

### 21.9) `moderation` Modülü Kuralları
- Canonical modül adı `moderation` olmalıdır.
- `moderation` modülü role bazlı veya kullanıcı bazlı scoped moderatör panelinin, vaka inceleme akışlarının ve moderatör iş yükü süreçlerinin sahibi olmalıdır.
- `moderation` modülü queue, assignment, case detail, moderator note, sınırlı aksiyon yürütme ve escalation akışlarını taşımalıdır.
- `moderation` modülü authorization, role veya permission sahipliği üretmemeli; moderator scope ve yetki kararları `access` ile korunmalıdır.
- `moderation` veri sahipliği; moderation case, assignment, moderator note, action summary ve escalation lifecycle alanları ile sınırlı kalmalıdır.
- `moderation` günlük scoped inceleme sahibidir; ancak `admin` tarafından aynı case üzerinde verilen override, reopen, freeze, reassignment veya final kararlar daha yüksek precedence taşır.
- `moderation` case hedefleri canonical olarak `docs/shared.md` dosyasındaki kayıtlarla hizalı olmalı; alt yüzey bilgisi `target_type` içine değil context verisine taşınmalıdır.
- `moderation` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - yorum, chapter ve manga yüzeyleri için scoped queue yapıları
  - comment moderator, chapter moderator veya manga moderator gibi role veya kullanıcı bazlı scope modelleri
  - moderator assignment, handoff ve escalation akışları
  - vaka timeline, moderator note ve karar özeti alanları
  - sınırlı moderator aksiyonları; örnek olarak hide, unhide, lock, unlock, escalate veya review complete yüzeyleri
  - admin tarafından yönetilen runtime ayarlar ile queue veya aksiyon yüzeylerinin ayrı ayrı açılıp kapatılabilmesi
  - moderasyon panelinin açık kalması, ancak merkezi settings ve kill switch yüzeylerinin admin modülünde kalması
  - admin tarafına workload, escalation ve audit sinyali üretebilen entegrasyon yüzeyleri

### 21.10) `notification` Modülü Kuralları
- Canonical modül adı `notification` olmalıdır.
- `notification` modülü sistem genelindeki bildirim üretimi, teslimi, kategori yönetimi ve detaylı bildirim tercihleri verisinin sahibi olmalıdır.
- `notification` modülü in-app inbox, read veya unread akışları, category, channel, template, delivery attempt ve suppression yüzeylerini taşımalıdır.
- `notification` modülü business event sahipliği veya authorization kararı üretmemelidir; bildirim olaylarını producer modüllerden almalı ve own-surface erişimini `access` ile korumalıdır.
- `notification` veri sahipliği; notification kaydı, delivery durumu, template veya category tanımı ve kullanıcı bildirim tercihleri ile sınırlı kalmalıdır.
- `notification` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - in-app inbox ve read veya unread davranışları
  - category, template ve channel yönetimi
  - module veya category bazlı delivery pause ve flood control yüzeyleri
  - quiet-hour, mute, digest veya batch delivery davranışları
  - detail bildirim tercihleri sahipliği; `user` modülünde yalnızca özet preference sinyali kalması
  - admin tarafından category, channel veya feature bazlı açma-kapama ve eşik yönetimi
  - `social`, `mission`, `royalpass`, `support`, `moderation` ve diğer producer modüller ile event kontratları

### 21.11) `social` Modülü Kuralları
- Canonical modül adı `social` olmalıdır.
- `social` modülü kullanıcılar arası arkadaşlık, takip, sosyal duvar, duvar altı etkileşim ve mesajlaşma iş alanlarının sahibi olmalıdır.
- `social` modülü friendship request, friendship state, follow relation, wall post veya wall reply ve direct message thread akışlarını taşımalıdır.
- `social` modülü manga veya chapter içerik yorumlarını sahiplenmemeli; bu alanlar `comment` modülünde kalmalıdır.
- `social` modülündeki wall reply yapısı social-native kabul edilmeli; `comment` thread sistemi ile örtük olarak birleştirilmemelidir.
- `social` veri sahipliği; sosyal ilişki kayıtları, social block veya mute ilişkileri, sosyal içerik kayıtları ve sosyal privacy sinyalleri ile sınırlı kalmalıdır.
- `social` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - friend request, accept, reject, remove ve friend list akışları
  - follow veya unfollow davranışları ve takip listeleri
  - profil sosyal duvarı, duvar post veya duvar reply yüzeyleri
  - direct message thread, mesaj gönderme ve mesaj görünürlük kuralları
  - social block veya mute sahipliği; bu alanların `user` modülünde veri sahipliği olarak tutulmaması
  - social block veya privacy deny sinyalinin final authorization kararında `access` tarafından deny precedence ile yorumlanması; mute sinyalinin ise ayrıca dokümante edilmedikçe teslim veya görünürlük sinyali olarak kalması
  - friendship, follow, messaging ve wall yüzeyleri için ayrı runtime control anahtarları
  - admin tarafından messaging, wall, follow veya friendship yüzeylerinin ayrı ayrı açılıp kapatılabilmesi
  - bildirim, anti-spam ve social privacy sinyallerinin birbirine karışmadan çalışması

### 21.12) `inventory` Modülü Kuralları
- Canonical modül adı `inventory` olmalıdır.
- `inventory` modülü item tanımı, kullanıcı envanteri, reward sahipliği ve final grant yürütümünün sahibi olmalıdır.
- `inventory` modülü ownable item definition, user inventory entry, grant, revoke, claim, consume ve equip akışlarını taşımalıdır.
- `inventory` modülü sellable shop product veya offer catalog sahipliği üretmemeli; `shop` ile ilişkisi product-to-item mapping ve final grant kontratı üzerinden kurulmalıdır.
- `inventory` modülü ödeme, görev ilerlemesi veya pass season sahipliği üretmemeli; yalnızca item sahipliği ve item durumunu taşımalıdır.
- `inventory` veri sahipliği; ownable item tanımı, quantity veya stack durumu, expiry, equip state ve source reference alanları ile sınırlı kalmalı; sellable catalog verisini içermemelidir.
- `inventory` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - stackable ve non-stackable item ayrımı
  - grant, reward teslim yürütümü, revoke, consume ve equip yüzeyleri
  - source reference ve idempotent grant davranışları
  - `user` modülünde seçilen kozmetik görünüm referanslarının `inventory` sahipliği ile hizalanması
  - admin tarafından inventory görünürlüğü, claim, consume veya equip yüzeylerinin açılıp kapatılması
  - `mission`, `royalpass` ve ileride gelecek diğer reward producer modülleri ile kontrollü grant kontratları

### 21.13) `mission` Modülü Kuralları
- Canonical modül adı `mission` olmalıdır.
- `mission` modülü günlük, haftalık, aylık, event ve seviye bazlı görev tanımlarının, görev ilerlemesinin ve reward için claim uygunluğu veya claim request yaşam döngüsünün sahibi olmalıdır.
- `mission` modülü mission definition, objective yapısı, progress kaydı, completion, claim eligibility ve reset pencerelerini taşımalıdır.
- `mission` modülü global EXP veya level sahipliğini tek başına üretmemeli; `user` modülündeki progression sinyallerini tüketerek görev değerlendirmesi yapabilmelidir.
- `mission` veri sahipliği; görev tanımı, görev kategorisi, progress state, claim eligibility state ve reward reference alanları ile sınırlı kalmalıdır.
- `mission` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - daily, weekly, monthly, event ve level-based mission tipleri
  - recurring reset, dönemsel yenileme ve gerektiğinde streak davranışları
  - producer event kontratları; örnek olarak okuma, yorum, sosyal etkileşim veya diğer görev sinyalleri
  - claim request yüzeyinin `inventory` içindeki final grant sahipliği ile hizalanması
  - claim, reset veya mission category yüzeylerinin admin tarafından ayrı ayrı açılıp kapatılabilmesi
  - `notification` ile görev tamamlama ve claim bildirim entegrasyonu
  - `royalpass` için görevden progress besleme yüzeyi

### 21.14) `royalpass` Modülü Kuralları
- Canonical modül adı `royalpass` olmalıdır.
- `royalpass` modülü aylık season yapısının, free veya premium track ilerlemesinin ve season ödülü için claim uygunluğu veya claim request yaşam döngüsünün sahibi olmalıdır.
- `royalpass` modülü season, tier, track, user season progress ve reward claim eligibility akışlarını taşımalıdır.
- `royalpass` modülü görev tanımı, item sahipliği veya ödeme sahipliği üretmemelidir; season içi progress ve claim eligibility sahipliği ile sınırlı kalmalıdır.
- `royalpass` veri sahipliği; season tanımı, track veya tier yapısı, user season claim eligibility state ve premium activation referansları ile sınırlı kalmalıdır.
- `royalpass` premium aktivasyonu ürünleşmiş satın alma akışında canonical olarak `shop` üzerinden başlamalı, gerçek para veya mana checkout veya bakiye doğruluğu gerekiyorsa `payment` tarafından tamamlanmalı ve final premium activation referansı `royalpass` tarafından tüketilmelidir.
- `royalpass` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - aylık season yaşam döngüsü ve season archive yapıları
  - free track ve premium track ayrımı
  - mission tabanlı progress veya puan besleme davranışları
  - claim request yüzeyinin `inventory` içindeki final grant sahipliği ile hizalanması
  - season görünürlüğü, claim yüzeyi veya premium track yüzeyinin admin tarafından ayrı ayrı açılıp kapatılabilmesi
  - ürünleşmiş premium pass satın alma zincirinde `shop` ürün orkestrasyonu, `payment` checkout veya bakiye doğruluğu ve `royalpass` entitlement sahipliğinin açıkça ayrılması
  - season pause veya sistem kaynaklı pasiflikte claim ve reward davranışlarının güvenli yönetimi
  - `notification` ile season başlangıcı, reward claim ve kalan tier bildirimi entegrasyonu

### 21.15) `history` Modülü Kuralları
- Canonical modül adı `history` olmalıdır.
- `history` modülü kullanıcıya ait continue reading, reading history, bookmark-library ve okuma devamlılığı kayıtlarının sahibi olmalıdır.
- `history` modülü user-manga library ilişkisi, user-chapter son okuma durumu, reading session checkpoint, resume anchor ve own history timeline akışlarını taşımalıdır.
- `history` modülü manga metadata, chapter page yapısı, kullanıcı profil verisi veya authorization kararı üretmemelidir.
- `history` veri sahipliği; kullanıcıya ait library entry kayıtları, bookmark veya favorite işaretleri, son okunan chapter veya page referansları, okuma progress snapshot'ları, history timeline verileri ve entry-level share metadata alanları ile sınırlı kalmalıdır.
- `history` modülü public veya shared library görünürlüğünü doğrudan kendi içinde karara bağlamamalı; global visibility default'ları `user`, permission kararları `access`, entry-level share metadata ise `history` içinde kalacak şekilde çalışmalıdır.
- `history` access entegrasyonunda en az `history.continue_reading.read.own`, `history.timeline.read.own`, `history.library.read.own`, `history.bookmark.write.own` ve gerektiğinde `history.library.read.public` gibi canonical permission örnekleri dokümante edilmelidir.
- `history` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - continue reading yüzeyi
  - reading history timeline veya own reading log görünümü
  - bookmark, favorite, library status ve gerektiğinde custom shelf benzeri kayıtlar
  - chapter read start, checkpoint, finish ve resume anchor entegrasyonu
  - cihazlar arası okuma devamlılığı ve duplicate progress yazımına karşı idempotent güncelleme davranışı
  - `manga`, `chapter`, `mission` ve ileride recommendation tarafına kontrollü okuma sinyali veya özet yüzeyi
  - continue reading, library, timeline veya bookmark write alt yüzeylerinin admin tarafından ayrı ayrı açılıp kapatılabilmesi

### 21.16) `ads` Modülü Kuralları
- Canonical modül adı `ads` olmalıdır.
- `ads` modülü reklam yerleşimi, kampanya, kreatif, teslim planı ve gösterim ölçümlemesi alanlarının sahibi olmalıdır.
- `ads` modülü placement, campaign, creative, active window, priority, frequency cap, impression ve click akışlarını taşımalıdır.
- `ads` modülü VIP reklamsız deneyim veya audience erişim kararını kendi içinde üretmemelidir; bu yorumlar `access` ile yapılmalıdır.
- `ads` veri sahipliği; placement tanımları, campaign yapılandırmaları, creative kayıtları, delivery metadata, gösterim veya tıklama logları ve reklam görünürlük state alanları ile sınırlı kalmalıdır.
- `ads` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - ana sayfa, listeleme, manga detay ve chapter çevresi placement yapıları
  - campaign aktiflik penceresi, öncelik ve frequency cap yönetimi
  - impression, click ve temel performans ölçümleme yüzeyleri
  - VIP reklamsız deneyim ile access policy hizası
  - surface, placement veya campaign bazlı admin runtime control anahtarları ve operasyon yüzeyleri

### 21.17) `shop` Modülü Kuralları
- Canonical modül adı `shop` olmalıdır.
- `shop` modülü sellable product veya offer kataloğu, teklif görünürlüğü, fiyatlandırma, satın alma orkestrasyonu ve ürün kullanım kurallarının sahibi olmalıdır.
- `shop` modülü sellable product, offer, price plan, purchase intent, purchase eligibility ve kullanım kuralı akışlarını taşımalıdır.
- `shop` modülü final item sahipliği, equip state veya ledger doğruluğunu kendi içinde üretmemelidir; sahiplik `inventory`, finansal bakiye ve işlem doğruluğu `payment` modülünde kalmalıdır. Stage 29 geçişinde kullanabileceği allowance bridge verisi canonical bakiye owner'lığı sayılmamalıdır.
- `shop` veri sahipliği; sellable product veya offer kataloğu, fiyat tanımı, indirim veya kampanya metadata'sı, purchase request veya order kayıtları, ürün görünürlüğü ve kullanım kısıtı verileri ile sınırlı kalmalıdır.
- `shop` modülü `payment` öncesi aşamada yalnızca purchase eligibility için geçici `seed_mana_allowance_snapshot` veya operasyonel allowance read modelini kullanabilir; bu köprü veri `payment` devreye girince kaldırılmalıdır.
- `shop` modülü ürünleşmiş RoyalPass veya benzeri entitlement ürünlerinde canonical purchase intent giriş noktası olabilir; ancak entitlement sahipliği ilgili hedef modülde, bakiye veya checkout doğruluğu ise `payment` içinde kalmalıdır.
- `shop` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - kozmetik ürün kataloğu, kategori ve slot uyumluluğu
  - mana bazlı fiyatlandırma ve kampanya görünürlüğü
  - satın alma isteği, idempotent purchase davranışı ve duplicate purchase koruması
  - VIP, level, RoyalPass veya görev kaynaklı ürün uygunluğu sinyalleri
  - `inventory` ile final grant veya equip kontratı ve `payment` ile bakiye düşüm mutabakatı
  - katalog, campaign, purchase veya belirli shop alt yüzeylerinin admin tarafından ayrı ayrı açılıp kapatılabilmesi

### 21.18) `payment` Modülü Kuralları
- Canonical modül adı `payment` olmalıdır.
- `payment` modülü mana satın alma, ödeme sağlayıcısı entegrasyonu, cüzdan veya ledger doğruluğu ve finansal işlem kayıtlarının sahibi olmalıdır.
- `payment` modülü mana package, checkout session, provider callback, transaction, ledger entry, balance snapshot ve refund veya reversal akışlarını taşımalıdır.
- `payment` modülü devreye girdiğinde `shop` içindeki geçici `seed_mana_allowance_snapshot` veya operasyonel allowance bridge yüzeyini devralmalı ve canonical bakiye owner'lığını tek başına üstlenmelidir.
- `payment` modülü ürün kataloğu, item sahipliği veya authorization kararı üretmemelidir; katalog `shop`, sahiplik `inventory`, erişim kararı `access` modülünde kalmalıdır.
- `payment` modülü checkout, mana purchase ve finansal doğruluğun sahibidir; ancak ürünleşmiş RoyalPass benzeri entitlement akışlarında final entitlement owner'lığına doğrudan geçmez, onaylanmış ödeme sonucunu ilgili modüle kontrollü kontrat ile aktarır.
- `payment` veri sahipliği; provider session kayıtları, purchase order veya transaction kayıtları, ledger hareketleri, bakiye snapshot'ları, fraud review state ve finansal audit metadata'sı ile sınırlı kalmalıdır.
- `payment` modülünün ana referans dokümanı `docs/modules.md` olmalıdır.
- Geliştirme ile gelen başlıca başlıklar şunlardır:
  - mana satın alma paketleri ve provider checkout oturumları
  - pending, success, failed, cancelled, refunded veya reversed işlem durumları
  - webhook veya callback doğrulaması ve idempotent işleme
  - `shop` ile bakiye düşüm veya mutabakat entegrasyonu
  - fraud review, audit ve finansal görünürlük yüzeyleri
  - mana purchase, checkout veya işlem görüntüleme alt yüzeylerinin admin tarafından kontrollü biçimde daraltılabilmesi

## 22) Dokümantasyon ve Çapraz Kesit Standartları
- Proje geneli dokümantasyon dili Türkçe tutulmalı; dosya başlıkları, bölüm adları ve tablo kolonları tutarlı yazılmalıdır.
- Modül dokümanları en az `Amaç`, `Sorumluluk Alanı`, `Bu Modül Neyi Yapmaz?`, `Veri Sahipliği`, `Bu Modül Hangi Verinin Sahibi Değildir?`, `Access Kontratı`, `API veya Event Sınırı`, `Bağımlılıklar`, `Settings Etkileri`, `Event Akışları`, `Audit ve İzleme`, `İdempotency ve Retry`, `State Yapısı` ve `Test Notları` bölümlerini taşımalıdır.
- Negatif sınır bölümleri boş geçilmemeli; modülün yapmadığı işler ve sahip olmadığı veriler açıkça yazılmalıdır.
- Çapraz kesit kararları için canonical shared dokümanlar `docs/shared.md`, `docs/shared.md`, `docs/shared.md`, `docs/shared.md`, `docs/shared.md`, `docs/shared.md` ve `docs/shared.md` olmalıdır.
- Teknik paket, cache/queue, media, search ve reporting/analytics kararları için aktif referanslar `docs/shared.md`, `docs/shared.md`, `docs/shared.md`, `docs/shared.md` ve `docs/shared.md` olmalıdır.
- Bir modül veya feature dokümanı bu yardımcı altyapı alanlarında aktif sistem kararı üretiyorsa ilgili shared doküman aynı değişiklik setinde güncellenmelidir; karar yalnızca modül içinde not olarak bırakılamaz.
- Enum veya karar sözlüğü niteliğindeki paylaşılan kayıtlar `docs/shared.md`, `docs/shared.md`, `docs/shared.md`, `docs/shared.md`, `docs/shared.md`, `docs/shared.md`, `docs/shared.md`, `docs/shared.md` ve `docs/shared.md` içinde tutulmalıdır.
- Runtime ayar yorumlama sırası `global kill switch -> module/surface availability -> audience selector -> entitlement etkisi -> action policy -> rate limit` biçiminde dokümante edilmeli ve `access` tarafından yorumlanan yüzeyler `docs/shared.md` ile hizalı kalmalıdır.
- `docs/shared.md` yaşayan dokümandır; yeni bir surface eklendiğinde availability anahtarı yanında rate limit, threshold, disabled behavior veya degrade davranışı da aynı değişiklikte yazılmalı ya da neden henüz `planned` kaldığı açıkça not edilmelidir.
- `docs/upgrade.md` ham öneri arşivi olarak değil, uygulanan, kısmi kalan ve bekleyen işlerin durumunu izleyen operasyonel takip belgesi olarak kullanılmalıdır.
- `support` intake tek başına otomatik moderation case sayılmamalı; `support -> moderation` ilişkisi açık handoff politikası ile tanımlanmalıdır.
- `admin` hard override kararı scoped moderator aksiyonunun üzerinde precedence taşır; bu kural modül dokümanlarında ve precedence matrisi içinde aynı şekilde yazılmalıdır.
- Projection veya event tabanlı read model kullanan modüller canonical write model, rebuild yolu, replay desteği ve kabul edilen eventual consistency penceresini `docs/shared.md` ile hizalı yazmalıdır.
- `payment`, `inventory`, `mission`, `royalpass`, `notification`, `support`, `moderation`, `history` ve benzeri event üreticisi modüllerde transactional outbox, retry ve dead-letter yaklaşımı plan dışı bırakılmamalıdır.
- Audit kaydı gereken modüller actor, target, action, result, reason, `correlation_id` ve `request_id` alan setini ortak kullanmalıdır.
- Request ID, correlation ID, rate limit, secret/config ayrımı, backup/restore/rollback ve PII retention kuralları `docs/shared.md` ile hizalı kalmalıdır.
- Test stratejisi, contract test zorunlulukları ve fixture standardı `docs/TESTING.md` içinde tutulmalı; modül dokümanları buradaki katmanlara referans vermelidir.


## 23) Aşama Başlatma ve Tamamlama Zorunlu Akışı
- Bir ajan herhangi bir aşamayı oluşturmaya veya uygulamaya başladığında, aşağıdaki adımları sırasıyla ve eksiksiz uygulamak zorundadır.
- Bu akış, diğer tüm genel kuralların yanında zorunlu operasyonel çalışma akışı olarak kabul edilmelidir.
- Zorunlu sıralı akış:
  - Önce kuralları okur.
  - Ardından proje yapısını inceler.
  - Yapılacak aşamayı docs/ROADMAP.md üzerinden okur.
  - İlgili aşama için uygulanabilir planı hazırlar.
  - Hazırlanan planı gerçekleştirir.
  - Aşamaya ait tüm testleri oluşturur.
  - Oluşturulan tüm testleri başarıyla tamamlar.
  - Docker build alır ve uygulamayı Docker içinde çalıştırır.
  - Versiyonlama işlemlerini bu dokümandaki sürüm/commit/branch kurallarına uygun şekilde uygular.
  - Tüm değişiklikleri Git'e yükler.
  - Son olarak planın tamamlandığını kontrol eder ve kısa, öz bir rapor ile sonucu iletir.
## 30) Dokümantasyon Yapısı ve Güncelleme Kuralı
- Projede aktif ana dokümantasyon seti aşağıdaki beş dosyadan oluşmalıdır:
  - `docs/rules.md`
  - `docs/roadmap.md`
  - `docs/changelog.md`
  - `docs/modules.md`
  - `docs/shared.md`
- Ayrı modül ve shared alt dokümanları yalnızca çalışma taslağı olarak tutulabilir; aktif referans seti yukarıdaki beş dosyadır.
- `rules.md` proje geneli bağlayıcı kuralları taşır.
- `roadmap.md` sistemin oluşturulma sırasını, fazlarını ve teslim sınırlarını taşır.
- `changelog.md` yalnızca projede gerçekten yapılan değişiklikleri kronolojik olarak kaydeder; plan, niyet veya gelecekte yapılacak işler changelog'a yazılmaz.
- `modules.md` tüm modül ownerlıkları, sınırlar, veri sahipliği, bağımlılıklar ve modül bazlı açıklamalar için tek referans dosyadır.
- `shared.md` shared sözlükler, ortak teknik kararlar, settings envanteri ve modüller üstü politikalar için tek referans dosyadır.
- Yeni modül eklendiğinde veya modül ownerlığı değiştiğinde aynı değişiklikte `modules.md` güncellenmelidir.
- Ortak enum, policy, precedence, ayar anahtarı veya operasyon standardı eklendiğinde aynı değişiklikte `shared.md` güncellenmelidir.
- Faz, sıra veya kapsam değiştiğinde aynı değişiklikte `roadmap.md` güncellenmelidir.
- Mimariyi, klasör yapısını, katman sınırlarını veya geliştirme kurallarını etkileyen değişikliklerde `rules.md` güncellenmelidir.
- Proje oluşturulurken veya geliştirme sürecinde gerçekten yapılan işlemler, eklenen dosyalar, çıkarılan dosyalar ve önemli kararlar `changelog.md` içinde kayıt altına alınmalıdır.
- Doküman güncellemesi olmadan yapılan mimari değişiklik tamamlanmış sayılmamalıdır.




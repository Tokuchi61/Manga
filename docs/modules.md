# Modüller

> Bu dosya proje icin canonical modul envanteri ve modul detay referansidir.

## Kullanım Kuralları
- Her modül için canonical ad tek olmalıdır.
- Modül ownerlığı, veri sahipliği ve “neyi yapmaz?” sınırı açık yazılmalıdır.
- Modüller arası ilişki yalnızca kontrollü contract veya event yüzeyi ile kurulmalıdır.
- Final authorization kararı gereken yerlerde `access` modülü referans alınmalıdır.
- Yeni modül eklendiğinde veya ownerlık değiştiğinde bu dosya aynı değişiklikte güncellenmelidir.


## Modül Envanteri

| Canonical Module | Domain Group | Status | Main Doc | Summary |
| --- | --- | --- | --- | --- |
| auth |  | active | docs/modules.md | Kimlik dogrulama, token, session ve hesap guvenligi akislarinin aktif modulu. |
| user |  | active | docs/modules.md | Kullanici hesabi, profil, tercih ve uyelik verisi modulu. |
| access |  | active | docs/modules.md | Merkezi authorization, policy ve erişim kararı modülü. |
| admin |  | planned | docs/modules.md | Yönetim, moderasyon denetimi ve operasyon use-case modülü. |
| manga |  | active | docs/modules.md | Ana içerik varlığı, metadata ve discovery modülü. |
| chapter |  | active | docs/modules.md | Bolum, sayfa ve okuma yuzeyi veri modulu. |
| comment |  | active | docs/modules.md | Icerik yorumlari ve thread etkilesim modulu. |
| support |  | active | docs/modules.md | Kullanıcı destek kayıtları, ticket ve report intake modülü. |
| moderation |  | active | docs/modules.md | Scoped moderasyon kuyrukları ve vaka yönetimi modülü. |
| notification |  | active | docs/modules.md | Bildirim üretimi, teslimi ve tercih yönetimi modülü. |
| social |  | active | docs/modules.md | Takip, arkadaşlık, duvar ve mesajlaşma modülü. |
| inventory |  | active | docs/modules.md | Item sahipliği, claim, consume ve equip modülü. |
| mission |  | active | docs/modules.md | Görev tanımı, ilerleme ve claim eligibility modülü. |
| royalpass |  | planned | docs/modules.md | Sezon, tier ve premium track ilerleme modülü. |
| history |  | active | docs/modules.md | Continue reading, kütüphane ve okuma geçmişi modülü. |
| ads |  | planned | docs/modules.md | Reklam placement, campaign ve ölçümleme modülü. |
| shop |  | planned | docs/modules.md | Ürün kataloğu, offer ve purchase orchestration modülü. |
| payment |  | planned | docs/modules.md | Checkout, ledger ve finansal işlem doğruluğu modülü. |


---

# Access Modülü

> Canonical modül adı: `access`

## Amaç
`access` modülünün amacı, sistemdeki tüm authorization, policy ve erişim kararlarını merkezi, tutarlı ve genişletilebilir şekilde yürütmektir.

## Sorumluluk Alanı
- role ve permission yapısı
- RBAC, policy ve context-aware access kararları
- guest, authenticated, vip ve own/any yorumları için helper sözlüğü
- endpoint guard ve modül bazlı authorization kontratları
- deny reason code standardı ve policy effect sözlüğü
- kısa ömürlü access decision cache planı ve invalidation yaklaşımı

## Bu Modül Neyi Yapmaz?
- kimlik doğrulama veya credential doğrulama yapmaz
- business veri owner'lığı taşıyan modül tablolarını sahiplenmez
- runtime ayar kayıtlarının canonical kaydını kendi içinde saklamaz

## Veri Sahipliği
- rol ve permission sözlüğü
- role-permission ilişkileri
- authorization policy kuralları
- erişim kararına temel olan sözleşme yapıları
- deny reason kodları ve policy effect sözlüğü

## Bu Modül Hangi Verinin Sahibi Değildir?
- auth session veya verification verileri
- user profile, VIP owner'lığı veya social block ilişkisinin ham kaydı
- settings envanteri, audit log ham kaydı veya modül business verileri

## Access Kontratı
`access` modülü access kararının kendisidir. `auth` tarafından doğrulanan kimliği ve `user`, `social`, `moderation` veya `admin` tarafından taşınan sinyalleri kullanarak karar üretir; veri sahipliği o modüllerde kalır. Runtime ayar yorumlama sırası `docs/shared.md`, çatışma çözümü `docs/shared.md`, policy çıktı sözlüğü ise `docs/shared.md` ile hizalı olmalıdır.

## API veya Event Sınırı
- guard, policy ve authorization yüzeyi dış modüller için resmi giriş noktası olmalıdır
- permission ve policy contract yüzeyi açık ve canonical isimlerle yönetilmelidir
- denial, override veya kritik authorization olayları gerektiğinde izlenebilir yüzey üretebilir

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `auth` kimlik doğrulama sinyalleri
- `user` kullanıcı durumu ve üyelik sinyalleri
- `social`, `moderation`, `admin` ve diğer modüllerden gelen scope veya override sinyalleri
- cache standardı için `docs/shared.md`

## Settings Etkileri
- `site.maintenance.enabled`
- `feature.*.enabled` biçimindeki kullanıcıya dönük availability anahtarları
- `feature.user.vip_benefits.enabled` gibi entitlement etkili ayarlar
- access decision cache varsa settings sürümü veya selector değişimi ile invalidation yapılmalıdır

## Event Akışları
- tüketir: auth identity, user visibility ve vip state, social block, moderation deny, admin override sinyalleri
- üretir: `access.policy.changed`, `access.decision.denied` veya metrik odaklı audit sinyalleri
- event tüketiminde replay ve duplicate handling `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- policy override, deny reason değişimi, emergency availability yorumu ve yetki bypass girişimleri auditlenmelidir
- denial metrikleri, cache hit oranı ve precedence kaynaklı ret sayıları izlenebilir olmalıdır

## İdempotency ve Retry
- access kararı salt değerlendirme niteliği taşıdığı için safe retry kabul edilir
- cache doldurma veya invalidation tekrar çalıştırıldığında business yan etki üretmemelidir
- distributed cache kullanılırsa key en az subject, surface, selector ve settings sürümünü kapsamalıdır

## State Yapısı
- role ve permission atama durumu
- policy yorumlama sonucu ve policy effect alanları
- scope ve ownership ayrımları
- deny reason veya availability kararları
- varsa decision cache state'i

## Test Notları
- policy ve permission testleri
- own/any, guest/authenticated ve vip ayrımı testleri
- precedence matrix ve deny reason code doğrulamaları
- decision cache invalidation testleri
- `auth -> access` ve `user -> access` kontrat testleri


---

# Admin Modülü

> Canonical modül adı: `admin`

## Amaç
`admin` modülünün amacı, sistemdeki yönetim, tam yetkili inceleme, merkezi ayar ve operasyon use-case'lerini tek bir yönetim yüzeyinde toplamaktır.

## Sorumluluk Alanı
- admin dashboard ve giriş noktaları
- kullanıcı yönetim akışları
- tam yetkili moderasyon veya support review inceleme akışları
- yüksek riskli handoff, escalation ve yönetimsel inceleme akışları
- merkezi ayarlar, modül açma-kapama ve özellik açma-kapama yönetim yüzeyleri
- risk seviyesi, double confirmation, impersonation ve export-friendly rapor yüzeyleri

## Bu Modül Neyi Yapmaz?
- günlük scoped moderatör akışlarının owner'lığına dönüşmez
- başka modüllerin business verisini kendi tablosuna taşımaz
- access guard olmadan kritik aksiyon çalıştırmaz

## Veri Sahipliği
- admin işlem kayıtları
- admin notları ve operasyonel görünüm verileri
- yönetim use-case akış bilgileri
- runtime ayar tanımları, ayar değişiklik geçmişi ve operasyonel kontrol kayıtları
- admin action risk seviyesi ve confirmation metadata'sı

## Bu Modül Hangi Verinin Sahibi Değildir?
- moderation case'in günlük karar owner'lığı
- support ticket içeriğinin ana owner'lığı
- kullanıcı credential veya ödeme ledger verisi
- audit log'un tek başına yerine geçecek serbest metin notlar

## Access Kontratı
`admin` modülü yetki kararı vermez. Tüm kritik admin akışları `access` modülünün guard veya policy kararları ile korunur. Admin override, reopen, reassignment, freeze veya final karar verdiğinde bu karar scoped moderatör aksiyonunun üzerinde precedence taşır; ancak günlük vaka owner'lığı yine ilgili modülde kalır. Impersonation açılırsa yüksek riskli, audit zorunlu ve zaman sınırlı yürütülmelidir.

## API veya Event Sınırı
- admin yüzeyi tam yetkili yönetim, operasyon ve yönetimsel inceleme use-case'lerini dışa açabilir
- scoped günlük moderatör use-case'leri `moderation` modülünde kalmalı; admin gerektiğinde aynı vaka verisi üzerinde override, handoff veya denetim yürütmelidir
- kritik admin aksiyonları audit veya log yüzeyi ile izlenebilir olmalıdır
- admin modülü diğer modüllerin veri sahipliğini devralmadan orkestrasyon yapmalıdır

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `auth` ile admin kimlik doğrulama entegrasyonu
- `user`, `access`, `moderation`, `support`, `payment` ve diğer modüller ile oversight entegrasyonu
- dashboard, summary ve export tasarımı için `docs/shared.md`

## Settings Etkileri
- settings envanterindeki kategori bazlı tüm admin-owned runtime ayarlar
- high-risk admin action yüzeyleri future key olarak ayrı dokümante edilmelidir
- env veya secret yönetimi gerektiren teknik config admin runtime ayarı olarak sunulmamalıdır

## Event Akışları
- üretir: `admin.setting.changed`, `admin.override.applied`, `admin.user.reviewed`
- tüketir: moderation escalation, support escalation, payment manual review ve operasyon sinyalleri
- admin orkestrasyonunda publish gereken olaylar `docs/shared.md` ile hizalı planlanmalıdır

## Audit ve İzleme
- settings değişikliği, override, impersonation, destructive action ve double confirmation gerektiren aksiyonlar immutable audit log üretmelidir
- admin note ile audit kaydı aynı veri modeli olmamalı; `docs/shared.md` ile ayrıştırılmalıdır

## İdempotency ve Retry
- destructive veya yüksek riskli admin aksiyonları request id ile duplicate koruması taşımalıdır
- aynı override isteği tekrarlandığında ikinci kez yan etki üretmemeli; ilk final state'e bağlanmalıdır
- batch operasyonlarda kısmi başarı ve recovery planı açıkça dokümante edilmelidir

## State Yapısı
- admin işlem durumu
- moderasyon veya support oversight karar yaşam döngüsü
- risk seviyesi veya confirmation metadata'sı
- aktif runtime ayarlar, modül durumu ve özellik durumu görünümleri

## Test Notları
- settings yönetimi ve runtime control testleri
- yetkisiz erişim ve forbidden senaryoları
- high-risk action, double confirmation ve impersonation doğrulamaları
- kritik admin aksiyonlarında audit doğrulamaları


---

# Ads Modülü

> Canonical modül adı: `ads`

## Amaç
`ads` modülünün amacı, reklam placement, kampanya, kreatif ve gösterim ölçümlemesini ayrı bir yaşam döngüsü altında yönetmektir.

## Sorumluluk Alanı
- placement, campaign ve creative yönetimi
- placement taxonomy, targeting ve active window kuralları
- priority, frequency cap ve delivery davranışları
- impression, click, invalid traffic koruması ve temel reklam performans ölçümlemesi
- reporting aggregate ve dashboard veri ihtiyacı
- admin tarafından yönetilen surface, campaign, placement veya click intake runtime kontrolleri ile uyumlu çalışma

## Bu Modül Neyi Yapmaz?
- VIP reklamsız deneyim kararını kendi içinde final hale getirmez
- kullanıcı profili, ödeme veya inventory owner verisini sahiplenmez
- authorization kararını veya genel ürün visibility kuralını tek başına üretmez

## Veri Sahipliği
- placement tanımları
- campaign ve creative kayıtları
- priority, active window ve delivery metadata alanları
- impression veya click logları
- campaign veya placement görünürlük state alanları
- aggregate raporlama için gerekli ham ölçüm metadata'sı

## Bu Modül Hangi Verinin Sahibi Değildir?
- kullanıcı VIP entitlement veya access policy owner verisi
- manga, chapter veya diğer placement hedeflerinin canonical içeriği
- finansal faturalama veya ödeme kayıtları

## Access Kontratı
`ads` yetki kararı vermez. VIP reklamsız deneyim, audience muafiyeti ve reklam yüzeyi görünürlüğü `access` ile yorumlanır. `ads` modülü yalnızca reklam teslimine temel olan veri ve ölçümleme kayıtlarını taşır. VIP no-ads precedence kuralı `docs/shared.md` ile hizalı kalmalıdır.

## API veya Event Sınırı
- placement resolve ve aktif campaign çözümleme yüzeyi
- impression veya click intake yüzeyi
- admin campaign veya placement yönetim yüzeyi
- temel performans raporlama veya export için kontrollü operasyon contract yüzeyi

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `access`, `admin`, `user`, `manga`, `chapter` ve reporting consumer'ları ile entegrasyon
- reporting katmanı için `docs/shared.md`
- cache veya async kararları için `docs/shared.md`

## Settings Etkileri
- `feature.ads.surface.enabled`
- `feature.ads.placement.enabled`
- `feature.ads.campaign.enabled`
- `feature.ads.click_intake.enabled`
- frequency cap metric'leri ayrı key alırsa settings envanteri güncellenmelidir

## Event Akışları
- üretir: `ads.impression.accepted`, `ads.click.accepted`, `ads.campaign.state_changed`
- tüketir: user entitlement veya access availability sinyalleri
- aggregate raporlama ve placement resolve cache'i `docs/shared.md` ile hizalı planlanmalıdır

## Audit ve İzleme
- campaign publish veya pause, targeting override, invalid click koruması ve runtime disable aksiyonları auditlenmelidir
- fill rate, click-through rate, frequency cap hit oranı ve invalid traffic oranı izlenmelidir

## İdempotency ve Retry
- impression veya click intake yüzeyleri tekrar işlendiğinde duplicate aggregate artışı üretmemelidir
- aggregation job tekrarlandığında sayımlar aynı final sonuca gelmelidir
- invalid traffic repair akışları geçmiş logları bozmadan telafi edebilmelidir

## State Yapısı
- draft, active, paused veya ended campaign durumu
- placement görünürlüğü veya delivery aktifliği
- frequency cap veya targeting metadata'sı
- campaign bazlı runtime disable durumu

## Test Notları
- placement çözümleme ve campaign öncelik testleri
- targeting, frequency cap ve invalid click doğrulamaları
- access no-ads precedence entegrasyonu
- aggregate rebuild ve cache davranışı testleri


---

# Auth Modulu

> Canonical modul adi: `auth`

## Amac
`auth` modulunun amaci, sistemde kimlik dogrulama ve oturum guvenligi akislarini merkezi ve guvenli sekilde yurutmektir.

## Sorumluluk Alani
- register, login ve logout akislarinin ownerligi
- access/refresh token yasam dongusu ve refresh rotation
- session list/revoke davranisi (`revoke_current`, `revoke_other_devices`, `revoke_all`)
- email verification, password reset ve password change akislarinin ownerligi
- login guvenlik esikleri, suspicious login sinyalleri ve auth guvenlik olaylari

## Bu Modul Neyi Yapmaz?
- authorization veya policy karari uretmez
- kullanici profili, VIP state, role/permission veya sosyal iliski verisi sahiplenmez
- bildirim teslimi, support vaka ownerligi veya moderation case ownerligi tasimaz

## Veri Sahipligi
- credential benzeri kimlik bilgileri
- auth session kayitlari
- verification/reset/revoke token ve guvenlik state kayitlari
- auth guvenlik olaylari ve operation metadata'si

## Bu Modul Hangi Verinin Sahibi Degildir?
- role, permission ve access policy kayitlari
- kullanici profili, gorunurluk tercihi ve VIP entitlement verisi
- runtime ayar envanterinin canonical kaydi (`docs/shared.md`)

## Access Kontrati
`auth` authorization karari vermez. Dogrulanmis kimlik, session ve guvenlik sinyalleri `access` modulune kontrollu kontrat yuzeyiyle aktarilir. Final allow/deny karari `access`te kalir.

## API veya Event Siniri
- Auth API yuzeyi register/login/logout/session/recovery/verification akislariyla sinirlidir.
- Disa acilan event veya contract payload'lari ownerlik sinirlarini asmaz.
- `user` hesap durumu sinyalleri auth akislarini etkileyebilir; fakat `user` ownerligi devralinmaz.

## Bagimliliklar
- `platform/config`, `platform/database`, `platform/logger`, `platform/validation`
- `shared/crypto/password`
- kontrollu sinirla `user`, `access`, `admin` kontratlari
- migration, mail ve operasyonel izleme altyapisi

## Settings Etkileri
- `auth.login.failed_attempt_limit_per_minute`
- `auth.login.cooldown_seconds`
- `auth.email.verification_resend_cooldown_seconds`
- Yeni auth runtime anahtari acildiginda ayni degisiklikte `docs/shared.md` guncellenir.

## Event Akislari
- Uretir: `auth.login.succeeded`, `auth.login.failed`, `auth.session.revoked`, `auth.email_verification.sent`, `auth.security.suspicious_login`
- Tuketir: `user.account.deactivated`, `admin.password_reset.forced` ve ilgili guvenlik operasyon sinyalleri
- Event publish akislarinda `docs/shared.md` ve idempotent consumer beklentisi zorunludur.

## Audit ve Izleme
- Kritik auth aksiyonlari (`login fail/success`, `revoke_all`, `password reset`, `forced reset`) auditlenir.
- Audit kayitlari ortak alan setini tasir: actor, target, action, result, reason, request_id, correlation_id.
- Guvenlik olaylari operasyonel loglardan ayrik sinifta izlenir.

## Idempotency ve Retry
- Verification resend ve password reset request akislarinda idempotent davranis zorunludur.
- Session revoke akislari tekrar calistiginda yeni yan etki uretmez.
- Login akislarinda duplicate session riski request scope/idempotency stratejisiyle kontrol edilir.

## State Yapisi
- account verification state
- session active/revoked state
- gecici guvenlik kisiti veya cooldown state
- suspicious login ve forced reset state

## Stage 4 Kapsam Haritasi
### 1) Veri modeli ve migration omurgasi
- Auth migration adlandirmasi: `YYYYMMDDHHMM_auth_<aksiyon>.up/down.sql`
- Credential/session/token ownerlik siniri auth modulu disina tasmaz.

### 2) Register/login/logout ve session akislari
- Register/login/logout akislarinin contract-first siniri korunur.
- Session listeleme ve revoke akislari (`current/others/all`) ayri use-case olarak tanimlanir.

### 3) Token/password/verification akislari
- Access/refresh token sure ve rotation kurallari netlestirilir.
- Forgot/reset/change password akislari session invalidate politikasiyla baglanir.
- Email verification baslangici ve resend cooldown standardi korunur.

### 4) Guvenlik, audit ve test gereksinimleri
- Login guvenlik limitleri ve cooldown ayarlari runtime key ile yorumlanir.
- Audit alan seti ve event izlenebilirligi zorunludur.
- Unit/service/contract/integration/e2e test katmanlari `docs/TESTING.md` ile hizali kalir.

## Test Notlari
- register/login/logout basarili/basarisiz senaryolari
- token rotation ve session revoke senaryolari
- forgot/reset/change password akislari
- verification + resend cooldown senaryolari
- `auth -> access` kontrat testleri
- repository integration ve migration smoke testleri


---

# Chapter Modülü

> Canonical modül adı: `chapter`

## Amaç
`chapter` modülünün amacı, manga içeriğinin okunabilir bölüm yapısını, bölüm sayfalarını ve okuma akışının veri temelini yönetmektir.

## Sorumluluk Alanı
- chapter CRUD, detail ve listing akışları
- navigation ve latest chapter görünümü
- page veya media ilişkileri
- preview, detail, tam read ve early access yüzeyleri
- history için gerekli read checkpoint anchor ve resume entegrasyon yüzeyleri
- storage stratejisi ve preview page count kuralı

## Bu Modül Neyi Yapmaz?
- kullanıcıya ait reading progress veya bookmark owner'lığı taşımaz
- ödeme, entitlement veya final visibility kararını tek başına üretmez
- yorum thread owner'lığı veya support kaydı sahiplenmez

## Veri Sahipliği
- chapter metadata alanları
- manga-chapter ilişkisi
- page_number, media referansı ve page yapısı
- access state, preview ve early access veri alanları
- kullanıcıya ait progress kaydı taşımadan resume için gereken canonical chapter veya page anchor bilgileri

## Bu Modül Hangi Verinin Sahibi Değildir?
- history içindeki continue reading veya last read state kayıtları
- payment, inventory veya access policy owner verisi
- support veya moderation vaka owner'lığı

## Access Kontratı
`chapter` erişimi etkileyen veriyi taşır; erişim kararını vermez. Guest, authenticated, vip ve early access kararları `access` tarafından yorumlanır. Admin tarafından yönetilen runtime ayarlar chapter okuma, preview veya early access alt yüzeylerini pasife alabilir; karar uygulaması yine `access` ile korunur. Kullanıcıya ait continue reading, reading history ve bookmark-library state'i `history` modülünde tutulmalıdır.

## API veya Event Sınırı
- chapter detail ve read yüzeyi
- navigation ve latest chapter yüzeyi
- yönetimsel chapter CRUD yüzeyi
- `history` için read start, checkpoint, finish ve resume anchor contract yüzeyi
- `comment` ve `support` için raporlanabilir hedef veya target relation yüzeyi

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `manga` ile parent içerik ilişkisi
- `access`, `admin`, `history`, `comment` ve `support` entegrasyonları
- media storage, attachment policy ve signed URL standardı için `docs/shared.md`

## Settings Etkileri
- `feature.chapter.preview.enabled`
- `feature.chapter.detail.enabled`
- `feature.chapter.read.enabled`
- `feature.chapter.early_access.enabled`
- preview page count gibi limitler ayrı metric key olarak dokümante edilmelidir

## Event Akışları
- üretir: `chapter.published`, `chapter.read.started`, `chapter.read.finished`, `chapter.reordered`
- tüketir: `manga.published_state.changed`, access veya entitlement sinyalleri
- resume anchor ve read event'leri `history` projection'larını besler

## Audit ve İzleme
- reorder, renumber, publish veya access state değişimleri auditlenmelidir
- read surface ile early access görünürlüğü arasındaki ayrım izlenebilir metrik taşımalıdır

## İdempotency ve Retry
- checkpoint intake contract'ı tekrar gönderildiğinde duplicate progress owner'lığı üretmemelidir; owner yine `history` içinde kalır
- renumber veya reorder işlemleri aynı request içinde tekrarlandığında stabil sıra üretmelidir
- media processing retry'ları yeni chapter kaydı yaratmamalıdır

## State Yapısı
- publish durumu
- preview veya early access state'i
- page yapısı ve media referansları
- resume anchor format sürümü

## Test Notları
- CRUD, detail ve read akışı testleri
- navigation ve latest chapter sıralama testleri
- preview, early access ve media doğrulamaları
- `chapter -> history` kontrat ve resume anchor doğrulamaları


---

# Comment Modülü

> Canonical modül adı: `comment`

## Amaç
`comment` modülünün amacı, hedefe bağlı yorum verisini, thread yapısını ve yorum yaşam döngüsünü yönetmektir.

## Sorumluluk Alanı
- comment create, edit, delete ve listing akışları
- root veya reply thread yapısı
- edit window, soft delete görünümü ve delete behavior kuralları
- spoiler, pin, lock ve moderation alanları
- nested reply depth limiti, anti-spam ve rate limit kuralları
- author shadowban veya moderation etkilerinin yorum görünümüne yansıması

## Bu Modül Neyi Yapmaz?
- social duvar post veya social-native reply owner'lığı taşımaz
- erişim veya yetki kararını tek başına vermez
- support veya moderation case kayıtlarını kendi içine dönüştürmez

## Veri Sahipliği
- yorum içeriği
- `target_type` ve `target_id` ilişkisi
- parent-child reply yapısı
- moderation, spoiler, pin ve lock alanları
- edit window ve delete görünümüne ait metadata

## Bu Modül Hangi Verinin Sahibi Değildir?
- hedef içerik varlığının owner verisi
- social messaging veya wall reply native kayıtları
- moderation case, support ticket veya access policy owner'lığı

## Hedef Tipi Sözlüğü
- `target_type` değerleri canonical olarak `docs/shared.md` dosyasındaki kayıtlarla hizalı olmalıdır.
- `comment` hedefleri bu aşamada içerik odaklı hedeflerdir; sosyal duvar post'u veya sosyal duvar reply zinciri `social` modülünde native kalmalı ve örtük olarak comment target'ına dönüştürülmemelidir.
- Yeni yorum hedefi eklendiğinde aynı değişiklik setinde hedef modül dokümanı, `comment` modül dokümanı ve canonical target type kaydı güncellenmelidir.

## Access Kontratı
`comment` yetki kararı vermez. Yorum oluşturma, düzenleme, silme ve görünürlük kararları `access` ile korunur; günlük moderasyon akışları `moderation`, yüksek riskli veya yönetimsel handoff akışları `admin` tarafında yürür. Site geneli, manga detayı veya chapter altı yorum alanlarının açılıp kapatılması ve yorum gönderme aralığı gibi runtime ayarlar admin üzerinden yönetilebilir.

## API veya Event Sınırı
- yorum listeleme ve thread yüzeyi
- yorum detay ve hedef ilişkisi yüzeyi
- moderation, support veya admin odaklı yorum işlemleri için kontrollü yüzey
- yorum modülü read ve write yüzeyleri gerektiğinde birbirinden bağımsız olarak kontrol edilebilir olmalıdır
- hedef tipi genişlemesi canonical target type sözlüğü ile uyumlu tutulmalıdır

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `manga`, `chapter`, `user`, `access`, `admin`, `support` ve `moderation` entegrasyonları
- profanity filter veya anti-spam yardımcıları için ortak altyapı ihtiyacı

## Settings Etkileri
- `feature.comment.read.enabled`
- `feature.comment.write.enabled`
- `comment.write.cooldown_seconds`
- edit window veya spoiler policy gibi ek yüzeyler eklenirse settings envanteri genişletilmelidir

## Event Akışları
- üretir: `comment.created`, `comment.edited`, `comment.deleted`, `comment.moderated`
- tüketir: moderation hide veya restore kararı, admin lock veya unlock kararı
- comment count ve engagement projection'ları `docs/shared.md` ile hizalı beslenmelidir

## Audit ve İzleme
- shadowban etkisi, bulk lock veya pin değişimi ve admin/moderation kaynaklı görünürlük değişimleri auditlenmelidir
- spam ve flood koruma tetiklenmeleri metrik olarak izlenmelidir

## İdempotency ve Retry
- duplicate create istekleri aynı request scope içinde ikinci yorum kaydı oluşturmamalıdır
- delete veya moderation hide aksiyonu tekrarlandığında aynı final görünürlük state'ini korumalıdır
- pagination veya thread rebuild retry'ları yeni yan etki üretmemelidir

## State Yapısı
- moderation_status
- spoiler_flag
- pin ve lock durumu
- visibility etkileyen alanlar
- edit window veya soft delete görünümü durumu

## Test Notları
- create, edit, delete ve list testleri
- thread, reply depth ve soft delete görünümü testleri
- target relation ve spoiler propagation testleri
- rate limit, anti-spam ve shadowban doğrulamaları
- visibility ve access entegrasyonu doğrulamaları


---

# History Modülü

> Canonical modül adı: `history`

## Amaç
`history` modülünün amacı, kullanıcıya ait continue reading, reading history, bookmark-library ve okuma devamlılığı kayıtlarını tek bir yaşam döngüsü altında yönetmektir.

## Sorumluluk Alanı
- continue reading ve resume akışları
- reading history timeline ve own history yüzeyleri
- user-manga library, bookmark veya favorite kayıtları
- resume checkpoint formatı, last read conflict çözümü ve multi-device merge policy
- timeline retention, cleanup ve compact history write stratejisi
- admin tarafından yönetilen continue reading, library, timeline veya bookmark write runtime kontrolleri ile uyumlu çalışma

## Bu Modül Neyi Yapmaz?
- manga veya chapter içeriğinin canonical owner'lığına dönüşmez
- global visibility default'unu tek başına belirlemez
- access policy, social ilişki veya notification preference owner'lığı taşımaz

## Veri Sahipliği
- user-manga library entry kayıtları
- bookmark, favorite veya reading status alanları
- user-chapter son okuma chapter veya page referansları
- progress snapshot, checkpoint ve history timeline verileri
- entry-level visibility veya sharing için gereken history tarafı metadata alanları
- global visibility default'u değil, library entry veya history event bazlı share metadata alanları

## Bu Modül Hangi Verinin Sahibi Değildir?
- manga metadata, chapter page yapısı veya kullanıcı profil verisi
- global user visibility preference sinyali
- access kararının kendisi veya recommendation owner algoritması

## Access Kontratı
`history` yetki kararı vermez. Kullanıcının kendi history veya library yüzeylerini görmesi `access` ile korunur. `access` tarafında en az `history.continue_reading.read.own`, `history.timeline.read.own`, `history.library.read.own`, `history.bookmark.write.own` ve gerektiğinde `history.library.read.public` gibi canonical permission örnekleri tanımlanmalıdır. Public library paylaşımı veya reading activity görünürlüğü gerektiğinde global görünürlük default'ları `user` modülündeki sinyallerden, entry-level paylaşım metadata'sı ise `history` içindeki kayıtlarından yorumlanmalıdır. Final görünürlükte `user` modülündeki global deny veya daha dar default üst sınırdır.

## API veya Event Sınırı
- own continue reading, own library ve own history timeline yüzeyi
- erişim uygunsa kontrollü public library read yüzeyi
- `chapter` modülünden read start, checkpoint, finish veya resume anchor intake kontratı
- `manga`, `mission` ve discovery tüketicileri için kontrollü okuma özeti veya signal surface
- admin tarafından history yüzeylerini daraltmak için kontrollü operasyon contract yüzeyi

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `user`, `manga`, `chapter`, `access`, `admin`, `mission` ve discovery consumer'ları ile entegrasyon
- projection updater ve compact history writer gibi yardımcı yapılar

## Settings Etkileri
- `feature.history.continue_reading.enabled`
- `feature.history.library.enabled`
- `feature.history.timeline.enabled`
- `feature.history.bookmark_write.enabled`
- retention veya cleanup metric'leri ayrı key olarak eklenirse settings envanterine yazılmalıdır

## Event Akışları
- üretir: `history.checkpoint.updated`, `history.chapter.finished`, `history.library.changed`
- tüketir: `chapter.read.started`, `chapter.read.finished`, user visibility sinyalleri
- continue reading projection ve discovery summary'leri `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- public share metadata değişimi, admin kaynaklı visibility daraltması ve bulk cleanup aksiyonları auditlenmelidir
- merge conflict oranı, checkpoint dedup sayısı ve projection lag izlenmelidir

## İdempotency ve Retry
- checkpoint write, merge ve multi-device sync akışları `docs/shared.md` ile hizalı idempotent olmalıdır
- aynı checkpoint tekrar geldiğinde yeni timeline gürültüsü üretmemelidir
- cleanup job veya rebuild retry'ları veri kaybı üretmeden tekrar çalışabilmelidir

## State Yapısı
- in_progress, completed veya dropped reading status
- bookmarked veya favorited library durumu
- last_read_at, resume snapshot ve merge conflict alanları
- timeline retention veya history pause durumu

## Test Notları
- continue reading çözümleme ve checkpoint idempotency testleri
- bookmark ve favorite ayrımı testleri
- duplicate history yazımı, merge conflict ve cleanup pencere doğrulamaları
- own history, public library visibility ve bookmark write permission testleri
- `chapter -> history` ve downstream signal entegrasyonu doğrulamaları


---

# Inventory Modülü

> Canonical modül adı: `inventory`

## Amaç
`inventory` modülünün amacı, kullanıcının sahip olduğu item, kozmetik, ödül ve benzeri varlıkları tek bir envanter yaşam döngüsü altında yönetmek ve ödül teslim yürütümünü merkezi biçimde sahiplenmektir.

## Sorumluluk Alanı
- ownable item definition ve item type ayrımları
- kullanıcı envanter kayıtları ve sahiplik durumları
- grant, reward teslim yürütümü, revoke, consume ve equip akışları
- stackable veya unique item kuralları
- expiry mantığı, equip slot kuralları ve selected cosmetic referansları
- reward source canonical enum ve idempotent grant kontrolü

## Bu Modül Neyi Yapmaz?
- sellable shop product veya offer catalog owner'lığı taşımaz
- ödeme, wallet veya ledger doğruluğu üretmez
- mission veya royalpass claim uygunluğu kararını tek başına vermez

## Veri Sahipliği
- ownable item definition alanları; sellable product veya offer catalog kayıtları değil
- user inventory entry kayıtları
- stack count, ownership state ve expiry bilgileri
- equip slot veya selected cosmetic referansları
- grant, consume ve revoke işlem logları
- reward source ve source reference metadata alanları

## Bu Modül Hangi Verinin Sahibi Değildir?
- shop ürün kataloğu ve fiyat planı verisi
- payment checkout, callback veya ledger kayıtları
- mission veya royalpass progression owner verisi

## Access Kontratı
`inventory` yetki kararı vermez. Kullanıcının kendi envanterini görmesi, reward teslimini tamamlaması, consume etmesi veya equip etmesi `access` ile korunur. `mission`, `royalpass` veya diğer producer modüller claim uygunluğu veya reward kaynağı üretebilir; ancak final grant yürütümü ve item sahipliği `inventory` modülünde kalır. Kaynak tipleri `docs/shared.md` ile hizalı olmalıdır.

## API veya Event Sınırı
- own inventory list, own item detail, own equip veya consume yüzeyi
- admin grant veya revoke orkestrasyonu için kontrollü yönetim yüzeyi
- reward producer modüller için açık grant execution veya reservation contract yüzeyi
- `shop` için sellable product -> grantable item definition çözümleme kontratı

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `user`, `access`, `admin`, `shop`, `notification`, `mission` ve `royalpass` entegrasyonları
- snapshot veya serializer yardımcıları için future-ready altyapı

## Settings Etkileri
- `feature.inventory.read.enabled`
- `feature.inventory.claim.enabled`
- `feature.inventory.equip.enabled`
- `feature.inventory.consume.enabled`
- category bazlı inventory görünürlüğü eklenirse settings envanteri genişletilmelidir

## Event Akışları
- üretir: `inventory.granted`, `inventory.revoked`, `inventory.consumed`, `inventory.equipped`
- tüketir: mission claim, royalpass claim, shop fulfillment, admin manual grant sinyalleri
- reward projection ve grant retry davranışı `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- manual grant, revoke, consume, equip slot değişimi ve correction aksiyonları auditlenmelidir
- duplicate grant koruması, expired item cleanup ve source mismatch oranı izlenmelidir

## İdempotency ve Retry
- grant akışları source reference + request id ile `docs/shared.md` kapsamında korunmalıdır
- consume ve equip tekrarlandığında ikinci yan etki üretmemeli; mevcut final state'e bağlanmalıdır
- recovery veya reconcile akışları yeni sahiplik çoğaltmamalıdır

## State Yapısı
- active, consumed, expired veya revoked item durumu
- stack balance veya quantity durumu
- equipped veya unequipped durumu
- claim veya grant pause durumu
- source type ve source reference alanları

## Test Notları
- grant, claim, consume ve equip testleri
- item type, stackable veya unique ve expiry doğrulamaları
- source idempotency ve duplicate grant testleri
- reward producer entegrasyonu ve access doğrulamaları
- equip slot ve selected cosmetic referans testleri


---

# Manga Modülü

> Canonical modül adı: `manga`

## Amaç
`manga` modülünün amacı, ana içerik varlığını, içerik metadata yapısını ve public içerik akışının veri temelini yönetmektir.

## Sorumluluk Alanı
- manga CRUD, listing ve detail akışları
- arama, filtreleme ve sıralama
- metadata, taxonomy ve slug ilişkileri
- publish, archive, moderation ve görünürlük verileri
- featured, recommendation, editoryal koleksiyon ve keşif yüzeyleri
- content versioning ve search index stratejisi için referans alanları

## Bu Modül Neyi Yapmaz?
- chapter page içeriğini veya kullanıcıya ait history kaydını owner olarak taşımaz
- yorum, support veya moderation vaka verisini kendi içine kopyalamaz
- recommendation verisini tek başına kara kutu algoritma olarak sahiplenmez; kaynak sinyal ve projection stratejisi dokümante edilmelidir

## Veri Sahipliği
- title, slug, summary ve görsel alanları
- taxonomy ilişkileri ve içerik metadata alanları
- publish ve moderation state alanları
- view_count, comment_count, chapter_count ve benzeri denormalize sayaç alanları
- editoryal collection, discovery placement ve recommendation metadata alanları

## Bu Modül Hangi Verinin Sahibi Değildir?
- chapter page veya page media owner verisi
- kullanıcıya ait continue reading, favorite veya bookmark kayıtları
- search provider altyapısının teknik config secret'ları

## Access Kontratı
`manga` erişim kararı vermez. Public görünürlük ve yönetimsel erişim kararları `access` tarafından yorumlanır. Admin tarafından yönetilen runtime ayarlar manga listeleme, detay, recommendation veya editoryal görünürlük yüzeylerini daraltabilir; veri sahipliği yine `manga` modülünde kalır. Visibility kavramı `docs/shared.md` ile hizalı tutulmalıdır.

## API veya Event Sınırı
- public listing ve public detail yüzeyi
- public recommendation, collection veya editoryal discovery yüzeyi
- yönetimsel CRUD ve metadata yüzeyi
- chapter varsayılan erişim verisi gerekiyorsa kontrollü contract yüzeyi
- sayaç güncelleme veya projection ihtiyacı varsa açık event veya counter contract yüzeyi
- `comment` ve `support` için raporlanabilir hedef veya target relation yüzeyi

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `access`, `admin`, `history`, `comment`, `chapter` ve discovery consumer'ları ile kontrollü entegrasyon
- arama ve index lifecycle için `docs/shared.md`

## Settings Etkileri
- `feature.manga.list.enabled`
- `feature.manga.detail.enabled`
- `feature.manga.discovery.enabled`
- arama veya editoryal surface için ayrı key açılırsa settings envanteri aynı değişiklikte güncellenmelidir

## Event Akışları
- üretir: `manga.published`, `manga.archived`, `manga.discovery.updated`
- tüketir: `comment.*`, `chapter.*`, `history.*` kaynaklı sayaç veya engagement sinyalleri
- search index ve recommendation projection'ları `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- publish, archive, visibility ve editoryal collection değişiklikleri auditlenmelidir
- projection lag, search index gecikmesi ve sayaç drift farkları izlenebilir olmalıdır

## İdempotency ve Retry
- publish veya archive transition'ları tekrar çalıştırıldığında çelişkili state üretmemelidir
- counter reconcile ve search reindex akışları replay-safe olmalıdır
- discovery placement güncellemeleri request id veya version alanı ile duplicate write üretmemelidir

## State Yapısı
- publish veya archive durumu
- visibility veya public görünürlük etkileyen alanlar
- collection veya discovery görünürlük state alanları
- taxonomy ve metadata sürüm bilgileri

## Test Notları
- CRUD, listing ve detail testleri
- search, filter, sort ve taxonomy doğrulamaları
- sayaç reconcile ve projection doğrulamaları
- public görünürlük ve access entegrasyonu doğrulamaları


---

# Mission Modülü

> Canonical modül adı: `mission`

## Amaç
`mission` modülünün amacı, günlük, haftalık, aylık, event ve seviye bazlı görev tanımlarını, kullanıcı ilerlemesini ve görev ödülü için claim uygunluğu veya claim request akışlarını yönetmektir.

## Sorumluluk Alanı
- mission definition ve mission category yönetimi
- daily, weekly, monthly, event ve level-based mission yapıları
- mission trigger source listesi ve progress accumulation kuralları
- kullanıcı progress, completion, claim uygunluğu ve claim request yaşam döngüsü
- reset window, streak veya dönemsel yenileme kuralları
- event mission ile season mission ilişkisinin dokümantasyonu

## Bu Modül Neyi Yapmaz?
- final reward sahipliğini kendi içinde yazmaz
- global EXP, level veya kullanıcı profil progression owner'lığına dönüşmez
- ödeme, inventory veya royalpass owner verisini kopyalamaz

## Veri Sahipliği
- mission tanımları ve category alanları
- user mission progress kayıtları
- completion, claim eligibility ve reset durum verileri
- event mission penceresi ve schedule alanları
- reward reference ve claim source bilgileri

## Bu Modül Hangi Verinin Sahibi Değildir?
- final inventory item kaydı veya grant sonucu
- payment veya shop purchase owner verisi
- royalpass season tier owner verisi

## Access Kontratı
`mission` yetki kararı vermez. Görev görüntüleme, claim ve yönetim yüzeyleri `access` ile korunur. Global EXP veya level sahipliği `user` modülünde kalabilir; `mission` bu verileri görev tamamlanma sinyali olarak tüketir. Claim uygunluğu ile final grant ayrımı transaction boundary olarak açık tutulmalıdır.

## API veya Event Sınırı
- own mission list, own progress ve own claim-request yüzeyi
- admin için mission definition, reset ve period kontrol yüzeyi
- producer modüllerden alınan progress event contract'ları ve `inventory` veya `notification` için reward event yüzeyi
- mission read, claim ve progress ingest yüzeyleri gerektiğinde birbirinden bağımsız runtime anahtarları ile yönetilebilir olmalıdır

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `user`, `access`, `admin`, `inventory`, `notification`, `royalpass` ve progress üreticisi modüller
- schedule evaluator, progress aggregator ve rule matcher yardımcıları

## Settings Etkileri
- `mission.daily.reset_hour_utc`
- `feature.mission.read.enabled`
- `feature.mission.claim.enabled`
- `feature.mission.progress_ingest.enabled`
- sezonluk veya event bazlı ek reset pencereleri açılırsa settings envanterinde görünür tutulmalıdır

## Event Akışları
- üretir: `mission.progressed`, `mission.completed`, `mission.claim.requested`, `mission.reset`
- tüketir: okuma, yorum, sosyal etkileşim ve diğer producer event'leri
- progress projection ve season handoff'ları `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- manual completion, reset override, claim rejection ve admin category müdahalesi auditlenmelidir
- stuck progress, duplicate claim ve reset drift oranları izlenebilir olmalıdır

## İdempotency ve Retry
- progress ingest ve claim request akışları `docs/shared.md` ile hizalı olmalıdır
- aynı görev ve dönem için tekrar claim isteği yeni grant zinciri başlatmamalıdır
- reset job tekrarlandığında aynı dönemi iki kez kapatmamalıdır

## State Yapısı
- active, completed, claimed veya expired mission durumu
- reset_pending veya recurring window durumu
- progress snapshot ve objective state alanları
- mission category veya claim surface kapanma durumu
- streak veya period metadata alanları

## Test Notları
- progress ve completion testleri
- trigger source ve accumulation kuralları testleri
- reset, recurring window ve streak testleri
- claim ve reward grant entegrasyonu testleri
- producer event, idempotency ve season ilişkisi doğrulamaları


---

# Moderation Modülü

> Canonical modül adı: `moderation`

## Amaç
`moderation` modülünün amacı, role bazlı veya kullanıcı bazlı scoped moderatörlerin günlük vaka yönetimi, inceleme ve sınırlı müdahale iş akışlarını `admin` modülünden ayrıştırılmış özel bir panel üzerinden yürütmektir.

## Sorumluluk Alanı
- moderatör paneli ve scoped queue yüzeyleri
- role veya kullanıcı bazlı moderatör scope matrisi
- vaka inceleme, assignment, case lifecycle ve stale policy akışları
- yorum, chapter ve manga yüzeyleri için scoped moderasyon iş akışları
- moderatör notu, karar özeti, evidence snapshot ve escalation akışları
- support kaynaklı report intake'lerden gerekli olduğunda linked moderation case açma entegrasyonu

## Bu Modül Neyi Yapmaz?
- global authorization, role veya permission owner'lığı üretmez
- merkezi settings ve kill switch owner'lığına dönüşmez
- admin override'ın üstüne çıkan final karar owner'lığı taşımaz

## Veri Sahipliği
- moderation case veya queue kayıtları
- moderatör assignment bilgileri
- moderatör notları, karar özeti ve vaka timeline verisi
- escalation, stale case ve handoff durum alanları
- evidence snapshot referansları

## Bu Modül Hangi Verinin Sahibi Değildir?
- target içeriğin canonical owner verisi
- access policy, admin settings veya support intake owner verisi
- payment, user veya inventory kaynaklı business kayıtlar

## Hedef Tipi Sözlüğü
- `moderation` case hedefleri canonical olarak `docs/shared.md` dosyasındaki kayıtlarla hizalı olmalıdır.
- Alt yüzey, ekran veya aksiyon bilgisi `target_type` içine gömülmemeli; context verisi ile taşınmalıdır.
- Yeni moderation hedefi eklendiğinde aynı değişiklik setinde hedef modül dokümanı, `moderation` modül dokümanı ve canonical target type kaydı güncellenmelidir.

## Access Kontratı
`moderation` yetki kararı vermez. Comment moderatörü, chapter moderatörü veya manga moderatörü gibi role bazlı ya da kullanıcı bazlı scope kararları `access` tarafından yorumlanır. Admin kullanıcıları aynı vaka verisine tam yönetim yetkisi ile erişebilir; merkezi settings, kill switch ve sistem operasyon yüzeyleri ise `admin` modülünde kalır. Günlük case sahipliği `moderation` içinde kalsa da admin override, reopen, reassignment veya freeze kararı oluştuğunda bu karar moderatör aksiyonunun üzerinde precedence taşır.

## API veya Event Sınırı
- moderatör queue, case detail ve sınırlı karar yürütme yüzeyi
- admin ile orkestrasyon gerektiren escalation veya yüksek riskli handoff yüzeyi
- `support` report kaydı ile linked moderation case ilişki yüzeyi; support kaydı ve moderation case aynı kayıt haline getirilmemelidir
- moderation action olayları gerektiğinde `admin`, `notification` veya hedef modüllere kontrollü event yüzeyi ile aktarılabilir

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `access`, `admin`, `support`, `comment`, `manga`, `chapter` ve `notification` entegrasyonları
- rule-based moderation helper veya content scoring altyapısı için geleceğe dönük ihtiyaç

## Settings Etkileri
- `feature.moderation.panel.enabled`
- `feature.moderation.queue.enabled`
- `feature.moderation.action.enabled`
- scope bazlı queue veya action yüzeyleri ayrıştığında settings envanteri genişletilmelidir

## Event Akışları
- üretir: `moderation.case.created`, `moderation.case.assigned`, `moderation.case.escalated`, `moderation.action.applied`
- tüketir: `support.moderation_handoff_requested`, admin override veya reopen sinyalleri
- queue summary ve action projection'ları `docs/shared.md` ile hizalı planlanmalıdır

## Audit ve İzleme
- case assignment, escalation, evidence snapshot erişimi ve action sonucu auditlenmelidir
- stale case birikimi, queue gecikmesi ve action override oranı izlenebilir olmalıdır

## İdempotency ve Retry
- linked case create akışı aynı support referansı için duplicate moderation case üretmemelidir
- aynı aksiyon tekrarlandığında ikinci kez içerik state'i bozulmamalıdır
- escalation retry'ları mevcut case referansı üzerinden güvenli biçimde yeniden denenmelidir

## State Yapısı
- `docs/shared.md` ile hizalı `case_status`
- `assignment_status`
- `escalation_status`
- `action_result` ve review lifecycle alanları
- queue veya action surface kapanma durumları

## Test Notları
- scoped queue görünürlüğü testleri
- case transition, assignment ve stale policy testleri
- support report'tan linked case açılışı doğrulamaları
- action audit ve admin handoff entegrasyonu testleri
- scope matrix ve precedence doğrulamaları


---

# Notification Modülü

> Canonical modül adı: `notification`

## Amaç
`notification` modülünün amacı, sistem genelindeki olaylardan bildirim üretmek, kullanıcıya uygun kanaldan teslim etmek ve bildirim tercihlerini merkezi şekilde yönetmektir.

## Sorumluluk Alanı
- in-app bildirim kutusu ve read/unread akışları
- bildirim category, template ve channel yönetimi
- in-app, email ve gelecekteki push delivery state ayrımı
- digest, retry, backoff ve dedup stratejileri
- kullanıcıya ait detaylı bildirim tercihleri, mute ve quiet-hour benzeri yüzeyler

## Bu Modül Neyi Yapmaz?
- business event owner'lığı veya authorization kararı üretmez
- kullanıcı profil verisini veya sosyal ilişki verisini sahiplenmez
- producer modül yerine onun iş kuralını kime bildirim gideceği konusunda yeniden tanımlamaz

## Veri Sahipliği
- notification kaydı ve delivery attempt verileri
- notification template, category ve channel tanımları
- kullanıcı bildirim tercihleri, category mute ve quiet-hour alanları
- suppression, dedup ve digest batch metadata'sı

## Bu Modül Hangi Verinin Sahibi Değildir?
- producer modüllerin kaynak event verisi
- kullanıcı hesabı, sosyal blok listesi veya access policy owner'lığı
- dış provider secret veya credential kayıtları

## Access Kontratı
`notification` yetki kararı vermez. Kendi bildirimlerini görme veya yönetme kararı `access` ile korunur. Hangi olayın kime bildirim üreteceği iş kuralı olarak `notification` içinde yorumlanabilir; ancak bu karar authorization yerine teslim kuralı niteliğindedir. Kategori adları `docs/shared.md` ile hizalı olmalıdır.

## API veya Event Sınırı
- own inbox, own preference ve own notification detail yüzeyi
- diğer modüllerden alınan producer event veya notification contract yüzeyi
- admin tarafından category, channel, digest veya delivery pause yönetimi için kontrollü operasyon yüzeyi

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `auth`, `user`, `access`, `admin`, `social`, `mission`, `royalpass`, `support`, `moderation`, `shop`, `payment` ve diğer producer modüller
- queue, retry ve delivery katmanı için `docs/shared.md`

## Settings Etkileri
- `feature.notification.inbox.enabled`
- `feature.notification.preference.enabled`
- `feature.notification.digest.enabled`
- `notification.delivery.paused`
- category veya channel bazlı selector genişlemeleri settings envanterinde açıkça gösterilmelidir

## Event Akışları
- tüketir: producer modüllerden gelen domain event'leri
- üretir: `notification.created`, `notification.delivered`, `notification.failed`, `notification.read`
- unread counter ve digest projection'ları `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- template veya category değişikliği, suppression kararı, delivery failure eşiği ve manual resend aksiyonları auditlenmelidir
- queue lag, unread counter drift ve provider failure oranı izlenmelidir

## İdempotency ve Retry
- aynı producer event aynı kullanıcı ve kategori için duplicate notification üretmemelidir
- delivery retry ve backoff davranışı `docs/shared.md` ile hizalı olmalıdır
- digest üretimi tekrarlandığında aynı batch için ikinci bildirim seti oluşturmamalıdır

## State Yapısı
- created, delivered, failed veya read durumu
- channel veya provider bazlı delivery state'i
- module veya category bazlı delivery pause durumu
- digest eligibility veya quiet-hour bilgisi

## Test Notları
- inbox ve preference akışları
- template, category ve channel çözümleme testleri
- delivery retry, digest ve dedup doğrulamaları
- producer event contract ve unread counter projection testleri


---

# Payment Modülü

> Canonical modül adı: `payment`

## Amaç
`payment` modülünün amacı, mana satın alma, ödeme sağlayıcısı entegrasyonu, ledger doğruluğu ve finansal işlem kayıtlarını merkezi biçimde yönetmektir.

## Sorumluluk Alanı
- mana package ve checkout session akışları
- provider callback veya webhook doğrulaması
- ledger-first işlem modeli, transaction, ledger entry ve balance snapshot yönetimi
- reconciliation job, refund, reversal ve fraud hold hazırlığı
- admin tarafından yönetilen mana purchase, checkout, transaction read veya callback intake runtime kontrolleri ile uyumlu çalışma
- provider adapter, webhook verifier ve money value object yaklaşımı

## Bu Modül Neyi Yapmaz?
- ürün kataloğu, offer görünürlüğü veya inventory sahipliği owner'lığı taşımaz
- RoyalPass gibi entitlement ürünlerinde final entitlement owner'lığına dönüşmez
- authorization kararını tek başına vermez

## Veri Sahipliği
- provider session ve checkout kayıtları
- purchase order ve transaction kayıtları
- ledger entry ve balance snapshot alanları
- provider reference, callback metadata ve audit alanları
- fraud review veya finansal inceleme state alanları
- reconciliation ve reversal metadata'sı

## Bu Modül Hangi Verinin Sahibi Değildir?
- shop ürün kataloğu ve fiyat planı verisi
- inventory item sahipliği veya final reward grant kaydı
- access policy owner verisi

## Access Kontratı
`payment` yetki kararı vermez. Mana satın alma ve işlem görüntüleme aksiyonları `access` ile korunur. Ürün kataloğu `shop`, final item sahipliği `inventory` modülünde kalır. `payment` devreye girdiğinde `shop` içindeki geçici allowance bridge yüzeyi kaldırılmalı ve canonical bakiye yalnızca `payment` içinde tutulmalıdır. Finansal yüzeylerde precedence ve intake pause davranışı `docs/shared.md` ile hizalı kalmalıdır.

## API veya Event Sınırı
- mana package listing ve checkout session başlatma yüzeyi
- provider callback veya webhook intake yüzeyi
- own transaction veya own wallet görünümü için kontrollü okuma yüzeyi
- `shop` ile bakiye düşüm veya mutabakat için kontrollü contract yüzeyi
- entitlement üreten modüller için onaylanmış ödeme veya bakiye güncelleme sonucunu aktaran kontrollü fulfillment contract yüzeyi

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `user`, `access`, `admin`, `shop`, `inventory`, `royalpass` ve payment provider entegrasyonları
- webhook verifier, provider adapter, reconcile job ve money helper altyapısı

## Settings Etkileri
- `feature.payment.mana_purchase.enabled`
- `feature.payment.checkout.enabled`
- `feature.payment.transaction_read.enabled`
- `payment.callback.intake.paused`
- provider bazlı throttling veya manual review surface'i eklenirse settings envanteri genişletilmelidir

## Event Akışları
- üretir: `payment.checkout.started`, `payment.callback.accepted`, `payment.transaction.settled`, `payment.refund.completed`
- tüketir: shop purchase intent, provider callback, admin manual review veya reconcile tetikleyicileri
- ledger, reconcile ve fulfillment publish akışları `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- checkout, callback, refund, reversal, manual review ve reconcile aksiyonları `docs/shared.md` ile auditlenmelidir
- callback retry oranı, ledger drift, snapshot mismatch ve reconcile backlog izlenmelidir

## İdempotency ve Retry
- callback veya webhook işleme, refund ve reversal akışları `docs/shared.md` ile hizalı olmalıdır
- aynı provider event ikinci kez finansal yan etki üretmemelidir
- reconcile job tekrar çalıştığında ledger-first doğruluğu bozmadan eksik state'i toparlayabilmelidir

## State Yapısı
- pending, success, failed, cancelled veya refunded transaction durumu
- authorized veya captured ödeme durumu gerekiyorsa ayrı state
- fraud_review, reversed veya reconciliation_required durumu
- mana purchase surface pause veya provider outage durumu
- balance snapshot ile ledger ilişki metadata'sı

## Test Notları
- checkout session ve callback doğrulama testleri
- ledger, balance snapshot ve mutabakat testleri
- idempotency, replay ve reconcile koruması doğrulamaları
- refund, reversal ve fraud hold senaryoları
- `shop`, `inventory` ve admin inceleme yüzeyi entegrasyon testleri


---

# RoyalPass Modülü

> Canonical modül adı: `royalpass`

## Amaç
`royalpass` modülünün amacı, aylık sezon bazlı pass yapısını, free veya premium track ilerlemesini ve sezon ödülü için claim uygunluğu veya claim request yaşam döngüsünü yönetmektir.

## Sorumluluk Alanı
- season, track ve tier yapıları
- free track ve premium track ayrımı
- season lifecycle, archive görünürlüğü ve tier unlock politikası
- user season progress, claim eligibility ve claim request kayıtları
- mission tabanlı royalpass puanı veya progress entegrasyonu
- premium activation source, paused veya frozen progress ve unclaimed reward davranışı
- cross-season carryover olup olmadığına dair açık kural seti

## Bu Modül Neyi Yapmaz?
- görev tanımı veya inventory sahipliği owner'lığı üretmez
- ödeme ledger'ı veya shop katalog owner'lığına dönüşmez
- entitlement satın alma akışının tamamını kendi içinde kapatmaz

## Veri Sahipliği
- season tanımı ve season state alanları
- tier reward ve track yapılandırmaları
- user season progress, claim eligibility ve premium activation referansları
- claim freeze veya season pause metadata alanları
- unclaimed reward ve carryover policy metadata'sı

## Bu Modül Hangi Verinin Sahibi Değildir?
- shop purchase intent veya payment checkout owner verisi
- final inventory grant kaydı
- mission definition owner verisi

## Access Kontratı
`royalpass` yetki kararı vermez. Season görünürlüğü, claim-request ve premium yüzey erişimi `access` ile korunur. Premium aktivasyonun satın alma kaynağı farklı modüllerden gelebilir; ancak ürünleşmiş satın alma akışında canonical purchase intent `shop` üzerinden başlamalı, checkout veya bakiye doğruluğu gerekiyorsa `payment` tarafından tamamlanmalı ve season içi claim uygunluğu ile progress sahipliği yine `royalpass` modülünde kalmalıdır.

## API veya Event Sınırı
- own season overview, own progress ve own reward claim-request yüzeyi
- admin için season yönetimi, tier yönetimi ve season pause yüzeyi
- `mission` progress tüketimi, `inventory` reward grant ve `notification` bildirim yüzeyleri için kontrollü contract veya event sınırı
- `shop`, `payment` veya admin grant akışlarından gelen premium activation referansları için kontrollü intake contract yüzeyi

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `user`, `access`, `admin`, `mission`, `inventory`, `notification`, `shop` ve `payment` entegrasyonları
- season scheduler, tier progress helper ve reward claim validator yardımcıları

## Settings Etkileri
- `feature.royalpass.claim.enabled`
- `feature.royalpass.season.enabled`
- `feature.royalpass.premium.enabled`
- pause veya carryover yüzeyi ayrı bir runtime kontrol alırsa settings envanteri genişletilmelidir

## Event Akışları
- üretir: `royalpass.progressed`, `royalpass.claim.requested`, `royalpass.season.started`, `royalpass.premium.activated`
- tüketir: `mission.progressed`, `shop.purchase.completed`, `payment.checkout.confirmed`
- tier snapshot projection'ları `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- premium activation, season pause, carryover override ve manual claim müdahaleleri auditlenmelidir
- season drift, stuck claim ve premium entitlement mismatch oranları izlenmelidir

## İdempotency ve Retry
- tier claim ve premium activation intake akışları `docs/shared.md` ile hizalı olmalıdır
- aynı tier için tekrar claim yeni grant üretmemelidir
- season rollover job tekrarlandığında aynı sezonu ikinci kez arşivlememeli veya yeni sezonu iki kez açmamalıdır

## State Yapısı
- draft, active, paused, ended veya archived season durumu
- free veya premium track erişim durumu
- tier claim durumu
- season veya claim surface kapanma durumu
- frozen progress veya carryover policy alanları

## Test Notları
- season ve track çözümleme testleri
- progress, tier unlock ve claim testleri
- premium activation source ve carryover kuralları testleri
- mission, inventory, shop, payment ve notification entegrasyonu testleri
- pause, freeze ve claim recovery doğrulamaları


---

# Shop Modülü

> Canonical modül adı: `shop`

## Amaç
`shop` modülünün amacı, ürün kataloğu, teklif görünürlüğü, fiyatlandırma ve satın alma orkestrasyonunu tek bir mağaza yaşam döngüsü altında yönetmektir.

## Sorumluluk Alanı
- sellable product, offer ve kategori yapıları
- product veya offer ayrımı ve time-limited offer kuralları
- mana bazlı fiyatlandırma, kampanya görünürlüğü ve discount evaluatör ihtiyacı
- purchase intent, eligibility, checkout handoff, fulfillment ve purchase recovery akışları
- already-owned item davranışı ve fail veya retry politikaları
- ürün kullanım kuralları, slot uyumluluğu, katalog görünürlüğü ve inventory item mapping

## Bu Modül Neyi Yapmaz?
- final item sahipliği, equip state veya inventory ledger owner'lığı taşımaz
- ödeme ledger'ı veya provider callback owner'lığına dönüşmez
- entitlement owner modülünün final state'ini kendi tablosunda canonical hale getirmez

## Veri Sahipliği
- sellable product, offer ve kategori kayıtları; ownable inventory item definition alanları değil
- fiyat planı, indirim veya kampanya metadata alanları
- purchase request veya order kayıtları
- VIP, level veya unlock gereksinimi gibi ürün uygunluk kuralları
- ürün görünürlük ve kullanılabilirlik state alanları
- purchase recovery ve orchestration metadata alanları

## Bu Modül Hangi Verinin Sahibi Değildir?
- inventory item sahipliği veya grant sonucu
- payment wallet, ledger, checkout ve callback kayıtları
- royalpass veya başka entitlement modülünün final activation owner verisi

## Access Kontratı
`shop` yetki kararı vermez. Ürün görüntüleme, satın alma ve yönetim aksiyonları `access` ile korunur. Final item sahipliği ve equip state `inventory`, bakiye veya ledger doğruluğu ise `payment` modülünde kalmalıdır. `shop` yalnızca purchase orkestrasyonu için gerekli geçici doğrulama köprülerini taşıyabilir. Satın alma kaynakları `docs/shared.md` ile hizalı olmalıdır.

## Geçiş Notu
- `payment` devreye girene kadar `shop`, yalnızca Stage 29 için tanımlanan geçici `seed_mana_allowance_snapshot` veya operasyonel allowance okuma yüzeyi ile purchase eligibility doğrulayabilir.
- Bu köprü veri canonical wallet, ledger veya gerçek bakiye owner'lığı sayılmaz; `payment` modülü açıldığında kaldırılmalı ve yerini `payment` kontratına bırakmalıdır.

## API veya Event Sınırı
- katalog listing, item detail ve purchase request yüzeyi
- admin katalog, fiyat ve görünürlük yönetim yüzeyi
- `inventory` için final grant veya teslim talep kontratı
- `payment` için bakiye düşüm, reserve veya mutabakat kontratı
- payment öncesi aşamada geçici `seed_mana_allowance_snapshot` okuma kontratı
- `royalpass` gibi entitlement modülleri için ürün bazlı premium activation veya fulfillment handoff kontratı

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `user`, `access`, `admin`, `inventory`, `payment`, `royalpass`, `mission` ve `notification` entegrasyonları
- pricing calculator, discount evaluator ve purchase orchestration helper ihtiyacı

## Settings Etkileri
- `feature.shop.catalog.enabled`
- `feature.shop.purchase.enabled`
- `feature.shop.campaign.enabled`
- recovery veya checkout handoff yüzeyi ayrı kontrol alırsa settings envanterine eklenmelidir

## Event Akışları
- üretir: `shop.purchase.intent.created`, `shop.purchase.completed`, `shop.purchase.recovery_requested`
- tüketir: payment checkout sonucu, inventory grant sonucu, entitlement activation sonucu
- purchase zinciri `docs/shared.md` ve `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- eligibility override, price veya campaign değişimi, manual recovery ve duplicate purchase koruması auditlenmelidir
- checkout handoff başarısı, already-owned red oranı ve recovery backlog izlenmelidir

## İdempotency ve Retry
- purchase intent, recovery replay ve fulfillment handoff akışları `docs/shared.md` ile hizalı olmalıdır
- aynı request için ikinci order veya ikinci fulfillment zinciri başlatılmamalıdır
- already-owned senaryosu retry sırasında yeni order state üretmemelidir

## State Yapısı
- draft, active veya archived product durumu
- visible, hidden veya campaign_only offer durumu
- purchasable, blocked veya sold_out benzeri purchase state alanları
- delivery_pending veya recovery_required satın alma durumu
- checkout handoff ve fulfillment sonuç durumu

## Test Notları
- katalog, product veya offer ayrımı ve fiyat çözümleme testleri
- purchase idempotency ve duplicate request doğrulamaları
- already-owned ve eligibility kuralları testleri
- `inventory`, geçici allowance bridge ve `payment` entegrasyon testleri
- recovery akışı ve runtime control doğrulamaları


---

# Social Modülü

> Canonical modül adı: `social`

## Amaç
`social` modülünün amacı, kullanıcılar arası sosyal ilişki, sosyal duvar ve mesajlaşma yüzeylerini ayrı bir iş alanı olarak yönetmektir.

## Sorumluluk Alanı
- arkadaşlık isteği, kabul, reddetme ve arkadaş listesi
- follow veya unfollow ilişkileri
- friend ve follow farkını ürün yüzeyinde koruyan kullanım kuralları
- sosyal duvar post ve duvar altı reply akışları
- direct message thread, unread state ve mesaj akışları
- sosyal privacy, block, mute veya restrict gibi ilişki kontrol alanları
- online state veya last active ihtiyacı için future-ready alanlar

## Bu Modül Neyi Yapmaz?
- manga veya chapter yorum thread owner'lığı taşımaz
- global authorization kararı üretmez
- notification delivery veya inventory ownership alanına girmez

## Veri Sahipliği
- friendship request ve friendship state verileri
- follow relation kayıtları
- social block, mute veya restrict ilişkileri
- wall post, wall reply, message thread ve message kayıtları
- sosyal görünürlük ve ilişki bazlı erişim sinyalleri

## Bu Modül Hangi Verinin Sahibi Değildir?
- comment modülündeki içerik yorumları
- access policy, user profile owner verisi veya notification preference detayları
- history, support veya moderation case kayıtları

## Access Kontratı
`social` yetki kararını kendi içinde üretmez. Profil duvarını kim görebilir, kimin kime mesaj atabileceği veya arkadaşlık yüzeyine erişim gibi kararlar `social` tarafından üretilen ilişki veya privacy sinyalleri kullanılarak `access` ile korunur. Block veya açık privacy deny sinyali oluştuğunda final sonuç `access` tarafından deny olarak yorumlanmalıdır. `mute` sinyali ise ayrıca dokümante edilmedikçe tek başına genel authorization deny sayılmamalı; teslim, görünürlük veya etkileşim azaltma sinyali olarak kalmalıdır.

## API veya Event Sınırı
- friendship, follow, wall ve direct message yüzeyleri
- sosyal duvar reply yapısı `social`-native kabul edilmeli; `comment` thread sistemine örtük olarak dönüştürülmemelidir
- `notification` ve `mission` için producer event yüzeyi
- admin tarafından ayrı ayrı açılıp kapatılabilen social alt yüzeyler için kontrollü contract

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `auth`, `user`, `access`, `notification`, `admin` ve gerektiğinde `mission` entegrasyonları
- anti-abuse rate limit ve unread counter projection yardımcıları

## Settings Etkileri
- `feature.social.friendship.enabled`
- `feature.social.follow.enabled`
- `feature.social.messaging.enabled`
- `feature.social.wall.enabled`
- restrict veya online presence yüzeyleri eklenirse settings envanteri genişletilmelidir

## Event Akışları
- üretir: `social.friendship.changed`, `social.follow.changed`, `social.message.sent`, `social.wall.posted`
- tüketir: user visibility veya moderation deny sinyalleri
- unread counter ve activity projection'ları `docs/shared.md` ile hizalı planlanmalıdır

## Audit ve İzleme
- block, restrict, messaging abuse ve admin kaynaklı privacy müdahaleleri auditlenmelidir
- spam, burst message ve block override denemeleri izlenebilir olmalıdır

## İdempotency ve Retry
- friendship request, follow veya direct message create akışları duplicate kayıt üretmemelidir
- message send retry'ı aynı message reference için ikinci kayıt oluşturmamalıdır
- block ve mute transition'ları tekrarlandığında aynı final state korunmalıdır

## State Yapısı
- friendship_status
- follow veya unfollow ilişkisi durumu
- wall visibility veya message availability durumu
- block, mute veya restrict ilişkisi durumu
- messaging, wall, friendship veya follow surface kapanma durumu

## Test Notları
- friendship ve follow akış testleri
- friend ve follow davranış ayrımı testleri
- wall post veya reply görünürlüğü testleri
- message thread, unread ve own-surface testleri
- block, mute, restrict ve runtime control entegrasyonu doğrulamaları


---

# Support Modülü

> Canonical modül adı: `support`

## Amaç
`support` modülünün amacı, kullanıcı iletişim taleplerini, destek biletlerini ve `manga`, `chapter`, `comment` hedefleri için içerik bildirimlerini tek bir yaşam döngüsü altında yönetmektir.

## Sorumluluk Alanı
- `communication` kaydı veya genel destek bileti oluşturma akışları
- `manga`, `chapter` ve `comment` için hedefe bağlı içerik bildirimi oluşturma akışları
- own support list, support detail ve support reply akışları
- `support_kind`, `category`, `priority`, `reason_code` ve isteğe bağlı hedef ilişkisi
- SLA, duplicate detection, attachment kuralları ve queue önceliklendirme ihtiyaçları
- internal note ile public reply ayrımı

## Bu Modül Neyi Yapmaz?
- report intake'i otomatik olarak moderation case'e dönüştürmez
- moderation kararının veya admin final override'ının owner'lığına dönüşmez
- hedef içerik veya kullanıcı profil verisinin ana owner'lığına girmez

## Veri Sahipliği
- `support_kind=communication`, `support_kind=ticket` veya `support_kind=report` ile taşınan support kaydı
- `category`, `priority`, `reason_code`, `reason_text` ve isteğe bağlı hedef ilişkisi
- hedefe bağlı kayıtlar için `target_type` ve `target_id`; hedefsiz iletişim veya ticket kayıtlarında boş hedef alanları
- internal note, requester reply ve çözüm metadata'sı

## Bu Modül Hangi Verinin Sahibi Değildir?
- moderation case karar akışı
- hedef içerik varlığının veya yorum kaydının owner verisi
- notification delivery state veya access policy owner'lığı

## Access Kontratı
`support` yetki kararı vermez. Oluşturma, own detail, own reply, review queue ve yönetimsel karar yüzeyleri `access` ile korunur; review ve karar yürütümü `admin` ile entegre çalışır. Communication kaydı, ticket oluşturma, report oluşturma, attachment kabulü ve intake davranışları admin tarafından yönetilen runtime ayarlar ile sınırlandırılabilir. Report kaydı varsayılan olarak moderation case ile aynı kayıt sayılmaz; `docs/shared.md` ile hizalı açık mapping politikası gerekir.

## API veya Event Sınırı
- support create, detail, own list ve own reply yüzeyi
- communication, ticket ve report create yüzeyleri için veri kontratı
- support review queue ve ticket yönetimi için veri yüzeyi
- linked moderation case açılışı için kontrollü handoff veya reference yüzeyi

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `access`, `admin`, `moderation`, `notification`, `manga`, `chapter` ve `comment` ile kontrollü entegrasyon
- attachment policy için `docs/shared.md`
- queue summary ve export yaklaşımı için `docs/shared.md`

## Settings Etkileri
- `feature.support.communication.enabled`
- `feature.support.ticket.enabled`
- `feature.support.report.enabled`
- `feature.support.attachment.enabled`
- `feature.support.internal_note.enabled`
- `support.intake.paused`

## Event Akışları
- üretir: `support.created`, `support.replied`, `support.resolved`, `support.moderation_handoff_requested`
- tüketir: moderation case sonucu, admin karar sinyali, notification template veya delivery sinyali
- linked handoff akışları `docs/shared.md` ve `docs/shared.md` ile hizalı olmalıdır

## Audit ve İzleme
- internal note, resolution, spam işaretleme, SLA breach ve moderation handoff aksiyonları auditlenmelidir
- PII ve attachment riski `docs/shared.md` ile hizalı retention notu taşımalıdır

## İdempotency ve Retry
- duplicate ticket veya report create istekleri requester, target ve request id kombinasyonu ile korunmalıdır
- moderation handoff tekrarlandığında ikinci vaka açılmamalı; linked case referansına bağlanmalıdır
- attachment processing retry'ları yeni support kaydı üretmemelidir

## State Yapısı
- open, pending, resolved veya spam durumu
- duplicate veya spam kontrolüne temel alanlar
- çözüm, reply ve inceleme yaşam döngüsü alanları
- intake pause veya create yüzeyi kapatma durumları

## Test Notları
- communication, ticket ve report ayrımı testleri
- target relation ve handoff doğrulamaları
- duplicate, spam, attachment ve SLA doğrulamaları
- `support -> moderation` ve `support -> notification` kontrat testleri


---

# User Modülü

> Canonical modül adı: `user`

## Amaç
`user` modülünün amacı, kullanıcı hesabı, profil, tercih, görünüm ve üyelik verisini taşıyan merkezi kullanıcı alanını oluşturmaktır.

## Sorumluluk Alanı
- kullanıcı hesabı ve profil alanları
- public veya private profil ayrımı ve profile visibility matrix
- hesap durumu, soft delete, deactivation ve ban etkileri
- üyelik, VIP lifecycle ve benefit freeze sinyalleri
- avatar, banner ve benzeri medya alanlarında ownership veya storage referans kuralları
- global user preference ile module-specific preference ayrımını korumak

## Bu Modül Neyi Yapmaz?
- credential, session veya password lifecycle owner'lığı taşımaz
- detaylı notification preference, social block veya mute listesi owner'lığı taşımaz
- continue reading, history timeline veya inventory sahipliği üretmez

## Veri Sahipliği
- username ve profil alanları
- hesap durumu alanları
- tercih ve görünüm alanları
- üyelik ve kullanıcıya ait hesap verileri
- VIP entitlement süresi ve dondurma veya devam bilgisi
- avatar, banner veya profil kozmetiği için metadata veya referans alanları
- public library veya reading activity görünürlüğünü etkileyen global preference sinyalleri

## Bu Modül Hangi Verinin Sahibi Değildir?
- auth credential, session ve verification token verileri
- social ilişki kayıtları, block veya mute listeleri
- notification category preference detayları
- history kayıtları, entry-level share metadata'sı veya inventory sahipliği

## Access Kontratı
`user` veriyi taşır; access kararı vermez. `access` modülünün yorumlayacağı kullanıcı durumu, üyelik ve görünürlük sinyalleri kontrollü kontrat ile dışa açılır. `user` içindeki global visibility veya sharing preference sinyali, history veya library paylaşımı için üst sınırı belirler; `history` içindeki entry-level metadata bu sınırı daraltabilir ama global deny kararını genişletemez. Visibility kararlarında `docs/shared.md` ve `docs/shared.md` ile hizalı kalınmalıdır.

## API veya Event Sınırı
- user okuma veya yazma yüzeyi profil ve hesap verisi ile sınırlı olmalıdır
- public ve private response yüzeyleri açıkça ayrılmalıdır
- kullanıcıya ait değişiklik olayları gerektiğinde diğer modüllere kontrollü yüzey ile aktarılabilir
- `history` için profile bağlı global visibility veya sharing default preference sinyalleri kontrollü contract ile dışa açılabilir

## Bağımlılıklar
- proje geneli altyapı aşamaları
- `auth` ile kullanıcı kimliği entegrasyonu
- `access`, `admin`, `inventory`, `notification` ve `history` için kontrollü veri okuma yüzeyleri
- medya ownership ve erişim politikası için `docs/shared.md`

## Settings Etkileri
- `feature.user.profile.enabled`
- `feature.user.vip_benefits.enabled`
- `feature.user.vip_badge.enabled`
- `feature.user.history_visibility_preference.enabled`
- yeni media veya moderation görünürlük yüzeyleri eklendiğinde settings envanteri aynı değişiklikte güncellenmelidir

## Event Akışları
- üretir: `user.profile.updated`, `user.visibility.changed`, `user.vip.changed`, `user.account.deactivated`
- tüketir: `auth.identity.created`, `inventory.cosmetic.selected`, `admin.user_state.changed`
- VIP veya profile görünürlük değişimleri `access` ve ilgili consumer modüllere kontrollü event veya contract ile aktarılmalıdır

## Audit ve İzleme
- VIP state değişimi, ban veya deactivation, görünürlük tercihi değişimi ve admin kaynaklı profil müdahalesi auditlenmelidir
- PII ve export riski taşıyan alanlar `docs/shared.md` ile hizalı maskeleme notu taşımalıdır

## İdempotency ve Retry
- profile update akışları request id veya optimistic concurrency ile duplicate write üretmemelidir
- VIP activation veya freeze akışları satın alma veya subscription referansı ile deduplicate edilmelidir
- kullanıcı state transition'ları tekrar çalıştırıldığında çelişkili son durum üretmemelidir

## State Yapısı
- hesap durumu
- profil görünürlüğü ve privacy katmanları
- üyelik veya VIP durum alanları
- sistem kaynaklı VIP pasifliğinde sürenin dondurulmasına ilişkin üyelik durumu

## Test Notları
- profil okuma ve güncelleme testleri
- public veya private response ayrımı testleri
- history visibility precedence doğrulamaları
- ban, deactivation ve VIP lifecycle doğrulamaları

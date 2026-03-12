# ModÃ¼ller

> Bu dosya `docs/modules/*` altÄ±ndaki tÃ¼m modÃ¼l aÃ§Ä±klamalarÄ±nÄ±n tek dosyada birleÅŸtirilmiÅŸ halidir. Yol haritasÄ±ndaki modÃ¼l aÅŸamalarÄ± iÃ§in bu dosya referanstÄ±r.

## KullanÄ±m KurallarÄ±
- Her modÃ¼l iÃ§in canonical ad tek olmalÄ±dÄ±r.
- ModÃ¼l ownerlÄ±ÄŸÄ±, veri sahipliÄŸi ve â€œneyi yapmaz?â€ sÄ±nÄ±rÄ± aÃ§Ä±k yazÄ±lmalÄ±dÄ±r.
- ModÃ¼ller arasÄ± iliÅŸki yalnÄ±zca kontrollÃ¼ contract veya event yÃ¼zeyi ile kurulmalÄ±dÄ±r.
- Final authorization kararÄ± gereken yerlerde `access` modÃ¼lÃ¼ referans alÄ±nmalÄ±dÄ±r.
- Yeni modÃ¼l eklendiÄŸinde veya ownerlÄ±k deÄŸiÅŸtiÄŸinde bu dosya aynÄ± deÄŸiÅŸiklikte gÃ¼ncellenmelidir.


## ModÃ¼l Envanteri

| Canonical Module | Domain Group | Status | Main Doc | Summary |
| --- | --- | --- | --- | --- |
| auth |  | active | docs/modules/auth.md | Kimlik dogrulama, token, session ve hesap guvenligi akislarinin aktif modulu. |
| user |  | active | docs/modules/user.md | Kullanici hesabi, profil, tercih ve uyelik verisi modulu. |
| access |  | active | docs/modules/access.md | Merkezi authorization, policy ve eriÅŸim kararÄ± modÃ¼lÃ¼. |
| admin |  | planned | docs/modules/admin.md | YÃ¶netim, moderasyon denetimi ve operasyon use-case modÃ¼lÃ¼. |
| manga |  | planned | docs/modules/manga.md | Ana iÃ§erik varlÄ±ÄŸÄ±, metadata ve discovery modÃ¼lÃ¼. |
| chapter |  | planned | docs/modules/chapter.md | BÃ¶lÃ¼m, sayfa ve okuma yÃ¼zeyi veri modÃ¼lÃ¼. |
| comment |  | planned | docs/modules/comment.md | Ä°Ã§erik yorumlarÄ± ve thread etkileÅŸim modÃ¼lÃ¼. |
| support |  | planned | docs/modules/support.md | KullanÄ±cÄ± destek kayÄ±tlarÄ±, ticket ve report intake modÃ¼lÃ¼. |
| moderation |  | planned | docs/modules/moderation.md | Scoped moderasyon kuyruklarÄ± ve vaka yÃ¶netimi modÃ¼lÃ¼. |
| notification |  | planned | docs/modules/notification.md | Bildirim Ã¼retimi, teslimi ve tercih yÃ¶netimi modÃ¼lÃ¼. |
| social |  | planned | docs/modules/social.md | Takip, arkadaÅŸlÄ±k, duvar ve mesajlaÅŸma modÃ¼lÃ¼. |
| inventory |  | planned | docs/modules/inventory.md | Item sahipliÄŸi, claim, consume ve equip modÃ¼lÃ¼. |
| mission |  | planned | docs/modules/mission.md | GÃ¶rev tanÄ±mÄ±, ilerleme ve claim eligibility modÃ¼lÃ¼. |
| royalpass |  | planned | docs/modules/royalpass.md | Sezon, tier ve premium track ilerleme modÃ¼lÃ¼. |
| history |  | planned | docs/modules/history.md | Continue reading, kÃ¼tÃ¼phane ve okuma geÃ§miÅŸi modÃ¼lÃ¼. |
| ads |  | planned | docs/modules/ads.md | Reklam placement, campaign ve Ã¶lÃ§Ã¼mleme modÃ¼lÃ¼. |
| shop |  | planned | docs/modules/shop.md | ÃœrÃ¼n kataloÄŸu, offer ve purchase orchestration modÃ¼lÃ¼. |
| payment |  | planned | docs/modules/payment.md | Checkout, ledger ve finansal iÅŸlem doÄŸruluÄŸu modÃ¼lÃ¼. |


---

# Access ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `access`

## AmaÃ§
`access` modÃ¼lÃ¼nÃ¼n amacÄ±, sistemdeki tÃ¼m authorization, policy ve eriÅŸim kararlarÄ±nÄ± merkezi, tutarlÄ± ve geniÅŸletilebilir ÅŸekilde yÃ¼rÃ¼tmektir.

## Sorumluluk AlanÄ±
- role ve permission yapÄ±sÄ±
- RBAC, policy ve context-aware access kararlarÄ±
- guest, authenticated, vip ve own/any yorumlarÄ± iÃ§in helper sÃ¶zlÃ¼ÄŸÃ¼
- endpoint guard ve modÃ¼l bazlÄ± authorization kontratlarÄ±
- deny reason code standardÄ± ve policy effect sÃ¶zlÃ¼ÄŸÃ¼
- kÄ±sa Ã¶mÃ¼rlÃ¼ access decision cache planÄ± ve invalidation yaklaÅŸÄ±mÄ±

## Bu ModÃ¼l Neyi Yapmaz?
- kimlik doÄŸrulama veya credential doÄŸrulama yapmaz
- business veri owner'lÄ±ÄŸÄ± taÅŸÄ±yan modÃ¼l tablolarÄ±nÄ± sahiplenmez
- runtime ayar kayÄ±tlarÄ±nÄ±n canonical kaydÄ±nÄ± kendi iÃ§inde saklamaz

## Veri SahipliÄŸi
- rol ve permission sÃ¶zlÃ¼ÄŸÃ¼
- role-permission iliÅŸkileri
- authorization policy kurallarÄ±
- eriÅŸim kararÄ±na temel olan sÃ¶zleÅŸme yapÄ±larÄ±
- deny reason kodlarÄ± ve policy effect sÃ¶zlÃ¼ÄŸÃ¼

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- auth session veya verification verileri
- user profile, VIP owner'lÄ±ÄŸÄ± veya social block iliÅŸkisinin ham kaydÄ±
- settings envanteri, audit log ham kaydÄ± veya modÃ¼l business verileri

## Access KontratÄ±
`access` modÃ¼lÃ¼ access kararÄ±nÄ±n kendisidir. `auth` tarafÄ±ndan doÄŸrulanan kimliÄŸi ve `user`, `social`, `moderation` veya `admin` tarafÄ±ndan taÅŸÄ±nan sinyalleri kullanarak karar Ã¼retir; veri sahipliÄŸi o modÃ¼llerde kalÄ±r. Runtime ayar yorumlama sÄ±rasÄ± `docs/settings/index.md`, Ã§atÄ±ÅŸma Ã§Ã¶zÃ¼mÃ¼ `docs/shared/precedence-rules.md`, policy Ã§Ä±ktÄ± sÃ¶zlÃ¼ÄŸÃ¼ ise `docs/shared/policy-effects.md` ile hizalÄ± olmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- guard, policy ve authorization yÃ¼zeyi dÄ±ÅŸ modÃ¼ller iÃ§in resmi giriÅŸ noktasÄ± olmalÄ±dÄ±r
- permission ve policy contract yÃ¼zeyi aÃ§Ä±k ve canonical isimlerle yÃ¶netilmelidir
- denial, override veya kritik authorization olaylarÄ± gerektiÄŸinde izlenebilir yÃ¼zey Ã¼retebilir

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `auth` kimlik doÄŸrulama sinyalleri
- `user` kullanÄ±cÄ± durumu ve Ã¼yelik sinyalleri
- `social`, `moderation`, `admin` ve diÄŸer modÃ¼llerden gelen scope veya override sinyalleri
- cache standardÄ± iÃ§in `docs/shared/cache-queue-strategy.md`

## Settings Etkileri
- `site.maintenance.enabled`
- `feature.*.enabled` biÃ§imindeki kullanÄ±cÄ±ya dÃ¶nÃ¼k availability anahtarlarÄ±
- `feature.user.vip_benefits.enabled` gibi entitlement etkili ayarlar
- access decision cache varsa settings sÃ¼rÃ¼mÃ¼ veya selector deÄŸiÅŸimi ile invalidation yapÄ±lmalÄ±dÄ±r

## Event AkÄ±ÅŸlarÄ±
- tÃ¼ketir: auth identity, user visibility ve vip state, social block, moderation deny, admin override sinyalleri
- Ã¼retir: `access.policy.changed`, `access.decision.denied` veya metrik odaklÄ± audit sinyalleri
- event tÃ¼ketiminde replay ve duplicate handling `docs/shared/idempotency-policy.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- policy override, deny reason deÄŸiÅŸimi, emergency availability yorumu ve yetki bypass giriÅŸimleri auditlenmelidir
- denial metrikleri, cache hit oranÄ± ve precedence kaynaklÄ± ret sayÄ±larÄ± izlenebilir olmalÄ±dÄ±r

## Ä°dempotency ve Retry
- access kararÄ± salt deÄŸerlendirme niteliÄŸi taÅŸÄ±dÄ±ÄŸÄ± iÃ§in safe retry kabul edilir
- cache doldurma veya invalidation tekrar Ã§alÄ±ÅŸtÄ±rÄ±ldÄ±ÄŸÄ±nda business yan etki Ã¼retmemelidir
- distributed cache kullanÄ±lÄ±rsa key en az subject, surface, selector ve settings sÃ¼rÃ¼mÃ¼nÃ¼ kapsamalÄ±dÄ±r

## State YapÄ±sÄ±
- role ve permission atama durumu
- policy yorumlama sonucu ve policy effect alanlarÄ±
- scope ve ownership ayrÄ±mlarÄ±
- deny reason veya availability kararlarÄ±
- varsa decision cache state'i

## Test NotlarÄ±
- policy ve permission testleri
- own/any, guest/authenticated ve vip ayrÄ±mÄ± testleri
- precedence matrix ve deny reason code doÄŸrulamalarÄ±
- decision cache invalidation testleri
- `auth -> access` ve `user -> access` kontrat testleri


---

# Admin ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `admin`

## AmaÃ§
`admin` modÃ¼lÃ¼nÃ¼n amacÄ±, sistemdeki yÃ¶netim, tam yetkili inceleme, merkezi ayar ve operasyon use-case'lerini tek bir yÃ¶netim yÃ¼zeyinde toplamaktÄ±r.

## Sorumluluk AlanÄ±
- admin dashboard ve giriÅŸ noktalarÄ±
- kullanÄ±cÄ± yÃ¶netim akÄ±ÅŸlarÄ±
- tam yetkili moderasyon veya support review inceleme akÄ±ÅŸlarÄ±
- yÃ¼ksek riskli handoff, escalation ve yÃ¶netimsel inceleme akÄ±ÅŸlarÄ±
- merkezi ayarlar, modÃ¼l aÃ§ma-kapama ve Ã¶zellik aÃ§ma-kapama yÃ¶netim yÃ¼zeyleri
- risk seviyesi, double confirmation, impersonation ve export-friendly rapor yÃ¼zeyleri

## Bu ModÃ¼l Neyi Yapmaz?
- gÃ¼nlÃ¼k scoped moderatÃ¶r akÄ±ÅŸlarÄ±nÄ±n owner'lÄ±ÄŸÄ±na dÃ¶nÃ¼ÅŸmez
- baÅŸka modÃ¼llerin business verisini kendi tablosuna taÅŸÄ±maz
- access guard olmadan kritik aksiyon Ã§alÄ±ÅŸtÄ±rmaz

## Veri SahipliÄŸi
- admin iÅŸlem kayÄ±tlarÄ±
- admin notlarÄ± ve operasyonel gÃ¶rÃ¼nÃ¼m verileri
- yÃ¶netim use-case akÄ±ÅŸ bilgileri
- runtime ayar tanÄ±mlarÄ±, ayar deÄŸiÅŸiklik geÃ§miÅŸi ve operasyonel kontrol kayÄ±tlarÄ±
- admin action risk seviyesi ve confirmation metadata'sÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- moderation case'in gÃ¼nlÃ¼k karar owner'lÄ±ÄŸÄ±
- support ticket iÃ§eriÄŸinin ana owner'lÄ±ÄŸÄ±
- kullanÄ±cÄ± credential veya Ã¶deme ledger verisi
- audit log'un tek baÅŸÄ±na yerine geÃ§ecek serbest metin notlar

## Access KontratÄ±
`admin` modÃ¼lÃ¼ yetki kararÄ± vermez. TÃ¼m kritik admin akÄ±ÅŸlarÄ± `access` modÃ¼lÃ¼nÃ¼n guard veya policy kararlarÄ± ile korunur. Admin override, reopen, reassignment, freeze veya final karar verdiÄŸinde bu karar scoped moderatÃ¶r aksiyonunun Ã¼zerinde precedence taÅŸÄ±r; ancak gÃ¼nlÃ¼k vaka owner'lÄ±ÄŸÄ± yine ilgili modÃ¼lde kalÄ±r. Impersonation aÃ§Ä±lÄ±rsa yÃ¼ksek riskli, audit zorunlu ve zaman sÄ±nÄ±rlÄ± yÃ¼rÃ¼tÃ¼lmelidir.

## API veya Event SÄ±nÄ±rÄ±
- admin yÃ¼zeyi tam yetkili yÃ¶netim, operasyon ve yÃ¶netimsel inceleme use-case'lerini dÄ±ÅŸa aÃ§abilir
- scoped gÃ¼nlÃ¼k moderatÃ¶r use-case'leri `moderation` modÃ¼lÃ¼nde kalmalÄ±; admin gerektiÄŸinde aynÄ± vaka verisi Ã¼zerinde override, handoff veya denetim yÃ¼rÃ¼tmelidir
- kritik admin aksiyonlarÄ± audit veya log yÃ¼zeyi ile izlenebilir olmalÄ±dÄ±r
- admin modÃ¼lÃ¼ diÄŸer modÃ¼llerin veri sahipliÄŸini devralmadan orkestrasyon yapmalÄ±dÄ±r

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `auth` ile admin kimlik doÄŸrulama entegrasyonu
- `user`, `access`, `moderation`, `support`, `payment` ve diÄŸer modÃ¼ller ile oversight entegrasyonu
- dashboard, summary ve export tasarÄ±mÄ± iÃ§in `docs/shared/reporting-analytics-strategy.md`

## Settings Etkileri
- settings envanterindeki kategori bazlÄ± tÃ¼m admin-owned runtime ayarlar
- high-risk admin action yÃ¼zeyleri future key olarak ayrÄ± dokÃ¼mante edilmelidir
- env veya secret yÃ¶netimi gerektiren teknik config admin runtime ayarÄ± olarak sunulmamalÄ±dÄ±r

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `admin.setting.changed`, `admin.override.applied`, `admin.user.reviewed`
- tÃ¼ketir: moderation escalation, support escalation, payment manual review ve operasyon sinyalleri
- admin orkestrasyonunda publish gereken olaylar `docs/shared/outbox-pattern.md` ile hizalÄ± planlanmalÄ±dÄ±r

## Audit ve Ä°zleme
- settings deÄŸiÅŸikliÄŸi, override, impersonation, destructive action ve double confirmation gerektiren aksiyonlar immutable audit log Ã¼retmelidir
- admin note ile audit kaydÄ± aynÄ± veri modeli olmamalÄ±; `docs/shared/audit-policy.md` ile ayrÄ±ÅŸtÄ±rÄ±lmalÄ±dÄ±r

## Ä°dempotency ve Retry
- destructive veya yÃ¼ksek riskli admin aksiyonlarÄ± request id ile duplicate korumasÄ± taÅŸÄ±malÄ±dÄ±r
- aynÄ± override isteÄŸi tekrarlandÄ±ÄŸÄ±nda ikinci kez yan etki Ã¼retmemeli; ilk final state'e baÄŸlanmalÄ±dÄ±r
- batch operasyonlarda kÄ±smi baÅŸarÄ± ve recovery planÄ± aÃ§Ä±kÃ§a dokÃ¼mante edilmelidir

## State YapÄ±sÄ±
- admin iÅŸlem durumu
- moderasyon veya support oversight karar yaÅŸam dÃ¶ngÃ¼sÃ¼
- risk seviyesi veya confirmation metadata'sÄ±
- aktif runtime ayarlar, modÃ¼l durumu ve Ã¶zellik durumu gÃ¶rÃ¼nÃ¼mleri

## Test NotlarÄ±
- settings yÃ¶netimi ve runtime control testleri
- yetkisiz eriÅŸim ve forbidden senaryolarÄ±
- high-risk action, double confirmation ve impersonation doÄŸrulamalarÄ±
- kritik admin aksiyonlarÄ±nda audit doÄŸrulamalarÄ±


---

# Ads ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `ads`

## AmaÃ§
`ads` modÃ¼lÃ¼nÃ¼n amacÄ±, reklam placement, kampanya, kreatif ve gÃ¶sterim Ã¶lÃ§Ã¼mlemesini ayrÄ± bir yaÅŸam dÃ¶ngÃ¼sÃ¼ altÄ±nda yÃ¶netmektir.

## Sorumluluk AlanÄ±
- placement, campaign ve creative yÃ¶netimi
- placement taxonomy, targeting ve active window kurallarÄ±
- priority, frequency cap ve delivery davranÄ±ÅŸlarÄ±
- impression, click, invalid traffic korumasÄ± ve temel reklam performans Ã¶lÃ§Ã¼mlemesi
- reporting aggregate ve dashboard veri ihtiyacÄ±
- admin tarafÄ±ndan yÃ¶netilen surface, campaign, placement veya click intake runtime kontrolleri ile uyumlu Ã§alÄ±ÅŸma

## Bu ModÃ¼l Neyi Yapmaz?
- VIP reklamsÄ±z deneyim kararÄ±nÄ± kendi iÃ§inde final hale getirmez
- kullanÄ±cÄ± profili, Ã¶deme veya inventory owner verisini sahiplenmez
- authorization kararÄ±nÄ± veya genel Ã¼rÃ¼n visibility kuralÄ±nÄ± tek baÅŸÄ±na Ã¼retmez

## Veri SahipliÄŸi
- placement tanÄ±mlarÄ±
- campaign ve creative kayÄ±tlarÄ±
- priority, active window ve delivery metadata alanlarÄ±
- impression veya click loglarÄ±
- campaign veya placement gÃ¶rÃ¼nÃ¼rlÃ¼k state alanlarÄ±
- aggregate raporlama iÃ§in gerekli ham Ã¶lÃ§Ã¼m metadata'sÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- kullanÄ±cÄ± VIP entitlement veya access policy owner verisi
- manga, chapter veya diÄŸer placement hedeflerinin canonical iÃ§eriÄŸi
- finansal faturalama veya Ã¶deme kayÄ±tlarÄ±

## Access KontratÄ±
`ads` yetki kararÄ± vermez. VIP reklamsÄ±z deneyim, audience muafiyeti ve reklam yÃ¼zeyi gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ `access` ile yorumlanÄ±r. `ads` modÃ¼lÃ¼ yalnÄ±zca reklam teslimine temel olan veri ve Ã¶lÃ§Ã¼mleme kayÄ±tlarÄ±nÄ± taÅŸÄ±r. VIP no-ads precedence kuralÄ± `docs/shared/precedence-rules.md` ile hizalÄ± kalmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- placement resolve ve aktif campaign Ã§Ã¶zÃ¼mleme yÃ¼zeyi
- impression veya click intake yÃ¼zeyi
- admin campaign veya placement yÃ¶netim yÃ¼zeyi
- temel performans raporlama veya export iÃ§in kontrollÃ¼ operasyon contract yÃ¼zeyi

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `access`, `admin`, `user`, `manga`, `chapter` ve reporting consumer'larÄ± ile entegrasyon
- reporting katmanÄ± iÃ§in `docs/shared/reporting-analytics-strategy.md`
- cache veya async kararlarÄ± iÃ§in `docs/shared/cache-queue-strategy.md`

## Settings Etkileri
- `feature.ads.surface.enabled`
- `feature.ads.placement.enabled`
- `feature.ads.campaign.enabled`
- `feature.ads.click_intake.enabled`
- frequency cap metric'leri ayrÄ± key alÄ±rsa settings envanteri gÃ¼ncellenmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `ads.impression.accepted`, `ads.click.accepted`, `ads.campaign.state_changed`
- tÃ¼ketir: user entitlement veya access availability sinyalleri
- aggregate raporlama ve placement resolve cache'i `docs/shared/projection-strategy.md` ile hizalÄ± planlanmalÄ±dÄ±r

## Audit ve Ä°zleme
- campaign publish veya pause, targeting override, invalid click korumasÄ± ve runtime disable aksiyonlarÄ± auditlenmelidir
- fill rate, click-through rate, frequency cap hit oranÄ± ve invalid traffic oranÄ± izlenmelidir

## Ä°dempotency ve Retry
- impression veya click intake yÃ¼zeyleri tekrar iÅŸlendiÄŸinde duplicate aggregate artÄ±ÅŸÄ± Ã¼retmemelidir
- aggregation job tekrarlandÄ±ÄŸÄ±nda sayÄ±mlar aynÄ± final sonuca gelmelidir
- invalid traffic repair akÄ±ÅŸlarÄ± geÃ§miÅŸ loglarÄ± bozmadan telafi edebilmelidir

## State YapÄ±sÄ±
- draft, active, paused veya ended campaign durumu
- placement gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ veya delivery aktifliÄŸi
- frequency cap veya targeting metadata'sÄ±
- campaign bazlÄ± runtime disable durumu

## Test NotlarÄ±
- placement Ã§Ã¶zÃ¼mleme ve campaign Ã¶ncelik testleri
- targeting, frequency cap ve invalid click doÄŸrulamalarÄ±
- access no-ads precedence entegrasyonu
- aggregate rebuild ve cache davranÄ±ÅŸÄ± testleri


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
- runtime ayar envanterinin canonical kaydi (`docs/settings/index.md`)

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
- Yeni auth runtime anahtari acildiginda ayni degisiklikte `docs/settings/index.md` guncellenir.

## Event Akislari
- Uretir: `auth.login.succeeded`, `auth.login.failed`, `auth.session.revoked`, `auth.email_verification.sent`, `auth.security.suspicious_login`
- Tuketir: `user.account.deactivated`, `admin.password_reset.forced` ve ilgili guvenlik operasyon sinyalleri
- Event publish akislarinda `docs/shared/outbox-pattern.md` ve idempotent consumer beklentisi zorunludur.

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

# Chapter ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `chapter`

## AmaÃ§
`chapter` modÃ¼lÃ¼nÃ¼n amacÄ±, manga iÃ§eriÄŸinin okunabilir bÃ¶lÃ¼m yapÄ±sÄ±nÄ±, bÃ¶lÃ¼m sayfalarÄ±nÄ± ve okuma akÄ±ÅŸÄ±nÄ±n veri temelini yÃ¶netmektir.

## Sorumluluk AlanÄ±
- chapter CRUD, detail ve listing akÄ±ÅŸlarÄ±
- navigation ve latest chapter gÃ¶rÃ¼nÃ¼mÃ¼
- page veya media iliÅŸkileri
- preview, detail, tam read ve early access yÃ¼zeyleri
- history iÃ§in gerekli read checkpoint anchor ve resume entegrasyon yÃ¼zeyleri
- storage stratejisi ve preview page count kuralÄ±

## Bu ModÃ¼l Neyi Yapmaz?
- kullanÄ±cÄ±ya ait reading progress veya bookmark owner'lÄ±ÄŸÄ± taÅŸÄ±maz
- Ã¶deme, entitlement veya final visibility kararÄ±nÄ± tek baÅŸÄ±na Ã¼retmez
- yorum thread owner'lÄ±ÄŸÄ± veya support kaydÄ± sahiplenmez

## Veri SahipliÄŸi
- chapter metadata alanlarÄ±
- manga-chapter iliÅŸkisi
- page_number, media referansÄ± ve page yapÄ±sÄ±
- access state, preview ve early access veri alanlarÄ±
- kullanÄ±cÄ±ya ait progress kaydÄ± taÅŸÄ±madan resume iÃ§in gereken canonical chapter veya page anchor bilgileri

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- history iÃ§indeki continue reading veya last read state kayÄ±tlarÄ±
- payment, inventory veya access policy owner verisi
- support veya moderation vaka owner'lÄ±ÄŸÄ±

## Access KontratÄ±
`chapter` eriÅŸimi etkileyen veriyi taÅŸÄ±r; eriÅŸim kararÄ±nÄ± vermez. Guest, authenticated, vip ve early access kararlarÄ± `access` tarafÄ±ndan yorumlanÄ±r. Admin tarafÄ±ndan yÃ¶netilen runtime ayarlar chapter okuma, preview veya early access alt yÃ¼zeylerini pasife alabilir; karar uygulamasÄ± yine `access` ile korunur. KullanÄ±cÄ±ya ait continue reading, reading history ve bookmark-library state'i `history` modÃ¼lÃ¼nde tutulmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- chapter detail ve read yÃ¼zeyi
- navigation ve latest chapter yÃ¼zeyi
- yÃ¶netimsel chapter CRUD yÃ¼zeyi
- `history` iÃ§in read start, checkpoint, finish ve resume anchor contract yÃ¼zeyi
- `comment` ve `support` iÃ§in raporlanabilir hedef veya target relation yÃ¼zeyi

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `manga` ile parent iÃ§erik iliÅŸkisi
- `access`, `admin`, `history`, `comment` ve `support` entegrasyonlarÄ±
- media storage, attachment policy ve signed URL standardÄ± iÃ§in `docs/shared/media-asset-strategy.md`

## Settings Etkileri
- `feature.chapter.preview.enabled`
- `feature.chapter.detail.enabled`
- `feature.chapter.read.enabled`
- `feature.chapter.early_access.enabled`
- preview page count gibi limitler ayrÄ± metric key olarak dokÃ¼mante edilmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `chapter.published`, `chapter.read.started`, `chapter.read.finished`, `chapter.reordered`
- tÃ¼ketir: `manga.published_state.changed`, access veya entitlement sinyalleri
- resume anchor ve read event'leri `history` projection'larÄ±nÄ± besler

## Audit ve Ä°zleme
- reorder, renumber, publish veya access state deÄŸiÅŸimleri auditlenmelidir
- read surface ile early access gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ arasÄ±ndaki ayrÄ±m izlenebilir metrik taÅŸÄ±malÄ±dÄ±r

## Ä°dempotency ve Retry
- checkpoint intake contract'Ä± tekrar gÃ¶nderildiÄŸinde duplicate progress owner'lÄ±ÄŸÄ± Ã¼retmemelidir; owner yine `history` iÃ§inde kalÄ±r
- renumber veya reorder iÅŸlemleri aynÄ± request iÃ§inde tekrarlandÄ±ÄŸÄ±nda stabil sÄ±ra Ã¼retmelidir
- media processing retry'larÄ± yeni chapter kaydÄ± yaratmamalÄ±dÄ±r

## State YapÄ±sÄ±
- publish durumu
- preview veya early access state'i
- page yapÄ±sÄ± ve media referanslarÄ±
- resume anchor format sÃ¼rÃ¼mÃ¼

## Test NotlarÄ±
- CRUD, detail ve read akÄ±ÅŸÄ± testleri
- navigation ve latest chapter sÄ±ralama testleri
- preview, early access ve media doÄŸrulamalarÄ±
- `chapter -> history` kontrat ve resume anchor doÄŸrulamalarÄ±


---

# Comment ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `comment`

## AmaÃ§
`comment` modÃ¼lÃ¼nÃ¼n amacÄ±, hedefe baÄŸlÄ± yorum verisini, thread yapÄ±sÄ±nÄ± ve yorum yaÅŸam dÃ¶ngÃ¼sÃ¼nÃ¼ yÃ¶netmektir.

## Sorumluluk AlanÄ±
- comment create, edit, delete ve listing akÄ±ÅŸlarÄ±
- root veya reply thread yapÄ±sÄ±
- edit window, soft delete gÃ¶rÃ¼nÃ¼mÃ¼ ve delete behavior kurallarÄ±
- spoiler, pin, lock ve moderation alanlarÄ±
- nested reply depth limiti, anti-spam ve rate limit kurallarÄ±
- author shadowban veya moderation etkilerinin yorum gÃ¶rÃ¼nÃ¼mÃ¼ne yansÄ±masÄ±

## Bu ModÃ¼l Neyi Yapmaz?
- social duvar post veya social-native reply owner'lÄ±ÄŸÄ± taÅŸÄ±maz
- eriÅŸim veya yetki kararÄ±nÄ± tek baÅŸÄ±na vermez
- support veya moderation case kayÄ±tlarÄ±nÄ± kendi iÃ§ine dÃ¶nÃ¼ÅŸtÃ¼rmez

## Veri SahipliÄŸi
- yorum iÃ§eriÄŸi
- `target_type` ve `target_id` iliÅŸkisi
- parent-child reply yapÄ±sÄ±
- moderation, spoiler, pin ve lock alanlarÄ±
- edit window ve delete gÃ¶rÃ¼nÃ¼mÃ¼ne ait metadata

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- hedef iÃ§erik varlÄ±ÄŸÄ±nÄ±n owner verisi
- social messaging veya wall reply native kayÄ±tlarÄ±
- moderation case, support ticket veya access policy owner'lÄ±ÄŸÄ±

## Hedef Tipi SÃ¶zlÃ¼ÄŸÃ¼
- `target_type` deÄŸerleri canonical olarak `docs/shared/target-types.md` dosyasÄ±ndaki kayÄ±tlarla hizalÄ± olmalÄ±dÄ±r.
- `comment` hedefleri bu aÅŸamada iÃ§erik odaklÄ± hedeflerdir; sosyal duvar post'u veya sosyal duvar reply zinciri `social` modÃ¼lÃ¼nde native kalmalÄ± ve Ã¶rtÃ¼k olarak comment target'Ä±na dÃ¶nÃ¼ÅŸtÃ¼rÃ¼lmemelidir.
- Yeni yorum hedefi eklendiÄŸinde aynÄ± deÄŸiÅŸiklik setinde hedef modÃ¼l dokÃ¼manÄ±, `comment` modÃ¼l dokÃ¼manÄ± ve canonical target type kaydÄ± gÃ¼ncellenmelidir.

## Access KontratÄ±
`comment` yetki kararÄ± vermez. Yorum oluÅŸturma, dÃ¼zenleme, silme ve gÃ¶rÃ¼nÃ¼rlÃ¼k kararlarÄ± `access` ile korunur; gÃ¼nlÃ¼k moderasyon akÄ±ÅŸlarÄ± `moderation`, yÃ¼ksek riskli veya yÃ¶netimsel handoff akÄ±ÅŸlarÄ± `admin` tarafÄ±nda yÃ¼rÃ¼r. Site geneli, manga detayÄ± veya chapter altÄ± yorum alanlarÄ±nÄ±n aÃ§Ä±lÄ±p kapatÄ±lmasÄ± ve yorum gÃ¶nderme aralÄ±ÄŸÄ± gibi runtime ayarlar admin Ã¼zerinden yÃ¶netilebilir.

## API veya Event SÄ±nÄ±rÄ±
- yorum listeleme ve thread yÃ¼zeyi
- yorum detay ve hedef iliÅŸkisi yÃ¼zeyi
- moderation, support veya admin odaklÄ± yorum iÅŸlemleri iÃ§in kontrollÃ¼ yÃ¼zey
- yorum modÃ¼lÃ¼ read ve write yÃ¼zeyleri gerektiÄŸinde birbirinden baÄŸÄ±msÄ±z olarak kontrol edilebilir olmalÄ±dÄ±r
- hedef tipi geniÅŸlemesi canonical target type sÃ¶zlÃ¼ÄŸÃ¼ ile uyumlu tutulmalÄ±dÄ±r

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `manga`, `chapter`, `user`, `access`, `admin`, `support` ve `moderation` entegrasyonlarÄ±
- profanity filter veya anti-spam yardÄ±mcÄ±larÄ± iÃ§in ortak altyapÄ± ihtiyacÄ±

## Settings Etkileri
- `feature.comment.read.enabled`
- `feature.comment.write.enabled`
- `comment.write.cooldown_seconds`
- edit window veya spoiler policy gibi ek yÃ¼zeyler eklenirse settings envanteri geniÅŸletilmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `comment.created`, `comment.edited`, `comment.deleted`, `comment.moderated`
- tÃ¼ketir: moderation hide veya restore kararÄ±, admin lock veya unlock kararÄ±
- comment count ve engagement projection'larÄ± `docs/shared/projection-strategy.md` ile hizalÄ± beslenmelidir

## Audit ve Ä°zleme
- shadowban etkisi, bulk lock veya pin deÄŸiÅŸimi ve admin/moderation kaynaklÄ± gÃ¶rÃ¼nÃ¼rlÃ¼k deÄŸiÅŸimleri auditlenmelidir
- spam ve flood koruma tetiklenmeleri metrik olarak izlenmelidir

## Ä°dempotency ve Retry
- duplicate create istekleri aynÄ± request scope iÃ§inde ikinci yorum kaydÄ± oluÅŸturmamalÄ±dÄ±r
- delete veya moderation hide aksiyonu tekrarlandÄ±ÄŸÄ±nda aynÄ± final gÃ¶rÃ¼nÃ¼rlÃ¼k state'ini korumalÄ±dÄ±r
- pagination veya thread rebuild retry'larÄ± yeni yan etki Ã¼retmemelidir

## State YapÄ±sÄ±
- moderation_status
- spoiler_flag
- pin ve lock durumu
- visibility etkileyen alanlar
- edit window veya soft delete gÃ¶rÃ¼nÃ¼mÃ¼ durumu

## Test NotlarÄ±
- create, edit, delete ve list testleri
- thread, reply depth ve soft delete gÃ¶rÃ¼nÃ¼mÃ¼ testleri
- target relation ve spoiler propagation testleri
- rate limit, anti-spam ve shadowban doÄŸrulamalarÄ±
- visibility ve access entegrasyonu doÄŸrulamalarÄ±


---

# History ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `history`

## AmaÃ§
`history` modÃ¼lÃ¼nÃ¼n amacÄ±, kullanÄ±cÄ±ya ait continue reading, reading history, bookmark-library ve okuma devamlÄ±lÄ±ÄŸÄ± kayÄ±tlarÄ±nÄ± tek bir yaÅŸam dÃ¶ngÃ¼sÃ¼ altÄ±nda yÃ¶netmektir.

## Sorumluluk AlanÄ±
- continue reading ve resume akÄ±ÅŸlarÄ±
- reading history timeline ve own history yÃ¼zeyleri
- user-manga library, bookmark veya favorite kayÄ±tlarÄ±
- resume checkpoint formatÄ±, last read conflict Ã§Ã¶zÃ¼mÃ¼ ve multi-device merge policy
- timeline retention, cleanup ve compact history write stratejisi
- admin tarafÄ±ndan yÃ¶netilen continue reading, library, timeline veya bookmark write runtime kontrolleri ile uyumlu Ã§alÄ±ÅŸma

## Bu ModÃ¼l Neyi Yapmaz?
- manga veya chapter iÃ§eriÄŸinin canonical owner'lÄ±ÄŸÄ±na dÃ¶nÃ¼ÅŸmez
- global visibility default'unu tek baÅŸÄ±na belirlemez
- access policy, social iliÅŸki veya notification preference owner'lÄ±ÄŸÄ± taÅŸÄ±maz

## Veri SahipliÄŸi
- user-manga library entry kayÄ±tlarÄ±
- bookmark, favorite veya reading status alanlarÄ±
- user-chapter son okuma chapter veya page referanslarÄ±
- progress snapshot, checkpoint ve history timeline verileri
- entry-level visibility veya sharing iÃ§in gereken history tarafÄ± metadata alanlarÄ±
- global visibility default'u deÄŸil, library entry veya history event bazlÄ± share metadata alanlarÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- manga metadata, chapter page yapÄ±sÄ± veya kullanÄ±cÄ± profil verisi
- global user visibility preference sinyali
- access kararÄ±nÄ±n kendisi veya recommendation owner algoritmasÄ±

## Access KontratÄ±
`history` yetki kararÄ± vermez. KullanÄ±cÄ±nÄ±n kendi history veya library yÃ¼zeylerini gÃ¶rmesi `access` ile korunur. `access` tarafÄ±nda en az `history.continue_reading.read.own`, `history.timeline.read.own`, `history.library.read.own`, `history.bookmark.write.own` ve gerektiÄŸinde `history.library.read.public` gibi canonical permission Ã¶rnekleri tanÄ±mlanmalÄ±dÄ±r. Public library paylaÅŸÄ±mÄ± veya reading activity gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ gerektiÄŸinde global gÃ¶rÃ¼nÃ¼rlÃ¼k default'larÄ± `user` modÃ¼lÃ¼ndeki sinyallerden, entry-level paylaÅŸÄ±m metadata'sÄ± ise `history` iÃ§indeki kayÄ±tlarÄ±ndan yorumlanmalÄ±dÄ±r. Final gÃ¶rÃ¼nÃ¼rlÃ¼kte `user` modÃ¼lÃ¼ndeki global deny veya daha dar default Ã¼st sÄ±nÄ±rdÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- own continue reading, own library ve own history timeline yÃ¼zeyi
- eriÅŸim uygunsa kontrollÃ¼ public library read yÃ¼zeyi
- `chapter` modÃ¼lÃ¼nden read start, checkpoint, finish veya resume anchor intake kontratÄ±
- `manga`, `mission` ve discovery tÃ¼keticileri iÃ§in kontrollÃ¼ okuma Ã¶zeti veya signal surface
- admin tarafÄ±ndan history yÃ¼zeylerini daraltmak iÃ§in kontrollÃ¼ operasyon contract yÃ¼zeyi

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `user`, `manga`, `chapter`, `access`, `admin`, `mission` ve discovery consumer'larÄ± ile entegrasyon
- projection updater ve compact history writer gibi yardÄ±mcÄ± yapÄ±lar

## Settings Etkileri
- `feature.history.continue_reading.enabled`
- `feature.history.library.enabled`
- `feature.history.timeline.enabled`
- `feature.history.bookmark_write.enabled`
- retention veya cleanup metric'leri ayrÄ± key olarak eklenirse settings envanterine yazÄ±lmalÄ±dÄ±r

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `history.checkpoint.updated`, `history.chapter.finished`, `history.library.changed`
- tÃ¼ketir: `chapter.read.started`, `chapter.read.finished`, user visibility sinyalleri
- continue reading projection ve discovery summary'leri `docs/shared/projection-strategy.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- public share metadata deÄŸiÅŸimi, admin kaynaklÄ± visibility daraltmasÄ± ve bulk cleanup aksiyonlarÄ± auditlenmelidir
- merge conflict oranÄ±, checkpoint dedup sayÄ±sÄ± ve projection lag izlenmelidir

## Ä°dempotency ve Retry
- checkpoint write, merge ve multi-device sync akÄ±ÅŸlarÄ± `docs/shared/idempotency-policy.md` ile hizalÄ± idempotent olmalÄ±dÄ±r
- aynÄ± checkpoint tekrar geldiÄŸinde yeni timeline gÃ¼rÃ¼ltÃ¼sÃ¼ Ã¼retmemelidir
- cleanup job veya rebuild retry'larÄ± veri kaybÄ± Ã¼retmeden tekrar Ã§alÄ±ÅŸabilmelidir

## State YapÄ±sÄ±
- in_progress, completed veya dropped reading status
- bookmarked veya favorited library durumu
- last_read_at, resume snapshot ve merge conflict alanlarÄ±
- timeline retention veya history pause durumu

## Test NotlarÄ±
- continue reading Ã§Ã¶zÃ¼mleme ve checkpoint idempotency testleri
- bookmark ve favorite ayrÄ±mÄ± testleri
- duplicate history yazÄ±mÄ±, merge conflict ve cleanup pencere doÄŸrulamalarÄ±
- own history, public library visibility ve bookmark write permission testleri
- `chapter -> history` ve downstream signal entegrasyonu doÄŸrulamalarÄ±


---

# Inventory ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `inventory`

## AmaÃ§
`inventory` modÃ¼lÃ¼nÃ¼n amacÄ±, kullanÄ±cÄ±nÄ±n sahip olduÄŸu item, kozmetik, Ã¶dÃ¼l ve benzeri varlÄ±klarÄ± tek bir envanter yaÅŸam dÃ¶ngÃ¼sÃ¼ altÄ±nda yÃ¶netmek ve Ã¶dÃ¼l teslim yÃ¼rÃ¼tÃ¼mÃ¼nÃ¼ merkezi biÃ§imde sahiplenmektir.

## Sorumluluk AlanÄ±
- ownable item definition ve item type ayrÄ±mlarÄ±
- kullanÄ±cÄ± envanter kayÄ±tlarÄ± ve sahiplik durumlarÄ±
- grant, reward teslim yÃ¼rÃ¼tÃ¼mÃ¼, revoke, consume ve equip akÄ±ÅŸlarÄ±
- stackable veya unique item kurallarÄ±
- expiry mantÄ±ÄŸÄ±, equip slot kurallarÄ± ve selected cosmetic referanslarÄ±
- reward source canonical enum ve idempotent grant kontrolÃ¼

## Bu ModÃ¼l Neyi Yapmaz?
- sellable shop product veya offer catalog owner'lÄ±ÄŸÄ± taÅŸÄ±maz
- Ã¶deme, wallet veya ledger doÄŸruluÄŸu Ã¼retmez
- mission veya royalpass claim uygunluÄŸu kararÄ±nÄ± tek baÅŸÄ±na vermez

## Veri SahipliÄŸi
- ownable item definition alanlarÄ±; sellable product veya offer catalog kayÄ±tlarÄ± deÄŸil
- user inventory entry kayÄ±tlarÄ±
- stack count, ownership state ve expiry bilgileri
- equip slot veya selected cosmetic referanslarÄ±
- grant, consume ve revoke iÅŸlem loglarÄ±
- reward source ve source reference metadata alanlarÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- shop Ã¼rÃ¼n kataloÄŸu ve fiyat planÄ± verisi
- payment checkout, callback veya ledger kayÄ±tlarÄ±
- mission veya royalpass progression owner verisi

## Access KontratÄ±
`inventory` yetki kararÄ± vermez. KullanÄ±cÄ±nÄ±n kendi envanterini gÃ¶rmesi, reward teslimini tamamlamasÄ±, consume etmesi veya equip etmesi `access` ile korunur. `mission`, `royalpass` veya diÄŸer producer modÃ¼ller claim uygunluÄŸu veya reward kaynaÄŸÄ± Ã¼retebilir; ancak final grant yÃ¼rÃ¼tÃ¼mÃ¼ ve item sahipliÄŸi `inventory` modÃ¼lÃ¼nde kalÄ±r. Kaynak tipleri `docs/shared/reward-source-types.md` ile hizalÄ± olmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- own inventory list, own item detail, own equip veya consume yÃ¼zeyi
- admin grant veya revoke orkestrasyonu iÃ§in kontrollÃ¼ yÃ¶netim yÃ¼zeyi
- reward producer modÃ¼ller iÃ§in aÃ§Ä±k grant execution veya reservation contract yÃ¼zeyi
- `shop` iÃ§in sellable product -> grantable item definition Ã§Ã¶zÃ¼mleme kontratÄ±

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `user`, `access`, `admin`, `shop`, `notification`, `mission` ve `royalpass` entegrasyonlarÄ±
- snapshot veya serializer yardÄ±mcÄ±larÄ± iÃ§in future-ready altyapÄ±

## Settings Etkileri
- `feature.inventory.read.enabled`
- `feature.inventory.claim.enabled`
- `feature.inventory.equip.enabled`
- `feature.inventory.consume.enabled`
- category bazlÄ± inventory gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ eklenirse settings envanteri geniÅŸletilmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `inventory.granted`, `inventory.revoked`, `inventory.consumed`, `inventory.equipped`
- tÃ¼ketir: mission claim, royalpass claim, shop fulfillment, admin manual grant sinyalleri
- reward projection ve grant retry davranÄ±ÅŸÄ± `docs/shared/outbox-pattern.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- manual grant, revoke, consume, equip slot deÄŸiÅŸimi ve correction aksiyonlarÄ± auditlenmelidir
- duplicate grant korumasÄ±, expired item cleanup ve source mismatch oranÄ± izlenmelidir

## Ä°dempotency ve Retry
- grant akÄ±ÅŸlarÄ± source reference + request id ile `docs/shared/idempotency-policy.md` kapsamÄ±nda korunmalÄ±dÄ±r
- consume ve equip tekrarlandÄ±ÄŸÄ±nda ikinci yan etki Ã¼retmemeli; mevcut final state'e baÄŸlanmalÄ±dÄ±r
- recovery veya reconcile akÄ±ÅŸlarÄ± yeni sahiplik Ã§oÄŸaltmamalÄ±dÄ±r

## State YapÄ±sÄ±
- active, consumed, expired veya revoked item durumu
- stack balance veya quantity durumu
- equipped veya unequipped durumu
- claim veya grant pause durumu
- source type ve source reference alanlarÄ±

## Test NotlarÄ±
- grant, claim, consume ve equip testleri
- item type, stackable veya unique ve expiry doÄŸrulamalarÄ±
- source idempotency ve duplicate grant testleri
- reward producer entegrasyonu ve access doÄŸrulamalarÄ±
- equip slot ve selected cosmetic referans testleri


---

# Manga ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `manga`

## AmaÃ§
`manga` modÃ¼lÃ¼nÃ¼n amacÄ±, ana iÃ§erik varlÄ±ÄŸÄ±nÄ±, iÃ§erik metadata yapÄ±sÄ±nÄ± ve public iÃ§erik akÄ±ÅŸÄ±nÄ±n veri temelini yÃ¶netmektir.

## Sorumluluk AlanÄ±
- manga CRUD, listing ve detail akÄ±ÅŸlarÄ±
- arama, filtreleme ve sÄ±ralama
- metadata, taxonomy ve slug iliÅŸkileri
- publish, archive, moderation ve gÃ¶rÃ¼nÃ¼rlÃ¼k verileri
- featured, recommendation, editoryal koleksiyon ve keÅŸif yÃ¼zeyleri
- content versioning ve search index stratejisi iÃ§in referans alanlarÄ±

## Bu ModÃ¼l Neyi Yapmaz?
- chapter page iÃ§eriÄŸini veya kullanÄ±cÄ±ya ait history kaydÄ±nÄ± owner olarak taÅŸÄ±maz
- yorum, support veya moderation vaka verisini kendi iÃ§ine kopyalamaz
- recommendation verisini tek baÅŸÄ±na kara kutu algoritma olarak sahiplenmez; kaynak sinyal ve projection stratejisi dokÃ¼mante edilmelidir

## Veri SahipliÄŸi
- title, slug, summary ve gÃ¶rsel alanlarÄ±
- taxonomy iliÅŸkileri ve iÃ§erik metadata alanlarÄ±
- publish ve moderation state alanlarÄ±
- view_count, comment_count, chapter_count ve benzeri denormalize sayaÃ§ alanlarÄ±
- editoryal collection, discovery placement ve recommendation metadata alanlarÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- chapter page veya page media owner verisi
- kullanÄ±cÄ±ya ait continue reading, favorite veya bookmark kayÄ±tlarÄ±
- search provider altyapÄ±sÄ±nÄ±n teknik config secret'larÄ±

## Access KontratÄ±
`manga` eriÅŸim kararÄ± vermez. Public gÃ¶rÃ¼nÃ¼rlÃ¼k ve yÃ¶netimsel eriÅŸim kararlarÄ± `access` tarafÄ±ndan yorumlanÄ±r. Admin tarafÄ±ndan yÃ¶netilen runtime ayarlar manga listeleme, detay, recommendation veya editoryal gÃ¶rÃ¼nÃ¼rlÃ¼k yÃ¼zeylerini daraltabilir; veri sahipliÄŸi yine `manga` modÃ¼lÃ¼nde kalÄ±r. Visibility kavramÄ± `docs/shared/visibility-states.md` ile hizalÄ± tutulmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- public listing ve public detail yÃ¼zeyi
- public recommendation, collection veya editoryal discovery yÃ¼zeyi
- yÃ¶netimsel CRUD ve metadata yÃ¼zeyi
- chapter varsayÄ±lan eriÅŸim verisi gerekiyorsa kontrollÃ¼ contract yÃ¼zeyi
- sayaÃ§ gÃ¼ncelleme veya projection ihtiyacÄ± varsa aÃ§Ä±k event veya counter contract yÃ¼zeyi
- `comment` ve `support` iÃ§in raporlanabilir hedef veya target relation yÃ¼zeyi

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `access`, `admin`, `history`, `comment`, `chapter` ve discovery consumer'larÄ± ile kontrollÃ¼ entegrasyon
- arama ve index lifecycle iÃ§in `docs/shared/search-strategy.md`

## Settings Etkileri
- `feature.manga.list.enabled`
- `feature.manga.detail.enabled`
- `feature.manga.discovery.enabled`
- arama veya editoryal surface iÃ§in ayrÄ± key aÃ§Ä±lÄ±rsa settings envanteri aynÄ± deÄŸiÅŸiklikte gÃ¼ncellenmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `manga.published`, `manga.archived`, `manga.discovery.updated`
- tÃ¼ketir: `comment.*`, `chapter.*`, `history.*` kaynaklÄ± sayaÃ§ veya engagement sinyalleri
- search index ve recommendation projection'larÄ± `docs/shared/projection-strategy.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- publish, archive, visibility ve editoryal collection deÄŸiÅŸiklikleri auditlenmelidir
- projection lag, search index gecikmesi ve sayaÃ§ drift farklarÄ± izlenebilir olmalÄ±dÄ±r

## Ä°dempotency ve Retry
- publish veya archive transition'larÄ± tekrar Ã§alÄ±ÅŸtÄ±rÄ±ldÄ±ÄŸÄ±nda Ã§eliÅŸkili state Ã¼retmemelidir
- counter reconcile ve search reindex akÄ±ÅŸlarÄ± replay-safe olmalÄ±dÄ±r
- discovery placement gÃ¼ncellemeleri request id veya version alanÄ± ile duplicate write Ã¼retmemelidir

## State YapÄ±sÄ±
- publish veya archive durumu
- visibility veya public gÃ¶rÃ¼nÃ¼rlÃ¼k etkileyen alanlar
- collection veya discovery gÃ¶rÃ¼nÃ¼rlÃ¼k state alanlarÄ±
- taxonomy ve metadata sÃ¼rÃ¼m bilgileri

## Test NotlarÄ±
- CRUD, listing ve detail testleri
- search, filter, sort ve taxonomy doÄŸrulamalarÄ±
- sayaÃ§ reconcile ve projection doÄŸrulamalarÄ±
- public gÃ¶rÃ¼nÃ¼rlÃ¼k ve access entegrasyonu doÄŸrulamalarÄ±


---

# Mission ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `mission`

## AmaÃ§
`mission` modÃ¼lÃ¼nÃ¼n amacÄ±, gÃ¼nlÃ¼k, haftalÄ±k, aylÄ±k, event ve seviye bazlÄ± gÃ¶rev tanÄ±mlarÄ±nÄ±, kullanÄ±cÄ± ilerlemesini ve gÃ¶rev Ã¶dÃ¼lÃ¼ iÃ§in claim uygunluÄŸu veya claim request akÄ±ÅŸlarÄ±nÄ± yÃ¶netmektir.

## Sorumluluk AlanÄ±
- mission definition ve mission category yÃ¶netimi
- daily, weekly, monthly, event ve level-based mission yapÄ±larÄ±
- mission trigger source listesi ve progress accumulation kurallarÄ±
- kullanÄ±cÄ± progress, completion, claim uygunluÄŸu ve claim request yaÅŸam dÃ¶ngÃ¼sÃ¼
- reset window, streak veya dÃ¶nemsel yenileme kurallarÄ±
- event mission ile season mission iliÅŸkisinin dokÃ¼mantasyonu

## Bu ModÃ¼l Neyi Yapmaz?
- final reward sahipliÄŸini kendi iÃ§inde yazmaz
- global EXP, level veya kullanÄ±cÄ± profil progression owner'lÄ±ÄŸÄ±na dÃ¶nÃ¼ÅŸmez
- Ã¶deme, inventory veya royalpass owner verisini kopyalamaz

## Veri SahipliÄŸi
- mission tanÄ±mlarÄ± ve category alanlarÄ±
- user mission progress kayÄ±tlarÄ±
- completion, claim eligibility ve reset durum verileri
- event mission penceresi ve schedule alanlarÄ±
- reward reference ve claim source bilgileri

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- final inventory item kaydÄ± veya grant sonucu
- payment veya shop purchase owner verisi
- royalpass season tier owner verisi

## Access KontratÄ±
`mission` yetki kararÄ± vermez. GÃ¶rev gÃ¶rÃ¼ntÃ¼leme, claim ve yÃ¶netim yÃ¼zeyleri `access` ile korunur. Global EXP veya level sahipliÄŸi `user` modÃ¼lÃ¼nde kalabilir; `mission` bu verileri gÃ¶rev tamamlanma sinyali olarak tÃ¼ketir. Claim uygunluÄŸu ile final grant ayrÄ±mÄ± transaction boundary olarak aÃ§Ä±k tutulmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- own mission list, own progress ve own claim-request yÃ¼zeyi
- admin iÃ§in mission definition, reset ve period kontrol yÃ¼zeyi
- producer modÃ¼llerden alÄ±nan progress event contract'larÄ± ve `inventory` veya `notification` iÃ§in reward event yÃ¼zeyi
- mission read, claim ve progress ingest yÃ¼zeyleri gerektiÄŸinde birbirinden baÄŸÄ±msÄ±z runtime anahtarlarÄ± ile yÃ¶netilebilir olmalÄ±dÄ±r

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `user`, `access`, `admin`, `inventory`, `notification`, `royalpass` ve progress Ã¼reticisi modÃ¼ller
- schedule evaluator, progress aggregator ve rule matcher yardÄ±mcÄ±larÄ±

## Settings Etkileri
- `mission.daily.reset_hour_utc`
- `feature.mission.read.enabled`
- `feature.mission.claim.enabled`
- `feature.mission.progress_ingest.enabled`
- sezonluk veya event bazlÄ± ek reset pencereleri aÃ§Ä±lÄ±rsa settings envanterinde gÃ¶rÃ¼nÃ¼r tutulmalÄ±dÄ±r

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `mission.progressed`, `mission.completed`, `mission.claim.requested`, `mission.reset`
- tÃ¼ketir: okuma, yorum, sosyal etkileÅŸim ve diÄŸer producer event'leri
- progress projection ve season handoff'larÄ± `docs/shared/projection-strategy.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- manual completion, reset override, claim rejection ve admin category mÃ¼dahalesi auditlenmelidir
- stuck progress, duplicate claim ve reset drift oranlarÄ± izlenebilir olmalÄ±dÄ±r

## Ä°dempotency ve Retry
- progress ingest ve claim request akÄ±ÅŸlarÄ± `docs/shared/idempotency-policy.md` ile hizalÄ± olmalÄ±dÄ±r
- aynÄ± gÃ¶rev ve dÃ¶nem iÃ§in tekrar claim isteÄŸi yeni grant zinciri baÅŸlatmamalÄ±dÄ±r
- reset job tekrarlandÄ±ÄŸÄ±nda aynÄ± dÃ¶nemi iki kez kapatmamalÄ±dÄ±r

## State YapÄ±sÄ±
- active, completed, claimed veya expired mission durumu
- reset_pending veya recurring window durumu
- progress snapshot ve objective state alanlarÄ±
- mission category veya claim surface kapanma durumu
- streak veya period metadata alanlarÄ±

## Test NotlarÄ±
- progress ve completion testleri
- trigger source ve accumulation kurallarÄ± testleri
- reset, recurring window ve streak testleri
- claim ve reward grant entegrasyonu testleri
- producer event, idempotency ve season iliÅŸkisi doÄŸrulamalarÄ±


---

# Moderation ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `moderation`

## AmaÃ§
`moderation` modÃ¼lÃ¼nÃ¼n amacÄ±, role bazlÄ± veya kullanÄ±cÄ± bazlÄ± scoped moderatÃ¶rlerin gÃ¼nlÃ¼k vaka yÃ¶netimi, inceleme ve sÄ±nÄ±rlÄ± mÃ¼dahale iÅŸ akÄ±ÅŸlarÄ±nÄ± `admin` modÃ¼lÃ¼nden ayrÄ±ÅŸtÄ±rÄ±lmÄ±ÅŸ Ã¶zel bir panel Ã¼zerinden yÃ¼rÃ¼tmektir.

## Sorumluluk AlanÄ±
- moderatÃ¶r paneli ve scoped queue yÃ¼zeyleri
- role veya kullanÄ±cÄ± bazlÄ± moderatÃ¶r scope matrisi
- vaka inceleme, assignment, case lifecycle ve stale policy akÄ±ÅŸlarÄ±
- yorum, chapter ve manga yÃ¼zeyleri iÃ§in scoped moderasyon iÅŸ akÄ±ÅŸlarÄ±
- moderatÃ¶r notu, karar Ã¶zeti, evidence snapshot ve escalation akÄ±ÅŸlarÄ±
- support kaynaklÄ± report intake'lerden gerekli olduÄŸunda linked moderation case aÃ§ma entegrasyonu

## Bu ModÃ¼l Neyi Yapmaz?
- global authorization, role veya permission owner'lÄ±ÄŸÄ± Ã¼retmez
- merkezi settings ve kill switch owner'lÄ±ÄŸÄ±na dÃ¶nÃ¼ÅŸmez
- admin override'Ä±n Ã¼stÃ¼ne Ã§Ä±kan final karar owner'lÄ±ÄŸÄ± taÅŸÄ±maz

## Veri SahipliÄŸi
- moderation case veya queue kayÄ±tlarÄ±
- moderatÃ¶r assignment bilgileri
- moderatÃ¶r notlarÄ±, karar Ã¶zeti ve vaka timeline verisi
- escalation, stale case ve handoff durum alanlarÄ±
- evidence snapshot referanslarÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- target iÃ§eriÄŸin canonical owner verisi
- access policy, admin settings veya support intake owner verisi
- payment, user veya inventory kaynaklÄ± business kayÄ±tlar

## Hedef Tipi SÃ¶zlÃ¼ÄŸÃ¼
- `moderation` case hedefleri canonical olarak `docs/shared/target-types.md` dosyasÄ±ndaki kayÄ±tlarla hizalÄ± olmalÄ±dÄ±r.
- Alt yÃ¼zey, ekran veya aksiyon bilgisi `target_type` iÃ§ine gÃ¶mÃ¼lmemeli; context verisi ile taÅŸÄ±nmalÄ±dÄ±r.
- Yeni moderation hedefi eklendiÄŸinde aynÄ± deÄŸiÅŸiklik setinde hedef modÃ¼l dokÃ¼manÄ±, `moderation` modÃ¼l dokÃ¼manÄ± ve canonical target type kaydÄ± gÃ¼ncellenmelidir.

## Access KontratÄ±
`moderation` yetki kararÄ± vermez. Comment moderatÃ¶rÃ¼, chapter moderatÃ¶rÃ¼ veya manga moderatÃ¶rÃ¼ gibi role bazlÄ± ya da kullanÄ±cÄ± bazlÄ± scope kararlarÄ± `access` tarafÄ±ndan yorumlanÄ±r. Admin kullanÄ±cÄ±larÄ± aynÄ± vaka verisine tam yÃ¶netim yetkisi ile eriÅŸebilir; merkezi settings, kill switch ve sistem operasyon yÃ¼zeyleri ise `admin` modÃ¼lÃ¼nde kalÄ±r. GÃ¼nlÃ¼k case sahipliÄŸi `moderation` iÃ§inde kalsa da admin override, reopen, reassignment veya freeze kararÄ± oluÅŸtuÄŸunda bu karar moderatÃ¶r aksiyonunun Ã¼zerinde precedence taÅŸÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- moderatÃ¶r queue, case detail ve sÄ±nÄ±rlÄ± karar yÃ¼rÃ¼tme yÃ¼zeyi
- admin ile orkestrasyon gerektiren escalation veya yÃ¼ksek riskli handoff yÃ¼zeyi
- `support` report kaydÄ± ile linked moderation case iliÅŸki yÃ¼zeyi; support kaydÄ± ve moderation case aynÄ± kayÄ±t haline getirilmemelidir
- moderation action olaylarÄ± gerektiÄŸinde `admin`, `notification` veya hedef modÃ¼llere kontrollÃ¼ event yÃ¼zeyi ile aktarÄ±labilir

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `access`, `admin`, `support`, `comment`, `manga`, `chapter` ve `notification` entegrasyonlarÄ±
- rule-based moderation helper veya content scoring altyapÄ±sÄ± iÃ§in geleceÄŸe dÃ¶nÃ¼k ihtiyaÃ§

## Settings Etkileri
- `feature.moderation.panel.enabled`
- `feature.moderation.queue.enabled`
- `feature.moderation.action.enabled`
- scope bazlÄ± queue veya action yÃ¼zeyleri ayrÄ±ÅŸtÄ±ÄŸÄ±nda settings envanteri geniÅŸletilmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `moderation.case.created`, `moderation.case.assigned`, `moderation.case.escalated`, `moderation.action.applied`
- tÃ¼ketir: `support.moderation_handoff_requested`, admin override veya reopen sinyalleri
- queue summary ve action projection'larÄ± `docs/shared/projection-strategy.md` ile hizalÄ± planlanmalÄ±dÄ±r

## Audit ve Ä°zleme
- case assignment, escalation, evidence snapshot eriÅŸimi ve action sonucu auditlenmelidir
- stale case birikimi, queue gecikmesi ve action override oranÄ± izlenebilir olmalÄ±dÄ±r

## Ä°dempotency ve Retry
- linked case create akÄ±ÅŸÄ± aynÄ± support referansÄ± iÃ§in duplicate moderation case Ã¼retmemelidir
- aynÄ± aksiyon tekrarlandÄ±ÄŸÄ±nda ikinci kez iÃ§erik state'i bozulmamalÄ±dÄ±r
- escalation retry'larÄ± mevcut case referansÄ± Ã¼zerinden gÃ¼venli biÃ§imde yeniden denenmelidir

## State YapÄ±sÄ±
- `docs/shared/moderation-statuses.md` ile hizalÄ± `case_status`
- `assignment_status`
- `escalation_status`
- `action_result` ve review lifecycle alanlarÄ±
- queue veya action surface kapanma durumlarÄ±

## Test NotlarÄ±
- scoped queue gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ testleri
- case transition, assignment ve stale policy testleri
- support report'tan linked case aÃ§Ä±lÄ±ÅŸÄ± doÄŸrulamalarÄ±
- action audit ve admin handoff entegrasyonu testleri
- scope matrix ve precedence doÄŸrulamalarÄ±


---

# Notification ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `notification`

## AmaÃ§
`notification` modÃ¼lÃ¼nÃ¼n amacÄ±, sistem genelindeki olaylardan bildirim Ã¼retmek, kullanÄ±cÄ±ya uygun kanaldan teslim etmek ve bildirim tercihlerini merkezi ÅŸekilde yÃ¶netmektir.

## Sorumluluk AlanÄ±
- in-app bildirim kutusu ve read/unread akÄ±ÅŸlarÄ±
- bildirim category, template ve channel yÃ¶netimi
- in-app, email ve gelecekteki push delivery state ayrÄ±mÄ±
- digest, retry, backoff ve dedup stratejileri
- kullanÄ±cÄ±ya ait detaylÄ± bildirim tercihleri, mute ve quiet-hour benzeri yÃ¼zeyler

## Bu ModÃ¼l Neyi Yapmaz?
- business event owner'lÄ±ÄŸÄ± veya authorization kararÄ± Ã¼retmez
- kullanÄ±cÄ± profil verisini veya sosyal iliÅŸki verisini sahiplenmez
- producer modÃ¼l yerine onun iÅŸ kuralÄ±nÄ± kime bildirim gideceÄŸi konusunda yeniden tanÄ±mlamaz

## Veri SahipliÄŸi
- notification kaydÄ± ve delivery attempt verileri
- notification template, category ve channel tanÄ±mlarÄ±
- kullanÄ±cÄ± bildirim tercihleri, category mute ve quiet-hour alanlarÄ±
- suppression, dedup ve digest batch metadata'sÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- producer modÃ¼llerin kaynak event verisi
- kullanÄ±cÄ± hesabÄ±, sosyal blok listesi veya access policy owner'lÄ±ÄŸÄ±
- dÄ±ÅŸ provider secret veya credential kayÄ±tlarÄ±

## Access KontratÄ±
`notification` yetki kararÄ± vermez. Kendi bildirimlerini gÃ¶rme veya yÃ¶netme kararÄ± `access` ile korunur. Hangi olayÄ±n kime bildirim Ã¼reteceÄŸi iÅŸ kuralÄ± olarak `notification` iÃ§inde yorumlanabilir; ancak bu karar authorization yerine teslim kuralÄ± niteliÄŸindedir. Kategori adlarÄ± `docs/shared/notification-categories.md` ile hizalÄ± olmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- own inbox, own preference ve own notification detail yÃ¼zeyi
- diÄŸer modÃ¼llerden alÄ±nan producer event veya notification contract yÃ¼zeyi
- admin tarafÄ±ndan category, channel, digest veya delivery pause yÃ¶netimi iÃ§in kontrollÃ¼ operasyon yÃ¼zeyi

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `auth`, `user`, `access`, `admin`, `social`, `mission`, `royalpass`, `support`, `moderation`, `shop`, `payment` ve diÄŸer producer modÃ¼ller
- queue, retry ve delivery katmanÄ± iÃ§in `docs/shared/cache-queue-strategy.md`

## Settings Etkileri
- `feature.notification.inbox.enabled`
- `feature.notification.preference.enabled`
- `feature.notification.digest.enabled`
- `notification.delivery.paused`
- category veya channel bazlÄ± selector geniÅŸlemeleri settings envanterinde aÃ§Ä±kÃ§a gÃ¶sterilmelidir

## Event AkÄ±ÅŸlarÄ±
- tÃ¼ketir: producer modÃ¼llerden gelen domain event'leri
- Ã¼retir: `notification.created`, `notification.delivered`, `notification.failed`, `notification.read`
- unread counter ve digest projection'larÄ± `docs/shared/projection-strategy.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- template veya category deÄŸiÅŸikliÄŸi, suppression kararÄ±, delivery failure eÅŸiÄŸi ve manual resend aksiyonlarÄ± auditlenmelidir
- queue lag, unread counter drift ve provider failure oranÄ± izlenmelidir

## Ä°dempotency ve Retry
- aynÄ± producer event aynÄ± kullanÄ±cÄ± ve kategori iÃ§in duplicate notification Ã¼retmemelidir
- delivery retry ve backoff davranÄ±ÅŸÄ± `docs/shared/idempotency-policy.md` ile hizalÄ± olmalÄ±dÄ±r
- digest Ã¼retimi tekrarlandÄ±ÄŸÄ±nda aynÄ± batch iÃ§in ikinci bildirim seti oluÅŸturmamalÄ±dÄ±r

## State YapÄ±sÄ±
- created, delivered, failed veya read durumu
- channel veya provider bazlÄ± delivery state'i
- module veya category bazlÄ± delivery pause durumu
- digest eligibility veya quiet-hour bilgisi

## Test NotlarÄ±
- inbox ve preference akÄ±ÅŸlarÄ±
- template, category ve channel Ã§Ã¶zÃ¼mleme testleri
- delivery retry, digest ve dedup doÄŸrulamalarÄ±
- producer event contract ve unread counter projection testleri


---

# Payment ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `payment`

## AmaÃ§
`payment` modÃ¼lÃ¼nÃ¼n amacÄ±, mana satÄ±n alma, Ã¶deme saÄŸlayÄ±cÄ±sÄ± entegrasyonu, ledger doÄŸruluÄŸu ve finansal iÅŸlem kayÄ±tlarÄ±nÄ± merkezi biÃ§imde yÃ¶netmektir.

## Sorumluluk AlanÄ±
- mana package ve checkout session akÄ±ÅŸlarÄ±
- provider callback veya webhook doÄŸrulamasÄ±
- ledger-first iÅŸlem modeli, transaction, ledger entry ve balance snapshot yÃ¶netimi
- reconciliation job, refund, reversal ve fraud hold hazÄ±rlÄ±ÄŸÄ±
- admin tarafÄ±ndan yÃ¶netilen mana purchase, checkout, transaction read veya callback intake runtime kontrolleri ile uyumlu Ã§alÄ±ÅŸma
- provider adapter, webhook verifier ve money value object yaklaÅŸÄ±mÄ±

## Bu ModÃ¼l Neyi Yapmaz?
- Ã¼rÃ¼n kataloÄŸu, offer gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ veya inventory sahipliÄŸi owner'lÄ±ÄŸÄ± taÅŸÄ±maz
- RoyalPass gibi entitlement Ã¼rÃ¼nlerinde final entitlement owner'lÄ±ÄŸÄ±na dÃ¶nÃ¼ÅŸmez
- authorization kararÄ±nÄ± tek baÅŸÄ±na vermez

## Veri SahipliÄŸi
- provider session ve checkout kayÄ±tlarÄ±
- purchase order ve transaction kayÄ±tlarÄ±
- ledger entry ve balance snapshot alanlarÄ±
- provider reference, callback metadata ve audit alanlarÄ±
- fraud review veya finansal inceleme state alanlarÄ±
- reconciliation ve reversal metadata'sÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- shop Ã¼rÃ¼n kataloÄŸu ve fiyat planÄ± verisi
- inventory item sahipliÄŸi veya final reward grant kaydÄ±
- access policy owner verisi

## Access KontratÄ±
`payment` yetki kararÄ± vermez. Mana satÄ±n alma ve iÅŸlem gÃ¶rÃ¼ntÃ¼leme aksiyonlarÄ± `access` ile korunur. ÃœrÃ¼n kataloÄŸu `shop`, final item sahipliÄŸi `inventory` modÃ¼lÃ¼nde kalÄ±r. `payment` devreye girdiÄŸinde `shop` iÃ§indeki geÃ§ici allowance bridge yÃ¼zeyi kaldÄ±rÄ±lmalÄ± ve canonical bakiye yalnÄ±zca `payment` iÃ§inde tutulmalÄ±dÄ±r. Finansal yÃ¼zeylerde precedence ve intake pause davranÄ±ÅŸÄ± `docs/shared/precedence-rules.md` ile hizalÄ± kalmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- mana package listing ve checkout session baÅŸlatma yÃ¼zeyi
- provider callback veya webhook intake yÃ¼zeyi
- own transaction veya own wallet gÃ¶rÃ¼nÃ¼mÃ¼ iÃ§in kontrollÃ¼ okuma yÃ¼zeyi
- `shop` ile bakiye dÃ¼ÅŸÃ¼m veya mutabakat iÃ§in kontrollÃ¼ contract yÃ¼zeyi
- entitlement Ã¼reten modÃ¼ller iÃ§in onaylanmÄ±ÅŸ Ã¶deme veya bakiye gÃ¼ncelleme sonucunu aktaran kontrollÃ¼ fulfillment contract yÃ¼zeyi

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `user`, `access`, `admin`, `shop`, `inventory`, `royalpass` ve payment provider entegrasyonlarÄ±
- webhook verifier, provider adapter, reconcile job ve money helper altyapÄ±sÄ±

## Settings Etkileri
- `feature.payment.mana_purchase.enabled`
- `feature.payment.checkout.enabled`
- `feature.payment.transaction_read.enabled`
- `payment.callback.intake.paused`
- provider bazlÄ± throttling veya manual review surface'i eklenirse settings envanteri geniÅŸletilmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `payment.checkout.started`, `payment.callback.accepted`, `payment.transaction.settled`, `payment.refund.completed`
- tÃ¼ketir: shop purchase intent, provider callback, admin manual review veya reconcile tetikleyicileri
- ledger, reconcile ve fulfillment publish akÄ±ÅŸlarÄ± `docs/shared/outbox-pattern.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- checkout, callback, refund, reversal, manual review ve reconcile aksiyonlarÄ± `docs/shared/audit-policy.md` ile auditlenmelidir
- callback retry oranÄ±, ledger drift, snapshot mismatch ve reconcile backlog izlenmelidir

## Ä°dempotency ve Retry
- callback veya webhook iÅŸleme, refund ve reversal akÄ±ÅŸlarÄ± `docs/shared/idempotency-policy.md` ile hizalÄ± olmalÄ±dÄ±r
- aynÄ± provider event ikinci kez finansal yan etki Ã¼retmemelidir
- reconcile job tekrar Ã§alÄ±ÅŸtÄ±ÄŸÄ±nda ledger-first doÄŸruluÄŸu bozmadan eksik state'i toparlayabilmelidir

## State YapÄ±sÄ±
- pending, success, failed, cancelled veya refunded transaction durumu
- authorized veya captured Ã¶deme durumu gerekiyorsa ayrÄ± state
- fraud_review, reversed veya reconciliation_required durumu
- mana purchase surface pause veya provider outage durumu
- balance snapshot ile ledger iliÅŸki metadata'sÄ±

## Test NotlarÄ±
- checkout session ve callback doÄŸrulama testleri
- ledger, balance snapshot ve mutabakat testleri
- idempotency, replay ve reconcile korumasÄ± doÄŸrulamalarÄ±
- refund, reversal ve fraud hold senaryolarÄ±
- `shop`, `inventory` ve admin inceleme yÃ¼zeyi entegrasyon testleri


---

# RoyalPass ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `royalpass`

## AmaÃ§
`royalpass` modÃ¼lÃ¼nÃ¼n amacÄ±, aylÄ±k sezon bazlÄ± pass yapÄ±sÄ±nÄ±, free veya premium track ilerlemesini ve sezon Ã¶dÃ¼lÃ¼ iÃ§in claim uygunluÄŸu veya claim request yaÅŸam dÃ¶ngÃ¼sÃ¼nÃ¼ yÃ¶netmektir.

## Sorumluluk AlanÄ±
- season, track ve tier yapÄ±larÄ±
- free track ve premium track ayrÄ±mÄ±
- season lifecycle, archive gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ve tier unlock politikasÄ±
- user season progress, claim eligibility ve claim request kayÄ±tlarÄ±
- mission tabanlÄ± royalpass puanÄ± veya progress entegrasyonu
- premium activation source, paused veya frozen progress ve unclaimed reward davranÄ±ÅŸÄ±
- cross-season carryover olup olmadÄ±ÄŸÄ±na dair aÃ§Ä±k kural seti

## Bu ModÃ¼l Neyi Yapmaz?
- gÃ¶rev tanÄ±mÄ± veya inventory sahipliÄŸi owner'lÄ±ÄŸÄ± Ã¼retmez
- Ã¶deme ledger'Ä± veya shop katalog owner'lÄ±ÄŸÄ±na dÃ¶nÃ¼ÅŸmez
- entitlement satÄ±n alma akÄ±ÅŸÄ±nÄ±n tamamÄ±nÄ± kendi iÃ§inde kapatmaz

## Veri SahipliÄŸi
- season tanÄ±mÄ± ve season state alanlarÄ±
- tier reward ve track yapÄ±landÄ±rmalarÄ±
- user season progress, claim eligibility ve premium activation referanslarÄ±
- claim freeze veya season pause metadata alanlarÄ±
- unclaimed reward ve carryover policy metadata'sÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- shop purchase intent veya payment checkout owner verisi
- final inventory grant kaydÄ±
- mission definition owner verisi

## Access KontratÄ±
`royalpass` yetki kararÄ± vermez. Season gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼, claim-request ve premium yÃ¼zey eriÅŸimi `access` ile korunur. Premium aktivasyonun satÄ±n alma kaynaÄŸÄ± farklÄ± modÃ¼llerden gelebilir; ancak Ã¼rÃ¼nleÅŸmiÅŸ satÄ±n alma akÄ±ÅŸÄ±nda canonical purchase intent `shop` Ã¼zerinden baÅŸlamalÄ±, checkout veya bakiye doÄŸruluÄŸu gerekiyorsa `payment` tarafÄ±ndan tamamlanmalÄ± ve season iÃ§i claim uygunluÄŸu ile progress sahipliÄŸi yine `royalpass` modÃ¼lÃ¼nde kalmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- own season overview, own progress ve own reward claim-request yÃ¼zeyi
- admin iÃ§in season yÃ¶netimi, tier yÃ¶netimi ve season pause yÃ¼zeyi
- `mission` progress tÃ¼ketimi, `inventory` reward grant ve `notification` bildirim yÃ¼zeyleri iÃ§in kontrollÃ¼ contract veya event sÄ±nÄ±rÄ±
- `shop`, `payment` veya admin grant akÄ±ÅŸlarÄ±ndan gelen premium activation referanslarÄ± iÃ§in kontrollÃ¼ intake contract yÃ¼zeyi

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `user`, `access`, `admin`, `mission`, `inventory`, `notification`, `shop` ve `payment` entegrasyonlarÄ±
- season scheduler, tier progress helper ve reward claim validator yardÄ±mcÄ±larÄ±

## Settings Etkileri
- `feature.royalpass.claim.enabled`
- `feature.royalpass.season.enabled`
- `feature.royalpass.premium.enabled`
- pause veya carryover yÃ¼zeyi ayrÄ± bir runtime kontrol alÄ±rsa settings envanteri geniÅŸletilmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `royalpass.progressed`, `royalpass.claim.requested`, `royalpass.season.started`, `royalpass.premium.activated`
- tÃ¼ketir: `mission.progressed`, `shop.purchase.completed`, `payment.checkout.confirmed`
- tier snapshot projection'larÄ± `docs/shared/projection-strategy.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- premium activation, season pause, carryover override ve manual claim mÃ¼dahaleleri auditlenmelidir
- season drift, stuck claim ve premium entitlement mismatch oranlarÄ± izlenmelidir

## Ä°dempotency ve Retry
- tier claim ve premium activation intake akÄ±ÅŸlarÄ± `docs/shared/idempotency-policy.md` ile hizalÄ± olmalÄ±dÄ±r
- aynÄ± tier iÃ§in tekrar claim yeni grant Ã¼retmemelidir
- season rollover job tekrarlandÄ±ÄŸÄ±nda aynÄ± sezonu ikinci kez arÅŸivlememeli veya yeni sezonu iki kez aÃ§mamalÄ±dÄ±r

## State YapÄ±sÄ±
- draft, active, paused, ended veya archived season durumu
- free veya premium track eriÅŸim durumu
- tier claim durumu
- season veya claim surface kapanma durumu
- frozen progress veya carryover policy alanlarÄ±

## Test NotlarÄ±
- season ve track Ã§Ã¶zÃ¼mleme testleri
- progress, tier unlock ve claim testleri
- premium activation source ve carryover kurallarÄ± testleri
- mission, inventory, shop, payment ve notification entegrasyonu testleri
- pause, freeze ve claim recovery doÄŸrulamalarÄ±


---

# Shop ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `shop`

## AmaÃ§
`shop` modÃ¼lÃ¼nÃ¼n amacÄ±, Ã¼rÃ¼n kataloÄŸu, teklif gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼, fiyatlandÄ±rma ve satÄ±n alma orkestrasyonunu tek bir maÄŸaza yaÅŸam dÃ¶ngÃ¼sÃ¼ altÄ±nda yÃ¶netmektir.

## Sorumluluk AlanÄ±
- sellable product, offer ve kategori yapÄ±larÄ±
- product veya offer ayrÄ±mÄ± ve time-limited offer kurallarÄ±
- mana bazlÄ± fiyatlandÄ±rma, kampanya gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ve discount evaluatÃ¶r ihtiyacÄ±
- purchase intent, eligibility, checkout handoff, fulfillment ve purchase recovery akÄ±ÅŸlarÄ±
- already-owned item davranÄ±ÅŸÄ± ve fail veya retry politikalarÄ±
- Ã¼rÃ¼n kullanÄ±m kurallarÄ±, slot uyumluluÄŸu, katalog gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ve inventory item mapping

## Bu ModÃ¼l Neyi Yapmaz?
- final item sahipliÄŸi, equip state veya inventory ledger owner'lÄ±ÄŸÄ± taÅŸÄ±maz
- Ã¶deme ledger'Ä± veya provider callback owner'lÄ±ÄŸÄ±na dÃ¶nÃ¼ÅŸmez
- entitlement owner modÃ¼lÃ¼nÃ¼n final state'ini kendi tablosunda canonical hale getirmez

## Veri SahipliÄŸi
- sellable product, offer ve kategori kayÄ±tlarÄ±; ownable inventory item definition alanlarÄ± deÄŸil
- fiyat planÄ±, indirim veya kampanya metadata alanlarÄ±
- purchase request veya order kayÄ±tlarÄ±
- VIP, level veya unlock gereksinimi gibi Ã¼rÃ¼n uygunluk kurallarÄ±
- Ã¼rÃ¼n gÃ¶rÃ¼nÃ¼rlÃ¼k ve kullanÄ±labilirlik state alanlarÄ±
- purchase recovery ve orchestration metadata alanlarÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- inventory item sahipliÄŸi veya grant sonucu
- payment wallet, ledger, checkout ve callback kayÄ±tlarÄ±
- royalpass veya baÅŸka entitlement modÃ¼lÃ¼nÃ¼n final activation owner verisi

## Access KontratÄ±
`shop` yetki kararÄ± vermez. ÃœrÃ¼n gÃ¶rÃ¼ntÃ¼leme, satÄ±n alma ve yÃ¶netim aksiyonlarÄ± `access` ile korunur. Final item sahipliÄŸi ve equip state `inventory`, bakiye veya ledger doÄŸruluÄŸu ise `payment` modÃ¼lÃ¼nde kalmalÄ±dÄ±r. `shop` yalnÄ±zca purchase orkestrasyonu iÃ§in gerekli geÃ§ici doÄŸrulama kÃ¶prÃ¼lerini taÅŸÄ±yabilir. SatÄ±n alma kaynaklarÄ± `docs/shared/purchase-source-types.md` ile hizalÄ± olmalÄ±dÄ±r.

## GeÃ§iÅŸ Notu
- `payment` devreye girene kadar `shop`, yalnÄ±zca Stage 29 iÃ§in tanÄ±mlanan geÃ§ici `seed_mana_allowance_snapshot` veya operasyonel allowance okuma yÃ¼zeyi ile purchase eligibility doÄŸrulayabilir.
- Bu kÃ¶prÃ¼ veri canonical wallet, ledger veya gerÃ§ek bakiye owner'lÄ±ÄŸÄ± sayÄ±lmaz; `payment` modÃ¼lÃ¼ aÃ§Ä±ldÄ±ÄŸÄ±nda kaldÄ±rÄ±lmalÄ± ve yerini `payment` kontratÄ±na bÄ±rakmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- katalog listing, item detail ve purchase request yÃ¼zeyi
- admin katalog, fiyat ve gÃ¶rÃ¼nÃ¼rlÃ¼k yÃ¶netim yÃ¼zeyi
- `inventory` iÃ§in final grant veya teslim talep kontratÄ±
- `payment` iÃ§in bakiye dÃ¼ÅŸÃ¼m, reserve veya mutabakat kontratÄ±
- payment Ã¶ncesi aÅŸamada geÃ§ici `seed_mana_allowance_snapshot` okuma kontratÄ±
- `royalpass` gibi entitlement modÃ¼lleri iÃ§in Ã¼rÃ¼n bazlÄ± premium activation veya fulfillment handoff kontratÄ±

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `user`, `access`, `admin`, `inventory`, `payment`, `royalpass`, `mission` ve `notification` entegrasyonlarÄ±
- pricing calculator, discount evaluator ve purchase orchestration helper ihtiyacÄ±

## Settings Etkileri
- `feature.shop.catalog.enabled`
- `feature.shop.purchase.enabled`
- `feature.shop.campaign.enabled`
- recovery veya checkout handoff yÃ¼zeyi ayrÄ± kontrol alÄ±rsa settings envanterine eklenmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `shop.purchase.intent.created`, `shop.purchase.completed`, `shop.purchase.recovery_requested`
- tÃ¼ketir: payment checkout sonucu, inventory grant sonucu, entitlement activation sonucu
- purchase zinciri `docs/shared/transaction-boundaries.md` ve `docs/shared/outbox-pattern.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- eligibility override, price veya campaign deÄŸiÅŸimi, manual recovery ve duplicate purchase korumasÄ± auditlenmelidir
- checkout handoff baÅŸarÄ±sÄ±, already-owned red oranÄ± ve recovery backlog izlenmelidir

## Ä°dempotency ve Retry
- purchase intent, recovery replay ve fulfillment handoff akÄ±ÅŸlarÄ± `docs/shared/idempotency-policy.md` ile hizalÄ± olmalÄ±dÄ±r
- aynÄ± request iÃ§in ikinci order veya ikinci fulfillment zinciri baÅŸlatÄ±lmamalÄ±dÄ±r
- already-owned senaryosu retry sÄ±rasÄ±nda yeni order state Ã¼retmemelidir

## State YapÄ±sÄ±
- draft, active veya archived product durumu
- visible, hidden veya campaign_only offer durumu
- purchasable, blocked veya sold_out benzeri purchase state alanlarÄ±
- delivery_pending veya recovery_required satÄ±n alma durumu
- checkout handoff ve fulfillment sonuÃ§ durumu

## Test NotlarÄ±
- katalog, product veya offer ayrÄ±mÄ± ve fiyat Ã§Ã¶zÃ¼mleme testleri
- purchase idempotency ve duplicate request doÄŸrulamalarÄ±
- already-owned ve eligibility kurallarÄ± testleri
- `inventory`, geÃ§ici allowance bridge ve `payment` entegrasyon testleri
- recovery akÄ±ÅŸÄ± ve runtime control doÄŸrulamalarÄ±


---

# Social ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `social`

## AmaÃ§
`social` modÃ¼lÃ¼nÃ¼n amacÄ±, kullanÄ±cÄ±lar arasÄ± sosyal iliÅŸki, sosyal duvar ve mesajlaÅŸma yÃ¼zeylerini ayrÄ± bir iÅŸ alanÄ± olarak yÃ¶netmektir.

## Sorumluluk AlanÄ±
- arkadaÅŸlÄ±k isteÄŸi, kabul, reddetme ve arkadaÅŸ listesi
- follow veya unfollow iliÅŸkileri
- friend ve follow farkÄ±nÄ± Ã¼rÃ¼n yÃ¼zeyinde koruyan kullanÄ±m kurallarÄ±
- sosyal duvar post ve duvar altÄ± reply akÄ±ÅŸlarÄ±
- direct message thread, unread state ve mesaj akÄ±ÅŸlarÄ±
- sosyal privacy, block, mute veya restrict gibi iliÅŸki kontrol alanlarÄ±
- online state veya last active ihtiyacÄ± iÃ§in future-ready alanlar

## Bu ModÃ¼l Neyi Yapmaz?
- manga veya chapter yorum thread owner'lÄ±ÄŸÄ± taÅŸÄ±maz
- global authorization kararÄ± Ã¼retmez
- notification delivery veya inventory ownership alanÄ±na girmez

## Veri SahipliÄŸi
- friendship request ve friendship state verileri
- follow relation kayÄ±tlarÄ±
- social block, mute veya restrict iliÅŸkileri
- wall post, wall reply, message thread ve message kayÄ±tlarÄ±
- sosyal gÃ¶rÃ¼nÃ¼rlÃ¼k ve iliÅŸki bazlÄ± eriÅŸim sinyalleri

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- comment modÃ¼lÃ¼ndeki iÃ§erik yorumlarÄ±
- access policy, user profile owner verisi veya notification preference detaylarÄ±
- history, support veya moderation case kayÄ±tlarÄ±

## Access KontratÄ±
`social` yetki kararÄ±nÄ± kendi iÃ§inde Ã¼retmez. Profil duvarÄ±nÄ± kim gÃ¶rebilir, kimin kime mesaj atabileceÄŸi veya arkadaÅŸlÄ±k yÃ¼zeyine eriÅŸim gibi kararlar `social` tarafÄ±ndan Ã¼retilen iliÅŸki veya privacy sinyalleri kullanÄ±larak `access` ile korunur. Block veya aÃ§Ä±k privacy deny sinyali oluÅŸtuÄŸunda final sonuÃ§ `access` tarafÄ±ndan deny olarak yorumlanmalÄ±dÄ±r. `mute` sinyali ise ayrÄ±ca dokÃ¼mante edilmedikÃ§e tek baÅŸÄ±na genel authorization deny sayÄ±lmamalÄ±; teslim, gÃ¶rÃ¼nÃ¼rlÃ¼k veya etkileÅŸim azaltma sinyali olarak kalmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- friendship, follow, wall ve direct message yÃ¼zeyleri
- sosyal duvar reply yapÄ±sÄ± `social`-native kabul edilmeli; `comment` thread sistemine Ã¶rtÃ¼k olarak dÃ¶nÃ¼ÅŸtÃ¼rÃ¼lmemelidir
- `notification` ve `mission` iÃ§in producer event yÃ¼zeyi
- admin tarafÄ±ndan ayrÄ± ayrÄ± aÃ§Ä±lÄ±p kapatÄ±labilen social alt yÃ¼zeyler iÃ§in kontrollÃ¼ contract

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `auth`, `user`, `access`, `notification`, `admin` ve gerektiÄŸinde `mission` entegrasyonlarÄ±
- anti-abuse rate limit ve unread counter projection yardÄ±mcÄ±larÄ±

## Settings Etkileri
- `feature.social.friendship.enabled`
- `feature.social.follow.enabled`
- `feature.social.messaging.enabled`
- `feature.social.wall.enabled`
- restrict veya online presence yÃ¼zeyleri eklenirse settings envanteri geniÅŸletilmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `social.friendship.changed`, `social.follow.changed`, `social.message.sent`, `social.wall.posted`
- tÃ¼ketir: user visibility veya moderation deny sinyalleri
- unread counter ve activity projection'larÄ± `docs/shared/projection-strategy.md` ile hizalÄ± planlanmalÄ±dÄ±r

## Audit ve Ä°zleme
- block, restrict, messaging abuse ve admin kaynaklÄ± privacy mÃ¼dahaleleri auditlenmelidir
- spam, burst message ve block override denemeleri izlenebilir olmalÄ±dÄ±r

## Ä°dempotency ve Retry
- friendship request, follow veya direct message create akÄ±ÅŸlarÄ± duplicate kayÄ±t Ã¼retmemelidir
- message send retry'Ä± aynÄ± message reference iÃ§in ikinci kayÄ±t oluÅŸturmamalÄ±dÄ±r
- block ve mute transition'larÄ± tekrarlandÄ±ÄŸÄ±nda aynÄ± final state korunmalÄ±dÄ±r

## State YapÄ±sÄ±
- friendship_status
- follow veya unfollow iliÅŸkisi durumu
- wall visibility veya message availability durumu
- block, mute veya restrict iliÅŸkisi durumu
- messaging, wall, friendship veya follow surface kapanma durumu

## Test NotlarÄ±
- friendship ve follow akÄ±ÅŸ testleri
- friend ve follow davranÄ±ÅŸ ayrÄ±mÄ± testleri
- wall post veya reply gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ testleri
- message thread, unread ve own-surface testleri
- block, mute, restrict ve runtime control entegrasyonu doÄŸrulamalarÄ±


---

# Support ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `support`

## AmaÃ§
`support` modÃ¼lÃ¼nÃ¼n amacÄ±, kullanÄ±cÄ± iletiÅŸim taleplerini, destek biletlerini ve `manga`, `chapter`, `comment` hedefleri iÃ§in iÃ§erik bildirimlerini tek bir yaÅŸam dÃ¶ngÃ¼sÃ¼ altÄ±nda yÃ¶netmektir.

## Sorumluluk AlanÄ±
- `communication` kaydÄ± veya genel destek bileti oluÅŸturma akÄ±ÅŸlarÄ±
- `manga`, `chapter` ve `comment` iÃ§in hedefe baÄŸlÄ± iÃ§erik bildirimi oluÅŸturma akÄ±ÅŸlarÄ±
- own support list, support detail ve support reply akÄ±ÅŸlarÄ±
- `support_kind`, `category`, `priority`, `reason_code` ve isteÄŸe baÄŸlÄ± hedef iliÅŸkisi
- SLA, duplicate detection, attachment kurallarÄ± ve queue Ã¶nceliklendirme ihtiyaÃ§larÄ±
- internal note ile public reply ayrÄ±mÄ±

## Bu ModÃ¼l Neyi Yapmaz?
- report intake'i otomatik olarak moderation case'e dÃ¶nÃ¼ÅŸtÃ¼rmez
- moderation kararÄ±nÄ±n veya admin final override'Ä±nÄ±n owner'lÄ±ÄŸÄ±na dÃ¶nÃ¼ÅŸmez
- hedef iÃ§erik veya kullanÄ±cÄ± profil verisinin ana owner'lÄ±ÄŸÄ±na girmez

## Veri SahipliÄŸi
- `support_kind=communication`, `support_kind=ticket` veya `support_kind=report` ile taÅŸÄ±nan support kaydÄ±
- `category`, `priority`, `reason_code`, `reason_text` ve isteÄŸe baÄŸlÄ± hedef iliÅŸkisi
- hedefe baÄŸlÄ± kayÄ±tlar iÃ§in `target_type` ve `target_id`; hedefsiz iletiÅŸim veya ticket kayÄ±tlarÄ±nda boÅŸ hedef alanlarÄ±
- internal note, requester reply ve Ã§Ã¶zÃ¼m metadata'sÄ±

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- moderation case karar akÄ±ÅŸÄ±
- hedef iÃ§erik varlÄ±ÄŸÄ±nÄ±n veya yorum kaydÄ±nÄ±n owner verisi
- notification delivery state veya access policy owner'lÄ±ÄŸÄ±

## Access KontratÄ±
`support` yetki kararÄ± vermez. OluÅŸturma, own detail, own reply, review queue ve yÃ¶netimsel karar yÃ¼zeyleri `access` ile korunur; review ve karar yÃ¼rÃ¼tÃ¼mÃ¼ `admin` ile entegre Ã§alÄ±ÅŸÄ±r. Communication kaydÄ±, ticket oluÅŸturma, report oluÅŸturma, attachment kabulÃ¼ ve intake davranÄ±ÅŸlarÄ± admin tarafÄ±ndan yÃ¶netilen runtime ayarlar ile sÄ±nÄ±rlandÄ±rÄ±labilir. Report kaydÄ± varsayÄ±lan olarak moderation case ile aynÄ± kayÄ±t sayÄ±lmaz; `docs/shared/precedence-rules.md` ile hizalÄ± aÃ§Ä±k mapping politikasÄ± gerekir.

## API veya Event SÄ±nÄ±rÄ±
- support create, detail, own list ve own reply yÃ¼zeyi
- communication, ticket ve report create yÃ¼zeyleri iÃ§in veri kontratÄ±
- support review queue ve ticket yÃ¶netimi iÃ§in veri yÃ¼zeyi
- linked moderation case aÃ§Ä±lÄ±ÅŸÄ± iÃ§in kontrollÃ¼ handoff veya reference yÃ¼zeyi

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `access`, `admin`, `moderation`, `notification`, `manga`, `chapter` ve `comment` ile kontrollÃ¼ entegrasyon
- attachment policy iÃ§in `docs/shared/media-asset-strategy.md`
- queue summary ve export yaklaÅŸÄ±mÄ± iÃ§in `docs/shared/reporting-analytics-strategy.md`

## Settings Etkileri
- `feature.support.communication.enabled`
- `feature.support.ticket.enabled`
- `feature.support.report.enabled`
- `feature.support.attachment.enabled`
- `feature.support.internal_note.enabled`
- `support.intake.paused`

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `support.created`, `support.replied`, `support.resolved`, `support.moderation_handoff_requested`
- tÃ¼ketir: moderation case sonucu, admin karar sinyali, notification template veya delivery sinyali
- linked handoff akÄ±ÅŸlarÄ± `docs/shared/transaction-boundaries.md` ve `docs/shared/idempotency-policy.md` ile hizalÄ± olmalÄ±dÄ±r

## Audit ve Ä°zleme
- internal note, resolution, spam iÅŸaretleme, SLA breach ve moderation handoff aksiyonlarÄ± auditlenmelidir
- PII ve attachment riski `docs/shared/operational-standards.md` ile hizalÄ± retention notu taÅŸÄ±malÄ±dÄ±r

## Ä°dempotency ve Retry
- duplicate ticket veya report create istekleri requester, target ve request id kombinasyonu ile korunmalÄ±dÄ±r
- moderation handoff tekrarlandÄ±ÄŸÄ±nda ikinci vaka aÃ§Ä±lmamalÄ±; linked case referansÄ±na baÄŸlanmalÄ±dÄ±r
- attachment processing retry'larÄ± yeni support kaydÄ± Ã¼retmemelidir

## State YapÄ±sÄ±
- open, pending, resolved veya spam durumu
- duplicate veya spam kontrolÃ¼ne temel alanlar
- Ã§Ã¶zÃ¼m, reply ve inceleme yaÅŸam dÃ¶ngÃ¼sÃ¼ alanlarÄ±
- intake pause veya create yÃ¼zeyi kapatma durumlarÄ±

## Test NotlarÄ±
- communication, ticket ve report ayrÄ±mÄ± testleri
- target relation ve handoff doÄŸrulamalarÄ±
- duplicate, spam, attachment ve SLA doÄŸrulamalarÄ±
- `support -> moderation` ve `support -> notification` kontrat testleri


---

# User ModÃ¼lÃ¼

> Canonical modÃ¼l adÄ±: `user`

## AmaÃ§
`user` modÃ¼lÃ¼nÃ¼n amacÄ±, kullanÄ±cÄ± hesabÄ±, profil, tercih, gÃ¶rÃ¼nÃ¼m ve Ã¼yelik verisini taÅŸÄ±yan merkezi kullanÄ±cÄ± alanÄ±nÄ± oluÅŸturmaktÄ±r.

## Sorumluluk AlanÄ±
- kullanÄ±cÄ± hesabÄ± ve profil alanlarÄ±
- public veya private profil ayrÄ±mÄ± ve profile visibility matrix
- hesap durumu, soft delete, deactivation ve ban etkileri
- Ã¼yelik, VIP lifecycle ve benefit freeze sinyalleri
- avatar, banner ve benzeri medya alanlarÄ±nda ownership veya storage referans kurallarÄ±
- global user preference ile module-specific preference ayrÄ±mÄ±nÄ± korumak

## Bu ModÃ¼l Neyi Yapmaz?
- credential, session veya password lifecycle owner'lÄ±ÄŸÄ± taÅŸÄ±maz
- detaylÄ± notification preference, social block veya mute listesi owner'lÄ±ÄŸÄ± taÅŸÄ±maz
- continue reading, history timeline veya inventory sahipliÄŸi Ã¼retmez

## Veri SahipliÄŸi
- username ve profil alanlarÄ±
- hesap durumu alanlarÄ±
- tercih ve gÃ¶rÃ¼nÃ¼m alanlarÄ±
- Ã¼yelik ve kullanÄ±cÄ±ya ait hesap verileri
- VIP entitlement sÃ¼resi ve dondurma veya devam bilgisi
- avatar, banner veya profil kozmetiÄŸi iÃ§in metadata veya referans alanlarÄ±
- public library veya reading activity gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼nÃ¼ etkileyen global preference sinyalleri

## Bu ModÃ¼l Hangi Verinin Sahibi DeÄŸildir?
- auth credential, session ve verification token verileri
- social iliÅŸki kayÄ±tlarÄ±, block veya mute listeleri
- notification category preference detaylarÄ±
- history kayÄ±tlarÄ±, entry-level share metadata'sÄ± veya inventory sahipliÄŸi

## Access KontratÄ±
`user` veriyi taÅŸÄ±r; access kararÄ± vermez. `access` modÃ¼lÃ¼nÃ¼n yorumlayacaÄŸÄ± kullanÄ±cÄ± durumu, Ã¼yelik ve gÃ¶rÃ¼nÃ¼rlÃ¼k sinyalleri kontrollÃ¼ kontrat ile dÄ±ÅŸa aÃ§Ä±lÄ±r. `user` iÃ§indeki global visibility veya sharing preference sinyali, history veya library paylaÅŸÄ±mÄ± iÃ§in Ã¼st sÄ±nÄ±rÄ± belirler; `history` iÃ§indeki entry-level metadata bu sÄ±nÄ±rÄ± daraltabilir ama global deny kararÄ±nÄ± geniÅŸletemez. Visibility kararlarÄ±nda `docs/shared/visibility-states.md` ve `docs/shared/precedence-rules.md` ile hizalÄ± kalÄ±nmalÄ±dÄ±r.

## API veya Event SÄ±nÄ±rÄ±
- user okuma veya yazma yÃ¼zeyi profil ve hesap verisi ile sÄ±nÄ±rlÄ± olmalÄ±dÄ±r
- public ve private response yÃ¼zeyleri aÃ§Ä±kÃ§a ayrÄ±lmalÄ±dÄ±r
- kullanÄ±cÄ±ya ait deÄŸiÅŸiklik olaylarÄ± gerektiÄŸinde diÄŸer modÃ¼llere kontrollÃ¼ yÃ¼zey ile aktarÄ±labilir
- `history` iÃ§in profile baÄŸlÄ± global visibility veya sharing default preference sinyalleri kontrollÃ¼ contract ile dÄ±ÅŸa aÃ§Ä±labilir

## BaÄŸÄ±mlÄ±lÄ±klar
- proje geneli altyapÄ± aÅŸamalarÄ±
- `auth` ile kullanÄ±cÄ± kimliÄŸi entegrasyonu
- `access`, `admin`, `inventory`, `notification` ve `history` iÃ§in kontrollÃ¼ veri okuma yÃ¼zeyleri
- medya ownership ve eriÅŸim politikasÄ± iÃ§in `docs/shared/media-asset-strategy.md`

## Settings Etkileri
- `feature.user.profile.enabled`
- `feature.user.vip_benefits.enabled`
- `feature.user.vip_badge.enabled`
- `feature.user.history_visibility_preference.enabled`
- yeni media veya moderation gÃ¶rÃ¼nÃ¼rlÃ¼k yÃ¼zeyleri eklendiÄŸinde settings envanteri aynÄ± deÄŸiÅŸiklikte gÃ¼ncellenmelidir

## Event AkÄ±ÅŸlarÄ±
- Ã¼retir: `user.profile.updated`, `user.visibility.changed`, `user.vip.changed`, `user.account.deactivated`
- tÃ¼ketir: `auth.identity.created`, `inventory.cosmetic.selected`, `admin.user_state.changed`
- VIP veya profile gÃ¶rÃ¼nÃ¼rlÃ¼k deÄŸiÅŸimleri `access` ve ilgili consumer modÃ¼llere kontrollÃ¼ event veya contract ile aktarÄ±lmalÄ±dÄ±r

## Audit ve Ä°zleme
- VIP state deÄŸiÅŸimi, ban veya deactivation, gÃ¶rÃ¼nÃ¼rlÃ¼k tercihi deÄŸiÅŸimi ve admin kaynaklÄ± profil mÃ¼dahalesi auditlenmelidir
- PII ve export riski taÅŸÄ±yan alanlar `docs/shared/operational-standards.md` ile hizalÄ± maskeleme notu taÅŸÄ±malÄ±dÄ±r

## Ä°dempotency ve Retry
- profile update akÄ±ÅŸlarÄ± request id veya optimistic concurrency ile duplicate write Ã¼retmemelidir
- VIP activation veya freeze akÄ±ÅŸlarÄ± satÄ±n alma veya subscription referansÄ± ile deduplicate edilmelidir
- kullanÄ±cÄ± state transition'larÄ± tekrar Ã§alÄ±ÅŸtÄ±rÄ±ldÄ±ÄŸÄ±nda Ã§eliÅŸkili son durum Ã¼retmemelidir

## State YapÄ±sÄ±
- hesap durumu
- profil gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ve privacy katmanlarÄ±
- Ã¼yelik veya VIP durum alanlarÄ±
- sistem kaynaklÄ± VIP pasifliÄŸinde sÃ¼renin dondurulmasÄ±na iliÅŸkin Ã¼yelik durumu

## Test NotlarÄ±
- profil okuma ve gÃ¼ncelleme testleri
- public veya private response ayrÄ±mÄ± testleri
- history visibility precedence doÄŸrulamalarÄ±
- ban, deactivation ve VIP lifecycle doÄŸrulamalarÄ±


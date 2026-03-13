# Shared ve Settings

> Bu dosya proje icin canonical shared kararlar ve runtime settings referansidir.

## Kullanım Kuralları
- Teknik stack, enum, policy, precedence, transaction ve operasyon kuralları burada canonical referanstır.
- Secret/config kararları ile runtime setting kararları birbirine karıştırılmamalıdır.
- Yeni shared enum, ortak policy veya runtime anahtarı eklendiğinde bu dosya aynı değişiklikte güncellenmelidir.
- Modül dokümanları bu dosyadaki shared kararlarla çelişmemelidir.




---

# Audit Event Tipleri

> Audit kayıtlarında `action` alanını sınıflandırmak veya kategori bazlı filtrelemek gerektiğinde bu sözlük kullanılmalıdır.

## Canonical Kayıtlar
| Değer | Açıklama |
| --- | --- |
| `security_auth` | Kimlik doğrulama ve oturum güvenliği olayları |
| `access_policy` | Policy değişimi, deny override ve availability yorumu |
| `admin_action` | Yönetimsel kritik aksiyon ve settings değişikliği |
| `moderation_action` | Case inceleme, escalation ve içerik aksiyonu |
| `support_case` | Ticket, report, çözüm ve internal note olayları |
| `payment_financial` | Checkout, callback, refund, reversal ve reconcile |
| `inventory_change` | Grant, revoke, consume, equip ve düzeltme |
| `shop_purchase` | Purchase intent, recovery ve eligibility override |
| `user_state` | Ban, deactivate, visibility veya VIP state değişimi |
| `notification_ops` | Template, category, suppression ve delivery failure eşiği |
| `ads_ops` | Campaign publish/pause, click protection ve aggregate müdahalesi |
| `system_ops` | Backup, restore, migration veya genel operasyon aksiyonları |


---

# Audit Politikası

> Güvenlik, yönetim, finansal doğruluk ve kritik state transition'lar için ortak audit formatı bu dokümanda tanımlanır.

## Amaç
- Audit ihtiyacını modüller arasında ortak alan seti ve ortak event sınıflarıyla yönetmek.
- Admin notu, operasyonel not ve immutable audit kaydını birbirine karıştırmamak.
- Güvenlik ve finansal olaylarda minimum izlenebilirlik standardı sağlamak.

## Zorunlu Alanlar
| Alan | Açıklama |
| --- | --- |
| `occurred_at` | Olayın kesin oluşma zamanı |
| `source_module` | Audit kaydını üreten canonical modül |
| `actor_type` / `actor_id` | İşlemi başlatan taraf |
| `target_type` / `target_id` | Etkilenen kullanıcı, içerik veya işlem hedefi |
| `action` | Yapılan aksiyonun canonical adı |
| `result` | `success`, `rejected`, `failed`, `reversed` gibi sonuç |
| `reason` | İnsan okunabilir kısa neden veya reason code |
| `request_id` | Gelen isteğin benzersiz izi |
| `correlation_id` | Akış boyunca taşınan ilişki kimliği |
| `risk_level` | `low`, `medium`, `high`, `critical` |
| `metadata_redaction_level` | PII masking seviyesini anlatan işaret |

## Audit Gerektiren Modül Alanları
| Modül | Minimum Audit Olayları |
| --- | --- |
| `auth` | login failure, suspicious login, password reset, session revoke |
| `access` | policy override, deny reason değişikliği, emergency availability yorumu |
| `admin` | settings değişikliği, override, impersonation, destructive action |
| `moderation` | case action, escalation, reassignment, evidence snapshot erişimi |
| `support` | karar değişimi, internal note, linked moderation handoff |
| `payment` | checkout, callback, refund, reversal, reconcile |
| `inventory` | grant, revoke, consume, manual correction |
| `shop` | purchase orchestration, recovery, eligibility override |
| `user` | VIP state değişimi, ban/deactivation, görünürlük değişimi |
| `ads` | campaign publish/pause, click protection müdahalesi |
| `notification` | template/category değişimi, suppression, delivery failure eşiği |

## Kurallar
- Admin, güvenlik ve finansal olaylar immutable log yaklaşımıyla tutulmalıdır.
- Admin note veya support internal note, audit kaydının yerine geçmemelidir; ayrı veri modeli olarak kalmalıdır.
- Audit export yüzeyi read-only olmalı ve PII masking `docs/shared.md` ile hizalı kalmalıdır.
- Event sınıfları için `docs/shared.md` içindeki canonical sözlük kullanılmalıdır.


---

# Cache ve Queue Stratejisi

> Bu doküman cache backend'i, async işleme baseline'ı ve broker'a geçiş kriterlerini somut karar seviyesinde tanımlar.

## Amaç
- Cache ve queue tarafında “ileride bakarız” belirsizliğini kaldırmak.
- Hangi teknolojinin bugünkü canonical seçim olduğunu netleştirmek.
- Async işleme ile source-of-truth sınırını ayırmak.

## Cache Kararı
- Cache zorunlu baseline değildir; source-of-truth her zaman owner modülün veritabanıdır.
- Cache ihtiyacı oluştuğunda canonical backend `Redis` olmalıdır.
- İlk cache adayları şunlardır:
- access decision cache
- manga listing veya discovery cache
- notification unread counter cache
- ads placement resolve cache

## Cache Kuralları
- Cache yokluğu business doğruluğunu bozmamalıdır; yalnızca performans etkisi yaratmalıdır.
- Cache key'leri subject, surface, selector ve gerekiyorsa settings sürümü ile ilişkilendirilmelidir.
- TTL ve invalidation davranışı modül dokümanı ile `docs/shared.md` notları içinde görünür tutulmalıdır.

## Queue / Async İşleme Kararı
- Bugünkü canonical async işleme baseline'ı PostgreSQL-backed jobs + transactional outbox yaklaşımıdır.
- Message broker baseline mimarinin parçası değildir.
- Outbox, retry, backoff ve dead-letter kuralları `docs/shared.md` ile birlikte değerlendirilmelidir.

## İlk Async İş Yükleri
- notification delivery
- projection rebuild
- payment reconciliation
- ads aggregate jobs
- mission reset
- royalpass season jobs

## Broker'a Geçiş Kriterleri
Aşağıdaki koşulların birden fazlası oluşmadan ayrı broker seçimine gidilmemelidir:
- modüller arası yüksek hacimli fan-out ihtiyacı
- aynı event için bağımsız consumer gruplarında kalıcı lag
- DB-backed job yaklaşımında operasyonel gözlem veya retry maliyetinin aşırı yükselmesi
- dış sistemler veya ayrı servislerle yoğun çift yönlü event entegrasyonu

## Redis'in Rolü
- Redis cache backend'i olarak canonical tercihtir.
- Redis, açık bir mimari güncelleme yapılmadıkça source-of-truth queue olarak kullanılmamalıdır.
- Queue hızlandırma veya ephemeral yardımcı akışlar ancak owner state yine DB'de kaldığı sürece düşünülebilir.


---

# İdempotency Politikası

> Tekrarlanan istek, callback, claim ve grant akışlarında yinelenen yan etkileri önlemek için canonical kurallar bu dokümanda tutulur.

## Amaç
- Aynı isteğin yeniden gelmesi durumunda tekrar ödeme, tekrar grant veya tekrar state transition oluşmasını engellemek.
- Producer ve consumer tarafında safe retry / unsafe retry ayrımını görünür kılmak.
- Modül bazlı özel uygulamaları tek çerçeveye bağlamak.

## Key Formatı
- Varsayılan key biçimi `module:operation:actor_or_scope:client_request_id` olmalıdır.
- Dış provider callback'lerinde client request yerine provider event veya provider transaction referansı kullanılabilir.
- Key ile birlikte payload hash, ilk sonuç durumu ve yan etki referansı saklanmalıdır.

## Kapsam ve Saklama
| Kapsam | Örnek İşlemler | TTL | Saklama Notu |
| --- | --- | --- | --- |
| HTTP write request | register, profile update, purchase intent | en az 24 saat | aynı client request tekrarında ilk sonuç döndürülür |
| Callback / webhook | payment provider callback, external delivery callback | en az 7 gün | provider event id bazlı tutulur |
| Claim / grant | mission claim, royalpass claim, inventory grant | en az 30 gün | reward veya ledger yan etkisi ile birlikte saklanır |
| Background consumer | queue consumer, outbox relay, projection updater | retry penceresi + gözlem süresi | event id veya message id bazlı tutulur |

## Duplicate Handling
- Aynı key ve aynı payload tekrar gelirse önceki sonuç veya mevcut pending durum döndürülmelidir.
- Aynı key fakat farklı payload gelirse istek reddedilmeli ve audit kaydı açılmalıdır.
- In-flight duplicate istekler yeni yan etki üretmemeli; mümkünse mevcut pending işlem referansına bağlanmalıdır.
- İdempotency kaydı durable olarak yazılmadan yan etki başlatılmamalıdır.

## Retry Kuralları
- Safe retry: read model update, callback doğrulama, claim uygunluğu kontrolü, delivery retry, projection rebuild.
- Unsafe retry: yeni ödeme başlatma, tekil reward grant, aynı entitlement'ı tekrar aktive etme, destructive admin action.
- Unsafe retry gerektiren akışlarda ya idempotency key zorunlu tutulmalı ya da manuel review kapısı açılmalıdır.

## Modül Uygulamaları
| Modül | Zorunlu İşlemler | Key Temeli | Not |
| --- | --- | --- | --- |
| `payment` | checkout callback, refund, reversal, reconcile write | provider event id veya checkout request id | ledger-first yaklaşım ile birlikte zorunludur |
| `inventory` | grant, consume, equip | source reference + actor + request id | duplicate grant koruması taşır |
| `mission` | claim request, progress ingest | user + mission + period + request id | period reset ile birlikte düşünülmelidir |
| `royalpass` | tier claim, premium activation intake | user + season + tier veya activation ref | cross-season çakışma önlenmelidir |
| `history` | checkpoint write, merge | user + manga/chapter + device/request | multi-device merge ile birlikte uygulanır |
| `notification` | dedup send, delivery retry | notification id veya producer event ref | category/channel bazlı dedup gerekir |
| `support` | ticket/report create, linked case handoff | requester + target + request id | duplicate spam/report koruması için kullanılır |
| `shop` | purchase intent, recovery replay | user + offer/product + request id | already-owned ve recovery senaryoları ile birlikte çalışır |


---

# Media ve Asset Stratejisi

> Media/asset ihtiyacı yalnızca not seviyesinde bırakılmamalı; bu doküman medya dosyaları, attachment akışları ve signed access politikaları için canonical yönü belirlemelidir.

## Amaç
- Chapter, user, manga ve support gibi alanlardaki medya ihtiyacını ortak kuralla yönetmek.
- Teknik medya yönetimi ile business relation owner'lığını ayırmak.
- Ayrı bir `media` modülünün ne zaman açılacağını netleştirmek.

## Bugünkü Karar
- Hemen ayrı bir leaf `media` modülü açılmayacaktır.
- Canonical yaklaşım, ortak media/asset altyapısının teknik olarak shared platform katmanında planlanmasıdır.
- Domain owner modül, media ilişkisini business referans olarak taşır; teknik metadata ve erişim politikası ortak media altyapısı ile yönetilir.

## Kapsam
- upload metadata
- ownership binding
- mime/type validation
- image dimension metadata
- attachment relation
- signed access policies

## Tüketici Alanlar
- `user`: avatar ve banner
- `manga`: cover, poster veya görsel metadata
- `chapter`: page media ve signed access ihtiyacı
- `support`: attachment intake ve scanning ilişkisi

## Kurallar
- Dosya sahipliği business modülün relation verisi olarak kalmalıdır; blob storage erişimi ortak teknik katmanda ele alınmalıdır.
- MIME, boyut, ölçü ve gerekiyorsa zararlı içerik taraması upload sınırında uygulanmalıdır.
- Support attachment yüzeyi ile chapter page media aynı validasyon seviyesini taşımak zorunda değildir; ortak altyapı farklı policy profilleri taşıyabilmelidir.
- Signed URL veya özel erişim gerekiyorsa token lifetime ve audit beklentisi açıkça dokümante edilmelidir.

## Ayrı Modüle Ayrışma Kriterleri
Aşağıdaki koşulların birden fazlası oluşmadan ayrı leaf `media` modülüne geçilmemelidir:
- birden çok modülün aynı upload, transform ve delivery pipeline'ını yoğun biçimde paylaşması
- media lifecycle'ının business modüllerden bağımsız operasyon yükü üretmesi
- storage, scanning, variant generation ve signed access kurallarının tek başına ayrı ekip veya ayrı servis sınırı gerektirmesi

## İlgili Referanslar
- Runtime ayar envanteri: `docs/shared.md`
- Operasyonel standartlar: `docs/shared.md`
- Cache ve queue stratejisi: `docs/shared.md`


---

# Moderation Durumları

> `moderation` modülündeki lifecycle alanları için canonical sözlük.

## Case Status
| Değer | Anlamı |
| --- | --- |
| `new` | Yeni açılmış ve henüz triage edilmemiş vaka |
| `queued` | İnceleme kuyruğunda bekleyen vaka |
| `assigned` | Moderatöre atanmış vaka |
| `in_review` | Aktif inceleme altında |
| `escalated` | Admin veya daha yüksek scope'a yükseltilmiş |
| `resolved` | Karar verilmiş ve kapanışa uygun |
| `rejected` | Geçersiz veya yanlış intake olarak reddedilmiş |
| `closed` | Finalize edilip kapatılmış |

## Assignment Status
| Değer | Anlamı |
| --- | --- |
| `unassigned` | Atanmamış |
| `assigned` | Bir moderatöre atanmış |
| `handoff_pending` | Scope veya kişi değişimi bekliyor |
| `released` | Atama boşaltılmış veya geri alınmış |


## Escalation Status
| Değer | Anlamı |
| --- | --- |
| `not_escalated` | Escalation uygulanmamış |
| `pending_admin` | Admin tarafından final karar bekleniyor |
| `escalated` | Escalation işlemi tetiklendi |
| `resolved` | Escalation kapatıldı |

## Action Result
| Değer | Anlamı |
| --- | --- |
| `none` | Henüz aksiyon uygulanmadı |
| `content_hidden` | İçerik görünürlüğü kapatıldı |
| `content_restored` | Önceki aksiyon geri alındı |
| `warning_sent` | Kullanıcıya uyarı uygulandı |
| `no_action` | İnceleme sonrası aksiyon gerekmedi |


---

# Notification Kategorileri

> Kullanıcıya dönük bildirim sınıfları için canonical sözlük.

## Canonical Kayıtlar
| Değer | Açıklama |
| --- | --- |
| `account_security` | Login, verification, password reset ve güvenlik uyarıları |
| `social` | Arkadaşlık, takip, direct message ve duvar olayları |
| `comment` | İçerik yorumu ve thread etkileşimleri |
| `support` | Ticket yanıtı, çözüm ve support güncellemeleri |
| `moderation` | Kullanıcıyı etkileyen moderasyon ve inceleme bildirimleri |
| `mission` | Görev ilerleme, tamamlanma ve claim uygunluğu |
| `royalpass` | Season başlangıcı, tier açılması ve reward claim |
| `shop` | Offer, katalog ve satın alma süreci bilgilendirmesi |
| `payment` | Checkout, ödeme sonucu, refund ve cüzdan hareketi |
| `system_ops` | Bakım, kill switch veya ürün genel duyuruları |


---

# Operasyonel Standartlar

> Request context, rate limit, secret yönetimi, backup ve veri saklama gibi modüller üstü işletim kararları bu dokümanda tutulur.

## Request ID ve Correlation ID
- Her dış istek en az bir `request_id` taşımalı veya giriş katmanında üretilmelidir.
- Asenkron zincirde aynı iş akışını bağlayan `correlation_id` korunmalıdır.
- Dış provider callback'lerinde provider event referansı ayrıca saklanmalı; ancak iç korelasyon alanının yerine geçmemelidir.

## Rate Limit Politikası
| Surface | Başlangıç Beklentisi |
| --- | --- |
| `auth.login` | başarısız giriş limiti ve cooldown |
| `comment.write` | flood ve anti-spam koruması |
| `support.intake` | duplicate/spam ve attachment koruması |
| `social.messaging` | anti-abuse ve burst kontrolü |
| `payment.callback` | provider callback replay koruması |
| `ads.click_intake` | invalid click ve bot koruması |

## Secret ve Config Ayrımı
- Env secret, provider credential, signing secret veya private key admin runtime ayarı olarak sunulmamalıdır.
- Runtime settings kullanıcıya dönük availability, threshold ve operasyon davranışı içindir; secret config ile aynı saklama düzeyinde değerlendirilmez.
- Secret rotation ve audit beklentisi ayrı güvenlik süreci olarak ele alınmalıdır.

## Backup, Restore ve Migration Rollback
- Destructive migration öncesinde doğrulanmış backup planı bulunmalıdır.
- Rollback yolu olmayan migration veya veri dönüşümü kritik alanlarda kabul edilmemelidir.
- Restore sonrası özellikle `payment`, projection ve outbox tabanlı modüller için reconcile adımı planlanmalıdır.

## PII ve Veri Saklama
| Modül | PII / Hassas Veri Sınıfı | Saklama / Maskeleme Notu |
| --- | --- | --- |
| `auth` | credential, güvenlik olayı, IP/device izi | uzun süreli güvenlik kaydı hariç maskeli saklama tercih edilir |
| `user` | profil ve üyelik bilgisi | public/private alanlar açık ayrılmalı, gereksiz kopya oluşturulmamalı |
| `support` | talep metni, attachment, kişisel açıklamalar | internal note ve public reply ayrımı zorunlu olmalıdır |
| `social` | direct message ve ilişki sinyalleri | export ve moderasyon yüzeylerinde minimizasyon uygulanmalıdır |
| `payment` | finansal referans ve provider metadata | finansal doğruluk korunurken kullanıcıya dönen yüzey maskeli olmalıdır |

## Ek Notlar
- Modül dokümanları veri retention, masking ve export risklerini kendi alanlarında ayrıca belirtmelidir.
- Audit kayıtları PII taşıyorsa `metadata_redaction_level` alanı zorunlu olmalıdır.


---

# Outbox Deseni

> Modüller arası event üretimi arttığında güvenli publish stratejisinin minimum omurgası bu dokümanda tutulur.

## Amaç
- “DB yaz + hemen publish” yarış durumlarını önlemek.
- Event publish başarısızlıklarında tekrar deneme ve dead-letter akışını standartlaştırmak.
- Replay ve projection rebuild ihtiyacına temel hazırlamak.

## Zorunlu Bileşenler
- Transaction ile birlikte yazılan bir outbox kaydı
- Durumu izlenebilir background publisher
- Retry ve backoff politikası
- Poison message veya kalıcı başarısızlık için dead-letter yaklaşımı
- Publish lag, retry count ve failure rate için gözlem metrikleri

## Öncelikli Modüller
- `payment`
- `inventory`
- `mission`
- `royalpass`
- `notification`
- `support`
- `moderation`
- `history`

## Mesaj Kuralları
- Her mesaj benzersiz `event_id`, `schema_version`, `request_id`, `correlation_id` ve mümkünse `causation_id` taşımalıdır.
- Consumer tarafı idempotent çalışacak biçimde tasarlanmalıdır; publish garantisi consumer tekrarından bağımsız düşünülmemelidir.
- Outbox payload'ı yalnızca consumer'ın ihtiyaç duyduğu kontrat alanlarını taşımalı, producer iç implementasyonunu sızdırmamalıdır.
- Dead-letter'a düşen mesajlar sessizce yok sayılmamalı; operasyonel görünüm ve replay planı taşımalıdır.


---

# Policy Etkileri

> Authorization ve availability kararlarının çıktısı olarak kullanılan canonical etki sözlüğü.

## Canonical Kayıtlar
| Değer | Açıklama |
| --- | --- |
| `allow` | İstek mevcut koşullarda izinli |
| `deny` | İstek kesin olarak reddedilir |
| `deny_soft` | İstek doğrudan reddedilmez; teslim/görünürlük azaltılır veya degrade edilir |
| `require_auth` | Erişim için kimlik doğrulama gerekir |
| `require_role` | Belirli rol veya scope zorunludur |
| `require_entitlement` | VIP, premium veya başka entitlement gerekir |
| `read_only` | Okuma açık kalırken yazma kapatılır |
| `write_off` | İlgili yazma yüzeyi kapalıdır |
| `mask` | Veri dönebilir ancak alan bazlı maskeleme uygulanır |
| `needs_review` | Otomatik karar verilmez; manuel inceleme gerekir |


---

# Precedence Kuralları

> Çapraz modül karar çatışmalarında hangi sinyalin baskın geleceğini bu doküman belirler. Modül dokümanları aynı çatışmayı farklı kelimelerle yeniden tanımlamamalı; burada yazılan kuralları referans almalıdır.

## Amaç
- Çakışan görünürlük, erişim, moderasyon ve operasyon kararlarını tek matris altında toplamak.
- `access`, `admin`, `moderation`, `support`, `history`, `social` ve entitlement akışları arasında ortak yorum sırası oluşturmak.
- Kod tarafında sezgisel kalan “hangisi baskın?” sorusunu dokümante etmek.

## Karar Sırası
1. Sistem genelindeki bakım veya acil durdurma kararı değerlendirilir.
2. İlgili modül veya surface için runtime kill switch ve availability anahtarı değerlendirilir.
3. Security, moderation veya admin override kaynaklı explicit deny/override kararları uygulanır.
4. Audience, role, entitlement ve ownership tabanlı access kararı yorumlanır.
5. Modülün kendi visibility/share metadata'sı yalnızca kalan izinli aralık içinde uygulanır.
6. Rate limit, cooldown veya backpressure kaynaklı geçici reddetme son aşamada uygulanır.

## Canonical Kurallar
| Çakışma | Kazanan | Gerekçe |
| --- | --- | --- |
| `system deny` vs `entitlement allow` | `system deny` | Bakım, güvenlik veya operasyonel kapatma ücretli haklardan daha üst precedence taşır. |
| `runtime kill switch` vs `normal access allow` | `runtime kill switch` | Surface kapalıysa permission tek başına erişim açamaz. |
| `admin hard override` vs `scoped moderator action` | `admin hard override` | Günlük vaka sahipliği `moderation` içinde kalsa da final yönetimsel karar `admin` tarafında baskındır. |
| `moderation block` vs `social visibility allow` | `moderation block` | İçerik veya kullanıcı güvenliği kaynaklı blok kararı sosyal görünürlükten üstündür. |
| `social block` vs `messaging/wall allow` | `social block` | Açık block ilişkisi direct message ve wall etkileşimini durdurur. |
| `social mute` vs `authorization allow` | `authorization allow` | `mute` varsayılan olarak teslim/görünürlük azaltma sinyalidir; tek başına final authorization deny sayılmaz. |
| `user global visibility deny` vs `history entry share opt-in` | `user global visibility deny` | `history` entry-level paylaşım kararı global deny tavanını aşamaz. |
| `support report` vs `moderation case` | Ayrı kayıtlar | Report intake otomatik olarak moderation case sayılmaz; açık mapping politikası gerekir. |
| `vip no-ads` veya başka entitlement muafiyeti vs `ads/payment/support` kill switch | Kill switch | Sistem güvenliği ve operasyon kararı ürün avantajının üstünde değerlendirilir. |
| `provider callback retry` vs `payment callback intake pause` | Intake pause | Provider tekrar denemesi açık olsa bile operasyonda intake geçici durdurulabilir. |

## Uygulama Notları
- `access` modülü kullanıcıya bakan availability ve permission kararlarının son yorumlayıcısıdır.
- Operasyonel pause, callback intake veya queue backpressure gibi kullanıcıya doğrudan görünmeyen kontroller ilgili `service` katmanı tarafından yorumlanabilir.
- Aynı precedence kuralı hem modül dokümanında hem settings envanterinde farklı ifadelerle yazılmamalıdır.


---

# Projection Stratejisi

> Event veya change-feed ile beslenen read model, counter ve denormalize özet yüzeyleri bu dokümanla hizalı olmalıdır.

## Amaç
- Canonical write model ile projection sorumluluğunu ayırmak.
- Eventual consistency beklentisini ve rebuild yolunu görünür kılmak.
- Counter, summary ve read model büyümesini modüller arasında tutarlı hale getirmek.

## Temel İlkeler
- Her projection için canonical write model owner modülde kalır; consumer modül owner tabloya doğrudan yazmaz.
- Projection güncellemeleri mümkün olduğunda event, outbox veya açık projection contract yüzeyi ile beslenir.
- Kabul edilen eventual consistency penceresi dokümante edilmeden denormalize alan açılmaz.
- Projection rebuild ve replay akışı en baştan planlanır; “yalnızca incremental çalışır” kabulü yapılmaz.
- Projection hataları ve lag durumu izlenebilir metrik üretebilmelidir.

## Canonical Projection Kayıtları
| Projection | Canonical Write Model | Event Kaynağı | Consumer Surface | Tutarlılık Penceresi | Rebuild Yolu | Replay |
| --- | --- | --- | --- | --- | --- | --- |
| `manga.comment_count` | `comment` | `comment.created`, `comment.deleted`, `comment.moderated` | manga detail ve listing | kısa | comment tablosundan recount + incremental catch-up | desteklenir |
| `manga.engagement_summary` | `history`, `comment` | read checkpoint ve engagement event'leri | discovery ve admin özetleri | orta | günlük batch rebuild + hedefli repair | desteklenir |
| `history.continue_reading_projection` | `history` | checkpoint ve finish event'leri | continue reading yüzeyi | kısa | son checkpoint'ten recompute | desteklenir |
| `notification.unread_counter` | `notification` | `notification.created`, `notification.read` | inbox badge ve header sayaçları | kısa | unread kayıtlardan recount | desteklenir |
| `support.queue_summary` | `support` | create, status_change, assignee_change | support operasyon paneli | kısa | status bazlı regroup | desteklenir |
| `moderation.queue_summary` | `moderation` | case create, assignment, resolution | moderation paneli | kısa | case durumlarından regroup | desteklenir |
| `ads.impression_aggregate` | `ads` | accepted impression ve click event'leri | reporting ve dashboard yüzeyi | orta | batch aggregation job | desteklenir |
| `mission.progress_projection` | `mission` | progress, claim, reset event'leri | mission liste ve progress özetleri | kısa | objective bazlı recompute | desteklenir |
| `royalpass.tier_progress_snapshot` | `royalpass` | progress ve claim event'leri | season overview ve tier görünümü | kısa | tier bazlı recompute | desteklenir |

## Uygulama Kuralları
- Projection rebuild işlemi idempotent olmalı ve aynı veri için tekrar çalıştırıldığında yeni yan etki üretmemelidir.
- Replay yapılacak event payload'ları şema sürümü, `request_id` ve `correlation_id` taşımalıdır.
- Counter veya summary güncellemeleri için doğrudan “owner tabloya başka modül yazsın” yaklaşımı kullanılmamalıdır.
- Event üreten modüller `docs/shared.md` ile hizalı transactional outbox yaklaşımı planlamalıdır.


---

# Purchase Source Tipleri

> Satın alma, checkout veya fulfillment başlatan canonical kaynak sözlüğü.

## Canonical Kayıtlar
| Değer | Tipik Owner | Açıklama |
| --- | --- | --- |
| `catalog_purchase` | `shop` | Normal ürün veya offer satın alma isteği |
| `premium_activation` | `shop`, `royalpass` | Premium pass veya entitlement aktivasyon akışı |
| `mana_wallet` | `payment` | İç bakiye kullanılarak yapılan harcama |
| `external_provider` | `payment` | Harici ödeme sağlayıcısı ile başlayan akış |
| `recovery_replay` | `shop`, `payment` | Başarısız veya kesilmiş akışın güvenli tekrar yürütümü |
| `admin_issue` | `admin` | Manuel operasyon kaynağı |
| `gift_code` | gelecekte ayrı modül | Kod veya kupon bazlı aktivasyon |


---

# Reporting ve Analytics Stratejisi

> Dashboard ve rapor ihtiyacı yalnızca modül içine dağılmış notlarla yönetilmemelidir. Bu doküman reporting/analytics read model yaklaşımını, export sınırlarını ve operasyon özetlerini canonical karara dönüştürür.

## Amaç
- Admin, ads, payment, support ve benzeri alanların raporlama ihtiyacını ortak bir modelde toplamak.
- Operasyon summary ile analytics aggregate katmanını ayırmak.
- Export-friendly query yüzeylerini write model ownership'ini bozmadan tanımlamak.

## Bugünkü Karar
- Ayrı bir analytics write modülü baseline değildir.
- Reporting yaklaşımı projection, aggregate read model ve kontrollü export query layer üstünden kurulmalıdır.
- Write model owner'lığı ilgili iş modülünde kalır; reporting yalnızca okuma odaklı görünüm üretir.

## Katmanlar
- Operasyon summary: admin ve ekiplerin anlık durum görmesi için düşük gecikmeli özetler
- Analytics aggregate: trend, hacim, funnel veya performans yorumları için periyodik projection'lar
- Export query layer: denetim, finans veya operasyon ihtiyaçları için kontrollü dışa aktarma yüzeyi

## Kurallar
- Reporting projection'ları canonical write model'in yerini almamalıdır.
- Export yüzeyleri yalnızca yetkili operasyon akışlarıyla açılmalı ve audit zorunluluğu taşımalıdır.
- `payment` reconciliation, refund ve fraud review summary; `ads` performans aggregate; `support` queue summary gibi yüzeyler aynı read-model prensipleriyle dokümante edilmelidir.
- Dashboard metriği ile karar verici business state birbirine karıştırılmamalıdır.

## Ne Zaman Ayrılaştırılır?
Aşağıdaki koşullar belirginleşmeden ayrı analytics servisi veya veri hattına geçilmemelidir:
- operasyon summary ile analytics sorgularının aynı DB yükünde sürdürülememesi
- büyük hacimli history ve event akışlarının ayrı depolama gerektirmesi
- ayrı ekip, ayrı erişim politikası veya ayrı veri saklama kuralları ihtiyacı oluşması

## İlgili Referanslar
- Projection stratejisi: `docs/shared.md`
- Audit politikası: `docs/shared.md`
- Cache ve queue stratejisi: `docs/shared.md`
- Admin modülü: `docs/modules.md`


---

# Reward Source Tipleri

> Ödül, grant veya entitlement benzeri sahiplik akışlarında canonical kaynak sözlüğü.

## Canonical Kayıtlar
| Değer | Tipik Owner | Açıklama |
| --- | --- | --- |
| `mission` | `mission` | Görev tamamlanması veya claim sonucu oluşan ödül |
| `royalpass` | `royalpass` | Season tier veya pass claim sonucu oluşan ödül |
| `shop` | `shop` | Satın alma orkestrasyonu sonucu gelen grant |
| `admin_grant` | `admin` | Manuel operasyon veya destek amaçlı grant |
| `compensation` | `admin`, `support` | Sistem hatası telafisi |
| `seasonal_event` | `mission`, `royalpass` | Dönemsel event kampanyası kaynağı |
| `referral` | gelecekte ayrı modül veya `user` | Referral benzeri teşvik kaynağı |
| `reconciliation_repair` | `payment`, `inventory` | Reconcile veya data repair sonucu düzeltme grant'i |


---

# Search Stratejisi

> Search katmanı modül içi kısa not olarak bırakılmamalı; baseline arama yaklaşımı, reindex süreci ve provider swap kriterleri bu dokümanda tutulmalıdır.

## Amaç
- Search için bugünkü canonical kararı netleştirmek.
- Reindex, fallback ve provider değişimi konularını modül notu olmaktan çıkarıp ortak kural haline getirmek.
- `manga` başta olmak üzere arama kullanan yüzeylerin aynı omurgaya bağlanmasını sağlamak.

## Bugünkü Karar
- Başlangıç search engine kararı `PostgreSQL full-text search` olmalıdır.
- Ayrı search provider veya servis baseline mimarinin parçası değildir.
- Search indeks yapısı write model'i gölgeleyen ikinci source-of-truth haline getirilmemelidir.

## Kapsam
- manga title, slug, summary ve taxonomy araması
- basic ranking ve filter kombinasyonları
- gerektiğinde editorial collection araması
- reindex ve drift recovery süreci

## Reindex ve Lifecycle Kuralları
- Search index kaynağı canonical write model'dir; owner veri `manga` modülünde kalır.
- Reindex tam rebuild veya scoped rebuild olarak çalışabilmelidir.
- Publish, archive veya taxonomy değişimi search projection güncellemesini tetiklemelidir.
- Reindex işlemleri replay-safe ve idempotent olmalıdır.

## Provider Swap Kriterleri
Aşağıdaki koşullar netleşmeden ayrı search provider'a geçilmemelidir:
- PostgreSQL full-text ile kabul edilebilir arama kalitesi veya latency sağlanamaması
- typo tolerance, weighted ranking, synonym expansion veya faceted search ihtiyacının baseline'ı aşması
- yüksek hacimli ayrı index operasyonlarının DB yükünü sürdürülemez hale getirmesi

## Fallback Davranışı
- Search provider veya projection geçici sorun yaşarsa listing fallback'i kontrollü biçimde daraltılabilmelidir.
- Hata durumunda owner içerik görünürlüğü bozulmamalı; gerekirse search yüzeyi geçici olarak pasife alınmalıdır.
- Search availability kararları `docs/shared.md` ile hizalı anahtarlar üzerinden yönetilmelidir.

## İlgili Referanslar
- Projection stratejisi: `docs/shared.md`
- Runtime ayar envanteri: `docs/shared.md`
- Manga modülü: `docs/modules.md`


---

# Support Durumları

> `support` modülündeki intake ve çözüm lifecycle'ı için canonical sözlük.

## Support Status
| Değer | Anlamı |
| --- | --- |
| `open` | Yeni açılmış kayıt |
| `triaged` | İlk sınıflandırması yapılmış |
| `waiting_user` | Kullanıcı yanıtı bekleniyor |
| `waiting_team` | İç ekip işlemi bekleniyor |
| `resolved` | Çözüm uygulanmış |
| `rejected` | Geçersiz veya kapsam dışı kayıt |
| `closed` | Süreç tamamlanmış ve kapatılmış |
| `spam` | Spam veya kötüye kullanım olarak işaretlenmiş |

## Reply Visibility
| Değer | Anlamı |
| --- | --- |
| `public_to_requester` | Talep sahibine görünür yanıt |
| `internal_only` | Sadece ekip içi not veya yanıt |


---

# Hedef Tipleri

> Canonical kayıt dosyası: `target_type` veya benzeri hedef tipi taşıyan tüm modüller bu dosyadaki değerleri kullanmalıdır.

## Amaç
- `target_type` değerlerinin modüller arasında tutarlı kalmasını sağlamak.
- Yeni hedef tipleri eklenirken owner modülü, consumer modülleri ve kullanım amacını görünür tutmak.
- `comment`, `moderation`, `support` ve ileride hedefe bağlı çalışacak diğer modüller için tek kaynak sunmak.

## Kullanım Kuralları
- `target_type` yalnızca hedefe bağlı kayıtlarda kullanılmalıdır; `support_kind=communication` veya hedefsiz `support_kind=ticket` kayıtlarında `target_type` ve `target_id` boş bırakılabilir.
- `target_type` değeri mümkün olduğunda canonical leaf modül adı ile aynı olmalıdır.
- Görünüm, ekran, aksiyon veya alt yüzey bilgisi `target_type` içine gömülmemelidir; bu bilgi ayrı alan veya context verisi ile taşınmalıdır.
- Yeni `target_type` değeri eklendiğinde aynı değişiklik setinde bu dosya, ilgili owner modül dokümanı ve ilgili consumer modül dokümanı güncellenmelidir.
- Kullanılmayan veya kaldırılan `target_type` değeri dokümanda durum notu ile işaretlenmeden sessizce silinmemelidir.

## Canonical Kayıtlar
| `target_type` | Owner Module | Başlangıç Consumer'lar | Açıklama | Status | Notes |
| --- | --- | --- | --- | --- | --- |
| `manga` | `manga` | `comment`, `moderation`, `support` | Manga veya seri düzeyindeki hedef varlık. | `active` | Public içerik ve bildirim hedefi. |
| `chapter` | `chapter` | `comment`, `moderation`, `support` | Okunabilir bölüm hedefi. | `active` | Okuma ve ihlal bildirimleri için kullanılır. |
| `comment` | `comment` | `moderation`, `support` | Yorum hedefi. | `active` | Yorum bildirimi ve inceleme akışları için kullanılır. |
| `social` | `social` | `-` | Sosyal duvar, ilişki ve mesajlaşma yüzeylerinin canonical hedef alanı için ayrılmış rezerv kayıt. | `active` | İleride moderation veya support ihtiyacı doğarsa consumer aynı değişiklik setinde eklenmelidir. |


---

# Teknik Stack ve Araç Seçimleri

> Öneri seviyesinde kalan paket ve araç seçimleri bu dokümanda canonical teknik karara dönüştürülür. Aynı sorumluluk için ikinci bir varsayılan araç seçilmemelidir.

## Amaç
- Repo geneline dağılmış teknik tercihleri tek bir aktif karara dönüştürmek.
- Kod yazımı başladığında hangi temel kütüphane ve araçların kullanılacağını netleştirmek.
- Setup, test, cache ve altyapı kararlarını ortak referansla bağlamak.

## Backend ve Uygulama Araçları
| Sorumluluk | Canonical Seçim | Not |
| --- | --- | --- |
| HTTP router | `chi` | Hafif, açık ve Go ekosistemiyle uyumlu tercih. |
| Config / env loader | `caarlos0/env` | Env tabanlı config standardı için kullanılır. |
| Structured logging | `zap` | Yüksek performanslı structured logging standardı. |
| UUID | `google/uuid` | Domain ve contract tarafında ortak UUID üretimi. |
| Input validation | `go-playground/validator/v10` | DTO ve request validation için canonical katman. |
| Migration | `golang-migrate` | Tüm DB migration akışlarının canonical aracı. |
| SQL erişimi | `pgx/v5` | PostgreSQL için canonical driver. |
| Connection pool | `pgxpool` | `pgx` ile birlikte canonical pooling katmanı. |
| Password hashing | `argon2id` | Credential güvenliği için default seçim. |
| Test assertion / helper | `testify` | Go stdlib test akışı yanında helper ve assertion desteği. |

## Platform ve Çalışma Kararları
- Docker-first çalışma standardı korunmalıdır.
- Main ve test veritabanı kesin olarak ayrı tutulmalıdır.
- Cache ihtiyaçlarında canonical backend `Redis` olmalıdır.
- Asenkron işleme baseline kararı PostgreSQL-backed jobs + transactional outbox olmalıdır.

## Uygulama Kuralları
- Aynı sorumluluk için farklı bir araç kullanılacaksa önce bu doküman güncellenmelidir.
- Wrapper veya adapter katmanı açılması canonical seçimi değiştirmez; alttaki teknik tercih aynı kalmalıdır.
- Logging, validation ve config erişimi platform katmanında ortaklaştırılmalıdır.
- Access policy değerlendirmesi için dış ağır policy engine varsayılan çözüm değildir; hafif in-house evaluator tercih edilmelidir.

## İlgili Referanslar
- Kurulum ve komut standardı: `docs/SETUP.md`
- Cache ve queue kararı: `docs/shared.md`
- Operasyon ve secret/config ayrımı: `docs/shared.md`


---

# Transaction Sınırları

> Hangi akışın tek transaction içinde, hangisinin event destekli çok aşamalı zincirle çalışması gerektiğini bu doküman tanımlar.

## Amaç
- Tek owner'lı akışlar ile çok modüllü orkestrasyonları ayırmak.
- Eventual consistency gereken akışlar için net sınır çizmek.
- Recovery, reconcile ve manuel müdahale ihtiyacını önceden görünür kılmak.

## Seçim Kuralları
- Tek modül, tek veritabanı transaction'ı ve tek owner aggregate içinde kalan akışlar mümkün olduğunda tek transaction kullanmalıdır.
- Dış provider, queue veya başka modül owner'lığı gerektiren akışlar çok aşamalı ve idempotent orkestrasyon ile yürütülmelidir.
- Bir transaction içinde “DB yaz + hemen publish” yaklaşımı yerine transactional outbox planı tercih edilmelidir.
- Kompanzasyon veya recovery gerektiren her akış dokümante edilmeden “sonra bakarız” yaklaşımıyla bırakılmamalıdır.

## Referans Akışlar
| Akış | Boundary | Gerekçe | Recovery / Kompanzasyon |
| --- | --- | --- | --- |
| `auth login + session create` | tek transaction | credential doğrulama ve session write aynı owner alanda kalır | session revoke ve güvenlik audit'i ile geri alınır |
| `shop purchase -> payment callback -> inventory grant` | çok aşamalı, event destekli | katalog, ödeme ve sahiplik üç ayrı owner alandır | reconcile, duplicate guard ve grant retry gerekir |
| `mission complete -> reward grant` | koordineli ama idempotent | claim uygunluğu ile final grant farklı owner'lardadır | claim replay ve grant dedup ile toparlanır |
| `support report -> moderation case` | policy'ye göre sync veya async | backpressure ve queue ihtiyacı oluşabilir | linked case reference ile tekrar denenir |
| `notification create -> channel delivery` | write + async delivery | delivery dış kanal ve retry gerektirir | backoff, suppression ve dead-letter gerekir |
| `royalpass premium activation` | çok aşamalı, event destekli | `shop`, `payment` ve `royalpass` owner'lıkları ayrıdır | activation ref, reconcile ve replay güvenli olmalıdır |

## Uygulama Kuralları
- Başka modülün tablosuna doğrudan yazı yapan transaction tasarlanmamalıdır.
- `request_id`, `correlation_id` ve idempotency key kritik geçişlerde boundary boyunca taşınmalıdır.
- Finansal veya kullanıcı hakkı etkileyen akışlarda manuel review kapısı açıkça işaretlenmelidir.


---

# Görünürlük Durumları

> Görünürlük anlamı taşıyan modüller mümkün olduğunda bu canonical sözlüğe hizalanmalıdır.

## Kullanım Kuralları
- Her modül tüm durumları kullanmak zorunda değildir.
- Modül içi özel state alanları bu sözlüğe map edilebiliyorsa modül dokümanında açıkça belirtilmelidir.
- Final erişim kararı yine `docs/shared.md` ve `access` yorumu ile verilir.

## Canonical Durumlar
| Durum | Anlamı | Tipik Modüller | Not |
| --- | --- | --- | --- |
| `public` | Herkese açık görünür yüzey | `manga`, `chapter`, `user`, `social` | access ve audience kontrolü ayrıca uygulanabilir |
| `limited` | Belirli audience veya entitlement ile görünür | `chapter`, `royalpass`, `ads` | VIP, early access veya scoped audience ile birlikte kullanılır |
| `private` | Yalnızca owner veya yetkili yüzeye görünür | `history`, `support`, `notification` | public default üretmez |
| `hidden` | Mevcut kayıt durur ama dış görünürlük kapanır | `comment`, `manga`, `social` | moderation veya runtime switch kaynaklı olabilir |
| `removed` | Dış yüzeyden kaldırılmış veya soft-deleted görünüm | `comment`, `support` | audit ve recovery notu taşımalıdır |
| `archived` | Aktif değil ama geçmiş kayıt olarak tutulur | `manga`, `royalpass`, `ads`, `shop` | read-only veya operasyonel görünümde kalabilir |


---

# Runtime Ayar Envanteri

> Canonical kayıt dosyası: admin tarafından yönetilen tüm runtime ayarlar, feature toggle'lar, kill switch yüzeyleri ve oran limitleri bu dokümanda tutulmalıdır.

## Amaç
- Runtime ayarların tek merkezden izlenmesini sağlamak.
- Ayar anahtarlarının çakışmadan büyümesini sağlamak.
- Hangi ayarın hangi modülde üretildiğini ve hangi katmanda yorumlandığını görünür tutmak.

## Yaşayan Doküman Notu
- Bu dosya kapanmış bir checklist değil, modül yüzeyleri büyüdükçe güncellenen canonical runtime envanteridir.
- Modül dokümanlarında runtime ile kontrol edilebilir olduğu yazılan her surface burada birebir temsil edilmeli veya neden henüz `planned` kaldığı açıkça not edilmelidir.

## Kayıt Kuralları
- Yeni runtime ayar, feature toggle, kill switch veya oran limiti eklendiğinde bu dosya aynı değişiklik setinde güncellenmelidir.
- Aynı `key + audience_kind + audience_selector + scope_kind + scope_selector` kombinasyonu için birden fazla aktif kayıt bırakılmamalıdır.
- Ücretli veya süreli avantajı etkileyen ayarlarda `entitlement_impact_policy` alanı zorunludur.
- `schedule_support` alanı en az `none`, `start_at` ve `time_window` modlarını ayırt edebilmelidir.
- `docs/RULES.md` ve ilgili modül dokümanı ile çelişen kayıt bırakılamaz.
- Bir modül dokümanında ayrı ayrı runtime kontrol edilebildiği yazılan yüzeyler settings envanterinde ya kendi canonical key kaydıyla ya da kapsadığı alt yüzeyleri açıkça listeleyen umbrella key kaydıyla temsil edilmelidir.
- Yeni bir surface için yalnızca availability anahtarı yazmak yeterli değildir; ilgili rate limit, threshold, cooldown, disabled behavior veya degrade davranışı da aynı değişiklikte eklenmeli ya da neden henüz `planned` kaldığı `notes` alanında belirtilmelidir.
- Service katmanında yorumlanan intake pause, callback gate, digest window veya benzeri operasyonel ayarlar için `consumer_layer` alanı zorunlu olarak doldurulmalı ve access availability ayarları ile karıştırılmamalıdır.

## Audience Sözlüğü
- `all`: tüm kullanıcılar.
- `guest`: giriş yapmamış ziyaretçi.
- `authenticated`: giriş yapmış tüm kullanıcılar.
- `authenticated_non_vip`: giriş yapmış ancak VIP olmayan kullanıcılar.
- `vip`: aktif VIP avantajı taşıyan kullanıcılar.

## Scope Sözlüğü
- `site`: tüm ürün yüzeyi.
- `module`: belirli modül geneli.
- `feature`: modül içindeki belirli özellik veya alt yüzey.
- `resource/context`: belirli kaynak, ekran veya özel bağlam.

## Key Grammar
- Boolean availability veya feature toggle anahtarları mümkün olduğunda `feature.<module>.<surface>.enabled` biçimini kullanmalıdır.
- Eşik, limit, cooldown veya davranış değeri taşıyan anahtarlar mümkün olduğunda `<module>.<surface>.<metric>` biçimini kullanmalıdır.
- Feature toggle niteliği taşımayan operasyonel pause, intake veya bakım bayrakları gerektiğinde `<module>.<surface>.<flag>` biçimini kullanabilir.
- Site geneli operasyon veya bakım anahtarları mümkün olduğunda `site.<surface>.<metric_or_flag>` biçimini kullanmalıdır.
- Audience, role, grup, kullanıcı veya resource bilgisi runtime key içine gömülmemelidir; bu bilgiler selector alanlarında taşınmalıdır.
- Key içindeki `<module>` bölümü her zaman canonical leaf modül adı ile aynı yazılmalıdır.

## Selector Grammar
- `scope_selector` için başlangıç canonical yaklaşım `-`, `<module>`, `<module>.<surface>`, `<module>.<surface>.<subsurface>` ve gerektiğinde `resource:<module>:<resource_kind>:<identifier>` biçimleridir.
- `audience_selector` için başlangıç canonical yaklaşım `-`, `role:<name>`, `group:<name>` ve `user:<id>` biçimleridir.
- Selector içindeki modül adı her zaman canonical leaf modül adı ile aynı yazılmalıdır.

## Disabled Behavior Sözlüğü
- `visibility_off`: yüzeyin dışarıya tamamen kapatılması.
- `read_only`: yalnızca okuma izni verilip yazma tarafının kapatılması.
- `write_off`: mevcut veri görünür kalırken yeni yazma aksiyonunun kapatılması.
- `intake_pause`: yeni kayıt alımının durdurulması.
- `read_only_intake`: mevcut kayıtların okunur kaldığı, ancak yeni intake/create akışının durdurulduğu durum.
- `attachment_off`: ana create akışı açık kalırken attachment kabulünün kapatılması.
- `preview_off`: yalnızca önizleme yüzeyinin kapatılması.
- `benefit_pause`: ücretli veya süreli avantajın geçici pasife alınıp sürenin dondurulması.

## Error Response Policy Sözlüğü
- `not_found`: dış yüzeyde kaynak veya yüzey yokmuş gibi davranılması.
- `forbidden`: yüzey görünür olsa bile ilgili aksiyonun reddedilmesi.
- `rate_limited`: eşik, cooldown veya throttling nedeniyle geçici reddedilme.
- `validation_error`: istek biçimi geçerli olsa bile ilgili alt yüzey kapalı olduğu için alan bazlı doğrulama hatası dönülmesi.
- `temporarily_unavailable`: sistem kaynaklı geçici pasiflik veya operasyonel durdurma.

## Entitlement Impact Policy Sözlüğü
- `none`: süreli hak veya entitlement üzerinde ek etki yoktur.
- `freeze_on_system_disable`: sistem kaynaklı pasiflikte kalan entitlement süresi dondurulur ve hak süresi korunur.

## Kayıt Alanları
| Alan | Açıklama |
| --- | --- |
| `key` | Canonical ayar anahtarı |
| `description` | Ayarın neyi kontrol ettiğini açıklayan kısa metin |
| `category` | `site`, `communication`, `operations`, `security_auth`, `access_availability`, `content`, `reading`, `engagement`, `support`, `membership`, `social`, `gamification`, `economy` |
| `owner_module` | Ayarın sahibi olan modül |
| `consumer_layer` | Ayarı yorumlayan katman: örnek `access`, `service`, `app` |
| `value_type` | `bool`, `int`, `duration`, `enum`, `string` vb. |
| `default_value` | Varsayılan değer |
| `allowed_range_or_enum` | İzin verilen aralık veya enum listesi |
| `scope_kind` | `site`, `module`, `feature`, `resource/context` |
| `scope_selector` | Scope'un somut hedefi; örnek `comment.write`, `manga.detail`, `chapter.preview`, `-` |
| `audience_kind` | Geçerli audience türleri |
| `audience_selector` | Audience'ın somut hedefi; örnek `role:comment_moderator`, `group:testers`, `-` |
| `sensitive` | Hassas veri içerip içermediği |
| `apply_mode` | `immediate`, `cache_refresh`, `scheduled` |
| `cache_strategy` | `none`, `ttl`, `manual_invalidate` |
| `schedule_support` | `none`, `start_at`, `time_window` |
| `audit_required` | Değişiklik audit zorunluluğu |
| `affected_surfaces` | Etkilenen yüzeyler |
| `disabled_behavior` | Varsa kapatma davranış tipi |
| `error_response_policy` | Dış API veya yüzey için hata/geri dönüş politikası |
| `entitlement_impact_policy` | Varsa süreli hak etkisi |
| `status` | `active`, `planned`, `deprecated` |
| `notes` | Ek açıklama |

## Access Yorumlama Modeli
- Kullanıcıya bakan availability, permission ve entitlement etkili runtime ayarları canonical olarak `access` tarafından yorumlanmalıdır.
- Yorum sırası `docs/shared.md` ile hizalı olmalı; `global -> module -> surface -> audience -> entitlement -> action -> rate limit` akışı korunmalıdır.
- Operasyonel pause, callback intake, queue backpressure veya delivery retry benzeri kullanıcıya görünmeyen ayarlar ilgili `service` katmanında yorumlanabilir; bu durum settings envanterinde `consumer_layer` ile görünür tutulmalıdır.
- Aynı surface için read ve write ayrımı varsa ayrı anahtar veya umbrella key notu bırakılmalıdır; örtük davranış bırakılmamalıdır.

## Kill Switch Seviyeleri
| Seviye | Açıklama | Tipik Kullanım |
| --- | --- | --- |
| `global` | Tüm ürün veya site yüzeyi etkilenir | bakım modu, genel güvenlik durdurması |
| `module-level` | Bir modülün ana yüzeyleri topluca kapanır | `support.intake.paused`, `notification.delivery.paused` |
| `surface-level` | Modül içindeki belirli feature kapanır | `feature.social.messaging.enabled` |
| `action-level` | Surface açık kalırken tek aksiyon kapanır | `feature.inventory.claim.enabled` |
| `intake-only` | Yeni kayıt alımı durur, mevcut veri okunabilir kalır | support intake, payment callback intake |
| `write-only` | Yazma kapalı, okuma açık kalır | comment write, bookmark write |
| `external-integration-only` | Dış provider veya delivery kanalı kapatılır | payment callback, notification channel send |

## Modül Yüzeyi Hizalama Matrisi
| Module | Surface | Subsurface | Canonical / Planned Key | Disabled Behavior | Error Response Policy | Audience Selector Support | Schedule Support | Entitlement Impact | Rate Limit Support | Note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `manga` | discovery | recommendation, collection, editorial | `feature.manga.discovery.enabled` | `visibility_off` | `not_found` | var | `none` | `none` | yok | Listing ve detail yüzeyinden bağımsız kalmalıdır. |
| `chapter` | read | preview, detail, early access | `feature.chapter.preview.enabled`, `feature.chapter.detail.enabled`, `feature.chapter.read.enabled`, `feature.chapter.early_access.enabled` | `preview_off`, `visibility_off` | `not_found` | var | `time_window` | `none` | yok | Early access görünürlüğü ayrı toggle taşımalıdır. |
| `comment` | write | create, edit, delete | `feature.comment.write.enabled`, `comment.write.cooldown_seconds` | `write_off` | `forbidden`, `rate_limited` | var | `none` | `none` | var | Anti-spam ve edit window service katmanında yorumlanır. |
| `history` | library | continue reading, timeline, bookmark.write | `feature.history.continue_reading.enabled`, `feature.history.library.enabled`, `feature.history.timeline.enabled`, `feature.history.bookmark_write.enabled` | `visibility_off`, `write_off` | `not_found`, `forbidden` | var | `none` | `none` | yok | Entry-level share metadata ayrı access yorumuna bağlanır. |
| `social` | interaction | friendship, follow, wall, messaging | `feature.social.friendship.enabled`, `feature.social.follow.enabled`, `feature.social.wall.enabled`, `feature.social.messaging.enabled` | `write_off`, `visibility_off` | `forbidden`, `not_found` | var | `none` | `none` | var | Block ve mute precedence access tarafında yorumlanır. |
| `notification` | delivery | inbox, preference, digest, channel send | `feature.notification.inbox.enabled`, `feature.notification.preference.enabled`, `feature.notification.digest.enabled`, `notification.delivery.paused` | `visibility_off`, `read_only` | `not_found`, `temporarily_unavailable` | var | `time_window` | `none` | var | Digest ve channel pause service tarafında da yorumlanır. |
| `inventory` | ownership | read, claim, equip, consume | `feature.inventory.read.enabled`, `feature.inventory.claim.enabled`, `feature.inventory.equip.enabled`, `feature.inventory.consume.enabled` | `visibility_off`, `write_off` | `not_found`, `forbidden` | var | `none` | `none` | yok | Grant tarafı idempotent çalışmalı, read/write ayrımı korunmalıdır. |
| `mission` | progression | read, claim, progress_ingest | `feature.mission.read.enabled`, `feature.mission.claim.enabled`, `feature.mission.progress_ingest.enabled` | `visibility_off`, `write_off`, `intake_pause` | `not_found`, `forbidden`, `temporarily_unavailable` | var | `time_window` | `none` | yok | Event ingest ayrı kapatılabilir olmalıdır. |
| `royalpass` | season | season, premium, claim | `feature.royalpass.season.enabled`, `feature.royalpass.premium.enabled`, `feature.royalpass.claim.enabled` | `benefit_pause`, `visibility_off`, `write_off` | `temporarily_unavailable`, `not_found`, `forbidden` | var | `time_window` | `freeze_on_system_disable` | yok | Premium entitlement ve season availability ayrı yorumlanmalıdır. |
| `shop` | purchase | catalog, campaign, purchase, recovery | `feature.shop.catalog.enabled`, `feature.shop.campaign.enabled`, `feature.shop.purchase.enabled` | `visibility_off`, `write_off` | `not_found`, `forbidden` | var | `time_window` | `none` | var | Recovery ve already-owned davranışı service katmanında tamamlanır. |
| `payment` | checkout | mana purchase, checkout, transaction read, callback intake | `feature.payment.mana_purchase.enabled`, `feature.payment.checkout.enabled`, `feature.payment.transaction_read.enabled`, `payment.callback.intake.paused` | `visibility_off`, `write_off`, `intake_pause` | `not_found`, `temporarily_unavailable` | var | `time_window` | `none` | var | Provider callback yüzeyi kullanıcı availability'sinden ayrı yönetilmelidir. |
| `ads` | delivery | surface, placement, campaign, click intake | `feature.ads.surface.enabled`, `feature.ads.placement.enabled`, `feature.ads.campaign.enabled`, `feature.ads.click_intake.enabled` | `visibility_off`, `intake_pause` | `not_found`, `temporarily_unavailable` | var | `time_window` | `none` | var | VIP no-ads precedence access ile yorumlanır. |
| `support` | intake | communication, ticket, report, attachment, internal_note | `feature.support.communication.enabled`, `feature.support.ticket.enabled`, `feature.support.report.enabled`, `feature.support.attachment.enabled`, `feature.support.internal_note.enabled`, `support.intake.paused` | `write_off`, `attachment_off`, `read_only_intake` | `forbidden`, `validation_error`, `temporarily_unavailable` | var | `time_window` | `none` | var | Report-to-case policy ve internal note görünürlüğü ayrı tutulmalıdır. |
## Başlangıç Referans Kayıtları
| Key | Description | Category | Owner Module | Consumer Layer | Value Type | Default | Allowed Range / Enum | Scope Kind | Scope Selector | Audience Kind | Audience Selector | Sensitive | Apply Mode | Cache Strategy | Schedule Support | Audit Required | Affected Surfaces | Disabled Behavior | Error Response Policy | Entitlement Impact | Status | Notes |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `site.maintenance.enabled` | Bakım modu anahtarı | `operations` | `admin` | `app` | `bool` | `false` | `true,false` | `site` | `-` | `all` | `-` | `false` | `immediate` | `none` | `time_window` | `true` | `site.*` | `visibility_off` | `temporarily_unavailable` | `none` | `planned` | Bakım modu için referans anahtar |
| `auth.login.failed_attempt_limit_per_minute` | Dakika başına başarısız giriş limiti | `security_auth` | `auth` | `service` | `int` | `5` | `1-20` | `site` | `-` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `auth.login` | `-` | `rate_limited` | `none` | `planned` | Başarısız giriş limiti |
| `auth.login.cooldown_seconds` | Başarısız giriş eşiği sonrası uygulanan cooldown süresi | `security_auth` | `auth` | `service` | `int` | `300` | `0-86400` | `site` | `-` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `auth.login` | `-` | `rate_limited` | `none` | `planned` | Failed login limiti ile birlikte çalışan geçici auth kısıtı |
| `auth.email.verification_resend_cooldown_seconds` | Verification tekrar gönderim bekleme süresi | `security_auth` | `auth` | `service` | `int` | `60` | `0-3600` | `site` | `-` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `auth.email_verification` | `-` | `rate_limited` | `none` | `planned` | Verification tekrar gönderim aralığı |
| `feature.comment.read.enabled` | Yorum okuma ve thread listeleme yüzeyini açma-kapama anahtarı | `access_availability` | `comment` | `access` | `bool` | `true` | `true,false` | `feature` | `comment.read` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `comment.read`,`comment.thread` | `visibility_off` | `not_found` | `none` | `planned` | Read yüzeyi write yüzeyinden bağımsız daraltılabilir |
| `comment.write.cooldown_seconds` | Yorum yazma aksiyonu için bekleme süresi | `engagement` | `comment` | `service` | `int` | `30` | `0-3600` | `feature` | `comment.write` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `comment.write` | `-` | `rate_limited` | `none` | `planned` | Yorum yazma aralığı |
| `feature.comment.write.enabled` | Yorum yazma yüzeyini açma-kapama anahtarı | `access_availability` | `comment` | `access` | `bool` | `true` | `true,false` | `feature` | `comment.write` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `comment.write` | `write_off` | `forbidden` | `none` | `planned` | Yorum yazma yüzeyi açma-kapama |
| `feature.user.profile.enabled` | Profil görünürlüğü ve profil detail yüzeyini açma-kapama anahtarı | `access_availability` | `user` | `access` | `bool` | `true` | `true,false` | `feature` | `user.profile` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `user.profile`,`user.profile.detail` | `visibility_off` | `not_found` | `none` | `planned` | Public profil yüzeyi own olmayan okumadan bağımsız daraltılabilir |
| `feature.user.vip_benefits.enabled` | VIP avantajlarını sistem genelinde açma-kapama anahtarı | `access_availability` | `user` | `access` | `bool` | `true` | `true,false` | `feature` | `user.vip_benefits` | `vip` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `user.vip_benefits`,`access.vip` | `benefit_pause` | `temporarily_unavailable` | `freeze_on_system_disable` | `planned` | VIP avantajının sistem kaynaklı pasifliğinde süre dondurulur |
| `feature.user.vip_badge.enabled` | VIP rozet veya VIP profil göstergesini açma-kapama anahtarı | `access_availability` | `user` | `access` | `bool` | `true` | `true,false` | `feature` | `user.vip_badge` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `user.vip_badge`,`user.profile` | `visibility_off` | `not_found` | `none` | `planned` | Görsel profil göstergesi VIP entitlement sahibinden ayrı yönetilebilir |
| `feature.user.history_visibility_preference.enabled` | Kullanıcının history veya library visibility preference yüzeyini açma-kapama anahtarı | `access_availability` | `user` | `access` | `bool` | `true` | `true,false` | `feature` | `user.history_visibility_preference` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `user.history_visibility_preference` | `write_off` | `forbidden` | `none` | `planned` | Preference yüzeyi kapansa bile history owner'lığı `history` modülünde kalır |
| `feature.manga.list.enabled` | Manga listing, search ve filter yüzeyini açma-kapama anahtarı | `access_availability` | `manga` | `access` | `bool` | `true` | `true,false` | `feature` | `manga.list` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `manga.list`,`manga.search`,`manga.filter` | `visibility_off` | `not_found` | `none` | `planned` | Listing yüzeyi discovery veya detail yüzeyinden bağımsız yönetilebilir |
| `feature.manga.detail.enabled` | Manga detail yüzeyini açma-kapama anahtarı | `access_availability` | `manga` | `access` | `bool` | `true` | `true,false` | `feature` | `manga.detail` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `manga.detail` | `visibility_off` | `not_found` | `none` | `planned` | Detail yüzeyi listing ve discovery'den bağımsız daraltılabilir |
| `feature.chapter.preview.enabled` | Chapter preview yüzeyini açma-kapama anahtarı | `access_availability` | `chapter` | `access` | `bool` | `true` | `true,false` | `feature` | `chapter.preview` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `chapter.preview` | `preview_off` | `not_found` | `none` | `planned` | Preview kapanınca detail veya tam read otomatik kapanmış sayılmaz |
| `feature.chapter.detail.enabled` | Chapter detail yüzeyini açma-kapama anahtarı | `access_availability` | `chapter` | `access` | `bool` | `true` | `true,false` | `feature` | `chapter.detail` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `chapter.detail` | `visibility_off` | `not_found` | `none` | `planned` | Detail yüzeyi preview ve tam read'den bağımsız yönetilebilir |
| `feature.chapter.read.enabled` | Chapter tam read yüzeyini açma-kapama anahtarı | `access_availability` | `chapter` | `access` | `bool` | `true` | `true,false` | `feature` | `chapter.read` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `chapter.read`,`chapter.resume` | `visibility_off` | `not_found` | `none` | `planned` | Tam read kapalıyken detail veya preview ayrı açık kalabilir |
| `feature.moderation.panel.enabled` | Moderation paneli genel görünürlük anahtarı | `access_availability` | `moderation` | `access` | `bool` | `true` | `true,false` | `feature` | `moderation.panel` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `moderation.panel` | `visibility_off` | `not_found` | `none` | `active` | Yetkili moderatör yüzeyi admin tarafından geçici olarak kapatılabilir |
| `feature.moderation.queue.enabled` | Moderation queue görünürlüğünü açma-kapama anahtarı | `access_availability` | `moderation` | `access` | `bool` | `true` | `true,false` | `feature` | `moderation.queue` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `moderation.queue` | `visibility_off` | `not_found` | `none` | `active` | Queue yüzeyi panel açık kalsa bile ayrı yönetilebilir |
| `feature.moderation.action.enabled` | Moderation aksiyon yüzeyini açma-kapama anahtarı | `access_availability` | `moderation` | `access` | `bool` | `true` | `true,false` | `feature` | `moderation.action` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `moderation.action` | `write_off` | `forbidden` | `none` | `active` | Sınırlı moderatör aksiyonları geçici olarak kapatılabilir |
| `feature.notification.inbox.enabled` | In-app bildirim kutusu yüzeyini açma-kapama anahtarı | `access_availability` | `notification` | `access` | `bool` | `true` | `true,false` | `feature` | `notification.inbox` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `notification.inbox` | `visibility_off` | `temporarily_unavailable` | `none` | `active` | Bildirim kutusu geçici olarak pasife alınabilir |
| `feature.notification.preference.enabled` | Bildirim preference yönetim yüzeyini açma-kapama anahtarı | `access_availability` | `notification` | `access` | `bool` | `true` | `true,false` | `feature` | `notification.preference` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `notification.preference` | `visibility_off` | `not_found` | `none` | `active` | Preference yüzeyi inbox görünürlüğünden bağımsız daraltılabilir |
| `notification.delivery.paused` | Notification delivery akışını geçici durdurma anahtarı | `operations` | `notification` | `service` | `bool` | `false` | `true,false` | `module` | `notification` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `notification.delivery`,`notification.channel`,`notification.category` | `read_only` | `temporarily_unavailable` | `none` | `active` | Category veya channel bazlı durdurma aynı key'in selector genişlemesiyle uygulanabilir |
| `feature.social.friendship.enabled` | Arkadaşlık yüzeyini açma-kapama anahtarı | `access_availability` | `social` | `access` | `bool` | `true` | `true,false` | `feature` | `social.friendship` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `social.friendship` | `write_off` | `forbidden` | `none` | `active` | Arkadaşlık isteği ve ilişki yönetimi yüzeyi ayrı olarak kapatılabilir |
| `feature.social.follow.enabled` | Takip yüzeyini açma-kapama anahtarı | `access_availability` | `social` | `access` | `bool` | `true` | `true,false` | `feature` | `social.follow` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `social.follow` | `write_off` | `forbidden` | `none` | `active` | Follow veya unfollow yazma yüzeyi ayrı olarak yönetilebilir |
| `feature.social.messaging.enabled` | Sosyal mesajlaşma gönderim yüzeyini açma-kapama anahtarı | `access_availability` | `social` | `access` | `bool` | `true` | `true,false` | `feature` | `social.messaging` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `social.messaging` | `write_off` | `forbidden` | `none` | `active` | Mesajlaşma yazma kapatılsa bile mevcut thread okuma davranışı ayrı ayarlanabilir |
| `feature.social.wall.enabled` | Sosyal duvar yüzeyini açma-kapama anahtarı | `access_availability` | `social` | `access` | `bool` | `true` | `true,false` | `feature` | `social.wall` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `social.wall`,`profile.wall` | `visibility_off` | `not_found` | `none` | `active` | Profil duvarı veya sosyal duvar yüzeyi yönetilebilir |
| `feature.support.communication.enabled` | İletişim kaydı oluşturma yüzeyini açma-kapama anahtarı | `access_availability` | `support` | `access` | `bool` | `true` | `true,false` | `feature` | `support.communication` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.communication`,`support.create` | `write_off` | `forbidden` | `none` | `planned` | Genel iletişim girişi geçici olarak kapatılabilir |
| `feature.support.ticket.enabled` | Destek bileti oluşturma yüzeyini açma-kapama anahtarı | `access_availability` | `support` | `access` | `bool` | `true` | `true,false` | `feature` | `support.ticket` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.ticket`,`support.create` | `write_off` | `forbidden` | `none` | `planned` | Hedefsiz destek bileti oluşturma yüzeyi geçici olarak kapatılabilir |
| `feature.support.report.enabled` | Hedefe bağlı içerik bildirimi yüzeyini açma-kapama anahtarı | `access_availability` | `support` | `access` | `bool` | `true` | `true,false` | `feature` | `support.report` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.report` | `write_off` | `forbidden` | `none` | `planned` | Manga, chapter veya comment için report oluşturma yüzeyi ayrı olarak yönetilebilir |
| `feature.support.attachment.enabled` | Support attachment kabulünü açma-kapama anahtarı | `support` | `support` | `service` | `bool` | `true` | `true,false` | `feature` | `support.attachment` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.communication`,`support.ticket`,`support.report` | `attachment_off` | `validation_error` | `none` | `planned` | Attachment kapalıysa kayıt akışı devam eder, ancak dosya kabul edilmez |
| `support.intake.paused` | Support yeni kayıt alımını durdurma anahtarı | `support` | `support` | `service` | `bool` | `false` | `true,false` | `module` | `support` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `support.create`,`support.communication`,`support.ticket`,`support.report` | `read_only_intake` | `temporarily_unavailable` | `none` | `planned` | Intake pause aktifken mevcut kayıtlar okunabilir kalırken yeni create yüzeyleri durdurulabilir |
| `feature.inventory.read.enabled` | Envanter list ve own item detail yüzeyini açma-kapama anahtarı | `access_availability` | `inventory` | `access` | `bool` | `true` | `true,false` | `feature` | `inventory.read` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `inventory.list`,`inventory.detail` | `visibility_off` | `not_found` | `none` | `active` | Liste ve detail görünürlüğü claim veya equip yüzeyinden bağımsız yönetilebilir |
| `feature.inventory.claim.enabled` | Envanter ödül claim yüzeyini açma-kapama anahtarı | `access_availability` | `inventory` | `access` | `bool` | `true` | `true,false` | `feature` | `inventory.claim` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `inventory.claim` | `write_off` | `forbidden` | `none` | `active` | Claim kapansa bile mevcut sahiplik kayıtları görünür kalabilir |
| `feature.inventory.equip.enabled` | Envanter equip yüzeyini açma-kapama anahtarı | `access_availability` | `inventory` | `access` | `bool` | `true` | `true,false` | `feature` | `inventory.equip` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `inventory.equip` | `write_off` | `forbidden` | `none` | `active` | Equip yüzeyi read veya claim yüzeyinden bağımsız kapatılabilir |
| `mission.daily.reset_hour_utc` | Günlük mission reset saati | `gamification` | `mission` | `service` | `int` | `0` | `0-23` | `module` | `mission` | `all` | `-` | `false` | `cache_refresh` | `manual_invalidate` | `none` | `true` | `mission.reset`,`mission.daily` | `-` | `-` | `none` | `active` | Dönemsel mission reset davranışını belirleyen referans eşik |
| `feature.mission.read.enabled` | Mission list ve progress görünümünü açma-kapama anahtarı | `access_availability` | `mission` | `access` | `bool` | `true` | `true,false` | `feature` | `mission.read` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `mission.read`,`mission.progress` | `visibility_off` | `not_found` | `none` | `active` | Görev görünümü claim yüzeyinden bağımsız yönetilebilir |
| `feature.mission.claim.enabled` | Mission claim-request yüzeyini açma-kapama anahtarı | `access_availability` | `mission` | `access` | `bool` | `true` | `true,false` | `feature` | `mission.claim` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `mission.claim` | `write_off` | `forbidden` | `none` | `active` | Görev ilerlemesi açık kalırken claim-request yüzeyi geçici olarak kapatılabilir |
| `feature.royalpass.claim.enabled` | RoyalPass claim-request yuzeyini acma-kapama anahtari | `access_availability` | `royalpass` | `access` | `bool` | `true` | `true,false` | `feature` | `royalpass.claim` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `royalpass.claim` | `write_off` | `forbidden` | `none` | `active` | Claim-request yuzeyi gecici olarak kapatilabilir; season ilerleme verisi ayri yonetilir |
| `feature.royalpass.season.enabled` | RoyalPass sezon avantajlarini sistem genelinde acma-kapama anahtari | `access_availability` | `royalpass` | `access` | `bool` | `true` | `true,false` | `feature` | `royalpass.season` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `royalpass.season`,`royalpass.track`,`royalpass.claim` | `benefit_pause` | `temporarily_unavailable` | `freeze_on_system_disable` | `active` | Sezon sistem kaynakli pasiflige alindiginda kalan hak ve sure guvenli bicimde dondurulur |
| `feature.royalpass.premium.enabled` | RoyalPass premium track yuzeyini acma-kapama anahtari | `access_availability` | `royalpass` | `access` | `bool` | `true` | `true,false` | `feature` | `royalpass.premium` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `royalpass.premium`,`royalpass.track` | `visibility_off` | `not_found` | `none` | `active` | Premium entitlement ownerligi degismeden premium gorunurluk ayri yonetilebilir |
| `feature.history.continue_reading.enabled` | Continue reading yüzeyini açma-kapama anahtarı | `access_availability` | `history` | `access` | `bool` | `true` | `true,false` | `feature` | `history.continue_reading` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `history.continue_reading` | `visibility_off` | `not_found` | `none` | `active` | Continue reading yüzeyi history timeline veya library yüzeyinden bağımsız kapatılabilir |
| `feature.history.library.enabled` | Kütüphane okuma yüzeyini açma-kapama anahtarı | `access_availability` | `history` | `access` | `bool` | `true` | `true,false` | `feature` | `history.library` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `history.library` | `visibility_off` | `not_found` | `none` | `active` | Library read yüzeyi continue reading ve timeline yüzeylerinden bağımsız yönetilebilir |
| `feature.history.timeline.enabled` | Okuma geçmişi timeline yüzeyini açma-kapama anahtarı | `access_availability` | `history` | `access` | `bool` | `true` | `true,false` | `feature` | `history.timeline` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `history.timeline` | `visibility_off` | `not_found` | `none` | `active` | Timeline yüzeyi continue reading ve library yüzeylerinden bağımsız kapatılabilir |
| `feature.history.bookmark_write.enabled` | Bookmark yazma yüzeyini açma-kapama anahtarı | `access_availability` | `history` | `access` | `bool` | `true` | `true,false` | `feature` | `history.bookmark.write` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `history.bookmark.write` | `write_off` | `forbidden` | `none` | `active` | Library görünür kalsa bile bookmark write yüzeyi ayrıca kapatılabilir |
| `feature.manga.discovery.enabled` | Recommendation, koleksiyon ve editoryal keşif yüzeylerini açma-kapama anahtarı | `access_availability` | `manga` | `access` | `bool` | `true` | `true,false` | `feature` | `manga.discovery` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `manga.recommendation`,`manga.collection`,`manga.discovery` | `visibility_off` | `not_found` | `none` | `planned` | Editoryal keşif kapatıldığında temel manga listing ve detail yüzeyi ayrı çalışabilir |
| `feature.ads.surface.enabled` | Reklam gösterim üst yüzeylerini açma-kapama anahtarı | `access_availability` | `ads` | `access` | `bool` | `true` | `true,false` | `feature` | `ads.surface` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `ads.surface`,`ads.placement`,`ads.delivery` | `visibility_off` | `not_found` | `none` | `planned` | Global ads switch placement ve campaign alt anahtarlarını override edebilir |
| `feature.ads.placement.enabled` | Placement çözümleme yüzeyini açma-kapama anahtarı | `access_availability` | `ads` | `access` | `bool` | `true` | `true,false` | `feature` | `ads.placement` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `ads.placement` | `visibility_off` | `not_found` | `none` | `planned` | Placement resolve yüzeyi campaign serve akışından bağımsız daraltılabilir |
| `feature.ads.campaign.enabled` | Campaign serve yüzeyini açma-kapama anahtarı | `access_availability` | `ads` | `access` | `bool` | `true` | `true,false` | `feature` | `ads.campaign` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `ads.campaign`,`ads.delivery` | `visibility_off` | `not_found` | `none` | `planned` | Placement tanımları kalsa bile campaign serve ayrı kapatılabilir |
| `feature.shop.catalog.enabled` | Shop katalog goruntuleme yuzeyini acma-kapama anahtari | `access_availability` | `shop` | `access` | `bool` | `true` | `true,false` | `feature` | `shop.catalog` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `shop.catalog`,`shop.item.view` | `visibility_off` | `not_found` | `none` | `active` | Katalog yuzeyi purchase aksiyonundan bagimsiz yonetilebilir |
| `feature.shop.purchase.enabled` | Shop satin alma yuzeyini acma-kapama anahtari | `access_availability` | `shop` | `access` | `bool` | `true` | `true,false` | `feature` | `shop.purchase` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `shop.purchase` | `write_off` | `forbidden` | `none` | `active` | Shop katalog gorunur kalirken satin alma aksiyonu ayri kapatilabilir |
| `feature.shop.campaign.enabled` | Shop campaign gorunurlugunu acma-kapama anahtari | `access_availability` | `shop` | `access` | `bool` | `true` | `true,false` | `feature` | `shop.campaign` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `shop.campaign`,`shop.catalog` | `visibility_off` | `not_found` | `none` | `active` | Kampanya badge veya spotlight yuzeyi katalog okumasindan bagimsiz daraltilabilir |
| `feature.payment.mana_purchase.enabled` | Mana package ve purchase entry yüzeyini açma-kapama anahtarı | `access_availability` | `payment` | `access` | `bool` | `true` | `true,false` | `feature` | `payment.mana_purchase` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `payment.mana_purchase` | `visibility_off` | `not_found` | `none` | `planned` | Checkout başlatma akışı ayrı anahtar ile yönetilir |
| `feature.payment.checkout.enabled` | Checkout session başlatma yüzeyini açma-kapama anahtarı | `access_availability` | `payment` | `access` | `bool` | `true` | `true,false` | `feature` | `payment.checkout` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `payment.checkout`,`payment.provider_session` | `write_off` | `temporarily_unavailable` | `none` | `planned` | Provider sorunu halinde checkout durdurulurken mana package görünürlüğü ayrı kalabilir |
| `feature.payment.transaction_read.enabled` | Kullanıcının kendi transaction veya wallet görünümünü açma-kapama anahtarı | `access_availability` | `payment` | `access` | `bool` | `true` | `true,false` | `feature` | `payment.transaction.read` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `payment.transaction.read`,`payment.wallet.read` | `visibility_off` | `not_found` | `none` | `planned` | İşlem geçmişi görünümü mana purchase ve checkout akışlarından bağımsız daraltılabilir |
| `feature.chapter.early_access.enabled` | Chapter early access görünürlüğünü açma-kapama anahtarı | `access_availability` | `chapter` | `access` | `bool` | `true` | `true,false` | `feature` | `chapter.early_access` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `chapter.early_access`,`chapter.read` | `visibility_off` | `not_found` | `none` | `planned` | Early access görünürlüğü preview ve tam read yüzeyinden bağımsız daraltılabilir |
| `feature.notification.digest.enabled` | Notification digest üretim ve teslim yüzeyini açma-kapama anahtarı | `access_availability` | `notification` | `service` | `bool` | `true` | `true,false` | `feature` | `notification.digest` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `notification.digest`,`notification.delivery` | `read_only` | `temporarily_unavailable` | `none` | `active` | Digest kapatılsa bile tekil in-app inbox açık kalabilir |
| `feature.inventory.consume.enabled` | Envanter consume yüzeyini açma-kapama anahtarı | `access_availability` | `inventory` | `access` | `bool` | `true` | `true,false` | `feature` | `inventory.consume` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `inventory.consume` | `write_off` | `forbidden` | `none` | `active` | Consume aksiyonu read veya equip yüzeyinden bağımsız kapatılabilir |
| `feature.mission.progress_ingest.enabled` | Mission progress event ingest yüzeyini açma-kapama anahtarı | `gamification` | `mission` | `service` | `bool` | `true` | `true,false` | `feature` | `mission.progress_ingest` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `mission.progress`,`mission.ingest` | `intake_pause` | `temporarily_unavailable` | `none` | `active` | Event ingest geçici kapatıldığında mission read yüzeyi açık kalabilir |
| `payment.callback.intake.paused` | Payment provider callback intake yüzeyini geçici durdurma anahtarı | `operations` | `payment` | `service` | `bool` | `false` | `true,false` | `feature` | `payment.callback` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `payment.callback`,`payment.reconcile` | `intake_pause` | `temporarily_unavailable` | `none` | `planned` | Provider kaynaklı geri bildirim geçici olarak durdurulabilir |
| `feature.ads.click_intake.enabled` | Ads click intake yüzeyini açma-kapama anahtarı | `operations` | `ads` | `service` | `bool` | `true` | `true,false` | `feature` | `ads.click_intake` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `ads.click`,`ads.aggregate` | `intake_pause` | `temporarily_unavailable` | `none` | `planned` | Invalid click koruması için click intake ayrı durdurulabilir |
| `feature.support.internal_note.enabled` | Support internal note yazma yüzeyini açma-kapama anahtarı | `support` | `support` | `access` | `bool` | `true` | `true,false` | `feature` | `support.internal_note` | `authenticated` | `role:support_agent` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.internal_note` | `write_off` | `forbidden` | `none` | `planned` | Internal note ile requester'a açık reply yüzeyi ayrı yönetilmelidir |

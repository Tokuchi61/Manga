# Shared ve Settings

> Bu dosya `docs/shared/*` ve `docs/settings/index.md` iÃ§eriklerinin tek dosyada birleÅŸtirilmiÅŸ halidir. ModÃ¼ller Ã¼stÃ¼ teknik kararlar, ortak sÃ¶zlÃ¼kler, policy belgeleri ve runtime settings envanteri burada tutulur.

## KullanÄ±m KurallarÄ±
- Teknik stack, enum, policy, precedence, transaction ve operasyon kurallarÄ± burada canonical referanstÄ±r.
- Secret/config kararlarÄ± ile runtime setting kararlarÄ± birbirine karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.
- Yeni shared enum, ortak policy veya runtime anahtarÄ± eklendiÄŸinde bu dosya aynÄ± deÄŸiÅŸiklikte gÃ¼ncellenmelidir.
- ModÃ¼l dokÃ¼manlarÄ± bu dosyadaki shared kararlarla Ã§eliÅŸmemelidir.




---

# Audit Event Tipleri

> Audit kayÄ±tlarÄ±nda `action` alanÄ±nÄ± sÄ±nÄ±flandÄ±rmak veya kategori bazlÄ± filtrelemek gerektiÄŸinde bu sÃ¶zlÃ¼k kullanÄ±lmalÄ±dÄ±r.

## Canonical KayÄ±tlar
| DeÄŸer | AÃ§Ä±klama |
| --- | --- |
| `security_auth` | Kimlik doÄŸrulama ve oturum gÃ¼venliÄŸi olaylarÄ± |
| `access_policy` | Policy deÄŸiÅŸimi, deny override ve availability yorumu |
| `admin_action` | YÃ¶netimsel kritik aksiyon ve settings deÄŸiÅŸikliÄŸi |
| `moderation_action` | Case inceleme, escalation ve iÃ§erik aksiyonu |
| `support_case` | Ticket, report, Ã§Ã¶zÃ¼m ve internal note olaylarÄ± |
| `payment_financial` | Checkout, callback, refund, reversal ve reconcile |
| `inventory_change` | Grant, revoke, consume, equip ve dÃ¼zeltme |
| `shop_purchase` | Purchase intent, recovery ve eligibility override |
| `user_state` | Ban, deactivate, visibility veya VIP state deÄŸiÅŸimi |
| `notification_ops` | Template, category, suppression ve delivery failure eÅŸiÄŸi |
| `ads_ops` | Campaign publish/pause, click protection ve aggregate mÃ¼dahalesi |
| `system_ops` | Backup, restore, migration veya genel operasyon aksiyonlarÄ± |


---

# Audit PolitikasÄ±

> GÃ¼venlik, yÃ¶netim, finansal doÄŸruluk ve kritik state transition'lar iÃ§in ortak audit formatÄ± bu dokÃ¼manda tanÄ±mlanÄ±r.

## AmaÃ§
- Audit ihtiyacÄ±nÄ± modÃ¼ller arasÄ±nda ortak alan seti ve ortak event sÄ±nÄ±flarÄ±yla yÃ¶netmek.
- Admin notu, operasyonel not ve immutable audit kaydÄ±nÄ± birbirine karÄ±ÅŸtÄ±rmamak.
- GÃ¼venlik ve finansal olaylarda minimum izlenebilirlik standardÄ± saÄŸlamak.

## Zorunlu Alanlar
| Alan | AÃ§Ä±klama |
| --- | --- |
| `occurred_at` | OlayÄ±n kesin oluÅŸma zamanÄ± |
| `source_module` | Audit kaydÄ±nÄ± Ã¼reten canonical modÃ¼l |
| `actor_type` / `actor_id` | Ä°ÅŸlemi baÅŸlatan taraf |
| `target_type` / `target_id` | Etkilenen kullanÄ±cÄ±, iÃ§erik veya iÅŸlem hedefi |
| `action` | YapÄ±lan aksiyonun canonical adÄ± |
| `result` | `success`, `rejected`, `failed`, `reversed` gibi sonuÃ§ |
| `reason` | Ä°nsan okunabilir kÄ±sa neden veya reason code |
| `request_id` | Gelen isteÄŸin benzersiz izi |
| `correlation_id` | AkÄ±ÅŸ boyunca taÅŸÄ±nan iliÅŸki kimliÄŸi |
| `risk_level` | `low`, `medium`, `high`, `critical` |
| `metadata_redaction_level` | PII masking seviyesini anlatan iÅŸaret |

## Audit Gerektiren ModÃ¼l AlanlarÄ±
| ModÃ¼l | Minimum Audit OlaylarÄ± |
| --- | --- |
| `auth` | login failure, suspicious login, password reset, session revoke |
| `access` | policy override, deny reason deÄŸiÅŸikliÄŸi, emergency availability yorumu |
| `admin` | settings deÄŸiÅŸikliÄŸi, override, impersonation, destructive action |
| `moderation` | case action, escalation, reassignment, evidence snapshot eriÅŸimi |
| `support` | karar deÄŸiÅŸimi, internal note, linked moderation handoff |
| `payment` | checkout, callback, refund, reversal, reconcile |
| `inventory` | grant, revoke, consume, manual correction |
| `shop` | purchase orchestration, recovery, eligibility override |
| `user` | VIP state deÄŸiÅŸimi, ban/deactivation, gÃ¶rÃ¼nÃ¼rlÃ¼k deÄŸiÅŸimi |
| `ads` | campaign publish/pause, click protection mÃ¼dahalesi |
| `notification` | template/category deÄŸiÅŸimi, suppression, delivery failure eÅŸiÄŸi |

## Kurallar
- Admin, gÃ¼venlik ve finansal olaylar immutable log yaklaÅŸÄ±mÄ±yla tutulmalÄ±dÄ±r.
- Admin note veya support internal note, audit kaydÄ±nÄ±n yerine geÃ§memelidir; ayrÄ± veri modeli olarak kalmalÄ±dÄ±r.
- Audit export yÃ¼zeyi read-only olmalÄ± ve PII masking `docs/shared/operational-standards.md` ile hizalÄ± kalmalÄ±dÄ±r.
- Event sÄ±nÄ±flarÄ± iÃ§in `docs/shared/audit-event-types.md` iÃ§indeki canonical sÃ¶zlÃ¼k kullanÄ±lmalÄ±dÄ±r.


---

# Cache ve Queue Stratejisi

> Bu dokÃ¼man cache backend'i, async iÅŸleme baseline'Ä± ve broker'a geÃ§iÅŸ kriterlerini somut karar seviyesinde tanÄ±mlar.

## AmaÃ§
- Cache ve queue tarafÄ±nda â€œileride bakarÄ±zâ€ belirsizliÄŸini kaldÄ±rmak.
- Hangi teknolojinin bugÃ¼nkÃ¼ canonical seÃ§im olduÄŸunu netleÅŸtirmek.
- Async iÅŸleme ile source-of-truth sÄ±nÄ±rÄ±nÄ± ayÄ±rmak.

## Cache KararÄ±
- Cache zorunlu baseline deÄŸildir; source-of-truth her zaman owner modÃ¼lÃ¼n veritabanÄ±dÄ±r.
- Cache ihtiyacÄ± oluÅŸtuÄŸunda canonical backend `Redis` olmalÄ±dÄ±r.
- Ä°lk cache adaylarÄ± ÅŸunlardÄ±r:
- access decision cache
- manga listing veya discovery cache
- notification unread counter cache
- ads placement resolve cache

## Cache KurallarÄ±
- Cache yokluÄŸu business doÄŸruluÄŸunu bozmamalÄ±dÄ±r; yalnÄ±zca performans etkisi yaratmalÄ±dÄ±r.
- Cache key'leri subject, surface, selector ve gerekiyorsa settings sÃ¼rÃ¼mÃ¼ ile iliÅŸkilendirilmelidir.
- TTL ve invalidation davranÄ±ÅŸÄ± modÃ¼l dokÃ¼manÄ± ile `docs/settings/index.md` notlarÄ± iÃ§inde gÃ¶rÃ¼nÃ¼r tutulmalÄ±dÄ±r.

## Queue / Async Ä°ÅŸleme KararÄ±
- BugÃ¼nkÃ¼ canonical async iÅŸleme baseline'Ä± PostgreSQL-backed jobs + transactional outbox yaklaÅŸÄ±mÄ±dÄ±r.
- Message broker baseline mimarinin parÃ§asÄ± deÄŸildir.
- Outbox, retry, backoff ve dead-letter kurallarÄ± `docs/shared/outbox-pattern.md` ile birlikte deÄŸerlendirilmelidir.

## Ä°lk Async Ä°ÅŸ YÃ¼kleri
- notification delivery
- projection rebuild
- payment reconciliation
- ads aggregate jobs
- mission reset
- royalpass season jobs

## Broker'a GeÃ§iÅŸ Kriterleri
AÅŸaÄŸÄ±daki koÅŸullarÄ±n birden fazlasÄ± oluÅŸmadan ayrÄ± broker seÃ§imine gidilmemelidir:
- modÃ¼ller arasÄ± yÃ¼ksek hacimli fan-out ihtiyacÄ±
- aynÄ± event iÃ§in baÄŸÄ±msÄ±z consumer gruplarÄ±nda kalÄ±cÄ± lag
- DB-backed job yaklaÅŸÄ±mÄ±nda operasyonel gÃ¶zlem veya retry maliyetinin aÅŸÄ±rÄ± yÃ¼kselmesi
- dÄ±ÅŸ sistemler veya ayrÄ± servislerle yoÄŸun Ã§ift yÃ¶nlÃ¼ event entegrasyonu

## Redis'in RolÃ¼
- Redis cache backend'i olarak canonical tercihtir.
- Redis, aÃ§Ä±k bir mimari gÃ¼ncelleme yapÄ±lmadÄ±kÃ§a source-of-truth queue olarak kullanÄ±lmamalÄ±dÄ±r.
- Queue hÄ±zlandÄ±rma veya ephemeral yardÄ±mcÄ± akÄ±ÅŸlar ancak owner state yine DB'de kaldÄ±ÄŸÄ± sÃ¼rece dÃ¼ÅŸÃ¼nÃ¼lebilir.


---

# Ä°dempotency PolitikasÄ±

> Tekrarlanan istek, callback, claim ve grant akÄ±ÅŸlarÄ±nda yinelenen yan etkileri Ã¶nlemek iÃ§in canonical kurallar bu dokÃ¼manda tutulur.

## AmaÃ§
- AynÄ± isteÄŸin yeniden gelmesi durumunda tekrar Ã¶deme, tekrar grant veya tekrar state transition oluÅŸmasÄ±nÄ± engellemek.
- Producer ve consumer tarafÄ±nda safe retry / unsafe retry ayrÄ±mÄ±nÄ± gÃ¶rÃ¼nÃ¼r kÄ±lmak.
- ModÃ¼l bazlÄ± Ã¶zel uygulamalarÄ± tek Ã§erÃ§eveye baÄŸlamak.

## Key FormatÄ±
- VarsayÄ±lan key biÃ§imi `module:operation:actor_or_scope:client_request_id` olmalÄ±dÄ±r.
- DÄ±ÅŸ provider callback'lerinde client request yerine provider event veya provider transaction referansÄ± kullanÄ±labilir.
- Key ile birlikte payload hash, ilk sonuÃ§ durumu ve yan etki referansÄ± saklanmalÄ±dÄ±r.

## Kapsam ve Saklama
| Kapsam | Ã–rnek Ä°ÅŸlemler | TTL | Saklama Notu |
| --- | --- | --- | --- |
| HTTP write request | register, profile update, purchase intent | en az 24 saat | aynÄ± client request tekrarÄ±nda ilk sonuÃ§ dÃ¶ndÃ¼rÃ¼lÃ¼r |
| Callback / webhook | payment provider callback, external delivery callback | en az 7 gÃ¼n | provider event id bazlÄ± tutulur |
| Claim / grant | mission claim, royalpass claim, inventory grant | en az 30 gÃ¼n | reward veya ledger yan etkisi ile birlikte saklanÄ±r |
| Background consumer | queue consumer, outbox relay, projection updater | retry penceresi + gÃ¶zlem sÃ¼resi | event id veya message id bazlÄ± tutulur |

## Duplicate Handling
- AynÄ± key ve aynÄ± payload tekrar gelirse Ã¶nceki sonuÃ§ veya mevcut pending durum dÃ¶ndÃ¼rÃ¼lmelidir.
- AynÄ± key fakat farklÄ± payload gelirse istek reddedilmeli ve audit kaydÄ± aÃ§Ä±lmalÄ±dÄ±r.
- In-flight duplicate istekler yeni yan etki Ã¼retmemeli; mÃ¼mkÃ¼nse mevcut pending iÅŸlem referansÄ±na baÄŸlanmalÄ±dÄ±r.
- Ä°dempotency kaydÄ± durable olarak yazÄ±lmadan yan etki baÅŸlatÄ±lmamalÄ±dÄ±r.

## Retry KurallarÄ±
- Safe retry: read model update, callback doÄŸrulama, claim uygunluÄŸu kontrolÃ¼, delivery retry, projection rebuild.
- Unsafe retry: yeni Ã¶deme baÅŸlatma, tekil reward grant, aynÄ± entitlement'Ä± tekrar aktive etme, destructive admin action.
- Unsafe retry gerektiren akÄ±ÅŸlarda ya idempotency key zorunlu tutulmalÄ± ya da manuel review kapÄ±sÄ± aÃ§Ä±lmalÄ±dÄ±r.

## ModÃ¼l UygulamalarÄ±
| ModÃ¼l | Zorunlu Ä°ÅŸlemler | Key Temeli | Not |
| --- | --- | --- | --- |
| `payment` | checkout callback, refund, reversal, reconcile write | provider event id veya checkout request id | ledger-first yaklaÅŸÄ±m ile birlikte zorunludur |
| `inventory` | grant, consume, equip | source reference + actor + request id | duplicate grant korumasÄ± taÅŸÄ±r |
| `mission` | claim request, progress ingest | user + mission + period + request id | period reset ile birlikte dÃ¼ÅŸÃ¼nÃ¼lmelidir |
| `royalpass` | tier claim, premium activation intake | user + season + tier veya activation ref | cross-season Ã§akÄ±ÅŸma Ã¶nlenmelidir |
| `history` | checkpoint write, merge | user + manga/chapter + device/request | multi-device merge ile birlikte uygulanÄ±r |
| `notification` | dedup send, delivery retry | notification id veya producer event ref | category/channel bazlÄ± dedup gerekir |
| `support` | ticket/report create, linked case handoff | requester + target + request id | duplicate spam/report korumasÄ± iÃ§in kullanÄ±lÄ±r |
| `shop` | purchase intent, recovery replay | user + offer/product + request id | already-owned ve recovery senaryolarÄ± ile birlikte Ã§alÄ±ÅŸÄ±r |


---

# Media ve Asset Stratejisi

> Media/asset ihtiyacÄ± yalnÄ±zca not seviyesinde bÄ±rakÄ±lmamalÄ±; bu dokÃ¼man medya dosyalarÄ±, attachment akÄ±ÅŸlarÄ± ve signed access politikalarÄ± iÃ§in canonical yÃ¶nÃ¼ belirlemelidir.

## AmaÃ§
- Chapter, user, manga ve support gibi alanlardaki medya ihtiyacÄ±nÄ± ortak kuralla yÃ¶netmek.
- Teknik medya yÃ¶netimi ile business relation owner'lÄ±ÄŸÄ±nÄ± ayÄ±rmak.
- AyrÄ± bir `media` modÃ¼lÃ¼nÃ¼n ne zaman aÃ§Ä±lacaÄŸÄ±nÄ± netleÅŸtirmek.

## BugÃ¼nkÃ¼ Karar
- Hemen ayrÄ± bir leaf `media` modÃ¼lÃ¼ aÃ§Ä±lmayacaktÄ±r.
- Canonical yaklaÅŸÄ±m, ortak media/asset altyapÄ±sÄ±nÄ±n teknik olarak shared platform katmanÄ±nda planlanmasÄ±dÄ±r.
- Domain owner modÃ¼l, media iliÅŸkisini business referans olarak taÅŸÄ±r; teknik metadata ve eriÅŸim politikasÄ± ortak media altyapÄ±sÄ± ile yÃ¶netilir.

## Kapsam
- upload metadata
- ownership binding
- mime/type validation
- image dimension metadata
- attachment relation
- signed access policies

## TÃ¼ketici Alanlar
- `user`: avatar ve banner
- `manga`: cover, poster veya gÃ¶rsel metadata
- `chapter`: page media ve signed access ihtiyacÄ±
- `support`: attachment intake ve scanning iliÅŸkisi

## Kurallar
- Dosya sahipliÄŸi business modÃ¼lÃ¼n relation verisi olarak kalmalÄ±dÄ±r; blob storage eriÅŸimi ortak teknik katmanda ele alÄ±nmalÄ±dÄ±r.
- MIME, boyut, Ã¶lÃ§Ã¼ ve gerekiyorsa zararlÄ± iÃ§erik taramasÄ± upload sÄ±nÄ±rÄ±nda uygulanmalÄ±dÄ±r.
- Support attachment yÃ¼zeyi ile chapter page media aynÄ± validasyon seviyesini taÅŸÄ±mak zorunda deÄŸildir; ortak altyapÄ± farklÄ± policy profilleri taÅŸÄ±yabilmelidir.
- Signed URL veya Ã¶zel eriÅŸim gerekiyorsa token lifetime ve audit beklentisi aÃ§Ä±kÃ§a dokÃ¼mante edilmelidir.

## AyrÄ± ModÃ¼le AyrÄ±ÅŸma Kriterleri
AÅŸaÄŸÄ±daki koÅŸullarÄ±n birden fazlasÄ± oluÅŸmadan ayrÄ± leaf `media` modÃ¼lÃ¼ne geÃ§ilmemelidir:
- birden Ã§ok modÃ¼lÃ¼n aynÄ± upload, transform ve delivery pipeline'Ä±nÄ± yoÄŸun biÃ§imde paylaÅŸmasÄ±
- media lifecycle'Ä±nÄ±n business modÃ¼llerden baÄŸÄ±msÄ±z operasyon yÃ¼kÃ¼ Ã¼retmesi
- storage, scanning, variant generation ve signed access kurallarÄ±nÄ±n tek baÅŸÄ±na ayrÄ± ekip veya ayrÄ± servis sÄ±nÄ±rÄ± gerektirmesi

## Ä°lgili Referanslar
- Runtime ayar envanteri: `docs/settings/index.md`
- Operasyonel standartlar: `docs/shared/operational-standards.md`
- Cache ve queue stratejisi: `docs/shared/cache-queue-strategy.md`


---

# Moderation DurumlarÄ±

> `moderation` modÃ¼lÃ¼ndeki lifecycle alanlarÄ± iÃ§in canonical sÃ¶zlÃ¼k.

## Case Status
| DeÄŸer | AnlamÄ± |
| --- | --- |
| `new` | Yeni aÃ§Ä±lmÄ±ÅŸ ve henÃ¼z triage edilmemiÅŸ vaka |
| `queued` | Ä°nceleme kuyruÄŸunda bekleyen vaka |
| `assigned` | ModeratÃ¶re atanmÄ±ÅŸ vaka |
| `in_review` | Aktif inceleme altÄ±nda |
| `escalated` | Admin veya daha yÃ¼ksek scope'a yÃ¼kseltilmiÅŸ |
| `resolved` | Karar verilmiÅŸ ve kapanÄ±ÅŸa uygun |
| `rejected` | GeÃ§ersiz veya yanlÄ±ÅŸ intake olarak reddedilmiÅŸ |
| `closed` | Finalize edilip kapatÄ±lmÄ±ÅŸ |

## Assignment Status
| DeÄŸer | AnlamÄ± |
| --- | --- |
| `unassigned` | AtanmamÄ±ÅŸ |
| `assigned` | Bir moderatÃ¶re atanmÄ±ÅŸ |
| `handoff_pending` | Scope veya kiÅŸi deÄŸiÅŸimi bekliyor |
| `released` | Atama boÅŸaltÄ±lmÄ±ÅŸ veya geri alÄ±nmÄ±ÅŸ |

## Action Result
| DeÄŸer | AnlamÄ± |
| --- | --- |
| `none` | HenÃ¼z aksiyon uygulanmadÄ± |
| `content_hidden` | Ä°Ã§erik gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ kapatÄ±ldÄ± |
| `content_restored` | Ã–nceki aksiyon geri alÄ±ndÄ± |
| `warning_sent` | KullanÄ±cÄ±ya uyarÄ± uygulandÄ± |
| `no_action` | Ä°nceleme sonrasÄ± aksiyon gerekmedi |


---

# Notification Kategorileri

> KullanÄ±cÄ±ya dÃ¶nÃ¼k bildirim sÄ±nÄ±flarÄ± iÃ§in canonical sÃ¶zlÃ¼k.

## Canonical KayÄ±tlar
| DeÄŸer | AÃ§Ä±klama |
| --- | --- |
| `account_security` | Login, verification, password reset ve gÃ¼venlik uyarÄ±larÄ± |
| `social` | ArkadaÅŸlÄ±k, takip, direct message ve duvar olaylarÄ± |
| `comment` | Ä°Ã§erik yorumu ve thread etkileÅŸimleri |
| `support` | Ticket yanÄ±tÄ±, Ã§Ã¶zÃ¼m ve support gÃ¼ncellemeleri |
| `moderation` | KullanÄ±cÄ±yÄ± etkileyen moderasyon ve inceleme bildirimleri |
| `mission` | GÃ¶rev ilerleme, tamamlanma ve claim uygunluÄŸu |
| `royalpass` | Season baÅŸlangÄ±cÄ±, tier aÃ§Ä±lmasÄ± ve reward claim |
| `shop` | Offer, katalog ve satÄ±n alma sÃ¼reci bilgilendirmesi |
| `payment` | Checkout, Ã¶deme sonucu, refund ve cÃ¼zdan hareketi |
| `system_ops` | BakÄ±m, kill switch veya Ã¼rÃ¼n genel duyurularÄ± |


---

# Operasyonel Standartlar

> Request context, rate limit, secret yÃ¶netimi, backup ve veri saklama gibi modÃ¼ller Ã¼stÃ¼ iÅŸletim kararlarÄ± bu dokÃ¼manda tutulur.

## Request ID ve Correlation ID
- Her dÄ±ÅŸ istek en az bir `request_id` taÅŸÄ±malÄ± veya giriÅŸ katmanÄ±nda Ã¼retilmelidir.
- Asenkron zincirde aynÄ± iÅŸ akÄ±ÅŸÄ±nÄ± baÄŸlayan `correlation_id` korunmalÄ±dÄ±r.
- DÄ±ÅŸ provider callback'lerinde provider event referansÄ± ayrÄ±ca saklanmalÄ±; ancak iÃ§ korelasyon alanÄ±nÄ±n yerine geÃ§memelidir.

## Rate Limit PolitikasÄ±
| Surface | BaÅŸlangÄ±Ã§ Beklentisi |
| --- | --- |
| `auth.login` | baÅŸarÄ±sÄ±z giriÅŸ limiti ve cooldown |
| `comment.write` | flood ve anti-spam korumasÄ± |
| `support.intake` | duplicate/spam ve attachment korumasÄ± |
| `social.messaging` | anti-abuse ve burst kontrolÃ¼ |
| `payment.callback` | provider callback replay korumasÄ± |
| `ads.click_intake` | invalid click ve bot korumasÄ± |

## Secret ve Config AyrÄ±mÄ±
- Env secret, provider credential, signing secret veya private key admin runtime ayarÄ± olarak sunulmamalÄ±dÄ±r.
- Runtime settings kullanÄ±cÄ±ya dÃ¶nÃ¼k availability, threshold ve operasyon davranÄ±ÅŸÄ± iÃ§indir; secret config ile aynÄ± saklama dÃ¼zeyinde deÄŸerlendirilmez.
- Secret rotation ve audit beklentisi ayrÄ± gÃ¼venlik sÃ¼reci olarak ele alÄ±nmalÄ±dÄ±r.

## Backup, Restore ve Migration Rollback
- Destructive migration Ã¶ncesinde doÄŸrulanmÄ±ÅŸ backup planÄ± bulunmalÄ±dÄ±r.
- Rollback yolu olmayan migration veya veri dÃ¶nÃ¼ÅŸÃ¼mÃ¼ kritik alanlarda kabul edilmemelidir.
- Restore sonrasÄ± Ã¶zellikle `payment`, projection ve outbox tabanlÄ± modÃ¼ller iÃ§in reconcile adÄ±mÄ± planlanmalÄ±dÄ±r.

## PII ve Veri Saklama
| ModÃ¼l | PII / Hassas Veri SÄ±nÄ±fÄ± | Saklama / Maskeleme Notu |
| --- | --- | --- |
| `auth` | credential, gÃ¼venlik olayÄ±, IP/device izi | uzun sÃ¼reli gÃ¼venlik kaydÄ± hariÃ§ maskeli saklama tercih edilir |
| `user` | profil ve Ã¼yelik bilgisi | public/private alanlar aÃ§Ä±k ayrÄ±lmalÄ±, gereksiz kopya oluÅŸturulmamalÄ± |
| `support` | talep metni, attachment, kiÅŸisel aÃ§Ä±klamalar | internal note ve public reply ayrÄ±mÄ± zorunlu olmalÄ±dÄ±r |
| `social` | direct message ve iliÅŸki sinyalleri | export ve moderasyon yÃ¼zeylerinde minimizasyon uygulanmalÄ±dÄ±r |
| `payment` | finansal referans ve provider metadata | finansal doÄŸruluk korunurken kullanÄ±cÄ±ya dÃ¶nen yÃ¼zey maskeli olmalÄ±dÄ±r |

## Ek Notlar
- ModÃ¼l dokÃ¼manlarÄ± veri retention, masking ve export risklerini kendi alanlarÄ±nda ayrÄ±ca belirtmelidir.
- Audit kayÄ±tlarÄ± PII taÅŸÄ±yorsa `metadata_redaction_level` alanÄ± zorunlu olmalÄ±dÄ±r.


---

# Outbox Deseni

> ModÃ¼ller arasÄ± event Ã¼retimi arttÄ±ÄŸÄ±nda gÃ¼venli publish stratejisinin minimum omurgasÄ± bu dokÃ¼manda tutulur.

## AmaÃ§
- â€œDB yaz + hemen publishâ€ yarÄ±ÅŸ durumlarÄ±nÄ± Ã¶nlemek.
- Event publish baÅŸarÄ±sÄ±zlÄ±klarÄ±nda tekrar deneme ve dead-letter akÄ±ÅŸÄ±nÄ± standartlaÅŸtÄ±rmak.
- Replay ve projection rebuild ihtiyacÄ±na temel hazÄ±rlamak.

## Zorunlu BileÅŸenler
- Transaction ile birlikte yazÄ±lan bir outbox kaydÄ±
- Durumu izlenebilir background publisher
- Retry ve backoff politikasÄ±
- Poison message veya kalÄ±cÄ± baÅŸarÄ±sÄ±zlÄ±k iÃ§in dead-letter yaklaÅŸÄ±mÄ±
- Publish lag, retry count ve failure rate iÃ§in gÃ¶zlem metrikleri

## Ã–ncelikli ModÃ¼ller
- `payment`
- `inventory`
- `mission`
- `royalpass`
- `notification`
- `support`
- `moderation`
- `history`

## Mesaj KurallarÄ±
- Her mesaj benzersiz `event_id`, `schema_version`, `request_id`, `correlation_id` ve mÃ¼mkÃ¼nse `causation_id` taÅŸÄ±malÄ±dÄ±r.
- Consumer tarafÄ± idempotent Ã§alÄ±ÅŸacak biÃ§imde tasarlanmalÄ±dÄ±r; publish garantisi consumer tekrarÄ±ndan baÄŸÄ±msÄ±z dÃ¼ÅŸÃ¼nÃ¼lmemelidir.
- Outbox payload'Ä± yalnÄ±zca consumer'Ä±n ihtiyaÃ§ duyduÄŸu kontrat alanlarÄ±nÄ± taÅŸÄ±malÄ±, producer iÃ§ implementasyonunu sÄ±zdÄ±rmamalÄ±dÄ±r.
- Dead-letter'a dÃ¼ÅŸen mesajlar sessizce yok sayÄ±lmamalÄ±; operasyonel gÃ¶rÃ¼nÃ¼m ve replay planÄ± taÅŸÄ±malÄ±dÄ±r.


---

# Policy Etkileri

> Authorization ve availability kararlarÄ±nÄ±n Ã§Ä±ktÄ±sÄ± olarak kullanÄ±lan canonical etki sÃ¶zlÃ¼ÄŸÃ¼.

## Canonical KayÄ±tlar
| DeÄŸer | AÃ§Ä±klama |
| --- | --- |
| `allow` | Ä°stek mevcut koÅŸullarda izinli |
| `deny` | Ä°stek kesin olarak reddedilir |
| `deny_soft` | Ä°stek doÄŸrudan reddedilmez; teslim/gÃ¶rÃ¼nÃ¼rlÃ¼k azaltÄ±lÄ±r veya degrade edilir |
| `require_auth` | EriÅŸim iÃ§in kimlik doÄŸrulama gerekir |
| `require_role` | Belirli rol veya scope zorunludur |
| `require_entitlement` | VIP, premium veya baÅŸka entitlement gerekir |
| `read_only` | Okuma aÃ§Ä±k kalÄ±rken yazma kapatÄ±lÄ±r |
| `write_off` | Ä°lgili yazma yÃ¼zeyi kapalÄ±dÄ±r |
| `mask` | Veri dÃ¶nebilir ancak alan bazlÄ± maskeleme uygulanÄ±r |
| `needs_review` | Otomatik karar verilmez; manuel inceleme gerekir |


---

# Precedence KurallarÄ±

> Ã‡apraz modÃ¼l karar Ã§atÄ±ÅŸmalarÄ±nda hangi sinyalin baskÄ±n geleceÄŸini bu dokÃ¼man belirler. ModÃ¼l dokÃ¼manlarÄ± aynÄ± Ã§atÄ±ÅŸmayÄ± farklÄ± kelimelerle yeniden tanÄ±mlamamalÄ±; burada yazÄ±lan kurallarÄ± referans almalÄ±dÄ±r.

## AmaÃ§
- Ã‡akÄ±ÅŸan gÃ¶rÃ¼nÃ¼rlÃ¼k, eriÅŸim, moderasyon ve operasyon kararlarÄ±nÄ± tek matris altÄ±nda toplamak.
- `access`, `admin`, `moderation`, `support`, `history`, `social` ve entitlement akÄ±ÅŸlarÄ± arasÄ±nda ortak yorum sÄ±rasÄ± oluÅŸturmak.
- Kod tarafÄ±nda sezgisel kalan â€œhangisi baskÄ±n?â€ sorusunu dokÃ¼mante etmek.

## Karar SÄ±rasÄ±
1. Sistem genelindeki bakÄ±m veya acil durdurma kararÄ± deÄŸerlendirilir.
2. Ä°lgili modÃ¼l veya surface iÃ§in runtime kill switch ve availability anahtarÄ± deÄŸerlendirilir.
3. Security, moderation veya admin override kaynaklÄ± explicit deny/override kararlarÄ± uygulanÄ±r.
4. Audience, role, entitlement ve ownership tabanlÄ± access kararÄ± yorumlanÄ±r.
5. ModÃ¼lÃ¼n kendi visibility/share metadata'sÄ± yalnÄ±zca kalan izinli aralÄ±k iÃ§inde uygulanÄ±r.
6. Rate limit, cooldown veya backpressure kaynaklÄ± geÃ§ici reddetme son aÅŸamada uygulanÄ±r.

## Canonical Kurallar
| Ã‡akÄ±ÅŸma | Kazanan | GerekÃ§e |
| --- | --- | --- |
| `system deny` vs `entitlement allow` | `system deny` | BakÄ±m, gÃ¼venlik veya operasyonel kapatma Ã¼cretli haklardan daha Ã¼st precedence taÅŸÄ±r. |
| `runtime kill switch` vs `normal access allow` | `runtime kill switch` | Surface kapalÄ±ysa permission tek baÅŸÄ±na eriÅŸim aÃ§amaz. |
| `admin hard override` vs `scoped moderator action` | `admin hard override` | GÃ¼nlÃ¼k vaka sahipliÄŸi `moderation` iÃ§inde kalsa da final yÃ¶netimsel karar `admin` tarafÄ±nda baskÄ±ndÄ±r. |
| `moderation block` vs `social visibility allow` | `moderation block` | Ä°Ã§erik veya kullanÄ±cÄ± gÃ¼venliÄŸi kaynaklÄ± blok kararÄ± sosyal gÃ¶rÃ¼nÃ¼rlÃ¼kten Ã¼stÃ¼ndÃ¼r. |
| `social block` vs `messaging/wall allow` | `social block` | AÃ§Ä±k block iliÅŸkisi direct message ve wall etkileÅŸimini durdurur. |
| `social mute` vs `authorization allow` | `authorization allow` | `mute` varsayÄ±lan olarak teslim/gÃ¶rÃ¼nÃ¼rlÃ¼k azaltma sinyalidir; tek baÅŸÄ±na final authorization deny sayÄ±lmaz. |
| `user global visibility deny` vs `history entry share opt-in` | `user global visibility deny` | `history` entry-level paylaÅŸÄ±m kararÄ± global deny tavanÄ±nÄ± aÅŸamaz. |
| `support report` vs `moderation case` | AyrÄ± kayÄ±tlar | Report intake otomatik olarak moderation case sayÄ±lmaz; aÃ§Ä±k mapping politikasÄ± gerekir. |
| `vip no-ads` veya baÅŸka entitlement muafiyeti vs `ads/payment/support` kill switch | Kill switch | Sistem gÃ¼venliÄŸi ve operasyon kararÄ± Ã¼rÃ¼n avantajÄ±nÄ±n Ã¼stÃ¼nde deÄŸerlendirilir. |
| `provider callback retry` vs `payment callback intake pause` | Intake pause | Provider tekrar denemesi aÃ§Ä±k olsa bile operasyonda intake geÃ§ici durdurulabilir. |

## Uygulama NotlarÄ±
- `access` modÃ¼lÃ¼ kullanÄ±cÄ±ya bakan availability ve permission kararlarÄ±nÄ±n son yorumlayÄ±cÄ±sÄ±dÄ±r.
- Operasyonel pause, callback intake veya queue backpressure gibi kullanÄ±cÄ±ya doÄŸrudan gÃ¶rÃ¼nmeyen kontroller ilgili `service` katmanÄ± tarafÄ±ndan yorumlanabilir.
- AynÄ± precedence kuralÄ± hem modÃ¼l dokÃ¼manÄ±nda hem settings envanterinde farklÄ± ifadelerle yazÄ±lmamalÄ±dÄ±r.


---

# Projection Stratejisi

> Event veya change-feed ile beslenen read model, counter ve denormalize Ã¶zet yÃ¼zeyleri bu dokÃ¼manla hizalÄ± olmalÄ±dÄ±r.

## AmaÃ§
- Canonical write model ile projection sorumluluÄŸunu ayÄ±rmak.
- Eventual consistency beklentisini ve rebuild yolunu gÃ¶rÃ¼nÃ¼r kÄ±lmak.
- Counter, summary ve read model bÃ¼yÃ¼mesini modÃ¼ller arasÄ±nda tutarlÄ± hale getirmek.

## Temel Ä°lkeler
- Her projection iÃ§in canonical write model owner modÃ¼lde kalÄ±r; consumer modÃ¼l owner tabloya doÄŸrudan yazmaz.
- Projection gÃ¼ncellemeleri mÃ¼mkÃ¼n olduÄŸunda event, outbox veya aÃ§Ä±k projection contract yÃ¼zeyi ile beslenir.
- Kabul edilen eventual consistency penceresi dokÃ¼mante edilmeden denormalize alan aÃ§Ä±lmaz.
- Projection rebuild ve replay akÄ±ÅŸÄ± en baÅŸtan planlanÄ±r; â€œyalnÄ±zca incremental Ã§alÄ±ÅŸÄ±râ€ kabulÃ¼ yapÄ±lmaz.
- Projection hatalarÄ± ve lag durumu izlenebilir metrik Ã¼retebilmelidir.

## Canonical Projection KayÄ±tlarÄ±
| Projection | Canonical Write Model | Event KaynaÄŸÄ± | Consumer Surface | TutarlÄ±lÄ±k Penceresi | Rebuild Yolu | Replay |
| --- | --- | --- | --- | --- | --- | --- |
| `manga.comment_count` | `comment` | `comment.created`, `comment.deleted`, `comment.moderated` | manga detail ve listing | kÄ±sa | comment tablosundan recount + incremental catch-up | desteklenir |
| `manga.engagement_summary` | `history`, `comment` | read checkpoint ve engagement event'leri | discovery ve admin Ã¶zetleri | orta | gÃ¼nlÃ¼k batch rebuild + hedefli repair | desteklenir |
| `history.continue_reading_projection` | `history` | checkpoint ve finish event'leri | continue reading yÃ¼zeyi | kÄ±sa | son checkpoint'ten recompute | desteklenir |
| `notification.unread_counter` | `notification` | `notification.created`, `notification.read` | inbox badge ve header sayaÃ§larÄ± | kÄ±sa | unread kayÄ±tlardan recount | desteklenir |
| `support.queue_summary` | `support` | create, status_change, assignee_change | support operasyon paneli | kÄ±sa | status bazlÄ± regroup | desteklenir |
| `moderation.queue_summary` | `moderation` | case create, assignment, resolution | moderation paneli | kÄ±sa | case durumlarÄ±ndan regroup | desteklenir |
| `ads.impression_aggregate` | `ads` | accepted impression ve click event'leri | reporting ve dashboard yÃ¼zeyi | orta | batch aggregation job | desteklenir |
| `mission.progress_projection` | `mission` | progress, claim, reset event'leri | mission liste ve progress Ã¶zetleri | kÄ±sa | objective bazlÄ± recompute | desteklenir |
| `royalpass.tier_progress_snapshot` | `royalpass` | progress ve claim event'leri | season overview ve tier gÃ¶rÃ¼nÃ¼mÃ¼ | kÄ±sa | tier bazlÄ± recompute | desteklenir |

## Uygulama KurallarÄ±
- Projection rebuild iÅŸlemi idempotent olmalÄ± ve aynÄ± veri iÃ§in tekrar Ã§alÄ±ÅŸtÄ±rÄ±ldÄ±ÄŸÄ±nda yeni yan etki Ã¼retmemelidir.
- Replay yapÄ±lacak event payload'larÄ± ÅŸema sÃ¼rÃ¼mÃ¼, `request_id` ve `correlation_id` taÅŸÄ±malÄ±dÄ±r.
- Counter veya summary gÃ¼ncellemeleri iÃ§in doÄŸrudan â€œowner tabloya baÅŸka modÃ¼l yazsÄ±nâ€ yaklaÅŸÄ±mÄ± kullanÄ±lmamalÄ±dÄ±r.
- Event Ã¼reten modÃ¼ller `docs/shared/outbox-pattern.md` ile hizalÄ± transactional outbox yaklaÅŸÄ±mÄ± planlamalÄ±dÄ±r.


---

# Purchase Source Tipleri

> SatÄ±n alma, checkout veya fulfillment baÅŸlatan canonical kaynak sÃ¶zlÃ¼ÄŸÃ¼.

## Canonical KayÄ±tlar
| DeÄŸer | Tipik Owner | AÃ§Ä±klama |
| --- | --- | --- |
| `catalog_purchase` | `shop` | Normal Ã¼rÃ¼n veya offer satÄ±n alma isteÄŸi |
| `premium_activation` | `shop`, `royalpass` | Premium pass veya entitlement aktivasyon akÄ±ÅŸÄ± |
| `mana_wallet` | `payment` | Ä°Ã§ bakiye kullanÄ±larak yapÄ±lan harcama |
| `external_provider` | `payment` | Harici Ã¶deme saÄŸlayÄ±cÄ±sÄ± ile baÅŸlayan akÄ±ÅŸ |
| `recovery_replay` | `shop`, `payment` | BaÅŸarÄ±sÄ±z veya kesilmiÅŸ akÄ±ÅŸÄ±n gÃ¼venli tekrar yÃ¼rÃ¼tÃ¼mÃ¼ |
| `admin_issue` | `admin` | Manuel operasyon kaynaÄŸÄ± |
| `gift_code` | gelecekte ayrÄ± modÃ¼l | Kod veya kupon bazlÄ± aktivasyon |


---

# Reporting ve Analytics Stratejisi

> Dashboard ve rapor ihtiyacÄ± yalnÄ±zca modÃ¼l iÃ§ine daÄŸÄ±lmÄ±ÅŸ notlarla yÃ¶netilmemelidir. Bu dokÃ¼man reporting/analytics read model yaklaÅŸÄ±mÄ±nÄ±, export sÄ±nÄ±rlarÄ±nÄ± ve operasyon Ã¶zetlerini canonical karara dÃ¶nÃ¼ÅŸtÃ¼rÃ¼r.

## AmaÃ§
- Admin, ads, payment, support ve benzeri alanlarÄ±n raporlama ihtiyacÄ±nÄ± ortak bir modelde toplamak.
- Operasyon summary ile analytics aggregate katmanÄ±nÄ± ayÄ±rmak.
- Export-friendly query yÃ¼zeylerini write model ownership'ini bozmadan tanÄ±mlamak.

## BugÃ¼nkÃ¼ Karar
- AyrÄ± bir analytics write modÃ¼lÃ¼ baseline deÄŸildir.
- Reporting yaklaÅŸÄ±mÄ± projection, aggregate read model ve kontrollÃ¼ export query layer Ã¼stÃ¼nden kurulmalÄ±dÄ±r.
- Write model owner'lÄ±ÄŸÄ± ilgili iÅŸ modÃ¼lÃ¼nde kalÄ±r; reporting yalnÄ±zca okuma odaklÄ± gÃ¶rÃ¼nÃ¼m Ã¼retir.

## Katmanlar
- Operasyon summary: admin ve ekiplerin anlÄ±k durum gÃ¶rmesi iÃ§in dÃ¼ÅŸÃ¼k gecikmeli Ã¶zetler
- Analytics aggregate: trend, hacim, funnel veya performans yorumlarÄ± iÃ§in periyodik projection'lar
- Export query layer: denetim, finans veya operasyon ihtiyaÃ§larÄ± iÃ§in kontrollÃ¼ dÄ±ÅŸa aktarma yÃ¼zeyi

## Kurallar
- Reporting projection'larÄ± canonical write model'in yerini almamalÄ±dÄ±r.
- Export yÃ¼zeyleri yalnÄ±zca yetkili operasyon akÄ±ÅŸlarÄ±yla aÃ§Ä±lmalÄ± ve audit zorunluluÄŸu taÅŸÄ±malÄ±dÄ±r.
- `payment` reconciliation, refund ve fraud review summary; `ads` performans aggregate; `support` queue summary gibi yÃ¼zeyler aynÄ± read-model prensipleriyle dokÃ¼mante edilmelidir.
- Dashboard metriÄŸi ile karar verici business state birbirine karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.

## Ne Zaman AyrÄ±laÅŸtÄ±rÄ±lÄ±r?
AÅŸaÄŸÄ±daki koÅŸullar belirginleÅŸmeden ayrÄ± analytics servisi veya veri hattÄ±na geÃ§ilmemelidir:
- operasyon summary ile analytics sorgularÄ±nÄ±n aynÄ± DB yÃ¼kÃ¼nde sÃ¼rdÃ¼rÃ¼lememesi
- bÃ¼yÃ¼k hacimli history ve event akÄ±ÅŸlarÄ±nÄ±n ayrÄ± depolama gerektirmesi
- ayrÄ± ekip, ayrÄ± eriÅŸim politikasÄ± veya ayrÄ± veri saklama kurallarÄ± ihtiyacÄ± oluÅŸmasÄ±

## Ä°lgili Referanslar
- Projection stratejisi: `docs/shared/projection-strategy.md`
- Audit politikasÄ±: `docs/shared/audit-policy.md`
- Cache ve queue stratejisi: `docs/shared/cache-queue-strategy.md`
- Admin modÃ¼lÃ¼: `docs/modules/admin.md`


---

# Reward Source Tipleri

> Ã–dÃ¼l, grant veya entitlement benzeri sahiplik akÄ±ÅŸlarÄ±nda canonical kaynak sÃ¶zlÃ¼ÄŸÃ¼.

## Canonical KayÄ±tlar
| DeÄŸer | Tipik Owner | AÃ§Ä±klama |
| --- | --- | --- |
| `mission` | `mission` | GÃ¶rev tamamlanmasÄ± veya claim sonucu oluÅŸan Ã¶dÃ¼l |
| `royalpass` | `royalpass` | Season tier veya pass claim sonucu oluÅŸan Ã¶dÃ¼l |
| `shop` | `shop` | SatÄ±n alma orkestrasyonu sonucu gelen grant |
| `admin_grant` | `admin` | Manuel operasyon veya destek amaÃ§lÄ± grant |
| `compensation` | `admin`, `support` | Sistem hatasÄ± telafisi |
| `seasonal_event` | `mission`, `royalpass` | DÃ¶nemsel event kampanyasÄ± kaynaÄŸÄ± |
| `referral` | gelecekte ayrÄ± modÃ¼l veya `user` | Referral benzeri teÅŸvik kaynaÄŸÄ± |
| `reconciliation_repair` | `payment`, `inventory` | Reconcile veya data repair sonucu dÃ¼zeltme grant'i |


---

# Search Stratejisi

> Search katmanÄ± modÃ¼l iÃ§i kÄ±sa not olarak bÄ±rakÄ±lmamalÄ±; baseline arama yaklaÅŸÄ±mÄ±, reindex sÃ¼reci ve provider swap kriterleri bu dokÃ¼manda tutulmalÄ±dÄ±r.

## AmaÃ§
- Search iÃ§in bugÃ¼nkÃ¼ canonical kararÄ± netleÅŸtirmek.
- Reindex, fallback ve provider deÄŸiÅŸimi konularÄ±nÄ± modÃ¼l notu olmaktan Ã§Ä±karÄ±p ortak kural haline getirmek.
- `manga` baÅŸta olmak Ã¼zere arama kullanan yÃ¼zeylerin aynÄ± omurgaya baÄŸlanmasÄ±nÄ± saÄŸlamak.

## BugÃ¼nkÃ¼ Karar
- BaÅŸlangÄ±Ã§ search engine kararÄ± `PostgreSQL full-text search` olmalÄ±dÄ±r.
- AyrÄ± search provider veya servis baseline mimarinin parÃ§asÄ± deÄŸildir.
- Search indeks yapÄ±sÄ± write model'i gÃ¶lgeleyen ikinci source-of-truth haline getirilmemelidir.

## Kapsam
- manga title, slug, summary ve taxonomy aramasÄ±
- basic ranking ve filter kombinasyonlarÄ±
- gerektiÄŸinde editorial collection aramasÄ±
- reindex ve drift recovery sÃ¼reci

## Reindex ve Lifecycle KurallarÄ±
- Search index kaynaÄŸÄ± canonical write model'dir; owner veri `manga` modÃ¼lÃ¼nde kalÄ±r.
- Reindex tam rebuild veya scoped rebuild olarak Ã§alÄ±ÅŸabilmelidir.
- Publish, archive veya taxonomy deÄŸiÅŸimi search projection gÃ¼ncellemesini tetiklemelidir.
- Reindex iÅŸlemleri replay-safe ve idempotent olmalÄ±dÄ±r.

## Provider Swap Kriterleri
AÅŸaÄŸÄ±daki koÅŸullar netleÅŸmeden ayrÄ± search provider'a geÃ§ilmemelidir:
- PostgreSQL full-text ile kabul edilebilir arama kalitesi veya latency saÄŸlanamamasÄ±
- typo tolerance, weighted ranking, synonym expansion veya faceted search ihtiyacÄ±nÄ±n baseline'Ä± aÅŸmasÄ±
- yÃ¼ksek hacimli ayrÄ± index operasyonlarÄ±nÄ±n DB yÃ¼kÃ¼nÃ¼ sÃ¼rdÃ¼rÃ¼lemez hale getirmesi

## Fallback DavranÄ±ÅŸÄ±
- Search provider veya projection geÃ§ici sorun yaÅŸarsa listing fallback'i kontrollÃ¼ biÃ§imde daraltÄ±labilmelidir.
- Hata durumunda owner iÃ§erik gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ bozulmamalÄ±; gerekirse search yÃ¼zeyi geÃ§ici olarak pasife alÄ±nmalÄ±dÄ±r.
- Search availability kararlarÄ± `docs/settings/index.md` ile hizalÄ± anahtarlar Ã¼zerinden yÃ¶netilmelidir.

## Ä°lgili Referanslar
- Projection stratejisi: `docs/shared/projection-strategy.md`
- Runtime ayar envanteri: `docs/settings/index.md`
- Manga modÃ¼lÃ¼: `docs/modules/manga.md`


---

# Support DurumlarÄ±

> `support` modÃ¼lÃ¼ndeki intake ve Ã§Ã¶zÃ¼m lifecycle'Ä± iÃ§in canonical sÃ¶zlÃ¼k.

## Support Status
| DeÄŸer | AnlamÄ± |
| --- | --- |
| `open` | Yeni aÃ§Ä±lmÄ±ÅŸ kayÄ±t |
| `triaged` | Ä°lk sÄ±nÄ±flandÄ±rmasÄ± yapÄ±lmÄ±ÅŸ |
| `waiting_user` | KullanÄ±cÄ± yanÄ±tÄ± bekleniyor |
| `waiting_team` | Ä°Ã§ ekip iÅŸlemi bekleniyor |
| `resolved` | Ã‡Ã¶zÃ¼m uygulanmÄ±ÅŸ |
| `rejected` | GeÃ§ersiz veya kapsam dÄ±ÅŸÄ± kayÄ±t |
| `closed` | SÃ¼reÃ§ tamamlanmÄ±ÅŸ ve kapatÄ±lmÄ±ÅŸ |
| `spam` | Spam veya kÃ¶tÃ¼ye kullanÄ±m olarak iÅŸaretlenmiÅŸ |

## Reply Visibility
| DeÄŸer | AnlamÄ± |
| --- | --- |
| `public_to_requester` | Talep sahibine gÃ¶rÃ¼nÃ¼r yanÄ±t |
| `internal_only` | Sadece ekip iÃ§i not veya yanÄ±t |


---

# Hedef Tipleri

> Canonical kayÄ±t dosyasÄ±: `target_type` veya benzeri hedef tipi taÅŸÄ±yan tÃ¼m modÃ¼ller bu dosyadaki deÄŸerleri kullanmalÄ±dÄ±r.

## AmaÃ§
- `target_type` deÄŸerlerinin modÃ¼ller arasÄ±nda tutarlÄ± kalmasÄ±nÄ± saÄŸlamak.
- Yeni hedef tipleri eklenirken owner modÃ¼lÃ¼, consumer modÃ¼lleri ve kullanÄ±m amacÄ±nÄ± gÃ¶rÃ¼nÃ¼r tutmak.
- `comment`, `moderation`, `support` ve ileride hedefe baÄŸlÄ± Ã§alÄ±ÅŸacak diÄŸer modÃ¼ller iÃ§in tek kaynak sunmak.

## KullanÄ±m KurallarÄ±
- `target_type` yalnÄ±zca hedefe baÄŸlÄ± kayÄ±tlarda kullanÄ±lmalÄ±dÄ±r; `support_kind=communication` veya hedefsiz `support_kind=ticket` kayÄ±tlarÄ±nda `target_type` ve `target_id` boÅŸ bÄ±rakÄ±labilir.
- `target_type` deÄŸeri mÃ¼mkÃ¼n olduÄŸunda canonical leaf modÃ¼l adÄ± ile aynÄ± olmalÄ±dÄ±r.
- GÃ¶rÃ¼nÃ¼m, ekran, aksiyon veya alt yÃ¼zey bilgisi `target_type` iÃ§ine gÃ¶mÃ¼lmemelidir; bu bilgi ayrÄ± alan veya context verisi ile taÅŸÄ±nmalÄ±dÄ±r.
- Yeni `target_type` deÄŸeri eklendiÄŸinde aynÄ± deÄŸiÅŸiklik setinde bu dosya, ilgili owner modÃ¼l dokÃ¼manÄ± ve ilgili consumer modÃ¼l dokÃ¼manÄ± gÃ¼ncellenmelidir.
- KullanÄ±lmayan veya kaldÄ±rÄ±lan `target_type` deÄŸeri dokÃ¼manda durum notu ile iÅŸaretlenmeden sessizce silinmemelidir.

## Canonical KayÄ±tlar
| `target_type` | Owner Module | BaÅŸlangÄ±Ã§ Consumer'lar | AÃ§Ä±klama | Status | Notes |
| --- | --- | --- | --- | --- | --- |
| `manga` | `manga` | `comment`, `moderation`, `support` | Manga veya seri dÃ¼zeyindeki hedef varlÄ±k. | `active` | Public iÃ§erik ve bildirim hedefi. |
| `chapter` | `chapter` | `comment`, `moderation`, `support` | Okunabilir bÃ¶lÃ¼m hedefi. | `active` | Okuma ve ihlal bildirimleri iÃ§in kullanÄ±lÄ±r. |
| `comment` | `comment` | `moderation`, `support` | Yorum hedefi. | `active` | Yorum bildirimi ve inceleme akÄ±ÅŸlarÄ± iÃ§in kullanÄ±lÄ±r. |
| `social` | `social` | `-` | Sosyal duvar, iliÅŸki ve mesajlaÅŸma yÃ¼zeylerinin canonical hedef alanÄ± iÃ§in ayrÄ±lmÄ±ÅŸ rezerv kayÄ±t. | `planned` | Ä°leride moderation veya support ihtiyacÄ± doÄŸarsa consumer aynÄ± deÄŸiÅŸiklik setinde eklenmelidir. |


---

# Teknik Stack ve AraÃ§ SeÃ§imleri

> Ã–neri seviyesinde kalan paket ve araÃ§ seÃ§imleri bu dokÃ¼manda canonical teknik karara dÃ¶nÃ¼ÅŸtÃ¼rÃ¼lÃ¼r. AynÄ± sorumluluk iÃ§in ikinci bir varsayÄ±lan araÃ§ seÃ§ilmemelidir.

## AmaÃ§
- Repo geneline daÄŸÄ±lmÄ±ÅŸ teknik tercihleri tek bir aktif karara dÃ¶nÃ¼ÅŸtÃ¼rmek.
- Kod yazÄ±mÄ± baÅŸladÄ±ÄŸÄ±nda hangi temel kÃ¼tÃ¼phane ve araÃ§larÄ±n kullanÄ±lacaÄŸÄ±nÄ± netleÅŸtirmek.
- Setup, test, cache ve altyapÄ± kararlarÄ±nÄ± ortak referansla baÄŸlamak.

## Backend ve Uygulama AraÃ§larÄ±
| Sorumluluk | Canonical SeÃ§im | Not |
| --- | --- | --- |
| HTTP router | `chi` | Hafif, aÃ§Ä±k ve Go ekosistemiyle uyumlu tercih. |
| Config / env loader | `caarlos0/env` | Env tabanlÄ± config standardÄ± iÃ§in kullanÄ±lÄ±r. |
| Structured logging | `zap` | YÃ¼ksek performanslÄ± structured logging standardÄ±. |
| UUID | `google/uuid` | Domain ve contract tarafÄ±nda ortak UUID Ã¼retimi. |
| Input validation | `go-playground/validator/v10` | DTO ve request validation iÃ§in canonical katman. |
| Migration | `golang-migrate` | TÃ¼m DB migration akÄ±ÅŸlarÄ±nÄ±n canonical aracÄ±. |
| SQL eriÅŸimi | `pgx/v5` | PostgreSQL iÃ§in canonical driver. |
| Connection pool | `pgxpool` | `pgx` ile birlikte canonical pooling katmanÄ±. |
| Password hashing | `argon2id` | Credential gÃ¼venliÄŸi iÃ§in default seÃ§im. |
| Test assertion / helper | `testify` | Go stdlib test akÄ±ÅŸÄ± yanÄ±nda helper ve assertion desteÄŸi. |

## Platform ve Ã‡alÄ±ÅŸma KararlarÄ±
- Docker-first Ã§alÄ±ÅŸma standardÄ± korunmalÄ±dÄ±r.
- Main ve test veritabanÄ± kesin olarak ayrÄ± tutulmalÄ±dÄ±r.
- Cache ihtiyaÃ§larÄ±nda canonical backend `Redis` olmalÄ±dÄ±r.
- Asenkron iÅŸleme baseline kararÄ± PostgreSQL-backed jobs + transactional outbox olmalÄ±dÄ±r.

## Uygulama KurallarÄ±
- AynÄ± sorumluluk iÃ§in farklÄ± bir araÃ§ kullanÄ±lacaksa Ã¶nce bu dokÃ¼man gÃ¼ncellenmelidir.
- Wrapper veya adapter katmanÄ± aÃ§Ä±lmasÄ± canonical seÃ§imi deÄŸiÅŸtirmez; alttaki teknik tercih aynÄ± kalmalÄ±dÄ±r.
- Logging, validation ve config eriÅŸimi platform katmanÄ±nda ortaklaÅŸtÄ±rÄ±lmalÄ±dÄ±r.
- Access policy deÄŸerlendirmesi iÃ§in dÄ±ÅŸ aÄŸÄ±r policy engine varsayÄ±lan Ã§Ã¶zÃ¼m deÄŸildir; hafif in-house evaluator tercih edilmelidir.

## Ä°lgili Referanslar
- Kurulum ve komut standardÄ±: `docs/SETUP.md`
- Cache ve queue kararÄ±: `docs/shared/cache-queue-strategy.md`
- Operasyon ve secret/config ayrÄ±mÄ±: `docs/shared/operational-standards.md`


---

# Transaction SÄ±nÄ±rlarÄ±

> Hangi akÄ±ÅŸÄ±n tek transaction iÃ§inde, hangisinin event destekli Ã§ok aÅŸamalÄ± zincirle Ã§alÄ±ÅŸmasÄ± gerektiÄŸini bu dokÃ¼man tanÄ±mlar.

## AmaÃ§
- Tek owner'lÄ± akÄ±ÅŸlar ile Ã§ok modÃ¼llÃ¼ orkestrasyonlarÄ± ayÄ±rmak.
- Eventual consistency gereken akÄ±ÅŸlar iÃ§in net sÄ±nÄ±r Ã§izmek.
- Recovery, reconcile ve manuel mÃ¼dahale ihtiyacÄ±nÄ± Ã¶nceden gÃ¶rÃ¼nÃ¼r kÄ±lmak.

## SeÃ§im KurallarÄ±
- Tek modÃ¼l, tek veritabanÄ± transaction'Ä± ve tek owner aggregate iÃ§inde kalan akÄ±ÅŸlar mÃ¼mkÃ¼n olduÄŸunda tek transaction kullanmalÄ±dÄ±r.
- DÄ±ÅŸ provider, queue veya baÅŸka modÃ¼l owner'lÄ±ÄŸÄ± gerektiren akÄ±ÅŸlar Ã§ok aÅŸamalÄ± ve idempotent orkestrasyon ile yÃ¼rÃ¼tÃ¼lmelidir.
- Bir transaction iÃ§inde â€œDB yaz + hemen publishâ€ yaklaÅŸÄ±mÄ± yerine transactional outbox planÄ± tercih edilmelidir.
- Kompanzasyon veya recovery gerektiren her akÄ±ÅŸ dokÃ¼mante edilmeden â€œsonra bakarÄ±zâ€ yaklaÅŸÄ±mÄ±yla bÄ±rakÄ±lmamalÄ±dÄ±r.

## Referans AkÄ±ÅŸlar
| AkÄ±ÅŸ | Boundary | GerekÃ§e | Recovery / Kompanzasyon |
| --- | --- | --- | --- |
| `auth login + session create` | tek transaction | credential doÄŸrulama ve session write aynÄ± owner alanda kalÄ±r | session revoke ve gÃ¼venlik audit'i ile geri alÄ±nÄ±r |
| `shop purchase -> payment callback -> inventory grant` | Ã§ok aÅŸamalÄ±, event destekli | katalog, Ã¶deme ve sahiplik Ã¼Ã§ ayrÄ± owner alandÄ±r | reconcile, duplicate guard ve grant retry gerekir |
| `mission complete -> reward grant` | koordineli ama idempotent | claim uygunluÄŸu ile final grant farklÄ± owner'lardadÄ±r | claim replay ve grant dedup ile toparlanÄ±r |
| `support report -> moderation case` | policy'ye gÃ¶re sync veya async | backpressure ve queue ihtiyacÄ± oluÅŸabilir | linked case reference ile tekrar denenir |
| `notification create -> channel delivery` | write + async delivery | delivery dÄ±ÅŸ kanal ve retry gerektirir | backoff, suppression ve dead-letter gerekir |
| `royalpass premium activation` | Ã§ok aÅŸamalÄ±, event destekli | `shop`, `payment` ve `royalpass` owner'lÄ±klarÄ± ayrÄ±dÄ±r | activation ref, reconcile ve replay gÃ¼venli olmalÄ±dÄ±r |

## Uygulama KurallarÄ±
- BaÅŸka modÃ¼lÃ¼n tablosuna doÄŸrudan yazÄ± yapan transaction tasarlanmamalÄ±dÄ±r.
- `request_id`, `correlation_id` ve idempotency key kritik geÃ§iÅŸlerde boundary boyunca taÅŸÄ±nmalÄ±dÄ±r.
- Finansal veya kullanÄ±cÄ± hakkÄ± etkileyen akÄ±ÅŸlarda manuel review kapÄ±sÄ± aÃ§Ä±kÃ§a iÅŸaretlenmelidir.


---

# GÃ¶rÃ¼nÃ¼rlÃ¼k DurumlarÄ±

> GÃ¶rÃ¼nÃ¼rlÃ¼k anlamÄ± taÅŸÄ±yan modÃ¼ller mÃ¼mkÃ¼n olduÄŸunda bu canonical sÃ¶zlÃ¼ÄŸe hizalanmalÄ±dÄ±r.

## KullanÄ±m KurallarÄ±
- Her modÃ¼l tÃ¼m durumlarÄ± kullanmak zorunda deÄŸildir.
- ModÃ¼l iÃ§i Ã¶zel state alanlarÄ± bu sÃ¶zlÃ¼ÄŸe map edilebiliyorsa modÃ¼l dokÃ¼manÄ±nda aÃ§Ä±kÃ§a belirtilmelidir.
- Final eriÅŸim kararÄ± yine `docs/shared/precedence-rules.md` ve `access` yorumu ile verilir.

## Canonical Durumlar
| Durum | AnlamÄ± | Tipik ModÃ¼ller | Not |
| --- | --- | --- | --- |
| `public` | Herkese aÃ§Ä±k gÃ¶rÃ¼nÃ¼r yÃ¼zey | `manga`, `chapter`, `user`, `social` | access ve audience kontrolÃ¼ ayrÄ±ca uygulanabilir |
| `limited` | Belirli audience veya entitlement ile gÃ¶rÃ¼nÃ¼r | `chapter`, `royalpass`, `ads` | VIP, early access veya scoped audience ile birlikte kullanÄ±lÄ±r |
| `private` | YalnÄ±zca owner veya yetkili yÃ¼zeye gÃ¶rÃ¼nÃ¼r | `history`, `support`, `notification` | public default Ã¼retmez |
| `hidden` | Mevcut kayÄ±t durur ama dÄ±ÅŸ gÃ¶rÃ¼nÃ¼rlÃ¼k kapanÄ±r | `comment`, `manga`, `social` | moderation veya runtime switch kaynaklÄ± olabilir |
| `removed` | DÄ±ÅŸ yÃ¼zeyden kaldÄ±rÄ±lmÄ±ÅŸ veya soft-deleted gÃ¶rÃ¼nÃ¼m | `comment`, `support` | audit ve recovery notu taÅŸÄ±malÄ±dÄ±r |
| `archived` | Aktif deÄŸil ama geÃ§miÅŸ kayÄ±t olarak tutulur | `manga`, `royalpass`, `ads`, `shop` | read-only veya operasyonel gÃ¶rÃ¼nÃ¼mde kalabilir |


---

# Runtime Ayar Envanteri

> Canonical kayÄ±t dosyasÄ±: admin tarafÄ±ndan yÃ¶netilen tÃ¼m runtime ayarlar, feature toggle'lar, kill switch yÃ¼zeyleri ve oran limitleri bu dokÃ¼manda tutulmalÄ±dÄ±r.

## AmaÃ§
- Runtime ayarlarÄ±n tek merkezden izlenmesini saÄŸlamak.
- Ayar anahtarlarÄ±nÄ±n Ã§akÄ±ÅŸmadan bÃ¼yÃ¼mesini saÄŸlamak.
- Hangi ayarÄ±n hangi modÃ¼lde Ã¼retildiÄŸini ve hangi katmanda yorumlandÄ±ÄŸÄ±nÄ± gÃ¶rÃ¼nÃ¼r tutmak.

## YaÅŸayan DokÃ¼man Notu
- Bu dosya kapanmÄ±ÅŸ bir checklist deÄŸil, modÃ¼l yÃ¼zeyleri bÃ¼yÃ¼dÃ¼kÃ§e gÃ¼ncellenen canonical runtime envanteridir.
- ModÃ¼l dokÃ¼manlarÄ±nda runtime ile kontrol edilebilir olduÄŸu yazÄ±lan her surface burada birebir temsil edilmeli veya neden henÃ¼z `planned` kaldÄ±ÄŸÄ± aÃ§Ä±kÃ§a not edilmelidir.

## KayÄ±t KurallarÄ±
- Yeni runtime ayar, feature toggle, kill switch veya oran limiti eklendiÄŸinde bu dosya aynÄ± deÄŸiÅŸiklik setinde gÃ¼ncellenmelidir.
- AynÄ± `key + audience_kind + audience_selector + scope_kind + scope_selector` kombinasyonu iÃ§in birden fazla aktif kayÄ±t bÄ±rakÄ±lmamalÄ±dÄ±r.
- Ãœcretli veya sÃ¼reli avantajÄ± etkileyen ayarlarda `entitlement_impact_policy` alanÄ± zorunludur.
- `schedule_support` alanÄ± en az `none`, `start_at` ve `time_window` modlarÄ±nÄ± ayÄ±rt edebilmelidir.
- `docs/RULES.md` ve ilgili modÃ¼l dokÃ¼manÄ± ile Ã§eliÅŸen kayÄ±t bÄ±rakÄ±lamaz.
- Bir modÃ¼l dokÃ¼manÄ±nda ayrÄ± ayrÄ± runtime kontrol edilebildiÄŸi yazÄ±lan yÃ¼zeyler settings envanterinde ya kendi canonical key kaydÄ±yla ya da kapsadÄ±ÄŸÄ± alt yÃ¼zeyleri aÃ§Ä±kÃ§a listeleyen umbrella key kaydÄ±yla temsil edilmelidir.
- Yeni bir surface iÃ§in yalnÄ±zca availability anahtarÄ± yazmak yeterli deÄŸildir; ilgili rate limit, threshold, cooldown, disabled behavior veya degrade davranÄ±ÅŸÄ± da aynÄ± deÄŸiÅŸiklikte eklenmeli ya da neden henÃ¼z `planned` kaldÄ±ÄŸÄ± `notes` alanÄ±nda belirtilmelidir.
- Service katmanÄ±nda yorumlanan intake pause, callback gate, digest window veya benzeri operasyonel ayarlar iÃ§in `consumer_layer` alanÄ± zorunlu olarak doldurulmalÄ± ve access availability ayarlarÄ± ile karÄ±ÅŸtÄ±rÄ±lmamalÄ±dÄ±r.

## Audience SÃ¶zlÃ¼ÄŸÃ¼
- `all`: tÃ¼m kullanÄ±cÄ±lar.
- `guest`: giriÅŸ yapmamÄ±ÅŸ ziyaretÃ§i.
- `authenticated`: giriÅŸ yapmÄ±ÅŸ tÃ¼m kullanÄ±cÄ±lar.
- `authenticated_non_vip`: giriÅŸ yapmÄ±ÅŸ ancak VIP olmayan kullanÄ±cÄ±lar.
- `vip`: aktif VIP avantajÄ± taÅŸÄ±yan kullanÄ±cÄ±lar.

## Scope SÃ¶zlÃ¼ÄŸÃ¼
- `site`: tÃ¼m Ã¼rÃ¼n yÃ¼zeyi.
- `module`: belirli modÃ¼l geneli.
- `feature`: modÃ¼l iÃ§indeki belirli Ã¶zellik veya alt yÃ¼zey.
- `resource/context`: belirli kaynak, ekran veya Ã¶zel baÄŸlam.

## Key Grammar
- Boolean availability veya feature toggle anahtarlarÄ± mÃ¼mkÃ¼n olduÄŸunda `feature.<module>.<surface>.enabled` biÃ§imini kullanmalÄ±dÄ±r.
- EÅŸik, limit, cooldown veya davranÄ±ÅŸ deÄŸeri taÅŸÄ±yan anahtarlar mÃ¼mkÃ¼n olduÄŸunda `<module>.<surface>.<metric>` biÃ§imini kullanmalÄ±dÄ±r.
- Feature toggle niteliÄŸi taÅŸÄ±mayan operasyonel pause, intake veya bakÄ±m bayraklarÄ± gerektiÄŸinde `<module>.<surface>.<flag>` biÃ§imini kullanabilir.
- Site geneli operasyon veya bakÄ±m anahtarlarÄ± mÃ¼mkÃ¼n olduÄŸunda `site.<surface>.<metric_or_flag>` biÃ§imini kullanmalÄ±dÄ±r.
- Audience, role, grup, kullanÄ±cÄ± veya resource bilgisi runtime key iÃ§ine gÃ¶mÃ¼lmemelidir; bu bilgiler selector alanlarÄ±nda taÅŸÄ±nmalÄ±dÄ±r.
- Key iÃ§indeki `<module>` bÃ¶lÃ¼mÃ¼ her zaman canonical leaf modÃ¼l adÄ± ile aynÄ± yazÄ±lmalÄ±dÄ±r.

## Selector Grammar
- `scope_selector` iÃ§in baÅŸlangÄ±Ã§ canonical yaklaÅŸÄ±m `-`, `<module>`, `<module>.<surface>`, `<module>.<surface>.<subsurface>` ve gerektiÄŸinde `resource:<module>:<resource_kind>:<identifier>` biÃ§imleridir.
- `audience_selector` iÃ§in baÅŸlangÄ±Ã§ canonical yaklaÅŸÄ±m `-`, `role:<name>`, `group:<name>` ve `user:<id>` biÃ§imleridir.
- Selector iÃ§indeki modÃ¼l adÄ± her zaman canonical leaf modÃ¼l adÄ± ile aynÄ± yazÄ±lmalÄ±dÄ±r.

## Disabled Behavior SÃ¶zlÃ¼ÄŸÃ¼
- `visibility_off`: yÃ¼zeyin dÄ±ÅŸarÄ±ya tamamen kapatÄ±lmasÄ±.
- `read_only`: yalnÄ±zca okuma izni verilip yazma tarafÄ±nÄ±n kapatÄ±lmasÄ±.
- `write_off`: mevcut veri gÃ¶rÃ¼nÃ¼r kalÄ±rken yeni yazma aksiyonunun kapatÄ±lmasÄ±.
- `intake_pause`: yeni kayÄ±t alÄ±mÄ±nÄ±n durdurulmasÄ±.
- `read_only_intake`: mevcut kayÄ±tlarÄ±n okunur kaldÄ±ÄŸÄ±, ancak yeni intake/create akÄ±ÅŸÄ±nÄ±n durdurulduÄŸu durum.
- `attachment_off`: ana create akÄ±ÅŸÄ± aÃ§Ä±k kalÄ±rken attachment kabulÃ¼nÃ¼n kapatÄ±lmasÄ±.
- `preview_off`: yalnÄ±zca Ã¶nizleme yÃ¼zeyinin kapatÄ±lmasÄ±.
- `benefit_pause`: Ã¼cretli veya sÃ¼reli avantajÄ±n geÃ§ici pasife alÄ±nÄ±p sÃ¼renin dondurulmasÄ±.

## Error Response Policy SÃ¶zlÃ¼ÄŸÃ¼
- `not_found`: dÄ±ÅŸ yÃ¼zeyde kaynak veya yÃ¼zey yokmuÅŸ gibi davranÄ±lmasÄ±.
- `forbidden`: yÃ¼zey gÃ¶rÃ¼nÃ¼r olsa bile ilgili aksiyonun reddedilmesi.
- `rate_limited`: eÅŸik, cooldown veya throttling nedeniyle geÃ§ici reddedilme.
- `validation_error`: istek biÃ§imi geÃ§erli olsa bile ilgili alt yÃ¼zey kapalÄ± olduÄŸu iÃ§in alan bazlÄ± doÄŸrulama hatasÄ± dÃ¶nÃ¼lmesi.
- `temporarily_unavailable`: sistem kaynaklÄ± geÃ§ici pasiflik veya operasyonel durdurma.

## Entitlement Impact Policy SÃ¶zlÃ¼ÄŸÃ¼
- `none`: sÃ¼reli hak veya entitlement Ã¼zerinde ek etki yoktur.
- `freeze_on_system_disable`: sistem kaynaklÄ± pasiflikte kalan entitlement sÃ¼resi dondurulur ve hak sÃ¼resi korunur.

## KayÄ±t AlanlarÄ±
| Alan | AÃ§Ä±klama |
| --- | --- |
| `key` | Canonical ayar anahtarÄ± |
| `description` | AyarÄ±n neyi kontrol ettiÄŸini aÃ§Ä±klayan kÄ±sa metin |
| `category` | `site`, `communication`, `operations`, `security_auth`, `access_availability`, `content`, `reading`, `engagement`, `support`, `membership`, `social`, `gamification`, `economy` |
| `owner_module` | AyarÄ±n sahibi olan modÃ¼l |
| `consumer_layer` | AyarÄ± yorumlayan katman: Ã¶rnek `access`, `service`, `app` |
| `value_type` | `bool`, `int`, `duration`, `enum`, `string` vb. |
| `default_value` | VarsayÄ±lan deÄŸer |
| `allowed_range_or_enum` | Ä°zin verilen aralÄ±k veya enum listesi |
| `scope_kind` | `site`, `module`, `feature`, `resource/context` |
| `scope_selector` | Scope'un somut hedefi; Ã¶rnek `comment.write`, `manga.detail`, `chapter.preview`, `-` |
| `audience_kind` | GeÃ§erli audience tÃ¼rleri |
| `audience_selector` | Audience'Ä±n somut hedefi; Ã¶rnek `role:comment_moderator`, `group:testers`, `-` |
| `sensitive` | Hassas veri iÃ§erip iÃ§ermediÄŸi |
| `apply_mode` | `immediate`, `cache_refresh`, `scheduled` |
| `cache_strategy` | `none`, `ttl`, `manual_invalidate` |
| `schedule_support` | `none`, `start_at`, `time_window` |
| `audit_required` | DeÄŸiÅŸiklik audit zorunluluÄŸu |
| `affected_surfaces` | Etkilenen yÃ¼zeyler |
| `disabled_behavior` | Varsa kapatma davranÄ±ÅŸ tipi |
| `error_response_policy` | DÄ±ÅŸ API veya yÃ¼zey iÃ§in hata/geri dÃ¶nÃ¼ÅŸ politikasÄ± |
| `entitlement_impact_policy` | Varsa sÃ¼reli hak etkisi |
| `status` | `active`, `planned`, `deprecated` |
| `notes` | Ek aÃ§Ä±klama |

## Access Yorumlama Modeli
- KullanÄ±cÄ±ya bakan availability, permission ve entitlement etkili runtime ayarlarÄ± canonical olarak `access` tarafÄ±ndan yorumlanmalÄ±dÄ±r.
- Yorum sÄ±rasÄ± `docs/shared/precedence-rules.md` ile hizalÄ± olmalÄ±; `global -> module -> surface -> audience -> entitlement -> action -> rate limit` akÄ±ÅŸÄ± korunmalÄ±dÄ±r.
- Operasyonel pause, callback intake, queue backpressure veya delivery retry benzeri kullanÄ±cÄ±ya gÃ¶rÃ¼nmeyen ayarlar ilgili `service` katmanÄ±nda yorumlanabilir; bu durum settings envanterinde `consumer_layer` ile gÃ¶rÃ¼nÃ¼r tutulmalÄ±dÄ±r.
- AynÄ± surface iÃ§in read ve write ayrÄ±mÄ± varsa ayrÄ± anahtar veya umbrella key notu bÄ±rakÄ±lmalÄ±dÄ±r; Ã¶rtÃ¼k davranÄ±ÅŸ bÄ±rakÄ±lmamalÄ±dÄ±r.

## Kill Switch Seviyeleri
| Seviye | AÃ§Ä±klama | Tipik KullanÄ±m |
| --- | --- | --- |
| `global` | TÃ¼m Ã¼rÃ¼n veya site yÃ¼zeyi etkilenir | bakÄ±m modu, genel gÃ¼venlik durdurmasÄ± |
| `module-level` | Bir modÃ¼lÃ¼n ana yÃ¼zeyleri topluca kapanÄ±r | `support.intake.paused`, `notification.delivery.paused` |
| `surface-level` | ModÃ¼l iÃ§indeki belirli feature kapanÄ±r | `feature.social.messaging.enabled` |
| `action-level` | Surface aÃ§Ä±k kalÄ±rken tek aksiyon kapanÄ±r | `feature.inventory.claim.enabled` |
| `intake-only` | Yeni kayÄ±t alÄ±mÄ± durur, mevcut veri okunabilir kalÄ±r | support intake, payment callback intake |
| `write-only` | Yazma kapalÄ±, okuma aÃ§Ä±k kalÄ±r | comment write, bookmark write |
| `external-integration-only` | DÄ±ÅŸ provider veya delivery kanalÄ± kapatÄ±lÄ±r | payment callback, notification channel send |

## ModÃ¼l YÃ¼zeyi Hizalama Matrisi
| Module | Surface | Subsurface | Canonical / Planned Key | Disabled Behavior | Error Response Policy | Audience Selector Support | Schedule Support | Entitlement Impact | Rate Limit Support | Note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `manga` | discovery | recommendation, collection, editorial | `feature.manga.discovery.enabled` | `visibility_off` | `not_found` | var | `none` | `none` | yok | Listing ve detail yÃ¼zeyinden baÄŸÄ±msÄ±z kalmalÄ±dÄ±r. |
| `chapter` | read | preview, detail, early access | `feature.chapter.preview.enabled`, `feature.chapter.detail.enabled`, `feature.chapter.read.enabled`, `feature.chapter.early_access.enabled` | `preview_off`, `visibility_off` | `not_found` | var | `time_window` | `none` | yok | Early access gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ayrÄ± toggle taÅŸÄ±malÄ±dÄ±r. |
| `comment` | write | create, edit, delete | `feature.comment.write.enabled`, `comment.write.cooldown_seconds` | `write_off` | `forbidden`, `rate_limited` | var | `none` | `none` | var | Anti-spam ve edit window service katmanÄ±nda yorumlanÄ±r. |
| `history` | library | continue reading, timeline, bookmark.write | `feature.history.continue_reading.enabled`, `feature.history.library.enabled`, `feature.history.timeline.enabled`, `feature.history.bookmark_write.enabled` | `visibility_off`, `write_off` | `not_found`, `forbidden` | var | `none` | `none` | yok | Entry-level share metadata ayrÄ± access yorumuna baÄŸlanÄ±r. |
| `social` | interaction | friendship, follow, wall, messaging | `feature.social.friendship.enabled`, `feature.social.follow.enabled`, `feature.social.wall.enabled`, `feature.social.messaging.enabled` | `write_off`, `visibility_off` | `forbidden`, `not_found` | var | `none` | `none` | var | Block ve mute precedence access tarafÄ±nda yorumlanÄ±r. |
| `notification` | delivery | inbox, preference, digest, channel send | `feature.notification.inbox.enabled`, `feature.notification.preference.enabled`, `feature.notification.digest.enabled`, `notification.delivery.paused` | `visibility_off`, `read_only` | `not_found`, `temporarily_unavailable` | var | `time_window` | `none` | var | Digest ve channel pause service tarafÄ±nda da yorumlanÄ±r. |
| `inventory` | ownership | read, claim, equip, consume | `feature.inventory.read.enabled`, `feature.inventory.claim.enabled`, `feature.inventory.equip.enabled`, `feature.inventory.consume.enabled` | `visibility_off`, `write_off` | `not_found`, `forbidden` | var | `none` | `none` | yok | Grant tarafÄ± idempotent Ã§alÄ±ÅŸmalÄ±, read/write ayrÄ±mÄ± korunmalÄ±dÄ±r. |
| `mission` | progression | read, claim, progress_ingest | `feature.mission.read.enabled`, `feature.mission.claim.enabled`, `feature.mission.progress_ingest.enabled` | `visibility_off`, `write_off`, `intake_pause` | `not_found`, `forbidden`, `temporarily_unavailable` | var | `time_window` | `none` | yok | Event ingest ayrÄ± kapatÄ±labilir olmalÄ±dÄ±r. |
| `royalpass` | season | season, premium, claim | `feature.royalpass.season.enabled`, `feature.royalpass.premium.enabled`, `feature.royalpass.claim.enabled` | `benefit_pause`, `visibility_off`, `write_off` | `temporarily_unavailable`, `not_found`, `forbidden` | var | `time_window` | `freeze_on_system_disable` | yok | Premium entitlement ve season availability ayrÄ± yorumlanmalÄ±dÄ±r. |
| `shop` | purchase | catalog, campaign, purchase, recovery | `feature.shop.catalog.enabled`, `feature.shop.campaign.enabled`, `feature.shop.purchase.enabled` | `visibility_off`, `write_off` | `not_found`, `forbidden` | var | `time_window` | `none` | var | Recovery ve already-owned davranÄ±ÅŸÄ± service katmanÄ±nda tamamlanÄ±r. |
| `payment` | checkout | mana purchase, checkout, transaction read, callback intake | `feature.payment.mana_purchase.enabled`, `feature.payment.checkout.enabled`, `feature.payment.transaction_read.enabled`, `payment.callback.intake.paused` | `visibility_off`, `write_off`, `intake_pause` | `not_found`, `temporarily_unavailable` | var | `time_window` | `none` | var | Provider callback yÃ¼zeyi kullanÄ±cÄ± availability'sinden ayrÄ± yÃ¶netilmelidir. |
| `ads` | delivery | surface, placement, campaign, click intake | `feature.ads.surface.enabled`, `feature.ads.placement.enabled`, `feature.ads.campaign.enabled`, `feature.ads.click_intake.enabled` | `visibility_off`, `intake_pause` | `not_found`, `temporarily_unavailable` | var | `time_window` | `none` | var | VIP no-ads precedence access ile yorumlanÄ±r. |
| `support` | intake | communication, ticket, report, attachment, internal_note | `feature.support.communication.enabled`, `feature.support.ticket.enabled`, `feature.support.report.enabled`, `feature.support.attachment.enabled`, `feature.support.internal_note.enabled`, `support.intake.paused` | `write_off`, `attachment_off`, `read_only_intake` | `forbidden`, `validation_error`, `temporarily_unavailable` | var | `time_window` | `none` | var | Report-to-case policy ve internal note gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ayrÄ± tutulmalÄ±dÄ±r. |
## BaÅŸlangÄ±Ã§ Referans KayÄ±tlarÄ±
| Key | Description | Category | Owner Module | Consumer Layer | Value Type | Default | Allowed Range / Enum | Scope Kind | Scope Selector | Audience Kind | Audience Selector | Sensitive | Apply Mode | Cache Strategy | Schedule Support | Audit Required | Affected Surfaces | Disabled Behavior | Error Response Policy | Entitlement Impact | Status | Notes |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `site.maintenance.enabled` | BakÄ±m modu anahtarÄ± | `operations` | `admin` | `app` | `bool` | `false` | `true,false` | `site` | `-` | `all` | `-` | `false` | `immediate` | `none` | `time_window` | `true` | `site.*` | `visibility_off` | `temporarily_unavailable` | `none` | `planned` | BakÄ±m modu iÃ§in referans anahtar |
| `auth.login.failed_attempt_limit_per_minute` | Dakika baÅŸÄ±na baÅŸarÄ±sÄ±z giriÅŸ limiti | `security_auth` | `auth` | `service` | `int` | `5` | `1-20` | `site` | `-` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `auth.login` | `-` | `rate_limited` | `none` | `planned` | BaÅŸarÄ±sÄ±z giriÅŸ limiti |
| `auth.login.cooldown_seconds` | BaÅŸarÄ±sÄ±z giriÅŸ eÅŸiÄŸi sonrasÄ± uygulanan cooldown sÃ¼resi | `security_auth` | `auth` | `service` | `int` | `300` | `0-86400` | `site` | `-` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `auth.login` | `-` | `rate_limited` | `none` | `planned` | Failed login limiti ile birlikte Ã§alÄ±ÅŸan geÃ§ici auth kÄ±sÄ±tÄ± |
| `auth.email.verification_resend_cooldown_seconds` | Verification tekrar gÃ¶nderim bekleme sÃ¼resi | `security_auth` | `auth` | `service` | `int` | `60` | `0-3600` | `site` | `-` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `auth.email_verification` | `-` | `rate_limited` | `none` | `planned` | Verification tekrar gÃ¶nderim aralÄ±ÄŸÄ± |
| `feature.comment.read.enabled` | Yorum okuma ve thread listeleme yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `comment` | `access` | `bool` | `true` | `true,false` | `feature` | `comment.read` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `comment.read`,`comment.thread` | `visibility_off` | `not_found` | `none` | `planned` | Read yÃ¼zeyi write yÃ¼zeyinden baÄŸÄ±msÄ±z daraltÄ±labilir |
| `comment.write.cooldown_seconds` | Yorum yazma aksiyonu iÃ§in bekleme sÃ¼resi | `engagement` | `comment` | `service` | `int` | `30` | `0-3600` | `feature` | `comment.write` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `comment.write` | `-` | `rate_limited` | `none` | `planned` | Yorum yazma aralÄ±ÄŸÄ± |
| `feature.comment.write.enabled` | Yorum yazma yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `comment` | `access` | `bool` | `true` | `true,false` | `feature` | `comment.write` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `comment.write` | `write_off` | `forbidden` | `none` | `planned` | Yorum yazma yÃ¼zeyi aÃ§ma-kapama |
| `feature.user.profile.enabled` | Profil gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ve profil detail yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `user` | `access` | `bool` | `true` | `true,false` | `feature` | `user.profile` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `user.profile`,`user.profile.detail` | `visibility_off` | `not_found` | `none` | `planned` | Public profil yÃ¼zeyi own olmayan okumadan baÄŸÄ±msÄ±z daraltÄ±labilir |
| `feature.user.vip_benefits.enabled` | VIP avantajlarÄ±nÄ± sistem genelinde aÃ§ma-kapama anahtarÄ± | `access_availability` | `user` | `access` | `bool` | `true` | `true,false` | `feature` | `user.vip_benefits` | `vip` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `user.vip_benefits`,`access.vip` | `benefit_pause` | `temporarily_unavailable` | `freeze_on_system_disable` | `planned` | VIP avantajÄ±nÄ±n sistem kaynaklÄ± pasifliÄŸinde sÃ¼re dondurulur |
| `feature.user.vip_badge.enabled` | VIP rozet veya VIP profil gÃ¶stergesini aÃ§ma-kapama anahtarÄ± | `access_availability` | `user` | `access` | `bool` | `true` | `true,false` | `feature` | `user.vip_badge` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `user.vip_badge`,`user.profile` | `visibility_off` | `not_found` | `none` | `planned` | GÃ¶rsel profil gÃ¶stergesi VIP entitlement sahibinden ayrÄ± yÃ¶netilebilir |
| `feature.user.history_visibility_preference.enabled` | KullanÄ±cÄ±nÄ±n history veya library visibility preference yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `user` | `access` | `bool` | `true` | `true,false` | `feature` | `user.history_visibility_preference` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `user.history_visibility_preference` | `write_off` | `forbidden` | `none` | `planned` | Preference yÃ¼zeyi kapansa bile history owner'lÄ±ÄŸÄ± `history` modÃ¼lÃ¼nde kalÄ±r |
| `feature.manga.list.enabled` | Manga listing, search ve filter yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `manga` | `access` | `bool` | `true` | `true,false` | `feature` | `manga.list` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `manga.list`,`manga.search`,`manga.filter` | `visibility_off` | `not_found` | `none` | `planned` | Listing yÃ¼zeyi discovery veya detail yÃ¼zeyinden baÄŸÄ±msÄ±z yÃ¶netilebilir |
| `feature.manga.detail.enabled` | Manga detail yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `manga` | `access` | `bool` | `true` | `true,false` | `feature` | `manga.detail` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `manga.detail` | `visibility_off` | `not_found` | `none` | `planned` | Detail yÃ¼zeyi listing ve discovery'den baÄŸÄ±msÄ±z daraltÄ±labilir |
| `feature.chapter.preview.enabled` | Chapter preview yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `chapter` | `access` | `bool` | `true` | `true,false` | `feature` | `chapter.preview` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `chapter.preview` | `preview_off` | `not_found` | `none` | `planned` | Preview kapanÄ±nca detail veya tam read otomatik kapanmÄ±ÅŸ sayÄ±lmaz |
| `feature.chapter.detail.enabled` | Chapter detail yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `chapter` | `access` | `bool` | `true` | `true,false` | `feature` | `chapter.detail` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `chapter.detail` | `visibility_off` | `not_found` | `none` | `planned` | Detail yÃ¼zeyi preview ve tam read'den baÄŸÄ±msÄ±z yÃ¶netilebilir |
| `feature.chapter.read.enabled` | Chapter tam read yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `chapter` | `access` | `bool` | `true` | `true,false` | `feature` | `chapter.read` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `chapter.read`,`chapter.resume` | `visibility_off` | `not_found` | `none` | `planned` | Tam read kapalÄ±yken detail veya preview ayrÄ± aÃ§Ä±k kalabilir |
| `feature.moderation.panel.enabled` | Moderation paneli genel gÃ¶rÃ¼nÃ¼rlÃ¼k anahtarÄ± | `access_availability` | `moderation` | `access` | `bool` | `true` | `true,false` | `feature` | `moderation.panel` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `moderation.panel` | `visibility_off` | `not_found` | `none` | `planned` | Yetkili moderatÃ¶r yÃ¼zeyi admin tarafÄ±ndan geÃ§ici olarak kapatÄ±labilir |
| `feature.moderation.queue.enabled` | Moderation queue gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼nÃ¼ aÃ§ma-kapama anahtarÄ± | `access_availability` | `moderation` | `access` | `bool` | `true` | `true,false` | `feature` | `moderation.queue` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `moderation.queue` | `visibility_off` | `not_found` | `none` | `planned` | Queue yÃ¼zeyi panel aÃ§Ä±k kalsa bile ayrÄ± yÃ¶netilebilir |
| `feature.moderation.action.enabled` | Moderation aksiyon yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `moderation` | `access` | `bool` | `true` | `true,false` | `feature` | `moderation.action` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `moderation.action` | `write_off` | `forbidden` | `none` | `planned` | SÄ±nÄ±rlÄ± moderatÃ¶r aksiyonlarÄ± geÃ§ici olarak kapatÄ±labilir |
| `feature.notification.inbox.enabled` | In-app bildirim kutusu yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `notification` | `access` | `bool` | `true` | `true,false` | `feature` | `notification.inbox` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `notification.inbox` | `visibility_off` | `temporarily_unavailable` | `none` | `planned` | Bildirim kutusu geÃ§ici olarak pasife alÄ±nabilir |
| `feature.notification.preference.enabled` | Bildirim preference yÃ¶netim yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `notification` | `access` | `bool` | `true` | `true,false` | `feature` | `notification.preference` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `notification.preference` | `visibility_off` | `not_found` | `none` | `planned` | Preference yÃ¼zeyi inbox gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼nden baÄŸÄ±msÄ±z daraltÄ±labilir |
| `notification.delivery.paused` | Notification delivery akÄ±ÅŸÄ±nÄ± geÃ§ici durdurma anahtarÄ± | `operations` | `notification` | `service` | `bool` | `false` | `true,false` | `module` | `notification` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `notification.delivery`,`notification.channel`,`notification.category` | `read_only` | `temporarily_unavailable` | `none` | `planned` | Category veya channel bazlÄ± durdurma aynÄ± key'in selector geniÅŸlemesiyle uygulanabilir |
| `feature.social.friendship.enabled` | ArkadaÅŸlÄ±k yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `social` | `access` | `bool` | `true` | `true,false` | `feature` | `social.friendship` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `social.friendship` | `write_off` | `forbidden` | `none` | `planned` | ArkadaÅŸlÄ±k isteÄŸi ve iliÅŸki yÃ¶netimi yÃ¼zeyi ayrÄ± olarak kapatÄ±labilir |
| `feature.social.follow.enabled` | Takip yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `social` | `access` | `bool` | `true` | `true,false` | `feature` | `social.follow` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `social.follow` | `write_off` | `forbidden` | `none` | `planned` | Follow veya unfollow yazma yÃ¼zeyi ayrÄ± olarak yÃ¶netilebilir |
| `feature.social.messaging.enabled` | Sosyal mesajlaÅŸma gÃ¶nderim yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `social` | `access` | `bool` | `true` | `true,false` | `feature` | `social.messaging` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `social.messaging` | `write_off` | `forbidden` | `none` | `planned` | MesajlaÅŸma yazma kapatÄ±lsa bile mevcut thread okuma davranÄ±ÅŸÄ± ayrÄ± ayarlanabilir |
| `feature.social.wall.enabled` | Sosyal duvar yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `social` | `access` | `bool` | `true` | `true,false` | `feature` | `social.wall` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `social.wall`,`profile.wall` | `visibility_off` | `not_found` | `none` | `planned` | Profil duvarÄ± veya sosyal duvar yÃ¼zeyi yÃ¶netilebilir |
| `feature.support.communication.enabled` | Ä°letiÅŸim kaydÄ± oluÅŸturma yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `support` | `access` | `bool` | `true` | `true,false` | `feature` | `support.communication` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.communication`,`support.create` | `write_off` | `forbidden` | `none` | `planned` | Genel iletiÅŸim giriÅŸi geÃ§ici olarak kapatÄ±labilir |
| `feature.support.ticket.enabled` | Destek bileti oluÅŸturma yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `support` | `access` | `bool` | `true` | `true,false` | `feature` | `support.ticket` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.ticket`,`support.create` | `write_off` | `forbidden` | `none` | `planned` | Hedefsiz destek bileti oluÅŸturma yÃ¼zeyi geÃ§ici olarak kapatÄ±labilir |
| `feature.support.report.enabled` | Hedefe baÄŸlÄ± iÃ§erik bildirimi yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `support` | `access` | `bool` | `true` | `true,false` | `feature` | `support.report` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.report` | `write_off` | `forbidden` | `none` | `planned` | Manga, chapter veya comment iÃ§in report oluÅŸturma yÃ¼zeyi ayrÄ± olarak yÃ¶netilebilir |
| `feature.support.attachment.enabled` | Support attachment kabulÃ¼nÃ¼ aÃ§ma-kapama anahtarÄ± | `support` | `support` | `service` | `bool` | `true` | `true,false` | `feature` | `support.attachment` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.communication`,`support.ticket`,`support.report` | `attachment_off` | `validation_error` | `none` | `planned` | Attachment kapalÄ±ysa kayÄ±t akÄ±ÅŸÄ± devam eder, ancak dosya kabul edilmez |
| `support.intake.paused` | Support yeni kayÄ±t alÄ±mÄ±nÄ± durdurma anahtarÄ± | `support` | `support` | `service` | `bool` | `false` | `true,false` | `module` | `support` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `support.create`,`support.communication`,`support.ticket`,`support.report` | `read_only_intake` | `temporarily_unavailable` | `none` | `planned` | Intake pause aktifken mevcut kayÄ±tlar okunabilir kalÄ±rken yeni create yÃ¼zeyleri durdurulabilir |
| `feature.inventory.read.enabled` | Envanter list ve own item detail yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `inventory` | `access` | `bool` | `true` | `true,false` | `feature` | `inventory.read` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `inventory.list`,`inventory.detail` | `visibility_off` | `not_found` | `none` | `planned` | Liste ve detail gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ claim veya equip yÃ¼zeyinden baÄŸÄ±msÄ±z yÃ¶netilebilir |
| `feature.inventory.claim.enabled` | Envanter Ã¶dÃ¼l claim yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `inventory` | `access` | `bool` | `true` | `true,false` | `feature` | `inventory.claim` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `inventory.claim` | `write_off` | `forbidden` | `none` | `planned` | Claim kapansa bile mevcut sahiplik kayÄ±tlarÄ± gÃ¶rÃ¼nÃ¼r kalabilir |
| `feature.inventory.equip.enabled` | Envanter equip yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `inventory` | `access` | `bool` | `true` | `true,false` | `feature` | `inventory.equip` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `inventory.equip` | `write_off` | `forbidden` | `none` | `planned` | Equip yÃ¼zeyi read veya claim yÃ¼zeyinden baÄŸÄ±msÄ±z kapatÄ±labilir |
| `mission.daily.reset_hour_utc` | GÃ¼nlÃ¼k mission reset saati | `gamification` | `mission` | `service` | `int` | `0` | `0-23` | `module` | `mission` | `all` | `-` | `false` | `cache_refresh` | `manual_invalidate` | `none` | `true` | `mission.reset`,`mission.daily` | `-` | `-` | `none` | `planned` | DÃ¶nemsel mission reset davranÄ±ÅŸÄ±nÄ± belirleyen referans eÅŸik |
| `feature.mission.read.enabled` | Mission list ve progress gÃ¶rÃ¼nÃ¼mÃ¼nÃ¼ aÃ§ma-kapama anahtarÄ± | `access_availability` | `mission` | `access` | `bool` | `true` | `true,false` | `feature` | `mission.read` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `mission.read`,`mission.progress` | `visibility_off` | `not_found` | `none` | `planned` | GÃ¶rev gÃ¶rÃ¼nÃ¼mÃ¼ claim yÃ¼zeyinden baÄŸÄ±msÄ±z yÃ¶netilebilir |
| `feature.mission.claim.enabled` | Mission claim-request yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `mission` | `access` | `bool` | `true` | `true,false` | `feature` | `mission.claim` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `mission.claim` | `write_off` | `forbidden` | `none` | `planned` | GÃ¶rev ilerlemesi aÃ§Ä±k kalÄ±rken claim-request yÃ¼zeyi geÃ§ici olarak kapatÄ±labilir |
| `feature.royalpass.claim.enabled` | RoyalPass claim-request yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `royalpass` | `access` | `bool` | `true` | `true,false` | `feature` | `royalpass.claim` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `royalpass.claim` | `write_off` | `forbidden` | `none` | `planned` | Claim-request yÃ¼zeyi geÃ§ici olarak kapatÄ±labilir; season ilerleme verisi ayrÄ± yÃ¶netilir |
| `feature.royalpass.season.enabled` | RoyalPass sezon avantajlarÄ±nÄ± sistem genelinde aÃ§ma-kapama anahtarÄ± | `access_availability` | `royalpass` | `access` | `bool` | `true` | `true,false` | `feature` | `royalpass.season` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `royalpass.season`,`royalpass.track`,`royalpass.claim` | `benefit_pause` | `temporarily_unavailable` | `freeze_on_system_disable` | `planned` | Sezon sistem kaynaklÄ± pasifliÄŸe alÄ±ndÄ±ÄŸÄ±nda kalan hak ve sÃ¼re gÃ¼venli biÃ§imde dondurulur |
| `feature.royalpass.premium.enabled` | RoyalPass premium track yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `royalpass` | `access` | `bool` | `true` | `true,false` | `feature` | `royalpass.premium` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `royalpass.premium`,`royalpass.track` | `visibility_off` | `not_found` | `none` | `planned` | Premium entitlement owner'lÄ±ÄŸÄ± deÄŸiÅŸmeden premium gÃ¶rÃ¼nÃ¼rlÃ¼k ayrÄ± yÃ¶netilebilir |
| `feature.history.continue_reading.enabled` | Continue reading yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `history` | `access` | `bool` | `true` | `true,false` | `feature` | `history.continue_reading` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `history.continue_reading` | `visibility_off` | `not_found` | `none` | `planned` | Continue reading yÃ¼zeyi history timeline veya library yÃ¼zeyinden baÄŸÄ±msÄ±z kapatÄ±labilir |
| `feature.history.library.enabled` | KÃ¼tÃ¼phane okuma yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `history` | `access` | `bool` | `true` | `true,false` | `feature` | `history.library` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `history.library` | `visibility_off` | `not_found` | `none` | `planned` | Library read yÃ¼zeyi continue reading ve timeline yÃ¼zeylerinden baÄŸÄ±msÄ±z yÃ¶netilebilir |
| `feature.history.timeline.enabled` | Okuma geÃ§miÅŸi timeline yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `history` | `access` | `bool` | `true` | `true,false` | `feature` | `history.timeline` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `history.timeline` | `visibility_off` | `not_found` | `none` | `planned` | Timeline yÃ¼zeyi continue reading ve library yÃ¼zeylerinden baÄŸÄ±msÄ±z kapatÄ±labilir |
| `feature.history.bookmark_write.enabled` | Bookmark yazma yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `history` | `access` | `bool` | `true` | `true,false` | `feature` | `history.bookmark.write` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `history.bookmark.write` | `write_off` | `forbidden` | `none` | `planned` | Library gÃ¶rÃ¼nÃ¼r kalsa bile bookmark write yÃ¼zeyi ayrÄ±ca kapatÄ±labilir |
| `feature.manga.discovery.enabled` | Recommendation, koleksiyon ve editoryal keÅŸif yÃ¼zeylerini aÃ§ma-kapama anahtarÄ± | `access_availability` | `manga` | `access` | `bool` | `true` | `true,false` | `feature` | `manga.discovery` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `manga.recommendation`,`manga.collection`,`manga.discovery` | `visibility_off` | `not_found` | `none` | `planned` | Editoryal keÅŸif kapatÄ±ldÄ±ÄŸÄ±nda temel manga listing ve detail yÃ¼zeyi ayrÄ± Ã§alÄ±ÅŸabilir |
| `feature.ads.surface.enabled` | Reklam gÃ¶sterim Ã¼st yÃ¼zeylerini aÃ§ma-kapama anahtarÄ± | `access_availability` | `ads` | `access` | `bool` | `true` | `true,false` | `feature` | `ads.surface` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `ads.surface`,`ads.placement`,`ads.delivery` | `visibility_off` | `not_found` | `none` | `planned` | Global ads switch placement ve campaign alt anahtarlarÄ±nÄ± override edebilir |
| `feature.ads.placement.enabled` | Placement Ã§Ã¶zÃ¼mleme yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `ads` | `access` | `bool` | `true` | `true,false` | `feature` | `ads.placement` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `ads.placement` | `visibility_off` | `not_found` | `none` | `planned` | Placement resolve yÃ¼zeyi campaign serve akÄ±ÅŸÄ±ndan baÄŸÄ±msÄ±z daraltÄ±labilir |
| `feature.ads.campaign.enabled` | Campaign serve yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `ads` | `access` | `bool` | `true` | `true,false` | `feature` | `ads.campaign` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `ads.campaign`,`ads.delivery` | `visibility_off` | `not_found` | `none` | `planned` | Placement tanÄ±mlarÄ± kalsa bile campaign serve ayrÄ± kapatÄ±labilir |
| `feature.shop.catalog.enabled` | Shop katalog gÃ¶rÃ¼ntÃ¼leme yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `shop` | `access` | `bool` | `true` | `true,false` | `feature` | `shop.catalog` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `shop.catalog`,`shop.item.view` | `visibility_off` | `not_found` | `none` | `planned` | Katalog yÃ¼zeyi purchase aksiyonundan baÄŸÄ±msÄ±z yÃ¶netilebilir |
| `feature.shop.purchase.enabled` | Shop satÄ±n alma yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `shop` | `access` | `bool` | `true` | `true,false` | `feature` | `shop.purchase` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `shop.purchase` | `write_off` | `forbidden` | `none` | `planned` | Shop katalog gÃ¶rÃ¼nÃ¼r kalÄ±rken satÄ±n alma aksiyonu ayrÄ± kapatÄ±labilir |
| `feature.shop.campaign.enabled` | Shop campaign gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼nÃ¼ aÃ§ma-kapama anahtarÄ± | `access_availability` | `shop` | `access` | `bool` | `true` | `true,false` | `feature` | `shop.campaign` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `shop.campaign`,`shop.catalog` | `visibility_off` | `not_found` | `none` | `planned` | Kampanya badge veya spotlight yÃ¼zeyi katalog okumasÄ±ndan baÄŸÄ±msÄ±z daraltÄ±labilir |
| `feature.payment.mana_purchase.enabled` | Mana package ve purchase entry yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `payment` | `access` | `bool` | `true` | `true,false` | `feature` | `payment.mana_purchase` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `payment.mana_purchase` | `visibility_off` | `not_found` | `none` | `planned` | Checkout baÅŸlatma akÄ±ÅŸÄ± ayrÄ± anahtar ile yÃ¶netilir |
| `feature.payment.checkout.enabled` | Checkout session baÅŸlatma yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `payment` | `access` | `bool` | `true` | `true,false` | `feature` | `payment.checkout` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `payment.checkout`,`payment.provider_session` | `write_off` | `temporarily_unavailable` | `none` | `planned` | Provider sorunu halinde checkout durdurulurken mana package gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ ayrÄ± kalabilir |
| `feature.payment.transaction_read.enabled` | KullanÄ±cÄ±nÄ±n kendi transaction veya wallet gÃ¶rÃ¼nÃ¼mÃ¼nÃ¼ aÃ§ma-kapama anahtarÄ± | `access_availability` | `payment` | `access` | `bool` | `true` | `true,false` | `feature` | `payment.transaction.read` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `payment.transaction.read`,`payment.wallet.read` | `visibility_off` | `not_found` | `none` | `planned` | Ä°ÅŸlem geÃ§miÅŸi gÃ¶rÃ¼nÃ¼mÃ¼ mana purchase ve checkout akÄ±ÅŸlarÄ±ndan baÄŸÄ±msÄ±z daraltÄ±labilir |
| `feature.chapter.early_access.enabled` | Chapter early access gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼nÃ¼ aÃ§ma-kapama anahtarÄ± | `access_availability` | `chapter` | `access` | `bool` | `true` | `true,false` | `feature` | `chapter.early_access` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `chapter.early_access`,`chapter.read` | `visibility_off` | `not_found` | `none` | `planned` | Early access gÃ¶rÃ¼nÃ¼rlÃ¼ÄŸÃ¼ preview ve tam read yÃ¼zeyinden baÄŸÄ±msÄ±z daraltÄ±labilir |
| `feature.notification.digest.enabled` | Notification digest Ã¼retim ve teslim yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `notification` | `service` | `bool` | `true` | `true,false` | `feature` | `notification.digest` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `notification.digest`,`notification.delivery` | `read_only` | `temporarily_unavailable` | `none` | `planned` | Digest kapatÄ±lsa bile tekil in-app inbox aÃ§Ä±k kalabilir |
| `feature.inventory.consume.enabled` | Envanter consume yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `access_availability` | `inventory` | `access` | `bool` | `true` | `true,false` | `feature` | `inventory.consume` | `authenticated` | `-` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `inventory.consume` | `write_off` | `forbidden` | `none` | `planned` | Consume aksiyonu read veya equip yÃ¼zeyinden baÄŸÄ±msÄ±z kapatÄ±labilir |
| `feature.mission.progress_ingest.enabled` | Mission progress event ingest yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `gamification` | `mission` | `service` | `bool` | `true` | `true,false` | `feature` | `mission.progress_ingest` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `mission.progress`,`mission.ingest` | `intake_pause` | `temporarily_unavailable` | `none` | `planned` | Event ingest geÃ§ici kapatÄ±ldÄ±ÄŸÄ±nda mission read yÃ¼zeyi aÃ§Ä±k kalabilir |
| `payment.callback.intake.paused` | Payment provider callback intake yÃ¼zeyini geÃ§ici durdurma anahtarÄ± | `operations` | `payment` | `service` | `bool` | `false` | `true,false` | `feature` | `payment.callback` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `payment.callback`,`payment.reconcile` | `intake_pause` | `temporarily_unavailable` | `none` | `planned` | Provider kaynaklÄ± geri bildirim geÃ§ici olarak durdurulabilir |
| `feature.ads.click_intake.enabled` | Ads click intake yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `operations` | `ads` | `service` | `bool` | `true` | `true,false` | `feature` | `ads.click_intake` | `all` | `-` | `false` | `immediate` | `cache_refresh` | `time_window` | `true` | `ads.click`,`ads.aggregate` | `intake_pause` | `temporarily_unavailable` | `none` | `planned` | Invalid click korumasÄ± iÃ§in click intake ayrÄ± durdurulabilir |
| `feature.support.internal_note.enabled` | Support internal note yazma yÃ¼zeyini aÃ§ma-kapama anahtarÄ± | `support` | `support` | `access` | `bool` | `true` | `true,false` | `feature` | `support.internal_note` | `authenticated` | `role:support_agent` | `false` | `immediate` | `cache_refresh` | `none` | `true` | `support.internal_note` | `write_off` | `forbidden` | `none` | `planned` | Internal note ile requester'a aÃ§Ä±k reply yÃ¼zeyi ayrÄ± yÃ¶netilmelidir |

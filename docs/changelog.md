# Changelog

Bu dosya yalnizca projede gercekte yapilan islemleri kaydeder.
Bu proje SemVer (`MAJOR.MINOR.PATCH`) standardini takip eder.

## [Unreleased]

## [0.13.0-alpha.1] - 2026-03-12

### Added
- `Asama 13` kapsaminda canonical `history` modulu eklendi: `apps/api/internal/modules/history`.
- History owner akis omurgasi eklendi:
  - own continue-reading listeleme
  - own library ve own timeline yuzeyleri
  - controlled public library read yuzeyi (entry-level `share_public`)
  - chapter read signal intake (`chapter.read.started/checkpoint/finished`)
  - admin runtime control (`continue-reading`, `library`, `timeline`, `bookmark-write`)
- History event sabitleri eklendi: `apps/api/internal/modules/history/events/events.go`.
- In-memory history repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/history/repository/memory_store.go`
  - `apps/api/internal/modules/history/repository/checkpoint_repository.go`
  - `apps/api/internal/modules/history/repository/library_repository.go`
  - `apps/api/internal/modules/history/repository/timeline_repository.go`
  - `apps/api/internal/modules/history/repository/runtime_repository.go`
  - `apps/api/internal/modules/history/repository/snapshot_store.go`
  - `apps/api/internal/modules/history/repository/memory_store_test.go`
- History service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/history/service/history_read_service.go`
  - `apps/api/internal/modules/history/service/history_update_service.go`
  - `apps/api/internal/modules/history/service/history_admin_service.go`
  - `apps/api/internal/modules/history/service/history_intake_service.go`
  - `apps/api/internal/modules/history/service/service_test.go`
- History HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/history/handler/*`
  - `apps/api/internal/modules/history/routes.go`
- History migration cifti eklendi:
  - `apps/api/migrations/202603120011_history_create_core_tables.up.sql`
  - `apps/api/migrations/202603120011_history_create_core_tables.down.sql`
- History stage testleri eklendi:
  - contract: `apps/api/tests/contract/history_events_contract_test.go`
  - integration: `apps/api/tests/integration/history_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/history_migration_integration_test.go`

### Changed
- API bootstrap'ta history module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- Snapshot persistence hedeflerine `history` eklendi.
- Chapter modulune history icin `GetResumeAnchor` ve `BuildReadSignal` kontrat ciktilari eklendi.
- `docs/modules.md` modul envanterinde `history` status'u `active` olarak guncellendi.
- `docs/shared.md` icindeki history feature key kayitlari `active` duruma cekildi.
- `VERSION`, `.env.example`, `README.md`, `docs/TESTING.md` ve `docs/upgrade.md` Asama 13 ile hizalandi.

### Release Notes
- Degisiklik Ozeti: Asama 13 history owner (continue-reading/library/timeline/intake/runtime-control) omurgasi kod seviyesine tasindi ve chapter->history kontrat zinciri aktif hale getirildi.
- Etkilenen Moduller: `history`, `chapter`, `app`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120011_history_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.12.0-alpha.1] - 2026-03-12

### Added
- `Asama 12` kapsaminda canonical `notification` modulu eklendi: `apps/api/internal/modules/notification`.
- Notification owner akis omurgasi eklendi:
  - own inbox listeleme + detail
  - read/unread lifecycle (`read` islemi)
  - own preference guncelleme (mute, quiet-hour, channel, digest)
  - support event intake ile bildirim olusturma
  - admin runtime control (`category-state`, `channel-state`, `digest-state`, `delivery-pause`)
- Notification event sabitleri eklendi: `apps/api/internal/modules/notification/events/events.go`.
- In-memory notification repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/notification/repository/memory_store.go`
  - `apps/api/internal/modules/notification/repository/notification_repository.go`
  - `apps/api/internal/modules/notification/repository/preference_repository.go`
  - `apps/api/internal/modules/notification/repository/admin_repository.go`
  - `apps/api/internal/modules/notification/repository/snapshot_store.go`
  - `apps/api/internal/modules/notification/repository/memory_store_test.go`
- Notification service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/notification/service/notification_inbox_service.go`
  - `apps/api/internal/modules/notification/service/notification_preference_service.go`
  - `apps/api/internal/modules/notification/service/notification_admin_service.go`
  - `apps/api/internal/modules/notification/service/notification_intake_service.go`
  - `apps/api/internal/modules/notification/service/service_test.go`
- Notification HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/notification/handler/*`
  - `apps/api/internal/modules/notification/routes.go`
- Notification migration cifti eklendi:
  - `apps/api/migrations/202603120010_notification_create_core_tables.up.sql`
  - `apps/api/migrations/202603120010_notification_create_core_tables.down.sql`
- Notification stage testleri eklendi:
  - contract: `apps/api/tests/contract/notification_events_contract_test.go`
  - integration: `apps/api/tests/integration/notification_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/notification_migration_integration_test.go`

### Changed
- API bootstrap'ta notification module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- Snapshot persistence hedeflerine `notification` eklendi.
- Support modulune `BuildNotificationSignal` kontrat ciktisi eklendi ve `notification` ile composition root'ta baglandi.
- `docs/modules.md` modul envanterinde `notification` status'u `active` olarak guncellendi.
- `docs/shared.md` icindeki notification feature/runtime key kayitlari `active` duruma cekildi.
- `VERSION`, `.env.example`, `README.md`, `docs/TESTING.md` ve `docs/upgrade.md` Asama 12 ile hizalandi.

### Release Notes
- Degisiklik Ozeti: Asama 12 notification owner (inbox/preference/intake/runtime-control) omurgasi kod seviyesine tasindi ve support->notification kontrat zinciri aktif hale getirildi.
- Etkilenen Moduller: `notification`, `support`, `app`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120010_notification_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.11.0-alpha.1] - 2026-03-12

### Added
- `Asama 11` kapsaminda canonical `moderation` modulu eklendi: `apps/api/internal/modules/moderation`.
- Moderation owner akis omurgasi eklendi:
  - support handoff'tan linked case olusturma
  - scoped queue listeleme ve case detail
  - moderator assignment/release
  - moderator note timeline
  - sinirli aksiyon uygulama (`hide`, `unhide`, `lock`, `unlock`, `warning`, `review_complete`)
  - admin handoff/escalation
- Moderation event sabitleri eklendi: `apps/api/internal/modules/moderation/events/events.go`.
- In-memory moderation repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/moderation/repository/memory_store.go`
  - `apps/api/internal/modules/moderation/repository/case_repository.go`
  - `apps/api/internal/modules/moderation/repository/memory_store_test.go`
- Moderation service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/moderation/service/moderation_create_service.go`
  - `apps/api/internal/modules/moderation/service/moderation_query_service.go`
  - `apps/api/internal/modules/moderation/service/moderation_assignment_service.go`
  - `apps/api/internal/modules/moderation/service/moderation_note_service.go`
  - `apps/api/internal/modules/moderation/service/moderation_action_service.go`
  - `apps/api/internal/modules/moderation/service/moderation_escalation_service.go`
  - `apps/api/internal/modules/moderation/service/service_test.go`
- Moderation HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/moderation/handler/*`
  - `apps/api/internal/modules/moderation/routes.go`
- Moderation migration cifti eklendi:
  - `apps/api/migrations/202603120009_moderation_create_core_tables.up.sql`
  - `apps/api/migrations/202603120009_moderation_create_core_tables.down.sql`
- Moderation stage testleri eklendi:
  - contract: `apps/api/tests/contract/moderation_events_contract_test.go`
  - integration: `apps/api/tests/integration/moderation_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/moderation_migration_integration_test.go`

### Changed
- API bootstrap'ta moderation module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- Support modulu `GetModerationHandoffReference` yaninda linked-case baglama kontrati (`LinkModerationCase`) saglayacak sekilde genisletildi.
- `docs/modules.md` modul envanterinde `moderation` status'u `active` olarak guncellendi.
- `VERSION`, `.env.example`, `README.md`, `docs/TESTING.md` ve `docs/upgrade.md` Asama 11 ile hizalandi.

### Docs
- `docs/shared.md` moderation lifecycle sozlugune `escalation_status` bolumu eklendi.
- `docs/shared.md` icindeki moderation feature key kayitlari `active` duruma cekildi.

### Release Notes
- Degisiklik Ozeti: Asama 11 moderation owner (queue/case/assignment/action/escalation) omurgasi kod seviyesine tasindi ve support->moderation linked-case akisi aktif hale getirildi.
- Etkilenen Moduller: `moderation`, `support`, `app`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120009_moderation_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.10.0-alpha.3] - 2026-03-12

### Changed
- Memory tabanli modul state'i icin kalici snapshot altyapisi eklendi; uygulama acilisinda restore, calisma sirasinda periyodik persist ve kapanista son persist akisi bootstrap'e baglandi.
- Access evaluate HTTP contract'i actor-context guven modelini koruyacak sekilde geriye uyumlu policy alanlariyla (`scope_kind`, `scope_selector`, `audience_selector`, `allow_override`) guncellendi.
- CI akisina dokuman encoding kontrol adimi eklendi (`scripts/check_utf8_no_bom.ps1`).

### Fixed
- ISS-004: Mimari persistence uyumsuzlugu snapshot store katmani ile kapatildi (Postgres ve file fallback).
- ISS-011: Access/support tarafindaki yari kalmis tekrarli yapilar sadelelestirildi; kullanilmayan `errors.go` dosyalari kaldirildi.
- ISS-013: Eksik e2e katmani tamamlandi (`apps/api/tests/e2e/critical_flow_e2e_test.go`).
- ISS-016: `normalize_utf8_no_bom.ps1` script hatasi giderildi ve dokumanlar UTF-8/BOM-suz normalize edildi.

### Docs
- `docs/issues.md` icindeki kalan acik bulgular temizlenerek acik bulgu sayisi sifira indirildi.

### Release Notes
- Degisiklik Ozeti: Kalan denetim bulgulari (ISS-004, ISS-011, ISS-013, ISS-016) kapatildi; test, CI ve docker dogrulama adimlariyla hotfix paketi yayina hazirlandi.
- Etkilenen Moduller: `app`, `platform/config`, `platform/snapshot`, `auth`, `user`, `access`, `manga`, `chapter`, `comment`, `support`, `tests`, `docs`, `ci`.
- Breaking Change: Yok.
- Migration Etkisi: Yok.

## [0.10.0-alpha.2] - 2026-03-12

### Changed
- HTTP actor kimligi ve rol guveni middleware/context tabanina tasindi (`X-Actor-*` basliklari) ve ilgili route guardlari aktive edildi.
- Auth, user, comment, support ve access handlerlari istemci tarafindan gelen actor/credential alanlari yerine server-side context kaynakli alanlari kullanacak sekilde duzenlendi.
- Moduller arasi bagimlilik yonleri netlestirildi; `auth->user`, `manga->chapter`, `manga/chapter->comment`, `manga/chapter/comment->support` dogrulama kontratlari composition root'ta baglandi.
- Dokuman referanslari mevcut tek dosya yerlesimine hizalandi (`docs/modules.md`, `docs/shared.md`) ve branch modeli kurali CI gercegi ile uyumlu hale getirildi (`codex/**` notu).

### Fixed
- Access evaluate endpointinde guven siniri sertlestirildi; `user_id`, `identity`, `user_signal` istemciden alinmiyor.
- 500 sinifindaki hata cevabi sizintilari `internal_error` defaultuna cekildi.
- Kritik persistence hatasi yutma noktalarinda login/security event akislari guvenli hale getirildi.
- Access bootstrap seed hatalari startup asamasinda fail-fast olacak sekilde duzeltildi.
- Config yukleme zorunluluklari sadelestirildi; runtime icin `DB_MAIN_DSN`/`DB_TEST_DSN` zorunlu degil.
- Chapter navigation akisi hard-limit tarama yerine hedefli komsu cozumu ile olceklenebilir hale getirildi.
- Integration test helperlari yeni auth/actor kontratina gore guncellendi.

### Docs
- `docs/issues.md` temizlendi; cozulen bulgular kayittan cikarildi, acik kalanlar ayrik listelendi.
- `docs/issue.md` kaldirildi ve tek issue kaynagi olarak `docs/issues.md` birakildi.

### Release Notes
- Degisiklik Ozeti: Guven siniri ve katman kurallari hotfix paketinde sertlestirildi; test ve docker dogrulamalariyla birlikte yayinlanabilir durum saglandi.
- Etkilenen Moduller: `auth`, `user`, `access`, `comment`, `support`, `manga`, `chapter`, `app`, `platform/config`, `docs`, `tests`.
- Breaking Change: HTTP actor/credential kimlik alanlari body/query yerine `X-Actor-*` headerlari ile alinmaktadir.
- Migration Etkisi: Yok.
## [0.10.0-alpha.1] - 2026-03-12

### Added
- `Asama 10` kapsaminda canonical `support` modulu eklendi: `apps/api/internal/modules/support`.
- Support owner akis omurgasi eklendi:
  - communication/ticket/report create
  - own support list + detail
  - public reply ve internal note ayrimi
  - review queue, status update ve resolve akislari
  - duplicate/idempotency (`request_id`) ve spam risk sinyali
  - report -> moderation handoff istegi
- `support -> moderation` ve `support -> notification` kontrat yuzeyleri eklendi:
  - `apps/api/internal/modules/support/contract/moderation_handoff_contract.go`
  - `apps/api/internal/modules/support/contract/notification_contract.go`
- Support event sabitleri eklendi: `apps/api/internal/modules/support/events/events.go`.
- In-memory support repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/support/repository/memory_store.go`
  - `apps/api/internal/modules/support/repository/support_repository.go`
  - `apps/api/internal/modules/support/repository/memory_store_test.go`
- Support service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/support/service/support_create_service.go`
  - `apps/api/internal/modules/support/service/support_query_service.go`
  - `apps/api/internal/modules/support/service/support_update_service.go`
  - `apps/api/internal/modules/support/service/support_contract_service.go`
  - `apps/api/internal/modules/support/service/service_test.go`
- Support HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/support/handler/*`
  - `apps/api/internal/modules/support/routes.go`
- Support migration cifti eklendi:
  - `apps/api/migrations/202603120008_support_create_core_tables.up.sql`
  - `apps/api/migrations/202603120008_support_create_core_tables.down.sql`
- Support stage testleri eklendi:
  - contract: `apps/api/tests/contract/support_moderation_contract_test.go`, `apps/api/tests/contract/support_notification_contract_test.go`
  - integration: `apps/api/tests/integration/support_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/support_migration_integration_test.go`

### Changed
- API bootstrap'ta support module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- `docs/modules.md` modul envanterinde `support` status'u `active` olarak guncellendi.
- `VERSION`, `.env.example`, `README.md`, `docs/TESTING.md` ve `docs/upgrade.md` Asama 10 ile hizalandi.
- Support modulu katmanlari coklu akis tasiyan tek dosyalari onleyecek sekilde islem bazli parcalandi (`routes.go` tek giris noktasi olarak korunarak).

### Fixed
- Yok.

### Docs
- Asama 10 support omurgasi ve versiyonlama guncellemeleri changelog ile izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 10 support owner (intake/review/reply/resolve/handoff) omurgasi kod seviyesine tasindi.
- Etkilenen Moduller: `support`, `app`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120008_support_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.9.0-alpha.1] - 2026-03-12

### Added
- `Asama 9` kapsaminda canonical `comment` modulu eklendi: `apps/api/internal/modules/comment`.
- Comment owner akis omurgasi eklendi:
  - comment create/edit/delete
  - target bazli root listing ve detail
  - thread/reply akisi (root + nested reply)
  - moderation state (`visible`, `hidden`, `flagged`) + spoiler/pin/lock/shadowban alanlari
  - soft delete gorunumu ve restore akisi
  - write cooldown, reply depth limiti ve edit/restore window kurallari
- `comment -> moderation/support` kontrat yuzeyi eklendi: `apps/api/internal/modules/comment/contract/comment_target_contract.go`.
- Comment event sabitleri eklendi: `apps/api/internal/modules/comment/events/events.go`.
- In-memory comment repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/comment/repository/memory_store.go`
  - `apps/api/internal/modules/comment/repository/comment_repository.go`
  - `apps/api/internal/modules/comment/repository/memory_store_test.go`
- Comment service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/comment/service/comment_create_service.go`
  - `apps/api/internal/modules/comment/service/comment_update_service.go`
  - `apps/api/internal/modules/comment/service/comment_query_service.go`
  - `apps/api/internal/modules/comment/service/comment_lifecycle_service.go`
  - `apps/api/internal/modules/comment/service/comment_moderation_service.go`
  - `apps/api/internal/modules/comment/service/comment_contract_service.go`
  - `apps/api/internal/modules/comment/service/service_test.go`
- Comment HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/comment/handler/*`
  - `apps/api/internal/modules/comment/routes.go`
- Comment migration cifti eklendi:
  - `apps/api/migrations/202603120007_comment_create_core_tables.up.sql`
  - `apps/api/migrations/202603120007_comment_create_core_tables.down.sql`
- Comment stage testleri eklendi:
  - contract: `apps/api/tests/contract/comment_target_contract_test.go`
  - integration: `apps/api/tests/integration/comment_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/comment_migration_integration_test.go`

### Changed
- API bootstrap'ta comment module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- `docs/modules.md` modul envanterinde `comment` status'u `active` olarak guncellendi.
- `VERSION`, `.env.example`, `README.md`, `docs/TESTING.md` ve `docs/upgrade.md` Asama 9 ile hizalandi.
- Comment modulu katmanlari coklu akis tasiyan tek dosyalari onleyecek sekilde islem bazli parcalandi (`routes.go` tek giris noktasi olarak korunarak).

### Fixed
- Yok.

### Docs
- Asama 9 comment omurgasi ve versiyonlama guncellemeleri changelog ile izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 9 comment owner (thread/reply/moderation/soft-delete) omurgasi kod seviyesine tasindi.
- Etkilenen Moduller: `comment`, `app`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120007_comment_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.8.0-alpha.1] - 2026-03-12

### Added
- `Asama 8` kapsaminda canonical `chapter` modulu eklendi: `apps/api/internal/modules/chapter`.
- Chapter owner akis omurgasi eklendi:
  - chapter create/update
  - manga bazli chapter listing + detail
  - read (`preview`, `full`) ve navigation (`previous`, `next`, `first`, `last`)
  - publish lifecycle (`draft`, `scheduled`, `published`, `archived`, `unpublished`)
  - read/early-access state (`read_access_level`, `vip_only`, `early_access_*`, fallback access)
  - media health ve integrity state guncellemeleri
  - soft delete + restore
- `chapter -> history` kontrat yuzeyi eklendi: `apps/api/internal/modules/chapter/contract/history_contract.go`.
- Chapter event sabitleri eklendi: `apps/api/internal/modules/chapter/events/events.go`.
- In-memory chapter repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/chapter/repository/memory_store.go`
  - `apps/api/internal/modules/chapter/repository/chapter_repository.go`
  - `apps/api/internal/modules/chapter/repository/memory_store_test.go`
- Chapter service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/chapter/service/chapter_create_service.go`
  - `apps/api/internal/modules/chapter/service/chapter_update_service.go`
  - `apps/api/internal/modules/chapter/service/chapter_query_service.go`
  - `apps/api/internal/modules/chapter/service/chapter_read_service.go`
  - `apps/api/internal/modules/chapter/service/chapter_publish_service.go`
  - `apps/api/internal/modules/chapter/service/chapter_access_service.go`
  - `apps/api/internal/modules/chapter/service/chapter_media_service.go`
  - `apps/api/internal/modules/chapter/service/chapter_lifecycle_service.go`
  - `apps/api/internal/modules/chapter/service/chapter_contract_service.go`
  - `apps/api/internal/modules/chapter/service/service_test.go`
- Chapter HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/chapter/handler/*`
  - `apps/api/internal/modules/chapter/routes.go`
- Chapter migration cifti eklendi:
  - `apps/api/migrations/202603120006_chapter_create_core_tables.up.sql`
  - `apps/api/migrations/202603120006_chapter_create_core_tables.down.sql`
- Chapter stage testleri eklendi:
  - contract: `apps/api/tests/contract/chapter_history_contract_test.go`
  - integration: `apps/api/tests/integration/chapter_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/chapter_migration_integration_test.go`

### Changed
- API bootstrap'ta chapter module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- `docs/modules.md` modul envanterinde `chapter` status'u `active` olarak guncellendi.
- `VERSION`, `.env.example`, `README.md`, `docs/TESTING.md` ve `docs/upgrade.md` Asama 8 ile hizalandi.
- Chapter modulu katmanlari coklu akis tasiyan tek dosyalari onleyecek sekilde islem bazli parcalandi (`routes.go` tek giris noktasi olarak korunarak).

### Fixed
- `apps/api/internal/modules/chapter/module_test.go` route mount testi `/manga/{manga_id}/chapters` yuzeyini dogrulayacak sekilde duzeltildi.

### Docs
- Asama 8 chapter omurgasi ve versiyonlama guncellemeleri changelog ile izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 8 chapter owner (chapter/page/read/navigation/release/early-access) omurgasi kod seviyesine tasindi.
- Etkilenen Moduller: `chapter`, `app`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120006_chapter_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.7.0-alpha.1] - 2026-03-12

### Added
- `Asama 7` kapsaminda canonical `manga` modulu eklendi: `apps/api/internal/modules/manga`.
- Manga owner akis omurgasi eklendi:
  - manga create/update
  - public listing + detail
  - search/filter/sort
  - publish lifecycle (`draft`, `scheduled`, `published`, `archived`, `unpublished`)
  - visibility (`public`, `hidden`) ve editorial/discovery (`featured`, `recommended`, `collection`)
  - denormalize counter sync (`chapter_count`, `comment_count`, `view_count`)
  - soft delete + restore
- `manga -> chapter` default access kontrat yuzeyi eklendi: `apps/api/internal/modules/manga/contract/chapter_defaults_contract.go`.
- Manga event sabitleri eklendi: `apps/api/internal/modules/manga/events/events.go`.
- In-memory manga repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/manga/repository/memory_store.go`
  - `apps/api/internal/modules/manga/repository/manga_repository.go`
  - `apps/api/internal/modules/manga/repository/memory_store_test.go`
- Manga service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/manga/service/manga_create_service.go`
  - `apps/api/internal/modules/manga/service/manga_update_service.go`
  - `apps/api/internal/modules/manga/service/manga_listing_service.go`
  - `apps/api/internal/modules/manga/service/manga_publish_service.go`
  - `apps/api/internal/modules/manga/service/manga_discovery_service.go`
  - `apps/api/internal/modules/manga/service/manga_lifecycle_service.go`
  - `apps/api/internal/modules/manga/service/service_test.go`
- Manga HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/manga/handler/*`
  - `apps/api/internal/modules/manga/routes.go`
- Manga migration cifti eklendi:
  - `apps/api/migrations/202603120005_manga_create_core_tables.up.sql`
  - `apps/api/migrations/202603120005_manga_create_core_tables.down.sql`
- Manga stage testleri eklendi:
  - contract: `apps/api/tests/contract/manga_access_contract_test.go`
  - integration: `apps/api/tests/integration/manga_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/manga_migration_integration_test.go`

### Changed
- API bootstrap'ta manga module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- Access canonical permission setine `manga.list.view` eklendi ve default role seed'lerine baglandi.
- `docs/modules.md` modul envanterinde `manga` status'u `active` olarak guncellendi.
- `VERSION`, `.env.example`, `README.md`, `docs/TESTING.md` ve `docs/upgrade.md` Asama 7 ile hizalandi.
- Manga modulu katmanlari coklu akis tasiyan tek dosyalari onleyecek sekilde islem bazli parcalandi (`routes.go` tek giris noktasi olarak korunarak).

### Fixed
- Yok.

### Docs
- Asama 7 manga omurgasi ve versiyonlama guncellemeleri changelog ile izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 7 manga icerik owner (metadata, listing/detail, discovery, lifecycle) omurgasi kod seviyesine tasindi.
- Etkilenen Moduller: `manga`, `access`, `app`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120005_manga_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.6.0-alpha.1] - 2026-03-12

### Added
- `Asama 6` kapsaminda canonical `access` modulu eklendi: `apps/api/internal/modules/access`.
- Access owner akis omurgasi eklendi:
  - role create/list
  - permission create/list
  - role-permission assignment
  - user-role assignment (sureli atama dahil)
  - temporary grant create/revoke
  - policy create/list ve conflict kontrolu
  - evaluate endpoint'i ile final allow/deny karar katmani
- `access` modul kontrat yuzeyi eklendi: `apps/api/internal/modules/access/contract/access_contract.go`.
- Access event sabitleri eklendi: `apps/api/internal/modules/access/events/events.go`.
- In-memory access repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/access/repository/memory_store.go`
  - `apps/api/internal/modules/access/repository/access_role_repository.go`
  - `apps/api/internal/modules/access/repository/access_permission_repository.go`
  - `apps/api/internal/modules/access/repository/access_assignment_repository.go`
  - `apps/api/internal/modules/access/repository/access_policy_repository.go`
  - `apps/api/internal/modules/access/repository/memory_store_test.go`
- Access service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/access/service/access_role_service.go`
  - `apps/api/internal/modules/access/service/access_permission_service.go`
  - `apps/api/internal/modules/access/service/access_assignment_service.go`
  - `apps/api/internal/modules/access/service/access_policy_service.go`
  - `apps/api/internal/modules/access/service/access_evaluate_service.go`
  - `apps/api/internal/modules/access/service/service_test.go`
- Access HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/access/handler/*`
  - `apps/api/internal/modules/access/routes.go`
- Access migration cifti eklendi:
  - `apps/api/migrations/202603120004_access_create_core_tables.up.sql`
  - `apps/api/migrations/202603120004_access_create_core_tables.down.sql`
- Access stage testleri eklendi:
  - contract: `apps/api/tests/contract/access_authorization_contract_test.go`
  - integration: `apps/api/tests/integration/access_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/access_migration_integration_test.go`

### Changed
- API bootstrap'ta access module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- `docs/modules.md` modul envanterinde `access` status'u `active` olarak guncellendi.
- `VERSION`, `.env.example`, `README.md`, `docs/TESTING.md` ve `docs/upgrade.md` Asama 6 ile hizalandi.
- Access modulu katmanlari coklu akis tasiyan tek dosyalari onleyecek sekilde islem bazli parcalandi (`routes.go` tek giris noktasi olarak korunarak).

### Fixed
- Yok.

### Docs
- Asama 6 access omurgasi ve versiyonlama guncellemeleri changelog ile izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 6 merkezi authorization/policy/degerlendirme omurgasi kod seviyesine tasindi.
- Etkilenen Moduller: `access`, `app`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120004_access_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.5.1-alpha.1] - 2026-03-12

### Changed
- `auth` modulu route kayitlari tek dosya ilkesine uygun olacak sekilde modul kokundeki `apps/api/internal/modules/auth/routes.go` icine toplandi.
- `user` modulu route kayitlari tek dosya ilkesine uygun olacak sekilde modul kokundeki `apps/api/internal/modules/user/routes.go` icine toplandi.
- Route kayit katmaninda bagimlilik kontrolu fail-fast olacak sekilde guncellendi (`router` veya `httpHandler` nil ise panic).

### Fixed
- Kural 7 uyumsuzlugu giderildi: route kayitlarinin harici `routes/` klasorunde parcalanmasi kaldirildi, modul kok `routes.go` standardina donuldu.
- Wiring hatalarinda endpointlerin sessizce mount edilmemesi (no-op) sorunu giderildi.

### Removed
- Asagidaki parcali route dosyalari kaldirildi:
  - `apps/api/internal/modules/auth/routes/auth_identity_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_password_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_session_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_verification_routes.go`
  - `apps/api/internal/modules/auth/routes/register_routes.go`
  - `apps/api/internal/modules/user/routes/register_routes.go`
  - `apps/api/internal/modules/user/routes/user_account_routes.go`
  - `apps/api/internal/modules/user/routes/user_profile_routes.go`

### Docs
- `VERSION`, `.env.example` ve `README.md` dosyalari yeni canonical surum `0.5.1-alpha.1` ile hizalandi.

### Release Notes
- Degisiklik Ozeti: Asama 4-5 route katmani modul kok `routes.go` standardina cekildi ve fail-fast guard'lar eklendi.
- Etkilenen Moduller: `auth`, `user`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: Yok.
- Operasyon Notu: Runtime versiyonu icin `APP_VERSION=0.5.1-alpha.1` kullanilmalidir.

## [0.5.0-alpha.1] - 2026-03-12

### Added
- `Asama 5` kapsaminda canonical `user` modulu eklendi: `apps/api/internal/modules/user`.
- User owner akis omurgasi eklendi:
  - user create ve profile read/update
  - public/private profile response ayrimi
  - profile visibility ve history visibility preference guncelleme
  - account state gecisleri (`active`, `deactivated`, `banned`)
  - VIP lifecycle (`activate`, `freeze`, `resume`, `deactivate`)
- `user -> access` kontrat yuzeyi eklendi: `apps/api/internal/modules/user/contract/access_contract.go`.
- User event sabitleri eklendi: `apps/api/internal/modules/user/events/events.go`.
- In-memory user repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/user/repository/memory_store.go`
  - `apps/api/internal/modules/user/repository/user_account_repository.go`
  - `apps/api/internal/modules/user/repository/memory_store_test.go`
- User service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/user/service/user_create_service.go`
  - `apps/api/internal/modules/user/service/user_profile_service.go`
  - `apps/api/internal/modules/user/service/user_visibility_service.go`
  - `apps/api/internal/modules/user/service/user_account_service.go`
  - `apps/api/internal/modules/user/service/user_vip_service.go`
  - `apps/api/internal/modules/user/service/service_test.go`
- User HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/user/handler/*`
  - `apps/api/internal/modules/user/routes/*`
- User migration cifti eklendi:
  - `apps/api/migrations/202603120003_user_create_core_tables.up.sql`
  - `apps/api/migrations/202603120003_user_create_core_tables.down.sql`
- User stage testleri eklendi:
  - contract: `apps/api/tests/contract/user_access_contract_test.go`
  - integration: `apps/api/tests/integration/user_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/user_migration_integration_test.go`

### Changed
- API bootstrap'ta user module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- `docs/modules.md` modul envanterinde `user` status'u `active` olarak guncellendi.
- `VERSION`, `.env.example`, `README.md`, `docs/TESTING.md` ve `docs/upgrade.md` Asama 5 ile hizalandi.
- User modulu katmanlari coklu akislari tek dosyada toplamayacak sekilde islem bazli parcalandi.

### Fixed
- Yok.

### Docs
- Asama 5 user omurgasi ve versiyonlama guncellemeleri changelog ile izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 5 user owner omurgasi (account/profile/visibility/membership/vip) kod seviyesine tasindi.
- Etkilenen Moduller: `user`, `app`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120003_user_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

## [0.4.0-alpha.1] - 2026-03-12

### Added
- `Asama 4` kapsaminda canonical `auth` modulu eklendi: `apps/api/internal/modules/auth`.
- Auth module owner akis omurgasi eklendi:
  - register/login/logout
  - session list + revoke (`current`, `others`, `all`)
  - refresh token rotation
  - forgot/reset/change password
  - email verification send/confirm + resend cooldown
  - failed login limit + cooldown davranisi
- `auth -> access` kontrat yuzeyi eklendi: `apps/api/internal/modules/auth/contract/identity_contract.go`.
- Auth event sabitleri eklendi: `apps/api/internal/modules/auth/events/events.go`.
- In-memory auth repository omurgasi ve testi eklendi:
  - `apps/api/internal/modules/auth/repository/memory_store.go`
  - `apps/api/internal/modules/auth/repository/auth_credential_repository.go`
  - `apps/api/internal/modules/auth/repository/auth_session_repository.go`
  - `apps/api/internal/modules/auth/repository/auth_token_repository.go`
  - `apps/api/internal/modules/auth/repository/auth_security_event_repository.go`
  - `apps/api/internal/modules/auth/repository/memory_store_test.go`
- Auth service use-case omurgasi ve kapsam testleri eklendi:
  - `apps/api/internal/modules/auth/service/auth_service.go`
  - `apps/api/internal/modules/auth/service/auth_register_service.go`
  - `apps/api/internal/modules/auth/service/auth_login_service.go`
  - `apps/api/internal/modules/auth/service/auth_refresh_token_service.go`
  - `apps/api/internal/modules/auth/service/auth_session_service.go`
  - `apps/api/internal/modules/auth/service/auth_password_service.go`
  - `apps/api/internal/modules/auth/service/auth_verification_service.go`
  - `apps/api/internal/modules/auth/service/auth_token_issue_service.go`
  - `apps/api/internal/modules/auth/service/auth_security_event_service.go`
  - `apps/api/internal/modules/auth/service/auth_request_meta_service.go`
  - `apps/api/internal/modules/auth/service/auth_session_parse_service.go`
  - `apps/api/internal/modules/auth/service/service_test.go`
- Auth HTTP handler ve route omurgasi eklendi:
  - `apps/api/internal/modules/auth/handler/auth_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_register_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_login_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_token_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_session_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_password_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_verification_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_request_meta_handler.go`
  - `apps/api/internal/modules/auth/handler/auth_response_handler.go`
- Auth route registration yapisi islem ailelerine gore parcalandi:
  - `apps/api/internal/modules/auth/routes.go`
  - `apps/api/internal/modules/auth/routes/register_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_identity_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_session_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_password_routes.go`
  - `apps/api/internal/modules/auth/routes/auth_verification_routes.go`
- Auth migration cifti eklendi:
  - `apps/api/migrations/202603120002_auth_create_core_tables.up.sql`
  - `apps/api/migrations/202603120002_auth_create_core_tables.down.sql`
- Auth stage testleri eklendi:
  - contract: `apps/api/tests/contract/auth_access_contract_test.go`
  - integration: `apps/api/tests/integration/auth_http_integration_test.go`
  - migration smoke: `apps/api/tests/integration/auth_migration_integration_test.go`
- Canonical validation ve password altyapisi eklendi:
  - `apps/api/internal/platform/validation/validator.go`
  - `apps/api/internal/shared/crypto/password/argon2id.go`

### Changed
- API bootstrap'ta auth module registry'ye baglandi: `apps/api/cmd/api/main.go`.
- Config'e auth runtime esikleri eklendi:
  - `AUTH_LOGIN_FAILED_ATTEMPT_LIMIT_PER_MINUTE`
  - `AUTH_LOGIN_COOLDOWN_SECONDS`
  - `AUTH_EMAIL_VERIFICATION_RESEND_COOLDOWN_SECONDS`
- `.env.example`, `README.md` ve `docs/upgrade.md` Asama 4 ile hizalandi.
- Auth modulu icinde coklu akis tasiyan tek dosyalar parcalama kuralina gore operasyon bazli ayrildi (`service.go` ve `http_handler.go` kaldirildi, route kayitlari `routes/` altina ayrildi).

### Fixed
- Yok.

### Docs
- Asama 4 auth omurgasi ve versiyonlama guncellemeleri changelog ile izlenebilir hale getirildi.

### Release Notes
- Degisiklik Ozeti: Asama 4 auth owner omurgasi (credential/session/token/verification/recovery) kod seviyesine tasindi.
- Etkilenen Moduller: `auth`, `platform/validation`, `shared/crypto/password`, `config`, `migrations`, `tests`, `docs`.
- Breaking Change: Yok.
- Migration Etkisi: `202603120002_auth_create_core_tables` migration cifti eklendi (uyumlu schema genislemesi).

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

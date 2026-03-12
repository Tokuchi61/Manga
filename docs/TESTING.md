# TESTING

Bu dokuman test katmanlarini ve Asama 0-16 dogrulama adimlarini listeler.

## Katmanlar

- Unit: config, servis ve saf is kurallari
- Integration: repository ve kritik HTTP akislarinin dogrulamasi
- Contract: moduller arasi acik kontratlarin dogrulamasi
- E2E: kritik uc-tan-uca akislar

## Asama 0 Temel Kontroller

- `cd apps/api && go test ./...`
- `cd apps/api && go build ./...`
- `docker build -f apps/api/Dockerfile -t novascans-api:local .`
- `docker compose -f deploy/docker-compose.yml up -d --build`
- `GET /health` ve `GET /version`

## Asama 1 Mimari Kontroller

- `cd apps/api && go test ./internal/modules/...`
- `cd apps/api && go test ./tests/contract/...`
- `cd apps/api && go test ./tests/integration/...`
- `cd apps/api && go test ./tests/e2e/...`
- Module registry duplicate-name ve canonical-name kurallarinin test edilmesi
- App bootstrap katmaninda module route mount davranisinin test edilmesi

## Asama 4 Auth Kontrolleri

- `cd apps/api && go test ./internal/modules/auth/...`
- `cd apps/api && go test ./tests/contract -run Auth`
- `cd apps/api && go test ./tests/integration -run Auth`
- Register/login/logout/session revoke ve token rotation akislarinin dogrulanmasi
- Forgot/reset/change password ile verification/resend cooldown senaryolarinin dogrulanmasi
- Auth migration smoke kontrolu (`202603120002_auth_create_core_tables`)

## Asama 5 User Kontrolleri

- `cd apps/api && go test ./internal/modules/user/...`
- `cd apps/api && go test ./tests/contract -run User`
- `cd apps/api && go test ./tests/integration -run User`
- Profil okuma-guncelleme ve public/private response ayriminin dogrulanmasi
- History visibility preference ust sinir (global deny) davranisinin dogrulanmasi
- Account state (deactivated/banned) ve VIP lifecycle (activate/freeze/resume/deactivate) senaryolarinin dogrulanmasi
- User migration smoke kontrolu (`202603120003_user_create_core_tables`)

## Asama 6 Access Kontrolleri

- `cd apps/api && go test ./internal/modules/access/...`
- `cd apps/api && go test ./tests/contract -run Access`
- `cd apps/api && go test ./tests/integration -run Access`
- Role/permission tanimlari, role-permission baglama ve user-role assignment akislarinin dogrulanmasi
- Policy create/update davranislarinda conflict kontrolu ve precedence yorumunun dogrulanmasi
- Guest/authenticated/vip/blocked ayrimlari ile own/any decision matrix senaryolarinin dogrulanmasi
- Access migration smoke kontrolu (`202603120004_access_create_core_tables`)

## Asama 7 Manga Kontrolleri

- `cd apps/api && go test ./internal/modules/manga/...`
- `cd apps/api && go test ./tests/contract -run Manga`
- `cd apps/api && go test ./tests/integration -run Manga`
- Manga create/update/list/detail akislarinin ve search/filter/sort davranisinin dogrulanmasi
- Publish/archive/unpublish/schedule lifecycle akislari ile visibility daraltma davranisinin dogrulanmasi
- Discovery/editorial (featured/recommended/collection) ve counter sync (chapter/comment/view) dogrulamalarinin yapilmasi
- Soft delete/restore davranisi ile public detail/list gorunurlugunun dogrulanmasi
- Manga migration smoke kontrolu (`202603120005_manga_create_core_tables`)

## Asama 8 Chapter Kontrolleri

- `cd apps/api && go test ./internal/modules/chapter/...`
- `cd apps/api && go test ./tests/contract -run Chapter`
- `cd apps/api && go test ./tests/integration -run Chapter`
- Chapter create/update/list/detail/read/navigation akislarinin dogrulanmasi
- Preview, early access penceresi ve fallback access alanlarinin dogrulanmasi
- Media health ve integrity state guncellemelerinin dogrulanmasi
- Soft delete/restore ve publish lifecycle gecislerinin dogrulanmasi
- `chapter -> history` resume anchor ve read signal kontratinin dogrulanmasi
- Chapter migration smoke kontrolu (`202603120006_chapter_create_core_tables`)

## Asama 9 Comment Kontrolleri

- `cd apps/api && go test ./internal/modules/comment/...`
- `cd apps/api && go test ./tests/contract -run Comment`
- `cd apps/api && go test ./tests/integration -run Comment`
- Comment create/edit/delete/list/detail/thread akislarinin dogrulanmasi
- Root/reply thread yapisi, reply depth limiti ve lock etkisinin dogrulanmasi
- Soft delete gorunumu, restore window ve edit window kurallarinin dogrulanmasi
- Moderation status (visible/hidden/flagged), spoiler/pin/lock ve shadowban gorunurluk kurallarinin dogrulanmasi
- Write cooldown ve spam risk sinyali davranisinin dogrulanmasi
- Comment migration smoke kontrolu (`202603120007_comment_create_core_tables`)

## Asama 10 Support Kontrolleri

- `cd apps/api && go test ./internal/modules/support/...`
- `cd apps/api && go test ./tests/contract -run Support`
- `cd apps/api && go test ./tests/integration -run Support`
- Communication/ticket/report create akislarinin ve own list/detail davranisinin dogrulanmasi
- Duplicate/idempotency (`request_id`) ve spam risk davranislarinin dogrulanmasi
- Public reply/internal note ayriminin ve status gecislerinin dogrulanmasi
- Resolve ve moderation handoff akislarinin dogrulanmasi
- `support -> moderation` ve `support -> notification` kontrat dogrulamalarinin calistirilmasi
- Support migration smoke kontrolu (`202603120008_support_create_core_tables`)

## Asama 11 Moderation Kontrolleri

- `cd apps/api && go test ./internal/modules/moderation/...`
- `cd apps/api && go test ./tests/contract -run Moderation`
- `cd apps/api && go test ./tests/integration -run Moderation`
- Moderation queue, case detail, assignment, note, action ve escalation akislarinin dogrulanmasi
- support report -> moderation linked case idempotency davranisinin dogrulanmasi
- Moderation migration smoke kontrolu (`202603120009_moderation_create_core_tables`)

## Asama 12 Notification Kontrolleri

- `cd apps/api && go test ./internal/modules/notification/...`
- `cd apps/api && go test ./tests/contract -run Notification`
- `cd apps/api && go test ./tests/integration -run Notification`
- Support event intake, inbox/detail/read ve preference akislarinin dogrulanmasi
- Delivery pause/category-channel runtime control davranislarinin dogrulanmasi
- Notification migration smoke kontrolu (`202603120010_notification_create_core_tables`)


## Asama 13 History Kontrolleri

- `cd apps/api && go test ./internal/modules/history/...`
- `cd apps/api && go test ./tests/contract -run History`
- `cd apps/api && go test ./tests/integration -run History`
- Chapter read intake, continue-reading, own/public library ve timeline akislarinin dogrulanmasi
- Bookmark/favorite/share metadata ve runtime control davranislarinin dogrulanmasi
- History migration smoke kontrolu (`202603120011_history_create_core_tables`)

## Asama 14 Social Kontrolleri

- `cd apps/api && go test ./internal/modules/social/...`
- `cd apps/api && go test ./tests/contract -run Social`
- `cd apps/api && go test ./tests/integration -run Social`
- Friendship request/accept/reject/remove ve friend list akislarinin dogrulanmasi
- Follow/unfollow ve followers/following liste akislarinin dogrulanmasi
- Wall post/reply/list ve include-replies davranisinin dogrulanmasi
- Direct message thread open/send/list/read akislarinin dogrulanmasi
- Block/mute/restrict iliski kontrolu ve runtime toggle etkilerinin dogrulanmasi
- Social migration smoke kontrolu (`202603120012_social_create_core_tables`)


## Asama 15 Inventory Kontrolleri

- `cd apps/api && go test ./internal/modules/inventory/...`
- `cd apps/api && go test ./tests/contract -run Inventory`
- `cd apps/api && go test ./tests/integration -run Inventory`
- Item definition, own inventory list/detail, claim ve admin grant/revoke akislarinin dogrulanmasi
- Consume/equip akislari ile stackable/non-stackable davranisinin dogrulanmasi
- Runtime control (read/claim/consume/equip) toggle etkilerinin dogrulanmasi
- Inventory migration smoke kontrolu (`202603120013_inventory_create_core_tables`)
## Asama 16 Mission Kontrolleri

- `cd apps/api && go test ./internal/modules/mission/...`
- `cd apps/api && go test ./tests/contract -run Mission`
- `cd apps/api && go test ./tests/integration -run Mission`
- Mission definition ownerligi, own mission list/detail ve progress ingest akislarinin dogrulanmasi
- Completion/claim-request idempotency davranisinin ve period key yorumunun dogrulanmasi
- Runtime control (read/claim/progress-ingest/reset-hour) toggle etkilerinin dogrulanmasi
- Mission migration smoke kontrolu (`202603120014_mission_create_core_tables`)

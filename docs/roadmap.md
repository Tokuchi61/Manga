# Yol Haritası

> Bu doküman doküman üretim planı değil, sistemin kurulma ve genişleme yol haritasıdır. Modül ve ortak sistem detayları için `modules.md` ve `shared.md` referans alınmalıdır.

## 1) Kullanım Şekli
- Bu dosya yalnızca aşama sırası, kapsam sınırı, bağımlılık ve tamamlanma kriterlerini taşır.
- Modül içeriği burada uzun uzun tekrar edilmez; ilgili modül için `modules.md` referans alınır.
- Shared policy, teknik stack, settings, precedence ve operasyon kararları için `shared.md` referans alınır.
- Bir aşama bitmiş sayılmadan o aşamanın çıktı, test ve dokümantasyon karşılığı kapatılmış olmalıdır.

## 2) Faz Mantığı
Yol haritası iki ana gruptan oluşur:
1. **Omurga fazları:** sistemi inşa etmeye başlamadan önce repo, mimari, veri ve operasyon omurgasını kurar.
2. **Modül fazları:** ürün modüllerini ownerlık sınırları korunarak sırayla devreye alır.

## 3) Omurga Fazları

### Aşama 0 — Temel Standartlar Başlangıç
**Amaç:** repo, branch, versiyonlama, docker, migration, temel dokümantasyon ve geliştirme standardını sabitlemek.  
**Kapsam:** proje kök yapısı, CI/PR disiplini, Docker-first çalışma, temel scriptler, README ve changelog omurgası.  
**Tamamlanma:** proje build/run standardı, repo iskeleti ve bağlayıcı doküman omurgası hazır olmalıdır.  
**Referans:** `rules.md`, `shared.md`

### Aşama 1 — Mimari Omurga ve Sınırlar
**Amaç:** backend katmanlarını, modül sınırlarını, shared/platform ayrımını ve bağımlılık yönünü netleştirmek.  
**Kapsam:** `apps/api`, `internal/app`, `internal/platform`, `internal/shared`, `internal/modules` omurgası; modül şablonu; dosya bölme ve sorumluluk kuralları.  
**Tamamlanma:** yeni bir modül açıldığında nereye yerleşeceği ve hangi katmanda hangi işin yapılacağı tartışmasız olmalıdır.  
**Referans:** `rules.md`, `modules.md`

### Aşama 2 — Çekirdek Ürün Hazırlığı
**Amaç:** ortak sözlükler, transaction kuralları, audit/idempotency/outbox gibi sistem çapı kuralları hazır hale getirmek.  
**Kapsam:** shared sözlükler, teknik stack, cache/queue, projection, media, reporting, search ve settings envanteri.  
**Tamamlanma:** ilk iş modülleri geliştirilmeden önce ortak policy ve teknik kararlar yazılı hale gelmiş olmalıdır.  
**Referans:** `shared.md`

### Aşama 3 — Genişleme ve Ölçeklenme Hazırlığı
**Amaç:** sistem büyüdükçe modülleşme, operasyon ve bakım disiplini bozulmadan ilerlemek.  
**Kapsam:** domain-group kullanımı, projection/read model stratejisi, reporting, reconcile, bakım ve refactor kuralları.  
**Tamamlanma:** yeni modüller ve çapraz akışlar eklendiğinde mevcut yapı bozulmadan genişleyebilir olmalıdır.  
**Referans:** `rules.md`, `shared.md`, `modules.md`

## 4) Modül Fazları

Aşağıdaki sıralama sistem kurulum sırasını gösterir. Her modülün detaylı kapsamı `modules.md` içindedir.

### Aşama 4 — Auth
Kimlik doğrulama, credential, session, token, verification, recovery ve auth güvenlik akışları.  
**Bağımlılık:** Aşama 0-3  
**Referans modül:** `auth`

### Aşama 5 — User
Kullanıcı hesabı, profil, görünürlük, üyelik ve VIP state omurgası.  
**Bağımlılık:** Auth  
**Referans modül:** `user`

### Aşama 6 — Access
Merkezi authorization, policy evaluation, feature availability ve final allow/deny katmanı.  
**Bağımlılık:** Auth, User  
**Referans modül:** `access`

### Aşama 7 — Manga
Ana içerik varlığı, metadata, discovery ve listing omurgası.  
**Bağımlılık:** Omurga fazları, Access  
**Referans modül:** `manga`

### Aşama 8 — Chapter
Chapter, page, release, early access ve okuma yüzeyi.  
**Bağımlılık:** Manga, Access  
**Referans modül:** `chapter`

### Aşama 9 — Comment
İçerik yorumları, thread, etkileşim ve anti-spam yazma yüzeyi.  
**Bağımlılık:** Auth, User, Manga/Chapter, Access  
**Referans modül:** `comment`

### Aşama 10 — Support
Destek talebi, report intake ve vaka açılış yüzeyleri.  
**Bağımlılık:** Auth, User, Access  
**Referans modül:** `support`

### Aşama 11 — Moderation
Moderation queue, case yönetimi ve karar uygulama sinyalleri.  
**Bağımlılık:** Support, Comment, User, Admin, Access  
**Referans modül:** `moderation`

### Aşama 12 — Notification
Bildirim üretimi, teslimi, kategori ve tercih yönetimi.  
**Bağımlılık:** Auth, User, Support, Moderation, Social, Mission  
**Referans modül:** `notification`

### Aşama 13 — History
Continue reading, kütüphane, timeline ve bookmark yüzeyleri.  
**Bağımlılık:** Auth, User, Manga, Chapter, Access  
**Referans modül:** `history`

### Aşama 14 — Social
Arkadaşlık, follow, duvar, mesajlaşma, block/mute/restrict akışları.  
**Bağımlılık:** Auth, User, Access, Notification  
**Referans modül:** `social`

### Aşama 15 — Inventory
Item sahipliği, claim, consume, equip ve kullanıcıya bağlı envanter state’i.  
**Bağımlılık:** User, Access, Notification  
**Referans modül:** `inventory`

### Aşama 16 — Mission
Görev tanımı, progress ingest, eligibility ve claim akışları.  
**Bağımlılık:** Auth, User, Inventory, Notification, History, Social  
**Referans modül:** `mission`

### Aşama 17 — RoyalPass
Sezon, tier, premium track ve reward progression.  
**Bağımlılık:** Mission, Inventory, Payment, Access  
**Referans modül:** `royalpass`

### Aşama 18 — Shop
Katalog, offer, ürün yüzeyi ve purchase orchestration.  
**Bağımlılık:** Inventory, Payment, Access  
**Referans modül:** `shop`

### Aşama 19 — Payment
Checkout, provider session, callback, ledger ve reconcile doğruluğu.  
**Bağımlılık:** Auth, User, Shop, RoyalPass, Access  
**Referans modül:** `payment`

### Aşama 20 — Ads
Placement, campaign, delivery, impression/click intake ve ölçümleme.  
**Bağımlılık:** Manga, Chapter, Access, Reporting/Analytics  
**Referans modül:** `ads`

### Aşama 21 — Admin
Yönetim yüzeyleri, operasyon kontrolleri, modül bazlı iç işlemler ve denetim araçları.  
**Bağımlılık:** sistemdeki ilgili modüller  
**Referans modül:** `admin`

## 5) Aşama Geçiş Kuralları
- Bir modül aşamasına geçmeden önce bağımlı olduğu önceki modüllerin ownerlığı ve kontratları yazılı olmalıdır.
- Modül detayı `modules.md` içinde tanımlanmadan modül implementasyonuna başlanmamalıdır.
- Yeni shared enum, policy veya ayar anahtarı gerekiyorsa aynı değişiklikte `shared.md` güncellenmelidir.
- Çoklu işlem yapan dosyalar tek dosyada toplanmamalıdır; handler, service, usecase, job ve benzeri alanlarda işlem bazlı parçalama korunmalıdır.
- Her aşama için minimum unit/integration/contract testi tanımlanmadan aşama tamamlanmış kabul edilmemelidir.

## 6) Bu Dosyada Neler Olmamalı?
- Uzun modül açıklamaları
- Shared enum tablolarının tamamı
- Settings envanterinin tüm satırları
- Satır satır klasör içeriği

Bu ayrıntılar ilgili referans dosyalarında tutulmalıdır.

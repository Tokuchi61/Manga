# Yol HaritasÄ±

> Bu dokÃ¼man dokÃ¼man Ã¼retim planÄ± deÄŸil, sistemin kurulma ve geniÅŸleme yol haritasÄ±dÄ±r. ModÃ¼l ve ortak sistem detaylarÄ± iÃ§in `modules.md` ve `shared.md` referans alÄ±nmalÄ±dÄ±r.

## 1) KullanÄ±m Åekli
- Bu dosya yalnÄ±zca aÅŸama sÄ±rasÄ±, kapsam sÄ±nÄ±rÄ±, baÄŸÄ±mlÄ±lÄ±k ve tamamlanma kriterlerini taÅŸÄ±r.
- ModÃ¼l iÃ§eriÄŸi burada uzun uzun tekrar edilmez; ilgili modÃ¼l iÃ§in `modules.md` referans alÄ±nÄ±r.
- Shared policy, teknik stack, settings, precedence ve operasyon kararlarÄ± iÃ§in `shared.md` referans alÄ±nÄ±r.
- Bir aÅŸama bitmiÅŸ sayÄ±lmadan o aÅŸamanÄ±n Ã§Ä±ktÄ±, test ve dokÃ¼mantasyon karÅŸÄ±lÄ±ÄŸÄ± kapatÄ±lmÄ±ÅŸ olmalÄ±dÄ±r.

## 2) Faz MantÄ±ÄŸÄ±
Yol haritasÄ± iki ana gruptan oluÅŸur:
1. **Omurga fazlarÄ±:** sistemi inÅŸa etmeye baÅŸlamadan Ã¶nce repo, mimari, veri ve operasyon omurgasÄ±nÄ± kurar.
2. **ModÃ¼l fazlarÄ±:** Ã¼rÃ¼n modÃ¼llerini ownerlÄ±k sÄ±nÄ±rlarÄ± korunarak sÄ±rayla devreye alÄ±r.

## 3) Omurga FazlarÄ±

### AÅŸama 0 â€” Temel Standartlar BaÅŸlangÄ±Ã§
**AmaÃ§:** repo, branch, versiyonlama, docker, migration, temel dokÃ¼mantasyon ve geliÅŸtirme standardÄ±nÄ± sabitlemek.  
**Kapsam:** proje kÃ¶k yapÄ±sÄ±, CI/PR disiplini, Docker-first Ã§alÄ±ÅŸma, temel scriptler, README ve changelog omurgasÄ±.  
**Tamamlanma:** proje build/run standardÄ±, repo iskeleti ve baÄŸlayÄ±cÄ± dokÃ¼man omurgasÄ± hazÄ±r olmalÄ±dÄ±r.  
**Referans:** `rules.md`, `shared.md`

### AÅŸama 1 â€” Mimari Omurga ve SÄ±nÄ±rlar
**AmaÃ§:** backend katmanlarÄ±nÄ±, modÃ¼l sÄ±nÄ±rlarÄ±nÄ±, shared/platform ayrÄ±mÄ±nÄ± ve baÄŸÄ±mlÄ±lÄ±k yÃ¶nÃ¼nÃ¼ netleÅŸtirmek.  
**Kapsam:** `apps/api`, `internal/app`, `internal/platform`, `internal/shared`, `internal/modules` omurgasÄ±; modÃ¼l ÅŸablonu; dosya bÃ¶lme ve sorumluluk kurallarÄ±.  
**Tamamlanma:** yeni bir modÃ¼l aÃ§Ä±ldÄ±ÄŸÄ±nda nereye yerleÅŸeceÄŸi ve hangi katmanda hangi iÅŸin yapÄ±lacaÄŸÄ± tartÄ±ÅŸmasÄ±z olmalÄ±dÄ±r.  
**Referans:** `rules.md`, `modules.md`

### AÅŸama 2 â€” Ã‡ekirdek ÃœrÃ¼n HazÄ±rlÄ±ÄŸÄ±
**AmaÃ§:** ortak sÃ¶zlÃ¼kler, transaction kurallarÄ±, audit/idempotency/outbox gibi sistem Ã§apÄ± kurallarÄ± hazÄ±r hale getirmek.  
**Kapsam:** shared sÃ¶zlÃ¼kler, teknik stack, cache/queue, projection, media, reporting, search ve settings envanteri.  
**Tamamlanma:** ilk iÅŸ modÃ¼lleri geliÅŸtirilmeden Ã¶nce ortak policy ve teknik kararlar yazÄ±lÄ± hale gelmiÅŸ olmalÄ±dÄ±r.  
**Referans:** `shared.md`

### AÅŸama 3 â€” GeniÅŸleme ve Ã–lÃ§eklenme HazÄ±rlÄ±ÄŸÄ±
**AmaÃ§:** sistem bÃ¼yÃ¼dÃ¼kÃ§e modÃ¼lleÅŸme, operasyon ve bakÄ±m disiplini bozulmadan ilerlemek.  
**Kapsam:** domain-group kullanÄ±mÄ±, projection/read model stratejisi, reporting, reconcile, bakÄ±m ve refactor kurallarÄ±.  
**Tamamlanma:** yeni modÃ¼ller ve Ã§apraz akÄ±ÅŸlar eklendiÄŸinde mevcut yapÄ± bozulmadan geniÅŸleyebilir olmalÄ±dÄ±r.  
**Referans:** `rules.md`, `shared.md`, `modules.md`

## 4) ModÃ¼l FazlarÄ±

AÅŸaÄŸÄ±daki sÄ±ralama sistem kurulum sÄ±rasÄ±nÄ± gÃ¶sterir. Her modÃ¼lÃ¼n detaylÄ± kapsamÄ± `modules.md` iÃ§indedir.

### AÅŸama 4 â€” Auth
Kimlik doÄŸrulama, credential, session, token, verification, recovery ve auth gÃ¼venlik akÄ±ÅŸlarÄ±.  
**BaÄŸÄ±mlÄ±lÄ±k:** AÅŸama 0-3  
**Referans modÃ¼l:** `auth`

### AÅŸama 5 â€” User
KullanÄ±cÄ± hesabÄ±, profil, gÃ¶rÃ¼nÃ¼rlÃ¼k, Ã¼yelik ve VIP state omurgasÄ±.  
**BaÄŸÄ±mlÄ±lÄ±k:** Auth  
**Referans modÃ¼l:** `user`

### AÅŸama 6 â€” Access
Merkezi authorization, policy evaluation, feature availability ve final allow/deny katmanÄ±.  
**BaÄŸÄ±mlÄ±lÄ±k:** Auth, User  
**Referans modÃ¼l:** `access`

### AÅŸama 7 â€” Manga
Ana iÃ§erik varlÄ±ÄŸÄ±, metadata, discovery ve listing omurgasÄ±.  
**BaÄŸÄ±mlÄ±lÄ±k:** Omurga fazlarÄ±, Access  
**Referans modÃ¼l:** `manga`

### AÅŸama 8 â€” Chapter
Chapter, page, release, early access ve okuma yÃ¼zeyi.  
**BaÄŸÄ±mlÄ±lÄ±k:** Manga, Access  
**Referans modÃ¼l:** `chapter`

### AÅŸama 9 â€” Comment
Ä°Ã§erik yorumlarÄ±, thread, etkileÅŸim ve anti-spam yazma yÃ¼zeyi.  
**BaÄŸÄ±mlÄ±lÄ±k:** Auth, User, Manga/Chapter, Access  
**Referans modÃ¼l:** `comment`

### AÅŸama 10 â€” Support
Destek talebi, report intake ve vaka aÃ§Ä±lÄ±ÅŸ yÃ¼zeyleri.  
**BaÄŸÄ±mlÄ±lÄ±k:** Auth, User, Access  
**Referans modÃ¼l:** `support`

### AÅŸama 11 â€” Moderation
Moderation queue, case yÃ¶netimi ve karar uygulama sinyalleri.  
**BaÄŸÄ±mlÄ±lÄ±k:** Support, Comment, User, Admin, Access  
**Referans modÃ¼l:** `moderation`

### AÅŸama 12 â€” Notification
Bildirim Ã¼retimi, teslimi, kategori ve tercih yÃ¶netimi.  
**BaÄŸÄ±mlÄ±lÄ±k:** Auth, User, Support, Moderation, Social, Mission  
**Referans modÃ¼l:** `notification`

### AÅŸama 13 â€” History
Continue reading, kÃ¼tÃ¼phane, timeline ve bookmark yÃ¼zeyleri.  
**BaÄŸÄ±mlÄ±lÄ±k:** Auth, User, Manga, Chapter, Access  
**Referans modÃ¼l:** `history`

### AÅŸama 14 â€” Social
ArkadaÅŸlÄ±k, follow, duvar, mesajlaÅŸma, block/mute/restrict akÄ±ÅŸlarÄ±.  
**BaÄŸÄ±mlÄ±lÄ±k:** Auth, User, Access, Notification  
**Referans modÃ¼l:** `social`

### AÅŸama 15 â€” Inventory
Item sahipliÄŸi, claim, consume, equip ve kullanÄ±cÄ±ya baÄŸlÄ± envanter stateâ€™i.  
**BaÄŸÄ±mlÄ±lÄ±k:** User, Access, Notification  
**Referans modÃ¼l:** `inventory`

### AÅŸama 16 â€” Mission
GÃ¶rev tanÄ±mÄ±, progress ingest, eligibility ve claim akÄ±ÅŸlarÄ±.  
**BaÄŸÄ±mlÄ±lÄ±k:** Auth, User, Inventory, Notification, History, Social  
**Referans modÃ¼l:** `mission`

### AÅŸama 17 â€” RoyalPass
Sezon, tier, premium track ve reward progression.  
**BaÄŸÄ±mlÄ±lÄ±k:** Mission, Inventory, Payment, Access  
**Referans modÃ¼l:** `royalpass`

### AÅŸama 18 â€” Shop
Katalog, offer, Ã¼rÃ¼n yÃ¼zeyi ve purchase orchestration.  
**BaÄŸÄ±mlÄ±lÄ±k:** Inventory, Payment, Access  
**Referans modÃ¼l:** `shop`

### AÅŸama 19 â€” Payment
Checkout, provider session, callback, ledger ve reconcile doÄŸruluÄŸu.  
**BaÄŸÄ±mlÄ±lÄ±k:** Auth, User, Shop, RoyalPass, Access  
**Referans modÃ¼l:** `payment`

### AÅŸama 20 â€” Ads
Placement, campaign, delivery, impression/click intake ve Ã¶lÃ§Ã¼mleme.  
**BaÄŸÄ±mlÄ±lÄ±k:** Manga, Chapter, Access, Reporting/Analytics  
**Referans modÃ¼l:** `ads`

### AÅŸama 21 â€” Admin
YÃ¶netim yÃ¼zeyleri, operasyon kontrolleri, modÃ¼l bazlÄ± iÃ§ iÅŸlemler ve denetim araÃ§larÄ±.  
**BaÄŸÄ±mlÄ±lÄ±k:** sistemdeki ilgili modÃ¼ller  
**Referans modÃ¼l:** `admin`

## 5) AÅŸama GeÃ§iÅŸ KurallarÄ±
- Bir modÃ¼l aÅŸamasÄ±na geÃ§meden Ã¶nce baÄŸÄ±mlÄ± olduÄŸu Ã¶nceki modÃ¼llerin ownerlÄ±ÄŸÄ± ve kontratlarÄ± yazÄ±lÄ± olmalÄ±dÄ±r.
- ModÃ¼l detayÄ± `modules.md` iÃ§inde tanÄ±mlanmadan modÃ¼l implementasyonuna baÅŸlanmamalÄ±dÄ±r.
- Yeni shared enum, policy veya ayar anahtarÄ± gerekiyorsa aynÄ± deÄŸiÅŸiklikte `shared.md` gÃ¼ncellenmelidir.
- Ã‡oklu iÅŸlem yapan dosyalar tek dosyada toplanmamalÄ±dÄ±r; handler, service, usecase, job ve benzeri alanlarda iÅŸlem bazlÄ± parÃ§alama korunmalÄ±dÄ±r.
- Her aÅŸama iÃ§in minimum unit/integration/contract testi tanÄ±mlanmadan aÅŸama tamamlanmÄ±ÅŸ kabul edilmemelidir.

## 6) Bu Dosyada Neler OlmamalÄ±?
- Uzun modÃ¼l aÃ§Ä±klamalarÄ±
- Shared enum tablolarÄ±nÄ±n tamamÄ±
- Settings envanterinin tÃ¼m satÄ±rlarÄ±
- SatÄ±r satÄ±r klasÃ¶r iÃ§eriÄŸi

Bu ayrÄ±ntÄ±lar ilgili referans dosyalarÄ±nda tutulmalÄ±dÄ±r.

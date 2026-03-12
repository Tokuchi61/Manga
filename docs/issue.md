# Proje Uyum Denetimi - Kod Bulgulari

Tarih: 2026-03-12
Kapsam: Asama 0-5 kod tabani (dokuman-path uyumsuzluklari bu kayittan bilincli olarak dislandi)

## [P1] User account state degisimi istemci girdisine gore `admin` gibi davranabiliyor
- Dosya: `apps/api/internal/modules/user/handler/user_account_handler.go:10`
- Dosya: `apps/api/internal/modules/user/service/user_account_service.go:27`
- Bulgu: `actor_scope` ve `actor_user_id` dogrudan request body'den alinip yetki kaynagi gibi kullaniliyor.
- Etki: Kimligi dogrulanmamis bir istek `actor_scope=admin` vererek baska kullanicilari `banned/deactivated` durumuna cekebilir.
- Oneri: `actor_scope` ve actor kimligi body'den alinmamali; auth/access context'inden derive edilmeli.

## [P1] Auth logout/session revoke akislari istemci sagladigi kimlik/session alanlarina guveniyor
- Dosya: `apps/api/internal/modules/auth/handler/auth_token_handler.go:23`
- Dosya: `apps/api/internal/modules/auth/service/auth_session_service.go:14`
- Dosya: `apps/api/internal/modules/auth/service/auth_session_service.go:55`
- Bulgu: `credential_id` + `session_id` body parametreleri ile oturum iptali yapiliyor; caller kimligi ile server-side bag kurulmamis.
- Etki: ID degerlerini bilen saldirgan, kendi olmayan session'lari revoke etmeyi deneyebilir.
- Oneri: Caller kimligi token/session context'ten alinmali, body'den gelen kimlik alanlari authoritative olmamali.

## [P1] User olusturma akisi Auth credential varligini dogrulamiyor
- Dosya: `apps/api/internal/modules/user/service/user_create_service.go:19`
- Bulgu: `credential_id` sadece UUID formatinda parse ediliyor; `auth` owner kaynagi ile varlik/uygunluk kontrolu yok.
- Etki: Sistemde bulunmayan credential ID ile user kaydi acilabilir (sahte baglanti riski).
- Oneri: `auth -> user` kontrati uzerinden credential varlik ve kullanilabilirlik dogrulamasi eklenmeli.

## [P2] `failed_attempt_limit_per_minute` semantigi kodda uygulanmiyor
- Dosya: `apps/api/internal/platform/config/config.go:21`
- Dosya: `apps/api/internal/modules/auth/service/auth_login_service.go:52`
- Referans not: `docs/shared.md:875` anahtar aciklamasi "dakika basina" limit olarak tanimli.
- Bulgu: Login tarafinda sadece toplam ardiskik deneme sayisi takip ediliyor; dakika penceresi (rolling/fixed window) yok.
- Etki: Ayar ismiyle davranis farkli; limit beklenenden daha agresif veya farkli calisabilir.
- Oneri: Denemeler icin zaman pencereli sayaç modeli eklenmeli ya da anahtar semantigi yeniden adlandirilmali.

## [P2] Username validasyonu normalize sonrasi bypass edilebiliyor
- Dosya: `apps/api/internal/modules/user/dto/user_create_dto.go:6`
- Dosya: `apps/api/internal/modules/user/repository/user_account_repository.go:15`
- Bulgu: `min=3` validasyonu trim/lower normalize oncesi calisiyor; repository'de trim edilince efektif username 3 altina inebiliyor.
- Etki: Is kurali min uzunluk beklentisi veri katmaninda korunmuyor.
- Oneri: Validasyon oncesi canonical normalize edilmeli veya normalize sonrasi ikinci dogrulama uygulanmali.

## [P2] Banned/deactivated hesaplar icin VIP state mutasyonu engellenmiyor
- Dosya: `apps/api/internal/modules/user/service/user_vip_service.go:13`
- Bulgu: VIP lifecycle akisinda hesap durumu kontrolu yok (`banned` / `deactivated` hesapta da activate/freeze/resume/deactivate calisabiliyor).
- Etki: Hesap durumu ile uyumsuz uyelik state gecisleri olusabilir.
- Oneri: VIP akisinin basinda account state guard eklenmeli.

## [P3] User event sabitleri tanimli ama yayin akisi yok
- Dosya: `apps/api/internal/modules/user/events/events.go:3`
- Bulgu: Stage 5 event sabitleri tanimli; ancak service tarafinda event publish/append mekanizmasi bulunmuyor.
- Etki: Downstream moduller icin beklenen sinyal akisinin isletimi garanti degil.
- Oneri: En azindan event append surface (veya outbox-ready abstraction) eklenmeli.



# Issues

Bu dosya proje genel denetim bulgularinin acik kalanlarini izler.

Son denetim tarihi: 2026-03-12
Durum: Kismen kapatildi

## Ozet

- Toplam acik bulgu: 4
- Kritiklik dagilimi: `P1=1`, `P2=3`
- Kapatilan basliklar: `ISS-001`, `ISS-002`, `ISS-003`, `ISS-005`, `ISS-006`, `ISS-007`, `ISS-008`, `ISS-009`, `ISS-010`, `ISS-012`, `ISS-014`, `ISS-015`

## Acik Bulgular

### [P1] ISS-004 - Mimari hedef ile calisan persistence mimarisi uyusmuyor

- Durum: Acik
- Not: Moduller halen varsayilan olarak memory store ile aciliyor; DB-backed repository gecisi tamamlanmadi.

### [P1] ISS-011 - Yari kalmis/tekrarli yapilar bakim maliyetini artiriyor

- Durum: Acik
- Not: Access tarafindaki kullanilmayan cache/config alanlari ve bazi tekrarli hata namespace alanlari sadelestirilmedi.

### [P2] ISS-013 - Test stratejisi belgede var, uygulamada eksik

- Durum: Acik
- Not: Entegrasyon/contract testleri guclendirildi ancak `apps/api/tests/e2e` katmani halen yok.

### [P2] ISS-016 - Dokuman encoding bozulmalari var (mojibake)

- Durum: Acik
- Not: Dokumanlarin UTF-8/BOM-suz normalize edilmesi ve CI lint adimi henuz tamamlanmadi.

## Son Uygulama Notu (2026-03-12)

- Kimlik/actor guven siniri request context'e tasindi.
- Access evaluate guven siniri ve role guardlari sikilastirildi.
- Moduller arasi varlik/sahiplik kontrat dogrulamalari eklendi.
- Hata/response sinirlari sertlestirildi (`internal_error` default).
- Config ve rate-limit semantigi duzeltildi.
- Chapter navigation akisi hard-limit taramasi yerine hedefli cozumlemeye cekildi.
- Dokuman yol referanslari mevcut tek dosya yapisina hizalandi ve `docs/issue.md` kaldirildi.

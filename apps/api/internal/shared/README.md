# shared Katmani

`apps/api/internal/shared/` sadece domain-agnostic ortak yapilar icindir.

- modullerden kod kacirma alani olarak kullanilmaz
- bir kod parcasi yalnizca iki modul kullaniyor diye otomatik shared'e alinmaz

## Asama 2 Canonical Paketleri

- `catalog/`: docs/shared.md icindeki canonical enum ve sozluk kayitlari
- `policy/`: transaction, outbox, projection, stack ve operasyonel ortak policy kayitlari
- `settings/`: runtime settings sozlugu, key grameri ve yorumlama sirasina ait ortak modeller

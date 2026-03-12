# shared Katmani

`apps/api/internal/shared/` sadece domain-agnostic ortak yapilar icindir.

- modullerden kod kacirma alani olarak kullanilmaz
- bir kod parcasi yalnizca iki modul kullaniyor diye otomatik shared'e alinmaz

## Asama 2 Canonical Paketleri

- `catalog/`: docs/shared.md icindeki canonical enum ve sozluk kayitlari
- `policy/`: transaction, outbox, projection, stack ve operasyonel ortak policy kayitlari
- `settings/`: runtime settings sozlugu, key grameri ve yorumlama sirasina ait ortak modeller

## Asama 3 Olceklenme Hazirligi

- `policy/scaling.go`: domain-group kullanim rehberi, module envanter kaydi ve status sozlugu
- `policy/readmodel.go`: projection/read model uygulama kurallari ve consistency penceresi
- `policy/reporting.go`: operasyon summary, analytics aggregate ve export query katmanlari
- `policy/reconcile.go`: reconcile gereken kritik akis referanslari ve guardrail kurallari
- `policy/maintenance.go`: bakim ve refactor disiplini icin canonical guardrail seti

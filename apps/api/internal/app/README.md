# app Katmani

`apps/api/internal/app/` yalnizca bootstrap ve composition root sorumlulugunu tasir.

- module wiring ve route mount islemleri burada yonetilir
- is kurali veya moduller arasi owner verisi burada tutulmaz
- moduller `app` paketini import etmez

# modules Katmani

`apps/api/internal/modules/` tum leaf modullerin canonical kokudur.

- varsayilan yapi: `apps/api/internal/modules/<module>/`
- opsiyonel domain group: `apps/api/internal/modules/<domain-group>/<module>/`
- module route mount islemleri app bootstrap tarafindan merkezi olarak yapilir

## Modul Iskeleti

Asagidaki script Stage 1 modul iskeletini olusturur:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/scaffold_module.ps1 -ModuleName auth
powershell -ExecutionPolicy Bypass -File scripts/scaffold_module.ps1 -ModuleName manga -DomainGroup content
```

Script su dosya/klasor omurgasini olusturur:

- `entity/`, `dto/`, `service/`, `repository/`, `handler/`, `middleware/`, `validator/`, `mapper/`, `contract/`, `events/`, `consumer/`, `producer/`, `jobs/`, `readmodel/`
- `module.go`, `routes.go`, `errors.go`

package admin

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/handler"
	adminrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires admin handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.AdminService
	snapshotter snapshot.Snapshotter
}

func New() Module {
	store := adminrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc, snapshotter: store}
}

func (m Module) Name() string {
	return "admin"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

func (m Module) MaintenanceEnabled(ctx context.Context) (bool, error) {
	if m.svc == nil {
		return false, nil
	}
	cfg, err := m.svc.GetRuntimeConfig(ctx)
	if err != nil {
		return false, err
	}
	return cfg.MaintenanceEnabled, nil
}

func (m *Module) Snapshot() ([]byte, error) {
	if m == nil || m.snapshotter == nil {
		return nil, nil
	}
	return m.snapshotter.Snapshot()
}

func (m *Module) RestoreSnapshot(data []byte) error {
	if m == nil || m.snapshotter == nil {
		return nil
	}
	return m.snapshotter.RestoreSnapshot(data)
}

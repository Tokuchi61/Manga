package manga

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/handler"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires manga handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.MangaService
	snapshotter snapshot.Snapshotter
}

func New() Module {
	store := mangarepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc, snapshotter: store}
}

func (m Module) Name() string {
	return "manga"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

func (m Module) TargetExists(ctx context.Context, mangaID string) (bool, error) {
	if m.svc == nil {
		return false, nil
	}
	return m.svc.TargetExists(ctx, mangaID)
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

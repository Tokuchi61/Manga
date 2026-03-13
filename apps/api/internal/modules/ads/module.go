package ads

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/handler"
	adsrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires ads handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	snapshotter snapshot.Snapshotter
}

func New() Module {
	store := adsrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), snapshotter: store}
}

func (m Module) Name() string {
	return "ads"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
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

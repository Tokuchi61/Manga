package history

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/handler"
	historyrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/history/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires history handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.HistoryService
	snapshotter snapshot.Snapshotter
}

func New() Module {
	store := historyrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc, snapshotter: store}
}

func (m *Module) SetChapterSignalProvider(provider service.ChapterSignalProvider) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetChapterSignalProvider(provider)
}

func (m Module) Name() string {
	return "history"
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

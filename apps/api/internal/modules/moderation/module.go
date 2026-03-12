package moderation

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/handler"
	moderationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires moderation handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.ModerationService
	snapshotter snapshot.Snapshotter
}

func New() Module {
	store := moderationrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc, snapshotter: store}
}

func (m *Module) SetSupportContracts(lookup service.SupportHandoffLookup, linker service.SupportCaseLinker) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetSupportContracts(lookup, linker)
}

func (m Module) Name() string {
	return "moderation"
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

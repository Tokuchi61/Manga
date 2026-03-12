package user

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/handler"
	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires user handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.UserService
	snapshotter snapshot.Snapshotter
}

func New() Module {
	store := userrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc, snapshotter: store}
}

func (m *Module) SetCredentialLookup(lookup service.CredentialLookup) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetCredentialLookup(lookup)
}

func (m Module) Name() string {
	return "user"
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

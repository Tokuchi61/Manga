package notification

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/handler"
	notificationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires notification handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.NotificationService
	snapshotter snapshot.Snapshotter
}

func New() Module {
	store := notificationrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc, snapshotter: store}
}

func (m *Module) SetSupportSignalProvider(provider service.SupportSignalProvider) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetSupportSignalProvider(provider)
}

func (m Module) Name() string {
	return "notification"
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

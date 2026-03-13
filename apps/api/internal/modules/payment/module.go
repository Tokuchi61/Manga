package payment

import (
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/handler"
	paymentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// RuntimeConfig maps stage-19 payment runtime keys.
type RuntimeConfig struct {
	CallbackSigningSecret  string
	CallbackTimestampSkew  time.Duration
	CallbackEventPublisher service.CallbackEventPublisher
}

// Module wires payment handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	snapshotter snapshot.Snapshotter
}

func New(cfg RuntimeConfig) Module {
	store := paymentrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator, service.Config{
		CallbackSigningSecret:  cfg.CallbackSigningSecret,
		CallbackTimestampSkew:  cfg.CallbackTimestampSkew,
		CallbackEventPublisher: cfg.CallbackEventPublisher,
	})

	return Module{httpHandler: handler.New(svc), snapshotter: store}
}

func (m Module) Name() string {
	return "payment"
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

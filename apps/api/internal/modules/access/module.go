package access

import (
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/handler"
	accessrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// RuntimeConfig maps stage-6 access runtime knobs.
type RuntimeConfig struct {
	DecisionCacheTTLSeconds int
}

// Module wires access handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
}

func New(cfg RuntimeConfig) Module {
	store := accessrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator, service.Config{
		DecisionCacheTTL: time.Duration(cfg.DecisionCacheTTLSeconds) * time.Second,
	})

	return Module{httpHandler: handler.New(svc)}
}

func (m Module) Name() string {
	return "access"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

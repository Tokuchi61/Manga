package manga

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/handler"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires manga handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
}

func New() Module {
	store := mangarepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc)}
}

func (m Module) Name() string {
	return "manga"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

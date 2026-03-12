package chapter

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/handler"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires chapter handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
}

func New() Module {
	store := chapterrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc)}
}

func (m Module) Name() string {
	return "chapter"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

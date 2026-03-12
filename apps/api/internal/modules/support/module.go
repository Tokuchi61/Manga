package support

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/handler"
	supportrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires support handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.SupportService
}

func New() Module {
	store := supportrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc}
}

func (m *Module) SetTargetLookups(mangaLookup service.MangaTargetLookup, chapterLookup service.ChapterTargetLookup, commentLookup service.CommentTargetLookup) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetTargetLookups(mangaLookup, chapterLookup, commentLookup)
}

func (m Module) Name() string {
	return "support"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

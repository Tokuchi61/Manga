package comment

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/handler"
	commentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires comment handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.CommentService
}

func New() Module {
	store := commentrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc}
}

func (m *Module) SetTargetLookups(mangaLookup service.MangaTargetLookup, chapterLookup service.ChapterTargetLookup) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetTargetLookups(mangaLookup, chapterLookup)
}

func (m Module) Name() string {
	return "comment"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

func (m Module) TargetExists(ctx context.Context, commentID string) (bool, error) {
	if m.svc == nil {
		return false, nil
	}
	return m.svc.TargetExists(ctx, commentID)
}

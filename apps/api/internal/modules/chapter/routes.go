package chapter

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/handler"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("chapter registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("chapter registerRoutes: http handler cannot be nil")
	}

	router.Get("/manga/{manga_id}/chapters", httpHandler.ListByManga)

	router.Route("/chapters", func(r chi.Router) {
		r.Post("/", httpHandler.CreateChapter)
		r.Get("/{chapter_id}", httpHandler.GetDetail)
		r.Get("/{chapter_id}/read", httpHandler.Read)
		r.Get("/{chapter_id}/navigation", httpHandler.Navigation)

		r.Patch("/{chapter_id}", httpHandler.UpdateChapter)
		r.Patch("/{chapter_id}/reorder", httpHandler.ReorderChapter)
		r.Patch("/{chapter_id}/publish", httpHandler.UpdatePublishState)
		r.Patch("/{chapter_id}/access", httpHandler.UpdateAccess)
		r.Patch("/{chapter_id}/media-health", httpHandler.UpdateMediaHealth)
		r.Patch("/{chapter_id}/integrity", httpHandler.UpdateIntegrity)
		r.Delete("/{chapter_id}", httpHandler.SoftDeleteChapter)
		r.Post("/{chapter_id}/restore", httpHandler.RestoreChapter)
	})
}

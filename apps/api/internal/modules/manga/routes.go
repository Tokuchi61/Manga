package manga

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/handler"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("manga registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("manga registerRoutes: http handler cannot be nil")
	}

	router.Route("/manga", func(r chi.Router) {
		r.Get("/", httpHandler.ListManga)
		r.Get("/discovery", httpHandler.ListDiscovery)
		r.Get("/{manga_id}", httpHandler.GetMangaDetail)

		r.Post("/", httpHandler.CreateManga)
		r.Patch("/{manga_id}", httpHandler.UpdateMangaMetadata)
		r.Patch("/{manga_id}/publish", httpHandler.UpdatePublishState)
		r.Patch("/{manga_id}/visibility", httpHandler.UpdateVisibility)
		r.Patch("/{manga_id}/editorial", httpHandler.UpdateEditorial)
		r.Patch("/{manga_id}/counters", httpHandler.SyncCounters)
		r.Delete("/{manga_id}", httpHandler.SoftDeleteManga)
		r.Post("/{manga_id}/restore", httpHandler.RestoreManga)
	})
}

package history

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("history registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("history registerRoutes: http handler cannot be nil")
	}

	router.Route("/history", func(r chi.Router) {
		r.With(identity.RequireUser).Get("/continue-reading", httpHandler.ListContinueReading)
		r.With(identity.RequireUser).Get("/library", httpHandler.ListLibrary)
		r.Get("/public/{user_id}/library", httpHandler.ListPublicLibrary)
		r.With(identity.RequireUser).Get("/timeline", httpHandler.ListTimeline)

		r.With(identity.RequireUser).Post("/intake/chapter", httpHandler.IngestChapterSignal)
		r.With(identity.RequireUser).Post("/library/{manga_id}/bookmark", httpHandler.UpdateBookmark)
		r.With(identity.RequireUser).Post("/library/{manga_id}/share", httpHandler.UpdateShare)

		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/continue-reading-state", httpHandler.UpdateContinueReadingState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/library-state", httpHandler.UpdateLibraryState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/timeline-state", httpHandler.UpdateTimelineState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/bookmark-write-state", httpHandler.UpdateBookmarkWriteState)
	})
}

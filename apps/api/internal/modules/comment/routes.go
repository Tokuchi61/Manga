package comment

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("comment registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("comment registerRoutes: http handler cannot be nil")
	}

	router.Get("/targets/{target_type}/{target_id}/comments", httpHandler.ListByTarget)

	router.Route("/comments", func(r chi.Router) {
		r.With(identity.RequireUser).Post("/", httpHandler.CreateComment)
		r.Get("/{comment_id}", httpHandler.GetDetail)
		r.Get("/{comment_id}/thread", httpHandler.GetThread)
		r.With(identity.RequireUser).Patch("/{comment_id}", httpHandler.UpdateComment)
		r.With(identity.RequireUser).Delete("/{comment_id}", httpHandler.DeleteComment)
		r.With(identity.RequireUser).Post("/{comment_id}/restore", httpHandler.RestoreComment)
		r.With(identity.RequireUser, identity.RequireAnyRole("moderator", "admin")).Patch("/{comment_id}/moderation", httpHandler.UpdateModeration)
	})
}

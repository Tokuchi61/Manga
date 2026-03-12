package comment

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/handler"
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
		r.Post("/", httpHandler.CreateComment)
		r.Get("/{comment_id}", httpHandler.GetDetail)
		r.Get("/{comment_id}/thread", httpHandler.GetThread)
		r.Patch("/{comment_id}", httpHandler.UpdateComment)
		r.Delete("/{comment_id}", httpHandler.DeleteComment)
		r.Post("/{comment_id}/restore", httpHandler.RestoreComment)
		r.Patch("/{comment_id}/moderation", httpHandler.UpdateModeration)
	})
}

package support

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/handler"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("support registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("support registerRoutes: http handler cannot be nil")
	}

	router.Route("/support", func(r chi.Router) {
		r.Get("/review/queue", httpHandler.ListReviewQueue)
		r.Get("/own", httpHandler.ListOwnSupport)

		r.Post("/communications", httpHandler.CreateCommunication)
		r.Post("/tickets", httpHandler.CreateTicket)
		r.Post("/reports", httpHandler.CreateReport)

		r.Get("/{support_id}", httpHandler.GetSupportDetail)
		r.Post("/{support_id}/replies", httpHandler.AddReply)
		r.Patch("/{support_id}/status", httpHandler.UpdateStatus)
		r.Post("/{support_id}/resolve", httpHandler.Resolve)
		r.Post("/{support_id}/handoff/moderation", httpHandler.RequestModerationHandoff)
	})
}

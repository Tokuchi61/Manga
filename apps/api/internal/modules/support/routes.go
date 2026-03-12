package support

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
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
		r.With(identity.RequireUser, identity.RequireAnyRole("support_agent", "moderator", "admin")).Get("/review/queue", httpHandler.ListReviewQueue)
		r.With(identity.RequireUser).Get("/own", httpHandler.ListOwnSupport)

		r.With(identity.RequireUser).Post("/communications", httpHandler.CreateCommunication)
		r.With(identity.RequireUser).Post("/tickets", httpHandler.CreateTicket)
		r.With(identity.RequireUser).Post("/reports", httpHandler.CreateReport)

		r.With(identity.RequireUser).Get("/{support_id}", httpHandler.GetSupportDetail)
		r.With(identity.RequireUser).Post("/{support_id}/replies", httpHandler.AddReply)
		r.With(identity.RequireUser, identity.RequireAnyRole("support_agent", "moderator", "admin")).Patch("/{support_id}/status", httpHandler.UpdateStatus)
		r.With(identity.RequireUser, identity.RequireAnyRole("support_agent", "moderator", "admin")).Post("/{support_id}/resolve", httpHandler.Resolve)
		r.With(identity.RequireUser, identity.RequireAnyRole("support_agent", "moderator", "admin")).Post("/{support_id}/handoff/moderation", httpHandler.RequestModerationHandoff)
	})
}

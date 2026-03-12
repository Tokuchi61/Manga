package moderation

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("moderation registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("moderation registerRoutes: http handler cannot be nil")
	}

	router.Route("/moderation", func(r chi.Router) {
		r.With(identity.RequireUser, identity.RequireAnyRole("moderator", "admin")).Get("/queue", httpHandler.ListQueue)
		r.With(identity.RequireUser, identity.RequireAnyRole("moderator", "admin", "support_agent")).Post("/cases/support-handoffs", httpHandler.CreateCaseFromSupportHandoff)
		r.With(identity.RequireUser, identity.RequireAnyRole("moderator", "admin")).Get("/cases/{case_id}", httpHandler.GetCaseDetail)
		r.With(identity.RequireUser, identity.RequireAnyRole("moderator", "admin")).Post("/cases/{case_id}/assign", httpHandler.AssignCase)
		r.With(identity.RequireUser, identity.RequireAnyRole("moderator", "admin")).Post("/cases/{case_id}/release", httpHandler.ReleaseCase)
		r.With(identity.RequireUser, identity.RequireAnyRole("moderator", "admin")).Post("/cases/{case_id}/notes", httpHandler.AddModeratorNote)
		r.With(identity.RequireUser, identity.RequireAnyRole("moderator", "admin")).Post("/cases/{case_id}/actions", httpHandler.ApplyAction)
		r.With(identity.RequireUser, identity.RequireAnyRole("moderator", "admin")).Post("/cases/{case_id}/escalate", httpHandler.EscalateCase)
	})
}

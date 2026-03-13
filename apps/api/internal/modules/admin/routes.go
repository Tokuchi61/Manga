package admin

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("admin registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("admin registerRoutes: http handler cannot be nil")
	}

	router.Route("/admin", func(r chi.Router) {
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/dashboard", httpHandler.GetDashboard)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/audit", httpHandler.ListAuditTrail)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/runtime/maintenance", httpHandler.UpdateMaintenanceState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/overrides", httpHandler.ListOverrides)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/overrides", httpHandler.ApplyOverride)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/user-reviews", httpHandler.ListUserReviews)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/user-reviews", httpHandler.ReviewUser)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/impersonations", httpHandler.ListImpersonationSessions)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/impersonations/start", httpHandler.StartImpersonation)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/impersonations/stop", httpHandler.StopImpersonation)
	})
}

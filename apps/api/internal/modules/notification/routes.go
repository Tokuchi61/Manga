package notification

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("notification registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("notification registerRoutes: http handler cannot be nil")
	}

	router.Route("/notifications", func(r chi.Router) {
		r.With(identity.RequireUser).Get("/inbox", httpHandler.ListInbox)
		r.With(identity.RequireUser).Get("/preferences", httpHandler.GetPreference)
		r.With(identity.RequireUser).Put("/preferences", httpHandler.UpdatePreference)
		r.With(identity.RequireUser, identity.RequireAnyRole("support_agent", "moderator", "admin")).Post("/intake/support", httpHandler.CreateFromSupportSignal)

		r.With(identity.RequireUser).Get("/{notification_id}", httpHandler.GetNotificationDetail)
		r.With(identity.RequireUser).Post("/{notification_id}/read", httpHandler.MarkRead)

		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/category-state", httpHandler.UpdateCategoryState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/channel-state", httpHandler.UpdateChannelState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/digest-state", httpHandler.UpdateDigestState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/delivery-pause", httpHandler.UpdateDeliveryPause)
	})
}

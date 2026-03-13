package ads

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("ads registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("ads registerRoutes: http handler cannot be nil")
	}

	router.Route("/ads", func(r chi.Router) {
		r.Get("/resolve", httpHandler.ResolvePlacements)
		r.Post("/impressions", httpHandler.IntakeImpression)
		r.Post("/clicks", httpHandler.IntakeClick)

		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/surface-state", httpHandler.UpdateSurfaceState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/placement-state", httpHandler.UpdatePlacementState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/campaign-state", httpHandler.UpdateCampaignState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/click-intake-state", httpHandler.UpdateClickIntakeState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/placements", httpHandler.ListPlacementDefinitions)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/placements", httpHandler.UpsertPlacementDefinition)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/campaigns", httpHandler.ListCampaignDefinitions)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/campaigns", httpHandler.UpsertCampaignDefinition)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/aggregate", httpHandler.ListCampaignAggregate)
	})
}

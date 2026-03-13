package royalpass

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("royalpass registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("royalpass registerRoutes: http handler cannot be nil")
	}

	router.Route("/royalpass", func(r chi.Router) {
		r.With(identity.RequireUser).Get("/overview", httpHandler.GetActorSeasonOverview)
		r.With(identity.RequireUser).Post("/progress/ingest", httpHandler.IngestRoyalPassProgress)
		r.With(identity.RequireUser).Post("/claims", httpHandler.ClaimTierReward)
		r.With(identity.RequireUser).Post("/premium/activate", httpHandler.ActivatePremiumTrack)

		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/season-state", httpHandler.UpdateSeasonState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/claim-state", httpHandler.UpdateClaimState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/premium-state", httpHandler.UpdatePremiumState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/seasons", httpHandler.ListSeasonDefinitions)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/seasons", httpHandler.UpsertSeasonDefinition)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/tiers", httpHandler.ListTierDefinitions)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/tiers", httpHandler.UpsertTierDefinition)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/reset-progress", httpHandler.ResetRoyalPassProgress)
	})
}

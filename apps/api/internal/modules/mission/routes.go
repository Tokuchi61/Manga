package mission

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("mission registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("mission registerRoutes: http handler cannot be nil")
	}

	router.Route("/missions", func(r chi.Router) {
		r.With(identity.RequireUser).Get("/", httpHandler.ListActorMissions)
		r.With(identity.RequireUser).Get("/{mission_id}", httpHandler.GetActorMissionDetail)
		r.With(identity.RequireUser).Post("/{mission_id}/progress/ingest", httpHandler.IngestMissionProgress)
		r.With(identity.RequireUser).Post("/{mission_id}/claim", httpHandler.ClaimMissionReward)

		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/read-state", httpHandler.UpdateReadState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/claim-state", httpHandler.UpdateClaimState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/progress-ingest-state", httpHandler.UpdateProgressIngestState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/reset-hour", httpHandler.UpdateDailyResetHour)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/definitions", httpHandler.ListMissionDefinitions)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/definitions", httpHandler.UpsertMissionDefinition)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/reset-progress", httpHandler.ResetMissionProgress)
	})
}

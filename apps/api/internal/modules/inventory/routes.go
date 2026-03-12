package inventory

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("inventory registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("inventory registerRoutes: http handler cannot be nil")
	}

	router.Route("/inventory", func(r chi.Router) {
		r.With(identity.RequireUser).Get("/items", httpHandler.ListInventoryEntries)
		r.With(identity.RequireUser).Get("/items/{item_id}", httpHandler.GetInventoryItemDetail)
		r.With(identity.RequireUser).Post("/claim", httpHandler.ClaimInventoryItem)
		r.With(identity.RequireUser).Post("/items/{item_id}/consume", httpHandler.ConsumeInventoryItem)
		r.With(identity.RequireUser).Post("/items/{item_id}/equip", httpHandler.UpdateEquipState)

		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/read-state", httpHandler.UpdateReadState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/claim-state", httpHandler.UpdateClaimState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/consume-state", httpHandler.UpdateConsumeState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/equip-state", httpHandler.UpdateEquipStateRuntime)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/items", httpHandler.ListItemDefinitions)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/items", httpHandler.UpsertItemDefinition)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/grants", httpHandler.AdminGrantItem)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/revokes", httpHandler.AdminRevokeItem)
	})
}

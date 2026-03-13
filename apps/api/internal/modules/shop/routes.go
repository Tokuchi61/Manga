package shop

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("shop registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("shop registerRoutes: http handler cannot be nil")
	}

	router.Route("/shop", func(r chi.Router) {
		r.With(identity.RequireUser).Get("/catalog", httpHandler.ListCatalog)
		r.With(identity.RequireUser).Get("/catalog/{product_id}", httpHandler.GetCatalogItem)
		r.With(identity.RequireUser).Post("/purchase/intents", httpHandler.CreatePurchaseIntent)
		r.With(identity.RequireUser).Post("/purchase/recovery", httpHandler.RequestPurchaseRecovery)

		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/catalog-state", httpHandler.UpdateCatalogState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/purchase-state", httpHandler.UpdatePurchaseState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/campaign-state", httpHandler.UpdateCampaignState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/products", httpHandler.ListProductDefinitions)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/products", httpHandler.UpsertProductDefinition)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/offers", httpHandler.ListOfferDefinitions)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/offers", httpHandler.UpsertOfferDefinition)
	})
}

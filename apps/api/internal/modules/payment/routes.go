package payment

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("payment registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("payment registerRoutes: http handler cannot be nil")
	}

	router.Route("/payment", func(r chi.Router) {
		r.With(identity.RequireUser).Get("/packages", httpHandler.ListManaPackages)
		r.With(identity.RequireUser).Post("/checkout/sessions", httpHandler.StartCheckoutSession)
		r.Post("/callback", httpHandler.ProcessProviderCallback)
		r.With(identity.RequireUser).Get("/wallet", httpHandler.GetOwnWallet)
		r.With(identity.RequireUser).Get("/transactions", httpHandler.ListOwnTransactions)

		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/mana-purchase-state", httpHandler.UpdateManaPurchaseState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/checkout-state", httpHandler.UpdateCheckoutState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/transaction-read-state", httpHandler.UpdateTransactionReadState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/callback-intake-state", httpHandler.UpdateCallbackIntakeState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/packages", httpHandler.ListAdminManaPackages)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/packages", httpHandler.UpsertManaPackage)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/reconcile", httpHandler.RunReconcile)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/refunds", httpHandler.ProcessRefund)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/reversals", httpHandler.ProcessReversal)
	})
}

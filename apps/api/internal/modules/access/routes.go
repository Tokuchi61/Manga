package access

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("access registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("access registerRoutes: http handler cannot be nil")
	}

	router.Route("/access", func(r chi.Router) {
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/roles", httpHandler.CreateRole)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/permissions", httpHandler.CreatePermission)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/roles/{role_id}/permissions", httpHandler.AssignPermissionToRole)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/users/{user_id}/roles", httpHandler.AssignRoleToUser)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/users/{user_id}/temporary-grants", httpHandler.CreateTemporaryGrant)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/policies", httpHandler.CreatePolicyRule)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin", "internal_service")).Post("/evaluate", httpHandler.Evaluate)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin", "internal_service")).Get("/contracts/permissions", httpHandler.ListCanonicalPermissions)
	})
}

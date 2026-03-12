package access

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/handler"
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
		r.Post("/roles", httpHandler.CreateRole)
		r.Post("/permissions", httpHandler.CreatePermission)
		r.Post("/roles/{role_id}/permissions", httpHandler.AssignPermissionToRole)
		r.Post("/users/{user_id}/roles", httpHandler.AssignRoleToUser)
		r.Post("/users/{user_id}/temporary-grants", httpHandler.CreateTemporaryGrant)
		r.Post("/policies", httpHandler.CreatePolicyRule)
		r.Post("/evaluate", httpHandler.Evaluate)
		r.Get("/contracts/permissions", httpHandler.ListCanonicalPermissions)
	})
}

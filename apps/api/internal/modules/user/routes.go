package user

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/handler"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("user registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("user registerRoutes: http handler cannot be nil")
	}

	router.Post("/users", httpHandler.CreateUser)
	router.Get("/users/{user_id}", httpHandler.GetPublicProfile)
	router.Get("/users/{user_id}/self", httpHandler.GetOwnProfile)
	router.Patch("/users/{user_id}/profile", httpHandler.UpdateProfile)
	router.Patch("/users/{user_id}/visibility", httpHandler.UpdateProfileVisibility)
	router.Patch("/users/{user_id}/history-visibility", httpHandler.UpdateHistoryVisibility)

	router.Patch("/users/{user_id}/account/state", httpHandler.UpdateAccountState)
	router.Patch("/users/{user_id}/vip", httpHandler.UpdateVIPState)
}

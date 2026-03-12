package routes

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/handler"
	"github.com/go-chi/chi/v5"
)

func registerAccountRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	router.Patch("/users/{user_id}/account/state", httpHandler.UpdateAccountState)
	router.Patch("/users/{user_id}/vip", httpHandler.UpdateVIPState)
}

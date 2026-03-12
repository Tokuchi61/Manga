package routes

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	"github.com/go-chi/chi/v5"
)

func registerIdentityRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	router.Post("/register", httpHandler.Register)
	router.Post("/login", httpHandler.Login)
	router.Post("/token/refresh", httpHandler.RefreshToken)
	router.Post("/logout", httpHandler.Logout)
}

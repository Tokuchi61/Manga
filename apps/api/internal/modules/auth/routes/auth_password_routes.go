package routes

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	"github.com/go-chi/chi/v5"
)

func registerPasswordRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	router.Post("/password/forgot", httpHandler.ForgotPassword)
	router.Post("/password/reset", httpHandler.ResetPassword)
	router.Post("/password/change", httpHandler.ChangePassword)
}

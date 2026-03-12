package routes

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	"github.com/go-chi/chi/v5"
)

func registerVerificationRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	router.Post("/email/verification/send", httpHandler.SendEmailVerification)
	router.Post("/email/verification/confirm", httpHandler.ConfirmEmailVerification)
}

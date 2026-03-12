package auth

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("auth registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("auth registerRoutes: http handler cannot be nil")
	}

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", httpHandler.Register)
		r.Post("/login", httpHandler.Login)
		r.Post("/token/refresh", httpHandler.RefreshToken)
		r.Post("/logout", httpHandler.Logout)

		r.Get("/sessions", httpHandler.ListSessions)
		r.Post("/sessions/revoke/current", httpHandler.RevokeCurrentSession)
		r.Post("/sessions/revoke/others", httpHandler.RevokeOtherSessions)
		r.Post("/sessions/revoke/all", httpHandler.RevokeAllSessions)

		r.Post("/password/forgot", httpHandler.ForgotPassword)
		r.Post("/password/reset", httpHandler.ResetPassword)
		r.Post("/password/change", httpHandler.ChangePassword)

		r.Post("/email/verification/send", httpHandler.SendEmailVerification)
		r.Post("/email/verification/confirm", httpHandler.ConfirmEmailVerification)
	})
}

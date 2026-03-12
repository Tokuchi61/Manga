package auth

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
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

		r.With(identity.RequireCredential).Post("/logout", httpHandler.Logout)
		r.With(identity.RequireCredential).Get("/sessions", httpHandler.ListSessions)
		r.With(identity.RequireCredential).Post("/sessions/revoke/current", httpHandler.RevokeCurrentSession)
		r.With(identity.RequireCredential).Post("/sessions/revoke/others", httpHandler.RevokeOtherSessions)
		r.With(identity.RequireCredential).Post("/sessions/revoke/all", httpHandler.RevokeAllSessions)

		r.Post("/password/forgot", httpHandler.ForgotPassword)
		r.Post("/password/reset", httpHandler.ResetPassword)
		r.With(identity.RequireCredential).Post("/password/change", httpHandler.ChangePassword)

		r.With(identity.RequireCredential).Post("/email/verification/send", httpHandler.SendEmailVerification)
		r.Post("/email/verification/confirm", httpHandler.ConfirmEmailVerification)
	})
}

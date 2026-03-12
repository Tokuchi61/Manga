package routes

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	"github.com/go-chi/chi/v5"
)

func Register(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil || httpHandler == nil {
		return
	}

	router.Route("/auth", func(r chi.Router) {
		registerIdentityRoutes(r, httpHandler)
		registerSessionRoutes(r, httpHandler)
		registerPasswordRoutes(r, httpHandler)
		registerVerificationRoutes(r, httpHandler)
	})
}

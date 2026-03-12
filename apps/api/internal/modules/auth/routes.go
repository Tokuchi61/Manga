package auth

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	authroutes "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/routes"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	authroutes.Register(router, httpHandler)
}

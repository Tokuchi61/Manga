package user

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/handler"
	userroutes "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/routes"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	userroutes.Register(router, httpHandler)
}

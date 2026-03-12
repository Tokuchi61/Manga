package routes

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/handler"
	"github.com/go-chi/chi/v5"
)

func Register(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil || httpHandler == nil {
		return
	}

	registerProfileRoutes(router, httpHandler)
	registerAccountRoutes(router, httpHandler)
}

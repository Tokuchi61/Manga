package routes

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	"github.com/go-chi/chi/v5"
)

func registerSessionRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	router.Get("/sessions", httpHandler.ListSessions)
	router.Post("/sessions/revoke/current", httpHandler.RevokeCurrentSession)
	router.Post("/sessions/revoke/others", httpHandler.RevokeOtherSessions)
	router.Post("/sessions/revoke/all", httpHandler.RevokeAllSessions)
}

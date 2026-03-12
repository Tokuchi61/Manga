package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/service"

// HTTPHandler exposes auth HTTP endpoints.
type HTTPHandler struct {
	service *service.AuthService
}

func New(svc *service.AuthService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

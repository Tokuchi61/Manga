package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/service"

// HTTPHandler exposes support HTTP endpoints.
type HTTPHandler struct {
	service *service.SupportService
}

func New(svc *service.SupportService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

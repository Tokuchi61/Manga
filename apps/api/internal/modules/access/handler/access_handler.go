package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/service"

// HTTPHandler exposes access HTTP endpoints.
type HTTPHandler struct {
	service *service.AccessService
}

func New(svc *service.AccessService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

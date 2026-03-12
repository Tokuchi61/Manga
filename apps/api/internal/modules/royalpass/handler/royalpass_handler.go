package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/service"

// HTTPHandler exposes royalpass HTTP endpoints.
type HTTPHandler struct {
	service *service.RoyalPassService
}

func New(svc *service.RoyalPassService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

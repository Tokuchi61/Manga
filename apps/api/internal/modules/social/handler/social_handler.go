package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/social/service"

// HTTPHandler exposes social HTTP endpoints.
type HTTPHandler struct {
	service *service.SocialService
}

func New(svc *service.SocialService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

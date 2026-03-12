package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/service"

// HTTPHandler exposes moderation HTTP endpoints.
type HTTPHandler struct {
	service *service.ModerationService
}

func New(svc *service.ModerationService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

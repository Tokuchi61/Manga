package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/service"

// HTTPHandler exposes notification HTTP endpoints.
type HTTPHandler struct {
	service *service.NotificationService
}

func New(svc *service.NotificationService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

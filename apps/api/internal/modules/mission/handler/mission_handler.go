package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/service"

// HTTPHandler exposes mission HTTP endpoints.
type HTTPHandler struct {
	service *service.MissionService
}

func New(svc *service.MissionService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

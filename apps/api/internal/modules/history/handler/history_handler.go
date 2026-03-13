package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/history/service"

// HTTPHandler exposes history HTTP endpoints.
type HTTPHandler struct {
	service *service.HistoryService
}

func New(svc *service.HistoryService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

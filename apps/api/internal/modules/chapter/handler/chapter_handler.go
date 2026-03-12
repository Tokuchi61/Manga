package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/service"

// HTTPHandler exposes chapter HTTP endpoints.
type HTTPHandler struct {
	service *service.ChapterService
}

func New(svc *service.ChapterService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/service"

// HTTPHandler exposes manga HTTP endpoints.
type HTTPHandler struct {
	service *service.MangaService
}

func New(svc *service.MangaService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

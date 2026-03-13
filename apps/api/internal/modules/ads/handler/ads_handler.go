package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/service"

// HTTPHandler exposes ads HTTP endpoints.
type HTTPHandler struct {
	service *service.AdsService
}

func New(svc *service.AdsService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

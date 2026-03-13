package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/service"

// HTTPHandler exposes shop HTTP endpoints.
type HTTPHandler struct {
	service *service.ShopService
}

func New(svc *service.ShopService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

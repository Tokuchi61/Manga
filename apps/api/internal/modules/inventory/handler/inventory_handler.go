package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/service"

// HTTPHandler exposes inventory HTTP endpoints.
type HTTPHandler struct {
	service *service.InventoryService
}

func New(svc *service.InventoryService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

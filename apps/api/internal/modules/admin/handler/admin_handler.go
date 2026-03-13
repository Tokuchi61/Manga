package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/service"

// HTTPHandler exposes admin HTTP endpoints.
type HTTPHandler struct {
	service *service.AdminService
}

func New(svc *service.AdminService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

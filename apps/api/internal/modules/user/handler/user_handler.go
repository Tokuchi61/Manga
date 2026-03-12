package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/service"

// HTTPHandler exposes user HTTP endpoints.
type HTTPHandler struct {
	service *service.UserService
}

func New(svc *service.UserService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/service"

// HTTPHandler exposes comment HTTP endpoints.
type HTTPHandler struct {
	service *service.CommentService
}

func New(svc *service.CommentService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

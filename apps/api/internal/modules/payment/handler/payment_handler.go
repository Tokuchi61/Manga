package handler

import "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/service"

// HTTPHandler exposes payment HTTP endpoints.
type HTTPHandler struct {
	service *service.PaymentService
}

func New(svc *service.PaymentService) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

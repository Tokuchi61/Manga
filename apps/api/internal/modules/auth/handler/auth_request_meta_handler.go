package handler

import (
	"net"
	"net/http"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/service"
	"github.com/go-chi/chi/v5/middleware"
)

func buildRequestMeta(r *http.Request) service.RequestMeta {
	requestID := middleware.GetReqID(r.Context())
	correlationID := strings.TrimSpace(r.Header.Get("X-Correlation-ID"))
	if correlationID == "" {
		correlationID = requestID
	}
	device := strings.TrimSpace(r.Header.Get("X-Device-ID"))
	if device == "" {
		device = strings.TrimSpace(r.UserAgent())
	}

	return service.RequestMeta{
		RequestID:     requestID,
		CorrelationID: correlationID,
		Device:        device,
		IP:            extractIP(r),
	}
}

func extractIP(r *http.Request) string {
	forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}
	return host
}

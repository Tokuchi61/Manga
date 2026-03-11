package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHealthEndpoint(t *testing.T) {
	handler := NewHTTPHandler(config.Config{AppVersion: "0.1.0-test"}, zap.NewNop(), nil, modules.EmptyRegistry())
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "\"status\":\"ok\"")
}

func TestVersionEndpoint(t *testing.T) {
	handler := NewHTTPHandler(config.Config{AppVersion: "0.1.0-test"}, zap.NewNop(), nil, modules.EmptyRegistry())
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "0.1.0-test")
}

func TestReadyEndpointWithoutDB(t *testing.T) {
	handler := NewHTTPHandler(config.Config{AppVersion: "0.1.0-test"}, zap.NewNop(), nil, modules.EmptyRegistry())
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.Contains(t, rec.Body.String(), "unavailable")
}

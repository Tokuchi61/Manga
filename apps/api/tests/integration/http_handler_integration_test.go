package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type integrationModule struct{}

func (integrationModule) Name() string {
	return "integration_sample"
}

func (integrationModule) RegisterRoutes(router chi.Router) {
	router.Get("/integration-sample/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"module_ok"}`))
	})
}

func TestHTTPHandlerMountsModuleRoutes(t *testing.T) {
	registry, err := modules.NewRegistry(integrationModule{})
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.1.0-test"}, zap.NewNop(), nil, registry)

	req := httptest.NewRequest(http.MethodGet, "/integration-sample/ping", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "module_ok")
}

func TestHTTPHandlerKeepsCoreRoutes(t *testing.T) {
	registry, err := modules.NewRegistry(integrationModule{})
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.1.0-test"}, zap.NewNop(), nil, registry)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

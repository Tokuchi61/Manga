package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	adminmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMaintenanceModeBlocksPublicRoutes(t *testing.T) {
	adminUserID := uuid.NewString()
	adminCredentialID := uuid.NewString()

	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		adminmodule.New(),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.21.1-test"}, zap.NewNop(), nil, registry)

	enableRec := performMaintenanceJSONRequest(t, handler, http.MethodPost, "/admin/runtime/maintenance", map[string]any{
		"request_id":         "req-maintenance-enable-1",
		"enabled":            true,
		"reason":             "maintenance_window",
		"risk_level":         "high",
		"double_confirmed":   true,
		"confirmation_token": "confirm-maintenance",
		"correlation_id":     "corr-maintenance-enable-1",
	}, adminUserID, adminCredentialID, "admin")
	require.Equal(t, http.StatusOK, enableRec.Code)

	registerRec := performMaintenanceJSONRequest(t, handler, http.MethodPost, "/auth/register", map[string]any{
		"email":    "maintenance@example.com",
		"password": "StrongPass123!",
	}, "", "", "")
	require.Equal(t, http.StatusServiceUnavailable, registerRec.Code)

	healthRec := performMaintenanceJSONRequest(t, handler, http.MethodGet, "/health", nil, "", "", "")
	require.Equal(t, http.StatusOK, healthRec.Code)

	adminRuntimeRec := performMaintenanceJSONRequest(t, handler, http.MethodGet, "/admin/runtime", nil, adminUserID, adminCredentialID, "admin")
	require.Equal(t, http.StatusOK, adminRuntimeRec.Code)
}

func performMaintenanceJSONRequest(t *testing.T, handler http.Handler, method string, path string, payload any, actorUserID string, actorCredentialID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	body := bytes.NewReader(nil)
	if payload != nil {
		raw, err := json.Marshal(payload)
		require.NoError(t, err)
		body = bytes.NewReader(raw)
	}

	req := httptest.NewRequest(method, path, body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if actorUserID != "" || actorCredentialID != "" || roles != "" {
		setActorHeaders(req, actorUserID, actorCredentialID, roles)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

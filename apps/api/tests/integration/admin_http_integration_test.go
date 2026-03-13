package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	accessmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/access"
	adminmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAdminHTTPRuntimeReviewOverrideImpersonationFlow(t *testing.T) {
	adminID := uuid.NewString()
	targetUserID := uuid.NewString()

	admin := adminmodule.New()
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		admin,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.21.0-test"}, zap.NewNop(), nil, registry)

	unauthorizedDashboardReq := performAdminRequest(t, handler, http.MethodGet, "/admin/dashboard", nil, "", "")
	require.Equal(t, http.StatusUnauthorized, unauthorizedDashboardReq.Code)

	forbiddenDashboardReq := performAdminRequest(t, handler, http.MethodGet, "/admin/dashboard", nil, adminID, "user")
	require.Equal(t, http.StatusForbidden, forbiddenDashboardReq.Code)

	dashboardReq := performAdminRequest(t, handler, http.MethodGet, "/admin/dashboard", nil, adminID, "admin")
	require.Equal(t, http.StatusOK, dashboardReq.Code)

	maintenanceNeedConfirmReq := performAdminJSONRequest(t, handler, http.MethodPost, "/admin/runtime/maintenance", map[string]any{
		"request_id": "req-maint-1",
		"enabled":    true,
		"reason":     "maintenance",
		"risk_level": "high",
	}, adminID, "admin")
	require.Equal(t, http.StatusPreconditionFailed, maintenanceNeedConfirmReq.Code)

	maintenanceReq := performAdminJSONRequest(t, handler, http.MethodPost, "/admin/runtime/maintenance", map[string]any{
		"request_id":         "req-maint-1",
		"enabled":            true,
		"reason":             "maintenance",
		"risk_level":         "high",
		"double_confirmed":   true,
		"confirmation_token": "confirm-maint",
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, maintenanceReq.Code)
	var maintenanceRes struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(maintenanceReq.Body.Bytes(), &maintenanceRes))
	require.Equal(t, "accepted", maintenanceRes.Status)

	runtimeReq := performAdminRequest(t, handler, http.MethodGet, "/admin/runtime", nil, adminID, "admin")
	require.Equal(t, http.StatusOK, runtimeReq.Code)
	var runtimeRes struct {
		MaintenanceEnabled bool `json:"maintenance_enabled"`
	}
	require.NoError(t, json.Unmarshal(runtimeReq.Body.Bytes(), &runtimeRes))
	require.True(t, runtimeRes.MaintenanceEnabled)

	reviewReq := performAdminJSONRequest(t, handler, http.MethodPost, "/admin/user-reviews", map[string]any{
		"request_id":     "req-review-1",
		"target_user_id": targetUserID,
		"decision":       "warning",
		"reason":         "policy",
		"risk_level":     "low",
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, reviewReq.Code)

	overrideReq := performAdminJSONRequest(t, handler, http.MethodPost, "/admin/overrides", map[string]any{
		"request_id":         "req-override-1",
		"target_module":      "moderation",
		"target_type":        "case",
		"target_id":          "case-1",
		"decision":           "freeze",
		"reason":             "critical",
		"risk_level":         "critical",
		"double_confirmed":   true,
		"confirmation_token": "confirm-override",
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, overrideReq.Code)
	var overrideRes struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(overrideReq.Body.Bytes(), &overrideRes))
	require.Equal(t, "accepted", overrideRes.Status)

	overrideIdempotentReq := performAdminJSONRequest(t, handler, http.MethodPost, "/admin/overrides", map[string]any{
		"request_id":         "req-override-1",
		"target_module":      "moderation",
		"target_type":        "case",
		"target_id":          "case-1",
		"decision":           "freeze",
		"reason":             "critical",
		"risk_level":         "critical",
		"double_confirmed":   true,
		"confirmation_token": "confirm-override",
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, overrideIdempotentReq.Code)
	require.NoError(t, json.Unmarshal(overrideIdempotentReq.Body.Bytes(), &overrideRes))
	require.Equal(t, "idempotent", overrideRes.Status)

	impersonationNeedConfirmReq := performAdminJSONRequest(t, handler, http.MethodPost, "/admin/impersonations/start", map[string]any{
		"request_id":     "req-imp-start-1",
		"target_user_id": targetUserID,
		"reason":         "investigation",
		"risk_level":     "high",
	}, adminID, "admin")
	require.Equal(t, http.StatusPreconditionFailed, impersonationNeedConfirmReq.Code)

	impersonationStartReq := performAdminJSONRequest(t, handler, http.MethodPost, "/admin/impersonations/start", map[string]any{
		"request_id":         "req-imp-start-1",
		"target_user_id":     targetUserID,
		"reason":             "investigation",
		"risk_level":         "high",
		"double_confirmed":   true,
		"confirmation_token": "confirm-start",
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, impersonationStartReq.Code)
	var impersonationStartRes struct {
		Status    string `json:"status"`
		SessionID string `json:"session_id"`
	}
	require.NoError(t, json.Unmarshal(impersonationStartReq.Body.Bytes(), &impersonationStartRes))
	require.Equal(t, "started", impersonationStartRes.Status)
	require.NotEmpty(t, impersonationStartRes.SessionID)

	impersonationListReq := performAdminRequest(t, handler, http.MethodGet, "/admin/impersonations?active_only=true", nil, adminID, "admin")
	require.Equal(t, http.StatusOK, impersonationListReq.Code)
	var impersonationListRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(impersonationListReq.Body.Bytes(), &impersonationListRes))
	require.Equal(t, 1, impersonationListRes.Count)

	impersonationStopReq := performAdminJSONRequest(t, handler, http.MethodPost, "/admin/impersonations/stop", map[string]any{
		"request_id":         "req-imp-stop-1",
		"session_id":         impersonationStartRes.SessionID,
		"reason":             "done",
		"risk_level":         "high",
		"double_confirmed":   true,
		"confirmation_token": "confirm-stop",
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, impersonationStopReq.Code)

	auditReq := performAdminRequest(t, handler, http.MethodGet, "/admin/audit", nil, adminID, "admin")
	require.Equal(t, http.StatusOK, auditReq.Code)
	var auditRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(auditReq.Body.Bytes(), &auditRes))
	require.GreaterOrEqual(t, auditRes.Count, 5)
}

func performAdminJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performAdminRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performAdminRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		reader = body
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if actorUserID != "" || roles != "" {
		setActorHeaders(req, actorUserID, "", roles)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

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
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	notificationmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/notification"
	supportmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/support"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNotificationHTTPSupportIntakeInboxPreferenceAndRuntimeFlow(t *testing.T) {
	requesterID := uuid.NewString()
	supportAgentID := uuid.NewString()
	adminID := uuid.NewString()

	support := supportmodule.New()
	notification := notificationmodule.New()
	notification.SetSupportSignalProvider(support)

	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		support,
		notification,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.12.0-test"}, zap.NewNop(), nil, registry)

	supportID := createSupportReportForNotificationTest(t, handler, requesterID, "req-stage12-support-1")

	intakeRec := performNotificationJSONRequest(t, handler, http.MethodPost, "/notifications/intake/support", map[string]any{
		"support_id": supportID,
	}, supportAgentID, "support_agent")
	require.Equal(t, http.StatusOK, intakeRec.Code)

	var intakeRes struct {
		NotificationID string `json:"notification_id"`
		Created        bool   `json:"created"`
	}
	require.NoError(t, json.Unmarshal(intakeRec.Body.Bytes(), &intakeRes))
	require.True(t, intakeRes.Created)
	require.NotEmpty(t, intakeRes.NotificationID)

	intakeAgainRec := performNotificationJSONRequest(t, handler, http.MethodPost, "/notifications/intake/support", map[string]any{
		"support_id": supportID,
	}, supportAgentID, "support_agent")
	require.Equal(t, http.StatusOK, intakeAgainRec.Code)

	var intakeAgainRes struct {
		NotificationID string `json:"notification_id"`
		Created        bool   `json:"created"`
	}
	require.NoError(t, json.Unmarshal(intakeAgainRec.Body.Bytes(), &intakeAgainRes))
	require.False(t, intakeAgainRes.Created)
	require.Equal(t, intakeRes.NotificationID, intakeAgainRes.NotificationID)

	inboxRec := performNotificationRequest(t, handler, http.MethodGet, "/notifications/inbox", nil, requesterID, "")
	require.Equal(t, http.StatusOK, inboxRec.Code)
	var inboxRes struct {
		Count       int `json:"count"`
		UnreadCount int `json:"unread_count"`
	}
	require.NoError(t, json.Unmarshal(inboxRec.Body.Bytes(), &inboxRes))
	require.Equal(t, 1, inboxRes.Count)
	require.Equal(t, 1, inboxRes.UnreadCount)

	markReadRec := performNotificationRequest(t, handler, http.MethodPost, "/notifications/"+intakeRes.NotificationID+"/read", bytes.NewReader([]byte("{}")), requesterID, "")
	require.Equal(t, http.StatusOK, markReadRec.Code)

	unreadInboxRec := performNotificationRequest(t, handler, http.MethodGet, "/notifications/inbox?unread_only=true", nil, requesterID, "")
	require.Equal(t, http.StatusOK, unreadInboxRec.Code)
	var unreadInboxRes struct {
		Count       int `json:"count"`
		UnreadCount int `json:"unread_count"`
	}
	require.NoError(t, json.Unmarshal(unreadInboxRec.Body.Bytes(), &unreadInboxRes))
	require.Equal(t, 0, unreadInboxRes.Count)
	require.Equal(t, 0, unreadInboxRes.UnreadCount)

	preferenceUpdateRec := performNotificationJSONRequest(t, handler, http.MethodPut, "/notifications/preferences", map[string]any{
		"muted_categories":   []string{"support"},
		"quiet_hours_enabled": true,
		"quiet_hours_start":   23,
		"quiet_hours_end":     7,
		"in_app_enabled":      true,
		"email_enabled":       true,
		"push_enabled":        false,
		"digest_enabled":      true,
	}, requesterID, "")
	require.Equal(t, http.StatusOK, preferenceUpdateRec.Code)

	supportMutedID := createSupportReportForNotificationTest(t, handler, requesterID, "req-stage12-support-2")
	mutedIntakeRec := performNotificationJSONRequest(t, handler, http.MethodPost, "/notifications/intake/support", map[string]any{
		"support_id": supportMutedID,
	}, supportAgentID, "support_agent")
	require.Equal(t, http.StatusConflict, mutedIntakeRec.Code)

	deliveryPauseRec := performNotificationJSONRequest(t, handler, http.MethodPost, "/notifications/admin/delivery-pause", map[string]any{
		"paused": true,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, deliveryPauseRec.Code)

	supportPausedID := createSupportReportForNotificationTest(t, handler, requesterID, "req-stage12-support-3")
	pausedIntakeRec := performNotificationJSONRequest(t, handler, http.MethodPost, "/notifications/intake/support", map[string]any{
		"support_id": supportPausedID,
	}, supportAgentID, "support_agent")
	require.Equal(t, http.StatusServiceUnavailable, pausedIntakeRec.Code)
}

func createSupportReportForNotificationTest(t *testing.T, handler http.Handler, requesterID string, requestID string) string {
	t.Helper()

	rec := performNotificationJSONRequest(t, handler, http.MethodPost, "/support/reports", map[string]any{
		"category":    "content",
		"priority":    "high",
		"reason_code": "abuse",
		"reason_text": "Reported abusive comment",
		"target_type": "comment",
		"target_id":   uuid.NewString(),
		"request_id":  requestID,
	}, requesterID, "")
	require.Equal(t, http.StatusCreated, rec.Code)

	var res struct {
		SupportID string `json:"support_id"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotEmpty(t, res.SupportID)
	return res.SupportID
}

func performNotificationJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performNotificationRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performNotificationRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	bodyReader := bytes.NewReader(nil)
	if body != nil {
		bodyReader = body
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	setActorHeaders(req, actorUserID, "", roles)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	return rec
}

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
	moderationmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation"
	supportmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/support"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestModerationHTTPSupportHandoffQueueLifecycle(t *testing.T) {
	requesterID := uuid.NewString()
	supportAgentID := uuid.NewString()
	moderatorID := uuid.NewString()
	targetID := uuid.NewString()

	support := supportmodule.New()
	moderation := moderationmodule.New()
	moderation.SetSupportContracts(support, support)

	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		support,
		moderation,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.11.0-test"}, zap.NewNop(), nil, registry)

	createReportRec := performModerationJSONRequest(t, handler, http.MethodPost, "/support/reports", map[string]any{
		"category":    "content",
		"priority":    "high",
		"reason_code": "abuse",
		"reason_text": "Reported abusive comment",
		"target_type": "comment",
		"target_id":   targetID,
		"request_id":  "req-stage11-support-1",
	}, requesterID, "")
	require.Equal(t, http.StatusCreated, createReportRec.Code)

	var reportRes struct {
		SupportID string `json:"support_id"`
	}
	require.NoError(t, json.Unmarshal(createReportRec.Body.Bytes(), &reportRes))
	require.NotEmpty(t, reportRes.SupportID)

	handoffRec := performModerationJSONRequest(t, handler, http.MethodPost, "/support/"+reportRes.SupportID+"/handoff/moderation", map[string]any{}, supportAgentID, "support_agent")
	require.Equal(t, http.StatusOK, handoffRec.Code)

	createCaseRec := performModerationJSONRequest(t, handler, http.MethodPost, "/moderation/cases/support-handoffs", map[string]any{
		"support_id": reportRes.SupportID,
	}, moderatorID, "moderator")
	require.Equal(t, http.StatusCreated, createCaseRec.Code)

	var createCaseRes struct {
		Created bool `json:"created"`
		Case    struct {
			CaseID string `json:"case_id"`
			Status string `json:"status"`
		} `json:"case"`
	}
	require.NoError(t, json.Unmarshal(createCaseRec.Body.Bytes(), &createCaseRes))
	require.True(t, createCaseRes.Created)
	require.Equal(t, "queued", createCaseRes.Case.Status)
	require.NotEmpty(t, createCaseRes.Case.CaseID)

	createCaseAgainRec := performModerationJSONRequest(t, handler, http.MethodPost, "/moderation/cases/support-handoffs", map[string]any{
		"support_id": reportRes.SupportID,
	}, moderatorID, "moderator")
	require.Equal(t, http.StatusOK, createCaseAgainRec.Code)

	queueRec := performModerationRequest(t, handler, http.MethodGet, "/moderation/queue?status=queued", nil, moderatorID, "moderator")
	require.Equal(t, http.StatusOK, queueRec.Code)

	assignRec := performModerationJSONRequest(t, handler, http.MethodPost, "/moderation/cases/"+createCaseRes.Case.CaseID+"/assign", map[string]any{}, moderatorID, "moderator")
	require.Equal(t, http.StatusOK, assignRec.Code)

	noteRec := performModerationJSONRequest(t, handler, http.MethodPost, "/moderation/cases/"+createCaseRes.Case.CaseID+"/notes", map[string]any{
		"body": "needs escalation",
	}, moderatorID, "moderator")
	require.Equal(t, http.StatusOK, noteRec.Code)

	actionRec := performModerationJSONRequest(t, handler, http.MethodPost, "/moderation/cases/"+createCaseRes.Case.CaseID+"/actions", map[string]any{
		"action_type": "hide",
		"reason_code": "abuse",
		"summary":     "hidden while reviewing",
	}, moderatorID, "moderator")
	require.Equal(t, http.StatusOK, actionRec.Code)

	escalateRec := performModerationJSONRequest(t, handler, http.MethodPost, "/moderation/cases/"+createCaseRes.Case.CaseID+"/escalate", map[string]any{
		"reason": "needs admin final decision",
	}, moderatorID, "moderator")
	require.Equal(t, http.StatusOK, escalateRec.Code)

	detailRec := performModerationRequest(t, handler, http.MethodGet, "/moderation/cases/"+createCaseRes.Case.CaseID, nil, moderatorID, "moderator")
	require.Equal(t, http.StatusOK, detailRec.Code)
	var detailRes struct {
		Status           string `json:"status"`
		EscalationStatus string `json:"escalation_status"`
	}
	require.NoError(t, json.Unmarshal(detailRec.Body.Bytes(), &detailRes))
	require.Equal(t, "escalated", detailRes.Status)
	require.Equal(t, "pending_admin", detailRes.EscalationStatus)
}

func performModerationJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performModerationRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performModerationRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
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

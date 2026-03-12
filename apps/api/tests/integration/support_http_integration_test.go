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
	supportmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/support"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSupportHTTPIntakeReviewLifecycleFlow(t *testing.T) {
	requesterID := uuid.NewString()
	agentID := uuid.NewString()
	reviewerID := uuid.NewString()
	targetID := uuid.NewString()

	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		supportmodule.New(),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.10.0-test"}, zap.NewNop(), nil, registry)

	createCommunicationRec := performSupportJSONRequest(t, handler, http.MethodPost, "/support/communications", map[string]any{
		"requester_user_id": requesterID,
		"category":          "communication",
		"priority":          "normal",
		"reason_text":       "General support inquiry",
		"request_id":        "req-support-http-1",
	})
	require.Equal(t, http.StatusCreated, createCommunicationRec.Code)

	createReportRec := performSupportJSONRequest(t, handler, http.MethodPost, "/support/reports", map[string]any{
		"requester_user_id": requesterID,
		"category":          "content",
		"priority":          "high",
		"reason_code":       "abuse",
		"reason_text":       "Reported abusive comment",
		"target_type":       "comment",
		"target_id":         targetID,
		"request_id":        "req-support-http-2",
	})
	require.Equal(t, http.StatusCreated, createReportRec.Code)

	var reportRes struct {
		SupportID string `json:"support_id"`
	}
	require.NoError(t, json.Unmarshal(createReportRec.Body.Bytes(), &reportRes))
	require.NotEmpty(t, reportRes.SupportID)

	listRec := performSupportRequest(t, handler, http.MethodGet, "/support/own?requester_user_id="+requesterID, nil)
	require.Equal(t, http.StatusOK, listRec.Code)
	var listRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(listRec.Body.Bytes(), &listRes))
	require.Equal(t, 2, listRes.Count)

	detailRec := performSupportRequest(t, handler, http.MethodGet, "/support/"+reportRes.SupportID+"?requester_user_id="+requesterID, nil)
	require.Equal(t, http.StatusOK, detailRec.Code)

	replyRec := performSupportJSONRequest(t, handler, http.MethodPost, "/support/"+reportRes.SupportID+"/replies", map[string]any{
		"actor_user_id": requesterID,
		"message":       "Any update from team?",
		"visibility":    "public_to_requester",
	})
	require.Equal(t, http.StatusOK, replyRec.Code)

	statusRec := performSupportJSONRequest(t, handler, http.MethodPatch, "/support/"+reportRes.SupportID+"/status", map[string]any{
		"status":              "triaged",
		"assignee_user_id":    agentID,
		"reviewed_by_user_id": reviewerID,
	})
	require.Equal(t, http.StatusOK, statusRec.Code)

	resolveRec := performSupportJSONRequest(t, handler, http.MethodPost, "/support/"+reportRes.SupportID+"/resolve", map[string]any{
		"reviewed_by_user_id": reviewerID,
		"resolution_note":     "Case resolved.",
	})
	require.Equal(t, http.StatusOK, resolveRec.Code)

	queueRec := performSupportRequest(t, handler, http.MethodGet, "/support/review/queue?status=resolved", nil)
	require.Equal(t, http.StatusOK, queueRec.Code)
	var queueRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(queueRec.Body.Bytes(), &queueRes))
	require.Equal(t, 1, queueRes.Count)

	handoffRec := performSupportJSONRequest(t, handler, http.MethodPost, "/support/"+reportRes.SupportID+"/handoff/moderation", map[string]any{})
	require.Equal(t, http.StatusOK, handoffRec.Code)

	handoffAgainRec := performSupportJSONRequest(t, handler, http.MethodPost, "/support/"+reportRes.SupportID+"/handoff/moderation", map[string]any{})
	require.Equal(t, http.StatusConflict, handoffAgainRec.Code)
}

func TestSupportHTTPCreateValidationFailure(t *testing.T) {
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		supportmodule.New(),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.10.0-test"}, zap.NewNop(), nil, registry)

	invalidCreateRec := performSupportJSONRequest(t, handler, http.MethodPost, "/support/reports", map[string]any{
		"requester_user_id": "not-a-uuid",
		"category":          "content",
		"reason_text":       "invalid input",
		"target_type":       "social",
		"target_id":         "not-a-uuid",
	})
	require.Equal(t, http.StatusBadRequest, invalidCreateRec.Code)
}

func performSupportJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performSupportRequest(t, handler, method, path, bytes.NewReader(payload))
}

func performSupportRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader) *httptest.ResponseRecorder {
	t.Helper()

	bodyReader := bytes.NewReader(nil)
	if body != nil {
		bodyReader = body
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	return rec
}

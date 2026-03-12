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
	missionmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/mission"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMissionHTTPProgressClaimRuntimeAndResetFlow(t *testing.T) {
	actorID := uuid.NewString()
	adminID := uuid.NewString()

	mission := missionmodule.New()
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		mission,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.16.0-test"}, zap.NewNop(), nil, registry)

	upsertDefinitionReq := performMissionJSONRequest(t, handler, http.MethodPost, "/missions/admin/definitions", map[string]any{
		"mission_id":      "daily_read_3",
		"category":        "daily",
		"title":           "Read 3 chapters",
		"objective_type":  "chapter_read",
		"target_count":    3,
		"reward_item_id":  "mana_potion",
		"reward_quantity": 1,
		"active":          true,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertDefinitionReq.Code)

	ingestReq := performMissionJSONRequest(t, handler, http.MethodPost, "/missions/daily_read_3/progress/ingest", map[string]any{
		"delta":       2,
		"source_type": "history",
		"request_id":  "req-stage16-progress-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, ingestReq.Code)

	ingestCompleteReq := performMissionJSONRequest(t, handler, http.MethodPost, "/missions/daily_read_3/progress/ingest", map[string]any{
		"delta":       1,
		"source_type": "history",
		"request_id":  "req-stage16-progress-2",
	}, actorID, "")
	require.Equal(t, http.StatusOK, ingestCompleteReq.Code)

	claimReq := performMissionJSONRequest(t, handler, http.MethodPost, "/missions/daily_read_3/claim", map[string]any{
		"request_id": "req-stage16-claim-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, claimReq.Code)

	listReq := performMissionRequest(t, handler, http.MethodGet, "/missions", nil, actorID, "")
	require.Equal(t, http.StatusOK, listReq.Code)
	var listRes struct {
		Count int `json:"count"`
		Items []struct {
			Status string `json:"status"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(listReq.Body.Bytes(), &listRes))
	require.Equal(t, 1, listRes.Count)
	require.Equal(t, "claimed", listRes.Items[0].Status)

	disableIngestReq := performMissionJSONRequest(t, handler, http.MethodPost, "/missions/admin/progress-ingest-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableIngestReq.Code)

	ingestBlockedReq := performMissionJSONRequest(t, handler, http.MethodPost, "/missions/daily_read_3/progress/ingest", map[string]any{
		"delta":       1,
		"source_type": "history",
		"request_id":  "req-stage16-progress-3",
	}, actorID, "")
	require.Equal(t, http.StatusForbidden, ingestBlockedReq.Code)

	resetReq := performMissionJSONRequest(t, handler, http.MethodPost, "/missions/admin/reset-progress", map[string]any{
		"target_user_id": actorID,
		"mission_id":     "daily_read_3",
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, resetReq.Code)

	disableReadReq := performMissionJSONRequest(t, handler, http.MethodPost, "/missions/admin/read-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableReadReq.Code)

	listBlockedReq := performMissionRequest(t, handler, http.MethodGet, "/missions", nil, actorID, "")
	require.Equal(t, http.StatusNotFound, listBlockedReq.Code)
}

func performMissionJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performMissionRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performMissionRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
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

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
	inventorymodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInventoryHTTPFlowAndRuntimeControls(t *testing.T) {
	actorID := uuid.NewString()
	adminID := uuid.NewString()

	inventory := inventorymodule.New()
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		inventory,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.15.0-test"}, zap.NewNop(), nil, registry)

	upsertDefReq := performInventoryJSONRequest(t, handler, http.MethodPost, "/inventory/admin/items", map[string]any{
		"item_id":    "mana_potion",
		"item_type":  "consumable",
		"stackable":  true,
		"equipable":  false,
		"consumable": true,
		"max_stack":  99,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertDefReq.Code)

	claimReq := performInventoryJSONRequest(t, handler, http.MethodPost, "/inventory/claim", map[string]any{
		"item_id":     "mana_potion",
		"quantity":    4,
		"source_type": "mission",
		"source_ref":  "daily_quest_1",
		"request_id":  "req-stage15-claim-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, claimReq.Code)

	listReq := performInventoryRequest(t, handler, http.MethodGet, "/inventory/items", nil, actorID, "")
	require.Equal(t, http.StatusOK, listReq.Code)
	var listRes struct {
		Count int `json:"count"`
		Items []struct {
			Quantity int `json:"quantity"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(listReq.Body.Bytes(), &listRes))
	require.Equal(t, 1, listRes.Count)
	require.Equal(t, 4, listRes.Items[0].Quantity)

	consumeReq := performInventoryJSONRequest(t, handler, http.MethodPost, "/inventory/items/mana_potion/consume", map[string]any{
		"quantity": 2,
	}, actorID, "")
	require.Equal(t, http.StatusOK, consumeReq.Code)

	disableClaimReq := performInventoryJSONRequest(t, handler, http.MethodPost, "/inventory/admin/claim-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableClaimReq.Code)

	claimBlockedReq := performInventoryJSONRequest(t, handler, http.MethodPost, "/inventory/claim", map[string]any{
		"item_id":     "mana_potion",
		"quantity":    1,
		"source_type": "mission",
		"source_ref":  "daily_quest_2",
		"request_id":  "req-stage15-claim-2",
	}, actorID, "")
	require.Equal(t, http.StatusForbidden, claimBlockedReq.Code)
}

func performInventoryJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performInventoryRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performInventoryRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
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

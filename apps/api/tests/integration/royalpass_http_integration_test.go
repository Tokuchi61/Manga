package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	accessmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/access"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	rpmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRoyalPassHTTPSeasonProgressClaimPremiumAndRuntimeFlow(t *testing.T) {
	actorID := uuid.NewString()
	adminID := uuid.NewString()

	rp := rpmodule.New()
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		rp,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.17.0-test"}, zap.NewNop(), nil, registry)

	startsAt := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
	endsAt := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)

	upsertSeasonReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/admin/seasons", map[string]any{
		"season_id": "season_2026_03",
		"title":     "Spring Launch",
		"state":     "active",
		"starts_at": startsAt,
		"ends_at":   endsAt,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertSeasonReq.Code)

	upsertFreeTierReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/admin/tiers", map[string]any{
		"season_id":       "season_2026_03",
		"tier_number":     1,
		"track":           "free",
		"required_points": 50,
		"reward_item_id":  "mana_100",
		"reward_quantity": 1,
		"active":          true,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertFreeTierReq.Code)

	upsertPremiumTierReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/admin/tiers", map[string]any{
		"season_id":       "season_2026_03",
		"tier_number":     1,
		"track":           "premium",
		"required_points": 50,
		"reward_item_id":  "avatar_frame_gold",
		"reward_quantity": 1,
		"active":          true,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertPremiumTierReq.Code)

	ingestReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/progress/ingest", map[string]any{
		"season_id":   "season_2026_03",
		"delta":       60,
		"source_type": "mission",
		"request_id":  "req-stage17-progress-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, ingestReq.Code)

	claimFreeReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/claims", map[string]any{
		"season_id":   "season_2026_03",
		"tier_number": 1,
		"track":       "free",
		"request_id":  "req-stage17-claim-free-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, claimFreeReq.Code)

	claimPremiumBlockedReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/claims", map[string]any{
		"season_id":   "season_2026_03",
		"tier_number": 1,
		"track":       "premium",
		"request_id":  "req-stage17-claim-premium-1",
	}, actorID, "")
	require.Equal(t, http.StatusForbidden, claimPremiumBlockedReq.Code)

	activatePremiumReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/premium/activate", map[string]any{
		"season_id":      "season_2026_03",
		"source_type":    "shop",
		"activation_ref": "shop-activation-1",
		"request_id":     "req-stage17-activate-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, activatePremiumReq.Code)

	claimPremiumReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/claims", map[string]any{
		"season_id":   "season_2026_03",
		"tier_number": 1,
		"track":       "premium",
		"request_id":  "req-stage17-claim-premium-2",
	}, actorID, "")
	require.Equal(t, http.StatusOK, claimPremiumReq.Code)

	overviewReq := performRoyalPassRequest(t, handler, http.MethodGet, "/royalpass/overview?season_id=season_2026_03", nil, actorID, "")
	require.Equal(t, http.StatusOK, overviewReq.Code)
	var overviewRes struct {
		Points           int  `json:"points"`
		PremiumActivated bool `json:"premium_activated"`
		Count            int  `json:"count"`
	}
	require.NoError(t, json.Unmarshal(overviewReq.Body.Bytes(), &overviewRes))
	require.Equal(t, 60, overviewRes.Points)
	require.True(t, overviewRes.PremiumActivated)
	require.Equal(t, 2, overviewRes.Count)

	disableClaimReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/admin/claim-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableClaimReq.Code)

	claimWhenDisabledReq := performRoyalPassJSONRequest(t, handler, http.MethodPost, "/royalpass/claims", map[string]any{
		"season_id":   "season_2026_03",
		"tier_number": 1,
		"track":       "free",
		"request_id":  "req-stage17-claim-disabled-1",
	}, actorID, "")
	require.Equal(t, http.StatusForbidden, claimWhenDisabledReq.Code)
}

func performRoyalPassJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performRoyalPassRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performRoyalPassRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
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

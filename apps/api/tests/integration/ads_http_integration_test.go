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
	adsmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAdsHTTPResolveIntakeRuntimeAndAggregateFlow(t *testing.T) {
	actorID := uuid.NewString()
	adminID := uuid.NewString()

	ads := adsmodule.New()
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		ads,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.20.0-test"}, zap.NewNop(), nil, registry)

	upsertPlacementReq := performAdsJSONRequest(t, handler, http.MethodPost, "/ads/admin/placements", map[string]any{
		"placement_id":  "home_top",
		"surface":       "home",
		"target_type":   "none",
		"visible":       true,
		"priority":      100,
		"frequency_cap": 2,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertPlacementReq.Code)

	now := time.Now().UTC()
	startsAt := now.Add(-time.Hour).Format(time.RFC3339)
	endsAt := now.Add(24 * time.Hour).Format(time.RFC3339)

	upsertCampaignReq := performAdsJSONRequest(t, handler, http.MethodPost, "/ads/admin/campaigns", map[string]any{
		"campaign_id":  "campaign_home_1",
		"placement_id": "home_top",
		"name":         "Home Top Campaign",
		"state":        "active",
		"creative_url": "https://cdn.example.com/ads/home-top.png",
		"click_url":    "https://example.com/c/home-top",
		"weight":       50,
		"starts_at":    startsAt,
		"ends_at":      endsAt,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertCampaignReq.Code)

	resolveReq := performAdsRequest(t, handler, http.MethodGet, "/ads/resolve?surface=home&target_type=none&session_id=session-1", nil, "", "")
	require.Equal(t, http.StatusOK, resolveReq.Code)
	var resolveRes struct {
		Count int `json:"count"`
		Items []struct {
			PlacementID string `json:"placement_id"`
			Campaigns   []struct {
				CampaignID string `json:"campaign_id"`
			} `json:"campaigns"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(resolveReq.Body.Bytes(), &resolveRes))
	require.Equal(t, 1, resolveRes.Count)
	require.Equal(t, "home_top", resolveRes.Items[0].PlacementID)
	require.Equal(t, "campaign_home_1", resolveRes.Items[0].Campaigns[0].CampaignID)

	impressionReq := performAdsJSONRequest(t, handler, http.MethodPost, "/ads/impressions", map[string]any{
		"request_id":   "req-ads-imp-1",
		"placement_id": "home_top",
		"campaign_id":  "campaign_home_1",
		"session_id":   "session-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, impressionReq.Code)
	var impressionRes struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(impressionReq.Body.Bytes(), &impressionRes))
	require.Equal(t, "accepted", impressionRes.Status)

	impressionIdempotentReq := performAdsJSONRequest(t, handler, http.MethodPost, "/ads/impressions", map[string]any{
		"request_id":   "req-ads-imp-1",
		"placement_id": "home_top",
		"campaign_id":  "campaign_home_1",
		"session_id":   "session-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, impressionIdempotentReq.Code)
	var impressionIdempotentRes struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(impressionIdempotentReq.Body.Bytes(), &impressionIdempotentRes))
	require.Equal(t, "idempotent", impressionIdempotentRes.Status)

	clickReq := performAdsJSONRequest(t, handler, http.MethodPost, "/ads/clicks", map[string]any{
		"request_id":   "req-ads-click-1",
		"placement_id": "home_top",
		"campaign_id":  "campaign_home_1",
		"session_id":   "session-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, clickReq.Code)
	var clickRes struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(clickReq.Body.Bytes(), &clickRes))
	require.Equal(t, "accepted", clickRes.Status)

	ignoredClickReq := performAdsJSONRequest(t, handler, http.MethodPost, "/ads/clicks", map[string]any{
		"request_id":      "req-ads-click-2",
		"placement_id":    "home_top",
		"campaign_id":     "campaign_home_1",
		"session_id":      "session-1",
		"invalid_traffic": true,
	}, actorID, "")
	require.Equal(t, http.StatusOK, ignoredClickReq.Code)
	var ignoredClickRes struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(ignoredClickReq.Body.Bytes(), &ignoredClickRes))
	require.Equal(t, "ignored_invalid_traffic", ignoredClickRes.Status)

	aggregateReq := performAdsRequest(t, handler, http.MethodGet, "/ads/admin/aggregate", nil, adminID, "admin")
	require.Equal(t, http.StatusOK, aggregateReq.Code)
	var aggregateRes struct {
		Count int `json:"count"`
		Items []struct {
			CampaignID      string `json:"campaign_id"`
			ImpressionCount int    `json:"impression_count"`
			ClickCount      int    `json:"click_count"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(aggregateReq.Body.Bytes(), &aggregateRes))
	require.Equal(t, 1, aggregateRes.Count)
	require.Equal(t, "campaign_home_1", aggregateRes.Items[0].CampaignID)
	require.Equal(t, 1, aggregateRes.Items[0].ImpressionCount)
	require.Equal(t, 1, aggregateRes.Items[0].ClickCount)

	disableClickIntakeReq := performAdsJSONRequest(t, handler, http.MethodPost, "/ads/admin/click-intake-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableClickIntakeReq.Code)

	clickWhenDisabledReq := performAdsJSONRequest(t, handler, http.MethodPost, "/ads/clicks", map[string]any{
		"request_id":   "req-ads-click-disabled",
		"placement_id": "home_top",
		"campaign_id":  "campaign_home_1",
		"session_id":   "session-1",
	}, actorID, "")
	require.Equal(t, http.StatusServiceUnavailable, clickWhenDisabledReq.Code)

	disableCampaignReq := performAdsJSONRequest(t, handler, http.MethodPost, "/ads/admin/campaign-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableCampaignReq.Code)

	resolveWhenDisabledReq := performAdsRequest(t, handler, http.MethodGet, "/ads/resolve?surface=home&target_type=none&session_id=session-1", nil, "", "")
	require.Equal(t, http.StatusNotFound, resolveWhenDisabledReq.Code)
}

func performAdsJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performAdsRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performAdsRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	bodyReader := bytes.NewReader(nil)
	if body != nil {
		bodyReader = body
	}

	req := httptest.NewRequest(method, path, bodyReader)
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

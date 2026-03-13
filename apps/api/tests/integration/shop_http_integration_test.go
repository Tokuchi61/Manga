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
	shopmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/shop"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestShopHTTPCatalogPurchaseRecoveryAndRuntimeFlow(t *testing.T) {
	actorID := uuid.NewString()
	adminID := uuid.NewString()

	shop := shopmodule.New()
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		shop,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.18.0-test"}, zap.NewNop(), nil, registry)

	upsertProductReq := performShopJSONRequest(t, handler, http.MethodPost, "/shop/admin/products", map[string]any{
		"product_id":        "product_avatar_01",
		"name":              "Avatar Frame",
		"category":          "cosmetic",
		"state":             "active",
		"inventory_item_id": "avatar_frame_gold",
		"single_purchase":   true,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertProductReq.Code)

	startsAt := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
	endsAt := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)

	upsertOfferReq := performShopJSONRequest(t, handler, http.MethodPost, "/shop/admin/offers", map[string]any{
		"offer_id":         "offer_avatar_01",
		"product_id":       "product_avatar_01",
		"title":            "Launch Offer",
		"visibility":       "visible",
		"price_mana":       500,
		"discount_percent": 10,
		"active":           true,
		"starts_at":        startsAt,
		"ends_at":          endsAt,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertOfferReq.Code)

	catalogReq := performShopRequest(t, handler, http.MethodGet, "/shop/catalog", nil, actorID, "")
	require.Equal(t, http.StatusOK, catalogReq.Code)
	var catalogRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(catalogReq.Body.Bytes(), &catalogRes))
	require.Equal(t, 1, catalogRes.Count)

	purchaseReq := performShopJSONRequest(t, handler, http.MethodPost, "/shop/purchase/intents", map[string]any{
		"product_id": "product_avatar_01",
		"request_id": "req-shop-stage18-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, purchaseReq.Code)
	var purchaseRes struct {
		Status   string `json:"status"`
		IntentID string `json:"intent_id"`
	}
	require.NoError(t, json.Unmarshal(purchaseReq.Body.Bytes(), &purchaseRes))
	require.Equal(t, "intent_created", purchaseRes.Status)
	require.NotEmpty(t, purchaseRes.IntentID)

	purchaseIdempotentReq := performShopJSONRequest(t, handler, http.MethodPost, "/shop/purchase/intents", map[string]any{
		"product_id": "product_avatar_01",
		"request_id": "req-shop-stage18-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, purchaseIdempotentReq.Code)

	purchaseDuplicateReq := performShopJSONRequest(t, handler, http.MethodPost, "/shop/purchase/intents", map[string]any{
		"product_id": "product_avatar_01",
		"request_id": "req-shop-stage18-2",
	}, actorID, "")
	require.Equal(t, http.StatusConflict, purchaseDuplicateReq.Code)

	recoveryReq := performShopJSONRequest(t, handler, http.MethodPost, "/shop/purchase/recovery", map[string]any{
		"intent_id": purchaseRes.IntentID,
	}, actorID, "")
	require.Equal(t, http.StatusOK, recoveryReq.Code)

	disablePurchaseReq := performShopJSONRequest(t, handler, http.MethodPost, "/shop/admin/purchase-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disablePurchaseReq.Code)

	purchaseWhenDisabledReq := performShopJSONRequest(t, handler, http.MethodPost, "/shop/purchase/intents", map[string]any{
		"product_id": "product_avatar_01",
		"request_id": "req-shop-stage18-disabled",
	}, actorID, "")
	require.Equal(t, http.StatusForbidden, purchaseWhenDisabledReq.Code)
}

func performShopJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performShopRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performShopRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
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

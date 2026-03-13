package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	accessmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/access"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	paymentmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment"
	paymenthandler "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/handler"
	paymentservice "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/service"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPaymentHTTPCheckoutCallbackWalletAndRuntimeFlow(t *testing.T) {
	actorID := uuid.NewString()
	adminID := uuid.NewString()

	payment := paymentmodule.New(paymentmodule.RuntimeConfig{})
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		payment,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.19.0-test"}, zap.NewNop(), nil, registry)

	upsertPackageReq := performPaymentJSONRequest(t, handler, http.MethodPost, "/payment/admin/packages", map[string]any{
		"package_id":     "mana_pack_small",
		"name":           "Small Mana Pack",
		"description":    "500 mana",
		"mana_amount":    500,
		"price_amount":   499,
		"price_currency": "TRY",
		"active":         true,
		"display_order":  10,
		"provider":       "mock_provider",
		"provider_sku":   "sku-small",
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, upsertPackageReq.Code)

	listPackagesReq := performPaymentRequest(t, handler, http.MethodGet, "/payment/packages", nil, actorID, "", nil)
	require.Equal(t, http.StatusOK, listPackagesReq.Code)
	var listPackagesRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(listPackagesReq.Body.Bytes(), &listPackagesRes))
	require.Equal(t, 1, listPackagesRes.Count)

	checkoutReq := performPaymentJSONRequest(t, handler, http.MethodPost, "/payment/checkout/sessions", map[string]any{
		"package_id": "mana_pack_small",
		"request_id": "req-payment-http-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, checkoutReq.Code)
	var checkoutRes struct {
		Status        string `json:"status"`
		TransactionID string `json:"transaction_id"`
		SessionID     string `json:"session_id"`
	}
	require.NoError(t, json.Unmarshal(checkoutReq.Body.Bytes(), &checkoutRes))
	require.Equal(t, "checkout_started", checkoutRes.Status)
	require.NotEmpty(t, checkoutRes.TransactionID)
	require.NotEmpty(t, checkoutRes.SessionID)

	callbackReq := performPaymentJSONRequest(t, handler, http.MethodPost, "/payment/callback", map[string]any{
		"provider_event_id":  "evt-payment-http-1",
		"session_id":         checkoutRes.SessionID,
		"provider_reference": "provider-ref-http-1",
		"status":             "success",
	}, "", "")
	require.Equal(t, http.StatusOK, callbackReq.Code)

	callbackIdempotentReq := performPaymentJSONRequest(t, handler, http.MethodPost, "/payment/callback", map[string]any{
		"provider_event_id": "evt-payment-http-1",
		"session_id":        checkoutRes.SessionID,
		"status":            "success",
	}, "", "")
	require.Equal(t, http.StatusOK, callbackIdempotentReq.Code)

	walletReq := performPaymentRequest(t, handler, http.MethodGet, "/payment/wallet", nil, actorID, "", nil)
	require.Equal(t, http.StatusOK, walletReq.Code)
	var walletRes struct {
		BalanceMana int `json:"balance_mana"`
	}
	require.NoError(t, json.Unmarshal(walletReq.Body.Bytes(), &walletRes))
	require.Equal(t, 500, walletRes.BalanceMana)

	transactionsReq := performPaymentRequest(t, handler, http.MethodGet, "/payment/transactions", nil, actorID, "", nil)
	require.Equal(t, http.StatusOK, transactionsReq.Code)
	var transactionsRes struct {
		Count int `json:"count"`
		Items []struct {
			Status string `json:"status"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(transactionsReq.Body.Bytes(), &transactionsRes))
	require.Equal(t, 1, transactionsRes.Count)
	require.Equal(t, "success", transactionsRes.Items[0].Status)

	disableReadReq := performPaymentJSONRequest(t, handler, http.MethodPost, "/payment/admin/transaction-read-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableReadReq.Code)

	walletWhenDisabledReq := performPaymentRequest(t, handler, http.MethodGet, "/payment/wallet", nil, actorID, "", nil)
	require.Equal(t, http.StatusNotFound, walletWhenDisabledReq.Code)

	enableReadReq := performPaymentJSONRequest(t, handler, http.MethodPost, "/payment/admin/transaction-read-state", map[string]any{
		"enabled": true,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, enableReadReq.Code)

	refundReq := performPaymentJSONRequest(t, handler, http.MethodPost, "/payment/admin/refunds", map[string]any{
		"transaction_id": checkoutRes.TransactionID,
		"reason_code":    "manual_refund",
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, refundReq.Code)

	walletAfterRefundReq := performPaymentRequest(t, handler, http.MethodGet, "/payment/wallet", nil, actorID, "", nil)
	require.Equal(t, http.StatusOK, walletAfterRefundReq.Code)
	var walletAfterRefundRes struct {
		BalanceMana int `json:"balance_mana"`
	}
	require.NoError(t, json.Unmarshal(walletAfterRefundReq.Body.Bytes(), &walletAfterRefundRes))
	require.Equal(t, 0, walletAfterRefundRes.BalanceMana)

	disableCheckoutReq := performPaymentJSONRequest(t, handler, http.MethodPost, "/payment/admin/checkout-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableCheckoutReq.Code)

	checkoutDisabledReq := performPaymentJSONRequest(t, handler, http.MethodPost, "/payment/checkout/sessions", map[string]any{
		"package_id": "mana_pack_small",
		"request_id": "req-payment-http-disabled",
	}, actorID, "")
	require.Equal(t, http.StatusServiceUnavailable, checkoutDisabledReq.Code)
}

func performPaymentJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	extraHeaders := map[string]string{}
	if path == "/payment/callback" {
		timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
		signature := paymentservice.SignProviderCallback(paymentservice.DefaultNonProdCallbackSigningSecret(), timestamp, payload)
		extraHeaders[paymenthandler.HeaderPaymentCallbackSignature] = signature
		extraHeaders[paymenthandler.HeaderPaymentCallbackTimestamp] = timestamp
	}

	return performPaymentRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles, extraHeaders)
}

func performPaymentRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string, extraHeaders map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	bodyReader := bytes.NewReader(nil)
	if body != nil {
		bodyReader = body
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range extraHeaders {
		req.Header.Set(key, value)
	}
	if actorUserID != "" || roles != "" {
		setActorHeaders(req, actorUserID, "", roles)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

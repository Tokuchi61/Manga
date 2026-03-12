package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAuthHTTPRegisterVerifyAndLoginFlow(t *testing.T) {
	registry, err := modules.NewRegistry(authmodule.New(authmodule.RuntimeConfig{}))
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.4.0-test"}, zap.NewNop(), nil, registry)

	registerBody := map[string]string{
		"email":    "http-auth@example.com",
		"password": "StrongPass123!",
	}
	registerRec := performJSONRequest(t, handler, http.MethodPost, "/auth/register", registerBody)
	require.Equal(t, http.StatusCreated, registerRec.Code)

	var registerRes struct {
		CredentialID      string `json:"credential_id"`
		VerificationToken string `json:"verification_token"`
	}
	require.NoError(t, json.Unmarshal(registerRec.Body.Bytes(), &registerRes))
	require.NotEmpty(t, registerRes.CredentialID)
	require.NotEmpty(t, registerRes.VerificationToken)

	confirmBody := map[string]string{"token": registerRes.VerificationToken}
	confirmRec := performJSONRequest(t, handler, http.MethodPost, "/auth/email/verification/confirm", confirmBody)
	require.Equal(t, http.StatusOK, confirmRec.Code)

	loginBody := map[string]string{
		"email":    "http-auth@example.com",
		"password": "StrongPass123!",
		"device":   "android",
		"ip":       "10.10.1.1",
	}
	loginRec := performJSONRequest(t, handler, http.MethodPost, "/auth/login", loginBody)
	require.Equal(t, http.StatusOK, loginRec.Code)

	var loginRes struct {
		CredentialID string `json:"credential_id"`
		SessionID    string `json:"session_id"`
		RefreshToken string `json:"refresh_token"`
	}
	require.NoError(t, json.Unmarshal(loginRec.Body.Bytes(), &loginRes))
	require.NotEmpty(t, loginRes.SessionID)
	require.NotEmpty(t, loginRes.RefreshToken)

	sessionsReq, err := http.NewRequest(http.MethodGet, "/auth/sessions?credential_id="+loginRes.CredentialID, nil)
	require.NoError(t, err)
	sessionsRec := httptest.NewRecorder()
	handler.ServeHTTP(sessionsRec, sessionsReq)
	require.Equal(t, http.StatusOK, sessionsRec.Code)
}

func TestAuthHTTPLoginCooldown(t *testing.T) {
	registry, err := modules.NewRegistry(authmodule.New(authmodule.RuntimeConfig{
		FailedAttemptLimitPerMinute:       2,
		LoginCooldownSeconds:              300,
		VerificationResendCooldownSeconds: 0,
	}))
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.4.0-test"}, zap.NewNop(), nil, registry)

	registerRec := performJSONRequest(t, handler, http.MethodPost, "/auth/register", map[string]string{
		"email":    "cooldown-http@example.com",
		"password": "StrongPass123!",
	})
	require.Equal(t, http.StatusCreated, registerRec.Code)

	var registerRes struct {
		VerificationToken string `json:"verification_token"`
	}
	require.NoError(t, json.Unmarshal(registerRec.Body.Bytes(), &registerRes))
	confirmRec := performJSONRequest(t, handler, http.MethodPost, "/auth/email/verification/confirm", map[string]string{"token": registerRes.VerificationToken})
	require.Equal(t, http.StatusOK, confirmRec.Code)

	failRec1 := performJSONRequest(t, handler, http.MethodPost, "/auth/login", map[string]string{
		"email":    "cooldown-http@example.com",
		"password": "WrongPass123!",
	})
	require.Equal(t, http.StatusUnauthorized, failRec1.Code)

	failRec2 := performJSONRequest(t, handler, http.MethodPost, "/auth/login", map[string]string{
		"email":    "cooldown-http@example.com",
		"password": "WrongPass123!",
	})
	require.Equal(t, http.StatusTooManyRequests, failRec2.Code)
}

func performJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	return rec
}

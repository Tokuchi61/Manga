package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestUserHTTPProfileAndVisibilityFlow(t *testing.T) {
	registry, err := modules.NewRegistry(authmodule.New(authmodule.RuntimeConfig{}), usermodule.New())
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.5.0-test"}, zap.NewNop(), nil, registry)

	createRec := performUserJSONRequest(t, handler, http.MethodPost, "/users", map[string]string{
		"credential_id": uuid.NewString(),
		"username":      "reader_http_one",
		"display_name":  "Reader HTTP One",
	})
	require.Equal(t, http.StatusCreated, createRec.Code)

	var createRes struct {
		UserID string `json:"user_id"`
	}
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &createRes))
	require.NotEmpty(t, createRes.UserID)

	publicRec := performUserRequest(t, handler, http.MethodGet, "/users/"+createRes.UserID, nil)
	require.Equal(t, http.StatusOK, publicRec.Code)

	visibilityRec := performUserJSONRequest(t, handler, http.MethodPatch, "/users/"+createRes.UserID+"/visibility?viewer_id="+createRes.UserID, map[string]string{
		"profile_visibility": "private",
	})
	require.Equal(t, http.StatusOK, visibilityRec.Code)

	publicAfterRec := performUserRequest(t, handler, http.MethodGet, "/users/"+createRes.UserID, nil)
	require.Equal(t, http.StatusForbidden, publicAfterRec.Code)

	ownRec := performUserRequest(t, handler, http.MethodGet, "/users/"+createRes.UserID+"/self?viewer_id="+createRes.UserID, nil)
	require.Equal(t, http.StatusOK, ownRec.Code)

	updateProfileRec := performUserJSONRequest(t, handler, http.MethodPatch, "/users/"+createRes.UserID+"/profile?viewer_id="+createRes.UserID, map[string]string{
		"display_name": "Reader Updated",
		"bio":          "Profile updated over HTTP",
	})
	require.Equal(t, http.StatusOK, updateProfileRec.Code)

	ownAfterUpdateRec := performUserRequest(t, handler, http.MethodGet, "/users/"+createRes.UserID+"/self?viewer_id="+createRes.UserID, nil)
	require.Equal(t, http.StatusOK, ownAfterUpdateRec.Code)

	var ownAfterUpdateRes struct {
		DisplayName string `json:"display_name"`
		Bio         string `json:"bio"`
	}
	require.NoError(t, json.Unmarshal(ownAfterUpdateRec.Body.Bytes(), &ownAfterUpdateRes))
	require.Equal(t, "Reader Updated", ownAfterUpdateRes.DisplayName)
	require.Equal(t, "Profile updated over HTTP", ownAfterUpdateRes.Bio)
}

func TestUserHTTPAccountStateAndVIPLifecycle(t *testing.T) {
	registry, err := modules.NewRegistry(authmodule.New(authmodule.RuntimeConfig{}), usermodule.New())
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.5.0-test"}, zap.NewNop(), nil, registry)

	createRec := performUserJSONRequest(t, handler, http.MethodPost, "/users", map[string]string{
		"credential_id": uuid.NewString(),
		"username":      "reader_http_two",
	})
	require.Equal(t, http.StatusCreated, createRec.Code)

	var createRes struct {
		UserID string `json:"user_id"`
	}
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &createRes))
	require.NotEmpty(t, createRes.UserID)

	vipActivateRec := performUserJSONRequest(t, handler, http.MethodPatch, "/users/"+createRes.UserID+"/vip", map[string]string{
		"action":  "activate",
		"ends_at": time.Date(2026, 4, 11, 10, 0, 0, 0, time.UTC).Format(time.RFC3339),
	})
	require.Equal(t, http.StatusOK, vipActivateRec.Code)

	vipFreezeRec := performUserJSONRequest(t, handler, http.MethodPatch, "/users/"+createRes.UserID+"/vip", map[string]string{
		"action":        "freeze",
		"freeze_reason": "system_pause",
	})
	require.Equal(t, http.StatusOK, vipFreezeRec.Code)

	vipResumeRec := performUserJSONRequest(t, handler, http.MethodPatch, "/users/"+createRes.UserID+"/vip", map[string]string{
		"action": "resume",
	})
	require.Equal(t, http.StatusOK, vipResumeRec.Code)

	deactivateRec := performUserJSONRequest(t, handler, http.MethodPatch, "/users/"+createRes.UserID+"/account/state", map[string]string{
		"actor_scope":   "self",
		"actor_user_id": createRes.UserID,
		"account_state": "deactivated",
	})
	require.Equal(t, http.StatusOK, deactivateRec.Code)

	updateProfileRec := performUserJSONRequest(t, handler, http.MethodPatch, "/users/"+createRes.UserID+"/profile?viewer_id="+createRes.UserID, map[string]string{
		"display_name": "Should Not Update",
	})
	require.Equal(t, http.StatusForbidden, updateProfileRec.Code)

	banRec := performUserJSONRequest(t, handler, http.MethodPatch, "/users/"+createRes.UserID+"/account/state", map[string]string{
		"actor_scope":   "admin",
		"account_state": "banned",
		"reason":        "high_risk",
	})
	require.Equal(t, http.StatusOK, banRec.Code)

	publicRec := performUserRequest(t, handler, http.MethodGet, "/users/"+createRes.UserID, nil)
	require.Equal(t, http.StatusForbidden, publicRec.Code)
}

func performUserJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performUserRequest(t, handler, method, path, bytes.NewReader(payload))
}

func performUserRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader) *httptest.ResponseRecorder {
	t.Helper()

	var bodyReader *bytes.Reader
	if body == nil {
		bodyReader = bytes.NewReader(nil)
	} else {
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

func TestUserRouteExists(t *testing.T) {
	registry, err := modules.NewRegistry(authmodule.New(authmodule.RuntimeConfig{}), usermodule.New())
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.5.0-test"}, zap.NewNop(), nil, registry)
	rec := performUserJSONRequest(t, handler, http.MethodPost, "/users", map[string]string{"invalid": "true"})
	require.NotEqual(t, http.StatusNotFound, rec.Code, fmt.Sprintf("response body: %s", rec.Body.String()))
}

package e2e_test

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
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCriticalAuthAccessUserFlow(t *testing.T) {
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.10.0-e2e"}, zap.NewNop(), nil, registry)

	registerRec := performJSONRequest(t, handler, http.MethodPost, "/auth/register", map[string]string{
		"email":    "e2e-reader@example.com",
		"password": "StrongPass123!",
	}, "", "", "")
	require.Equal(t, http.StatusCreated, registerRec.Code)

	var registerRes struct {
		CredentialID      string `json:"credential_id"`
		VerificationToken string `json:"verification_token"`
	}
	require.NoError(t, json.Unmarshal(registerRec.Body.Bytes(), &registerRes))
	require.NotEmpty(t, registerRes.CredentialID)
	require.NotEmpty(t, registerRes.VerificationToken)

	confirmRec := performJSONRequest(t, handler, http.MethodPost, "/auth/email/verification/confirm", map[string]string{
		"token": registerRes.VerificationToken,
	}, "", "", "")
	require.Equal(t, http.StatusOK, confirmRec.Code)

	createUserRec := performJSONRequest(t, handler, http.MethodPost, "/users", map[string]string{
		"credential_id": registerRes.CredentialID,
		"username":      "e2e_reader",
		"display_name":  "E2E Reader",
	}, "", "", "")
	require.Equal(t, http.StatusCreated, createUserRec.Code)

	var createUserRes struct {
		UserID string `json:"user_id"`
	}
	require.NoError(t, json.Unmarshal(createUserRec.Body.Bytes(), &createUserRes))
	require.NotEmpty(t, createUserRes.UserID)

	selfRec := performJSONRequest(t, handler, http.MethodGet, "/users/"+createUserRes.UserID+"/self", nil, createUserRes.UserID, "", "")
	require.Equal(t, http.StatusOK, selfRec.Code)

	adminActorID := uuid.NewString()
	adminCredentialID := uuid.NewString()
	permissionName := "user.profile.read.any"

	createRoleRec := performJSONRequest(t, handler, http.MethodPost, "/access/roles", map[string]any{
		"name":     "e2e_profile_reader",
		"priority": 42,
	}, adminActorID, adminCredentialID, "admin")
	require.Equal(t, http.StatusCreated, createRoleRec.Code)

	var roleRes struct {
		RoleID string `json:"role_id"`
	}
	require.NoError(t, json.Unmarshal(createRoleRec.Body.Bytes(), &roleRes))
	require.NotEmpty(t, roleRes.RoleID)

	createPermissionRec := performJSONRequest(t, handler, http.MethodPost, "/access/permissions", map[string]any{
		"name":          permissionName,
		"module":        "user",
		"surface":       "profile",
		"action":        "read",
		"audience_kind": "all",
	}, adminActorID, adminCredentialID, "admin")
	require.Equal(t, http.StatusCreated, createPermissionRec.Code)

	attachPermissionRec := performJSONRequest(t, handler, http.MethodPost, "/access/roles/"+roleRes.RoleID+"/permissions", map[string]any{
		"permission_name": permissionName,
	}, adminActorID, adminCredentialID, "admin")
	require.Equal(t, http.StatusOK, attachPermissionRec.Code)

	assignRoleRec := performJSONRequest(t, handler, http.MethodPost, "/access/users/"+createUserRes.UserID+"/roles", map[string]any{
		"role_name": "e2e_profile_reader",
	}, adminActorID, adminCredentialID, "admin")
	require.Equal(t, http.StatusOK, assignRoleRec.Code)

	evaluateRec := performJSONRequest(t, handler, http.MethodPost, "/access/evaluate", map[string]any{
		"permission": permissionName,
	}, createUserRes.UserID, registerRes.CredentialID, "admin")
	require.Equal(t, http.StatusOK, evaluateRec.Code)

	var evaluateRes struct {
		Allowed bool `json:"allowed"`
	}
	require.NoError(t, json.Unmarshal(evaluateRec.Body.Bytes(), &evaluateRes))
	require.True(t, evaluateRes.Allowed)
}

func performJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, actorCredentialID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload := bytes.NewReader(nil)
	if body != nil {
		serialized, err := json.Marshal(body)
		require.NoError(t, err)
		payload = bytes.NewReader(serialized)
	}

	req := httptest.NewRequest(method, path, payload)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if actorUserID != "" {
		req.Header.Set(identity.HeaderActorUserID, actorUserID)
	}
	if actorCredentialID != "" {
		req.Header.Set(identity.HeaderActorCredentialID, actorCredentialID)
	}
	if roles != "" {
		req.Header.Set(identity.HeaderActorRoles, roles)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

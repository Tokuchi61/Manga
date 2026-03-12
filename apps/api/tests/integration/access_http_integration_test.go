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
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAccessHTTPRolePermissionAndEvaluationFlow(t *testing.T) {
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.6.0-test"}, zap.NewNop(), nil, registry)

	roleRec := performAccessJSONRequest(t, handler, http.MethodPost, "/access/roles", map[string]any{
		"name":     "content_moderator",
		"priority": 45,
	})
	require.Equal(t, http.StatusCreated, roleRec.Code)

	var roleRes struct {
		RoleID string `json:"role_id"`
	}
	require.NoError(t, json.Unmarshal(roleRec.Body.Bytes(), &roleRes))
	require.NotEmpty(t, roleRes.RoleID)

	permissionName := "comment.manage.any"
	permissionRec := performAccessJSONRequest(t, handler, http.MethodPost, "/access/permissions", map[string]any{
		"name":          permissionName,
		"module":        "comment",
		"surface":       "manage",
		"action":        "write",
		"audience_kind": "authenticated",
	})
	require.Equal(t, http.StatusCreated, permissionRec.Code)

	assignPermissionRec := performAccessJSONRequest(t, handler, http.MethodPost, "/access/roles/"+roleRes.RoleID+"/permissions", map[string]any{
		"permission_name": permissionName,
	})
	require.Equal(t, http.StatusOK, assignPermissionRec.Code)

	userID := uuid.NewString()
	assignRoleRec := performAccessJSONRequest(t, handler, http.MethodPost, "/access/users/"+userID+"/roles", map[string]any{
		"role_name": "content_moderator",
	})
	require.Equal(t, http.StatusOK, assignRoleRec.Code)

	evaluateAllowedRec := performAccessJSONRequest(t, handler, http.MethodPost, "/access/evaluate", map[string]any{
		"user_id":    userID,
		"permission": permissionName,
		"identity": map[string]any{
			"credential_id":  uuid.NewString(),
			"email_verified": true,
		},
		"user_signal": map[string]any{
			"account_state": "active",
		},
	})
	require.Equal(t, http.StatusOK, evaluateAllowedRec.Code)

	var evaluateAllowedRes struct {
		Allowed bool `json:"allowed"`
	}
	require.NoError(t, json.Unmarshal(evaluateAllowedRec.Body.Bytes(), &evaluateAllowedRes))
	require.True(t, evaluateAllowedRes.Allowed)

	evaluateGuestDenyRec := performAccessJSONRequest(t, handler, http.MethodPost, "/access/evaluate", map[string]any{
		"permission": "chapter.read.authenticated",
	})
	require.Equal(t, http.StatusOK, evaluateGuestDenyRec.Code)

	var evaluateGuestDenyRes struct {
		Allowed    bool   `json:"allowed"`
		ReasonCode string `json:"reason_code"`
	}
	require.NoError(t, json.Unmarshal(evaluateGuestDenyRec.Body.Bytes(), &evaluateGuestDenyRes))
	require.False(t, evaluateGuestDenyRes.Allowed)
	require.Equal(t, "chapter_requires_authenticated", evaluateGuestDenyRes.ReasonCode)
}

func TestAccessHTTPEmergencyDenyPolicy(t *testing.T) {
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.6.0-test"}, zap.NewNop(), nil, registry)

	policyRec := performAccessJSONRequest(t, handler, http.MethodPost, "/access/policies", map[string]any{
		"key":               "feature.chapter.read.enabled",
		"effect":            "emergency_deny",
		"audience_kind":     "all",
		"audience_selector": "-",
		"scope_kind":        "feature",
		"scope_selector":    "chapter.read",
	})
	require.Equal(t, http.StatusCreated, policyRec.Code)

	evaluateRec := performAccessJSONRequest(t, handler, http.MethodPost, "/access/evaluate", map[string]any{
		"user_id":        uuid.NewString(),
		"permission":     "chapter.read.authenticated",
		"scope_kind":     "feature",
		"scope_selector": "chapter.read",
		"identity": map[string]any{
			"credential_id":  uuid.NewString(),
			"email_verified": true,
		},
		"user_signal": map[string]any{
			"account_state": "active",
		},
	})
	require.Equal(t, http.StatusOK, evaluateRec.Code)

	var evaluateRes struct {
		Allowed    bool   `json:"allowed"`
		ReasonCode string `json:"reason_code"`
	}
	require.NoError(t, json.Unmarshal(evaluateRec.Body.Bytes(), &evaluateRes))
	require.False(t, evaluateRes.Allowed)
	require.Equal(t, "emergency_deny", evaluateRes.ReasonCode)
}

func performAccessJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	return rec
}

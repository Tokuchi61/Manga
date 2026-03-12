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
	socialmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/social"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSocialHTTPFriendshipWallMessageAndRuntimeFlow(t *testing.T) {
	actorID := uuid.NewString()
	targetID := uuid.NewString()
	adminID := uuid.NewString()

	social := socialmodule.New()
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		social,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.14.0-test"}, zap.NewNop(), nil, registry)

	friendReq := performSocialJSONRequest(t, handler, http.MethodPost, "/social/friendships/requests", map[string]any{
		"target_user_id": targetID,
		"request_id":     "req-stage14-friend-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, friendReq.Code)

	respondReq := performSocialJSONRequest(t, handler, http.MethodPost, "/social/friendships/requests/"+actorID+"/respond", map[string]any{
		"action": "accept",
	}, targetID, "")
	require.Equal(t, http.StatusOK, respondReq.Code)

	friendsReq := performSocialRequest(t, handler, http.MethodGet, "/social/friendships", nil, actorID, "")
	require.Equal(t, http.StatusOK, friendsReq.Code)
	var friendsRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(friendsReq.Body.Bytes(), &friendsRes))
	require.Equal(t, 1, friendsRes.Count)

	followReq := performSocialRequest(t, handler, http.MethodPost, "/social/follow/"+targetID+"?request_id=req-stage14-follow-1", nil, actorID, "")
	require.Equal(t, http.StatusOK, followReq.Code)

	postReq := performSocialJSONRequest(t, handler, http.MethodPost, "/social/wall/posts", map[string]any{
		"body":       "hello social wall",
		"request_id": "req-stage14-wall-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, postReq.Code)
	var postRes struct {
		PostID string `json:"post_id"`
	}
	require.NoError(t, json.Unmarshal(postReq.Body.Bytes(), &postRes))
	require.NotEmpty(t, postRes.PostID)

	replyReq := performSocialJSONRequest(t, handler, http.MethodPost, "/social/wall/posts/"+postRes.PostID+"/replies", map[string]any{
		"body":       "reply social",
		"request_id": "req-stage14-reply-1",
	}, targetID, "")
	require.Equal(t, http.StatusOK, replyReq.Code)

	listWallReq := performSocialRequest(t, handler, http.MethodGet, "/social/wall/posts?owner_user_id="+actorID+"&include_replies=true", nil, actorID, "")
	require.Equal(t, http.StatusOK, listWallReq.Code)
	var wallRes struct {
		Count int `json:"count"`
		Items []struct {
			Replies []any `json:"replies"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(listWallReq.Body.Bytes(), &wallRes))
	require.Equal(t, 1, wallRes.Count)
	require.Len(t, wallRes.Items[0].Replies, 1)

	openThreadReq := performSocialJSONRequest(t, handler, http.MethodPost, "/social/messages/threads", map[string]any{
		"target_user_id": targetID,
	}, actorID, "")
	require.Equal(t, http.StatusOK, openThreadReq.Code)
	var openThreadRes struct {
		ThreadID string `json:"thread_id"`
	}
	require.NoError(t, json.Unmarshal(openThreadReq.Body.Bytes(), &openThreadRes))
	require.NotEmpty(t, openThreadRes.ThreadID)

	sendReq := performSocialJSONRequest(t, handler, http.MethodPost, "/social/messages/threads/"+openThreadRes.ThreadID+"/messages", map[string]any{
		"body":           "hello direct message",
		"request_id":     "req-stage14-msg-1",
		"correlation_id": "corr-stage14-msg-1",
	}, actorID, "")
	require.Equal(t, http.StatusOK, sendReq.Code)

	messageListReq := performSocialRequest(t, handler, http.MethodGet, "/social/messages/threads/"+openThreadRes.ThreadID+"/messages", nil, targetID, "")
	require.Equal(t, http.StatusOK, messageListReq.Code)
	var messageListRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(messageListReq.Body.Bytes(), &messageListRes))
	require.Equal(t, 1, messageListRes.Count)

	blockReq := performSocialJSONRequest(t, handler, http.MethodPost, "/social/relations/block/"+targetID, map[string]any{
		"enabled": true,
	}, actorID, "")
	require.Equal(t, http.StatusOK, blockReq.Code)

	openBlockedReq := performSocialJSONRequest(t, handler, http.MethodPost, "/social/messages/threads", map[string]any{
		"target_user_id": targetID,
	}, actorID, "")
	require.Equal(t, http.StatusForbidden, openBlockedReq.Code)

	disableFollowReq := performSocialJSONRequest(t, handler, http.MethodPost, "/social/admin/follow-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableFollowReq.Code)

	followBlockedReq := performSocialRequest(t, handler, http.MethodPost, "/social/follow/"+targetID+"?request_id=req-stage14-follow-2", nil, actorID, "")
	require.Equal(t, http.StatusForbidden, followBlockedReq.Code)
}

func performSocialJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performSocialRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performSocialRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
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

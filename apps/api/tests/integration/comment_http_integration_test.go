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
	chaptermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter"
	commentmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment"
	mangamodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCommentHTTPCreateThreadModerationLifecycleFlow(t *testing.T) {
	targetID := uuid.NewString()
	authorRoot := uuid.NewString()
	authorReply := uuid.NewString()
	moderatorID := uuid.NewString()

	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		mangamodule.New(),
		chaptermodule.New(),
		commentmodule.New(),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.9.0-test"}, zap.NewNop(), nil, registry)

	createRootRec := performCommentJSONRequest(t, handler, http.MethodPost, "/comments", map[string]any{
		"target_type": "manga",
		"target_id":   targetID,
		"content":     "Root comment",
	}, authorRoot, "")
	require.Equal(t, http.StatusCreated, createRootRec.Code)

	var rootRes struct {
		CommentID string `json:"comment_id"`
	}
	require.NoError(t, json.Unmarshal(createRootRec.Body.Bytes(), &rootRes))
	require.NotEmpty(t, rootRes.CommentID)

	createReplyRec := performCommentJSONRequest(t, handler, http.MethodPost, "/comments", map[string]any{
		"target_type":       "manga",
		"target_id":         targetID,
		"parent_comment_id": rootRes.CommentID,
		"content":           "Reply comment",
	}, authorReply, "")
	require.Equal(t, http.StatusCreated, createReplyRec.Code)

	var replyRes struct {
		CommentID string `json:"comment_id"`
	}
	require.NoError(t, json.Unmarshal(createReplyRec.Body.Bytes(), &replyRes))
	require.NotEmpty(t, replyRes.CommentID)

	listRec := performCommentRequest(t, handler, http.MethodGet, "/targets/manga/"+targetID+"/comments", nil, "", "")
	require.Equal(t, http.StatusOK, listRec.Code)
	var listRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(listRec.Body.Bytes(), &listRes))
	require.Equal(t, 1, listRes.Count)

	threadRec := performCommentRequest(t, handler, http.MethodGet, "/comments/"+rootRes.CommentID+"/thread", nil, "", "")
	require.Equal(t, http.StatusOK, threadRec.Code)
	var threadRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(threadRec.Body.Bytes(), &threadRes))
	require.Equal(t, 1, threadRes.Count)

	updateRec := performCommentJSONRequest(t, handler, http.MethodPatch, "/comments/"+rootRes.CommentID, map[string]any{
		"content": "Root edited",
	}, authorRoot, "")
	require.Equal(t, http.StatusOK, updateRec.Code)

	moderationRec := performCommentJSONRequest(t, handler, http.MethodPatch, "/comments/"+replyRes.CommentID+"/moderation", map[string]any{
		"moderation_status": "hidden",
	}, moderatorID, "moderator")
	require.Equal(t, http.StatusOK, moderationRec.Code)

	replyDetailHiddenRec := performCommentRequest(t, handler, http.MethodGet, "/comments/"+replyRes.CommentID, nil, "", "")
	require.Equal(t, http.StatusNotFound, replyDetailHiddenRec.Code)

	replyDetailVisibleForHiddenRec := performCommentRequest(t, handler, http.MethodGet, "/comments/"+replyRes.CommentID+"?include_hidden=true", nil, "", "")
	require.Equal(t, http.StatusOK, replyDetailVisibleForHiddenRec.Code)

	deleteRec := performCommentJSONRequest(t, handler, http.MethodDelete, "/comments/"+rootRes.CommentID, map[string]any{
		"reason": "cleanup",
	}, authorRoot, "")
	require.Equal(t, http.StatusOK, deleteRec.Code)

	rootDetailAfterDeleteRec := performCommentRequest(t, handler, http.MethodGet, "/comments/"+rootRes.CommentID, nil, "", "")
	require.Equal(t, http.StatusOK, rootDetailAfterDeleteRec.Code)
	var detailAfterDelete struct {
		Deleted bool   `json:"deleted"`
		Content string `json:"content"`
	}
	require.NoError(t, json.Unmarshal(rootDetailAfterDeleteRec.Body.Bytes(), &detailAfterDelete))
	require.True(t, detailAfterDelete.Deleted)
	require.Equal(t, "[deleted]", detailAfterDelete.Content)

	restoreRec := performCommentJSONRequest(t, handler, http.MethodPost, "/comments/"+rootRes.CommentID+"/restore", map[string]any{}, authorRoot, "")
	require.Equal(t, http.StatusOK, restoreRec.Code)
}

func TestCommentHTTPCreateValidationFailure(t *testing.T) {
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		mangamodule.New(),
		chaptermodule.New(),
		commentmodule.New(),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.9.0-test"}, zap.NewNop(), nil, registry)

	invalidCreateRec := performCommentJSONRequest(t, handler, http.MethodPost, "/comments", map[string]any{
		"target_type": "social",
		"target_id":   "not-a-uuid",
		"content":     "",
	}, uuid.NewString(), "")
	require.Equal(t, http.StatusBadRequest, invalidCreateRec.Code)
}

func performCommentJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performCommentRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performCommentRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
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

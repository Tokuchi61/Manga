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
	mangamodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestChapterHTTPCreatePublishReadNavigationLifecycleFlow(t *testing.T) {
	mangaID := uuid.NewString()

	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		mangamodule.New(),
		chaptermodule.New(),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.8.0-test"}, zap.NewNop(), nil, registry)

	createRec := performChapterJSONRequest(t, handler, http.MethodPost, "/chapters", map[string]any{
		"manga_id":           mangaID,
		"title":              "Chapter HTTP 1",
		"sequence_no":        1,
		"preview_enabled":    true,
		"preview_page_count": 1,
		"pages": []map[string]any{
			{
				"page_number": 1,
				"media_url":   "https://cdn.example.com/chapter-http-1.jpg",
				"width":       1200,
				"height":      1800,
			},
			{
				"page_number": 2,
				"media_url":   "https://cdn.example.com/chapter-http-2.jpg",
				"width":       1200,
				"height":      1800,
			},
		},
	})
	require.Equal(t, http.StatusCreated, createRec.Code)

	var createRes struct {
		ChapterID string `json:"chapter_id"`
	}
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &createRes))
	require.NotEmpty(t, createRes.ChapterID)

	listDraftRec := performChapterRequest(t, handler, http.MethodGet, "/manga/"+mangaID+"/chapters", nil)
	require.Equal(t, http.StatusOK, listDraftRec.Code)
	var listDraftRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(listDraftRec.Body.Bytes(), &listDraftRes))
	require.Equal(t, 0, listDraftRes.Count)

	publishRec := performChapterJSONRequest(t, handler, http.MethodPatch, "/chapters/"+createRes.ChapterID+"/publish", map[string]any{
		"action": "publish",
	})
	require.Equal(t, http.StatusOK, publishRec.Code)

	detailRec := performChapterRequest(t, handler, http.MethodGet, "/chapters/"+createRes.ChapterID, nil)
	require.Equal(t, http.StatusOK, detailRec.Code)

	readPreviewRec := performChapterRequest(t, handler, http.MethodGet, "/chapters/"+createRes.ChapterID+"/read?mode=preview", nil)
	require.Equal(t, http.StatusOK, readPreviewRec.Code)
	var readPreviewRes struct {
		PageCount int `json:"page_count"`
	}
	require.NoError(t, json.Unmarshal(readPreviewRec.Body.Bytes(), &readPreviewRes))
	require.Equal(t, 1, readPreviewRes.PageCount)

	navigationRec := performChapterRequest(t, handler, http.MethodGet, "/chapters/"+createRes.ChapterID+"/navigation", nil)
	require.Equal(t, http.StatusOK, navigationRec.Code)

	accessRec := performChapterJSONRequest(t, handler, http.MethodPatch, "/chapters/"+createRes.ChapterID+"/access", map[string]any{
		"vip_only": true,
	})
	require.Equal(t, http.StatusOK, accessRec.Code)

	mediaRec := performChapterJSONRequest(t, handler, http.MethodPatch, "/chapters/"+createRes.ChapterID+"/media-health", map[string]any{
		"media_health_status": "degraded",
	})
	require.Equal(t, http.StatusOK, mediaRec.Code)

	integrityRec := performChapterJSONRequest(t, handler, http.MethodPatch, "/chapters/"+createRes.ChapterID+"/integrity", map[string]any{
		"integrity_status": "passed",
	})
	require.Equal(t, http.StatusOK, integrityRec.Code)

	deleteRec := performChapterRequest(t, handler, http.MethodDelete, "/chapters/"+createRes.ChapterID, nil)
	require.Equal(t, http.StatusOK, deleteRec.Code)

	detailAfterDeleteRec := performChapterRequest(t, handler, http.MethodGet, "/chapters/"+createRes.ChapterID, nil)
	require.Equal(t, http.StatusNotFound, detailAfterDeleteRec.Code)

	restoreRec := performChapterRequest(t, handler, http.MethodPost, "/chapters/"+createRes.ChapterID+"/restore", bytes.NewReader(nil))
	require.Equal(t, http.StatusOK, restoreRec.Code)

	republishRec := performChapterJSONRequest(t, handler, http.MethodPatch, "/chapters/"+createRes.ChapterID+"/publish", map[string]any{
		"action": "publish",
	})
	require.Equal(t, http.StatusOK, republishRec.Code)

	detailAfterRestoreRec := performChapterRequest(t, handler, http.MethodGet, "/chapters/"+createRes.ChapterID, nil)
	require.Equal(t, http.StatusOK, detailAfterRestoreRec.Code)
}

func TestChapterHTTPCreateValidationFailure(t *testing.T) {
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		mangamodule.New(),
		chaptermodule.New(),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.8.0-test"}, zap.NewNop(), nil, registry)

	invalidCreateRec := performChapterJSONRequest(t, handler, http.MethodPost, "/chapters", map[string]any{
		"manga_id":    "not-a-uuid",
		"title":       "",
		"sequence_no": 0,
		"pages":       []any{},
	})
	require.Equal(t, http.StatusBadRequest, invalidCreateRec.Code)
}

func performChapterJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performChapterRequest(t, handler, method, path, bytes.NewReader(payload))
}

func performChapterRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader) *httptest.ResponseRecorder {
	t.Helper()

	bodyReader := bytes.NewReader(nil)
	if body != nil {
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

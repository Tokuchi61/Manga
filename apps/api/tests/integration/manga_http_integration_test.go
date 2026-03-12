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
	mangamodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMangaHTTPCreatePublishDiscoveryLifecycleFlow(t *testing.T) {
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		mangamodule.New(),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.7.0-test"}, zap.NewNop(), nil, registry)

	createRec := performMangaJSONRequest(t, handler, http.MethodPost, "/manga", map[string]any{
		"title":   "Crystal Blade",
		"summary": "A sword and world-building saga",
		"genres":  []string{"action", "fantasy"},
		"tags":    []string{"adventure"},
	})
	require.Equal(t, http.StatusCreated, createRec.Code)

	var createRes struct {
		MangaID string `json:"manga_id"`
	}
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &createRes))
	require.NotEmpty(t, createRes.MangaID)

	listDraftRec := performMangaRequest(t, handler, http.MethodGet, "/manga", nil)
	require.Equal(t, http.StatusOK, listDraftRec.Code)
	var listDraftRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(listDraftRec.Body.Bytes(), &listDraftRes))
	require.Equal(t, 0, listDraftRes.Count)

	publishRec := performMangaJSONRequest(t, handler, http.MethodPatch, "/manga/"+createRes.MangaID+"/publish", map[string]any{
		"action": "publish",
	})
	require.Equal(t, http.StatusOK, publishRec.Code)

	editorialRec := performMangaJSONRequest(t, handler, http.MethodPatch, "/manga/"+createRes.MangaID+"/editorial", map[string]any{
		"recommended":     true,
		"collection_keys": []string{"spring_pick"},
	})
	require.Equal(t, http.StatusOK, editorialRec.Code)

	countersRec := performMangaJSONRequest(t, handler, http.MethodPatch, "/manga/"+createRes.MangaID+"/counters", map[string]any{
		"chapter_count": 3,
		"comment_count": 14,
		"view_count":    150,
	})
	require.Equal(t, http.StatusOK, countersRec.Code)

	listRec := performMangaRequest(t, handler, http.MethodGet, "/manga?search=crystal", nil)
	require.Equal(t, http.StatusOK, listRec.Code)
	var listRes struct {
		Count int `json:"count"`
		Items []struct {
			MangaID      string `json:"manga_id"`
			ChapterCount int64  `json:"chapter_count"`
			CommentCount int64  `json:"comment_count"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(listRec.Body.Bytes(), &listRes))
	require.Equal(t, 1, listRes.Count)
	require.Equal(t, createRes.MangaID, listRes.Items[0].MangaID)
	require.EqualValues(t, 3, listRes.Items[0].ChapterCount)
	require.EqualValues(t, 14, listRes.Items[0].CommentCount)

	detailRec := performMangaRequest(t, handler, http.MethodGet, "/manga/"+createRes.MangaID, nil)
	require.Equal(t, http.StatusOK, detailRec.Code)

	discoveryRec := performMangaRequest(t, handler, http.MethodGet, "/manga/discovery?mode=recommended", nil)
	require.Equal(t, http.StatusOK, discoveryRec.Code)
	var discoveryRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(discoveryRec.Body.Bytes(), &discoveryRes))
	require.Equal(t, 1, discoveryRes.Count)

	deleteRec := performMangaRequest(t, handler, http.MethodDelete, "/manga/"+createRes.MangaID, nil)
	require.Equal(t, http.StatusOK, deleteRec.Code)

	detailAfterDeleteRec := performMangaRequest(t, handler, http.MethodGet, "/manga/"+createRes.MangaID, nil)
	require.Equal(t, http.StatusNotFound, detailAfterDeleteRec.Code)

	restoreRec := performMangaRequest(t, handler, http.MethodPost, "/manga/"+createRes.MangaID+"/restore", bytes.NewReader(nil))
	require.Equal(t, http.StatusOK, restoreRec.Code)
}

func TestMangaHTTPSortByPopular(t *testing.T) {
	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		mangamodule.New(),
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.7.0-test"}, zap.NewNop(), nil, registry)

	createAndPublish := func(title string, viewCount int) string {
		createRec := performMangaJSONRequest(t, handler, http.MethodPost, "/manga", map[string]any{
			"title":   title,
			"summary": title + " summary",
			"genres":  []string{"action"},
		})
		require.Equal(t, http.StatusCreated, createRec.Code)
		var createRes struct {
			MangaID string `json:"manga_id"`
		}
		require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &createRes))

		publishRec := performMangaJSONRequest(t, handler, http.MethodPatch, "/manga/"+createRes.MangaID+"/publish", map[string]any{"action": "publish"})
		require.Equal(t, http.StatusOK, publishRec.Code)

		syncRec := performMangaJSONRequest(t, handler, http.MethodPatch, "/manga/"+createRes.MangaID+"/counters", map[string]any{"view_count": viewCount})
		require.Equal(t, http.StatusOK, syncRec.Code)

		return createRes.MangaID
	}

	lowID := createAndPublish("Low View", 5)
	highID := createAndPublish("High View", 500)

	listRec := performMangaRequest(t, handler, http.MethodGet, "/manga?sort=popular", nil)
	require.Equal(t, http.StatusOK, listRec.Code)
	var listRes struct {
		Items []struct {
			MangaID string `json:"manga_id"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(listRec.Body.Bytes(), &listRes))
	require.GreaterOrEqual(t, len(listRes.Items), 2)
	require.Equal(t, highID, listRes.Items[0].MangaID)
	require.NotEqual(t, lowID, listRes.Items[0].MangaID)
}

func performMangaJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performMangaRequest(t, handler, method, path, bytes.NewReader(payload))
}

func performMangaRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader) *httptest.ResponseRecorder {
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

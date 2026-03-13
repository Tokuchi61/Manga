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
	historymodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/history"
	mangamodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHistoryHTTPIntakeLibraryTimelineAndRuntimeFlow(t *testing.T) {
	userID := uuid.NewString()
	adminID := uuid.NewString()

	manga := mangamodule.New()
	chapter := chaptermodule.New()
	chapter.SetMangaLookup(manga)

	history := historymodule.New()
	history.SetChapterSignalProvider(chapter)

	registry, err := modules.NewRegistry(
		authmodule.New(authmodule.RuntimeConfig{}),
		usermodule.New(),
		accessmodule.New(accessmodule.RuntimeConfig{}),
		manga,
		chapter,
		history,
	)
	require.NoError(t, err)

	handler := app.NewHTTPHandler(config.Config{AppVersion: "0.13.0-test"}, zap.NewNop(), nil, registry)

	mangaID := createMangaForHistoryTest(t, handler)
	chapterID := createChapterForHistoryTest(t, handler, mangaID)

	intakeRec := performHistoryJSONRequest(t, handler, http.MethodPost, "/history/intake/chapter", map[string]any{
		"chapter_id":     chapterID,
		"event":          "chapter.read.checkpoint",
		"page_number":    2,
		"request_id":     "req-stage13-1",
		"correlation_id": "corr-stage13-1",
	}, userID, "")
	require.Equal(t, http.StatusOK, intakeRec.Code)

	var intakeRes struct {
		LibraryEntryID string `json:"library_entry_id"`
		Created        bool   `json:"created"`
	}
	require.NoError(t, json.Unmarshal(intakeRec.Body.Bytes(), &intakeRes))
	require.True(t, intakeRes.Created)
	require.NotEmpty(t, intakeRes.LibraryEntryID)

	intakeAgainRec := performHistoryJSONRequest(t, handler, http.MethodPost, "/history/intake/chapter", map[string]any{
		"chapter_id":     chapterID,
		"event":          "chapter.read.checkpoint",
		"page_number":    2,
		"request_id":     "req-stage13-1",
		"correlation_id": "corr-stage13-1",
	}, userID, "")
	require.Equal(t, http.StatusOK, intakeAgainRec.Code)

	var intakeAgainRes struct {
		LibraryEntryID string `json:"library_entry_id"`
		Created        bool   `json:"created"`
	}
	require.NoError(t, json.Unmarshal(intakeAgainRec.Body.Bytes(), &intakeAgainRes))
	require.False(t, intakeAgainRes.Created)
	require.Equal(t, intakeRes.LibraryEntryID, intakeAgainRes.LibraryEntryID)

	continueRec := performHistoryRequest(t, handler, http.MethodGet, "/history/continue-reading", nil, userID, "")
	require.Equal(t, http.StatusOK, continueRec.Code)
	var continueRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(continueRec.Body.Bytes(), &continueRes))
	require.Equal(t, 1, continueRes.Count)

	timelineRec := performHistoryRequest(t, handler, http.MethodGet, "/history/timeline", nil, userID, "")
	require.Equal(t, http.StatusOK, timelineRec.Code)
	var timelineRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(timelineRec.Body.Bytes(), &timelineRes))
	require.Equal(t, 1, timelineRes.Count)

	bookmarkRec := performHistoryJSONRequest(t, handler, http.MethodPost, "/history/library/"+mangaID+"/bookmark", map[string]any{
		"bookmarked": true,
		"favorited":  true,
	}, userID, "")
	require.Equal(t, http.StatusOK, bookmarkRec.Code)

	shareRec := performHistoryJSONRequest(t, handler, http.MethodPost, "/history/library/"+mangaID+"/share", map[string]any{
		"share_public": true,
	}, userID, "")
	require.Equal(t, http.StatusOK, shareRec.Code)

	publicLibraryRec := performHistoryRequest(t, handler, http.MethodGet, "/history/public/"+userID+"/library", nil, "", "")
	require.Equal(t, http.StatusOK, publicLibraryRec.Code)
	var publicLibraryRes struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(publicLibraryRec.Body.Bytes(), &publicLibraryRes))
	require.Equal(t, 1, publicLibraryRes.Count)

	disableBookmarkWriteRec := performHistoryJSONRequest(t, handler, http.MethodPost, "/history/admin/bookmark-write-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableBookmarkWriteRec.Code)

	bookmarkBlockedRec := performHistoryJSONRequest(t, handler, http.MethodPost, "/history/library/"+mangaID+"/bookmark", map[string]any{
		"bookmarked": false,
		"favorited":  false,
	}, userID, "")
	require.Equal(t, http.StatusForbidden, bookmarkBlockedRec.Code)

	disableLibraryRec := performHistoryJSONRequest(t, handler, http.MethodPost, "/history/admin/library-state", map[string]any{
		"enabled": false,
	}, adminID, "admin")
	require.Equal(t, http.StatusOK, disableLibraryRec.Code)

	libraryHiddenRec := performHistoryRequest(t, handler, http.MethodGet, "/history/library", nil, userID, "")
	require.Equal(t, http.StatusNotFound, libraryHiddenRec.Code)
}

func createMangaForHistoryTest(t *testing.T, handler http.Handler) string {
	t.Helper()

	rec := performHistoryJSONRequest(t, handler, http.MethodPost, "/manga", map[string]any{
		"title":   "History Test Manga",
		"summary": "History stage coverage",
		"genres":  []string{"action"},
		"tags":    []string{"test"},
	}, "", "")
	require.Equal(t, http.StatusCreated, rec.Code)

	var res struct {
		MangaID string `json:"manga_id"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotEmpty(t, res.MangaID)
	return res.MangaID
}

func createChapterForHistoryTest(t *testing.T, handler http.Handler, mangaID string) string {
	t.Helper()

	rec := performHistoryJSONRequest(t, handler, http.MethodPost, "/chapters", map[string]any{
		"manga_id":           mangaID,
		"title":              "History Chapter",
		"sequence_no":        1,
		"preview_enabled":    true,
		"preview_page_count": 1,
		"pages": []map[string]any{
			{"page_number": 1, "media_url": "https://cdn.example.com/history-1.jpg", "width": 1200, "height": 1800},
			{"page_number": 2, "media_url": "https://cdn.example.com/history-2.jpg", "width": 1200, "height": 1800},
			{"page_number": 3, "media_url": "https://cdn.example.com/history-3.jpg", "width": 1200, "height": 1800},
		},
	}, "", "")
	require.Equal(t, http.StatusCreated, rec.Code)

	var res struct {
		ChapterID string `json:"chapter_id"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotEmpty(t, res.ChapterID)
	return res.ChapterID
}

func performHistoryJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, actorUserID string, roles string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	return performHistoryRequest(t, handler, method, path, bytes.NewReader(payload), actorUserID, roles)
}

func performHistoryRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, actorUserID string, roles string) *httptest.ResponseRecorder {
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

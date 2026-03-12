package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListContinueReading(w http.ResponseWriter, r *http.Request) {
	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.ListContinueReading(r.Context(), dto.ListContinueReadingRequest{
		UserID: userID,
		Limit:  parseIntQuery(r, "limit"),
		Offset: parseIntQuery(r, "offset"),
		SortBy: strings.TrimSpace(r.URL.Query().Get("sort_by")),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListLibrary(w http.ResponseWriter, r *http.Request) {
	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.ListLibrary(r.Context(), dto.ListLibraryRequest{
		UserID:     userID,
		Status:     strings.TrimSpace(r.URL.Query().Get("status")),
		Bookmarked: parseBoolPointerQuery(r, "bookmarked"),
		Favorited:  parseBoolPointerQuery(r, "favorited"),
		SharedOnly: parseBoolQuery(r, "shared_only"),
		Limit:      parseIntQuery(r, "limit"),
		Offset:     parseIntQuery(r, "offset"),
		SortBy:     strings.TrimSpace(r.URL.Query().Get("sort_by")),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListPublicLibrary(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListPublicLibrary(r.Context(), dto.ListPublicLibraryRequest{
		OwnerUserID: chi.URLParam(r, "user_id"),
		Limit:       parseIntQuery(r, "limit"),
		Offset:      parseIntQuery(r, "offset"),
		SortBy:      strings.TrimSpace(r.URL.Query().Get("sort_by")),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListTimeline(w http.ResponseWriter, r *http.Request) {
	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.ListTimeline(r.Context(), dto.ListTimelineRequest{
		UserID: userID,
		Event:  strings.TrimSpace(r.URL.Query().Get("event")),
		Limit:  parseIntQuery(r, "limit"),
		Offset: parseIntQuery(r, "offset"),
		SortBy: strings.TrimSpace(r.URL.Query().Get("sort_by")),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetRuntimeConfig(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetRuntimeConfig(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func parseIntQuery(r *http.Request, key string) int {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return value
}

func parseBoolQuery(r *http.Request, key string) bool {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return false
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false
	}
	return value
}

func parseBoolPointerQuery(r *http.Request, key string) *bool {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil
	}
	return &value
}

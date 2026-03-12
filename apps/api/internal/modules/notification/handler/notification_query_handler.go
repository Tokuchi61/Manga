package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListInbox(w http.ResponseWriter, r *http.Request) {
	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.ListInbox(r.Context(), dto.ListInboxRequest{
		UserID:     userID,
		Category:   strings.TrimSpace(r.URL.Query().Get("category")),
		State:      strings.TrimSpace(r.URL.Query().Get("state")),
		UnreadOnly: parseBoolQuery(r, "unread_only"),
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

func (h *HTTPHandler) GetNotificationDetail(w http.ResponseWriter, r *http.Request) {
	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.GetNotificationDetail(r.Context(), dto.GetNotificationDetailRequest{
		UserID:         userID,
		NotificationID: chi.URLParam(r, "notification_id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetPreference(w http.ResponseWriter, r *http.Request) {
	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.GetPreference(r.Context(), dto.GetPreferenceRequest{UserID: userID})
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

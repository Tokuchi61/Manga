package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListOwnSupport(w http.ResponseWriter, r *http.Request) {
	requesterID := strings.TrimSpace(r.URL.Query().Get("requester_user_id"))
	limit := parseIntQuery(r, "limit")
	offset := parseIntQuery(r, "offset")

	res, err := h.service.ListOwnSupport(r.Context(), dto.ListOwnSupportRequest{
		RequesterUserID: requesterID,
		Status:          strings.TrimSpace(r.URL.Query().Get("status")),
		SortBy:          strings.TrimSpace(r.URL.Query().Get("sort_by")),
		Limit:           limit,
		Offset:          offset,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetSupportDetail(w http.ResponseWriter, r *http.Request) {
	supportID := chi.URLParam(r, "support_id")
	requesterID := strings.TrimSpace(r.URL.Query().Get("requester_user_id"))

	res, err := h.service.GetSupportDetail(r.Context(), dto.GetSupportDetailRequest{
		SupportID:       supportID,
		RequesterUserID: requesterID,
		IncludeInternal: parseBoolQuery(r, "include_internal"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListReviewQueue(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListReviewQueue(r.Context(), dto.ListReviewQueueRequest{
		Status:   strings.TrimSpace(r.URL.Query().Get("status")),
		Priority: strings.TrimSpace(r.URL.Query().Get("priority")),
		Limit:    parseIntQuery(r, "limit"),
		Offset:   parseIntQuery(r, "offset"),
	})
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

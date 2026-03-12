package handler

import (
	"net/http"
	"strconv"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListQueue(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	response, err := h.service.ListQueue(r.Context(), dto.ListQueueRequest{
		Status:                  r.URL.Query().Get("status"),
		TargetType:              r.URL.Query().Get("target_type"),
		AssignedModeratorUserID: r.URL.Query().Get("assigned_moderator_user_id"),
		SortBy:                  r.URL.Query().Get("sort_by"),
		Limit:                   limit,
		Offset:                  offset,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *HTTPHandler) GetCaseDetail(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	response, err := h.service.GetCaseDetail(r.Context(), dto.GetCaseDetailRequest{CaseID: caseID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

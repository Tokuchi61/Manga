package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) UpdateProfileVisibility(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateProfileVisibilityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.UserID = chi.URLParam(r, "user_id")
	req.ViewerID = r.URL.Query().Get("viewer_id")

	res, err := h.service.UpdateProfileVisibility(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateHistoryVisibility(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateHistoryVisibilityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.UserID = chi.URLParam(r, "user_id")
	req.ViewerID = r.URL.Query().Get("viewer_id")

	res, err := h.service.UpdateHistoryVisibility(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

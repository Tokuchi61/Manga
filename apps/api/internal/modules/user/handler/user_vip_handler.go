package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) UpdateVIPState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateVIPStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.UserID = chi.URLParam(r, "user_id")

	res, err := h.service.UpdateVIPState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

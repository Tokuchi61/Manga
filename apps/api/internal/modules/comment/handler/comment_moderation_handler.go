package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) UpdateModeration(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateModerationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.CommentID = chi.URLParam(r, "comment_id")

	res, err := h.service.UpdateModeration(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

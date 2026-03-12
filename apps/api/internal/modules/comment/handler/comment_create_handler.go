package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
)

func (h *HTTPHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCommentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateComment(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, res)
}

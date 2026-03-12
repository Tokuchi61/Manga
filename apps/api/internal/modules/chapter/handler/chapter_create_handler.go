package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
)

func (h *HTTPHandler) CreateChapter(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateChapterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateChapter(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

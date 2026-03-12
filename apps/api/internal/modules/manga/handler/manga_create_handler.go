package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
)

func (h *HTTPHandler) CreateManga(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateMangaRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateManga(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

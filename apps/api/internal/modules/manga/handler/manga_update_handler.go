package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) UpdateMangaMetadata(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateMangaMetadataRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.MangaID = chi.URLParam(r, "manga_id")

	res, err := h.service.UpdateMangaMetadata(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

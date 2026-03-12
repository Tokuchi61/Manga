package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) SoftDeleteManga(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.SoftDeleteManga(r.Context(), dto.SoftDeleteMangaRequest{
		MangaID: chi.URLParam(r, "manga_id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) RestoreManga(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.RestoreManga(r.Context(), dto.RestoreMangaRequest{
		MangaID: chi.URLParam(r, "manga_id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

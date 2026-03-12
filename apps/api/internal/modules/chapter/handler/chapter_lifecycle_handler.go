package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) SoftDeleteChapter(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.SoftDeleteChapter(r.Context(), dto.SoftDeleteChapterRequest{ChapterID: chi.URLParam(r, "chapter_id")})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) RestoreChapter(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.RestoreChapter(r.Context(), dto.RestoreChapterRequest{ChapterID: chi.URLParam(r, "chapter_id")})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

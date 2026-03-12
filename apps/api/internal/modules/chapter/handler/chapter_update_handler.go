package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) UpdateChapter(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateChapterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ChapterID = chi.URLParam(r, "chapter_id")

	res, err := h.service.UpdateChapter(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ReorderChapter(w http.ResponseWriter, r *http.Request) {
	var req dto.ReorderChapterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ChapterID = chi.URLParam(r, "chapter_id")

	res, err := h.service.ReorderChapter(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

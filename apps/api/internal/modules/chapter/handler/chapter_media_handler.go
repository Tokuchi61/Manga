package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) UpdateMediaHealth(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateMediaHealthRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ChapterID = chi.URLParam(r, "chapter_id")

	res, err := h.service.UpdateMediaHealth(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateIntegrity(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateIntegrityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ChapterID = chi.URLParam(r, "chapter_id")

	res, err := h.service.UpdateIntegrity(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

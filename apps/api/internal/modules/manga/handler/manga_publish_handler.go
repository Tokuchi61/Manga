package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) UpdatePublishState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePublishStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.MangaID = chi.URLParam(r, "manga_id")

	res, err := h.service.UpdatePublishState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateVisibility(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateVisibilityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.MangaID = chi.URLParam(r, "manga_id")

	res, err := h.service.UpdateVisibility(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateEditorial(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEditorialRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.MangaID = chi.URLParam(r, "manga_id")

	res, err := h.service.UpdateEditorial(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) SyncCounters(w http.ResponseWriter, r *http.Request) {
	var req dto.SyncCountersRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.MangaID = chi.URLParam(r, "manga_id")

	res, err := h.service.SyncCounters(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

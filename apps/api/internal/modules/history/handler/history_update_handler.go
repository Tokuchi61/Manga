package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) UpdateBookmark(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateBookmarkRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.UserID = userID
	req.MangaID = chi.URLParam(r, "manga_id")

	res, err := h.service.UpdateBookmark(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateShare(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateShareRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.UserID = userID
	req.MangaID = chi.URLParam(r, "manga_id")

	res, err := h.service.UpdateShare(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateContinueReadingState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateContinueReadingStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateContinueReadingState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateLibraryState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateLibraryStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateLibraryState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateTimelineState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateTimelineStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateTimelineState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateBookmarkWriteState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateBookmarkWriteStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateBookmarkWriteState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

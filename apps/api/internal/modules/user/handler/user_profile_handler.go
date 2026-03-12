package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetPublicProfile(r.Context(), dto.GetPublicProfileRequest{
		UserID: chi.URLParam(r, "user_id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetOwnProfile(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetOwnProfile(r.Context(), dto.GetOwnProfileRequest{
		UserID:   chi.URLParam(r, "user_id"),
		ViewerID: r.URL.Query().Get("viewer_id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateProfileRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.UserID = chi.URLParam(r, "user_id")
	req.ViewerID = r.URL.Query().Get("viewer_id")

	res, err := h.service.UpdateProfile(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

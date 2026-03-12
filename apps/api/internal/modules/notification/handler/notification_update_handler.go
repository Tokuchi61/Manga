package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.MarkRead(r.Context(), dto.MarkReadRequest{
		UserID:         userID,
		NotificationID: chi.URLParam(r, "notification_id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdatePreference(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePreferenceRequest
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

	res, err := h.service.UpdatePreference(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateCategoryState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCategoryStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateCategoryState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateChannelState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateChannelStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateChannelState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateDigestState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateDigestStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateDigestState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateDeliveryPause(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateDeliveryPauseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateDeliveryPause(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

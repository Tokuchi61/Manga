package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
)

func (h *HTTPHandler) SendEmailVerification(w http.ResponseWriter, r *http.Request) {
	var req dto.SendVerificationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.service.SendEmailVerification(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, res)
}

func (h *HTTPHandler) ConfirmEmailVerification(w http.ResponseWriter, r *http.Request) {
	var req dto.ConfirmVerificationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.service.ConfirmEmailVerification(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

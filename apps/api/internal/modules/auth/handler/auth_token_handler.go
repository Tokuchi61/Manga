package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.service.RefreshToken(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	credentialID, ok := identity.CredentialID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_credential_id")
		return
	}
	req.CredentialID = credentialID

	res, err := h.service.Logout(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

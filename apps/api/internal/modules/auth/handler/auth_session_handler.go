package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	credentialID, ok := identity.CredentialID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_credential_id")
		return
	}

	req := dto.ListSessionsRequest{CredentialID: credentialID}
	res, err := h.service.ListSessions(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) RevokeCurrentSession(w http.ResponseWriter, r *http.Request) {
	var req dto.RevokeCurrentSessionRequest
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

	res, err := h.service.RevokeCurrentSession(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) RevokeOtherSessions(w http.ResponseWriter, r *http.Request) {
	var req dto.RevokeOtherSessionsRequest
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

	res, err := h.service.RevokeOtherSessions(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	var req dto.RevokeAllSessionsRequest
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

	res, err := h.service.RevokeAllSessions(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

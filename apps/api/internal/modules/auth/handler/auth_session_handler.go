package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
)

func (h *HTTPHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	req := dto.ListSessionsRequest{CredentialID: r.URL.Query().Get("credential_id")}
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
	res, err := h.service.RevokeAllSessions(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

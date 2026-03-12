package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	var req dto.EvaluateRequest
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
	if credentialID, ok := identity.CredentialID(r.Context()); ok {
		req.Identity = &dto.EvaluateIdentity{CredentialID: credentialID, EmailVerified: true}
	} else {
		req.Identity = nil
	}
	req.UserSignal = nil

	res, err := h.service.Evaluate(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

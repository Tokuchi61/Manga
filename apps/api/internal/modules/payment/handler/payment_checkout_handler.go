package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) StartCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var req dto.StartCheckoutSessionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.ActorUserID = actorUserID

	res, err := h.service.StartCheckoutSession(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

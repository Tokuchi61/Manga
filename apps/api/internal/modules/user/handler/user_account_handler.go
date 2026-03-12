package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) UpdateAccountState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateAccountStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	req.UserID = chi.URLParam(r, "user_id")
	req.ActorUserID = actorUserID
	if identity.HasAnyRole(r.Context(), "admin") {
		req.ActorScope = "admin"
	} else {
		req.ActorScope = "self"
	}

	res, err := h.service.UpdateAccountState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

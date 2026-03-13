package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) ListManaPackages(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.ListManaPackages(r.Context(), dto.ListManaPackagesRequest{ActorUserID: actorUserID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

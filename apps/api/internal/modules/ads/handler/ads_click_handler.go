package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) IntakeClick(w http.ResponseWriter, r *http.Request) {
	var req dto.IntakeClickRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if actorUserID, ok := identity.UserID(r.Context()); ok {
		req.UserID = actorUserID
	}

	res, err := h.service.IntakeClick(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

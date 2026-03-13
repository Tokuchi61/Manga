package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) GetRuntimeConfig(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetRuntimeConfig(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateMaintenanceState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateMaintenanceStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_or_invalid_actor_user_id")
		return
	}

	res, err := h.service.UpdateMaintenanceState(r.Context(), actorUserID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

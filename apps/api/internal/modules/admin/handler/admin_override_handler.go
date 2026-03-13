package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) ApplyOverride(w http.ResponseWriter, r *http.Request) {
	var req dto.ApplyOverrideRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_or_invalid_actor_user_id")
		return
	}

	res, err := h.service.ApplyOverride(r.Context(), actorUserID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListOverrides(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListOverrides(r.Context(), dto.ListOverridesRequest{
		TargetModule: parseStringQuery(r, "target_module"),
		Limit:        parseIntQuery(r, "limit"),
		Offset:       parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

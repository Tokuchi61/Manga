package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) IngestMissionProgress(w http.ResponseWriter, r *http.Request) {
	var req dto.IngestMissionProgressRequest
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
	req.MissionID = chi.URLParam(r, "mission_id")

	res, err := h.service.IngestMissionProgress(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

package handler

import (
	"net/http"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListActorMissions(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.ListActorMissions(r.Context(), dto.ListActorMissionsRequest{
		ActorUserID: actorUserID,
		Category:    strings.TrimSpace(r.URL.Query().Get("category")),
		State:       strings.TrimSpace(r.URL.Query().Get("state")),
		Limit:       parseIntQuery(r, "limit"),
		Offset:      parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetActorMissionDetail(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.GetActorMissionDetail(r.Context(), dto.GetActorMissionDetailRequest{
		ActorUserID: actorUserID,
		MissionID:   chi.URLParam(r, "mission_id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

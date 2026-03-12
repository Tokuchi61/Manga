package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ApplyAction(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	caseID := chi.URLParam(r, "case_id")
	var body struct {
		ActionType string `json:"action_type"`
		ReasonCode string `json:"reason_code,omitempty"`
		Summary    string `json:"summary,omitempty"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.ApplyAction(r.Context(), dto.ApplyActionRequest{
		CaseID:      caseID,
		ActorUserID: actorUserID,
		ActionType:  body.ActionType,
		ReasonCode:  body.ReasonCode,
		Summary:     body.Summary,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

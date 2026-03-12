package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) EscalateCase(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	caseID := chi.URLParam(r, "case_id")
	var body struct {
		Reason string `json:"reason"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.EscalateCase(r.Context(), dto.EscalateCaseRequest{
		CaseID:      caseID,
		ActorUserID: actorUserID,
		Reason:      body.Reason,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

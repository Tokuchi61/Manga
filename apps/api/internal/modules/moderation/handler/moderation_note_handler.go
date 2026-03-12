package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) AddModeratorNote(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	caseID := chi.URLParam(r, "case_id")
	var body struct {
		Body         string `json:"body"`
		InternalOnly bool   `json:"internal_only"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.AddModeratorNote(r.Context(), dto.AddModeratorNoteRequest{
		CaseID:       caseID,
		ActorUserID:  actorUserID,
		Body:         body.Body,
		InternalOnly: body.InternalOnly,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

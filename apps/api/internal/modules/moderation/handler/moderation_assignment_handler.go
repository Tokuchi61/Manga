package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) AssignCase(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	caseID := chi.URLParam(r, "case_id")
	var body struct {
		AssigneeUserID string `json:"assignee_user_id,omitempty"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.AssignCase(r.Context(), dto.AssignCaseRequest{
		CaseID:         caseID,
		ActorUserID:    actorUserID,
		AssigneeUserID: body.AssigneeUserID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *HTTPHandler) ReleaseCase(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	caseID := chi.URLParam(r, "case_id")
	response, err := h.service.ReleaseCase(r.Context(), dto.ReleaseCaseRequest{
		CaseID:      caseID,
		ActorUserID: actorUserID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

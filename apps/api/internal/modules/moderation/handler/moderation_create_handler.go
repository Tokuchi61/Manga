package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) CreateCaseFromSupportHandoff(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	var body struct {
		SupportID     string `json:"support_id"`
		RequestID     string `json:"request_id,omitempty"`
		CorrelationID string `json:"correlation_id,omitempty"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.CreateCaseFromSupportHandoff(r.Context(), dto.CreateCaseFromSupportHandoffRequest{
		SupportID:     body.SupportID,
		RequestID:     body.RequestID,
		CorrelationID: body.CorrelationID,
		ActorUserID:   actorUserID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	status := http.StatusCreated
	if !response.Created {
		status = http.StatusOK
	}
	writeJSON(w, status, response)
}

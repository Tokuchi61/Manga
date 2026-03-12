package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) IngestChapterSignal(w http.ResponseWriter, r *http.Request) {
	var req dto.IngestChapterSignalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.UserID = userID

	res, err := h.service.IngestChapterSignal(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

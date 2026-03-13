package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) StartImpersonation(w http.ResponseWriter, r *http.Request) {
	var req dto.StartImpersonationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_or_invalid_actor_user_id")
		return
	}

	res, err := h.service.StartImpersonation(r.Context(), actorUserID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) StopImpersonation(w http.ResponseWriter, r *http.Request) {
	var req dto.StopImpersonationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_or_invalid_actor_user_id")
		return
	}

	res, err := h.service.StopImpersonation(r.Context(), actorUserID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListImpersonationSessions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListImpersonationSessions(r.Context(), dto.ListImpersonationSessionsRequest{
		ActiveOnly: parseBoolQuery(r, "active_only"),
		Limit:      parseIntQuery(r, "limit"),
		Offset:     parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

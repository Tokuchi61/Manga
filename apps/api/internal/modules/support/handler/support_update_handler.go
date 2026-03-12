package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) AddReply(w http.ResponseWriter, r *http.Request) {
	supportID := chi.URLParam(r, "support_id")

	var req dto.AddSupportReplyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	isTeamActor := identity.HasAnyRole(r.Context(), "support_agent", "moderator", "admin")
	if !isTeamActor {
		req.Visibility = "public_to_requester"
	}

	req.SupportID = supportID
	req.ActorUserID = actorUserID
	req.ActorIsTeam = isTeamActor

	res, err := h.service.AddReply(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	supportID := chi.URLParam(r, "support_id")

	var req dto.UpdateSupportStatusRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.SupportID = supportID
	req.ReviewedByUserID = &actorUserID

	res, err := h.service.UpdateStatus(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	supportID := chi.URLParam(r, "support_id")

	var req dto.ResolveSupportRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.SupportID = supportID
	req.ReviewedByUserID = actorUserID

	res, err := h.service.Resolve(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) RequestModerationHandoff(w http.ResponseWriter, r *http.Request) {
	supportID := chi.URLParam(r, "support_id")

	var req dto.RequestModerationHandoffRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.SupportID = supportID

	res, err := h.service.RequestModerationHandoff(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

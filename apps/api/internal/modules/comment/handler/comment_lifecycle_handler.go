package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	var req dto.DeleteCommentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.CommentID = chi.URLParam(r, "comment_id")
	req.ActorUserID = actorUserID

	res, err := h.service.DeleteComment(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) RestoreComment(w http.ResponseWriter, r *http.Request) {
	var req dto.RestoreCommentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.CommentID = chi.URLParam(r, "comment_id")
	req.ActorUserID = actorUserID

	res, err := h.service.RestoreComment(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

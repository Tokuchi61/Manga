package handler

import (
	"net/http"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) CreateWallPost(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWallPostRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.ActorUserID = actorUserID

	res, err := h.service.CreateWallPost(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) CreateWallReply(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWallReplyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.ActorUserID = actorUserID
	req.PostID = chi.URLParam(r, "post_id")

	res, err := h.service.CreateWallReply(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListWallPosts(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.ListWallPosts(r.Context(), dto.ListWallPostsRequest{
		ActorUserID:    actorUserID,
		OwnerUserID:    strings.TrimSpace(r.URL.Query().Get("owner_user_id")),
		IncludeReplies: parseBoolQuery(r, "include_replies"),
		Limit:          parseIntQuery(r, "limit"),
		Offset:         parseIntQuery(r, "offset"),
		SortBy:         strings.TrimSpace(r.URL.Query().Get("sort_by")),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

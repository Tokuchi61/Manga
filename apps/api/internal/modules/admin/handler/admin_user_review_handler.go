package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) ReviewUser(w http.ResponseWriter, r *http.Request) {
	var req dto.ReviewUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_or_invalid_actor_user_id")
		return
	}

	res, err := h.service.ReviewUser(r.Context(), actorUserID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListUserReviews(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListUserReviews(r.Context(), dto.ListUserReviewsRequest{
		TargetUserID: parseStringQuery(r, "target_user_id"),
		Limit:        parseIntQuery(r, "limit"),
		Offset:       parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

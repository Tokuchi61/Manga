package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListCatalog(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.ListCatalog(r.Context(), dto.ListCatalogRequest{
		ActorUserID:     actorUserID,
		IncludeCampaign: parseBoolQuery(r, "include_campaign"),
		Limit:           parseIntQuery(r, "limit"),
		Offset:          parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetCatalogItem(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.GetCatalogItem(r.Context(), dto.GetCatalogItemRequest{
		ActorUserID:     actorUserID,
		ProductID:       chi.URLParam(r, "product_id"),
		IncludeCampaign: parseBoolQuery(r, "include_campaign"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

package handler

import (
	"net/http"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListInventoryEntries(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.ListInventoryEntries(r.Context(), dto.ListInventoryEntriesRequest{
		ActorUserID:  actorUserID,
		ItemType:     strings.TrimSpace(r.URL.Query().Get("item_type")),
		EquippedOnly: parseBoolQuery(r, "equipped_only"),
		SortBy:       strings.TrimSpace(r.URL.Query().Get("sort")),
		Limit:        parseIntQuery(r, "limit"),
		Offset:       parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetInventoryItemDetail(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	res, err := h.service.GetInventoryItemDetail(r.Context(), dto.GetInventoryItemDetailRequest{
		ActorUserID: actorUserID,
		ItemID:      chi.URLParam(r, "item_id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

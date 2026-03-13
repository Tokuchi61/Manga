package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/dto"
)

func (h *HTTPHandler) ResolvePlacements(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ResolvePlacements(r.Context(), dto.ResolvePlacementsRequest{
		Surface:    parseStringQuery(r, "surface"),
		TargetType: parseStringQuery(r, "target_type"),
		TargetID:   parseStringQuery(r, "target_id"),
		SessionID:  parseStringQuery(r, "session_id"),
		NoAds:      parseBoolQuery(r, "no_ads"),
		Limit:      parseIntQuery(r, "limit"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

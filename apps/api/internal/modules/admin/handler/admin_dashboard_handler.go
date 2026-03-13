package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
)

func (h *HTTPHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetDashboard(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListAuditTrail(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListAuditTrail(r.Context(), dto.ListAuditTrailRequest{
		ActionType: parseStringQuery(r, "action_type"),
		RiskLevel:  parseStringQuery(r, "risk_level"),
		Limit:      parseIntQuery(r, "limit"),
		Offset:     parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

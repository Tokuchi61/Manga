package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
)

func (h *HTTPHandler) CreatePolicyRule(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePolicyRuleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreatePolicyRule(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

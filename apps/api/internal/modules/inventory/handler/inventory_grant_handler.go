package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/dto"
)

func (h *HTTPHandler) AdminGrantItem(w http.ResponseWriter, r *http.Request) {
	var req dto.AdminGrantItemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.AdminGrantItem(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) AdminRevokeItem(w http.ResponseWriter, r *http.Request) {
	var req dto.AdminRevokeItemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.AdminRevokeItem(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

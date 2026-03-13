package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
)

func (h *HTTPHandler) ProcessProviderCallback(w http.ResponseWriter, r *http.Request) {
	var req dto.ProcessProviderCallbackRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.ProcessProviderCallback(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

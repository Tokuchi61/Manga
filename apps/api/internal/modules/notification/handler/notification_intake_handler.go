package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/dto"
)

func (h *HTTPHandler) CreateFromSupportSignal(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFromSupportSignalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateFromSupportSignal(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

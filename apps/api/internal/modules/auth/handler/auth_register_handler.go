package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
)

func (h *HTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Register(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

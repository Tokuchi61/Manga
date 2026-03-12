package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
)

func (h *HTTPHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePermissionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreatePermission(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *HTTPHandler) ListCanonicalPermissions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListCanonicalPermissions(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

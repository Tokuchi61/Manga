package handler

import (
	"net/http"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
)

func (h *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(req.Device) == "" {
		req.Device = r.UserAgent()
	}
	if strings.TrimSpace(req.IP) == "" {
		req.IP = extractIP(r)
	}

	res, err := h.service.Login(r.Context(), req, buildRequestMeta(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

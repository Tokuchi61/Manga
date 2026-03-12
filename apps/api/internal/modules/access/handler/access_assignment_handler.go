package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) AssignPermissionToRole(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignPermissionToRoleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.AssignPermissionToRole(r.Context(), chi.URLParam(r, "role_id"), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignRoleToUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.AssignRoleToUser(r.Context(), chi.URLParam(r, "user_id"), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) CreateTemporaryGrant(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTemporaryGrantRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateTemporaryGrant(r.Context(), chi.URLParam(r, "user_id"), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

package handler

import (
	"net/http"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/dto"
)

func (h *HTTPHandler) GetRuntimeConfig(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetRuntimeConfig(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateReadState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateReadStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.service.UpdateReadState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateClaimState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateClaimStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.service.UpdateClaimState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateConsumeState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateConsumeStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.service.UpdateConsumeState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateEquipStateRuntime(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEquipStateRuntimeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.service.UpdateEquipStateRuntime(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpsertItemDefinition(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertItemDefinitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.service.UpsertItemDefinition(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListItemDefinitions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListItemDefinitions(r.Context(), strings.TrimSpace(r.URL.Query().Get("item_type")))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

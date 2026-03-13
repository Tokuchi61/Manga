package handler

import (
	"net/http"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/dto"
)

func (h *HTTPHandler) GetRuntimeConfig(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetRuntimeConfig(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateSeasonState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateSeasonStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateSeasonState(r.Context(), req)
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

func (h *HTTPHandler) UpdatePremiumState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePremiumStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdatePremiumState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListSeasonDefinitions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListSeasonDefinitions(r.Context(), dto.ListSeasonDefinitionsRequest{
		State:  strings.TrimSpace(r.URL.Query().Get("state")),
		Limit:  parseIntQuery(r, "limit"),
		Offset: parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpsertSeasonDefinition(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertSeasonDefinitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpsertSeasonDefinition(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListTierDefinitions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListTierDefinitions(r.Context(), dto.ListTierDefinitionsRequest{
		SeasonID:   strings.TrimSpace(r.URL.Query().Get("season_id")),
		Track:      strings.TrimSpace(r.URL.Query().Get("track")),
		ActiveOnly: parseBoolQuery(r, "active_only"),
		Limit:      parseIntQuery(r, "limit"),
		Offset:     parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpsertTierDefinition(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertTierDefinitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpsertTierDefinition(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ResetRoyalPassProgress(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetRoyalPassProgressRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.ResetRoyalPassProgress(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

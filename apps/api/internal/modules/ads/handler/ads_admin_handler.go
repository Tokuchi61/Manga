package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/dto"
)

func (h *HTTPHandler) GetRuntimeConfig(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetRuntimeConfig(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateSurfaceState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateSurfaceStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateSurfaceState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdatePlacementState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePlacementStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdatePlacementState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateCampaignState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCampaignStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateCampaignState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateClickIntakeState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateClickIntakeStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateClickIntakeState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListPlacementDefinitions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListPlacementDefinitions(r.Context(), dto.ListPlacementDefinitionsRequest{
		Surface: parseStringQuery(r, "surface"),
		Visible: parseBoolQuery(r, "visible_only"),
		Limit:   parseIntQuery(r, "limit"),
		Offset:  parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpsertPlacementDefinition(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertPlacementDefinitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpsertPlacementDefinition(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListCampaignDefinitions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListCampaignDefinitions(r.Context(), dto.ListCampaignDefinitionsRequest{
		PlacementID: parseStringQuery(r, "placement_id"),
		State:       parseStringQuery(r, "state"),
		Limit:       parseIntQuery(r, "limit"),
		Offset:      parseIntQuery(r, "offset"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpsertCampaignDefinition(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertCampaignDefinitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpsertCampaignDefinition(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListCampaignAggregate(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListCampaignAggregate(r.Context(), parseIntQuery(r, "limit"), parseIntQuery(r, "offset"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

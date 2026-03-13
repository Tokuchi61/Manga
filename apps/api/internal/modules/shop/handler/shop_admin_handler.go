package handler

import (
	"net/http"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/dto"
)

func (h *HTTPHandler) GetRuntimeConfig(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetRuntimeConfig(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdateCatalogState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCatalogStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateCatalogState(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) UpdatePurchaseState(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePurchaseStateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdatePurchaseState(r.Context(), req)
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

func (h *HTTPHandler) ListProductDefinitions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListProductDefinitions(r.Context(), dto.ListProductDefinitionsRequest{
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

func (h *HTTPHandler) UpsertProductDefinition(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertProductDefinitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpsertProductDefinition(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListOfferDefinitions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListOfferDefinitions(r.Context(), dto.ListOfferDefinitionsRequest{
		ProductID:  strings.TrimSpace(r.URL.Query().Get("product_id")),
		Visibility: strings.TrimSpace(r.URL.Query().Get("visibility")),
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

func (h *HTTPHandler) UpsertOfferDefinition(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertOfferDefinitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpsertOfferDefinition(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

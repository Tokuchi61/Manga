package handler

import (
	"encoding/json"
	"errors"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/httpx"
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/service"
)

func decodeJSON(r *http.Request, out any) error {
	return httpx.DecodeJSON(r, out)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrNotFound), errors.Is(err, service.ErrReadDisabled):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrForbiddenAction),
		errors.Is(err, service.ErrClaimDisabled),
		errors.Is(err, service.ErrConsumeDisabled),
		errors.Is(err, service.ErrEquipDisabled):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, service.ErrConflict),
		errors.Is(err, service.ErrInsufficientQuantity),
		errors.Is(err, service.ErrItemAlreadyOwned):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

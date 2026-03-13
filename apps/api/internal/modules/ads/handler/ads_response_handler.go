package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/service"
)

func decodeJSON(r *http.Request, out any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrNotFound), errors.Is(err, service.ErrSurfaceDisabled), errors.Is(err, service.ErrPlacementDisabled), errors.Is(err, service.ErrCampaignDisabled):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrClickIntakeDisabled):
		writeError(w, http.StatusServiceUnavailable, err.Error())
	case errors.Is(err, service.ErrConflict):
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

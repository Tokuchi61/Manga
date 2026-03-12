package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/service"
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
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrForbiddenAction):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, service.ErrCategoryDisabled),
		errors.Is(err, service.ErrCategoryMuted),
		errors.Is(err, service.ErrChannelDisabled),
		errors.Is(err, service.ErrDigestDisabled),
		errors.Is(err, service.ErrSupportSignalInvalid):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrDeliveryPaused), errors.Is(err, service.ErrSupportSignalUnavailable):
		writeError(w, http.StatusServiceUnavailable, err.Error())
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

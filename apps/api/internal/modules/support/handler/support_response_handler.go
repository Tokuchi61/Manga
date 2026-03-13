package handler

import (
	"encoding/json"
	"errors"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/httpx"
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/service"
)

func decodeJSON(r *http.Request, out any) error {
	return httpx.DecodeJSON(r, out)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrValidation), errors.Is(err, service.ErrTargetRequired), errors.Is(err, service.ErrInvalidTarget):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrSupportAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrSupportNotFound), errors.Is(err, service.ErrSupportNotVisible):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrForbiddenAction):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, service.ErrInvalidStateTransition), errors.Is(err, service.ErrModerationHandoffNotAllowed), errors.Is(err, service.ErrAlreadyHandedOff):
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

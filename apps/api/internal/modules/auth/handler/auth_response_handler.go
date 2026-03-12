package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/service"
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
	case errors.Is(err, service.ErrCredentialAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrEmailNotVerified),
		errors.Is(err, service.ErrCredentialSuspended),
		errors.Is(err, service.ErrCredentialBanned):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, service.ErrLoginCooldownActive), errors.Is(err, service.ErrVerificationResendCooldown):
		writeError(w, http.StatusTooManyRequests, err.Error())
	case errors.Is(err, service.ErrSessionNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrTokenInvalidOrExpired):
		writeError(w, http.StatusBadRequest, err.Error())
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

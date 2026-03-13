package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/httpx"
)

func decodeJSON(r *http.Request, out any) error {
	return httpx.DecodeJSON(r, out)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrNotFound), errors.Is(err, service.ErrManaPurchaseDisabled), errors.Is(err, service.ErrTransactionReadDisabled):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrCheckoutDisabled), errors.Is(err, service.ErrCallbackIntakePaused):
		writeError(w, http.StatusServiceUnavailable, err.Error())
	case errors.Is(err, service.ErrCallbackSignature), errors.Is(err, service.ErrCallbackTimestamp):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrConflict), errors.Is(err, service.ErrForbiddenAction):
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

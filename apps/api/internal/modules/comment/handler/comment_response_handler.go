package handler

import (
	"encoding/json"
	"errors"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/httpx"
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/service"
)

func decodeJSON(r *http.Request, out any) error {
	return httpx.DecodeJSON(r, out)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrCommentAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrCommentNotFound), errors.Is(err, service.ErrCommentNotVisible), errors.Is(err, service.ErrTargetNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrForbiddenAction):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, service.ErrRateLimited):
		writeError(w, http.StatusTooManyRequests, err.Error())
	case errors.Is(err, service.ErrEditWindowExpired), errors.Is(err, service.ErrRestoreWindowExpired), errors.Is(err, service.ErrReplyDepthExceeded), errors.Is(err, service.ErrCommentLocked), errors.Is(err, service.ErrInvalidStateTransition):
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

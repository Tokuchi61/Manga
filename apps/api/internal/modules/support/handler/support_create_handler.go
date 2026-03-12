package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) CreateCommunication(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCommunicationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	requesterID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.RequesterUserID = requesterID

	res, err := h.service.CreateCommunication(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *HTTPHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTicketRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	requesterID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.RequesterUserID = requesterID

	res, err := h.service.CreateTicket(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *HTTPHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateReportRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	requesterID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}
	req.RequesterUserID = requesterID

	res, err := h.service.CreateReport(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

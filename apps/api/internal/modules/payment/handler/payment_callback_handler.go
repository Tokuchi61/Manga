package handler

import (
	"net/http"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/httpx"
)

const (
	HeaderPaymentCallbackSignature = "X-Payment-Signature"
	HeaderPaymentCallbackTimestamp = "X-Payment-Timestamp"
)

func (h *HTTPHandler) ProcessProviderCallback(w http.ResponseWriter, r *http.Request) {
	payload, err := httpx.ReadBody(r, httpx.DefaultMaxBodyBytes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	signature := r.Header.Get(HeaderPaymentCallbackSignature)
	timestamp := r.Header.Get(HeaderPaymentCallbackTimestamp)
	if err := h.service.VerifyProviderCallbackSignature(payload, signature, timestamp); err != nil {
		writeServiceError(w, err)
		return
	}

	var req dto.ProcessProviderCallbackRequest
	if err := httpx.DecodeJSONBytes(payload, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.ProcessProviderCallback(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

package dto

import "time"

// ProcessProviderCallbackRequest consumes provider callback payload.
type ProcessProviderCallbackRequest struct {
	ProviderEventID   string `json:"provider_event_id" validate:"required,max=128"`
	SessionID         string `json:"session_id" validate:"required,max=64"`
	ProviderReference string `json:"provider_reference,omitempty" validate:"omitempty,max=128"`
	Status            string `json:"status" validate:"required,oneof=success failed cancelled"`
}

// ProcessProviderCallbackResponse returns callback processing result.
type ProcessProviderCallbackResponse struct {
	Status            string    `json:"status"`
	TransactionID     string    `json:"transaction_id"`
	TransactionStatus string    `json:"transaction_status"`
	ProcessedAt       time.Time `json:"processed_at"`
}

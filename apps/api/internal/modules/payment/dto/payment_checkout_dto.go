package dto

import "time"

// StartCheckoutSessionRequest starts a provider checkout session.
type StartCheckoutSessionRequest struct {
	ActorUserID   string `json:"-" validate:"required,max=64"`
	PackageID     string `json:"package_id" validate:"required,max=64"`
	RequestID     string `json:"request_id" validate:"required,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
	Source        string `json:"source,omitempty" validate:"omitempty,oneof=external_provider recovery_replay"`
}

// StartCheckoutSessionResponse returns checkout-start result.
type StartCheckoutSessionResponse struct {
	Status        string     `json:"status"`
	TransactionID string     `json:"transaction_id"`
	SessionID     string     `json:"session_id"`
	Provider      string     `json:"provider"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	AmountMana    int        `json:"amount_mana"`
	PriceAmount   int        `json:"price_amount"`
	PriceCurrency string     `json:"price_currency"`
}

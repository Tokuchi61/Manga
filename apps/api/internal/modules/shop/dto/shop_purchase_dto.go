package dto

// CreatePurchaseIntentRequest creates idempotent purchase intent.
type CreatePurchaseIntentRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	ProductID     string `json:"product_id,omitempty" validate:"omitempty,max=64"`
	OfferID       string `json:"offer_id,omitempty" validate:"omitempty,max=64"`
	ActorLevel    int    `json:"actor_level,omitempty" validate:"omitempty,min=0,max=1000"`
	ActorVIP      bool   `json:"actor_vip,omitempty"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// CreatePurchaseIntentResponse returns purchase intent payload.
type CreatePurchaseIntentResponse struct {
	Status         string `json:"status"`
	IntentID       string `json:"intent_id"`
	ProductID      string `json:"product_id"`
	OfferID        string `json:"offer_id"`
	FinalPriceMana int    `json:"final_price_mana"`
	Currency       string `json:"currency"`
	PurchaseStatus string `json:"purchase_status"`
}

// RequestPurchaseRecoveryRequest requests purchase recovery for pending delivery.
type RequestPurchaseRecoveryRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
	IntentID    string `json:"intent_id" validate:"required,uuid4"`
}

// RequestPurchaseRecoveryResponse returns recovery request payload.
type RequestPurchaseRecoveryResponse struct {
	Status         string `json:"status"`
	IntentID       string `json:"intent_id"`
	PurchaseStatus string `json:"purchase_status"`
}

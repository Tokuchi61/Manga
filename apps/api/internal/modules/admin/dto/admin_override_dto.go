package dto

import "time"

// ApplyOverrideRequest applies high-privilege override action.
type ApplyOverrideRequest struct {
	RequestID         string     `json:"request_id" validate:"required,max=128"`
	CorrelationID     string     `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
	TargetModule      string     `json:"target_module" validate:"required,max=64"`
	TargetType        string     `json:"target_type" validate:"required,max=64"`
	TargetID          string     `json:"target_id" validate:"required,max=128"`
	Decision          string     `json:"decision" validate:"required,oneof=allow deny freeze reopen"`
	Reason            string     `json:"reason" validate:"required,max=512"`
	RiskLevel         string     `json:"risk_level" validate:"required,oneof=low medium high critical"`
	DoubleConfirmed   bool       `json:"double_confirmed"`
	ConfirmationToken string     `json:"confirmation_token,omitempty" validate:"omitempty,max=256"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
}

// OverrideResponse returns override payload.
type OverrideResponse struct {
	Status       string     `json:"status"`
	ActionID     string     `json:"action_id"`
	OverrideID   string     `json:"override_id"`
	TargetModule string     `json:"target_module"`
	TargetType   string     `json:"target_type"`
	TargetID     string     `json:"target_id"`
	Decision     string     `json:"decision"`
	RiskLevel    string     `json:"risk_level"`
	Active       bool       `json:"active"`
	Event        string     `json:"event"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ListOverridesRequest resolves override list payload.
type ListOverridesRequest struct {
	TargetModule string `json:"-" validate:"omitempty,max=64"`
	Limit        int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset       int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListOverridesResponse wraps override list payload.
type ListOverridesResponse struct {
	Items []OverrideResponse `json:"items"`
	Count int                `json:"count"`
}

package dto

import "time"

// UpdateVIPStateRequest updates VIP membership lifecycle.
type UpdateVIPStateRequest struct {
	UserID       string     `json:"-" validate:"required,uuid4"`
	Action       string     `json:"action" validate:"required,oneof=activate freeze resume deactivate"`
	EndsAt       *time.Time `json:"ends_at,omitempty"`
	FreezeReason string     `json:"freeze_reason,omitempty" validate:"omitempty,max=120"`
}

// VIPStateResponse returns current VIP state.
type VIPStateResponse struct {
	Status    string     `json:"status"`
	VIPActive bool       `json:"vip_active"`
	VIPFrozen bool       `json:"vip_frozen"`
	VIPEndsAt *time.Time `json:"vip_ends_at,omitempty"`
}

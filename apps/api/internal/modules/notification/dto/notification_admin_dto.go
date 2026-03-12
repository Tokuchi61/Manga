package dto

import "time"

// UpdateCategoryStateRequest updates category availability.
type UpdateCategoryStateRequest struct {
	Category string `json:"category" validate:"required,max=64"`
	Enabled  bool   `json:"enabled"`
}

// UpdateChannelStateRequest updates channel availability.
type UpdateChannelStateRequest struct {
	Channel string `json:"channel" validate:"required,oneof=in_app email push"`
	Enabled bool   `json:"enabled"`
}

// UpdateDigestStateRequest updates digest runtime switch.
type UpdateDigestStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateDeliveryPauseRequest updates module-level delivery pause state.
type UpdateDeliveryPauseRequest struct {
	Paused bool `json:"paused"`
}

// RuntimeConfigResponse is admin runtime control payload.
type RuntimeConfigResponse struct {
	CategoryEnabled map[string]bool `json:"category_enabled"`
	ChannelEnabled  map[string]bool `json:"channel_enabled"`
	DigestEnabled   bool            `json:"digest_enabled"`
	DeliveryPaused  bool            `json:"delivery_paused"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

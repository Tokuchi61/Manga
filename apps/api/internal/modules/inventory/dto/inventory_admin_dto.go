package dto

import "time"

// UpdateReadStateRequest updates inventory read runtime state.
type UpdateReadStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateClaimStateRequest updates inventory claim runtime state.
type UpdateClaimStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateConsumeStateRequest updates inventory consume runtime state.
type UpdateConsumeStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateEquipStateRuntimeRequest updates inventory equip runtime state.
type UpdateEquipStateRuntimeRequest struct {
	Enabled bool `json:"enabled"`
}

// RuntimeConfigResponse is inventory runtime control payload.
type RuntimeConfigResponse struct {
	ReadEnabled    bool      `json:"read_enabled"`
	ClaimEnabled   bool      `json:"claim_enabled"`
	ConsumeEnabled bool      `json:"consume_enabled"`
	EquipEnabled   bool      `json:"equip_enabled"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// UpsertItemDefinitionRequest creates or updates inventory item definition.
type UpsertItemDefinitionRequest struct {
	ItemID     string `json:"item_id" validate:"required,max=128"`
	ItemType   string `json:"item_type" validate:"required,max=64"`
	Stackable  bool   `json:"stackable"`
	Equipable  bool   `json:"equipable"`
	Consumable bool   `json:"consumable"`
	MaxStack   int    `json:"max_stack" validate:"omitempty,min=1,max=100000"`
}

// ItemDefinitionResponse returns inventory item definition.
type ItemDefinitionResponse struct {
	ItemID     string    `json:"item_id"`
	ItemType   string    `json:"item_type"`
	Stackable  bool      `json:"stackable"`
	Equipable  bool      `json:"equipable"`
	Consumable bool      `json:"consumable"`
	MaxStack   int       `json:"max_stack"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ListItemDefinitionsResponse wraps definition list response.
type ListItemDefinitionsResponse struct {
	Items []ItemDefinitionResponse `json:"items"`
	Count int                      `json:"count"`
}

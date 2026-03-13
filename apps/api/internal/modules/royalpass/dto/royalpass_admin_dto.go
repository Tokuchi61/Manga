package dto

import "time"

// UpdateSeasonStateRequest updates season surface runtime state.
type UpdateSeasonStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateClaimStateRequest updates claim surface runtime state.
type UpdateClaimStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdatePremiumStateRequest updates premium surface runtime state.
type UpdatePremiumStateRequest struct {
	Enabled bool `json:"enabled"`
}

// RuntimeConfigResponse is royalpass runtime control payload.
type RuntimeConfigResponse struct {
	SeasonEnabled  bool      `json:"season_enabled"`
	ClaimEnabled   bool      `json:"claim_enabled"`
	PremiumEnabled bool      `json:"premium_enabled"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// UpsertSeasonDefinitionRequest creates or updates season definition.
type UpsertSeasonDefinitionRequest struct {
	SeasonID string     `json:"season_id" validate:"required,max=64"`
	Title    string     `json:"title" validate:"required,max=160"`
	State    string     `json:"state" validate:"required,oneof=draft active paused ended archived"`
	StartsAt *time.Time `json:"starts_at,omitempty"`
	EndsAt   *time.Time `json:"ends_at,omitempty"`
}

// SeasonDefinitionResponse returns season definition payload.
type SeasonDefinitionResponse struct {
	SeasonID  string     `json:"season_id"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	StartsAt  *time.Time `json:"starts_at,omitempty"`
	EndsAt    *time.Time `json:"ends_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ListSeasonDefinitionsRequest resolves season definitions.
type ListSeasonDefinitionsRequest struct {
	State  string `json:"-" validate:"omitempty,oneof=draft active paused ended archived"`
	Limit  int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListSeasonDefinitionsResponse wraps season definition list.
type ListSeasonDefinitionsResponse struct {
	Items []SeasonDefinitionResponse `json:"items"`
	Count int                        `json:"count"`
}

// UpsertTierDefinitionRequest creates or updates tier definition.
type UpsertTierDefinitionRequest struct {
	SeasonID       string `json:"season_id" validate:"required,max=64"`
	TierNumber     int    `json:"tier_number" validate:"required,min=1,max=100000"`
	Track          string `json:"track" validate:"required,oneof=free premium"`
	RequiredPoints int    `json:"required_points" validate:"required,min=1,max=1000000"`
	RewardItemID   string `json:"reward_item_id,omitempty" validate:"omitempty,max=128"`
	RewardQuantity int    `json:"reward_quantity" validate:"required,min=1,max=1000000"`
	Active         bool   `json:"active"`
}

// TierDefinitionResponse returns tier definition payload.
type TierDefinitionResponse struct {
	SeasonID       string    `json:"season_id"`
	TierNumber     int       `json:"tier_number"`
	Track          string    `json:"track"`
	RequiredPoints int       `json:"required_points"`
	RewardItemID   string    `json:"reward_item_id,omitempty"`
	RewardQuantity int       `json:"reward_quantity"`
	Active         bool      `json:"active"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ListTierDefinitionsRequest resolves tier definitions.
type ListTierDefinitionsRequest struct {
	SeasonID   string `json:"-" validate:"omitempty,max=64"`
	Track      string `json:"-" validate:"omitempty,oneof=free premium"`
	ActiveOnly bool   `json:"-"`
	Limit      int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset     int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListTierDefinitionsResponse wraps tier definitions.
type ListTierDefinitionsResponse struct {
	Items []TierDefinitionResponse `json:"items"`
	Count int                      `json:"count"`
}

// ResetRoyalPassProgressRequest resets actor progress for season or all seasons.
type ResetRoyalPassProgressRequest struct {
	TargetUserID string `json:"target_user_id" validate:"required,uuid4"`
	SeasonID     string `json:"season_id,omitempty" validate:"omitempty,max=64"`
}

// ResetRoyalPassProgressResponse returns reset operation result.
type ResetRoyalPassProgressResponse struct {
	Status       string `json:"status"`
	DeletedCount int    `json:"deleted_count"`
}

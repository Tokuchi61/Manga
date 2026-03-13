package dto

import "time"

// UpdateSurfaceStateRequest updates ads surface runtime state.
type UpdateSurfaceStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdatePlacementStateRequest updates placement resolve runtime state.
type UpdatePlacementStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateCampaignStateRequest updates campaign serve runtime state.
type UpdateCampaignStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateClickIntakeStateRequest updates click intake runtime state.
type UpdateClickIntakeStateRequest struct {
	Enabled bool `json:"enabled"`
}

// RuntimeConfigResponse is ads runtime control payload.
type RuntimeConfigResponse struct {
	SurfaceEnabled     bool      `json:"surface_enabled"`
	PlacementEnabled   bool      `json:"placement_enabled"`
	CampaignEnabled    bool      `json:"campaign_enabled"`
	ClickIntakeEnabled bool      `json:"click_intake_enabled"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// UpsertPlacementDefinitionRequest creates or updates placement definition.
type UpsertPlacementDefinitionRequest struct {
	PlacementID  string `json:"placement_id" validate:"required,max=64"`
	Surface      string `json:"surface" validate:"required,oneof=home listing manga chapter"`
	TargetType   string `json:"target_type" validate:"required,oneof=none manga chapter"`
	TargetID     string `json:"target_id,omitempty" validate:"omitempty,max=64"`
	Visible      bool   `json:"visible"`
	Priority     int    `json:"priority" validate:"omitempty,min=0,max=1000"`
	FrequencyCap int    `json:"frequency_cap" validate:"omitempty,min=0,max=500"`
}

// PlacementDefinitionResponse returns placement payload.
type PlacementDefinitionResponse struct {
	PlacementID  string    `json:"placement_id"`
	Surface      string    `json:"surface"`
	TargetType   string    `json:"target_type"`
	TargetID     string    `json:"target_id,omitempty"`
	Visible      bool      `json:"visible"`
	Priority     int       `json:"priority"`
	FrequencyCap int       `json:"frequency_cap"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListPlacementDefinitionsRequest resolves placement list.
type ListPlacementDefinitionsRequest struct {
	Surface string `json:"-" validate:"omitempty,oneof=home listing manga chapter"`
	Visible bool   `json:"-"`
	Limit   int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset  int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListPlacementDefinitionsResponse wraps placement list payload.
type ListPlacementDefinitionsResponse struct {
	Items []PlacementDefinitionResponse `json:"items"`
	Count int                           `json:"count"`
}

// UpsertCampaignDefinitionRequest creates or updates campaign definition.
type UpsertCampaignDefinitionRequest struct {
	CampaignID  string     `json:"campaign_id" validate:"required,max=64"`
	PlacementID string     `json:"placement_id" validate:"required,max=64"`
	Name        string     `json:"name" validate:"required,max=160"`
	State       string     `json:"state" validate:"required,oneof=draft active paused ended"`
	CreativeURL string     `json:"creative_url" validate:"required,max=400"`
	ClickURL    string     `json:"click_url" validate:"required,max=400"`
	Weight      int        `json:"weight" validate:"omitempty,min=1,max=1000"`
	StartsAt    *time.Time `json:"starts_at,omitempty"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
}

// CampaignDefinitionResponse returns campaign payload.
type CampaignDefinitionResponse struct {
	CampaignID  string     `json:"campaign_id"`
	PlacementID string     `json:"placement_id"`
	Name        string     `json:"name"`
	State       string     `json:"state"`
	CreativeURL string     `json:"creative_url"`
	ClickURL    string     `json:"click_url"`
	Weight      int        `json:"weight"`
	StartsAt    *time.Time `json:"starts_at,omitempty"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ListCampaignDefinitionsRequest resolves campaign list.
type ListCampaignDefinitionsRequest struct {
	PlacementID string `json:"-" validate:"omitempty,max=64"`
	State       string `json:"-" validate:"omitempty,oneof=draft active paused ended"`
	Limit       int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset      int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListCampaignDefinitionsResponse wraps campaign list payload.
type ListCampaignDefinitionsResponse struct {
	Items []CampaignDefinitionResponse `json:"items"`
	Count int                          `json:"count"`
}

// CampaignAggregateResponse returns campaign counter payload.
type CampaignAggregateResponse struct {
	CampaignID      string    `json:"campaign_id"`
	ImpressionCount int       `json:"impression_count"`
	ClickCount      int       `json:"click_count"`
	CTRPercent      float64   `json:"ctr_percent"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ListCampaignAggregateResponse wraps aggregate list payload.
type ListCampaignAggregateResponse struct {
	Items []CampaignAggregateResponse `json:"items"`
	Count int                         `json:"count"`
}

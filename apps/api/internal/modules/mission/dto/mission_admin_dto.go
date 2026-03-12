package dto

import "time"

// UpdateReadStateRequest updates mission read runtime state.
type UpdateReadStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateClaimStateRequest updates mission claim runtime state.
type UpdateClaimStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateProgressIngestStateRequest updates mission ingest runtime state.
type UpdateProgressIngestStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateDailyResetHourRequest updates mission daily reset hour config.
type UpdateDailyResetHourRequest struct {
	Hour int `json:"hour" validate:"required,min=0,max=23"`
}

// RuntimeConfigResponse is mission runtime control payload.
type RuntimeConfigResponse struct {
	ReadEnabled           bool      `json:"read_enabled"`
	ClaimEnabled          bool      `json:"claim_enabled"`
	ProgressIngestEnabled bool      `json:"progress_ingest_enabled"`
	DailyResetHourUTC     int       `json:"daily_reset_hour_utc"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// UpsertMissionDefinitionRequest creates or updates mission definition.
type UpsertMissionDefinitionRequest struct {
	MissionID      string     `json:"mission_id" validate:"required,max=128"`
	Category       string     `json:"category" validate:"required,oneof=daily weekly monthly event level"`
	Title          string     `json:"title" validate:"required,max=160"`
	ObjectiveType  string     `json:"objective_type" validate:"required,max=64"`
	TargetCount    int        `json:"target_count" validate:"required,min=1,max=1000000"`
	RewardItemID   string     `json:"reward_item_id,omitempty" validate:"omitempty,max=128"`
	RewardQuantity int        `json:"reward_quantity" validate:"required,min=1,max=1000000"`
	Active         bool       `json:"active"`
	StartsAt       *time.Time `json:"starts_at,omitempty"`
	EndsAt         *time.Time `json:"ends_at,omitempty"`
}

// MissionDefinitionResponse returns mission definition payload.
type MissionDefinitionResponse struct {
	MissionID      string     `json:"mission_id"`
	Category       string     `json:"category"`
	Title          string     `json:"title"`
	ObjectiveType  string     `json:"objective_type"`
	TargetCount    int        `json:"target_count"`
	RewardItemID   string     `json:"reward_item_id,omitempty"`
	RewardQuantity int        `json:"reward_quantity"`
	Active         bool       `json:"active"`
	StartsAt       *time.Time `json:"starts_at,omitempty"`
	EndsAt         *time.Time `json:"ends_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ListMissionDefinitionsRequest resolves mission definition list.
type ListMissionDefinitionsRequest struct {
	Category   string `json:"-" validate:"omitempty,oneof=daily weekly monthly event level"`
	ActiveOnly bool   `json:"-"`
	Limit      int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset     int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListMissionDefinitionsResponse wraps mission definition list response.
type ListMissionDefinitionsResponse struct {
	Items []MissionDefinitionResponse `json:"items"`
	Count int                         `json:"count"`
}

// ResetMissionProgressRequest resets mission progress for a target actor.
type ResetMissionProgressRequest struct {
	TargetUserID string `json:"target_user_id" validate:"required,uuid4"`
	MissionID    string `json:"mission_id,omitempty" validate:"omitempty,max=128"`
}

// ResetMissionProgressResponse returns reset operation result.
type ResetMissionProgressResponse struct {
	Status       string `json:"status"`
	DeletedCount int    `json:"deleted_count"`
}

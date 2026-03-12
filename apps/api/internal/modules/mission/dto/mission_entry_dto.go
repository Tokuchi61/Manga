package dto

import "time"

// ListActorMissionsRequest resolves actor mission list.
type ListActorMissionsRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
	Category    string `json:"-" validate:"omitempty,oneof=daily weekly monthly event level"`
	State       string `json:"-" validate:"omitempty,oneof=active completed claimed expired"`
	Limit       int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset      int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// GetActorMissionDetailRequest resolves own mission detail.
type GetActorMissionDetailRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
	MissionID   string `json:"-" validate:"required,max=128"`
}

// MissionProgressItemResponse is mission progress payload.
type MissionProgressItemResponse struct {
	MissionID      string     `json:"mission_id"`
	Category       string     `json:"category"`
	Title          string     `json:"title"`
	ObjectiveType  string     `json:"objective_type"`
	TargetCount    int        `json:"target_count"`
	ProgressCount  int        `json:"progress_count"`
	RewardItemID   string     `json:"reward_item_id,omitempty"`
	RewardQuantity int        `json:"reward_quantity"`
	Status         string     `json:"status"`
	PeriodKey      string     `json:"period_key"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	ClaimedAt      *time.Time `json:"claimed_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ListActorMissionsResponse wraps actor mission listing.
type ListActorMissionsResponse struct {
	Items []MissionProgressItemResponse `json:"items"`
	Count int                           `json:"count"`
}

// MissionDetailResponse returns own mission detail.
type MissionDetailResponse struct {
	MissionID      string     `json:"mission_id"`
	Category       string     `json:"category"`
	Title          string     `json:"title"`
	ObjectiveType  string     `json:"objective_type"`
	TargetCount    int        `json:"target_count"`
	ProgressCount  int        `json:"progress_count"`
	RewardItemID   string     `json:"reward_item_id,omitempty"`
	RewardQuantity int        `json:"reward_quantity"`
	Status         string     `json:"status"`
	PeriodKey      string     `json:"period_key"`
	StartsAt       *time.Time `json:"starts_at,omitempty"`
	EndsAt         *time.Time `json:"ends_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	ClaimedAt      *time.Time `json:"claimed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

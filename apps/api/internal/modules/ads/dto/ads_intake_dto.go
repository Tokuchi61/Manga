package dto

import "time"

// IntakeImpressionRequest consumes impression intake payload.
type IntakeImpressionRequest struct {
	RequestID   string `json:"request_id" validate:"required,max=128"`
	PlacementID string `json:"placement_id" validate:"required,max=64"`
	CampaignID  string `json:"campaign_id" validate:"required,max=64"`
	SessionID   string `json:"session_id,omitempty" validate:"omitempty,max=128"`
	UserID      string `json:"user_id,omitempty" validate:"omitempty,max=64"`
}

// IntakeImpressionResponse returns intake result payload.
type IntakeImpressionResponse struct {
	Status       string    `json:"status"`
	ImpressionID string    `json:"impression_id,omitempty"`
	CampaignID   string    `json:"campaign_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// IntakeClickRequest consumes click intake payload.
type IntakeClickRequest struct {
	RequestID      string `json:"request_id" validate:"required,max=128"`
	PlacementID    string `json:"placement_id" validate:"required,max=64"`
	CampaignID     string `json:"campaign_id" validate:"required,max=64"`
	SessionID      string `json:"session_id,omitempty" validate:"omitempty,max=128"`
	UserID         string `json:"user_id,omitempty" validate:"omitempty,max=64"`
	InvalidTraffic bool   `json:"invalid_traffic"`
}

// IntakeClickResponse returns click intake result payload.
type IntakeClickResponse struct {
	Status     string    `json:"status"`
	ClickID    string    `json:"click_id,omitempty"`
	CampaignID string    `json:"campaign_id"`
	CreatedAt  time.Time `json:"created_at"`
}

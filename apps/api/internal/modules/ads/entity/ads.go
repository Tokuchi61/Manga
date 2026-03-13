package entity

import "time"

const (
	CampaignStateDraft  = "draft"
	CampaignStateActive = "active"
	CampaignStatePaused = "paused"
	CampaignStateEnded  = "ended"
)

const (
	ImpressionStatusAccepted = "accepted"
	ImpressionStatusIgnored  = "ignored"
)

const (
	ClickStatusAccepted = "accepted"
	ClickStatusIgnored  = "ignored"
)

// PlacementDefinition stores stage-20 placement metadata.
type PlacementDefinition struct {
	PlacementID  string
	Surface      string
	TargetType   string
	TargetID     string
	Visible      bool
	Priority     int
	FrequencyCap int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CampaignDefinition stores stage-20 campaign metadata.
type CampaignDefinition struct {
	CampaignID  string
	PlacementID string
	Name        string
	State       string
	CreativeURL string
	ClickURL    string
	Weight      int
	StartsAt    *time.Time
	EndsAt      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ImpressionLog stores accepted impression intake entries.
type ImpressionLog struct {
	ImpressionID string
	RequestID    string
	PlacementID  string
	CampaignID   string
	SessionID    string
	UserID       string
	Status       string
	CreatedAt    time.Time
}

// ClickLog stores accepted click intake entries.
type ClickLog struct {
	ClickID        string
	RequestID      string
	PlacementID    string
	CampaignID     string
	SessionID      string
	UserID         string
	Status         string
	InvalidTraffic bool
	CreatedAt      time.Time
}

// CampaignAggregate stores simple campaign counters.
type CampaignAggregate struct {
	CampaignID      string
	ImpressionCount int
	ClickCount      int
	UpdatedAt       time.Time
}

// RuntimeConfig stores stage-20 runtime controls.
type RuntimeConfig struct {
	SurfaceEnabled     bool
	PlacementEnabled   bool
	CampaignEnabled    bool
	ClickIntakeEnabled bool
	UpdatedAt          time.Time
}

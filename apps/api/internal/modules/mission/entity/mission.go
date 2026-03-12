package entity

import "time"

// MissionDefinition stores stage-16 mission definition metadata.
type MissionDefinition struct {
	MissionID      string
	Category       string
	Title          string
	ObjectiveType  string
	TargetCount    int
	RewardItemID   string
	RewardQuantity int
	Active         bool
	StartsAt       *time.Time
	EndsAt         *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// MissionProgress stores actor mission progression for a period.
type MissionProgress struct {
	UserID            string
	MissionID         string
	PeriodKey         string
	ProgressCount     int
	Completed         bool
	Claimed           bool
	LastRequestID     string
	LastCorrelationID string
	CompletedAt       *time.Time
	ClaimedAt         *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// RuntimeConfig stores stage-16 mission runtime controls.
type RuntimeConfig struct {
	ReadEnabled           bool
	ClaimEnabled          bool
	ProgressIngestEnabled bool
	DailyResetHourUTC     int
	UpdatedAt             time.Time
}

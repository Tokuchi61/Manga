package entity

import "time"

const (
	TrackFree    = "free"
	TrackPremium = "premium"
)

const (
	SeasonStateDraft    = "draft"
	SeasonStateActive   = "active"
	SeasonStatePaused   = "paused"
	SeasonStateEnded    = "ended"
	SeasonStateArchived = "archived"
)

// SeasonDefinition stores stage-17 season metadata.
type SeasonDefinition struct {
	SeasonID  string
	Title     string
	State     string
	StartsAt  *time.Time
	EndsAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TierDefinition stores stage-17 track tier metadata.
type TierDefinition struct {
	SeasonID       string
	TierNumber     int
	Track          string
	RequiredPoints int
	RewardItemID   string
	RewardQuantity int
	Active         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// UserProgress stores actor progress and claim eligibility state.
type UserProgress struct {
	UserID                  string
	SeasonID                string
	Points                  int
	PremiumActivated        bool
	PremiumActivationSource string
	PremiumActivationRef    string
	ClaimedTiers            map[string]time.Time
	LastRequestID           string
	LastCorrelationID       string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// RuntimeConfig stores stage-17 royalpass runtime controls.
type RuntimeConfig struct {
	SeasonEnabled  bool
	ClaimEnabled   bool
	PremiumEnabled bool
	UpdatedAt      time.Time
}

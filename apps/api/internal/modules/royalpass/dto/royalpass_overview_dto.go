package dto

import "time"

// GetActorSeasonOverviewRequest resolves actor season overview.
type GetActorSeasonOverviewRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
	SeasonID    string `json:"-" validate:"omitempty,max=64"`
}

// TierProgressResponse returns free/premium tier status payload.
type TierProgressResponse struct {
	TierNumber     int        `json:"tier_number"`
	Track          string     `json:"track"`
	RequiredPoints int        `json:"required_points"`
	RewardItemID   string     `json:"reward_item_id,omitempty"`
	RewardQuantity int        `json:"reward_quantity"`
	Active         bool       `json:"active"`
	Unlocked       bool       `json:"unlocked"`
	Claimed        bool       `json:"claimed"`
	ClaimedAt      *time.Time `json:"claimed_at,omitempty"`
}

// SeasonOverviewResponse returns actor royalpass season overview.
type SeasonOverviewResponse struct {
	SeasonID         string                 `json:"season_id"`
	Title            string                 `json:"title"`
	State            string                 `json:"state"`
	StartsAt         *time.Time             `json:"starts_at,omitempty"`
	EndsAt           *time.Time             `json:"ends_at,omitempty"`
	Points           int                    `json:"points"`
	PremiumActivated bool                   `json:"premium_activated"`
	Items            []TierProgressResponse `json:"items"`
	Count            int                    `json:"count"`
}

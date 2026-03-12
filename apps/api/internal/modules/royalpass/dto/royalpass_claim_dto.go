package dto

// ClaimTierRewardRequest triggers tier claim request.
type ClaimTierRewardRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	SeasonID      string `json:"season_id,omitempty" validate:"omitempty,max=64"`
	TierNumber    int    `json:"tier_number" validate:"required,min=1,max=100000"`
	Track         string `json:"track" validate:"required,oneof=free premium"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// ClaimTierRewardResponse returns tier claim request result payload.
type ClaimTierRewardResponse struct {
	Status         string `json:"status"`
	SeasonID       string `json:"season_id"`
	TierNumber     int    `json:"tier_number"`
	Track          string `json:"track"`
	RewardItemID   string `json:"reward_item_id,omitempty"`
	RewardQuantity int    `json:"reward_quantity"`
}

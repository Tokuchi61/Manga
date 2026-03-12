package dto

// ClaimMissionRequest triggers mission reward claim request.
type ClaimMissionRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	MissionID     string `json:"-" validate:"required,max=128"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// ClaimMissionResponse returns mission claim request result.
type ClaimMissionResponse struct {
	Status         string                      `json:"status"`
	RewardItemID   string                      `json:"reward_item_id,omitempty"`
	RewardQuantity int                         `json:"reward_quantity"`
	Mission        MissionProgressItemResponse `json:"mission"`
}

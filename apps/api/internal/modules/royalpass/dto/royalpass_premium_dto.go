package dto

// ActivatePremiumTrackRequest ingests premium activation reference.
type ActivatePremiumTrackRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	SeasonID      string `json:"season_id,omitempty" validate:"omitempty,max=64"`
	SourceType    string `json:"source_type" validate:"required,max=64"`
	ActivationRef string `json:"activation_ref" validate:"required,max=128"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// ActivatePremiumTrackResponse returns premium activation result.
type ActivatePremiumTrackResponse struct {
	Status           string `json:"status"`
	SeasonID         string `json:"season_id"`
	PremiumActivated bool   `json:"premium_activated"`
}

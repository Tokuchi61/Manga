package dto

// CreateCaseFromSupportHandoffRequest defines support->moderation linked case creation payload.
type CreateCaseFromSupportHandoffRequest struct {
	SupportID     string `json:"support_id" validate:"required,uuid"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,min=3,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,min=3,max=128"`
	ActorUserID   string `json:"-"`
}

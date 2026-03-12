package dto

// AddSupportReplyRequest adds support reply/internal note.
type AddSupportReplyRequest struct {
	SupportID   string `json:"-" validate:"required,uuid4"`
	ActorUserID string `json:"actor_user_id" validate:"required,uuid4"`
	Message     string `json:"message" validate:"required,min=1,max=5000"`
	Visibility  string `json:"visibility" validate:"required,oneof=public_to_requester internal_only"`
}

// UpdateSupportStatusRequest updates support status and review fields.
type UpdateSupportStatusRequest struct {
	SupportID        string  `json:"-" validate:"required,uuid4"`
	Status           string  `json:"status" validate:"required,oneof=open triaged waiting_user waiting_team resolved rejected closed spam"`
	AssigneeUserID   *string `json:"assignee_user_id,omitempty" validate:"omitempty,uuid4"`
	ReviewedByUserID *string `json:"reviewed_by_user_id,omitempty" validate:"omitempty,uuid4"`
}

// ResolveSupportRequest resolves support with resolution note.
type ResolveSupportRequest struct {
	SupportID        string `json:"-" validate:"required,uuid4"`
	ReviewedByUserID string `json:"reviewed_by_user_id" validate:"required,uuid4"`
	ResolutionNote   string `json:"resolution_note" validate:"required,min=1,max=5000"`
}

// RequestModerationHandoffRequest requests linked moderation case handoff.
type RequestModerationHandoffRequest struct {
	SupportID     string `json:"-" validate:"required,uuid4"`
	ReasonCode    string `json:"reason_code,omitempty" validate:"omitempty,max=64"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// OperationResponse is generic operation result payload.
type OperationResponse struct {
	Status string `json:"status"`
}

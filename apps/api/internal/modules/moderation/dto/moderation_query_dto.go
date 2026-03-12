package dto

import "time"

type ListQueueRequest struct {
	Status                 string `json:"status,omitempty" validate:"omitempty,min=2,max=64"`
	TargetType             string `json:"target_type,omitempty" validate:"omitempty,min=3,max=32"`
	AssignedModeratorUserID string `json:"assigned_moderator_user_id,omitempty" validate:"omitempty,uuid"`
	SortBy                 string `json:"sort_by,omitempty" validate:"omitempty,min=3,max=32"`
	Limit                  int    `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Offset                 int    `json:"offset,omitempty" validate:"omitempty,min=0,max=10000"`
}

type GetCaseDetailRequest struct {
	CaseID string `json:"-" validate:"required,uuid"`
}

type NoteResponse struct {
	NoteID      string    `json:"note_id"`
	AuthorUserID string    `json:"author_user_id"`
	Body        string    `json:"body"`
	InternalOnly bool      `json:"internal_only"`
	CreatedAt   time.Time `json:"created_at"`
}

type ActionResponse struct {
	ActionID     string    `json:"action_id"`
	ActorUserID  string    `json:"actor_user_id"`
	ActionType   string    `json:"action_type"`
	ReasonCode   string    `json:"reason_code"`
	Summary      string    `json:"summary"`
	ActionResult string    `json:"action_result"`
	CreatedAt    time.Time `json:"created_at"`
}

type ModerationCaseResponse struct {
	CaseID                  string           `json:"case_id"`
	Source                  string           `json:"source"`
	SourceRefID             *string          `json:"source_ref_id,omitempty"`
	RequestID               string           `json:"request_id"`
	CorrelationID           string           `json:"correlation_id,omitempty"`
	TargetType              string           `json:"target_type"`
	TargetID                string           `json:"target_id"`
	Status                  string           `json:"status"`
	AssignmentStatus        string           `json:"assignment_status"`
	AssignedModeratorUserID *string          `json:"assigned_moderator_user_id,omitempty"`
	EscalationStatus        string           `json:"escalation_status"`
	EscalationReason        string           `json:"escalation_reason,omitempty"`
	ActionResult            string           `json:"action_result"`
	CreatedAt               time.Time        `json:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at"`
	Notes                   []NoteResponse   `json:"notes,omitempty"`
	Actions                 []ActionResponse `json:"actions,omitempty"`
}

type ListQueueResponse struct {
	Cases []ModerationCaseResponse `json:"cases"`
	Count int                      `json:"count"`
}

type CreateCaseFromSupportHandoffResponse struct {
	Case    ModerationCaseResponse `json:"case"`
	Created bool                   `json:"created"`
}

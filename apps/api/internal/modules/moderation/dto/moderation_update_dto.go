package dto

type AssignCaseRequest struct {
	CaseID         string `json:"-" validate:"required,uuid"`
	ActorUserID    string `json:"-" validate:"required,uuid"`
	AssigneeUserID string `json:"assignee_user_id,omitempty" validate:"omitempty,uuid"`
}

type ReleaseCaseRequest struct {
	CaseID      string `json:"-" validate:"required,uuid"`
	ActorUserID string `json:"-" validate:"required,uuid"`
}

type AddModeratorNoteRequest struct {
	CaseID      string `json:"-" validate:"required,uuid"`
	ActorUserID string `json:"-" validate:"required,uuid"`
	Body        string `json:"body" validate:"required,min=3,max=4096"`
	InternalOnly bool   `json:"internal_only"`
}

type ApplyActionRequest struct {
	CaseID      string `json:"-" validate:"required,uuid"`
	ActorUserID string `json:"-" validate:"required,uuid"`
	ActionType  string `json:"action_type" validate:"required,min=3,max=64"`
	ReasonCode  string `json:"reason_code,omitempty" validate:"omitempty,max=64"`
	Summary     string `json:"summary,omitempty" validate:"omitempty,max=512"`
}

type EscalateCaseRequest struct {
	CaseID      string `json:"-" validate:"required,uuid"`
	ActorUserID string `json:"-" validate:"required,uuid"`
	Reason      string `json:"reason" validate:"required,min=3,max=512"`
}

type OperationResponse struct {
	Status string `json:"status"`
}

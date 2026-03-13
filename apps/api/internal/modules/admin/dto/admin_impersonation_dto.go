package dto

import "time"

// StartImpersonationRequest starts high-risk impersonation session.
type StartImpersonationRequest struct {
	RequestID         string `json:"request_id" validate:"required,max=128"`
	CorrelationID     string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
	TargetUserID      string `json:"target_user_id" validate:"required,max=64"`
	Reason            string `json:"reason" validate:"required,max=512"`
	RiskLevel         string `json:"risk_level" validate:"required,oneof=low medium high critical"`
	DurationMinutes   int    `json:"duration_minutes" validate:"omitempty,min=1,max=120"`
	DoubleConfirmed   bool   `json:"double_confirmed"`
	ConfirmationToken string `json:"confirmation_token,omitempty" validate:"omitempty,max=256"`
}

// StopImpersonationRequest stops active impersonation session.
type StopImpersonationRequest struct {
	RequestID         string `json:"request_id" validate:"required,max=128"`
	CorrelationID     string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
	SessionID         string `json:"session_id" validate:"required,max=64"`
	Reason            string `json:"reason" validate:"required,max=512"`
	RiskLevel         string `json:"risk_level" validate:"required,oneof=low medium high critical"`
	DoubleConfirmed   bool   `json:"double_confirmed"`
	ConfirmationToken string `json:"confirmation_token,omitempty" validate:"omitempty,max=256"`
}

// ImpersonationSessionResponse returns impersonation session payload.
type ImpersonationSessionResponse struct {
	Status       string     `json:"status"`
	ActionID     string     `json:"action_id"`
	SessionID    string     `json:"session_id"`
	ActorUserID  string     `json:"actor_user_id"`
	TargetUserID string     `json:"target_user_id"`
	RiskLevel    string     `json:"risk_level"`
	Reason       string     `json:"reason"`
	Active       bool       `json:"active"`
	StartedAt    time.Time  `json:"started_at"`
	EndsAt       *time.Time `json:"ends_at,omitempty"`
	ExpiresAt    time.Time  `json:"expires_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ListImpersonationSessionsRequest resolves impersonation session list payload.
type ListImpersonationSessionsRequest struct {
	ActiveOnly bool `json:"-"`
	Limit      int  `json:"-" validate:"omitempty,min=1,max=100"`
	Offset     int  `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListImpersonationSessionsResponse wraps impersonation session list payload.
type ListImpersonationSessionsResponse struct {
	Items []ImpersonationSessionResponse `json:"items"`
	Count int                            `json:"count"`
}

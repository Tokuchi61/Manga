package entity

import "time"

// ImpersonationSession stores stage-21 admin impersonation lifecycle state.
type ImpersonationSession struct {
	SessionID    string
	ActionID     string
	ActorUserID  string
	TargetUserID string
	Reason       string
	RiskLevel    string
	Active       bool
	StartedAt    time.Time
	EndedAt      *time.Time
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

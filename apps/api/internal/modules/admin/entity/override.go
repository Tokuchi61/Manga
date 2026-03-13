package entity

import "time"

const (
	OverrideDecisionAllow  = "allow"
	OverrideDecisionDeny   = "deny"
	OverrideDecisionFreeze = "freeze"
	OverrideDecisionReopen = "reopen"
)

// OverrideRecord stores stage-21 admin hard override entries.
type OverrideRecord struct {
	OverrideID   string
	ActionID     string
	TargetModule string
	TargetType   string
	TargetID     string
	Decision     string
	Reason       string
	RiskLevel    string
	Active       bool
	ExpiresAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

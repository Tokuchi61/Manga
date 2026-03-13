package entity

import "time"

const (
	ActionTypeSettingChanged     = "setting_changed"
	ActionTypeOverrideApplied    = "override_applied"
	ActionTypeUserReviewed       = "user_reviewed"
	ActionTypeImpersonationStart = "impersonation_started"
	ActionTypeImpersonationStop  = "impersonation_stopped"
)

const (
	ActionResultSuccess = "success"
)

const (
	RiskLevelLow      = "low"
	RiskLevelMedium   = "medium"
	RiskLevelHigh     = "high"
	RiskLevelCritical = "critical"
)

// AdminAction stores immutable admin audit action metadata.
type AdminAction struct {
	ActionID                   string
	ActionType                 string
	ActorUserID                string
	TargetUserID               string
	TargetModule               string
	TargetType                 string
	TargetID                   string
	RequestID                  string
	CorrelationID              string
	Reason                     string
	Result                     string
	RiskLevel                  string
	RequiresDoubleConfirmation bool
	DoubleConfirmed            bool
	ConfirmationToken          string
	Metadata                   map[string]string
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

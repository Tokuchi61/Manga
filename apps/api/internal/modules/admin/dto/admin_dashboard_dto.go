package dto

import "time"

// AuditActionResponse returns immutable admin action payload.
type AuditActionResponse struct {
	ActionID                   string            `json:"action_id"`
	ActionType                 string            `json:"action_type"`
	ActorUserID                string            `json:"actor_user_id,omitempty"`
	TargetUserID               string            `json:"target_user_id,omitempty"`
	TargetModule               string            `json:"target_module,omitempty"`
	TargetType                 string            `json:"target_type,omitempty"`
	TargetID                   string            `json:"target_id,omitempty"`
	RequestID                  string            `json:"request_id,omitempty"`
	CorrelationID              string            `json:"correlation_id,omitempty"`
	Reason                     string            `json:"reason"`
	Result                     string            `json:"result"`
	RiskLevel                  string            `json:"risk_level"`
	RequiresDoubleConfirmation bool              `json:"requires_double_confirmation"`
	DoubleConfirmed            bool              `json:"double_confirmed"`
	Metadata                   map[string]string `json:"metadata,omitempty"`
	CreatedAt                  time.Time         `json:"created_at"`
}

// ListAuditTrailRequest resolves audit list payload.
type ListAuditTrailRequest struct {
	ActionType string `json:"-" validate:"omitempty,max=64"`
	RiskLevel  string `json:"-" validate:"omitempty,oneof=low medium high critical"`
	Limit      int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset     int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListAuditTrailResponse wraps audit list payload.
type ListAuditTrailResponse struct {
	Items []AuditActionResponse `json:"items"`
	Count int                   `json:"count"`
}

// DashboardResponse returns stage-21 admin summary.
type DashboardResponse struct {
	MaintenanceEnabled  bool                  `json:"maintenance_enabled"`
	TotalActions        int                   `json:"total_actions"`
	TotalOverrides      int                   `json:"total_overrides"`
	TotalUserReviews    int                   `json:"total_user_reviews"`
	ActiveImpersonation int                   `json:"active_impersonation"`
	LatestActions       []AuditActionResponse `json:"latest_actions"`
	GeneratedAt         time.Time             `json:"generated_at"`
}

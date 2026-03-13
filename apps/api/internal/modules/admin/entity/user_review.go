package entity

import "time"

const (
	UserReviewDecisionWarning     = "warning"
	UserReviewDecisionRestriction = "restriction"
	UserReviewDecisionSuspend     = "suspend"
	UserReviewDecisionBan         = "ban"
	UserReviewDecisionClear       = "clear"
)

// UserReviewRecord stores stage-21 admin user review actions.
type UserReviewRecord struct {
	ReviewID     string
	ActionID     string
	TargetUserID string
	Decision     string
	Reason       string
	RiskLevel    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

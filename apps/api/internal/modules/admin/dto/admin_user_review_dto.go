package dto

import "time"

// ReviewUserRequest applies user review decision.
type ReviewUserRequest struct {
	RequestID         string `json:"request_id" validate:"required,max=128"`
	CorrelationID     string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
	TargetUserID      string `json:"target_user_id" validate:"required,max=64"`
	Decision          string `json:"decision" validate:"required,oneof=warning restriction suspend ban clear"`
	Reason            string `json:"reason" validate:"required,max=512"`
	RiskLevel         string `json:"risk_level" validate:"required,oneof=low medium high critical"`
	DoubleConfirmed   bool   `json:"double_confirmed"`
	ConfirmationToken string `json:"confirmation_token,omitempty" validate:"omitempty,max=256"`
}

// UserReviewResponse returns review payload.
type UserReviewResponse struct {
	Status       string    `json:"status"`
	ActionID     string    `json:"action_id"`
	ReviewID     string    `json:"review_id"`
	TargetUserID string    `json:"target_user_id"`
	Decision     string    `json:"decision"`
	RiskLevel    string    `json:"risk_level"`
	Event        string    `json:"event"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListUserReviewsRequest resolves user review list payload.
type ListUserReviewsRequest struct {
	TargetUserID string `json:"-" validate:"omitempty,max=64"`
	Limit        int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset       int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListUserReviewsResponse wraps review list payload.
type ListUserReviewsResponse struct {
	Items []UserReviewResponse `json:"items"`
	Count int                  `json:"count"`
}

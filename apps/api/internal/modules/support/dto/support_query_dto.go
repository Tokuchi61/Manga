package dto

import "time"

// ListOwnSupportRequest defines requester-owned support listing controls.
type ListOwnSupportRequest struct {
	RequesterUserID string `json:"-" validate:"required,uuid4"`
	Status          string `json:"-" validate:"omitempty,oneof=open triaged waiting_user waiting_team resolved rejected closed spam"`
	SortBy          string `json:"-" validate:"omitempty,oneof=newest oldest priority"`
	Limit           int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset          int    `json:"-" validate:"omitempty,min=0"`
}

// GetSupportDetailRequest resolves support detail.
type GetSupportDetailRequest struct {
	SupportID       string `json:"-" validate:"required,uuid4"`
	RequesterUserID string `json:"-" validate:"required,uuid4"`
	IncludeInternal bool   `json:"-"`
}

// ListReviewQueueRequest resolves support review queue.
type ListReviewQueueRequest struct {
	Status   string `json:"-" validate:"omitempty,oneof=open triaged waiting_user waiting_team resolved rejected closed spam"`
	Priority string `json:"-" validate:"omitempty,oneof=low normal high urgent"`
	Limit    int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset   int    `json:"-" validate:"omitempty,min=0"`
}

// SupportReplyResponse is support reply payload item.
type SupportReplyResponse struct {
	ReplyID      string    `json:"reply_id"`
	AuthorUserID string    `json:"author_user_id"`
	Message      string    `json:"message"`
	Visibility   string    `json:"visibility"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SupportListItemResponse is support listing item payload.
type SupportListItemResponse struct {
	SupportID                    string     `json:"support_id"`
	RequesterUserID              string     `json:"requester_user_id"`
	SupportKind                  string     `json:"support_kind"`
	Category                     string     `json:"category"`
	Priority                     string     `json:"priority"`
	ReasonCode                   string     `json:"reason_code,omitempty"`
	Status                       string     `json:"status"`
	TargetType                   *string    `json:"target_type,omitempty"`
	TargetID                     *string    `json:"target_id,omitempty"`
	DuplicateOfSupportID         *string    `json:"duplicate_of_support_id,omitempty"`
	AssigneeUserID               *string    `json:"assignee_user_id,omitempty"`
	ReviewedByUserID             *string    `json:"reviewed_by_user_id,omitempty"`
	ResolvedAt                   *time.Time `json:"resolved_at,omitempty"`
	ModerationHandoffRequestedAt *time.Time `json:"moderation_handoff_requested_at,omitempty"`
	CreatedAt                    time.Time  `json:"created_at"`
	UpdatedAt                    time.Time  `json:"updated_at"`
}

// ListOwnSupportResponse wraps support list payload.
type ListOwnSupportResponse struct {
	Items []SupportListItemResponse `json:"items"`
	Count int                       `json:"count"`
}

// SupportDetailResponse is support detail payload.
type SupportDetailResponse struct {
	SupportID                    string                 `json:"support_id"`
	RequesterUserID              string                 `json:"requester_user_id"`
	SupportKind                  string                 `json:"support_kind"`
	Category                     string                 `json:"category"`
	Priority                     string                 `json:"priority"`
	ReasonCode                   string                 `json:"reason_code,omitempty"`
	ReasonText                   string                 `json:"reason_text"`
	Status                       string                 `json:"status"`
	TargetType                   *string                `json:"target_type,omitempty"`
	TargetID                     *string                `json:"target_id,omitempty"`
	DuplicateOfSupportID         *string                `json:"duplicate_of_support_id,omitempty"`
	SpamRiskScore                int                    `json:"spam_risk_score"`
	Attachments                  []string               `json:"attachments"`
	Replies                      []SupportReplyResponse `json:"replies"`
	ResolutionNote               string                 `json:"resolution_note,omitempty"`
	AssigneeUserID               *string                `json:"assignee_user_id,omitempty"`
	ReviewedByUserID             *string                `json:"reviewed_by_user_id,omitempty"`
	ResolvedAt                   *time.Time             `json:"resolved_at,omitempty"`
	ClosedAt                     *time.Time             `json:"closed_at,omitempty"`
	ModerationHandoffRequestedAt *time.Time             `json:"moderation_handoff_requested_at,omitempty"`
	LinkedModerationCaseID       *string                `json:"linked_moderation_case_id,omitempty"`
	CreatedAt                    time.Time              `json:"created_at"`
	UpdatedAt                    time.Time              `json:"updated_at"`
}

// ListReviewQueueResponse wraps support queue payload.
type ListReviewQueueResponse struct {
	Items []SupportListItemResponse `json:"items"`
	Count int                       `json:"count"`
}

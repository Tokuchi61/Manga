package dto

import "time"

// ListCommentsRequest defines target-based root comment listing controls.
type ListCommentsRequest struct {
	TargetType    string `json:"-" validate:"required,oneof=manga chapter"`
	TargetID      string `json:"-" validate:"required,uuid4"`
	SortBy        string `json:"-" validate:"omitempty,oneof=newest oldest popular"`
	Limit         int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset        int    `json:"-" validate:"omitempty,min=0"`
	IncludeHidden bool   `json:"-"`
}

// GetCommentDetailRequest resolves a single comment detail.
type GetCommentDetailRequest struct {
	CommentID     string `json:"-" validate:"required,uuid4"`
	IncludeHidden bool   `json:"-"`
}

// GetCommentThreadRequest resolves comment thread responses.
type GetCommentThreadRequest struct {
	CommentID     string `json:"-" validate:"required,uuid4"`
	SortBy        string `json:"-" validate:"omitempty,oneof=newest oldest popular"`
	Limit         int    `json:"-" validate:"omitempty,min=1,max=500"`
	Offset        int    `json:"-" validate:"omitempty,min=0"`
	IncludeHidden bool   `json:"-"`
}

// CommentListItemResponse is the list payload item.
type CommentListItemResponse struct {
	CommentID        string     `json:"comment_id"`
	TargetType       string     `json:"target_type"`
	TargetID         string     `json:"target_id"`
	AuthorUserID     string     `json:"author_user_id"`
	ParentCommentID  *string    `json:"parent_comment_id,omitempty"`
	RootCommentID    *string    `json:"root_comment_id,omitempty"`
	Depth            int        `json:"depth"`
	Content          string     `json:"content"`
	Spoiler          bool       `json:"spoiler"`
	Pinned           bool       `json:"pinned"`
	Locked           bool       `json:"locked"`
	ModerationStatus string     `json:"moderation_status"`
	Shadowbanned     bool       `json:"shadowbanned"`
	Deleted          bool       `json:"deleted"`
	ReplyCount       int64      `json:"reply_count"`
	LikeCount        int64      `json:"like_count"`
	EditCount        int        `json:"edit_count"`
	EditedAt         *time.Time `json:"edited_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ListCommentsResponse wraps listing response payload.
type ListCommentsResponse struct {
	Items []CommentListItemResponse `json:"items"`
	Count int                       `json:"count"`
}

// CommentDetailResponse is the detail payload.
type CommentDetailResponse struct {
	CommentID        string     `json:"comment_id"`
	TargetType       string     `json:"target_type"`
	TargetID         string     `json:"target_id"`
	AuthorUserID     string     `json:"author_user_id"`
	ParentCommentID  *string    `json:"parent_comment_id,omitempty"`
	RootCommentID    *string    `json:"root_comment_id,omitempty"`
	Depth            int        `json:"depth"`
	Content          string     `json:"content"`
	Attachments      []string   `json:"attachments"`
	Spoiler          bool       `json:"spoiler"`
	Pinned           bool       `json:"pinned"`
	Locked           bool       `json:"locked"`
	ModerationStatus string     `json:"moderation_status"`
	Shadowbanned     bool       `json:"shadowbanned"`
	SpamRiskScore    int        `json:"spam_risk_score"`
	Deleted          bool       `json:"deleted"`
	DeleteReason     string     `json:"delete_reason,omitempty"`
	ReplyCount       int64      `json:"reply_count"`
	LikeCount        int64      `json:"like_count"`
	EditCount        int        `json:"edit_count"`
	EditedAt         *time.Time `json:"edited_at,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// CommentThreadResponse returns root + replies thread payload.
type CommentThreadResponse struct {
	Root    CommentDetailResponse     `json:"root"`
	Replies []CommentListItemResponse `json:"replies"`
	Count   int                       `json:"count"`
}

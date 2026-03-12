package dto

import "time"

// CreateWallPostRequest creates own wall post.
type CreateWallPostRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	Body          string `json:"body" validate:"required,max=2000"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// CreateWallReplyRequest creates social-native reply under post.
type CreateWallReplyRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	PostID        string `json:"-" validate:"required,uuid4"`
	Body          string `json:"body" validate:"required,max=2000"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// ListWallPostsRequest resolves wall feed for owner.
type ListWallPostsRequest struct {
	ActorUserID    string `json:"-" validate:"required,uuid4"`
	OwnerUserID    string `json:"-" validate:"omitempty,uuid4"`
	IncludeReplies bool   `json:"-"`
	Limit          int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset         int    `json:"-" validate:"omitempty,min=0"`
	SortBy         string `json:"-" validate:"omitempty,oneof=newest oldest"`
}

// WallReplyItemResponse is wall reply payload.
type WallReplyItemResponse struct {
	ReplyID     string    `json:"reply_id"`
	PostID      string    `json:"post_id"`
	OwnerUserID string    `json:"owner_user_id"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WallPostItemResponse is wall post payload.
type WallPostItemResponse struct {
	PostID      string                  `json:"post_id"`
	OwnerUserID string                  `json:"owner_user_id"`
	Body        string                  `json:"body"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
	Replies     []WallReplyItemResponse `json:"replies,omitempty"`
}

// WallWriteResponse returns wall write result.
type WallWriteResponse struct {
	PostID      string    `json:"post_id,omitempty"`
	ReplyID     string    `json:"reply_id,omitempty"`
	OwnerUserID string    `json:"owner_user_id"`
	Created     bool      `json:"created"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListWallPostsResponse wraps wall feed payload.
type ListWallPostsResponse struct {
	Items []WallPostItemResponse `json:"items"`
	Count int                    `json:"count"`
}

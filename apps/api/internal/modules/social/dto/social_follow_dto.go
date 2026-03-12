package dto

import "time"

// FollowUserRequest follows target user.
type FollowUserRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	TargetUserID string `json:"-" validate:"required,uuid4"`
	RequestID    string `json:"request_id,omitempty" validate:"omitempty,max=128"`
}

// UnfollowUserRequest unfollows target user.
type UnfollowUserRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	TargetUserID string `json:"-" validate:"required,uuid4"`
}

// FollowUserResponse returns follow operation result.
type FollowUserResponse struct {
	FollowerUserID string    `json:"follower_user_id"`
	FolloweeUserID string    `json:"followee_user_id"`
	Following      bool      `json:"following"`
	Created        bool      `json:"created"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ListFollowersRequest resolves followers of actor.
type ListFollowersRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
}

// ListFollowingRequest resolves following list of actor.
type ListFollowingRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
}

// FollowItemResponse is follow list payload.
type FollowItemResponse struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListFollowsResponse wraps follow list payload.
type ListFollowsResponse struct {
	Items []FollowItemResponse `json:"items"`
	Count int                  `json:"count"`
}

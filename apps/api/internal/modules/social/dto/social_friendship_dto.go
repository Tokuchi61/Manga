package dto

import "time"

// CreateFriendRequestRequest sends friend request to target user.
type CreateFriendRequestRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	TargetUserID string `json:"target_user_id" validate:"required,uuid4"`
	RequestID    string `json:"request_id,omitempty" validate:"omitempty,max=128"`
}

// CreateFriendRequestResponse returns friend request operation result.
type CreateFriendRequestResponse struct {
	RequesterUserID string    `json:"requester_user_id"`
	TargetUserID    string    `json:"target_user_id"`
	Status          string    `json:"status"`
	Created         bool      `json:"created"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RespondFriendRequestRequest accepts or rejects incoming request.
type RespondFriendRequestRequest struct {
	ActorUserID     string `json:"-" validate:"required,uuid4"`
	RequesterUserID string `json:"-" validate:"required,uuid4"`
	Action          string `json:"action" validate:"required,oneof=accept reject"`
}

// RemoveFriendRequest removes existing friendship relation.
type RemoveFriendRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	TargetUserID string `json:"-" validate:"required,uuid4"`
}

// ListFriendRequestsRequest resolves incoming/outgoing pending requests.
type ListFriendRequestsRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
	Direction   string `json:"-" validate:"omitempty,oneof=incoming outgoing"`
}

// ListFriendsRequest resolves accepted friend list.
type ListFriendsRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
}

// FriendRequestItemResponse is pending request payload.
type FriendRequestItemResponse struct {
	RequesterUserID string    `json:"requester_user_id"`
	TargetUserID    string    `json:"target_user_id"`
	Status          string    `json:"status"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ListFriendRequestsResponse wraps pending request list.
type ListFriendRequestsResponse struct {
	Items []FriendRequestItemResponse `json:"items"`
	Count int                         `json:"count"`
}

// FriendItemResponse is accepted friendship payload.
type FriendItemResponse struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListFriendsResponse wraps accepted friend list.
type ListFriendsResponse struct {
	Items []FriendItemResponse `json:"items"`
	Count int                  `json:"count"`
}

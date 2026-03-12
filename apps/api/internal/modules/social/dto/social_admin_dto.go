package dto

import "time"

// UpdateFriendshipStateRequest updates friendship runtime state.
type UpdateFriendshipStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateFollowStateRequest updates follow runtime state.
type UpdateFollowStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateWallStateRequest updates wall runtime state.
type UpdateWallStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateMessagingStateRequest updates messaging runtime state.
type UpdateMessagingStateRequest struct {
	Enabled bool `json:"enabled"`
}

// RuntimeConfigResponse is social runtime control payload.
type RuntimeConfigResponse struct {
	FriendshipEnabled bool      `json:"friendship_enabled"`
	FollowEnabled     bool      `json:"follow_enabled"`
	WallEnabled       bool      `json:"wall_enabled"`
	MessagingEnabled  bool      `json:"messaging_enabled"`
	UpdatedAt         time.Time `json:"updated_at"`
}

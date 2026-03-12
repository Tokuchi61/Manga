package entity

import "time"

// FriendshipRequestStatus defines lifecycle of a friend request.
type FriendshipRequestStatus string

const (
	FriendshipRequestPending  FriendshipRequestStatus = "pending"
	FriendshipRequestRejected FriendshipRequestStatus = "rejected"
)

// Friendship models an accepted friendship relationship.
type Friendship struct {
	UserAID   string
	UserBID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FriendshipRequest models directional friendship request.
type FriendshipRequest struct {
	RequesterUserID string
	TargetUserID    string
	Status          FriendshipRequestStatus
	RequestID       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FollowRelation models directional follow relationship.
type FollowRelation struct {
	FollowerUserID string
	FolloweeUserID string
	RequestID      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// WallPost models social-native wall post.
type WallPost struct {
	ID          string
	OwnerUserID string
	Body        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WallReply models social-native reply under wall post.
type WallReply struct {
	ID          string
	PostID      string
	OwnerUserID string
	Body        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MessageThread models direct messaging thread.
type MessageThread struct {
	ID            string
	UserAID       string
	UserBID       string
	LastMessageID string
	UnreadByA     int
	UnreadByB     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Message models direct thread message.
type Message struct {
	ID            string
	ThreadID      string
	SenderUserID  string
	Body          string
	RequestID     string
	CorrelationID string
	CreatedAt     time.Time
}

// RuntimeConfig stores stage-14 social runtime controls.
type RuntimeConfig struct {
	FriendshipEnabled bool
	FollowEnabled     bool
	WallEnabled       bool
	MessagingEnabled  bool
	UpdatedAt         time.Time
}

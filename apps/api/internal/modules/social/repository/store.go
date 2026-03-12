package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
)

var (
	ErrNotFound = errors.New("social_repository_not_found")
	ErrConflict = errors.New("social_repository_conflict")
)

// RelationEntry holds actor-target relation tuple.
type RelationEntry struct {
	ActorUserID  string
	TargetUserID string
	UpdatedAt    time.Time
}

// Store defines social persistence boundary.
type Store interface {
	CreateFriendRequest(ctx context.Context, request entity.FriendshipRequest, dedupKey string) (entity.FriendshipRequest, bool, error)
	GetFriendRequest(ctx context.Context, requesterUserID string, targetUserID string) (entity.FriendshipRequest, error)
	DeleteFriendRequest(ctx context.Context, requesterUserID string, targetUserID string) error
	ListFriendRequests(ctx context.Context, userID string, direction string) ([]entity.FriendshipRequest, error)
	UpsertFriendship(ctx context.Context, friendship entity.Friendship) error
	DeleteFriendship(ctx context.Context, userAID string, userBID string) (bool, error)
	ListFriends(ctx context.Context, userID string) ([]entity.Friendship, error)
	AreFriends(ctx context.Context, userAID string, userBID string) (bool, error)

	UpsertFollow(ctx context.Context, follow entity.FollowRelation, dedupKey string) (entity.FollowRelation, bool, error)
	DeleteFollow(ctx context.Context, followerUserID string, followeeUserID string) (bool, error)
	ListFollowers(ctx context.Context, userID string) ([]entity.FollowRelation, error)
	ListFollowing(ctx context.Context, userID string) ([]entity.FollowRelation, error)

	CreateWallPost(ctx context.Context, post entity.WallPost, dedupKey string) (entity.WallPost, bool, error)
	GetWallPostByID(ctx context.Context, postID string) (entity.WallPost, error)
	CreateWallReply(ctx context.Context, reply entity.WallReply, dedupKey string) (entity.WallReply, bool, error)
	ListWallPosts(ctx context.Context, ownerUserID string, sortBy string, limit int, offset int) ([]entity.WallPost, error)
	ListWallRepliesByPostID(ctx context.Context, postID string, sortBy string) ([]entity.WallReply, error)

	OpenThread(ctx context.Context, thread entity.MessageThread) (entity.MessageThread, bool, error)
	GetThreadByID(ctx context.Context, threadID string) (entity.MessageThread, error)
	UpdateThread(ctx context.Context, thread entity.MessageThread) error
	ListThreadsByUser(ctx context.Context, userID string) ([]entity.MessageThread, error)
	CreateMessage(ctx context.Context, message entity.Message, dedupKey string) (entity.Message, bool, error)
	ListMessagesByThreadID(ctx context.Context, threadID string, sortBy string, limit int, offset int) ([]entity.Message, error)

	SetBlock(ctx context.Context, actorUserID string, targetUserID string, enabled bool, updatedAt time.Time) error
	SetMute(ctx context.Context, actorUserID string, targetUserID string, enabled bool, updatedAt time.Time) error
	SetRestrict(ctx context.Context, actorUserID string, targetUserID string, enabled bool, updatedAt time.Time) error
	ListBlocked(ctx context.Context, actorUserID string) ([]RelationEntry, error)
	ListMuted(ctx context.Context, actorUserID string) ([]RelationEntry, error)
	ListRestricted(ctx context.Context, actorUserID string) ([]RelationEntry, error)
	IsBlockedEither(ctx context.Context, userAID string, userBID string) (bool, error)
	IsRestrictedEither(ctx context.Context, userAID string, userBID string) (bool, error)
	IsMutedEither(ctx context.Context, userAID string, userBID string) (bool, error)

	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error
}

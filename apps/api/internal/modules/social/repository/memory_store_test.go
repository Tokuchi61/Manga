package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreFriendshipAndFollowDedup(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 13, 1, 0, 0, 0, time.UTC)
	userA := uuid.NewString()
	userB := uuid.NewString()

	request, created, err := store.CreateFriendRequest(context.Background(), entity.FriendshipRequest{
		RequesterUserID: userA,
		TargetUserID:    userB,
		Status:          entity.FriendshipRequestPending,
		RequestID:       "req-social-friend-1",
		CreatedAt:       now,
		UpdatedAt:       now,
	}, "req-social-friend-1")
	require.NoError(t, err)
	require.True(t, created)
	require.Equal(t, userA, request.RequesterUserID)

	again, created, err := store.CreateFriendRequest(context.Background(), entity.FriendshipRequest{
		RequesterUserID: userA,
		TargetUserID:    userB,
		Status:          entity.FriendshipRequestPending,
		RequestID:       "req-social-friend-1",
		CreatedAt:       now,
		UpdatedAt:       now,
	}, "req-social-friend-1")
	require.NoError(t, err)
	require.False(t, created)
	require.Equal(t, request.RequesterUserID, again.RequesterUserID)

	follow, followCreated, err := store.UpsertFollow(context.Background(), entity.FollowRelation{
		FollowerUserID: userA,
		FolloweeUserID: userB,
		RequestID:      "req-social-follow-1",
		CreatedAt:      now,
		UpdatedAt:      now,
	}, "req-social-follow-1")
	require.NoError(t, err)
	require.True(t, followCreated)
	require.Equal(t, userA, follow.FollowerUserID)

	followAgain, followCreated, err := store.UpsertFollow(context.Background(), entity.FollowRelation{
		FollowerUserID: userA,
		FolloweeUserID: userB,
		RequestID:      "req-social-follow-1",
		CreatedAt:      now,
		UpdatedAt:      now,
	}, "req-social-follow-1")
	require.NoError(t, err)
	require.False(t, followCreated)
	require.Equal(t, follow.FollowerUserID, followAgain.FollowerUserID)
}

func TestMemoryStoreThreadAndMessageDedupFlow(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 13, 1, 5, 0, 0, time.UTC)
	userA := uuid.NewString()
	userB := uuid.NewString()
	threadID := uuid.NewString()

	thread, created, err := store.OpenThread(context.Background(), entity.MessageThread{
		ID:        threadID,
		UserAID:   userA,
		UserBID:   userB,
		CreatedAt: now,
		UpdatedAt: now,
	})
	require.NoError(t, err)
	require.True(t, created)
	require.Equal(t, threadID, thread.ID)

	again, created, err := store.OpenThread(context.Background(), entity.MessageThread{
		ID:        uuid.NewString(),
		UserAID:   userB,
		UserBID:   userA,
		CreatedAt: now,
		UpdatedAt: now,
	})
	require.NoError(t, err)
	require.False(t, created)
	require.Equal(t, threadID, again.ID)

	messageID := uuid.NewString()
	message, created, err := store.CreateMessage(context.Background(), entity.Message{
		ID:           messageID,
		ThreadID:     threadID,
		SenderUserID: userA,
		Body:         "hello world",
		RequestID:    "req-social-msg-1",
		CreatedAt:    now,
	}, "req-social-msg-1")
	require.NoError(t, err)
	require.True(t, created)
	require.Equal(t, messageID, message.ID)

	againMessage, created, err := store.CreateMessage(context.Background(), entity.Message{
		ID:           uuid.NewString(),
		ThreadID:     threadID,
		SenderUserID: userA,
		Body:         "hello world",
		RequestID:    "req-social-msg-1",
		CreatedAt:    now,
	}, "req-social-msg-1")
	require.NoError(t, err)
	require.False(t, created)
	require.Equal(t, messageID, againMessage.ID)

	messages, err := store.ListMessagesByThreadID(context.Background(), threadID, "newest", 10, 0)
	require.NoError(t, err)
	require.Len(t, messages, 1)
	require.Equal(t, messageID, messages[0].ID)
}

func TestMemoryStoreRelationsAndRuntimeConfig(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 13, 1, 10, 0, 0, time.UTC)
	userA := uuid.NewString()
	userB := uuid.NewString()

	require.NoError(t, store.SetBlock(context.Background(), userA, userB, true, now))
	require.NoError(t, store.SetMute(context.Background(), userA, userB, true, now))
	require.NoError(t, store.SetRestrict(context.Background(), userB, userA, true, now))

	blocked, err := store.IsBlockedEither(context.Background(), userA, userB)
	require.NoError(t, err)
	require.True(t, blocked)

	muted, err := store.IsMutedEither(context.Background(), userA, userB)
	require.NoError(t, err)
	require.True(t, muted)

	restricted, err := store.IsRestrictedEither(context.Background(), userA, userB)
	require.NoError(t, err)
	require.True(t, restricted)

	blockedList, err := store.ListBlocked(context.Background(), userA)
	require.NoError(t, err)
	require.Len(t, blockedList, 1)
	require.Equal(t, userB, blockedList[0].TargetUserID)

	cfg, err := store.GetRuntimeConfig(context.Background())
	require.NoError(t, err)
	require.True(t, cfg.WallEnabled)
	cfg.WallEnabled = false
	cfg.UpdatedAt = now
	require.NoError(t, store.UpdateRuntimeConfig(context.Background(), cfg))

	persisted, err := store.GetRuntimeConfig(context.Background())
	require.NoError(t, err)
	require.False(t, persisted.WallEnabled)
}

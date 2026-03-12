package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/dto"
	socialrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/social/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSocialServiceFriendshipFollowWallMessageFlow(t *testing.T) {
	store := socialrepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 13, 1, 30, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	actorID := uuid.NewString()
	targetID := uuid.NewString()

	friendRequest, err := svc.CreateFriendRequest(context.Background(), dto.CreateFriendRequestRequest{
		ActorUserID:  actorID,
		TargetUserID: targetID,
		RequestID:    "req-social-friend-1",
	})
	require.NoError(t, err)
	require.True(t, friendRequest.Created)
	require.Equal(t, "pending", friendRequest.Status)

	acceptResult, err := svc.RespondFriendRequest(context.Background(), dto.RespondFriendRequestRequest{
		ActorUserID:     targetID,
		RequesterUserID: actorID,
		Action:          "accept",
	})
	require.NoError(t, err)
	require.Equal(t, "accepted", acceptResult.Status)

	friends, err := svc.ListFriends(context.Background(), dto.ListFriendsRequest{ActorUserID: actorID})
	require.NoError(t, err)
	require.Equal(t, 1, friends.Count)
	require.Equal(t, targetID, friends.Items[0].UserID)

	follow, err := svc.FollowUser(context.Background(), dto.FollowUserRequest{
		ActorUserID:  actorID,
		TargetUserID: targetID,
		RequestID:    "req-social-follow-1",
	})
	require.NoError(t, err)
	require.True(t, follow.Following)

	post, err := svc.CreateWallPost(context.Background(), dto.CreateWallPostRequest{
		ActorUserID: actorID,
		Body:        "hello wall",
		RequestID:   "req-social-wall-1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, post.PostID)

	reply, err := svc.CreateWallReply(context.Background(), dto.CreateWallReplyRequest{
		ActorUserID: targetID,
		PostID:      post.PostID,
		Body:        "hi there",
		RequestID:   "req-social-reply-1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, reply.ReplyID)

	wall, err := svc.ListWallPosts(context.Background(), dto.ListWallPostsRequest{
		ActorUserID:    actorID,
		OwnerUserID:    actorID,
		IncludeReplies: true,
	})
	require.NoError(t, err)
	require.Equal(t, 1, wall.Count)
	require.Len(t, wall.Items[0].Replies, 1)

	thread, err := svc.OpenThread(context.Background(), dto.OpenThreadRequest{ActorUserID: actorID, TargetUserID: targetID})
	require.NoError(t, err)
	require.NotEmpty(t, thread.ThreadID)

	send, err := svc.SendMessage(context.Background(), dto.SendMessageRequest{
		ActorUserID:   actorID,
		ThreadID:      thread.ThreadID,
		Body:          "hello dm",
		RequestID:     "req-social-message-1",
		CorrelationID: "corr-social-message-1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, send.MessageID)

	threads, err := svc.ListThreads(context.Background(), dto.ListThreadsRequest{ActorUserID: actorID})
	require.NoError(t, err)
	require.Equal(t, 1, threads.Count)
	require.Equal(t, 0, threads.Items[0].UnreadCount)

	messages, err := svc.ListThreadMessages(context.Background(), dto.ListThreadMessagesRequest{
		ActorUserID: targetID,
		ThreadID:    thread.ThreadID,
	})
	require.NoError(t, err)
	require.Equal(t, 1, messages.Count)

	_, err = svc.MarkThreadRead(context.Background(), dto.MarkThreadReadRequest{ActorUserID: targetID, ThreadID: thread.ThreadID})
	require.NoError(t, err)
}

func TestSocialServiceBlockRemovesRelationsAndBlocksMessaging(t *testing.T) {
	store := socialrepository.NewMemoryStore()
	svc := New(store, validation.New())

	actorID := uuid.NewString()
	targetID := uuid.NewString()

	_, err := svc.CreateFriendRequest(context.Background(), dto.CreateFriendRequestRequest{
		ActorUserID:  actorID,
		TargetUserID: targetID,
		RequestID:    "req-social-friend-2",
	})
	require.NoError(t, err)

	_, err = svc.RespondFriendRequest(context.Background(), dto.RespondFriendRequestRequest{
		ActorUserID:     targetID,
		RequesterUserID: actorID,
		Action:          "accept",
	})
	require.NoError(t, err)

	_, err = svc.FollowUser(context.Background(), dto.FollowUserRequest{
		ActorUserID:  actorID,
		TargetUserID: targetID,
		RequestID:    "req-social-follow-2",
	})
	require.NoError(t, err)

	_, err = svc.UpdateBlock(context.Background(), dto.UpdateBlockRequest{ActorUserID: actorID, TargetUserID: targetID, Enabled: true})
	require.NoError(t, err)

	friends, err := svc.ListFriends(context.Background(), dto.ListFriendsRequest{ActorUserID: actorID})
	require.NoError(t, err)
	require.Equal(t, 0, friends.Count)

	following, err := svc.ListFollowing(context.Background(), dto.ListFollowingRequest{ActorUserID: actorID})
	require.NoError(t, err)
	require.Equal(t, 0, following.Count)

	_, err = svc.OpenThread(context.Background(), dto.OpenThreadRequest{ActorUserID: actorID, TargetUserID: targetID})
	require.ErrorIs(t, err, ErrForbiddenAction)
}

func TestSocialServiceRuntimeTogglesAffectWrites(t *testing.T) {
	store := socialrepository.NewMemoryStore()
	svc := New(store, validation.New())

	actorID := uuid.NewString()
	targetID := uuid.NewString()

	_, err := svc.UpdateFriendshipState(context.Background(), dto.UpdateFriendshipStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.CreateFriendRequest(context.Background(), dto.CreateFriendRequestRequest{ActorUserID: actorID, TargetUserID: targetID})
	require.True(t, errors.Is(err, ErrFriendshipDisabled))

	_, err = svc.UpdateFollowState(context.Background(), dto.UpdateFollowStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.FollowUser(context.Background(), dto.FollowUserRequest{ActorUserID: actorID, TargetUserID: targetID})
	require.True(t, errors.Is(err, ErrFollowDisabled))

	_, err = svc.UpdateWallState(context.Background(), dto.UpdateWallStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.ListWallPosts(context.Background(), dto.ListWallPostsRequest{ActorUserID: actorID})
	require.True(t, errors.Is(err, ErrWallDisabled))

	_, err = svc.UpdateMessagingState(context.Background(), dto.UpdateMessagingStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.OpenThread(context.Background(), dto.OpenThreadRequest{ActorUserID: actorID, TargetUserID: targetID})
	require.True(t, errors.Is(err, ErrMessagingDisabled))
}

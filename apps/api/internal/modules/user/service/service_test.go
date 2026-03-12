package service

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func newTestService(nowRef *time.Time) *UserService {
	svc := New(userrepository.NewMemoryStore(), validation.New())
	svc.now = func() time.Time { return nowRef.UTC() }
	return svc
}

func TestCreateAndPublicPrivateProfileFlow(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	createRes, err := svc.CreateUser(ctx, dto.CreateUserRequest{
		CredentialID: uuid.NewString(),
		Username:     "reader_one",
		DisplayName:  "Reader One",
	})
	require.NoError(t, err)
	require.NotEmpty(t, createRes.UserID)

	publicRes, err := svc.GetPublicProfile(ctx, dto.GetPublicProfileRequest{UserID: createRes.UserID})
	require.NoError(t, err)
	require.Equal(t, "reader_one", publicRes.Username)

	_, err = svc.UpdateProfileVisibility(ctx, dto.UpdateProfileVisibilityRequest{
		UserID:            createRes.UserID,
		ViewerID:          createRes.UserID,
		ProfileVisibility: "private",
	})
	require.NoError(t, err)

	_, err = svc.GetPublicProfile(ctx, dto.GetPublicProfileRequest{UserID: createRes.UserID})
	require.ErrorIs(t, err, ErrProfileNotVisible)

	ownRes, err := svc.GetOwnProfile(ctx, dto.GetOwnProfileRequest{UserID: createRes.UserID, ViewerID: createRes.UserID})
	require.NoError(t, err)
	require.Equal(t, "private", ownRes.ProfileVisibility)
}

func TestUpdateProfileAndHistoryVisibilityPrecedence(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	createRes, err := svc.CreateUser(ctx, dto.CreateUserRequest{
		CredentialID: uuid.NewString(),
		Username:     "reader_two",
	})
	require.NoError(t, err)

	display := "Reader Updated"
	bio := "Manga lover"
	_, err = svc.UpdateProfile(ctx, dto.UpdateProfileRequest{
		UserID:      createRes.UserID,
		ViewerID:    createRes.UserID,
		DisplayName: &display,
		Bio:         &bio,
	})
	require.NoError(t, err)

	_, err = svc.UpdateHistoryVisibility(ctx, dto.UpdateHistoryVisibilityRequest{
		UserID:                      createRes.UserID,
		ViewerID:                    createRes.UserID,
		HistoryVisibilityPreference: "private",
	})
	require.NoError(t, err)

	ownRes, err := svc.GetOwnProfile(ctx, dto.GetOwnProfileRequest{UserID: createRes.UserID, ViewerID: createRes.UserID})
	require.NoError(t, err)
	require.Equal(t, "Reader Updated", ownRes.DisplayName)
	require.Equal(t, "private", ownRes.HistoryVisibilityPreference)

	require.False(t, historyVisibleForViewer(entity.HistoryVisibilityPrivate, entity.HistoryVisibilityPublic, false))
	require.True(t, historyVisibleForViewer(entity.HistoryVisibilityPublic, entity.HistoryVisibilityPublic, false))
	require.False(t, historyVisibleForViewer(entity.HistoryVisibilityPublic, entity.HistoryVisibilityPrivate, false))
	require.True(t, historyVisibleForViewer(entity.HistoryVisibilityPrivate, entity.HistoryVisibilityPrivate, true))
}

func TestAccountStateEffects(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	createRes, err := svc.CreateUser(ctx, dto.CreateUserRequest{
		CredentialID: uuid.NewString(),
		Username:     "reader_three",
	})
	require.NoError(t, err)

	_, err = svc.UpdateAccountState(ctx, dto.UpdateAccountStateRequest{
		UserID:       createRes.UserID,
		ActorUserID:  createRes.UserID,
		ActorScope:   "self",
		AccountState: "deactivated",
	})
	require.NoError(t, err)

	display := "Should Fail"
	_, err = svc.UpdateProfile(ctx, dto.UpdateProfileRequest{
		UserID:      createRes.UserID,
		ViewerID:    createRes.UserID,
		DisplayName: &display,
	})
	require.ErrorIs(t, err, ErrAccountDeactivated)

	_, err = svc.UpdateAccountState(ctx, dto.UpdateAccountStateRequest{
		UserID:       createRes.UserID,
		ActorScope:   "admin",
		AccountState: "banned",
	})
	require.NoError(t, err)

	_, err = svc.GetPublicProfile(ctx, dto.GetPublicProfileRequest{UserID: createRes.UserID})
	require.ErrorIs(t, err, ErrProfileNotVisible)
}

func TestVIPLifecycle(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	createRes, err := svc.CreateUser(ctx, dto.CreateUserRequest{
		CredentialID: uuid.NewString(),
		Username:     "reader_four",
	})
	require.NoError(t, err)

	endAt := now.Add(30 * 24 * time.Hour)
	vipRes, err := svc.UpdateVIPState(ctx, dto.UpdateVIPStateRequest{
		UserID: createRes.UserID,
		Action: "activate",
		EndsAt: &endAt,
	})
	require.NoError(t, err)
	require.True(t, vipRes.VIPActive)

	vipRes, err = svc.UpdateVIPState(ctx, dto.UpdateVIPStateRequest{
		UserID:       createRes.UserID,
		Action:       "freeze",
		FreezeReason: "system_pause",
	})
	require.NoError(t, err)
	require.True(t, vipRes.VIPFrozen)

	vipRes, err = svc.UpdateVIPState(ctx, dto.UpdateVIPStateRequest{
		UserID: createRes.UserID,
		Action: "resume",
	})
	require.NoError(t, err)
	require.False(t, vipRes.VIPFrozen)

	vipRes, err = svc.UpdateVIPState(ctx, dto.UpdateVIPStateRequest{
		UserID: createRes.UserID,
		Action: "deactivate",
	})
	require.NoError(t, err)
	require.False(t, vipRes.VIPActive)
}

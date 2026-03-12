package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreUserLifecycle(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)

	user := entity.UserAccount{
		ID:                          uuid.NewString(),
		CredentialID:                uuid.NewString(),
		Username:                    "reader_one",
		DisplayName:                 "Reader One",
		ProfileVisibility:           entity.ProfileVisibilityPublic,
		HistoryVisibilityPreference: entity.HistoryVisibilityPrivate,
		AccountState:                entity.AccountStateActive,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}

	require.NoError(t, store.CreateUser(ctx, user))

	loadedByID, err := store.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, normalizeID(user.ID), loadedByID.ID)

	loadedByCredential, err := store.GetUserByCredentialID(ctx, user.CredentialID)
	require.NoError(t, err)
	require.Equal(t, normalizeID(user.ID), loadedByCredential.ID)

	loadedByID.DisplayName = "Updated"
	require.NoError(t, store.UpdateUser(ctx, loadedByID))

	updated, err := store.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.DisplayName)
}

func TestMemoryStoreRejectsDuplicateCredentialOrUsername(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)

	credentialID := uuid.NewString()
	userA := entity.UserAccount{
		ID:                          uuid.NewString(),
		CredentialID:                credentialID,
		Username:                    "reader_one",
		DisplayName:                 "Reader One",
		ProfileVisibility:           entity.ProfileVisibilityPublic,
		HistoryVisibilityPreference: entity.HistoryVisibilityPrivate,
		AccountState:                entity.AccountStateActive,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}
	require.NoError(t, store.CreateUser(ctx, userA))

	userSameCredential := userA
	userSameCredential.ID = uuid.NewString()
	userSameCredential.Username = "reader_two"
	require.ErrorIs(t, store.CreateUser(ctx, userSameCredential), ErrConflict)

	userSameUsername := userA
	userSameUsername.ID = uuid.NewString()
	userSameUsername.CredentialID = uuid.NewString()
	require.ErrorIs(t, store.CreateUser(ctx, userSameUsername), ErrConflict)
}

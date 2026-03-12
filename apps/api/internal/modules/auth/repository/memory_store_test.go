package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreCredentialAndSessionLifecycle(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)

	credential := entity.Credential{
		ID:           uuid.New(),
		Email:        "store@example.com",
		PasswordHash: "hash",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, store.CreateCredential(ctx, credential))

	loadedCredential, err := store.GetCredentialByEmail(ctx, "store@example.com")
	require.NoError(t, err)
	require.Equal(t, credential.ID, loadedCredential.ID)

	session := entity.Session{
		ID:           uuid.New(),
		CredentialID: credential.ID,
		CreatedAt:    now,
		LastSeenAt:   now,
	}
	require.NoError(t, store.CreateSession(ctx, session))

	sessions, err := store.ListSessionsByCredential(ctx, credential.ID)
	require.NoError(t, err)
	require.Len(t, sessions, 1)

	require.NoError(t, store.RevokeSession(ctx, session.ID, now.Add(time.Minute).Unix()))
	revokedSession, err := store.GetSessionByID(ctx, session.ID)
	require.NoError(t, err)
	require.True(t, revokedSession.IsRevoked())
}

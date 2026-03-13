package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreSeasonTierRoundTrip(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	start := now.Add(-2 * time.Hour)
	end := now.Add(24 * time.Hour)

	err := store.UpsertSeasonDefinition(context.Background(), entity.SeasonDefinition{
		SeasonID:  "season_2026_03",
		Title:     "Spring Launch",
		State:     entity.SeasonStateActive,
		StartsAt:  &start,
		EndsAt:    &end,
		CreatedAt: now,
		UpdatedAt: now,
	})
	require.NoError(t, err)

	err = store.UpsertTierDefinition(context.Background(), entity.TierDefinition{
		SeasonID:       "season_2026_03",
		TierNumber:     1,
		Track:          entity.TrackFree,
		RequiredPoints: 100,
		RewardItemID:   "mana_100",
		RewardQuantity: 1,
		Active:         true,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	require.NoError(t, err)

	season, err := store.GetSeasonDefinition(context.Background(), "season_2026_03")
	require.NoError(t, err)
	require.Equal(t, "Spring Launch", season.Title)

	tier, err := store.GetTierDefinition(context.Background(), "season_2026_03", 1, entity.TrackFree)
	require.NoError(t, err)
	require.Equal(t, 100, tier.RequiredPoints)

	resolved, err := store.ResolveCurrentSeason(context.Background(), now)
	require.NoError(t, err)
	require.Equal(t, "season_2026_03", resolved.SeasonID)
}

func TestMemoryStoreProgressDedupAndReset(t *testing.T) {
	store := NewMemoryStore()
	userID := uuid.NewString()
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)

	progress := entity.UserProgress{
		UserID:           userID,
		SeasonID:         "season_2026_03",
		Points:           120,
		PremiumActivated: true,
		ClaimedTiers: map[string]time.Time{
			"free:1": now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.NoError(t, store.UpsertUserProgress(context.Background(), progress))

	readProgress, err := store.GetUserProgress(context.Background(), userID, "season_2026_03")
	require.NoError(t, err)
	require.Equal(t, 120, readProgress.Points)
	require.True(t, readProgress.PremiumActivated)
	require.Contains(t, readProgress.ClaimedTiers, "free:1")

	require.NoError(t, store.PutProgressDedup(context.Background(), "dedup-progress", readProgress))
	require.NoError(t, store.PutClaimDedup(context.Background(), "dedup-claim", readProgress))
	require.NoError(t, store.PutPremiumActivationDedup(context.Background(), "dedup-premium", readProgress))

	_, err = store.GetProgressDedup(context.Background(), "dedup-progress")
	require.NoError(t, err)
	_, err = store.GetClaimDedup(context.Background(), "dedup-claim")
	require.NoError(t, err)
	_, err = store.GetPremiumActivationDedup(context.Background(), "dedup-premium")
	require.NoError(t, err)

	deleted, err := store.DeleteUserProgressByUserSeason(context.Background(), userID, "season_2026_03")
	require.NoError(t, err)
	require.Equal(t, 1, deleted)

	_, err = store.GetUserProgress(context.Background(), userID, "season_2026_03")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestMemoryStoreSnapshotRoundTrip(t *testing.T) {
	store := NewMemoryStore()
	userID := uuid.NewString()
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	start := now.Add(-time.Hour)
	end := now.Add(24 * time.Hour)

	require.NoError(t, store.UpsertSeasonDefinition(context.Background(), entity.SeasonDefinition{
		SeasonID:  "season_2026_03",
		Title:     "Spring Launch",
		State:     entity.SeasonStateActive,
		StartsAt:  &start,
		EndsAt:    &end,
		CreatedAt: now,
		UpdatedAt: now,
	}))

	require.NoError(t, store.UpsertUserProgress(context.Background(), entity.UserProgress{
		UserID:           userID,
		SeasonID:         "season_2026_03",
		Points:           50,
		PremiumActivated: false,
		ClaimedTiers:     map[string]time.Time{},
		CreatedAt:        now,
		UpdatedAt:        now,
	}))

	payload, err := store.Snapshot()
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	restored := NewMemoryStore()
	require.NoError(t, restored.RestoreSnapshot(payload))

	season, err := restored.GetSeasonDefinition(context.Background(), "season_2026_03")
	require.NoError(t, err)
	require.Equal(t, "Spring Launch", season.Title)

	progress, err := restored.GetUserProgress(context.Background(), userID, "season_2026_03")
	require.NoError(t, err)
	require.Equal(t, 50, progress.Points)
}

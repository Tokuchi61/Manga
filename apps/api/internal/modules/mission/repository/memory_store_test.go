package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreMissionDefinitionRoundTrip(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

	err := store.UpsertMissionDefinition(context.Background(), entity.MissionDefinition{
		MissionID:      "daily_read_3",
		Category:       "daily",
		Title:          "Read 3 chapters",
		ObjectiveType:  "chapter_read",
		TargetCount:    3,
		RewardItemID:   "mana_potion",
		RewardQuantity: 1,
		Active:         true,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	require.NoError(t, err)

	definition, err := store.GetMissionDefinition(context.Background(), "daily_read_3")
	require.NoError(t, err)
	require.Equal(t, "daily", definition.Category)
	require.Equal(t, 3, definition.TargetCount)

	items, err := store.ListMissionDefinitions(context.Background(), "daily", true, 10, 0)
	require.NoError(t, err)
	require.Len(t, items, 1)
}

func TestMemoryStoreMissionProgressResetAndDedup(t *testing.T) {
	store := NewMemoryStore()
	userID := uuid.NewString()
	now := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

	progress := entity.MissionProgress{
		UserID:            userID,
		MissionID:         "daily_read_3",
		PeriodKey:         "2026-03-15",
		ProgressCount:     2,
		Completed:         false,
		Claimed:           false,
		LastRequestID:     "req-progress-1",
		LastCorrelationID: "corr-progress-1",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	require.NoError(t, store.UpsertMissionProgress(context.Background(), progress))

	readProgress, err := store.GetMissionProgress(context.Background(), userID, "daily_read_3", "2026-03-15")
	require.NoError(t, err)
	require.Equal(t, 2, readProgress.ProgressCount)

	require.NoError(t, store.PutProgressDedup(context.Background(), "dedup-progress", readProgress))
	require.NoError(t, store.PutClaimDedup(context.Background(), "dedup-claim", readProgress))

	_, err = store.GetProgressDedup(context.Background(), "dedup-progress")
	require.NoError(t, err)
	_, err = store.GetClaimDedup(context.Background(), "dedup-claim")
	require.NoError(t, err)

	deleted, err := store.DeleteMissionProgressByUserMission(context.Background(), userID, "daily_read_3")
	require.NoError(t, err)
	require.Equal(t, 1, deleted)

	_, err = store.GetMissionProgress(context.Background(), userID, "daily_read_3", "2026-03-15")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestMemoryStoreMissionSnapshotRoundTrip(t *testing.T) {
	store := NewMemoryStore()
	userID := uuid.NewString()
	now := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

	require.NoError(t, store.UpsertMissionDefinition(context.Background(), entity.MissionDefinition{
		MissionID:      "weekly_comment_5",
		Category:       "weekly",
		Title:          "Write 5 comments",
		ObjectiveType:  "comment_create",
		TargetCount:    5,
		RewardQuantity: 2,
		Active:         true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}))

	require.NoError(t, store.UpsertMissionProgress(context.Background(), entity.MissionProgress{
		UserID:        userID,
		MissionID:     "weekly_comment_5",
		PeriodKey:     "2026-W11",
		ProgressCount: 3,
		CreatedAt:     now,
		UpdatedAt:     now,
	}))

	snapshotPayload, err := store.Snapshot()
	require.NoError(t, err)
	require.NotEmpty(t, snapshotPayload)

	restored := NewMemoryStore()
	require.NoError(t, restored.RestoreSnapshot(snapshotPayload))

	definition, err := restored.GetMissionDefinition(context.Background(), "weekly_comment_5")
	require.NoError(t, err)
	require.Equal(t, "weekly", definition.Category)

	progress, err := restored.GetMissionProgress(context.Background(), userID, "weekly_comment_5", "2026-W11")
	require.NoError(t, err)
	require.Equal(t, 3, progress.ProgressCount)
}

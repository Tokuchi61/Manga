package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreCheckpointDedupAndContinueReadingList(t *testing.T) {
	store := NewMemoryStore()
	userID := uuid.NewString()
	mangaID := uuid.NewString()
	chapterID := uuid.NewString()
	now := time.Date(2026, 3, 12, 23, 30, 0, 0, time.UTC)

	entry, timelineID, created, err := store.UpsertCheckpoint(context.Background(), entity.Checkpoint{
		UserID:     userID,
		MangaID:    mangaID,
		ChapterID:  chapterID,
		Event:      "chapter.read.checkpoint",
		PageNumber: 6,
		PageCount:  20,
		OccurredAt: now,
	}, "req-1")
	require.NoError(t, err)
	require.True(t, created)
	require.NotEmpty(t, timelineID)
	require.Equal(t, 6, entry.LastPageNumber)
	require.Equal(t, entity.ReadingStatusInProgress, entry.Status)

	again, againTimelineID, created, err := store.UpsertCheckpoint(context.Background(), entity.Checkpoint{
		UserID:     userID,
		MangaID:    mangaID,
		ChapterID:  chapterID,
		Event:      "chapter.read.checkpoint",
		PageNumber: 6,
		PageCount:  20,
		OccurredAt: now,
	}, "req-1")
	require.NoError(t, err)
	require.False(t, created)
	require.Equal(t, entry.ID, again.ID)
	require.Equal(t, timelineID, againTimelineID)

	items, err := store.ListContinueReading(context.Background(), ContinueReadingQuery{UserID: userID, Limit: 10})
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, entry.ID, items[0].ID)
}

func TestMemoryStoreLibraryFiltersAndRuntimeConfig(t *testing.T) {
	store := NewMemoryStore()
	userID := uuid.NewString()
	mangaID := uuid.NewString()
	now := time.Date(2026, 3, 12, 23, 40, 0, 0, time.UTC)

	_, err := store.UpsertLibraryEntry(context.Background(), entity.LibraryEntry{
		ID:             uuid.NewString(),
		UserID:         userID,
		MangaID:        mangaID,
		Status:         entity.ReadingStatusCompleted,
		Bookmarked:     true,
		Favorited:      true,
		SharePublic:    true,
		LastReadAt:     now,
		CreatedAt:      now,
		UpdatedAt:      now,
		LastChapterID:  uuid.NewString(),
		PageCount:      20,
		LastPageNumber: 20,
	})
	require.NoError(t, err)

	bookmarked := true
	items, err := store.ListLibrary(context.Background(), LibraryQuery{UserID: userID, Bookmarked: &bookmarked, SharedOnly: true, Limit: 10})
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.True(t, items[0].SharePublic)

	runtimeCfg, err := store.GetRuntimeConfig(context.Background())
	require.NoError(t, err)
	require.True(t, runtimeCfg.LibraryEnabled)
	runtimeCfg.BookmarkWriteEnabled = false
	require.NoError(t, store.UpdateRuntimeConfig(context.Background(), runtimeCfg))

	persisted, err := store.GetRuntimeConfig(context.Background())
	require.NoError(t, err)
	require.False(t, persisted.BookmarkWriteEnabled)
}

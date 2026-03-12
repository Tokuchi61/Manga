package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreCreateConflictBySlugAndSequence(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 15, 0, 0, 0, time.UTC)
	mangaID := uuid.NewString()

	first := entity.Chapter{
		ID:                        uuid.NewString(),
		MangaID:                   mangaID,
		Slug:                      "chapter-1",
		Title:                     "Chapter 1",
		SequenceNo:                1,
		DisplayNumber:             "1",
		PublishState:              entity.PublishStatePublished,
		ReadAccessLevel:           entity.ReadAccessAuthenticated,
		EarlyAccessFallbackAccess: entity.ReadAccessAuthenticated,
		MediaHealthStatus:         entity.MediaHealthHealthy,
		IntegrityStatus:           entity.IntegrityUnknown,
		PageCount:                 1,
		Pages:                     []entity.ChapterPage{{PageNumber: 1, MediaURL: "https://cdn.example.com/1.jpg", Width: 1200, Height: 1800, CDNHealthy: true, UpdatedAt: now}},
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}
	require.NoError(t, store.CreateChapter(context.Background(), first))

	duplicateSlug := first
	duplicateSlug.ID = uuid.NewString()
	duplicateSlug.SequenceNo = 2
	err := store.CreateChapter(context.Background(), duplicateSlug)
	require.ErrorIs(t, err, ErrConflict)

	duplicateSequence := first
	duplicateSequence.ID = uuid.NewString()
	duplicateSequence.Slug = "chapter-2"
	err = store.CreateChapter(context.Background(), duplicateSequence)
	require.ErrorIs(t, err, ErrConflict)
}

func TestMemoryStoreListByMangaSortAndPagination(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 15, 0, 0, 0, time.UTC)
	mangaID := uuid.NewString()

	for i := 1; i <= 3; i++ {
		chapter := entity.Chapter{
			ID:                        uuid.NewString(),
			MangaID:                   mangaID,
			Slug:                      fmt.Sprintf("chapter-%d", i),
			Title:                     "Chapter",
			SequenceNo:                i,
			DisplayNumber:             "",
			PublishState:              entity.PublishStatePublished,
			ReadAccessLevel:           entity.ReadAccessAuthenticated,
			EarlyAccessFallbackAccess: entity.ReadAccessAuthenticated,
			MediaHealthStatus:         entity.MediaHealthHealthy,
			IntegrityStatus:           entity.IntegrityUnknown,
			PageCount:                 1,
			Pages:                     []entity.ChapterPage{{PageNumber: 1, MediaURL: "https://cdn.example.com/a.jpg", Width: 1200, Height: 1800, CDNHealthy: true, UpdatedAt: now}},
			CreatedAt:                 now.Add(time.Duration(i) * time.Minute),
			UpdatedAt:                 now.Add(time.Duration(i) * time.Minute),
		}
		require.NoError(t, store.CreateChapter(context.Background(), chapter))
	}

	items, err := store.ListChaptersByManga(context.Background(), ListQuery{MangaID: mangaID, SortBy: "sequence_oldest", Limit: 2, Offset: 1})
	require.NoError(t, err)
	require.Len(t, items, 2)
	require.Equal(t, 2, items[0].SequenceNo)
	require.Equal(t, 3, items[1].SequenceNo)
}

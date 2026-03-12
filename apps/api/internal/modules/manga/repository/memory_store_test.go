package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreCreateConflictAndGetBySlug(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 12, 0, 0, 0, time.UTC)

	item := entity.Manga{
		ID:             uuid.NewString(),
		Slug:           "tower-of-dawn",
		Title:          "Tower of Dawn",
		Summary:        "Summary",
		ShortSummary:   "Summary",
		Genres:         []string{"action"},
		PublishState:   entity.PublishStatePublished,
		Visibility:     entity.VisibilityPublic,
		ContentVersion: 1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, store.CreateManga(context.Background(), item))

	item2 := item
	item2.ID = uuid.NewString()
	item2.Title = "Tower of Dawn 2"
	err := store.CreateManga(context.Background(), item2)
	require.ErrorIs(t, err, ErrConflict)

	resolved, err := store.GetMangaBySlug(context.Background(), "tower-of-dawn")
	require.NoError(t, err)
	require.Equal(t, item.ID, resolved.ID)
}

func TestMemoryStoreListFiltersAndSort(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 12, 0, 0, 0, time.UTC)

	items := []entity.Manga{
		{
			ID:             uuid.NewString(),
			Slug:           "alpha",
			Title:          "Alpha",
			Summary:        "A summary",
			ShortSummary:   "A short",
			Genres:         []string{"action"},
			Tags:           []string{"fantasy"},
			PublishState:   entity.PublishStatePublished,
			Visibility:     entity.VisibilityPublic,
			ViewCount:      10,
			ContentVersion: 1,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.NewString(),
			Slug:           "beta",
			Title:          "Beta",
			Summary:        "B summary",
			ShortSummary:   "B short",
			Genres:         []string{"drama"},
			Tags:           []string{"slice_of_life"},
			PublishState:   entity.PublishStatePublished,
			Visibility:     entity.VisibilityPublic,
			ViewCount:      50,
			ContentVersion: 1,
			CreatedAt:      now.Add(1 * time.Hour),
			UpdatedAt:      now.Add(1 * time.Hour),
		},
	}

	for _, item := range items {
		require.NoError(t, store.CreateManga(context.Background(), item))
	}

	byGenre, err := store.ListManga(context.Background(), ListQuery{Genre: "action", IncludeUnpublished: false, IncludeHidden: false, Limit: 10})
	require.NoError(t, err)
	require.Len(t, byGenre, 1)
	require.Equal(t, "alpha", byGenre[0].Slug)

	byPopular, err := store.ListManga(context.Background(), ListQuery{SortBy: "popular", IncludeUnpublished: false, IncludeHidden: false, Limit: 10})
	require.NoError(t, err)
	require.Len(t, byPopular, 2)
	require.Equal(t, "beta", byPopular[0].Slug)
}

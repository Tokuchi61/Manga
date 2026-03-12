package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func strPtr(v string) *string {
	return &v
}

func TestMemoryStoreCreateConflictByID(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 19, 0, 0, 0, time.UTC)

	comment := entity.Comment{
		ID:               uuid.NewString(),
		TargetType:       entity.TargetTypeManga,
		TargetID:         uuid.NewString(),
		AuthorUserID:     uuid.NewString(),
		Depth:            0,
		Content:          "First",
		SanitizedContent: "First",
		ModerationStatus: entity.ModerationStatusVisible,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	require.NoError(t, store.CreateComment(context.Background(), comment))

	err := store.CreateComment(context.Background(), comment)
	require.ErrorIs(t, err, ErrConflict)
}

func TestMemoryStoreListAndThreadFilters(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 19, 30, 0, 0, time.UTC)
	targetID := uuid.NewString()
	rootID := uuid.NewString()
	replyID := uuid.NewString()
	hiddenID := uuid.NewString()

	root := entity.Comment{
		ID:               rootID,
		TargetType:       entity.TargetTypeChapter,
		TargetID:         targetID,
		AuthorUserID:     uuid.NewString(),
		Depth:            0,
		Content:          "Root",
		SanitizedContent: "Root",
		ModerationStatus: entity.ModerationStatusVisible,
		ReplyCount:       1,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	reply := entity.Comment{
		ID:               replyID,
		TargetType:       entity.TargetTypeChapter,
		TargetID:         targetID,
		AuthorUserID:     uuid.NewString(),
		ParentCommentID:  strPtr(rootID),
		RootCommentID:    strPtr(rootID),
		Depth:            1,
		Content:          "Reply",
		SanitizedContent: "Reply",
		ModerationStatus: entity.ModerationStatusVisible,
		CreatedAt:        now.Add(1 * time.Minute),
		UpdatedAt:        now.Add(1 * time.Minute),
	}
	hidden := entity.Comment{
		ID:               hiddenID,
		TargetType:       entity.TargetTypeChapter,
		TargetID:         targetID,
		AuthorUserID:     uuid.NewString(),
		Depth:            0,
		Content:          "Hidden",
		SanitizedContent: "Hidden",
		ModerationStatus: entity.ModerationStatusHidden,
		CreatedAt:        now.Add(2 * time.Minute),
		UpdatedAt:        now.Add(2 * time.Minute),
	}

	require.NoError(t, store.CreateComment(context.Background(), root))
	require.NoError(t, store.CreateComment(context.Background(), reply))
	require.NoError(t, store.CreateComment(context.Background(), hidden))

	items, err := store.ListCommentsByTarget(context.Background(), ListQuery{
		TargetType:     "chapter",
		TargetID:       targetID,
		ParentOnly:     true,
		SortBy:         "newest",
		Limit:          10,
		Offset:         0,
		IncludeHidden:  false,
		IncludeDeleted: true,
	})
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, rootID, items[0].ID)

	itemsWithHidden, err := store.ListCommentsByTarget(context.Background(), ListQuery{
		TargetType:     "chapter",
		TargetID:       targetID,
		ParentOnly:     true,
		SortBy:         "newest",
		Limit:          10,
		Offset:         0,
		IncludeHidden:  true,
		IncludeDeleted: true,
	})
	require.NoError(t, err)
	require.Len(t, itemsWithHidden, 2)
	require.Equal(t, hiddenID, itemsWithHidden[0].ID)

	thread, err := store.ListCommentsByRoot(context.Background(), ThreadQuery{
		RootCommentID:  rootID,
		SortBy:         "oldest",
		Limit:          10,
		Offset:         0,
		IncludeHidden:  false,
		IncludeDeleted: true,
	})
	require.NoError(t, err)
	require.Len(t, thread, 1)
	require.Equal(t, replyID, thread[0].ID)
}

package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
)

var (
	ErrNotFound = errors.New("chapter_repository_not_found")
	ErrConflict = errors.New("chapter_repository_conflict")
)

// ListQuery defines repository-level chapter listing controls.
type ListQuery struct {
	MangaID            string
	SortBy             string
	Limit              int
	Offset             int
	IncludeUnpublished bool
	IncludeDeleted     bool
}

// NavigationResult represents chapter navigation neighbors in a manga scope.
type NavigationResult struct {
	FoundCurrent bool
	FirstID      *string
	LastID       *string
	PreviousID   *string
	NextID       *string
}

// Store defines chapter persistence boundary.
type Store interface {
	CreateChapter(ctx context.Context, chapter entity.Chapter) error
	GetChapterByID(ctx context.Context, chapterID string) (entity.Chapter, error)
	GetChapterBySlug(ctx context.Context, mangaID string, slug string) (entity.Chapter, error)
	ListChaptersByManga(ctx context.Context, query ListQuery) ([]entity.Chapter, error)
	ResolveNavigation(ctx context.Context, mangaID string, chapterID string) (NavigationResult, error)
	UpdateChapter(ctx context.Context, chapter entity.Chapter) error
}

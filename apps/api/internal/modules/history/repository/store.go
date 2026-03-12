package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
)

var (
	ErrNotFound = errors.New("history_repository_not_found")
	ErrConflict = errors.New("history_repository_conflict")
)

// ContinueReadingQuery defines requester-owned continue-reading listing controls.
type ContinueReadingQuery struct {
	UserID string
	Limit  int
	Offset int
	SortBy string
}

// LibraryQuery defines requester-owned library listing controls.
type LibraryQuery struct {
	UserID     string
	Status     string
	Bookmarked *bool
	Favorited  *bool
	SharedOnly bool
	Limit      int
	Offset     int
	SortBy     string
}

// TimelineQuery defines requester-owned timeline listing controls.
type TimelineQuery struct {
	UserID string
	Event  string
	Limit  int
	Offset int
	SortBy string
}

// Store defines history persistence boundary.
type Store interface {
	UpsertCheckpoint(ctx context.Context, checkpoint entity.Checkpoint, dedupKey string) (entity.LibraryEntry, string, bool, error)
	GetLibraryEntry(ctx context.Context, userID string, mangaID string) (entity.LibraryEntry, error)
	UpsertLibraryEntry(ctx context.Context, entry entity.LibraryEntry) (entity.LibraryEntry, error)

	ListContinueReading(ctx context.Context, query ContinueReadingQuery) ([]entity.LibraryEntry, error)
	ListLibrary(ctx context.Context, query LibraryQuery) ([]entity.LibraryEntry, error)
	ListTimeline(ctx context.Context, query TimelineQuery) ([]entity.TimelineEvent, error)

	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error
}

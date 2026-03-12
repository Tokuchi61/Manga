package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
)

var (
	ErrNotFound = errors.New("manga_repository_not_found")
	ErrConflict = errors.New("manga_repository_conflict")
)

// ListQuery defines repository-level listing controls.
type ListQuery struct {
	Search             string
	Genre              string
	Tag                string
	Theme              string
	ContentWarning     string
	SortBy             string
	Limit              int
	Offset             int
	IncludeUnpublished bool
	IncludeHidden      bool
	IncludeDeleted     bool
}

// Store defines manga persistence boundary.
type Store interface {
	CreateManga(ctx context.Context, manga entity.Manga) error
	GetMangaByID(ctx context.Context, mangaID string) (entity.Manga, error)
	GetMangaBySlug(ctx context.Context, slug string) (entity.Manga, error)
	ListManga(ctx context.Context, query ListQuery) ([]entity.Manga, error)
	UpdateManga(ctx context.Context, manga entity.Manga) error
}

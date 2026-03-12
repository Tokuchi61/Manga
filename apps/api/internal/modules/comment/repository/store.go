package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/entity"
)

var (
	ErrNotFound = errors.New("comment_repository_not_found")
	ErrConflict = errors.New("comment_repository_conflict")
)

// ListQuery defines target-based listing controls.
type ListQuery struct {
	TargetType     string
	TargetID       string
	ParentOnly     bool
	SortBy         string
	Limit          int
	Offset         int
	IncludeHidden  bool
	IncludeDeleted bool
}

// ThreadQuery defines root thread listing controls.
type ThreadQuery struct {
	RootCommentID  string
	SortBy         string
	Limit          int
	Offset         int
	IncludeHidden  bool
	IncludeDeleted bool
}

// Store defines comment persistence boundary.
type Store interface {
	CreateComment(ctx context.Context, comment entity.Comment) error
	GetCommentByID(ctx context.Context, commentID string) (entity.Comment, error)
	GetLatestCommentByAuthor(ctx context.Context, authorUserID string) (entity.Comment, error)
	ListCommentsByTarget(ctx context.Context, query ListQuery) ([]entity.Comment, error)
	ListCommentsByRoot(ctx context.Context, query ThreadQuery) ([]entity.Comment, error)
	UpdateComment(ctx context.Context, comment entity.Comment) error
}

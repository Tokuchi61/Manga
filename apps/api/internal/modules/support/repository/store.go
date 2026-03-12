package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
)

var (
	ErrNotFound = errors.New("support_repository_not_found")
	ErrConflict = errors.New("support_repository_conflict")
)

// ListQuery defines requester-owned support listing controls.
type ListQuery struct {
	RequesterUserID string
	Status          string
	SortBy          string
	Limit           int
	Offset          int
}

// QueueQuery defines support review queue controls.
type QueueQuery struct {
	Status   string
	Priority string
	Limit    int
	Offset   int
}

// Store defines support persistence boundary.
type Store interface {
	CreateCase(ctx context.Context, support entity.SupportCase) error
	GetCaseByID(ctx context.Context, supportID string) (entity.SupportCase, error)
	ListCasesByRequester(ctx context.Context, query ListQuery) ([]entity.SupportCase, error)
	ListReviewQueue(ctx context.Context, query QueueQuery) ([]entity.SupportCase, error)
	UpdateCase(ctx context.Context, support entity.SupportCase) error
	FindCaseByRequesterRequestID(ctx context.Context, requesterUserID string, requestID string) (entity.SupportCase, error)
	FindRecentSimilarCase(ctx context.Context, requesterUserID string, kind entity.SupportKind, targetType *entity.SupportTargetType, targetID *string, reasonCode string, since time.Time) (entity.SupportCase, error)
}

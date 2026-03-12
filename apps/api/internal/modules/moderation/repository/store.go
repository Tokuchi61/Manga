package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
)

var (
	ErrNotFound = errors.New("moderation_repository_not_found")
	ErrConflict = errors.New("moderation_repository_conflict")
)

type QueueQuery struct {
	Status                 string
	TargetType             string
	AssignedModeratorUserID string
	SortBy                 string
	Limit                  int
	Offset                 int
}

// Store defines moderation persistence boundary.
type Store interface {
	CreateCase(ctx context.Context, moderationCase entity.Case) error
	GetCaseByID(ctx context.Context, caseID string) (entity.Case, error)
	GetCaseBySourceRef(ctx context.Context, source entity.CaseSource, sourceRefID string) (entity.Case, error)
	ListQueue(ctx context.Context, query QueueQuery) ([]entity.Case, error)
	UpdateCase(ctx context.Context, moderationCase entity.Case) error
}

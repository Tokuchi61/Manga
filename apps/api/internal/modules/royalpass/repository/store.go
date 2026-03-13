package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
)

var (
	ErrNotFound = errors.New("royalpass_repository_not_found")
)

// Store defines royalpass persistence boundary.
type Store interface {
	UpsertSeasonDefinition(ctx context.Context, season entity.SeasonDefinition) error
	GetSeasonDefinition(ctx context.Context, seasonID string) (entity.SeasonDefinition, error)
	ListSeasonDefinitions(ctx context.Context, state string, limit int, offset int) ([]entity.SeasonDefinition, error)
	ResolveCurrentSeason(ctx context.Context, now time.Time) (entity.SeasonDefinition, error)

	UpsertTierDefinition(ctx context.Context, tier entity.TierDefinition) error
	GetTierDefinition(ctx context.Context, seasonID string, tierNumber int, track string) (entity.TierDefinition, error)
	ListTierDefinitions(ctx context.Context, seasonID string, track string, activeOnly bool, limit int, offset int) ([]entity.TierDefinition, error)

	GetUserProgress(ctx context.Context, userID string, seasonID string) (entity.UserProgress, error)
	UpsertUserProgress(ctx context.Context, progress entity.UserProgress) error
	DeleteUserProgressByUser(ctx context.Context, userID string) (int, error)
	DeleteUserProgressByUserSeason(ctx context.Context, userID string, seasonID string) (int, error)

	GetProgressDedup(ctx context.Context, dedupKey string) (entity.UserProgress, error)
	PutProgressDedup(ctx context.Context, dedupKey string, progress entity.UserProgress) error
	GetClaimDedup(ctx context.Context, dedupKey string) (entity.UserProgress, error)
	PutClaimDedup(ctx context.Context, dedupKey string, progress entity.UserProgress) error
	GetPremiumActivationDedup(ctx context.Context, dedupKey string) (entity.UserProgress, error)
	PutPremiumActivationDedup(ctx context.Context, dedupKey string, progress entity.UserProgress) error

	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error
}

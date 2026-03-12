package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
)

var (
	ErrNotFound = errors.New("mission_repository_not_found")
)

// Store defines mission persistence boundary.
type Store interface {
	UpsertMissionDefinition(ctx context.Context, definition entity.MissionDefinition) error
	GetMissionDefinition(ctx context.Context, missionID string) (entity.MissionDefinition, error)
	ListMissionDefinitions(ctx context.Context, category string, activeOnly bool, limit int, offset int) ([]entity.MissionDefinition, error)

	GetMissionProgress(ctx context.Context, userID string, missionID string, periodKey string) (entity.MissionProgress, error)
	UpsertMissionProgress(ctx context.Context, progress entity.MissionProgress) error
	DeleteMissionProgressByUser(ctx context.Context, userID string) (int, error)
	DeleteMissionProgressByUserMission(ctx context.Context, userID string, missionID string) (int, error)

	GetProgressDedup(ctx context.Context, dedupKey string) (entity.MissionProgress, error)
	PutProgressDedup(ctx context.Context, dedupKey string, progress entity.MissionProgress) error
	GetClaimDedup(ctx context.Context, dedupKey string) (entity.MissionProgress, error)
	PutClaimDedup(ctx context.Context, dedupKey string, progress entity.MissionProgress) error

	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error
}

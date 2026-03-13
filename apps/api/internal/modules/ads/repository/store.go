package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
)

var (
	ErrNotFound = errors.New("ads_repository_not_found")
)

// Store defines ads persistence boundary.
type Store interface {
	UpsertPlacementDefinition(ctx context.Context, placement entity.PlacementDefinition) error
	GetPlacementDefinition(ctx context.Context, placementID string) (entity.PlacementDefinition, error)
	ListPlacementDefinitions(ctx context.Context, surface string, visibleOnly bool, limit int, offset int) ([]entity.PlacementDefinition, error)

	UpsertCampaignDefinition(ctx context.Context, campaign entity.CampaignDefinition) error
	GetCampaignDefinition(ctx context.Context, campaignID string) (entity.CampaignDefinition, error)
	ListCampaignDefinitions(ctx context.Context, placementID string, state string, limit int, offset int) ([]entity.CampaignDefinition, error)

	CreateImpressionLog(ctx context.Context, log entity.ImpressionLog) error
	CreateClickLog(ctx context.Context, log entity.ClickLog) error
	CountImpressionsBySessionPlacement(ctx context.Context, sessionID string, placementID string) (int, error)

	GetImpressionDedup(ctx context.Context, dedupKey string) (entity.ImpressionLog, error)
	PutImpressionDedup(ctx context.Context, dedupKey string, log entity.ImpressionLog) error
	GetClickDedup(ctx context.Context, dedupKey string) (entity.ClickLog, error)
	PutClickDedup(ctx context.Context, dedupKey string, log entity.ClickLog) error

	GetCampaignAggregate(ctx context.Context, campaignID string) (entity.CampaignAggregate, error)
	UpsertCampaignAggregate(ctx context.Context, aggregate entity.CampaignAggregate) error
	ListCampaignAggregate(ctx context.Context, limit int, offset int) ([]entity.CampaignAggregate, error)

	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error
}

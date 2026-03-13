package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
)

// MemoryStore is stage-20 bootstrap persistence for ads flows.
type MemoryStore struct {
	mu sync.RWMutex

	placementsByID       map[string]entity.PlacementDefinition
	campaignsByID        map[string]entity.CampaignDefinition
	impressionsByID      map[string]entity.ImpressionLog
	clicksByID           map[string]entity.ClickLog
	impressionDedupByKey map[string]entity.ImpressionLog
	clickDedupByKey      map[string]entity.ClickLog
	aggregateByCampaign  map[string]entity.CampaignAggregate

	runtimeConfig entity.RuntimeConfig
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	return &MemoryStore{
		placementsByID:       make(map[string]entity.PlacementDefinition),
		campaignsByID:        make(map[string]entity.CampaignDefinition),
		impressionsByID:      make(map[string]entity.ImpressionLog),
		clicksByID:           make(map[string]entity.ClickLog),
		impressionDedupByKey: make(map[string]entity.ImpressionLog),
		clickDedupByKey:      make(map[string]entity.ClickLog),
		aggregateByCampaign:  make(map[string]entity.CampaignAggregate),
		runtimeConfig: entity.RuntimeConfig{
			SurfaceEnabled:     true,
			PlacementEnabled:   true,
			CampaignEnabled:    true,
			ClickIntakeEnabled: true,
			UpdatedAt:          now,
		},
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func clonePlacement(in entity.PlacementDefinition) entity.PlacementDefinition {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneCampaign(in entity.CampaignDefinition) entity.CampaignDefinition {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	if in.StartsAt != nil {
		startsAt := in.StartsAt.UTC()
		out.StartsAt = &startsAt
	}
	if in.EndsAt != nil {
		endsAt := in.EndsAt.UTC()
		out.EndsAt = &endsAt
	}
	return out
}

func cloneImpression(in entity.ImpressionLog) entity.ImpressionLog {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	return out
}

func cloneClick(in entity.ClickLog) entity.ClickLog {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	return out
}

func cloneAggregate(in entity.CampaignAggregate) entity.CampaignAggregate {
	out := in
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneRuntimeConfig(in entity.RuntimeConfig) entity.RuntimeConfig {
	out := in
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func applyOffsetLimit[T any](items []T, offset int, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []T{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]T(nil), items[offset:end]...)
}

func sortByUpdatedDesc[T any](items []T, less func(i int, j int) bool) {
	sort.Slice(items, less)
}

package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
)

// MemoryStore is stage-16 bootstrap persistence for mission flows.
type MemoryStore struct {
	mu sync.RWMutex

	missionDefinitionsByID map[string]entity.MissionDefinition
	missionProgressByKey   map[string]entity.MissionProgress
	progressDedupByKey     map[string]entity.MissionProgress
	claimDedupByKey        map[string]entity.MissionProgress

	runtimeConfig entity.RuntimeConfig
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	return &MemoryStore{
		missionDefinitionsByID: make(map[string]entity.MissionDefinition),
		missionProgressByKey:   make(map[string]entity.MissionProgress),
		progressDedupByKey:     make(map[string]entity.MissionProgress),
		claimDedupByKey:        make(map[string]entity.MissionProgress),
		runtimeConfig: entity.RuntimeConfig{
			ReadEnabled:           true,
			ClaimEnabled:          true,
			ProgressIngestEnabled: true,
			DailyResetHourUTC:     0,
			UpdatedAt:             now,
		},
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func missionProgressKey(userID string, missionID string, periodKey string) string {
	return normalizeValue(userID) + ":" + normalizeValue(missionID) + ":" + normalizeValue(periodKey)
}

func cloneDefinition(in entity.MissionDefinition) entity.MissionDefinition {
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

func cloneProgress(in entity.MissionProgress) entity.MissionProgress {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	if in.CompletedAt != nil {
		completedAt := in.CompletedAt.UTC()
		out.CompletedAt = &completedAt
	}
	if in.ClaimedAt != nil {
		claimedAt := in.ClaimedAt.UTC()
		out.ClaimedAt = &claimedAt
	}
	return out
}

func cloneRuntimeConfig(in entity.RuntimeConfig) entity.RuntimeConfig {
	out := in
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func sortByTimeDesc[T any](items []T, less func(i int, j int) bool) {
	sort.Slice(items, less)
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

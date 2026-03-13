package repository

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
)

// MemoryStore is stage-17 bootstrap persistence for royalpass flows.
type MemoryStore struct {
	mu sync.RWMutex

	seasonDefinitionsByID       map[string]entity.SeasonDefinition
	tierDefinitionsByKey        map[string]entity.TierDefinition
	userProgressByKey           map[string]entity.UserProgress
	progressDedupByKey          map[string]entity.UserProgress
	claimDedupByKey             map[string]entity.UserProgress
	premiumActivationDedupByKey map[string]entity.UserProgress

	runtimeConfig entity.RuntimeConfig
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	return &MemoryStore{
		seasonDefinitionsByID:       make(map[string]entity.SeasonDefinition),
		tierDefinitionsByKey:        make(map[string]entity.TierDefinition),
		userProgressByKey:           make(map[string]entity.UserProgress),
		progressDedupByKey:          make(map[string]entity.UserProgress),
		claimDedupByKey:             make(map[string]entity.UserProgress),
		premiumActivationDedupByKey: make(map[string]entity.UserProgress),
		runtimeConfig: entity.RuntimeConfig{
			SeasonEnabled:  true,
			ClaimEnabled:   true,
			PremiumEnabled: true,
			UpdatedAt:      now,
		},
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func tierKey(seasonID string, tierNumber int, track string) string {
	return normalizeValue(seasonID) + ":" + normalizeValue(track) + ":" + normalizeValue(strconv.Itoa(tierNumber))
}

func progressKey(userID string, seasonID string) string {
	return normalizeValue(userID) + ":" + normalizeValue(seasonID)
}

func cloneSeason(in entity.SeasonDefinition) entity.SeasonDefinition {
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

func cloneTier(in entity.TierDefinition) entity.TierDefinition {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneClaimedTiers(in map[string]time.Time) map[string]time.Time {
	if in == nil {
		return make(map[string]time.Time)
	}
	out := make(map[string]time.Time, len(in))
	for key, value := range in {
		out[key] = value.UTC()
	}
	return out
}

func cloneUserProgress(in entity.UserProgress) entity.UserProgress {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	out.ClaimedTiers = cloneClaimedTiers(in.ClaimedTiers)
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

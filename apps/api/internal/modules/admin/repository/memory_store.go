package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
)

// MemoryStore is stage-21 bootstrap persistence for admin flows.
type MemoryStore struct {
	mu sync.RWMutex

	runtimeConfig entity.RuntimeConfig

	actionsByID       map[string]entity.AdminAction
	actionDedupByKey  map[string]entity.AdminAction
	overridesByID     map[string]entity.OverrideRecord
	userReviewsByID   map[string]entity.UserReviewRecord
	impersonationByID map[string]entity.ImpersonationSession
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	return &MemoryStore{
		runtimeConfig: entity.RuntimeConfig{
			MaintenanceEnabled: false,
			UpdatedAt:          now,
		},
		actionsByID:       make(map[string]entity.AdminAction),
		actionDedupByKey:  make(map[string]entity.AdminAction),
		overridesByID:     make(map[string]entity.OverrideRecord),
		userReviewsByID:   make(map[string]entity.UserReviewRecord),
		impersonationByID: make(map[string]entity.ImpersonationSession),
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func cloneRuntimeConfig(in entity.RuntimeConfig) entity.RuntimeConfig {
	out := in
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneAction(in entity.AdminAction) entity.AdminAction {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	if len(in.Metadata) > 0 {
		out.Metadata = make(map[string]string, len(in.Metadata))
		for key, value := range in.Metadata {
			out.Metadata[key] = value
		}
	}
	return out
}

func cloneOverride(in entity.OverrideRecord) entity.OverrideRecord {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	if in.ExpiresAt != nil {
		expiresAt := in.ExpiresAt.UTC()
		out.ExpiresAt = &expiresAt
	}
	return out
}

func cloneUserReview(in entity.UserReviewRecord) entity.UserReviewRecord {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneImpersonationSession(in entity.ImpersonationSession) entity.ImpersonationSession {
	out := in
	out.StartedAt = in.StartedAt.UTC()
	out.ExpiresAt = in.ExpiresAt.UTC()
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	if in.EndedAt != nil {
		endedAt := in.EndedAt.UTC()
		out.EndedAt = &endedAt
	}
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
		limit = len(items)
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]T(nil), items[offset:end]...)
}

func sortByCreatedAtDesc[T any](items []T, lessFn func(i int, j int) bool) {
	sort.Slice(items, lessFn)
}

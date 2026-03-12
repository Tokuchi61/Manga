package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/entity"
)

// MemoryStore is stage-15 bootstrap persistence for inventory flows.
type MemoryStore struct {
	mu sync.RWMutex

	itemDefinitionsByID map[string]entity.ItemDefinition
	inventoryByKey      map[string]entity.InventoryEntry
	grantDedupByKey     map[string]entity.InventoryEntry

	runtimeConfig entity.RuntimeConfig
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	return &MemoryStore{
		itemDefinitionsByID: make(map[string]entity.ItemDefinition),
		inventoryByKey:      make(map[string]entity.InventoryEntry),
		grantDedupByKey:     make(map[string]entity.InventoryEntry),
		runtimeConfig: entity.RuntimeConfig{
			ReadEnabled:    true,
			ClaimEnabled:   true,
			ConsumeEnabled: true,
			EquipEnabled:   true,
			UpdatedAt:      now,
		},
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func inventoryKey(userID string, itemID string) string {
	return normalizeValue(userID) + ":" + normalizeValue(itemID)
}

func cloneItemDefinition(in entity.ItemDefinition) entity.ItemDefinition {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneEntry(in entity.InventoryEntry) entity.InventoryEntry {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	if in.ExpiresAt != nil {
		expiresAt := in.ExpiresAt.UTC()
		out.ExpiresAt = &expiresAt
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

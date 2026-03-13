package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
)

// MemoryStore is stage-18 bootstrap persistence for shop flows.
type MemoryStore struct {
	mu sync.RWMutex

	productDefinitionsByID map[string]entity.ProductDefinition
	offerDefinitionsByID   map[string]entity.OfferDefinition
	purchaseIntentsByID    map[string]entity.PurchaseIntent
	purchaseDedupByKey     map[string]entity.PurchaseIntent

	runtimeConfig entity.RuntimeConfig
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	return &MemoryStore{
		productDefinitionsByID: make(map[string]entity.ProductDefinition),
		offerDefinitionsByID:   make(map[string]entity.OfferDefinition),
		purchaseIntentsByID:    make(map[string]entity.PurchaseIntent),
		purchaseDedupByKey:     make(map[string]entity.PurchaseIntent),
		runtimeConfig: entity.RuntimeConfig{
			CatalogEnabled:  true,
			PurchaseEnabled: true,
			CampaignEnabled: true,
			UpdatedAt:       now,
		},
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func cloneProduct(in entity.ProductDefinition) entity.ProductDefinition {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneOffer(in entity.OfferDefinition) entity.OfferDefinition {
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

func clonePurchaseIntent(in entity.PurchaseIntent) entity.PurchaseIntent {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
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

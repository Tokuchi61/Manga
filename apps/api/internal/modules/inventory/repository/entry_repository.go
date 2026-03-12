package repository

import (
	"context"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/entity"
)

func (s *MemoryStore) GetInventoryEntry(_ context.Context, userID string, itemID string) (entity.InventoryEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.inventoryByKey[inventoryKey(userID, itemID)]
	if !ok {
		return entity.InventoryEntry{}, ErrNotFound
	}
	return cloneEntry(entry), nil
}

func (s *MemoryStore) UpsertInventoryEntry(_ context.Context, entry entity.InventoryEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.inventoryByKey[inventoryKey(entry.UserID, entry.ItemID)] = cloneEntry(entry)
	return nil
}

func (s *MemoryStore) DeleteInventoryEntry(_ context.Context, userID string, itemID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := inventoryKey(userID, itemID)
	if _, exists := s.inventoryByKey[key]; !exists {
		return ErrNotFound
	}
	delete(s.inventoryByKey, key)
	return nil
}

func (s *MemoryStore) ListInventoryEntries(_ context.Context, userID string, itemType string, equippedOnly bool, sortBy string, limit int, offset int) ([]entity.InventoryEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedUserID := normalizeValue(userID)
	normalizedType := normalizeValue(itemType)
	items := make([]entity.InventoryEntry, 0)
	for _, entry := range s.inventoryByKey {
		if normalizeValue(entry.UserID) != normalizedUserID {
			continue
		}
		if equippedOnly && !entry.Equipped {
			continue
		}
		if normalizedType != "" {
			definition, ok := s.itemDefinitionsByID[normalizeValue(entry.ItemID)]
			if !ok || normalizeValue(definition.ItemType) != normalizedType {
				continue
			}
		}
		items = append(items, cloneEntry(entry))
	}

	normalizedSort := normalizeValue(sortBy)
	if normalizedSort == "oldest" {
		sortByTimeDesc(items, func(i int, j int) bool {
			if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
				return strings.Compare(items[i].ItemID, items[j].ItemID) < 0
			}
			return items[i].UpdatedAt.Before(items[j].UpdatedAt)
		})
	} else {
		sortByTimeDesc(items, func(i int, j int) bool {
			if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
				return strings.Compare(items[i].ItemID, items[j].ItemID) < 0
			}
			return items[i].UpdatedAt.After(items[j].UpdatedAt)
		})
	}

	return applyOffsetLimit(items, offset, limit), nil
}

func (s *MemoryStore) GetGrantByDedup(_ context.Context, dedupKey string) (entity.InventoryEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.grantDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.InventoryEntry{}, ErrNotFound
	}
	return cloneEntry(entry), nil
}

func (s *MemoryStore) PutGrantDedup(_ context.Context, dedupKey string, entry entity.InventoryEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.grantDedupByKey[normalizeValue(dedupKey)] = cloneEntry(entry)
	return nil
}

package repository

import (
	"context"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/entity"
)

func (s *MemoryStore) UpsertItemDefinition(_ context.Context, definition entity.ItemDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := normalizeValue(definition.ItemID)
	s.itemDefinitionsByID[key] = cloneItemDefinition(definition)
	return nil
}

func (s *MemoryStore) GetItemDefinition(_ context.Context, itemID string) (entity.ItemDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	definition, ok := s.itemDefinitionsByID[normalizeValue(itemID)]
	if !ok {
		return entity.ItemDefinition{}, ErrNotFound
	}
	return cloneItemDefinition(definition), nil
}

func (s *MemoryStore) ListItemDefinitions(_ context.Context, itemType string) ([]entity.ItemDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedType := normalizeValue(itemType)
	items := make([]entity.ItemDefinition, 0, len(s.itemDefinitionsByID))
	for _, definition := range s.itemDefinitionsByID {
		if normalizedType != "" && normalizeValue(definition.ItemType) != normalizedType {
			continue
		}
		items = append(items, cloneItemDefinition(definition))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return strings.Compare(items[i].ItemID, items[j].ItemID) < 0
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return items, nil
}

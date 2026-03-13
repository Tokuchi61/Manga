package repository

import (
	"context"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
)

func (s *MemoryStore) UpsertProductDefinition(_ context.Context, product entity.ProductDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := normalizeValue(product.ProductID)
	if existing, ok := s.productDefinitionsByID[key]; ok {
		product.CreatedAt = existing.CreatedAt
	}
	if product.CreatedAt.IsZero() {
		product.CreatedAt = time.Now().UTC()
	}
	s.productDefinitionsByID[key] = cloneProduct(product)
	return nil
}

func (s *MemoryStore) GetProductDefinition(_ context.Context, productID string) (entity.ProductDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, ok := s.productDefinitionsByID[normalizeValue(productID)]
	if !ok {
		return entity.ProductDefinition{}, ErrNotFound
	}
	return cloneProduct(product), nil
}

func (s *MemoryStore) ListProductDefinitions(_ context.Context, state string, limit int, offset int) ([]entity.ProductDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedState := normalizeValue(state)
	items := make([]entity.ProductDefinition, 0, len(s.productDefinitionsByID))
	for _, product := range s.productDefinitionsByID {
		if normalizedState != "" && normalizeValue(product.State) != normalizedState {
			continue
		}
		items = append(items, cloneProduct(product))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return strings.Compare(items[i].ProductID, items[j].ProductID) < 0
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return applyOffsetLimit(items, offset, limit), nil
}

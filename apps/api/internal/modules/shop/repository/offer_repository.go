package repository

import (
	"context"
	"sort"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
)

func (s *MemoryStore) UpsertOfferDefinition(_ context.Context, offer entity.OfferDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := normalizeValue(offer.OfferID)
	if existing, ok := s.offerDefinitionsByID[key]; ok {
		offer.CreatedAt = existing.CreatedAt
	}
	if offer.CreatedAt.IsZero() {
		offer.CreatedAt = offer.UpdatedAt
	}
	s.offerDefinitionsByID[key] = cloneOffer(offer)
	return nil
}

func (s *MemoryStore) GetOfferDefinition(_ context.Context, offerID string) (entity.OfferDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	offer, ok := s.offerDefinitionsByID[normalizeValue(offerID)]
	if !ok {
		return entity.OfferDefinition{}, ErrNotFound
	}
	return cloneOffer(offer), nil
}

func (s *MemoryStore) ListOfferDefinitions(_ context.Context, productID string, visibility string, activeOnly bool, limit int, offset int) ([]entity.OfferDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedProductID := normalizeValue(productID)
	normalizedVisibility := normalizeValue(visibility)
	items := make([]entity.OfferDefinition, 0, len(s.offerDefinitionsByID))
	for _, offer := range s.offerDefinitionsByID {
		if normalizedProductID != "" && normalizeValue(offer.ProductID) != normalizedProductID {
			continue
		}
		if normalizedVisibility != "" && normalizeValue(offer.Visibility) != normalizedVisibility {
			continue
		}
		if activeOnly && !offer.Active {
			continue
		}
		items = append(items, cloneOffer(offer))
	}

	sort.Slice(items, func(i int, j int) bool {
		if items[i].ProductID == items[j].ProductID {
			if items[i].PriceMana == items[j].PriceMana {
				return items[i].UpdatedAt.After(items[j].UpdatedAt)
			}
			return items[i].PriceMana < items[j].PriceMana
		}
		return items[i].ProductID < items[j].ProductID
	})

	return applyOffsetLimit(items, offset, limit), nil
}

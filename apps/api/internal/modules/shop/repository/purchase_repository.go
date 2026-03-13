package repository

import (
	"context"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
)

func (s *MemoryStore) CreatePurchaseIntent(_ context.Context, intent entity.PurchaseIntent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.purchaseIntentsByID[normalizeValue(intent.IntentID)] = clonePurchaseIntent(intent)
	return nil
}

func (s *MemoryStore) GetPurchaseIntent(_ context.Context, intentID string) (entity.PurchaseIntent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	intent, ok := s.purchaseIntentsByID[normalizeValue(intentID)]
	if !ok {
		return entity.PurchaseIntent{}, ErrNotFound
	}
	return clonePurchaseIntent(intent), nil
}

func (s *MemoryStore) UpdatePurchaseIntent(_ context.Context, intent entity.PurchaseIntent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := normalizeValue(intent.IntentID)
	if _, ok := s.purchaseIntentsByID[key]; !ok {
		return ErrNotFound
	}
	s.purchaseIntentsByID[key] = clonePurchaseIntent(intent)
	return nil
}

func (s *MemoryStore) ListPurchaseIntentsByUser(_ context.Context, userID string, limit int, offset int) ([]entity.PurchaseIntent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedUserID := normalizeValue(userID)
	items := make([]entity.PurchaseIntent, 0)
	for _, intent := range s.purchaseIntentsByID {
		if normalizeValue(intent.UserID) != normalizedUserID {
			continue
		}
		items = append(items, clonePurchaseIntent(intent))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return strings.Compare(items[i].IntentID, items[j].IntentID) < 0
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return applyOffsetLimit(items, offset, limit), nil
}

func (s *MemoryStore) FindLatestPurchaseByUserProduct(_ context.Context, userID string, productID string) (entity.PurchaseIntent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedUserID := normalizeValue(userID)
	normalizedProductID := normalizeValue(productID)
	var found *entity.PurchaseIntent
	for _, intent := range s.purchaseIntentsByID {
		if normalizeValue(intent.UserID) != normalizedUserID {
			continue
		}
		if normalizeValue(intent.ProductID) != normalizedProductID {
			continue
		}
		candidate := clonePurchaseIntent(intent)
		if found == nil || candidate.UpdatedAt.After(found.UpdatedAt) {
			copyCandidate := candidate
			found = &copyCandidate
		}
	}
	if found == nil {
		return entity.PurchaseIntent{}, ErrNotFound
	}
	return *found, nil
}

func (s *MemoryStore) GetPurchaseDedup(_ context.Context, dedupKey string) (entity.PurchaseIntent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	intent, ok := s.purchaseDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.PurchaseIntent{}, ErrNotFound
	}
	return clonePurchaseIntent(intent), nil
}

func (s *MemoryStore) PutPurchaseDedup(_ context.Context, dedupKey string, intent entity.PurchaseIntent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.purchaseDedupByKey[normalizeValue(dedupKey)] = clonePurchaseIntent(intent)
	return nil
}

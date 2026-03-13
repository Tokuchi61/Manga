package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
)

func (s *MemoryStore) CreateTransaction(_ context.Context, tx entity.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.transactionsByID[normalizeValue(tx.TransactionID)] = cloneTransaction(tx)
	return nil
}

func (s *MemoryStore) GetTransaction(_ context.Context, transactionID string) (entity.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.transactionsByID[normalizeValue(transactionID)]
	if !ok {
		return entity.Transaction{}, ErrNotFound
	}
	return cloneTransaction(tx), nil
}

func (s *MemoryStore) UpdateTransaction(_ context.Context, tx entity.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := normalizeValue(tx.TransactionID)
	if _, ok := s.transactionsByID[key]; !ok {
		return ErrNotFound
	}
	s.transactionsByID[key] = cloneTransaction(tx)
	return nil
}

func (s *MemoryStore) ListTransactionsByUser(_ context.Context, userID string, status string, limit int, offset int) ([]entity.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedUserID := normalizeValue(userID)
	normalizedStatus := normalizeValue(status)

	items := make([]entity.Transaction, 0, len(s.transactionsByID))
	for _, tx := range s.transactionsByID {
		if normalizeValue(tx.UserID) != normalizedUserID {
			continue
		}
		if normalizedStatus != "" && normalizeValue(tx.Status) != normalizedStatus {
			continue
		}
		items = append(items, cloneTransaction(tx))
	}

	sortByUpdatedDesc(items, func(i int, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return applyOffsetLimit(items, offset, limit), nil
}

func (s *MemoryStore) GetCheckoutDedup(_ context.Context, dedupKey string) (entity.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.checkoutDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.Transaction{}, ErrNotFound
	}
	return cloneTransaction(tx), nil
}

func (s *MemoryStore) PutCheckoutDedup(_ context.Context, dedupKey string, tx entity.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.checkoutDedupByKey[normalizeValue(dedupKey)] = cloneTransaction(tx)
	return nil
}

func (s *MemoryStore) GetCallbackDedup(_ context.Context, providerEventID string) (entity.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.callbackDedupByEvent[normalizeValue(providerEventID)]
	if !ok {
		return entity.Transaction{}, ErrNotFound
	}
	return cloneTransaction(tx), nil
}

func (s *MemoryStore) PutCallbackDedup(_ context.Context, providerEventID string, tx entity.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.callbackDedupByEvent[normalizeValue(providerEventID)] = cloneTransaction(tx)
	return nil
}

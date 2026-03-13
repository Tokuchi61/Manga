package repository

import (
	"context"
	"sort"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
)

func (s *MemoryStore) CreateLedgerEntry(_ context.Context, entry entity.LedgerEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ledgerEntriesByID[normalizeValue(entry.EntryID)] = cloneLedgerEntry(entry)
	return nil
}

func (s *MemoryStore) ListLedgerEntriesByUser(_ context.Context, userID string) ([]entity.LedgerEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedUserID := normalizeValue(userID)
	items := make([]entity.LedgerEntry, 0, len(s.ledgerEntriesByID))
	for _, entry := range s.ledgerEntriesByID {
		if normalizeValue(entry.UserID) != normalizedUserID {
			continue
		}
		items = append(items, cloneLedgerEntry(entry))
	}

	sort.Slice(items, func(i int, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	return items, nil
}

func (s *MemoryStore) ListLedgerUsers(_ context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uniq := make(map[string]string)
	for _, entry := range s.ledgerEntriesByID {
		key := normalizeValue(entry.UserID)
		if key == "" {
			continue
		}
		uniq[key] = entry.UserID
	}

	users := make([]string, 0, len(uniq))
	for _, userID := range uniq {
		users = append(users, userID)
	}
	sort.Strings(users)
	return users, nil
}

func (s *MemoryStore) GetBalanceSnapshot(_ context.Context, userID string) (entity.BalanceSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot, ok := s.balanceByUserID[normalizeValue(userID)]
	if !ok {
		return entity.BalanceSnapshot{}, ErrNotFound
	}
	return cloneBalanceSnapshot(snapshot), nil
}

func (s *MemoryStore) UpsertBalanceSnapshot(_ context.Context, snapshot entity.BalanceSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.balanceByUserID[normalizeValue(snapshot.UserID)] = cloneBalanceSnapshot(snapshot)
	return nil
}

package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
)

func (s *MemoryStore) CreateOverrideRecord(ctx context.Context, record entity.OverrideRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.overridesByID[record.OverrideID] = cloneOverride(record)
	return nil
}

func (s *MemoryStore) ListOverrideRecords(ctx context.Context, targetModule string, limit int, offset int) ([]entity.OverrideRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filtered := make([]entity.OverrideRecord, 0, len(s.overridesByID))
	for _, record := range s.overridesByID {
		if targetModule != "" && normalizeValue(record.TargetModule) != normalizeValue(targetModule) {
			continue
		}
		filtered = append(filtered, cloneOverride(record))
	}

	sortByCreatedAtDesc(filtered, func(i int, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	return applyOffsetLimit(filtered, offset, limit), nil
}

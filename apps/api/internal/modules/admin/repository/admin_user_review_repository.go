package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
)

func (s *MemoryStore) CreateUserReviewRecord(ctx context.Context, record entity.UserReviewRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userReviewsByID[record.ReviewID] = cloneUserReview(record)
	return nil
}

func (s *MemoryStore) ListUserReviewRecords(ctx context.Context, targetUserID string, limit int, offset int) ([]entity.UserReviewRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filtered := make([]entity.UserReviewRecord, 0, len(s.userReviewsByID))
	for _, record := range s.userReviewsByID {
		if targetUserID != "" && normalizeValue(record.TargetUserID) != normalizeValue(targetUserID) {
			continue
		}
		filtered = append(filtered, cloneUserReview(record))
	}

	sortByCreatedAtDesc(filtered, func(i int, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	return applyOffsetLimit(filtered, offset, limit), nil
}

package repository

import (
	"context"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
)

func (s *MemoryStore) CreateCase(_ context.Context, support entity.SupportCase) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.casesByID[support.ID]; exists {
		return ErrConflict
	}

	stored := cloneSupportCase(support)
	stored.RequesterUserID = normalizeValue(support.RequesterUserID)
	stored.ReasonCode = normalizeValue(support.ReasonCode)
	if stored.TargetType != nil {
		normalizedTargetType := entity.SupportTargetType(normalizeValue(string(*stored.TargetType)))
		stored.TargetType = &normalizedTargetType
	}
	if stored.TargetID != nil {
		normalizedTargetID := normalizeValue(*stored.TargetID)
		stored.TargetID = &normalizedTargetID
	}

	requestIDKey := normalizeValue(stored.RequestID)
	if requestIDKey != "" {
		requesterKey := normalizeValue(stored.RequesterUserID)
		if _, ok := s.requestIDIndex[requesterKey]; !ok {
			s.requestIDIndex[requesterKey] = make(map[string]string)
		}
		if _, exists := s.requestIDIndex[requesterKey][requestIDKey]; exists {
			return ErrConflict
		}
		s.requestIDIndex[requesterKey][requestIDKey] = stored.ID
	}

	s.casesByID[stored.ID] = stored
	return nil
}

func (s *MemoryStore) GetCaseByID(_ context.Context, supportID string) (entity.SupportCase, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	support, ok := s.casesByID[supportID]
	if !ok {
		return entity.SupportCase{}, ErrNotFound
	}
	return cloneSupportCase(support), nil
}

func (s *MemoryStore) ListCasesByRequester(_ context.Context, query ListQuery) ([]entity.SupportCase, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	requesterKey := normalizeValue(query.RequesterUserID)
	statusKey := normalizeValue(query.Status)

	items := make([]entity.SupportCase, 0)
	for _, support := range s.casesByID {
		if normalizeValue(support.RequesterUserID) != requesterKey {
			continue
		}
		if statusKey != "" && normalizeValue(string(support.Status)) != statusKey {
			continue
		}
		items = append(items, cloneSupportCase(support))
	}

	sortSupports(items, query.SortBy)
	return applyOffsetLimit(items, query.Offset, query.Limit), nil
}

func (s *MemoryStore) ListReviewQueue(_ context.Context, query QueueQuery) ([]entity.SupportCase, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statusKey := normalizeValue(query.Status)
	priorityKey := normalizeValue(query.Priority)

	items := make([]entity.SupportCase, 0)
	for _, support := range s.casesByID {
		if statusKey != "" {
			if normalizeValue(string(support.Status)) != statusKey {
				continue
			}
		} else if !defaultQueueStatus(support.Status) {
			continue
		}
		if priorityKey != "" && normalizeValue(string(support.Priority)) != priorityKey {
			continue
		}
		items = append(items, cloneSupportCase(support))
	}

	sortSupports(items, "priority")
	return applyOffsetLimit(items, query.Offset, query.Limit), nil
}

func (s *MemoryStore) UpdateCase(_ context.Context, support entity.SupportCase) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.casesByID[support.ID]
	if !exists {
		return ErrNotFound
	}

	updated := cloneSupportCase(support)
	updated.RequesterUserID = normalizeValue(updated.RequesterUserID)
	updated.ReasonCode = normalizeValue(updated.ReasonCode)
	if updated.TargetType != nil {
		normalizedTargetType := entity.SupportTargetType(normalizeValue(string(*updated.TargetType)))
		updated.TargetType = &normalizedTargetType
	}
	if updated.TargetID != nil {
		normalizedTargetID := normalizeValue(*updated.TargetID)
		updated.TargetID = &normalizedTargetID
	}

	currentRequestKey := normalizeValue(current.RequestID)
	updatedRequestKey := normalizeValue(updated.RequestID)
	requesterKey := normalizeValue(updated.RequesterUserID)
	if currentRequestKey != updatedRequestKey {
		if currentRequestKey != "" {
			delete(s.requestIDIndex[requesterKey], currentRequestKey)
		}
		if updatedRequestKey != "" {
			if _, ok := s.requestIDIndex[requesterKey]; !ok {
				s.requestIDIndex[requesterKey] = make(map[string]string)
			}
			if ownerID, ownerExists := s.requestIDIndex[requesterKey][updatedRequestKey]; ownerExists && ownerID != support.ID {
				return ErrConflict
			}
			s.requestIDIndex[requesterKey][updatedRequestKey] = support.ID
		}
	}

	s.casesByID[support.ID] = updated
	return nil
}

func (s *MemoryStore) FindCaseByRequesterRequestID(_ context.Context, requesterUserID string, requestID string) (entity.SupportCase, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	requesterKey := normalizeValue(requesterUserID)
	requestIDKey := normalizeValue(requestID)
	if requestIDKey == "" {
		return entity.SupportCase{}, ErrNotFound
	}

	supportID, ok := s.requestIDIndex[requesterKey][requestIDKey]
	if !ok {
		return entity.SupportCase{}, ErrNotFound
	}
	support, ok := s.casesByID[supportID]
	if !ok {
		return entity.SupportCase{}, ErrNotFound
	}
	return cloneSupportCase(support), nil
}

func (s *MemoryStore) FindRecentSimilarCase(_ context.Context, requesterUserID string, kind entity.SupportKind, targetType *entity.SupportTargetType, targetID *string, reasonCode string, since time.Time) (entity.SupportCase, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var latest entity.SupportCase
	found := false
	for _, support := range s.casesByID {
		if !similarCaseMatch(support, requesterUserID, kind, targetType, targetID, reasonCode, since) {
			continue
		}
		if !found || support.CreatedAt.After(latest.CreatedAt) {
			latest = support
			found = true
		}
	}
	if !found {
		return entity.SupportCase{}, ErrNotFound
	}
	return cloneSupportCase(latest), nil
}

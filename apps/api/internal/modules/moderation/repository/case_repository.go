package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
)

func (s *MemoryStore) CreateCase(_ context.Context, moderationCase entity.Case) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.casesByID[moderationCase.ID]; exists {
		return ErrConflict
	}

	stored := cloneCase(moderationCase)
	stored.Source = entity.CaseSource(normalizeValue(string(stored.Source)))
	stored.TargetType = entity.TargetType(normalizeValue(string(stored.TargetType)))
	stored.Status = entity.CaseStatus(normalizeValue(string(stored.Status)))
	stored.AssignmentStatus = entity.AssignmentStatus(normalizeValue(string(stored.AssignmentStatus)))
	stored.EscalationStatus = entity.EscalationStatus(normalizeValue(string(stored.EscalationStatus)))
	stored.ActionResult = entity.ActionResult(normalizeValue(string(stored.ActionResult)))
	stored.EscalationReason = normalizeValue(stored.EscalationReason)
	stored.TargetID = normalizeValue(stored.TargetID)
	stored.RequestID = normalizeValue(stored.RequestID)
	stored.CorrelationID = normalizeValue(stored.CorrelationID)
	if stored.SourceRefID != nil {
		normalizedRef := normalizeValue(*stored.SourceRefID)
		stored.SourceRefID = &normalizedRef
	}
	if stored.ReporterUserID != nil {
		normalizedReporter := normalizeValue(*stored.ReporterUserID)
		stored.ReporterUserID = &normalizedReporter
	}
	if stored.AssignedModeratorUserID != nil {
		normalizedAssignee := normalizeValue(*stored.AssignedModeratorUserID)
		stored.AssignedModeratorUserID = &normalizedAssignee
	}

	if stored.SourceRefID != nil && *stored.SourceRefID != "" {
		key := sourceRefKey(stored.Source, *stored.SourceRefID)
		if _, exists := s.sourceRefIndex[key]; exists {
			return ErrConflict
		}
		s.sourceRefIndex[key] = stored.ID
	}

	s.casesByID[stored.ID] = stored
	return nil
}

func (s *MemoryStore) GetCaseByID(_ context.Context, caseID string) (entity.Case, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	moderationCase, ok := s.casesByID[caseID]
	if !ok {
		return entity.Case{}, ErrNotFound
	}
	return cloneCase(moderationCase), nil
}

func (s *MemoryStore) GetCaseBySourceRef(_ context.Context, source entity.CaseSource, sourceRefID string) (entity.Case, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := sourceRefKey(source, sourceRefID)
	caseID, ok := s.sourceRefIndex[key]
	if !ok {
		return entity.Case{}, ErrNotFound
	}
	moderationCase, ok := s.casesByID[caseID]
	if !ok {
		return entity.Case{}, ErrNotFound
	}
	return cloneCase(moderationCase), nil
}

func (s *MemoryStore) ListQueue(_ context.Context, query QueueQuery) ([]entity.Case, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statusKey := normalizeValue(query.Status)
	targetTypeKey := normalizeValue(query.TargetType)
	assignedModeratorKey := normalizeValue(query.AssignedModeratorUserID)

	items := make([]entity.Case, 0)
	for _, moderationCase := range s.casesByID {
		if statusKey != "" && normalizeValue(string(moderationCase.Status)) != statusKey {
			continue
		}
		if targetTypeKey != "" && normalizeValue(string(moderationCase.TargetType)) != targetTypeKey {
			continue
		}
		if assignedModeratorKey != "" {
			if moderationCase.AssignedModeratorUserID == nil || normalizeValue(*moderationCase.AssignedModeratorUserID) != assignedModeratorKey {
				continue
			}
		}
		items = append(items, cloneCase(moderationCase))
	}

	sortCases(items, query.SortBy)
	return applyOffsetLimit(items, query.Offset, query.Limit), nil
}

func (s *MemoryStore) UpdateCase(_ context.Context, moderationCase entity.Case) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.casesByID[moderationCase.ID]
	if !exists {
		return ErrNotFound
	}

	updated := cloneCase(moderationCase)
	updated.Source = entity.CaseSource(normalizeValue(string(updated.Source)))
	updated.TargetType = entity.TargetType(normalizeValue(string(updated.TargetType)))
	updated.Status = entity.CaseStatus(normalizeValue(string(updated.Status)))
	updated.AssignmentStatus = entity.AssignmentStatus(normalizeValue(string(updated.AssignmentStatus)))
	updated.EscalationStatus = entity.EscalationStatus(normalizeValue(string(updated.EscalationStatus)))
	updated.ActionResult = entity.ActionResult(normalizeValue(string(updated.ActionResult)))
	updated.EscalationReason = normalizeValue(updated.EscalationReason)
	updated.TargetID = normalizeValue(updated.TargetID)
	updated.RequestID = normalizeValue(updated.RequestID)
	updated.CorrelationID = normalizeValue(updated.CorrelationID)
	if updated.SourceRefID != nil {
		normalizedRef := normalizeValue(*updated.SourceRefID)
		updated.SourceRefID = &normalizedRef
	}
	if updated.ReporterUserID != nil {
		normalizedReporter := normalizeValue(*updated.ReporterUserID)
		updated.ReporterUserID = &normalizedReporter
	}
	if updated.AssignedModeratorUserID != nil {
		normalizedAssignee := normalizeValue(*updated.AssignedModeratorUserID)
		updated.AssignedModeratorUserID = &normalizedAssignee
	}

	currentRefKey := ""
	if current.SourceRefID != nil && *current.SourceRefID != "" {
		currentRefKey = sourceRefKey(current.Source, *current.SourceRefID)
	}
	updatedRefKey := ""
	if updated.SourceRefID != nil && *updated.SourceRefID != "" {
		updatedRefKey = sourceRefKey(updated.Source, *updated.SourceRefID)
	}
	if currentRefKey != updatedRefKey {
		if currentRefKey != "" {
			delete(s.sourceRefIndex, currentRefKey)
		}
		if updatedRefKey != "" {
			if ownerID, exists := s.sourceRefIndex[updatedRefKey]; exists && ownerID != updated.ID {
				return ErrConflict
			}
			s.sourceRefIndex[updatedRefKey] = updated.ID
		}
	}

	s.casesByID[updated.ID] = updated
	return nil
}

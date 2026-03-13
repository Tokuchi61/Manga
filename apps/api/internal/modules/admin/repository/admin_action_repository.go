package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
)

func (s *MemoryStore) CreateAdminAction(ctx context.Context, action entity.AdminAction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.actionsByID[action.ActionID] = cloneAction(action)
	return nil
}

func (s *MemoryStore) ListAdminActions(ctx context.Context, actionType string, riskLevel string, limit int, offset int) ([]entity.AdminAction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filtered := make([]entity.AdminAction, 0, len(s.actionsByID))
	for _, action := range s.actionsByID {
		if actionType != "" && normalizeValue(action.ActionType) != normalizeValue(actionType) {
			continue
		}
		if riskLevel != "" && normalizeValue(action.RiskLevel) != normalizeValue(riskLevel) {
			continue
		}
		filtered = append(filtered, cloneAction(action))
	}

	sortByCreatedAtDesc(filtered, func(i int, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	return applyOffsetLimit(filtered, offset, limit), nil
}

func (s *MemoryStore) GetActionDedup(ctx context.Context, dedupKey string) (entity.AdminAction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	action, ok := s.actionDedupByKey[dedupKey]
	if !ok {
		return entity.AdminAction{}, ErrNotFound
	}
	return cloneAction(action), nil
}

func (s *MemoryStore) PutActionDedup(ctx context.Context, dedupKey string, action entity.AdminAction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.actionDedupByKey[dedupKey] = cloneAction(action)
	return nil
}

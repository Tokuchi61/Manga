package repository

import (
	"context"
	"sort"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	"github.com/google/uuid"
)

func (s *MemoryStore) CreatePolicyRule(_ context.Context, rule entity.PolicyRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rule.Key = normalizePolicyKey(rule.Key)
	conflictKey := rule.ConflictKey()
	if rule.Active {
		if _, exists := s.activePolicyConflict[conflictKey]; exists {
			return ErrConflict
		}
		s.activePolicyConflict[conflictKey] = rule.ID
	}

	s.policyRulesByID[rule.ID] = rule
	if _, ok := s.policyByKeyIndex[rule.Key]; !ok {
		s.policyByKeyIndex[rule.Key] = make(map[uuid.UUID]struct{})
	}
	s.policyByKeyIndex[rule.Key][rule.ID] = struct{}{}
	return nil
}

func (s *MemoryStore) ListPolicyRulesByKey(_ context.Context, key string) ([]entity.PolicyRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedKey := normalizePolicyKey(key)
	ruleIDs := s.policyByKeyIndex[normalizedKey]
	result := make([]entity.PolicyRule, 0, len(ruleIDs))
	for ruleID := range ruleIDs {
		rule, ok := s.policyRulesByID[ruleID]
		if ok {
			result = append(result, rule)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Version == result[j].Version {
			return result[i].CreatedAt.After(result[j].CreatedAt)
		}
		return result[i].Version > result[j].Version
	})
	return result, nil
}

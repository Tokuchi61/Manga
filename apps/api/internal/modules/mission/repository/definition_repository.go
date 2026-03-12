package repository

import (
	"context"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
)

func (s *MemoryStore) UpsertMissionDefinition(_ context.Context, definition entity.MissionDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := normalizeValue(definition.MissionID)
	if existing, ok := s.missionDefinitionsByID[key]; ok {
		definition.CreatedAt = existing.CreatedAt
	}
	if definition.CreatedAt.IsZero() {
		definition.CreatedAt = time.Now().UTC()
	}
	s.missionDefinitionsByID[key] = cloneDefinition(definition)
	return nil
}

func (s *MemoryStore) GetMissionDefinition(_ context.Context, missionID string) (entity.MissionDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	definition, ok := s.missionDefinitionsByID[normalizeValue(missionID)]
	if !ok {
		return entity.MissionDefinition{}, ErrNotFound
	}
	return cloneDefinition(definition), nil
}

func (s *MemoryStore) ListMissionDefinitions(_ context.Context, category string, activeOnly bool, limit int, offset int) ([]entity.MissionDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedCategory := normalizeValue(category)
	now := time.Now().UTC()
	items := make([]entity.MissionDefinition, 0, len(s.missionDefinitionsByID))
	for _, definition := range s.missionDefinitionsByID {
		if normalizedCategory != "" && normalizeValue(definition.Category) != normalizedCategory {
			continue
		}
		if activeOnly {
			if !definition.Active {
				continue
			}
			if definition.StartsAt != nil && now.Before(definition.StartsAt.UTC()) {
				continue
			}
			if definition.EndsAt != nil && now.After(definition.EndsAt.UTC()) {
				continue
			}
		}
		items = append(items, cloneDefinition(definition))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return strings.Compare(items[i].MissionID, items[j].MissionID) < 0
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return applyOffsetLimit(items, offset, limit), nil
}

package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
)

func (s *MemoryStore) UpsertPlacementDefinition(_ context.Context, placement entity.PlacementDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.placementsByID[normalizeValue(placement.PlacementID)] = clonePlacement(placement)
	return nil
}

func (s *MemoryStore) GetPlacementDefinition(_ context.Context, placementID string) (entity.PlacementDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	placement, ok := s.placementsByID[normalizeValue(placementID)]
	if !ok {
		return entity.PlacementDefinition{}, ErrNotFound
	}
	return clonePlacement(placement), nil
}

func (s *MemoryStore) ListPlacementDefinitions(_ context.Context, surface string, visibleOnly bool, limit int, offset int) ([]entity.PlacementDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedSurface := normalizeValue(surface)
	items := make([]entity.PlacementDefinition, 0, len(s.placementsByID))
	for _, placement := range s.placementsByID {
		if normalizedSurface != "" && normalizeValue(placement.Surface) != normalizedSurface {
			continue
		}
		if visibleOnly && !placement.Visible {
			continue
		}
		items = append(items, clonePlacement(placement))
	}

	sortByUpdatedDesc(items, func(i int, j int) bool {
		if items[i].Priority == items[j].Priority {
			return items[i].UpdatedAt.After(items[j].UpdatedAt)
		}
		return items[i].Priority > items[j].Priority
	})

	return applyOffsetLimit(items, offset, limit), nil
}

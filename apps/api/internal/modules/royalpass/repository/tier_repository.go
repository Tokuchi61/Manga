package repository

import (
	"context"
	"sort"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
)

func (s *MemoryStore) UpsertTierDefinition(_ context.Context, tier entity.TierDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := tierKey(tier.SeasonID, tier.TierNumber, tier.Track)
	if existing, ok := s.tierDefinitionsByKey[key]; ok {
		tier.CreatedAt = existing.CreatedAt
	}
	if tier.CreatedAt.IsZero() {
		tier.CreatedAt = tier.UpdatedAt
	}
	s.tierDefinitionsByKey[key] = cloneTier(tier)
	return nil
}

func (s *MemoryStore) GetTierDefinition(_ context.Context, seasonID string, tierNumber int, track string) (entity.TierDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tier, ok := s.tierDefinitionsByKey[tierKey(seasonID, tierNumber, track)]
	if !ok {
		return entity.TierDefinition{}, ErrNotFound
	}
	return cloneTier(tier), nil
}

func (s *MemoryStore) ListTierDefinitions(_ context.Context, seasonID string, track string, activeOnly bool, limit int, offset int) ([]entity.TierDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedSeasonID := normalizeValue(seasonID)
	normalizedTrack := normalizeValue(track)
	items := make([]entity.TierDefinition, 0, len(s.tierDefinitionsByKey))
	for _, tier := range s.tierDefinitionsByKey {
		if normalizedSeasonID != "" && normalizeValue(tier.SeasonID) != normalizedSeasonID {
			continue
		}
		if normalizedTrack != "" && normalizeValue(tier.Track) != normalizedTrack {
			continue
		}
		if activeOnly && !tier.Active {
			continue
		}
		items = append(items, cloneTier(tier))
	}

	sort.Slice(items, func(i int, j int) bool {
		left := items[i]
		right := items[j]
		if left.TierNumber == right.TierNumber {
			if left.Track == right.Track {
				return left.UpdatedAt.After(right.UpdatedAt)
			}
			return left.Track < right.Track
		}
		return left.TierNumber < right.TierNumber
	})

	return applyOffsetLimit(items, offset, limit), nil
}

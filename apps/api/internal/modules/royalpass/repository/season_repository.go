package repository

import (
	"context"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
)

func (s *MemoryStore) UpsertSeasonDefinition(_ context.Context, season entity.SeasonDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := normalizeValue(season.SeasonID)
	if existing, ok := s.seasonDefinitionsByID[key]; ok {
		season.CreatedAt = existing.CreatedAt
	}
	if season.CreatedAt.IsZero() {
		season.CreatedAt = time.Now().UTC()
	}
	s.seasonDefinitionsByID[key] = cloneSeason(season)
	return nil
}

func (s *MemoryStore) GetSeasonDefinition(_ context.Context, seasonID string) (entity.SeasonDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	season, ok := s.seasonDefinitionsByID[normalizeValue(seasonID)]
	if !ok {
		return entity.SeasonDefinition{}, ErrNotFound
	}
	return cloneSeason(season), nil
}

func (s *MemoryStore) ListSeasonDefinitions(_ context.Context, state string, limit int, offset int) ([]entity.SeasonDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedState := normalizeValue(state)
	items := make([]entity.SeasonDefinition, 0, len(s.seasonDefinitionsByID))
	for _, season := range s.seasonDefinitionsByID {
		if normalizedState != "" && normalizeValue(season.State) != normalizedState {
			continue
		}
		items = append(items, cloneSeason(season))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return strings.Compare(items[i].SeasonID, items[j].SeasonID) < 0
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return applyOffsetLimit(items, offset, limit), nil
}

func (s *MemoryStore) ResolveCurrentSeason(_ context.Context, now time.Time) (entity.SeasonDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	utcNow := now.UTC()
	candidates := make([]entity.SeasonDefinition, 0)
	for _, season := range s.seasonDefinitionsByID {
		if normalizeValue(season.State) != entity.SeasonStateActive {
			continue
		}
		if season.StartsAt != nil && utcNow.Before(season.StartsAt.UTC()) {
			continue
		}
		if season.EndsAt != nil && utcNow.After(season.EndsAt.UTC()) {
			continue
		}
		candidates = append(candidates, cloneSeason(season))
	}

	if len(candidates) == 0 {
		return entity.SeasonDefinition{}, ErrNotFound
	}

	sortByTimeDesc(candidates, func(i int, j int) bool {
		left := candidates[i]
		right := candidates[j]
		if left.StartsAt != nil && right.StartsAt != nil && !left.StartsAt.Equal(*right.StartsAt) {
			return left.StartsAt.After(*right.StartsAt)
		}
		if left.UpdatedAt.Equal(right.UpdatedAt) {
			return strings.Compare(left.SeasonID, right.SeasonID) < 0
		}
		return left.UpdatedAt.After(right.UpdatedAt)
	})

	return candidates[0], nil
}

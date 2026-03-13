package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
)

func (s *MemoryStore) CreateImpressionLog(_ context.Context, log entity.ImpressionLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.impressionsByID[normalizeValue(log.ImpressionID)] = cloneImpression(log)
	return nil
}

func (s *MemoryStore) CreateClickLog(_ context.Context, log entity.ClickLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clicksByID[normalizeValue(log.ClickID)] = cloneClick(log)
	return nil
}

func (s *MemoryStore) CountImpressionsBySessionPlacement(_ context.Context, sessionID string, placementID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedSessionID := normalizeValue(sessionID)
	normalizedPlacementID := normalizeValue(placementID)
	count := 0
	for _, log := range s.impressionsByID {
		if normalizedSessionID != "" && normalizeValue(log.SessionID) != normalizedSessionID {
			continue
		}
		if normalizeValue(log.PlacementID) != normalizedPlacementID {
			continue
		}
		if normalizeValue(log.Status) != entity.ImpressionStatusAccepted {
			continue
		}
		count++
	}
	return count, nil
}

func (s *MemoryStore) GetImpressionDedup(_ context.Context, dedupKey string) (entity.ImpressionLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log, ok := s.impressionDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.ImpressionLog{}, ErrNotFound
	}
	return cloneImpression(log), nil
}

func (s *MemoryStore) PutImpressionDedup(_ context.Context, dedupKey string, log entity.ImpressionLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.impressionDedupByKey[normalizeValue(dedupKey)] = cloneImpression(log)
	return nil
}

func (s *MemoryStore) GetClickDedup(_ context.Context, dedupKey string) (entity.ClickLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log, ok := s.clickDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.ClickLog{}, ErrNotFound
	}
	return cloneClick(log), nil
}

func (s *MemoryStore) PutClickDedup(_ context.Context, dedupKey string, log entity.ClickLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clickDedupByKey[normalizeValue(dedupKey)] = cloneClick(log)
	return nil
}

package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
)

func (s *MemoryStore) GetMissionProgress(_ context.Context, userID string, missionID string, periodKey string) (entity.MissionProgress, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	progress, ok := s.missionProgressByKey[missionProgressKey(userID, missionID, periodKey)]
	if !ok {
		return entity.MissionProgress{}, ErrNotFound
	}
	return cloneProgress(progress), nil
}

func (s *MemoryStore) UpsertMissionProgress(_ context.Context, progress entity.MissionProgress) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.missionProgressByKey[missionProgressKey(progress.UserID, progress.MissionID, progress.PeriodKey)] = cloneProgress(progress)
	return nil
}

func (s *MemoryStore) DeleteMissionProgressByUser(_ context.Context, userID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedUserID := normalizeValue(userID)
	deleted := 0
	for key, progress := range s.missionProgressByKey {
		if normalizeValue(progress.UserID) == normalizedUserID {
			delete(s.missionProgressByKey, key)
			deleted++
		}
	}
	return deleted, nil
}

func (s *MemoryStore) DeleteMissionProgressByUserMission(_ context.Context, userID string, missionID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedUserID := normalizeValue(userID)
	normalizedMissionID := normalizeValue(missionID)
	deleted := 0
	for key, progress := range s.missionProgressByKey {
		if normalizeValue(progress.UserID) == normalizedUserID && normalizeValue(progress.MissionID) == normalizedMissionID {
			delete(s.missionProgressByKey, key)
			deleted++
		}
	}
	return deleted, nil
}

func (s *MemoryStore) GetProgressDedup(_ context.Context, dedupKey string) (entity.MissionProgress, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	progress, ok := s.progressDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.MissionProgress{}, ErrNotFound
	}
	return cloneProgress(progress), nil
}

func (s *MemoryStore) PutProgressDedup(_ context.Context, dedupKey string, progress entity.MissionProgress) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.progressDedupByKey[normalizeValue(dedupKey)] = cloneProgress(progress)
	return nil
}

func (s *MemoryStore) GetClaimDedup(_ context.Context, dedupKey string) (entity.MissionProgress, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	progress, ok := s.claimDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.MissionProgress{}, ErrNotFound
	}
	return cloneProgress(progress), nil
}

func (s *MemoryStore) PutClaimDedup(_ context.Context, dedupKey string, progress entity.MissionProgress) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.claimDedupByKey[normalizeValue(dedupKey)] = cloneProgress(progress)
	return nil
}

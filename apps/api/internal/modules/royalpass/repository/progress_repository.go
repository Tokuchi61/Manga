package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
)

func (s *MemoryStore) GetUserProgress(_ context.Context, userID string, seasonID string) (entity.UserProgress, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	progress, ok := s.userProgressByKey[progressKey(userID, seasonID)]
	if !ok {
		return entity.UserProgress{}, ErrNotFound
	}
	return cloneUserProgress(progress), nil
}

func (s *MemoryStore) UpsertUserProgress(_ context.Context, progress entity.UserProgress) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userProgressByKey[progressKey(progress.UserID, progress.SeasonID)] = cloneUserProgress(progress)
	return nil
}

func (s *MemoryStore) DeleteUserProgressByUser(_ context.Context, userID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedUserID := normalizeValue(userID)
	deleted := 0
	for key, progress := range s.userProgressByKey {
		if normalizeValue(progress.UserID) == normalizedUserID {
			delete(s.userProgressByKey, key)
			deleted++
		}
	}
	return deleted, nil
}

func (s *MemoryStore) DeleteUserProgressByUserSeason(_ context.Context, userID string, seasonID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedUserID := normalizeValue(userID)
	normalizedSeasonID := normalizeValue(seasonID)
	deleted := 0
	for key, progress := range s.userProgressByKey {
		if normalizeValue(progress.UserID) == normalizedUserID && normalizeValue(progress.SeasonID) == normalizedSeasonID {
			delete(s.userProgressByKey, key)
			deleted++
		}
	}
	return deleted, nil
}

func (s *MemoryStore) GetProgressDedup(_ context.Context, dedupKey string) (entity.UserProgress, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	progress, ok := s.progressDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.UserProgress{}, ErrNotFound
	}
	return cloneUserProgress(progress), nil
}

func (s *MemoryStore) PutProgressDedup(_ context.Context, dedupKey string, progress entity.UserProgress) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.progressDedupByKey[normalizeValue(dedupKey)] = cloneUserProgress(progress)
	return nil
}

func (s *MemoryStore) GetClaimDedup(_ context.Context, dedupKey string) (entity.UserProgress, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	progress, ok := s.claimDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.UserProgress{}, ErrNotFound
	}
	return cloneUserProgress(progress), nil
}

func (s *MemoryStore) PutClaimDedup(_ context.Context, dedupKey string, progress entity.UserProgress) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.claimDedupByKey[normalizeValue(dedupKey)] = cloneUserProgress(progress)
	return nil
}

func (s *MemoryStore) GetPremiumActivationDedup(_ context.Context, dedupKey string) (entity.UserProgress, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	progress, ok := s.premiumActivationDedupByKey[normalizeValue(dedupKey)]
	if !ok {
		return entity.UserProgress{}, ErrNotFound
	}
	return cloneUserProgress(progress), nil
}

func (s *MemoryStore) PutPremiumActivationDedup(_ context.Context, dedupKey string, progress entity.UserProgress) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.premiumActivationDedupByKey[normalizeValue(dedupKey)] = cloneUserProgress(progress)
	return nil
}

package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
)

func (s *MemoryStore) GetPreferenceByUserID(_ context.Context, userID string) (entity.Preference, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedUserID := normalizeValue(userID)
	preference, exists := s.preferencesByUser[normalizedUserID]
	if !exists {
		return defaultPreference(normalizedUserID, s.runtimeConfig.UpdatedAt), nil
	}
	return clonePreference(preference), nil
}

func (s *MemoryStore) UpsertPreference(_ context.Context, preference entity.Preference) (entity.Preference, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stored := clonePreference(preference)
	stored.UserID = normalizeValue(preference.UserID)
	s.preferencesByUser[stored.UserID] = stored
	return clonePreference(stored), nil
}

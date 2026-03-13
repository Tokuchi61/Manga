package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
)

func (s *MemoryStore) CreateProviderSession(_ context.Context, session entity.ProviderSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.providerSessionsByID[normalizeValue(session.SessionID)] = cloneProviderSession(session)
	return nil
}

func (s *MemoryStore) GetProviderSession(_ context.Context, sessionID string) (entity.ProviderSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.providerSessionsByID[normalizeValue(sessionID)]
	if !ok {
		return entity.ProviderSession{}, ErrNotFound
	}
	return cloneProviderSession(session), nil
}

func (s *MemoryStore) UpdateProviderSession(_ context.Context, session entity.ProviderSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := normalizeValue(session.SessionID)
	if _, ok := s.providerSessionsByID[key]; !ok {
		return ErrNotFound
	}
	s.providerSessionsByID[key] = cloneProviderSession(session)
	return nil
}

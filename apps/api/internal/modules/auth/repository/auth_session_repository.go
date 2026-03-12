package repository

import (
	"context"
	"sort"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/google/uuid"
)

func (s *MemoryStore) CreateSession(_ context.Context, session entity.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessionsByID[session.ID] = session
	if _, ok := s.sessionCredentialIdx[session.CredentialID]; !ok {
		s.sessionCredentialIdx[session.CredentialID] = make(map[uuid.UUID]struct{})
	}
	s.sessionCredentialIdx[session.CredentialID][session.ID] = struct{}{}
	return nil
}

func (s *MemoryStore) GetSessionByID(_ context.Context, sessionID uuid.UUID) (entity.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessionsByID[sessionID]
	if !ok {
		return entity.Session{}, ErrNotFound
	}
	return session, nil
}

func (s *MemoryStore) ListSessionsByCredential(_ context.Context, credentialID uuid.UUID) ([]entity.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessionIDs := s.sessionCredentialIdx[credentialID]
	result := make([]entity.Session, 0, len(sessionIDs))
	for sessionID := range sessionIDs {
		if session, ok := s.sessionsByID[sessionID]; ok {
			result = append(result, session)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

func (s *MemoryStore) UpdateSession(_ context.Context, session entity.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessionsByID[session.ID]; !ok {
		return ErrNotFound
	}
	s.sessionsByID[session.ID] = session
	return nil
}

func (s *MemoryStore) RevokeSession(_ context.Context, sessionID uuid.UUID, revokedAtUnix int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessionsByID[sessionID]
	if !ok {
		return ErrNotFound
	}
	revokedAt := time.Unix(revokedAtUnix, 0).UTC()
	session.RevokedAt = &revokedAt
	session.LastSeenAt = revokedAt
	s.sessionsByID[sessionID] = session
	return nil
}

func (s *MemoryStore) RevokeOtherSessions(_ context.Context, credentialID uuid.UUID, currentSessionID uuid.UUID, revokedAtUnix int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	revokedAt := time.Unix(revokedAtUnix, 0).UTC()
	for sessionID := range s.sessionCredentialIdx[credentialID] {
		if sessionID == currentSessionID {
			continue
		}
		session := s.sessionsByID[sessionID]
		session.RevokedAt = &revokedAt
		session.LastSeenAt = revokedAt
		s.sessionsByID[sessionID] = session
	}
	return nil
}

func (s *MemoryStore) RevokeAllSessions(_ context.Context, credentialID uuid.UUID, revokedAtUnix int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	revokedAt := time.Unix(revokedAtUnix, 0).UTC()
	for sessionID := range s.sessionCredentialIdx[credentialID] {
		session := s.sessionsByID[sessionID]
		session.RevokedAt = &revokedAt
		session.LastSeenAt = revokedAt
		s.sessionsByID[sessionID] = session
	}
	return nil
}

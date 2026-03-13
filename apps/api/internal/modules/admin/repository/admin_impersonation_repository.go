package repository

import (
	"context"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
)

func (s *MemoryStore) CreateImpersonationSession(ctx context.Context, session entity.ImpersonationSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.impersonationByID[session.SessionID] = cloneImpersonationSession(session)
	return nil
}

func (s *MemoryStore) GetImpersonationSession(ctx context.Context, sessionID string) (entity.ImpersonationSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.impersonationByID[sessionID]
	if !ok {
		return entity.ImpersonationSession{}, ErrNotFound
	}
	return cloneImpersonationSession(session), nil
}

func (s *MemoryStore) UpdateImpersonationSession(ctx context.Context, session entity.ImpersonationSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.impersonationByID[session.SessionID]; !ok {
		return ErrNotFound
	}
	s.impersonationByID[session.SessionID] = cloneImpersonationSession(session)
	return nil
}

func (s *MemoryStore) FindActiveImpersonation(ctx context.Context, actorUserID string, targetUserID string) (entity.ImpersonationSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now().UTC()
	for _, session := range s.impersonationByID {
		if !session.Active {
			continue
		}
		if session.ExpiresAt.Before(now) {
			continue
		}
		if actorUserID != "" && normalizeValue(session.ActorUserID) != normalizeValue(actorUserID) {
			continue
		}
		if targetUserID != "" && normalizeValue(session.TargetUserID) != normalizeValue(targetUserID) {
			continue
		}
		return cloneImpersonationSession(session), nil
	}

	return entity.ImpersonationSession{}, ErrNotFound
}

func (s *MemoryStore) ListImpersonationSessions(ctx context.Context, activeOnly bool, limit int, offset int) ([]entity.ImpersonationSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now().UTC()
	filtered := make([]entity.ImpersonationSession, 0, len(s.impersonationByID))
	for _, session := range s.impersonationByID {
		if activeOnly {
			if !session.Active {
				continue
			}
			if session.ExpiresAt.Before(now) {
				continue
			}
		}
		filtered = append(filtered, cloneImpersonationSession(session))
	}

	sortByCreatedAtDesc(filtered, func(i int, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	return applyOffsetLimit(filtered, offset, limit), nil
}

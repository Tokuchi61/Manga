package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/google/uuid"
)

func (s *MemoryStore) AppendSecurityEvent(_ context.Context, event entity.SecurityEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.securityEvents = append(s.securityEvents, event)
	return nil
}

func (s *MemoryStore) ListSecurityEvents(_ context.Context, credentialID uuid.UUID) ([]entity.SecurityEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]entity.SecurityEvent, 0)
	for _, event := range s.securityEvents {
		if event.CredentialID != nil && *event.CredentialID == credentialID {
			result = append(result, event)
		}
	}
	return result, nil
}

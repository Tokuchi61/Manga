package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/google/uuid"
)

func (s *MemoryStore) CreateCredential(_ context.Context, credential entity.Credential) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	email := normalizeEmail(credential.Email)
	if _, exists := s.credentialEmailIdx[email]; exists {
		return ErrConflict
	}

	credential.Email = email
	s.credentialsByID[credential.ID] = credential
	s.credentialEmailIdx[email] = credential.ID
	return nil
}

func (s *MemoryStore) GetCredentialByEmail(_ context.Context, email string) (entity.Credential, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	credentialID, ok := s.credentialEmailIdx[normalizeEmail(email)]
	if !ok {
		return entity.Credential{}, ErrNotFound
	}
	credential, ok := s.credentialsByID[credentialID]
	if !ok {
		return entity.Credential{}, ErrNotFound
	}
	return credential, nil
}

func (s *MemoryStore) GetCredentialByID(_ context.Context, credentialID uuid.UUID) (entity.Credential, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	credential, ok := s.credentialsByID[credentialID]
	if !ok {
		return entity.Credential{}, ErrNotFound
	}
	return credential, nil
}

func (s *MemoryStore) UpdateCredential(_ context.Context, credential entity.Credential) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.credentialsByID[credential.ID]; !ok {
		return ErrNotFound
	}
	credential.Email = normalizeEmail(credential.Email)
	s.credentialsByID[credential.ID] = credential
	s.credentialEmailIdx[credential.Email] = credential.ID
	return nil
}

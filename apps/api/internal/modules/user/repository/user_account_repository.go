package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
)

func (s *MemoryStore) CreateUser(_ context.Context, user entity.UserAccount) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID := normalizeID(user.ID)
	credentialID := normalizeID(user.CredentialID)
	username := normalizeUsername(user.Username)

	if _, exists := s.usersByID[userID]; exists {
		return ErrConflict
	}
	if _, exists := s.credentialUserIndex[credentialID]; exists {
		return ErrConflict
	}
	if _, exists := s.usernameUserIndex[username]; exists {
		return ErrConflict
	}

	user.ID = userID
	user.CredentialID = credentialID
	user.Username = username
	s.usersByID[userID] = user
	s.credentialUserIndex[credentialID] = userID
	s.usernameUserIndex[username] = userID
	return nil
}

func (s *MemoryStore) GetUserByID(_ context.Context, userID string) (entity.UserAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByID[normalizeID(userID)]
	if !ok {
		return entity.UserAccount{}, ErrNotFound
	}
	return user, nil
}

func (s *MemoryStore) GetUserByCredentialID(_ context.Context, credentialID string) (entity.UserAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, ok := s.credentialUserIndex[normalizeID(credentialID)]
	if !ok {
		return entity.UserAccount{}, ErrNotFound
	}
	user, ok := s.usersByID[userID]
	if !ok {
		return entity.UserAccount{}, ErrNotFound
	}
	return user, nil
}

func (s *MemoryStore) UpdateUser(_ context.Context, user entity.UserAccount) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID := normalizeID(user.ID)
	current, ok := s.usersByID[userID]
	if !ok {
		return ErrNotFound
	}

	nextUsername := normalizeUsername(user.Username)
	if nextUsername != normalizeUsername(current.Username) {
		if existingUserID, exists := s.usernameUserIndex[nextUsername]; exists && existingUserID != userID {
			return ErrConflict
		}
		delete(s.usernameUserIndex, normalizeUsername(current.Username))
		s.usernameUserIndex[nextUsername] = userID
	}

	user.ID = userID
	user.CredentialID = normalizeID(user.CredentialID)
	user.Username = nextUsername
	s.usersByID[userID] = user
	s.credentialUserIndex[user.CredentialID] = userID
	return nil
}

package repository

import (
	"context"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
)

func (s *MemoryStore) CreateToken(_ context.Context, token entity.Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	indexKey := tokenIndexKey(token.TokenHash, token.Type)
	if _, exists := s.tokenHashTypeIndex[indexKey]; exists {
		return fmt.Errorf("duplicate token hash")
	}
	s.tokensByID[token.ID] = token
	s.tokenHashTypeIndex[indexKey] = token.ID
	return nil
}

func (s *MemoryStore) GetTokenByHashAndType(_ context.Context, hash string, tokenType entity.TokenType) (entity.Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tokenID, ok := s.tokenHashTypeIndex[tokenIndexKey(hash, tokenType)]
	if !ok {
		return entity.Token{}, ErrNotFound
	}
	token, ok := s.tokensByID[tokenID]
	if !ok {
		return entity.Token{}, ErrNotFound
	}
	return token, nil
}

func (s *MemoryStore) UpdateToken(_ context.Context, token entity.Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tokensByID[token.ID]; !ok {
		return ErrNotFound
	}
	s.tokensByID[token.ID] = token
	s.tokenHashTypeIndex[tokenIndexKey(token.TokenHash, token.Type)] = token.ID
	return nil
}

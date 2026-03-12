package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/google/uuid"
)

func (s *AuthService) issueToken(ctx context.Context, credentialID uuid.UUID, sessionID *uuid.UUID, tokenType entity.TokenType, ttl time.Duration) (string, error) {
	rawToken, err := generateOpaqueToken(40)
	if err != nil {
		return "", err
	}

	now := s.now().UTC()
	token := entity.Token{
		ID:           uuid.New(),
		CredentialID: credentialID,
		SessionID:    sessionID,
		Type:         tokenType,
		TokenHash:    hashOpaqueToken(rawToken),
		ExpiresAt:    now.Add(ttl),
		CreatedAt:    now,
	}
	if err := s.store.CreateToken(ctx, token); err != nil {
		return "", err
	}
	return rawToken, nil
}

func hashOpaqueToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func generateOpaqueToken(size int) (string, error) {
	if size <= 0 {
		size = 32
	}
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

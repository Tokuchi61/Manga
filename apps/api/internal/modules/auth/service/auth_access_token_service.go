package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (s *AuthService) issueAccessToken(ctx context.Context, credentialID string, expiresAt time.Time) (string, string, error) {
	if s == nil {
		return "", "", fmt.Errorf("%w: service nil", ErrValidation)
	}

	userID, err := s.resolveUserID(ctx, credentialID)
	if err != nil {
		return "", "", err
	}

	token, err := identity.IssueAccessTokenWithSecret(s.cfg.AccessTokenSecret, identity.TokenClaims{
		UserID:       userID,
		CredentialID: credentialID,
		ExpiresAt:    expiresAt,
		Issuer:       s.cfg.AccessTokenIssuer,
	})
	if err != nil {
		return "", "", err
	}

	return token, userID, nil
}

package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
)

func (s *AuthService) RefreshToken(ctx context.Context, request dto.RefreshTokenRequest, _ RequestMeta) (dto.RefreshTokenResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	tokenHash := hashOpaqueToken(request.RefreshToken)
	token, err := s.store.GetTokenByHashAndType(ctx, tokenHash, entity.TokenTypeRefresh)
	if err != nil {
		if errors.Is(err, authrepository.ErrNotFound) {
			return dto.RefreshTokenResponse{}, ErrTokenInvalidOrExpired
		}
		return dto.RefreshTokenResponse{}, err
	}

	now := s.now().UTC()
	if token.IsConsumed() || token.IsExpired(now) {
		return dto.RefreshTokenResponse{}, ErrTokenInvalidOrExpired
	}
	if token.SessionID == nil {
		return dto.RefreshTokenResponse{}, ErrTokenInvalidOrExpired
	}

	session, err := s.store.GetSessionByID(ctx, *token.SessionID)
	if err != nil {
		if errors.Is(err, authrepository.ErrNotFound) {
			return dto.RefreshTokenResponse{}, ErrSessionNotFound
		}
		return dto.RefreshTokenResponse{}, err
	}
	if session.IsRevoked() {
		return dto.RefreshTokenResponse{}, ErrSessionNotFound
	}

	token.ConsumedAt = &now
	if err := s.store.UpdateToken(ctx, token); err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	refreshToken, err := s.issueToken(ctx, token.CredentialID, token.SessionID, entity.TokenTypeRefresh, s.cfg.RefreshTokenTTL)
	if err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	session.LastSeenAt = now
	if err := s.store.UpdateSession(ctx, session); err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	accessTokenExpiresAt := now.Add(s.cfg.AccessTokenTTL)
	accessToken, userID, err := s.issueAccessToken(ctx, token.CredentialID.String(), accessTokenExpiresAt)
	if err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	return dto.RefreshTokenResponse{
		UserID:               userID,
		SessionID:            session.ID.String(),
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
	}, nil
}

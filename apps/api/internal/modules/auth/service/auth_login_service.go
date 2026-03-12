package service

import (
	"context"
	"errors"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/events"
	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/crypto/password"
	"github.com/google/uuid"
)

func (s *AuthService) Login(ctx context.Context, request dto.LoginRequest, meta RequestMeta) (dto.LoginResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.LoginResponse{}, err
	}

	meta = normalizeMeta(meta)
	credential, err := s.store.GetCredentialByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, authrepository.ErrNotFound) {
			return dto.LoginResponse{}, ErrInvalidCredentials
		}
		return dto.LoginResponse{}, err
	}

	now := s.now().UTC()
	if credential.Banned {
		s.appendSecurityEvent(ctx, credential.ID, events.EventLoginFailed, "rejected", "credential_banned", meta)
		return dto.LoginResponse{}, ErrCredentialBanned
	}
	if credential.Suspended {
		s.appendSecurityEvent(ctx, credential.ID, events.EventLoginFailed, "rejected", "credential_suspended", meta)
		return dto.LoginResponse{}, ErrCredentialSuspended
	}
	if !credential.EmailVerified {
		s.appendSecurityEvent(ctx, credential.ID, events.EventLoginFailed, "rejected", "email_not_verified", meta)
		return dto.LoginResponse{}, ErrEmailNotVerified
	}
	if credential.LoginCooldownUntil != nil && now.Before(*credential.LoginCooldownUntil) {
		s.appendSecurityEvent(ctx, credential.ID, events.EventSecuritySuspiciousLogin, "rejected", "login_cooldown_active", meta)
		return dto.LoginResponse{}, ErrLoginCooldownActive
	}

	ok, err := password.Verify(credential.PasswordHash, request.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if !ok {
		if credential.FailedLoginWindowStartedAt == nil || now.Sub(*credential.FailedLoginWindowStartedAt) >= time.Minute {
			windowStart := now
			credential.FailedLoginWindowStartedAt = &windowStart
			credential.FailedLoginAttempts = 0
		}
		credential.FailedLoginAttempts++
		if credential.FailedLoginAttempts >= s.cfg.FailedAttemptLimit {
			cooldownUntil := now.Add(s.cfg.LoginCooldown)
			credential.LoginCooldownUntil = &cooldownUntil
			credential.FailedLoginAttempts = 0
			credential.FailedLoginWindowStartedAt = nil
			s.appendSecurityEvent(ctx, credential.ID, events.EventSecuritySuspiciousLogin, "rejected", "failed_login_limit_reached", meta)
		}
		credential.UpdatedAt = now
		if err := s.store.UpdateCredential(ctx, credential); err != nil {
			return dto.LoginResponse{}, err
		}
		s.appendSecurityEvent(ctx, credential.ID, events.EventLoginFailed, "failed", "invalid_password", meta)
		if credential.LoginCooldownUntil != nil && now.Before(*credential.LoginCooldownUntil) {
			return dto.LoginResponse{}, ErrLoginCooldownActive
		}
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	credential.FailedLoginAttempts = 0
	credential.FailedLoginWindowStartedAt = nil
	credential.LoginCooldownUntil = nil
	credential.UpdatedAt = now
	if err := s.store.UpdateCredential(ctx, credential); err != nil {
		return dto.LoginResponse{}, err
	}

	session := entity.Session{
		ID:           uuid.New(),
		CredentialID: credential.ID,
		Device:       coalesce(request.Device, meta.Device),
		IP:           coalesce(request.IP, meta.IP),
		CreatedAt:    now,
		LastSeenAt:   now,
	}
	if err := s.store.CreateSession(ctx, session); err != nil {
		return dto.LoginResponse{}, err
	}

	refreshToken, err := s.issueToken(ctx, credential.ID, &session.ID, entity.TokenTypeRefresh, s.cfg.RefreshTokenTTL)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	accessToken, err := generateOpaqueToken(32)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	s.appendSecurityEvent(ctx, credential.ID, events.EventLoginSucceeded, "success", "login_succeeded", meta)

	return dto.LoginResponse{
		CredentialID:         credential.ID.String(),
		SessionID:            session.ID.String(),
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenExpiresAt: now.Add(s.cfg.AccessTokenTTL),
	}, nil
}

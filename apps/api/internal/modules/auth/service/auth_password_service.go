package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/crypto/password"
	"github.com/google/uuid"
)

func (s *AuthService) ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest, meta RequestMeta) (dto.TokenDispatchResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.TokenDispatchResponse{}, err
	}

	credential, err := s.store.GetCredentialByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, authrepository.ErrNotFound) {
			return dto.TokenDispatchResponse{Status: "requested"}, nil
		}
		return dto.TokenDispatchResponse{}, err
	}

	token, err := s.issueToken(ctx, credential.ID, nil, entity.TokenTypePasswordReset, s.cfg.PasswordResetTokenTTL)
	if err != nil {
		return dto.TokenDispatchResponse{}, err
	}

	meta = normalizeMeta(meta)
	s.appendSecurityEvent(ctx, credential.ID, "auth.password_reset.requested", "success", "password_reset_requested", meta)
	return dto.TokenDispatchResponse{Status: "requested", Token: token}, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest, meta RequestMeta) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	tokenHash := hashOpaqueToken(request.Token)
	token, err := s.store.GetTokenByHashAndType(ctx, tokenHash, entity.TokenTypePasswordReset)
	if err != nil {
		if errors.Is(err, authrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrTokenInvalidOrExpired
		}
		return dto.OperationResponse{}, err
	}

	now := s.now().UTC()
	if token.IsConsumed() || token.IsExpired(now) {
		return dto.OperationResponse{}, ErrTokenInvalidOrExpired
	}

	credential, err := s.store.GetCredentialByID(ctx, token.CredentialID)
	if err != nil {
		return dto.OperationResponse{}, err
	}

	hash, err := password.Hash(request.NewPassword)
	if err != nil {
		return dto.OperationResponse{}, err
	}
	credential.PasswordHash = hash
	credential.UpdatedAt = now
	if err := s.store.UpdateCredential(ctx, credential); err != nil {
		return dto.OperationResponse{}, err
	}

	token.ConsumedAt = &now
	if err := s.store.UpdateToken(ctx, token); err != nil {
		return dto.OperationResponse{}, err
	}
	if err := s.store.RevokeAllSessions(ctx, credential.ID, now.Unix()); err != nil {
		return dto.OperationResponse{}, err
	}

	meta = normalizeMeta(meta)
	s.appendSecurityEvent(ctx, credential.ID, "auth.password_reset.applied", "success", "password_reset_applied", meta)
	return dto.OperationResponse{Status: "password_reset"}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, request dto.ChangePasswordRequest, meta RequestMeta) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	credentialID, err := uuid.Parse(request.CredentialID)
	if err != nil {
		return dto.OperationResponse{}, fmt.Errorf("%w: invalid credential id", ErrValidation)
	}

	credential, err := s.store.GetCredentialByID(ctx, credentialID)
	if err != nil {
		if errors.Is(err, authrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrInvalidCredentials
		}
		return dto.OperationResponse{}, err
	}

	ok, err := password.Verify(credential.PasswordHash, request.OldPassword)
	if err != nil {
		return dto.OperationResponse{}, err
	}
	if !ok {
		return dto.OperationResponse{}, ErrInvalidCredentials
	}

	now := s.now().UTC()
	hash, err := password.Hash(request.NewPassword)
	if err != nil {
		return dto.OperationResponse{}, err
	}
	credential.PasswordHash = hash
	credential.UpdatedAt = now
	if err := s.store.UpdateCredential(ctx, credential); err != nil {
		return dto.OperationResponse{}, err
	}
	if err := s.store.RevokeAllSessions(ctx, credential.ID, now.Unix()); err != nil {
		return dto.OperationResponse{}, err
	}

	meta = normalizeMeta(meta)
	s.appendSecurityEvent(ctx, credential.ID, "auth.password.changed", "success", "password_changed", meta)
	return dto.OperationResponse{Status: "password_changed"}, nil
}

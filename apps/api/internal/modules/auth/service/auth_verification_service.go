package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/events"
	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	"github.com/google/uuid"
)

func (s *AuthService) SendEmailVerification(ctx context.Context, request dto.SendVerificationRequest, meta RequestMeta) (dto.TokenDispatchResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.TokenDispatchResponse{}, err
	}

	credentialID, err := uuid.Parse(request.CredentialID)
	if err != nil {
		return dto.TokenDispatchResponse{}, fmt.Errorf("%w: invalid credential id", ErrValidation)
	}

	credential, err := s.store.GetCredentialByID(ctx, credentialID)
	if err != nil {
		if errors.Is(err, authrepository.ErrNotFound) {
			return dto.TokenDispatchResponse{}, ErrInvalidCredentials
		}
		return dto.TokenDispatchResponse{}, err
	}
	if credential.EmailVerified {
		return dto.TokenDispatchResponse{Status: "already_verified"}, nil
	}

	now := s.now().UTC()
	if credential.VerificationResendAvailableAt != nil && now.Before(*credential.VerificationResendAvailableAt) {
		return dto.TokenDispatchResponse{}, ErrVerificationResendCooldown
	}

	token, err := s.issueToken(ctx, credential.ID, nil, entity.TokenTypeEmailVerification, s.cfg.EmailVerificationTokenTTL)
	if err != nil {
		return dto.TokenDispatchResponse{}, err
	}

	nextResend := now.Add(s.cfg.VerificationResendCooldown)
	credential.VerificationResendAvailableAt = &nextResend
	credential.UpdatedAt = now
	if err := s.store.UpdateCredential(ctx, credential); err != nil {
		return dto.TokenDispatchResponse{}, err
	}

	meta = normalizeMeta(meta)
	s.appendSecurityEvent(ctx, credential.ID, events.EventEmailVerificationSent, "success", "verification_resend", meta)
	return dto.TokenDispatchResponse{Status: "sent", Token: token}, nil
}

func (s *AuthService) ConfirmEmailVerification(ctx context.Context, request dto.ConfirmVerificationRequest, meta RequestMeta) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	tokenHash := hashOpaqueToken(request.Token)
	token, err := s.store.GetTokenByHashAndType(ctx, tokenHash, entity.TokenTypeEmailVerification)
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
	credential.EmailVerified = true
	credential.UpdatedAt = now
	if err := s.store.UpdateCredential(ctx, credential); err != nil {
		return dto.OperationResponse{}, err
	}

	token.ConsumedAt = &now
	if err := s.store.UpdateToken(ctx, token); err != nil {
		return dto.OperationResponse{}, err
	}

	meta = normalizeMeta(meta)
	s.appendSecurityEvent(ctx, credential.ID, "auth.email_verification.confirmed", "success", "email_verified", meta)
	return dto.OperationResponse{Status: "verified"}, nil
}

package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/events"
	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/crypto/password"
	"github.com/google/uuid"
)

func (s *AuthService) Register(ctx context.Context, request dto.RegisterRequest, meta RequestMeta) (dto.RegisterResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RegisterResponse{}, err
	}

	if _, err := s.store.GetCredentialByEmail(ctx, request.Email); err == nil {
		return dto.RegisterResponse{}, ErrCredentialAlreadyExists
	} else if !errors.Is(err, authrepository.ErrNotFound) {
		return dto.RegisterResponse{}, err
	}

	hash, err := password.Hash(request.Password)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	now := s.now().UTC()
	credential := entity.Credential{
		ID:                  uuid.New(),
		Email:               request.Email,
		PasswordHash:        hash,
		EmailVerified:       false,
		FailedLoginAttempts: 0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.store.CreateCredential(ctx, credential); err != nil {
		if errors.Is(err, authrepository.ErrConflict) {
			return dto.RegisterResponse{}, ErrCredentialAlreadyExists
		}
		return dto.RegisterResponse{}, err
	}

	verificationToken, err := s.issueToken(ctx, credential.ID, nil, entity.TokenTypeEmailVerification, s.cfg.EmailVerificationTokenTTL)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	nextResend := now.Add(s.cfg.VerificationResendCooldown)
	credential.VerificationResendAvailableAt = &nextResend
	credential.UpdatedAt = now
	if err := s.store.UpdateCredential(ctx, credential); err != nil {
		return dto.RegisterResponse{}, err
	}

	meta = normalizeMeta(meta)
	s.appendSecurityEvent(ctx, credential.ID, events.EventEmailVerificationSent, "success", "register_initial_verification", meta)

	return dto.RegisterResponse{
		CredentialID:      credential.ID.String(),
		Email:             credential.Email,
		EmailVerified:     credential.EmailVerified,
		VerificationToken: verificationToken,
	}, nil
}

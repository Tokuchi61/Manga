package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/events"
	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	"github.com/google/uuid"
)

func (s *AuthService) Logout(ctx context.Context, request dto.LogoutRequest, meta RequestMeta) (dto.OperationResponse, error) {
	_, err := s.RevokeCurrentSession(ctx, dto.RevokeCurrentSessionRequest{
		CredentialID: request.CredentialID,
		SessionID:    request.SessionID,
	}, meta)
	if err != nil {
		return dto.OperationResponse{}, err
	}
	return dto.OperationResponse{Status: "logged_out"}, nil
}

func (s *AuthService) ListSessions(ctx context.Context, request dto.ListSessionsRequest) (dto.ListSessionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListSessionsResponse{}, err
	}

	credentialID, err := uuid.Parse(request.CredentialID)
	if err != nil {
		return dto.ListSessionsResponse{}, fmt.Errorf("%w: invalid credential id", ErrValidation)
	}

	sessions, err := s.store.ListSessionsByCredential(ctx, credentialID)
	if err != nil {
		return dto.ListSessionsResponse{}, err
	}

	result := make([]dto.SessionInfo, 0, len(sessions))
	for _, session := range sessions {
		result = append(result, dto.SessionInfo{
			SessionID:  session.ID.String(),
			Device:     session.Device,
			IP:         session.IP,
			CreatedAt:  session.CreatedAt,
			LastSeenAt: session.LastSeenAt,
			RevokedAt:  session.RevokedAt,
		})
	}

	return dto.ListSessionsResponse{Sessions: result}, nil
}

func (s *AuthService) RevokeCurrentSession(ctx context.Context, request dto.RevokeCurrentSessionRequest, meta RequestMeta) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	credentialID, sessionID, err := parseCredentialAndSession(request.CredentialID, request.SessionID)
	if err != nil {
		return dto.OperationResponse{}, err
	}

	session, err := s.store.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, authrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrSessionNotFound
		}
		return dto.OperationResponse{}, err
	}
	if session.CredentialID != credentialID {
		return dto.OperationResponse{}, ErrSessionNotFound
	}

	now := s.now().UTC()
	if !session.IsRevoked() {
		if err := s.store.RevokeSession(ctx, sessionID, now.Unix()); err != nil {
			return dto.OperationResponse{}, err
		}
	}

	meta = normalizeMeta(meta)
	s.appendSecurityEvent(ctx, credentialID, events.EventSessionRevoked, "success", "revoke_current", meta)
	return dto.OperationResponse{Status: "revoked_current"}, nil
}

func (s *AuthService) RevokeOtherSessions(ctx context.Context, request dto.RevokeOtherSessionsRequest, meta RequestMeta) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	credentialID, sessionID, err := parseCredentialAndSession(request.CredentialID, request.SessionID)
	if err != nil {
		return dto.OperationResponse{}, err
	}

	now := s.now().UTC()
	if err := s.store.RevokeOtherSessions(ctx, credentialID, sessionID, now.Unix()); err != nil {
		return dto.OperationResponse{}, err
	}

	meta = normalizeMeta(meta)
	s.appendSecurityEvent(ctx, credentialID, events.EventSessionRevoked, "success", "revoke_others", meta)
	return dto.OperationResponse{Status: "revoked_others"}, nil
}

func (s *AuthService) RevokeAllSessions(ctx context.Context, request dto.RevokeAllSessionsRequest, meta RequestMeta) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	credentialID, err := uuid.Parse(request.CredentialID)
	if err != nil {
		return dto.OperationResponse{}, fmt.Errorf("%w: invalid credential id", ErrValidation)
	}

	now := s.now().UTC()
	if err := s.store.RevokeAllSessions(ctx, credentialID, now.Unix()); err != nil {
		return dto.OperationResponse{}, err
	}

	meta = normalizeMeta(meta)
	s.appendSecurityEvent(ctx, credentialID, events.EventSessionRevoked, "success", "revoke_all", meta)
	return dto.OperationResponse{Status: "revoked_all"}, nil
}

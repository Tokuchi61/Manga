package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
	adminrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/repository"
	"github.com/google/uuid"
)

func (s *AdminService) StartImpersonation(ctx context.Context, actorUserID string, request dto.StartImpersonationRequest) (dto.ImpersonationSessionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ImpersonationSessionResponse{}, err
	}
	if strings.TrimSpace(actorUserID) == "" {
		return dto.ImpersonationSessionResponse{}, ErrForbiddenAction
	}
	if !requiresDoubleConfirmation(request.RiskLevel) {
		return dto.ImpersonationSessionResponse{}, ErrValidation
	}
	if err := requireDoubleConfirmation(request.RiskLevel, request.DoubleConfirmed, request.ConfirmationToken); err != nil {
		return dto.ImpersonationSessionResponse{}, err
	}

	dedupKey := buildActionDedupKey("impersonation_start", request.RequestID)
	dedupAction, err := s.store.GetActionDedup(ctx, dedupKey)
	if err == nil {
		sessionID := dedupAction.Metadata["session_id"]
		session, getErr := s.store.GetImpersonationSession(ctx, sessionID)
		if getErr != nil {
			return dto.ImpersonationSessionResponse{}, getErr
		}
		return toImpersonationResponse("idempotent", dedupAction.ActionID, session), nil
	}
	if err != nil && !errors.Is(err, adminrepository.ErrNotFound) {
		return dto.ImpersonationSessionResponse{}, err
	}

	if _, err := s.store.FindActiveImpersonation(ctx, actorUserID, request.TargetUserID); err == nil {
		return dto.ImpersonationSessionResponse{}, ErrConflict
	} else if !errors.Is(err, adminrepository.ErrNotFound) {
		return dto.ImpersonationSessionResponse{}, err
	}

	durationMinutes := request.DurationMinutes
	if durationMinutes <= 0 {
		durationMinutes = 15
	}

	now := s.now().UTC()
	sessionID := uuid.NewString()
	action := entity.AdminAction{
		ActionID:                   uuid.NewString(),
		ActionType:                 entity.ActionTypeImpersonationStart,
		ActorUserID:                actorUserID,
		TargetUserID:               request.TargetUserID,
		TargetModule:               "auth",
		TargetType:                 "impersonation",
		TargetID:                   request.TargetUserID,
		RequestID:                  request.RequestID,
		CorrelationID:              request.CorrelationID,
		Reason:                     request.Reason,
		Result:                     entity.ActionResultSuccess,
		RiskLevel:                  request.RiskLevel,
		RequiresDoubleConfirmation: true,
		DoubleConfirmed:            request.DoubleConfirmed,
		ConfirmationToken:          request.ConfirmationToken,
		Metadata: map[string]string{
			"session_id":       sessionID,
			"duration_minutes": strconv.Itoa(durationMinutes),
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	session := entity.ImpersonationSession{
		SessionID:    sessionID,
		ActionID:     action.ActionID,
		ActorUserID:  actorUserID,
		TargetUserID: request.TargetUserID,
		Reason:       request.Reason,
		RiskLevel:    request.RiskLevel,
		Active:       true,
		StartedAt:    now,
		ExpiresAt:    now.Add(time.Duration(durationMinutes) * time.Minute),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.store.CreateAdminAction(ctx, action); err != nil {
		return dto.ImpersonationSessionResponse{}, err
	}
	if err := s.store.CreateImpersonationSession(ctx, session); err != nil {
		return dto.ImpersonationSessionResponse{}, err
	}
	if err := s.store.PutActionDedup(ctx, dedupKey, action); err != nil {
		return dto.ImpersonationSessionResponse{}, err
	}

	return toImpersonationResponse("started", action.ActionID, session), nil
}

func (s *AdminService) StopImpersonation(ctx context.Context, actorUserID string, request dto.StopImpersonationRequest) (dto.ImpersonationSessionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ImpersonationSessionResponse{}, err
	}
	if strings.TrimSpace(actorUserID) == "" {
		return dto.ImpersonationSessionResponse{}, ErrForbiddenAction
	}
	if !requiresDoubleConfirmation(request.RiskLevel) {
		return dto.ImpersonationSessionResponse{}, ErrValidation
	}
	if err := requireDoubleConfirmation(request.RiskLevel, request.DoubleConfirmed, request.ConfirmationToken); err != nil {
		return dto.ImpersonationSessionResponse{}, err
	}

	dedupKey := buildActionDedupKey("impersonation_stop", request.RequestID)
	dedupAction, err := s.store.GetActionDedup(ctx, dedupKey)
	if err == nil {
		sessionID := dedupAction.Metadata["session_id"]
		session, getErr := s.store.GetImpersonationSession(ctx, sessionID)
		if getErr != nil {
			return dto.ImpersonationSessionResponse{}, getErr
		}
		return toImpersonationResponse("idempotent", dedupAction.ActionID, session), nil
	}
	if err != nil && !errors.Is(err, adminrepository.ErrNotFound) {
		return dto.ImpersonationSessionResponse{}, err
	}

	session, err := s.store.GetImpersonationSession(ctx, request.SessionID)
	if err != nil {
		if errors.Is(err, adminrepository.ErrNotFound) {
			return dto.ImpersonationSessionResponse{}, ErrNotFound
		}
		return dto.ImpersonationSessionResponse{}, err
	}
	if !session.Active {
		return dto.ImpersonationSessionResponse{}, ErrConflict
	}

	now := s.now().UTC()
	action := entity.AdminAction{
		ActionID:                   uuid.NewString(),
		ActionType:                 entity.ActionTypeImpersonationStop,
		ActorUserID:                actorUserID,
		TargetUserID:               session.TargetUserID,
		TargetModule:               "auth",
		TargetType:                 "impersonation",
		TargetID:                   session.TargetUserID,
		RequestID:                  request.RequestID,
		CorrelationID:              request.CorrelationID,
		Reason:                     request.Reason,
		Result:                     entity.ActionResultSuccess,
		RiskLevel:                  request.RiskLevel,
		RequiresDoubleConfirmation: true,
		DoubleConfirmed:            request.DoubleConfirmed,
		ConfirmationToken:          request.ConfirmationToken,
		Metadata: map[string]string{
			"session_id": session.SessionID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	session.Active = false
	session.EndedAt = &now
	session.UpdatedAt = now

	if err := s.store.CreateAdminAction(ctx, action); err != nil {
		return dto.ImpersonationSessionResponse{}, err
	}
	if err := s.store.UpdateImpersonationSession(ctx, session); err != nil {
		if errors.Is(err, adminrepository.ErrNotFound) {
			return dto.ImpersonationSessionResponse{}, ErrNotFound
		}
		return dto.ImpersonationSessionResponse{}, err
	}
	if err := s.store.PutActionDedup(ctx, dedupKey, action); err != nil {
		return dto.ImpersonationSessionResponse{}, err
	}

	return toImpersonationResponse("stopped", action.ActionID, session), nil
}

func (s *AdminService) ListImpersonationSessions(ctx context.Context, request dto.ListImpersonationSessionsRequest) (dto.ListImpersonationSessionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListImpersonationSessionsResponse{}, err
	}

	sessions, err := s.store.ListImpersonationSessions(ctx, request.ActiveOnly, request.Limit, request.Offset)
	if err != nil {
		return dto.ListImpersonationSessionsResponse{}, err
	}

	items := make([]dto.ImpersonationSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, toImpersonationResponse("listed", session.ActionID, session))
	}

	return dto.ListImpersonationSessionsResponse{Items: items, Count: len(items)}, nil
}

package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
	moderationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/repository"
	"github.com/google/uuid"
)

func (s *ModerationService) EscalateCase(ctx context.Context, request dto.EscalateCaseRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	caseID, err := parseID(request.CaseID, "case_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	actorID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	reason, err := sanitizeText(request.Reason, "reason")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	moderationCase, err := s.store.GetCaseByID(ctx, caseID)
	if err != nil {
		if errors.Is(err, moderationrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCaseNotFound
		}
		return dto.OperationResponse{}, err
	}
	if moderationCase.IsTerminal() {
		return dto.OperationResponse{}, ErrInvalidStateTransition
	}
	if moderationCase.EscalationStatus == entity.EscalationStatusPendingAdmin || moderationCase.EscalationStatus == entity.EscalationStatusEscalated {
		return dto.OperationResponse{}, ErrAlreadyEscalated
	}

	now := s.now().UTC()
	moderationCase.Status = entity.CaseStatusEscalated
	moderationCase.AssignmentStatus = entity.AssignmentStatusHandoffPending
	moderationCase.AssignedModeratorUserID = nil
	moderationCase.EscalationStatus = entity.EscalationStatusPendingAdmin
	moderationCase.EscalationReason = reason
	moderationCase.EscalatedAt = &now
	moderationCase.UpdatedAt = now
	moderationCase.Actions = append(moderationCase.Actions, entity.CaseAction{
		ID:           uuid.NewString(),
		CaseID:        moderationCase.ID,
		ActorUserID:  actorID,
		ActionType:   entity.ActionTypeEscalate,
		ReasonCode:   "escalated_to_admin",
		Summary:      reason,
		ActionResult: entity.ActionResultNoAction,
		CreatedAt:    now,
	})

	if err := s.store.UpdateCase(ctx, moderationCase); err != nil {
		if errors.Is(err, moderationrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCaseNotFound
		}
		return dto.OperationResponse{}, err
	}
	return dto.OperationResponse{Status: "escalated"}, nil
}

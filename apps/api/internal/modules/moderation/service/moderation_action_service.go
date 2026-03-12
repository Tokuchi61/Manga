package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
	moderationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/repository"
	"github.com/google/uuid"
)

func (s *ModerationService) ApplyAction(ctx context.Context, request dto.ApplyActionRequest) (dto.OperationResponse, error) {
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
	actionType, err := parseActionType(request.ActionType)
	if err != nil {
		return dto.OperationResponse{}, err
	}
	summary := strings.TrimSpace(request.Summary)
	if summary == "" {
		summary = string(actionType)
	} else {
		summary, err = sanitizeText(summary, "summary")
		if err != nil {
			return dto.OperationResponse{}, err
		}
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

	now := s.now().UTC()
	result := mapActionResult(actionType)
	moderationCase.Actions = append(moderationCase.Actions, entity.CaseAction{
		ID:           uuid.NewString(),
		CaseID:        moderationCase.ID,
		ActorUserID:  actorID,
		ActionType:   actionType,
		ReasonCode:   normalizeValue(request.ReasonCode),
		Summary:      summary,
		ActionResult: result,
		CreatedAt:    now,
	})
	moderationCase.ActionResult = result
	moderationCase.LastActionAt = &now

	switch actionType {
	case entity.ActionTypeReviewComplete:
		moderationCase.Status = entity.CaseStatusResolved
		moderationCase.AssignmentStatus = entity.AssignmentStatusReleased
		moderationCase.AssignedModeratorUserID = nil
		if moderationCase.EscalationStatus == entity.EscalationStatusPendingAdmin {
			moderationCase.EscalationStatus = entity.EscalationStatusResolved
		}
	default:
		if moderationCase.Status == entity.CaseStatusNew || moderationCase.Status == entity.CaseStatusQueued || moderationCase.Status == entity.CaseStatusAssigned {
			moderationCase.Status = entity.CaseStatusInReview
		}
	}

	if moderationCase.AssignmentStatus == entity.AssignmentStatusUnassigned || moderationCase.AssignmentStatus == entity.AssignmentStatusReleased {
		moderationCase.AssignmentStatus = entity.AssignmentStatusAssigned
		moderationCase.AssignedModeratorUserID = &actorID
	}
	moderationCase.UpdatedAt = now

	if err := s.store.UpdateCase(ctx, moderationCase); err != nil {
		if errors.Is(err, moderationrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCaseNotFound
		}
		return dto.OperationResponse{}, err
	}
	return dto.OperationResponse{Status: "action_applied"}, nil
}

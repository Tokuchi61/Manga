package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
	moderationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/repository"
)

func (s *ModerationService) AssignCase(ctx context.Context, request dto.AssignCaseRequest) (dto.OperationResponse, error) {
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
	assigneeID := actorID
	if request.AssigneeUserID != "" {
		assigneeID, err = parseID(request.AssigneeUserID, "assignee_user_id")
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
	if moderationCase.EscalationStatus == entity.EscalationStatusEscalated || moderationCase.EscalationStatus == entity.EscalationStatusPendingAdmin {
		return dto.OperationResponse{}, ErrInvalidStateTransition
	}

	now := s.now().UTC()
	moderationCase.AssignedModeratorUserID = &assigneeID
	moderationCase.AssignmentStatus = entity.AssignmentStatusAssigned
	if moderationCase.Status == entity.CaseStatusNew || moderationCase.Status == entity.CaseStatusQueued {
		moderationCase.Status = entity.CaseStatusAssigned
	}
	moderationCase.UpdatedAt = now

	if err := s.store.UpdateCase(ctx, moderationCase); err != nil {
		if errors.Is(err, moderationrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCaseNotFound
		}
		if errors.Is(err, moderationrepository.ErrConflict) {
			return dto.OperationResponse{}, ErrCaseAlreadyExists
		}
		return dto.OperationResponse{}, err
	}
	return dto.OperationResponse{Status: "assigned"}, nil
}

func (s *ModerationService) ReleaseCase(ctx context.Context, request dto.ReleaseCaseRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	caseID, err := parseID(request.CaseID, "case_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	_, err = parseID(request.ActorUserID, "actor_user_id")
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

	now := s.now().UTC()
	moderationCase.AssignedModeratorUserID = nil
	moderationCase.AssignmentStatus = entity.AssignmentStatusReleased
	if moderationCase.Status == entity.CaseStatusAssigned || moderationCase.Status == entity.CaseStatusInReview {
		moderationCase.Status = entity.CaseStatusQueued
	}
	moderationCase.UpdatedAt = now

	if err := s.store.UpdateCase(ctx, moderationCase); err != nil {
		if errors.Is(err, moderationrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCaseNotFound
		}
		return dto.OperationResponse{}, err
	}
	return dto.OperationResponse{Status: "released"}, nil
}

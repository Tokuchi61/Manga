package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
	moderationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/repository"
	"github.com/google/uuid"
)

func (s *ModerationService) AddModeratorNote(ctx context.Context, request dto.AddModeratorNoteRequest) (dto.OperationResponse, error) {
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
	noteBody, err := sanitizeText(request.Body, "body")
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
	moderationCase.Notes = append(moderationCase.Notes, entity.ModeratorNote{
		ID:           uuid.NewString(),
		CaseID:        moderationCase.ID,
		AuthorUserID: actorID,
		Body:          noteBody,
		InternalOnly:  request.InternalOnly,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if moderationCase.Status == entity.CaseStatusNew || moderationCase.Status == entity.CaseStatusQueued || moderationCase.Status == entity.CaseStatusAssigned {
		moderationCase.Status = entity.CaseStatusInReview
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
	return dto.OperationResponse{Status: "note_added"}, nil
}

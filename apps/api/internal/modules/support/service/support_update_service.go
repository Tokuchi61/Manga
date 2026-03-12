package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
	supportrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/repository"
	"github.com/google/uuid"
)

func (s *SupportService) AddReply(ctx context.Context, request dto.AddSupportReplyRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	supportID, err := parseID(request.SupportID, "support_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	actorID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	visibility, err := parseReplyVisibility(request.Visibility)
	if err != nil {
		return dto.OperationResponse{}, err
	}
	message, err := sanitizeText(request.Message, "message")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	supportCase, err := s.store.GetCaseByID(ctx, supportID)
	if err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrSupportNotFound
		}
		return dto.OperationResponse{}, err
	}

	if supportCase.Status == entity.SupportStatusClosed || supportCase.Status == entity.SupportStatusSpam {
		return dto.OperationResponse{}, ErrInvalidStateTransition
	}
	if !request.ActorIsTeam && actorID != supportCase.RequesterUserID {
		return dto.OperationResponse{}, ErrForbiddenAction
	}
	if visibility == entity.ReplyVisibilityInternalOnly && actorID == supportCase.RequesterUserID {
		return dto.OperationResponse{}, ErrForbiddenAction
	}

	now := s.now().UTC()
	supportCase.Replies = append(supportCase.Replies, entity.SupportReply{
		ID:            uuid.NewString(),
		SupportID:     supportCase.ID,
		AuthorUserID:  actorID,
		Message:       strings.TrimSpace(request.Message),
		SanitizedBody: message,
		Visibility:    visibility,
		CreatedAt:     now,
		UpdatedAt:     now,
	})

	if visibility == entity.ReplyVisibilityPublicToRequester {
		if actorID == supportCase.RequesterUserID {
			supportCase.Status = entity.SupportStatusWaitingTeam
		} else {
			supportCase.Status = entity.SupportStatusWaitingUser
		}
	}

	supportCase.UpdatedAt = now
	if err := s.store.UpdateCase(ctx, supportCase); err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrSupportNotFound
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "reply_added"}, nil
}

func (s *SupportService) UpdateStatus(ctx context.Context, request dto.UpdateSupportStatusRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	supportID, err := parseID(request.SupportID, "support_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	status, err := parseStatus(request.Status)
	if err != nil {
		return dto.OperationResponse{}, err
	}

	supportCase, err := s.store.GetCaseByID(ctx, supportID)
	if err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrSupportNotFound
		}
		return dto.OperationResponse{}, err
	}
	if !isValidStatusTransition(supportCase.Status, status) {
		return dto.OperationResponse{}, ErrInvalidStateTransition
	}

	if request.AssigneeUserID != nil {
		if strings.TrimSpace(*request.AssigneeUserID) == "" {
			supportCase.AssigneeUserID = nil
		} else {
			assigneeID, parseErr := parseID(*request.AssigneeUserID, "assignee_user_id")
			if parseErr != nil {
				return dto.OperationResponse{}, parseErr
			}
			supportCase.AssigneeUserID = &assigneeID
		}
	}
	if request.ReviewedByUserID != nil {
		if strings.TrimSpace(*request.ReviewedByUserID) == "" {
			supportCase.ReviewedByUserID = nil
		} else {
			reviewerID, parseErr := parseID(*request.ReviewedByUserID, "reviewed_by_user_id")
			if parseErr != nil {
				return dto.OperationResponse{}, parseErr
			}
			supportCase.ReviewedByUserID = &reviewerID
		}
	}

	now := s.now().UTC()
	supportCase.Status = status
	switch status {
	case entity.SupportStatusResolved:
		supportCase.ResolvedAt = &now
	case entity.SupportStatusClosed:
		supportCase.ClosedAt = &now
	}
	if status != entity.SupportStatusResolved {
		supportCase.ResolvedAt = nil
	}
	if status != entity.SupportStatusClosed {
		supportCase.ClosedAt = nil
	}

	supportCase.UpdatedAt = now
	if err := s.store.UpdateCase(ctx, supportCase); err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrSupportNotFound
		}
		if errors.Is(err, supportrepository.ErrConflict) {
			return dto.OperationResponse{}, ErrSupportAlreadyExists
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "status_updated"}, nil
}

func (s *SupportService) Resolve(ctx context.Context, request dto.ResolveSupportRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	supportID, err := parseID(request.SupportID, "support_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	reviewerID, err := parseID(request.ReviewedByUserID, "reviewed_by_user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	resolutionNote, err := sanitizeText(request.ResolutionNote, "resolution_note")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	supportCase, err := s.store.GetCaseByID(ctx, supportID)
	if err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrSupportNotFound
		}
		return dto.OperationResponse{}, err
	}
	if supportCase.Status == entity.SupportStatusClosed || supportCase.Status == entity.SupportStatusSpam {
		return dto.OperationResponse{}, ErrInvalidStateTransition
	}

	now := s.now().UTC()
	supportCase.Status = entity.SupportStatusResolved
	supportCase.ResolutionNote = resolutionNote
	supportCase.ReviewedByUserID = &reviewerID
	supportCase.ResolvedAt = &now
	supportCase.ClosedAt = nil
	supportCase.UpdatedAt = now

	if err := s.store.UpdateCase(ctx, supportCase); err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrSupportNotFound
		}
		if errors.Is(err, supportrepository.ErrConflict) {
			return dto.OperationResponse{}, ErrSupportAlreadyExists
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "resolved"}, nil
}

func (s *SupportService) RequestModerationHandoff(ctx context.Context, request dto.RequestModerationHandoffRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	supportID, err := parseID(request.SupportID, "support_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	supportCase, err := s.store.GetCaseByID(ctx, supportID)
	if err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrSupportNotFound
		}
		return dto.OperationResponse{}, err
	}
	if supportCase.Kind != entity.SupportKindReport || supportCase.TargetType == nil || supportCase.TargetID == nil {
		return dto.OperationResponse{}, ErrModerationHandoffNotAllowed
	}
	if supportCase.LinkedModerationCaseID != nil || supportCase.ModerationHandoffRequestedAt != nil {
		return dto.OperationResponse{}, ErrAlreadyHandedOff
	}

	now := s.now().UTC()
	supportCase.ModerationHandoffRequestedAt = &now
	supportCase.Status = entity.SupportStatusTriaged
	supportCase.UpdatedAt = now

	if err := s.store.UpdateCase(ctx, supportCase); err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrSupportNotFound
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "moderation_handoff_requested"}, nil
}

func isValidStatusTransition(from entity.SupportStatus, to entity.SupportStatus) bool {
	if from == to {
		return true
	}

	switch from {
	case entity.SupportStatusClosed:
		return false
	case entity.SupportStatusSpam:
		return to == entity.SupportStatusTriaged || to == entity.SupportStatusOpen || to == entity.SupportStatusRejected
	default:
		return true
	}
}

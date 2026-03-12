package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
	moderationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/repository"
)

func (s *ModerationService) ListQueue(ctx context.Context, request dto.ListQueueRequest) (dto.ListQueueResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListQueueResponse{}, err
	}

	status, err := parseOptionalCaseStatus(request.Status)
	if err != nil {
		return dto.ListQueueResponse{}, err
	}
	targetType, err := parseOptionalTargetType(request.TargetType)
	if err != nil {
		return dto.ListQueueResponse{}, err
	}

	assignedModerator := ""
	if request.AssignedModeratorUserID != "" {
		assignedModerator, err = parseID(request.AssignedModeratorUserID, "assigned_moderator_user_id")
		if err != nil {
			return dto.ListQueueResponse{}, err
		}
	}

	items, err := s.store.ListQueue(ctx, moderationrepository.QueueQuery{
		Status:                  string(status),
		TargetType:              string(targetType),
		AssignedModeratorUserID: assignedModerator,
		SortBy:                  parseSortBy(request.SortBy, "newest"),
		Limit:                   request.Limit,
		Offset:                  request.Offset,
	})
	if err != nil {
		return dto.ListQueueResponse{}, err
	}

	response := dto.ListQueueResponse{Cases: make([]dto.ModerationCaseResponse, 0, len(items)), Count: len(items)}
	for _, moderationCase := range items {
		response.Cases = append(response.Cases, toCaseResponse(moderationCase, false))
	}
	return response, nil
}

func (s *ModerationService) GetCaseDetail(ctx context.Context, request dto.GetCaseDetailRequest) (dto.ModerationCaseResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ModerationCaseResponse{}, err
	}

	caseID, err := parseID(request.CaseID, "case_id")
	if err != nil {
		return dto.ModerationCaseResponse{}, err
	}

	moderationCase, err := s.store.GetCaseByID(ctx, caseID)
	if err != nil {
		if errors.Is(err, moderationrepository.ErrNotFound) {
			return dto.ModerationCaseResponse{}, ErrCaseNotFound
		}
		return dto.ModerationCaseResponse{}, err
	}
	return toCaseResponse(moderationCase, true), nil
}

func toCaseResponse(moderationCase entity.Case, includeTimeline bool) dto.ModerationCaseResponse {
	response := dto.ModerationCaseResponse{
		CaseID:                   moderationCase.ID,
		Source:                   string(moderationCase.Source),
		SourceRefID:              moderationCase.SourceRefID,
		RequestID:                moderationCase.RequestID,
		CorrelationID:            moderationCase.CorrelationID,
		TargetType:               string(moderationCase.TargetType),
		TargetID:                 moderationCase.TargetID,
		Status:                   string(moderationCase.Status),
		AssignmentStatus:         string(moderationCase.AssignmentStatus),
		AssignedModeratorUserID:  moderationCase.AssignedModeratorUserID,
		EscalationStatus:         string(moderationCase.EscalationStatus),
		EscalationReason:         moderationCase.EscalationReason,
		ActionResult:             string(moderationCase.ActionResult),
		CreatedAt:                moderationCase.CreatedAt,
		UpdatedAt:                moderationCase.UpdatedAt,
	}

	if !includeTimeline {
		return response
	}

	if len(moderationCase.Notes) > 0 {
		response.Notes = make([]dto.NoteResponse, 0, len(moderationCase.Notes))
		for _, note := range moderationCase.Notes {
			response.Notes = append(response.Notes, dto.NoteResponse{
				NoteID:       note.ID,
				AuthorUserID: note.AuthorUserID,
				Body:         note.Body,
				InternalOnly: note.InternalOnly,
				CreatedAt:    note.CreatedAt,
			})
		}
	}
	if len(moderationCase.Actions) > 0 {
		response.Actions = make([]dto.ActionResponse, 0, len(moderationCase.Actions))
		for _, action := range moderationCase.Actions {
			response.Actions = append(response.Actions, dto.ActionResponse{
				ActionID:     action.ID,
				ActorUserID:  action.ActorUserID,
				ActionType:   string(action.ActionType),
				ReasonCode:   action.ReasonCode,
				Summary:      action.Summary,
				ActionResult: string(action.ActionResult),
				CreatedAt:    action.CreatedAt,
			})
		}
	}

	return response
}

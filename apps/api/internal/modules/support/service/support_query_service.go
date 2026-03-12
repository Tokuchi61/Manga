package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
	supportrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/repository"
)

func (s *SupportService) ListOwnSupport(ctx context.Context, request dto.ListOwnSupportRequest) (dto.ListOwnSupportResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListOwnSupportResponse{}, err
	}

	requesterID, err := parseID(request.RequesterUserID, "requester_user_id")
	if err != nil {
		return dto.ListOwnSupportResponse{}, err
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	items, err := s.store.ListCasesByRequester(ctx, supportrepository.ListQuery{
		RequesterUserID: requesterID,
		Status:          normalizeValue(request.Status),
		SortBy:          parseSortBy(request.SortBy, "newest"),
		Limit:           limit,
		Offset:          request.Offset,
	})
	if err != nil {
		return dto.ListOwnSupportResponse{}, err
	}

	result := make([]dto.SupportListItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toListItem(item))
	}

	return dto.ListOwnSupportResponse{Items: result, Count: len(result)}, nil
}

func (s *SupportService) GetSupportDetail(ctx context.Context, request dto.GetSupportDetailRequest) (dto.SupportDetailResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.SupportDetailResponse{}, err
	}

	supportID, err := parseID(request.SupportID, "support_id")
	if err != nil {
		return dto.SupportDetailResponse{}, err
	}
	requesterID, err := parseID(request.RequesterUserID, "requester_user_id")
	if err != nil {
		return dto.SupportDetailResponse{}, err
	}

	supportCase, err := s.store.GetCaseByID(ctx, supportID)
	if err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return dto.SupportDetailResponse{}, ErrSupportNotFound
		}
		return dto.SupportDetailResponse{}, err
	}
	if supportCase.RequesterUserID != requesterID {
		return dto.SupportDetailResponse{}, ErrForbiddenAction
	}

	return toDetailResponse(supportCase, request.IncludeInternal), nil
}

func (s *SupportService) ListReviewQueue(ctx context.Context, request dto.ListReviewQueueRequest) (dto.ListReviewQueueResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListReviewQueueResponse{}, err
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	items, err := s.store.ListReviewQueue(ctx, supportrepository.QueueQuery{
		Status:   normalizeValue(request.Status),
		Priority: normalizeValue(request.Priority),
		Limit:    limit,
		Offset:   request.Offset,
	})
	if err != nil {
		return dto.ListReviewQueueResponse{}, err
	}

	result := make([]dto.SupportListItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toListItem(item))
	}

	return dto.ListReviewQueueResponse{Items: result, Count: len(result)}, nil
}

func toListItem(supportCase entity.SupportCase) dto.SupportListItemResponse {
	return dto.SupportListItemResponse{
		SupportID:                    supportCase.ID,
		RequesterUserID:              supportCase.RequesterUserID,
		SupportKind:                  string(supportCase.Kind),
		Category:                     supportCase.Category,
		Priority:                     string(supportCase.Priority),
		ReasonCode:                   supportCase.ReasonCode,
		Status:                       string(supportCase.Status),
		TargetType:                   toStringPtr(supportCase.TargetType),
		TargetID:                     supportCase.TargetID,
		DuplicateOfSupportID:         supportCase.DuplicateOfSupportID,
		AssigneeUserID:               supportCase.AssigneeUserID,
		ReviewedByUserID:             supportCase.ReviewedByUserID,
		ResolvedAt:                   supportCase.ResolvedAt,
		ModerationHandoffRequestedAt: supportCase.ModerationHandoffRequestedAt,
		CreatedAt:                    supportCase.CreatedAt,
		UpdatedAt:                    supportCase.UpdatedAt,
	}
}

func toDetailResponse(supportCase entity.SupportCase, includeInternal bool) dto.SupportDetailResponse {
	replies := make([]dto.SupportReplyResponse, 0, len(supportCase.Replies))
	for _, reply := range supportCase.Replies {
		if !includeInternal && reply.Visibility == entity.ReplyVisibilityInternalOnly {
			continue
		}
		replies = append(replies, dto.SupportReplyResponse{
			ReplyID:      reply.ID,
			AuthorUserID: reply.AuthorUserID,
			Message:      reply.SanitizedBody,
			Visibility:   string(reply.Visibility),
			CreatedAt:    reply.CreatedAt,
			UpdatedAt:    reply.UpdatedAt,
		})
	}

	attachments := append([]string(nil), supportCase.Attachments...)

	return dto.SupportDetailResponse{
		SupportID:                    supportCase.ID,
		RequesterUserID:              supportCase.RequesterUserID,
		SupportKind:                  string(supportCase.Kind),
		Category:                     supportCase.Category,
		Priority:                     string(supportCase.Priority),
		ReasonCode:                   supportCase.ReasonCode,
		ReasonText:                   supportCase.ReasonText,
		Status:                       string(supportCase.Status),
		TargetType:                   toStringPtr(supportCase.TargetType),
		TargetID:                     supportCase.TargetID,
		DuplicateOfSupportID:         supportCase.DuplicateOfSupportID,
		SpamRiskScore:                supportCase.SpamRiskScore,
		Attachments:                  attachments,
		Replies:                      replies,
		ResolutionNote:               supportCase.ResolutionNote,
		AssigneeUserID:               supportCase.AssigneeUserID,
		ReviewedByUserID:             supportCase.ReviewedByUserID,
		ResolvedAt:                   supportCase.ResolvedAt,
		ClosedAt:                     supportCase.ClosedAt,
		ModerationHandoffRequestedAt: supportCase.ModerationHandoffRequestedAt,
		LinkedModerationCaseID:       supportCase.LinkedModerationCaseID,
		CreatedAt:                    supportCase.CreatedAt,
		UpdatedAt:                    supportCase.UpdatedAt,
	}
}

package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
	adminevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/events"
	adminrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/repository"
	"github.com/google/uuid"
)

func (s *AdminService) ReviewUser(ctx context.Context, actorUserID string, request dto.ReviewUserRequest) (dto.UserReviewResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.UserReviewResponse{}, err
	}
	if strings.TrimSpace(actorUserID) == "" {
		return dto.UserReviewResponse{}, ErrForbiddenAction
	}
	if err := requireDoubleConfirmation(request.RiskLevel, request.DoubleConfirmed, request.ConfirmationToken); err != nil {
		return dto.UserReviewResponse{}, err
	}

	dedupKey := buildActionDedupKey("user_review", request.RequestID)
	dedupAction, err := s.store.GetActionDedup(ctx, dedupKey)
	if err == nil {
		reviews, listErr := s.store.ListUserReviewRecords(ctx, request.TargetUserID, 0, 0)
		if listErr != nil {
			return dto.UserReviewResponse{}, listErr
		}
		reviewID := dedupAction.Metadata["review_id"]
		for _, record := range reviews {
			if record.ReviewID != reviewID {
				continue
			}
			return toUserReviewResponse("idempotent", adminevents.EventAdminUserReviewed, dedupAction.ActionID, record), nil
		}
		return dto.UserReviewResponse{}, ErrNotFound
	}
	if err != nil && !errors.Is(err, adminrepository.ErrNotFound) {
		return dto.UserReviewResponse{}, err
	}

	now := s.now().UTC()
	reviewID := uuid.NewString()
	action := entity.AdminAction{
		ActionID:                   uuid.NewString(),
		ActionType:                 entity.ActionTypeUserReviewed,
		ActorUserID:                actorUserID,
		TargetUserID:               request.TargetUserID,
		TargetModule:               "user",
		TargetType:                 "user_review",
		TargetID:                   request.TargetUserID,
		RequestID:                  request.RequestID,
		CorrelationID:              request.CorrelationID,
		Reason:                     request.Reason,
		Result:                     entity.ActionResultSuccess,
		RiskLevel:                  request.RiskLevel,
		RequiresDoubleConfirmation: requiresDoubleConfirmation(request.RiskLevel),
		DoubleConfirmed:            request.DoubleConfirmed,
		ConfirmationToken:          request.ConfirmationToken,
		Metadata: map[string]string{
			"review_id": reviewID,
			"decision":  request.Decision,
			"event":     adminevents.EventAdminUserReviewed,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	record := entity.UserReviewRecord{
		ReviewID:     reviewID,
		ActionID:     action.ActionID,
		TargetUserID: request.TargetUserID,
		Decision:     request.Decision,
		Reason:       request.Reason,
		RiskLevel:    request.RiskLevel,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.store.CreateAdminAction(ctx, action); err != nil {
		return dto.UserReviewResponse{}, err
	}
	if err := s.store.CreateUserReviewRecord(ctx, record); err != nil {
		return dto.UserReviewResponse{}, err
	}
	if err := s.store.PutActionDedup(ctx, dedupKey, action); err != nil {
		return dto.UserReviewResponse{}, err
	}

	return toUserReviewResponse("accepted", adminevents.EventAdminUserReviewed, action.ActionID, record), nil
}

func (s *AdminService) ListUserReviews(ctx context.Context, request dto.ListUserReviewsRequest) (dto.ListUserReviewsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListUserReviewsResponse{}, err
	}

	records, err := s.store.ListUserReviewRecords(ctx, request.TargetUserID, request.Limit, request.Offset)
	if err != nil {
		return dto.ListUserReviewsResponse{}, err
	}

	items := make([]dto.UserReviewResponse, 0, len(records))
	for _, record := range records {
		items = append(items, toUserReviewResponse("listed", adminevents.EventAdminUserReviewed, record.ActionID, record))
	}

	return dto.ListUserReviewsResponse{Items: items, Count: len(items)}, nil
}

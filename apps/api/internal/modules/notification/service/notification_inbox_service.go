package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
	notificationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/repository"
)

func (s *NotificationService) ListInbox(ctx context.Context, request dto.ListInboxRequest) (dto.ListInboxResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListInboxResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.ListInboxResponse{}, err
	}
	category, err := parseOptionalCategory(request.Category)
	if err != nil {
		return dto.ListInboxResponse{}, err
	}
	state, err := parseOptionalState(request.State)
	if err != nil {
		return dto.ListInboxResponse{}, err
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	items, err := s.store.ListInbox(ctx, notificationrepository.InboxQuery{
		UserID:     userID,
		Category:   string(category),
		State:      string(state),
		UnreadOnly: request.UnreadOnly,
		SortBy:     parseSortBy(request.SortBy, "newest"),
		Limit:      limit,
		Offset:     request.Offset,
	})
	if err != nil {
		return dto.ListInboxResponse{}, err
	}

	unreadItems, err := s.store.ListInbox(ctx, notificationrepository.InboxQuery{
		UserID:     userID,
		UnreadOnly: true,
		SortBy:     "newest",
		Limit:      100000,
		Offset:     0,
	})
	if err != nil {
		return dto.ListInboxResponse{}, err
	}

	result := make([]dto.NotificationListItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapListItem(item))
	}

	return dto.ListInboxResponse{
		Items:       result,
		Count:       len(result),
		UnreadCount: len(unreadItems),
	}, nil
}

func (s *NotificationService) GetNotificationDetail(ctx context.Context, request dto.GetNotificationDetailRequest) (dto.NotificationDetailResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.NotificationDetailResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.NotificationDetailResponse{}, err
	}
	notificationID, err := parseID(request.NotificationID, "notification_id")
	if err != nil {
		return dto.NotificationDetailResponse{}, err
	}

	notification, err := s.store.GetNotificationByID(ctx, notificationID)
	if err != nil {
		if errors.Is(err, notificationrepository.ErrNotFound) {
			return dto.NotificationDetailResponse{}, ErrNotFound
		}
		return dto.NotificationDetailResponse{}, err
	}
	if notification.UserID != userID {
		return dto.NotificationDetailResponse{}, ErrForbiddenAction
	}

	return mapDetail(notification), nil
}

func (s *NotificationService) MarkRead(ctx context.Context, request dto.MarkReadRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	notificationID, err := parseID(request.NotificationID, "notification_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	notification, err := s.store.GetNotificationByID(ctx, notificationID)
	if err != nil {
		if errors.Is(err, notificationrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrNotFound
		}
		return dto.OperationResponse{}, err
	}
	if notification.UserID != userID {
		return dto.OperationResponse{}, ErrForbiddenAction
	}

	if notification.State == entity.NotificationStateRead {
		return dto.OperationResponse{Status: "read"}, nil
	}

	now := s.now().UTC()
	notification.State = entity.NotificationStateRead
	notification.ReadAt = &now
	notification.UpdatedAt = now

	if err := s.store.UpdateNotification(ctx, notification); err != nil {
		if errors.Is(err, notificationrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrNotFound
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "read"}, nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
	notificationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/catalog"
	"github.com/google/uuid"
)

func (s *NotificationService) CreateFromSupportSignal(ctx context.Context, request dto.CreateFromSupportSignalRequest) (dto.CreateFromSupportSignalResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateFromSupportSignalResponse{}, err
	}
	if s.supportSignalProvider == nil {
		return dto.CreateFromSupportSignalResponse{}, ErrSupportSignalUnavailable
	}

	supportID, err := parseID(request.SupportID, "support_id")
	if err != nil {
		return dto.CreateFromSupportSignalResponse{}, err
	}
	channel, err := parseChannel(request.Channel, entity.DeliveryChannelInApp)
	if err != nil {
		return dto.CreateFromSupportSignalResponse{}, err
	}

	signal, err := s.supportSignalProvider.BuildNotificationSignal(
		ctx,
		supportID,
		strings.TrimSpace(request.Event),
		strings.TrimSpace(request.RequestID),
		strings.TrimSpace(request.CorrelationID),
	)
	if err != nil {
		return dto.CreateFromSupportSignalResponse{}, fmt.Errorf("%w: %v", ErrSupportSignalInvalid, err)
	}

	userID, err := parseID(signal.RequesterUserID, "requester_user_id")
	if err != nil {
		return dto.CreateFromSupportSignalResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.CreateFromSupportSignalResponse{}, err
	}
	if cfg.DeliveryPaused {
		return dto.CreateFromSupportSignalResponse{}, ErrDeliveryPaused
	}

	category := catalog.NotificationSupport
	if !cfg.CategoryEnabled[category] {
		return dto.CreateFromSupportSignalResponse{}, ErrCategoryDisabled
	}
	if !cfg.ChannelEnabled[channel] {
		return dto.CreateFromSupportSignalResponse{}, ErrChannelDisabled
	}

	preference, err := s.store.GetPreferenceByUserID(ctx, userID)
	if err != nil {
		return dto.CreateFromSupportSignalResponse{}, err
	}
	if hasMutedCategory(preference.MutedCategories, category) {
		return dto.CreateFromSupportSignalResponse{}, ErrCategoryMuted
	}
	if !preference.DigestEnabled && channel == entity.DeliveryChannelEmail {
		return dto.CreateFromSupportSignalResponse{}, ErrDigestDisabled
	}
	if channel == entity.DeliveryChannelInApp && !preference.InAppEnabled {
		return dto.CreateFromSupportSignalResponse{}, ErrChannelDisabled
	}
	if channel == entity.DeliveryChannelEmail && !preference.EmailEnabled {
		return dto.CreateFromSupportSignalResponse{}, ErrChannelDisabled
	}
	if channel == entity.DeliveryChannelPush && !preference.PushEnabled {
		return dto.CreateFromSupportSignalResponse{}, ErrChannelDisabled
	}

	eventName := strings.TrimSpace(signal.Event)
	if eventName == "" {
		eventName = "support.created"
	}
	dedupKey := strings.ToLower(fmt.Sprintf("support:%s:%s:%s:%s", eventName, supportID, userID, category))

	title := strings.TrimSpace(request.Title)
	if title == "" {
		title = "Support Guncellemesi"
	}
	title, err = sanitizeText(title, "title", 180)
	if err != nil {
		return dto.CreateFromSupportSignalResponse{}, err
	}

	body := strings.TrimSpace(request.Body)
	if body == "" {
		body = fmt.Sprintf("Support talebiniz (%s) durumu: %s", supportID, signal.Status)
	}
	body, err = sanitizeText(body, "body", 4000)
	if err != nil {
		return dto.CreateFromSupportSignalResponse{}, err
	}

	now := s.now().UTC()
	state := entity.NotificationStateCreated
	attemptCount := 0
	if channel == entity.DeliveryChannelInApp {
		state = entity.NotificationStateDelivered
		attemptCount = 1
	}

	notification := entity.Notification{
		ID:                   uuid.NewString(),
		UserID:               userID,
		Category:             category,
		Channel:              channel,
		TemplateKey:          "support_status_update",
		Title:                title,
		Body:                 body,
		State:                state,
		SourceEvent:          eventName,
		SourceRefID:          supportID,
		RequestID:            strings.TrimSpace(signal.RequestID),
		CorrelationID:        strings.TrimSpace(signal.CorrelationID),
		DeliveryAttemptCount: attemptCount,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	persisted, created, err := s.store.CreateNotification(ctx, notification, dedupKey)
	if err != nil {
		if errors.Is(err, notificationrepository.ErrConflict) {
			return dto.CreateFromSupportSignalResponse{}, ErrSupportSignalInvalid
		}
		return dto.CreateFromSupportSignalResponse{}, err
	}

	return dto.CreateFromSupportSignalResponse{
		NotificationID: persisted.ID,
		Created:        created,
		Category:       string(persisted.Category),
		State:          string(persisted.State),
	}, nil
}

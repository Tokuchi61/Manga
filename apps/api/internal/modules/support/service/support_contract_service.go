package service

import (
	"context"
	"errors"

	supportcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
	supportrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/repository"
)

// GetModerationHandoffReference exposes support->moderation linked-case contract.
func (s *SupportService) GetModerationHandoffReference(ctx context.Context, supportID string, requestID string, correlationID string) (supportcontract.ModerationHandoffReference, error) {
	parsedID, err := parseID(supportID, "support_id")
	if err != nil {
		return supportcontract.ModerationHandoffReference{}, err
	}

	supportCase, err := s.store.GetCaseByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return supportcontract.ModerationHandoffReference{}, ErrSupportNotFound
		}
		return supportcontract.ModerationHandoffReference{}, err
	}
	if supportCase.Kind != entity.SupportKindReport || supportCase.TargetType == nil || supportCase.TargetID == nil {
		return supportcontract.ModerationHandoffReference{}, ErrModerationHandoffNotAllowed
	}

	requestedAt := s.now().UTC()
	if supportCase.ModerationHandoffRequestedAt != nil {
		requestedAt = supportCase.ModerationHandoffRequestedAt.UTC()
	}
	if requestID == "" {
		requestID = supportCase.RequestID
	}

	return supportcontract.ModerationHandoffReference{
		SupportID:     supportCase.ID,
		SupportKind:   string(supportCase.Kind),
		TargetType:    string(*supportCase.TargetType),
		TargetID:      *supportCase.TargetID,
		ReasonCode:    supportCase.ReasonCode,
		RequestedAt:   requestedAt,
		RequestID:     requestID,
		CorrelationID: correlationID,
	}, nil
}

// LinkModerationCase stores linked moderation case id for a support report.
func (s *SupportService) LinkModerationCase(ctx context.Context, supportID string, moderationCaseID string) error {
	parsedSupportID, err := parseID(supportID, "support_id")
	if err != nil {
		return err
	}
	parsedModerationCaseID, err := parseID(moderationCaseID, "moderation_case_id")
	if err != nil {
		return err
	}

	supportCase, err := s.store.GetCaseByID(ctx, parsedSupportID)
	if err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return ErrSupportNotFound
		}
		return err
	}
	if supportCase.Kind != entity.SupportKindReport || supportCase.TargetType == nil || supportCase.TargetID == nil {
		return ErrModerationHandoffNotAllowed
	}
	if supportCase.LinkedModerationCaseID != nil {
		if *supportCase.LinkedModerationCaseID == parsedModerationCaseID {
			return nil
		}
		return ErrAlreadyHandedOff
	}

	now := s.now().UTC()
	supportCase.LinkedModerationCaseID = &parsedModerationCaseID
	if supportCase.ModerationHandoffRequestedAt == nil {
		supportCase.ModerationHandoffRequestedAt = &now
	}
	if supportCase.Status == entity.SupportStatusOpen {
		supportCase.Status = entity.SupportStatusTriaged
	}
	supportCase.UpdatedAt = now

	if err := s.store.UpdateCase(ctx, supportCase); err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return ErrSupportNotFound
		}
		if errors.Is(err, supportrepository.ErrConflict) {
			return ErrSupportAlreadyExists
		}
		return err
	}
	return nil
}

// BuildNotificationSignal creates stable support->notification signal payload.
func (s *SupportService) BuildNotificationSignal(ctx context.Context, supportID string, event string, requestID string, correlationID string) (supportcontract.NotificationSignal, error) {
	parsedID, err := parseID(supportID, "support_id")
	if err != nil {
		return supportcontract.NotificationSignal{}, err
	}

	supportCase, err := s.store.GetCaseByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, supportrepository.ErrNotFound) {
			return supportcontract.NotificationSignal{}, ErrSupportNotFound
		}
		return supportcontract.NotificationSignal{}, err
	}

	if event == "" {
		event = supportcontract.EventSupportCreated
	}
	if requestID == "" {
		requestID = supportCase.RequestID
	}

	return supportcontract.NotificationSignal{
		Event:           event,
		SupportID:       supportCase.ID,
		RequesterUserID: supportCase.RequesterUserID,
		Status:          string(supportCase.Status),
		OccurredAt:      s.now().UTC(),
		RequestID:       requestID,
		CorrelationID:   correlationID,
	}, nil
}

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

func (s *AdminService) ApplyOverride(ctx context.Context, actorUserID string, request dto.ApplyOverrideRequest) (dto.OverrideResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OverrideResponse{}, err
	}
	if strings.TrimSpace(actorUserID) == "" {
		return dto.OverrideResponse{}, ErrForbiddenAction
	}
	if err := requireDoubleConfirmation(request.RiskLevel, request.DoubleConfirmed, request.ConfirmationToken); err != nil {
		return dto.OverrideResponse{}, err
	}

	dedupKey := buildActionDedupKey("override", request.RequestID)
	dedupAction, err := s.store.GetActionDedup(ctx, dedupKey)
	if err == nil {
		overrides, listErr := s.store.ListOverrideRecords(ctx, "", 0, 0)
		if listErr != nil {
			return dto.OverrideResponse{}, listErr
		}
		overrideID := dedupAction.Metadata["override_id"]
		for _, record := range overrides {
			if record.OverrideID != overrideID {
				continue
			}
			return toOverrideResponse("idempotent", adminevents.EventAdminOverrideApplied, dedupAction.ActionID, record), nil
		}
		return dto.OverrideResponse{}, ErrNotFound
	}
	if err != nil && !errors.Is(err, adminrepository.ErrNotFound) {
		return dto.OverrideResponse{}, err
	}

	now := s.now().UTC()
	overrideID := uuid.NewString()
	action := entity.AdminAction{
		ActionID:                   uuid.NewString(),
		ActionType:                 entity.ActionTypeOverrideApplied,
		ActorUserID:                actorUserID,
		TargetModule:               request.TargetModule,
		TargetType:                 request.TargetType,
		TargetID:                   request.TargetID,
		RequestID:                  request.RequestID,
		CorrelationID:              request.CorrelationID,
		Reason:                     request.Reason,
		Result:                     entity.ActionResultSuccess,
		RiskLevel:                  request.RiskLevel,
		RequiresDoubleConfirmation: requiresDoubleConfirmation(request.RiskLevel),
		DoubleConfirmed:            request.DoubleConfirmed,
		ConfirmationToken:          request.ConfirmationToken,
		Metadata: map[string]string{
			"override_id": overrideID,
			"decision":    request.Decision,
			"event":       adminevents.EventAdminOverrideApplied,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	record := entity.OverrideRecord{
		OverrideID:   overrideID,
		ActionID:     action.ActionID,
		TargetModule: request.TargetModule,
		TargetType:   request.TargetType,
		TargetID:     request.TargetID,
		Decision:     request.Decision,
		Reason:       request.Reason,
		RiskLevel:    request.RiskLevel,
		Active:       true,
		ExpiresAt:    request.ExpiresAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.store.CreateAdminAction(ctx, action); err != nil {
		return dto.OverrideResponse{}, err
	}
	if err := s.store.CreateOverrideRecord(ctx, record); err != nil {
		return dto.OverrideResponse{}, err
	}
	if err := s.store.PutActionDedup(ctx, dedupKey, action); err != nil {
		return dto.OverrideResponse{}, err
	}

	return toOverrideResponse("accepted", adminevents.EventAdminOverrideApplied, action.ActionID, record), nil
}

func (s *AdminService) ListOverrides(ctx context.Context, request dto.ListOverridesRequest) (dto.ListOverridesResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListOverridesResponse{}, err
	}

	records, err := s.store.ListOverrideRecords(ctx, request.TargetModule, request.Limit, request.Offset)
	if err != nil {
		return dto.ListOverridesResponse{}, err
	}

	items := make([]dto.OverrideResponse, 0, len(records))
	for _, record := range records {
		items = append(items, toOverrideResponse("listed", adminevents.EventAdminOverrideApplied, record.ActionID, record))
	}

	return dto.ListOverridesResponse{Items: items, Count: len(items)}, nil
}

package service

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
	adminevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/events"
	adminrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/repository"
	"github.com/google/uuid"
)

func (s *AdminService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	return dto.RuntimeConfigResponse{
		MaintenanceEnabled: cfg.MaintenanceEnabled,
		UpdatedAt:          cfg.UpdatedAt,
	}, nil
}

func (s *AdminService) UpdateMaintenanceState(ctx context.Context, actorUserID string, request dto.UpdateMaintenanceStateRequest) (dto.SettingChangeResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.SettingChangeResponse{}, err
	}
	if strings.TrimSpace(actorUserID) == "" {
		return dto.SettingChangeResponse{}, ErrForbiddenAction
	}
	if err := requireDoubleConfirmation(request.RiskLevel, request.DoubleConfirmed, request.ConfirmationToken); err != nil {
		return dto.SettingChangeResponse{}, err
	}

	dedupKey := buildActionDedupKey("maintenance", request.RequestID)
	dedupAction, err := s.store.GetActionDedup(ctx, dedupKey)
	if err == nil {
		cfg, cfgErr := s.store.GetRuntimeConfig(ctx)
		if cfgErr != nil {
			return dto.SettingChangeResponse{}, cfgErr
		}
		return dto.SettingChangeResponse{
			Status:    "idempotent",
			ActionID:  dedupAction.ActionID,
			Key:       "site.maintenance.enabled",
			Enabled:   cfg.MaintenanceEnabled,
			Event:     adminevents.EventAdminSettingChanged,
			UpdatedAt: cfg.UpdatedAt,
		}, nil
	}
	if err != nil && !errors.Is(err, adminrepository.ErrNotFound) {
		return dto.SettingChangeResponse{}, err
	}

	now := s.now().UTC()
	action := entity.AdminAction{
		ActionID:                   uuid.NewString(),
		ActionType:                 entity.ActionTypeSettingChanged,
		ActorUserID:                actorUserID,
		TargetModule:               "site",
		TargetType:                 "setting",
		TargetID:                   "site.maintenance.enabled",
		RequestID:                  request.RequestID,
		CorrelationID:              request.CorrelationID,
		Reason:                     request.Reason,
		Result:                     entity.ActionResultSuccess,
		RiskLevel:                  request.RiskLevel,
		RequiresDoubleConfirmation: requiresDoubleConfirmation(request.RiskLevel),
		DoubleConfirmed:            request.DoubleConfirmed,
		ConfirmationToken:          request.ConfirmationToken,
		Metadata: map[string]string{
			"key":     "site.maintenance.enabled",
			"enabled": strconv.FormatBool(request.Enabled),
			"event":   adminevents.EventAdminSettingChanged,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.CreateAdminAction(ctx, action); err != nil {
		return dto.SettingChangeResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.SettingChangeResponse{}, err
	}
	cfg.MaintenanceEnabled = request.Enabled
	cfg.UpdatedAt = now
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.SettingChangeResponse{}, err
	}
	if err := s.store.PutActionDedup(ctx, dedupKey, action); err != nil {
		return dto.SettingChangeResponse{}, err
	}

	return dto.SettingChangeResponse{
		Status:    "accepted",
		ActionID:  action.ActionID,
		Key:       "site.maintenance.enabled",
		Enabled:   cfg.MaintenanceEnabled,
		Event:     adminevents.EventAdminSettingChanged,
		UpdatedAt: cfg.UpdatedAt,
	}, nil
}

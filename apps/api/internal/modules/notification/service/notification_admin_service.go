package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/catalog"
)

func (s *NotificationService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *NotificationService) UpdateCategoryState(ctx context.Context, request dto.UpdateCategoryStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	category, err := parseOptionalCategory(request.Category)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	if category == "" {
		return dto.RuntimeConfigResponse{}, ErrValidation
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	if cfg.CategoryEnabled == nil {
		cfg.CategoryEnabled = make(map[catalog.NotificationCategory]bool)
	}
	cfg.CategoryEnabled[category] = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *NotificationService) UpdateChannelState(ctx context.Context, request dto.UpdateChannelStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	channel, err := parseChannel(request.Channel, "")
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	if cfg.ChannelEnabled == nil {
		cfg.ChannelEnabled = make(map[entity.DeliveryChannel]bool)
	}
	cfg.ChannelEnabled[channel] = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *NotificationService) UpdateDigestState(ctx context.Context, request dto.UpdateDigestStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.DigestEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *NotificationService) UpdateDeliveryPause(ctx context.Context, request dto.UpdateDeliveryPauseRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.DeliveryPaused = request.Paused
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

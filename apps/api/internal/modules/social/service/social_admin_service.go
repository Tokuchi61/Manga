package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/dto"
)

func (s *SocialService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *SocialService) UpdateFriendshipState(ctx context.Context, request dto.UpdateFriendshipStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.FriendshipEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *SocialService) UpdateFollowState(ctx context.Context, request dto.UpdateFollowStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.FollowEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *SocialService) UpdateWallState(ctx context.Context, request dto.UpdateWallStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.WallEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *SocialService) UpdateMessagingState(ctx context.Context, request dto.UpdateMessagingStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.MessagingEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

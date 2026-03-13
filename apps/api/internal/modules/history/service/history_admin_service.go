package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/dto"
)

func (s *HistoryService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *HistoryService) UpdateContinueReadingState(ctx context.Context, request dto.UpdateContinueReadingStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.ContinueReadingEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *HistoryService) UpdateLibraryState(ctx context.Context, request dto.UpdateLibraryStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.LibraryEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *HistoryService) UpdateTimelineState(ctx context.Context, request dto.UpdateTimelineStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.TimelineEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *HistoryService) UpdateBookmarkWriteState(ctx context.Context, request dto.UpdateBookmarkWriteStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.BookmarkWriteEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

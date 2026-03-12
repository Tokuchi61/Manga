package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
	missionrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/repository"
)

func (s *MissionService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return dto.RuntimeConfigResponse{
		ReadEnabled:           cfg.ReadEnabled,
		ClaimEnabled:          cfg.ClaimEnabled,
		ProgressIngestEnabled: cfg.ProgressIngestEnabled,
		DailyResetHourUTC:     cfg.DailyResetHourUTC,
		UpdatedAt:             cfg.UpdatedAt,
	}, nil
}

func (s *MissionService) UpdateReadState(ctx context.Context, request dto.UpdateReadStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.ReadEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *MissionService) UpdateClaimState(ctx context.Context, request dto.UpdateClaimStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.ClaimEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *MissionService) UpdateProgressIngestState(ctx context.Context, request dto.UpdateProgressIngestStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.ProgressIngestEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *MissionService) UpdateDailyResetHour(ctx context.Context, request dto.UpdateDailyResetHourRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.DailyResetHourUTC = request.Hour
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *MissionService) UpsertMissionDefinition(ctx context.Context, request dto.UpsertMissionDefinitionRequest) (dto.MissionDefinitionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.MissionDefinitionResponse{}, err
	}
	if request.StartsAt != nil && request.EndsAt != nil && request.StartsAt.After(*request.EndsAt) {
		return dto.MissionDefinitionResponse{}, fmt.Errorf("%w: starts_at_after_ends_at", ErrValidation)
	}

	now := s.now().UTC()
	definition := entity.MissionDefinition{
		MissionID:      request.MissionID,
		Category:       request.Category,
		Title:          request.Title,
		ObjectiveType:  request.ObjectiveType,
		TargetCount:    request.TargetCount,
		RewardItemID:   request.RewardItemID,
		RewardQuantity: request.RewardQuantity,
		Active:         request.Active,
		StartsAt:       request.StartsAt,
		EndsAt:         request.EndsAt,
		UpdatedAt:      now,
	}

	existing, err := s.store.GetMissionDefinition(ctx, request.MissionID)
	if err == nil {
		definition.CreatedAt = existing.CreatedAt
	} else if err != nil && !errors.Is(err, missionrepository.ErrNotFound) {
		return dto.MissionDefinitionResponse{}, err
	}
	if definition.CreatedAt.IsZero() {
		definition.CreatedAt = now
	}

	if err := s.store.UpsertMissionDefinition(ctx, definition); err != nil {
		return dto.MissionDefinitionResponse{}, err
	}
	return toDefinitionResponse(definition), nil
}

func (s *MissionService) ListMissionDefinitions(ctx context.Context, request dto.ListMissionDefinitionsRequest) (dto.ListMissionDefinitionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListMissionDefinitionsResponse{}, err
	}

	definitions, err := s.store.ListMissionDefinitions(ctx, request.Category, request.ActiveOnly, request.Limit, request.Offset)
	if err != nil {
		return dto.ListMissionDefinitionsResponse{}, err
	}

	items := make([]dto.MissionDefinitionResponse, 0, len(definitions))
	for _, definition := range definitions {
		items = append(items, toDefinitionResponse(definition))
	}
	return dto.ListMissionDefinitionsResponse{Items: items, Count: len(items)}, nil
}

func (s *MissionService) ResetMissionProgress(ctx context.Context, request dto.ResetMissionProgressRequest) (dto.ResetMissionProgressResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ResetMissionProgressResponse{}, err
	}

	deletedCount := 0
	var err error
	if request.MissionID != "" {
		deletedCount, err = s.store.DeleteMissionProgressByUserMission(ctx, request.TargetUserID, request.MissionID)
	} else {
		deletedCount, err = s.store.DeleteMissionProgressByUser(ctx, request.TargetUserID)
	}
	if err != nil {
		return dto.ResetMissionProgressResponse{}, err
	}

	return dto.ResetMissionProgressResponse{Status: "reset", DeletedCount: deletedCount}, nil
}

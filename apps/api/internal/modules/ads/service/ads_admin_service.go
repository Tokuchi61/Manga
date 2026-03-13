package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
	adsrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/repository"
)

func (s *AdsService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return dto.RuntimeConfigResponse{
		SurfaceEnabled:     cfg.SurfaceEnabled,
		PlacementEnabled:   cfg.PlacementEnabled,
		CampaignEnabled:    cfg.CampaignEnabled,
		ClickIntakeEnabled: cfg.ClickIntakeEnabled,
		UpdatedAt:          cfg.UpdatedAt,
	}, nil
}

func (s *AdsService) UpdateSurfaceState(ctx context.Context, request dto.UpdateSurfaceStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.SurfaceEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *AdsService) UpdatePlacementState(ctx context.Context, request dto.UpdatePlacementStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.PlacementEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *AdsService) UpdateCampaignState(ctx context.Context, request dto.UpdateCampaignStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.CampaignEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *AdsService) UpdateClickIntakeState(ctx context.Context, request dto.UpdateClickIntakeStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.ClickIntakeEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *AdsService) UpsertPlacementDefinition(ctx context.Context, request dto.UpsertPlacementDefinitionRequest) (dto.PlacementDefinitionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.PlacementDefinitionResponse{}, err
	}
	if normalizeValue(request.TargetType) == "none" {
		request.TargetID = ""
	}

	now := s.now().UTC()
	placement := entity.PlacementDefinition{
		PlacementID:  request.PlacementID,
		Surface:      request.Surface,
		TargetType:   request.TargetType,
		TargetID:     request.TargetID,
		Visible:      request.Visible,
		Priority:     request.Priority,
		FrequencyCap: request.FrequencyCap,
		UpdatedAt:    now,
	}

	existing, err := s.store.GetPlacementDefinition(ctx, request.PlacementID)
	if err == nil {
		placement.CreatedAt = existing.CreatedAt
	} else if err != nil && !errors.Is(err, adsrepository.ErrNotFound) {
		return dto.PlacementDefinitionResponse{}, err
	}
	if placement.CreatedAt.IsZero() {
		placement.CreatedAt = now
	}

	if err := s.store.UpsertPlacementDefinition(ctx, placement); err != nil {
		return dto.PlacementDefinitionResponse{}, err
	}
	return toPlacementDefinitionResponse(placement), nil
}

func (s *AdsService) ListPlacementDefinitions(ctx context.Context, request dto.ListPlacementDefinitionsRequest) (dto.ListPlacementDefinitionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListPlacementDefinitionsResponse{}, err
	}

	placements, err := s.store.ListPlacementDefinitions(ctx, request.Surface, request.Visible, request.Limit, request.Offset)
	if err != nil {
		return dto.ListPlacementDefinitionsResponse{}, err
	}

	items := make([]dto.PlacementDefinitionResponse, 0, len(placements))
	for _, placement := range placements {
		items = append(items, toPlacementDefinitionResponse(placement))
	}
	return dto.ListPlacementDefinitionsResponse{Items: items, Count: len(items)}, nil
}

func (s *AdsService) UpsertCampaignDefinition(ctx context.Context, request dto.UpsertCampaignDefinitionRequest) (dto.CampaignDefinitionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CampaignDefinitionResponse{}, err
	}
	if request.StartsAt != nil && request.EndsAt != nil && request.StartsAt.After(*request.EndsAt) {
		return dto.CampaignDefinitionResponse{}, fmt.Errorf("%w: starts_at_after_ends_at", ErrValidation)
	}

	if _, err := s.store.GetPlacementDefinition(ctx, request.PlacementID); err != nil {
		if errors.Is(err, adsrepository.ErrNotFound) {
			return dto.CampaignDefinitionResponse{}, ErrNotFound
		}
		return dto.CampaignDefinitionResponse{}, err
	}

	now := s.now().UTC()
	campaign := entity.CampaignDefinition{
		CampaignID:  request.CampaignID,
		PlacementID: request.PlacementID,
		Name:        request.Name,
		State:       request.State,
		CreativeURL: request.CreativeURL,
		ClickURL:    request.ClickURL,
		Weight:      request.Weight,
		StartsAt:    request.StartsAt,
		EndsAt:      request.EndsAt,
		UpdatedAt:   now,
	}

	existing, err := s.store.GetCampaignDefinition(ctx, request.CampaignID)
	if err == nil {
		campaign.CreatedAt = existing.CreatedAt
	} else if err != nil && !errors.Is(err, adsrepository.ErrNotFound) {
		return dto.CampaignDefinitionResponse{}, err
	}
	if campaign.CreatedAt.IsZero() {
		campaign.CreatedAt = now
	}

	if err := s.store.UpsertCampaignDefinition(ctx, campaign); err != nil {
		return dto.CampaignDefinitionResponse{}, err
	}
	return toCampaignDefinitionResponse(campaign), nil
}

func (s *AdsService) ListCampaignDefinitions(ctx context.Context, request dto.ListCampaignDefinitionsRequest) (dto.ListCampaignDefinitionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListCampaignDefinitionsResponse{}, err
	}

	campaigns, err := s.store.ListCampaignDefinitions(ctx, request.PlacementID, request.State, request.Limit, request.Offset)
	if err != nil {
		return dto.ListCampaignDefinitionsResponse{}, err
	}

	items := make([]dto.CampaignDefinitionResponse, 0, len(campaigns))
	for _, campaign := range campaigns {
		items = append(items, toCampaignDefinitionResponse(campaign))
	}
	return dto.ListCampaignDefinitionsResponse{Items: items, Count: len(items)}, nil
}

func (s *AdsService) ListCampaignAggregate(ctx context.Context, limit int, offset int) (dto.ListCampaignAggregateResponse, error) {
	aggregates, err := s.store.ListCampaignAggregate(ctx, limit, offset)
	if err != nil {
		return dto.ListCampaignAggregateResponse{}, err
	}

	items := make([]dto.CampaignAggregateResponse, 0, len(aggregates))
	for _, aggregate := range aggregates {
		ctr := 0.0
		if aggregate.ImpressionCount > 0 {
			ctr = (float64(aggregate.ClickCount) / float64(aggregate.ImpressionCount)) * 100
		}
		items = append(items, dto.CampaignAggregateResponse{
			CampaignID:      aggregate.CampaignID,
			ImpressionCount: aggregate.ImpressionCount,
			ClickCount:      aggregate.ClickCount,
			CTRPercent:      ctr,
			UpdatedAt:       aggregate.UpdatedAt,
		})
	}

	return dto.ListCampaignAggregateResponse{Items: items, Count: len(items)}, nil
}

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
	rprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/repository"
)

func (s *RoyalPassService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return dto.RuntimeConfigResponse{
		SeasonEnabled:  cfg.SeasonEnabled,
		ClaimEnabled:   cfg.ClaimEnabled,
		PremiumEnabled: cfg.PremiumEnabled,
		UpdatedAt:      cfg.UpdatedAt,
	}, nil
}

func (s *RoyalPassService) UpdateSeasonState(ctx context.Context, request dto.UpdateSeasonStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.SeasonEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *RoyalPassService) UpdateClaimState(ctx context.Context, request dto.UpdateClaimStateRequest) (dto.RuntimeConfigResponse, error) {
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

func (s *RoyalPassService) UpdatePremiumState(ctx context.Context, request dto.UpdatePremiumStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.PremiumEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *RoyalPassService) UpsertSeasonDefinition(ctx context.Context, request dto.UpsertSeasonDefinitionRequest) (dto.SeasonDefinitionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.SeasonDefinitionResponse{}, err
	}
	if request.StartsAt != nil && request.EndsAt != nil && request.StartsAt.After(*request.EndsAt) {
		return dto.SeasonDefinitionResponse{}, fmt.Errorf("%w: starts_at_after_ends_at", ErrValidation)
	}

	now := s.now().UTC()
	season := entity.SeasonDefinition{
		SeasonID:  request.SeasonID,
		Title:     request.Title,
		State:     request.State,
		StartsAt:  request.StartsAt,
		EndsAt:    request.EndsAt,
		UpdatedAt: now,
	}

	existing, err := s.store.GetSeasonDefinition(ctx, request.SeasonID)
	if err == nil {
		season.CreatedAt = existing.CreatedAt
	} else if err != nil && !errors.Is(err, rprepository.ErrNotFound) {
		return dto.SeasonDefinitionResponse{}, err
	}
	if season.CreatedAt.IsZero() {
		season.CreatedAt = now
	}

	if err := s.store.UpsertSeasonDefinition(ctx, season); err != nil {
		return dto.SeasonDefinitionResponse{}, err
	}
	return toSeasonDefinitionResponse(season), nil
}

func (s *RoyalPassService) ListSeasonDefinitions(ctx context.Context, request dto.ListSeasonDefinitionsRequest) (dto.ListSeasonDefinitionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListSeasonDefinitionsResponse{}, err
	}

	seasons, err := s.store.ListSeasonDefinitions(ctx, request.State, request.Limit, request.Offset)
	if err != nil {
		return dto.ListSeasonDefinitionsResponse{}, err
	}

	items := make([]dto.SeasonDefinitionResponse, 0, len(seasons))
	for _, season := range seasons {
		items = append(items, toSeasonDefinitionResponse(season))
	}
	return dto.ListSeasonDefinitionsResponse{Items: items, Count: len(items)}, nil
}

func (s *RoyalPassService) UpsertTierDefinition(ctx context.Context, request dto.UpsertTierDefinitionRequest) (dto.TierDefinitionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.TierDefinitionResponse{}, err
	}

	if _, err := s.store.GetSeasonDefinition(ctx, request.SeasonID); err != nil {
		if errors.Is(err, rprepository.ErrNotFound) {
			return dto.TierDefinitionResponse{}, ErrNotFound
		}
		return dto.TierDefinitionResponse{}, err
	}

	now := s.now().UTC()
	tier := entity.TierDefinition{
		SeasonID:       request.SeasonID,
		TierNumber:     request.TierNumber,
		Track:          request.Track,
		RequiredPoints: request.RequiredPoints,
		RewardItemID:   request.RewardItemID,
		RewardQuantity: request.RewardQuantity,
		Active:         request.Active,
		UpdatedAt:      now,
	}

	existing, err := s.store.GetTierDefinition(ctx, request.SeasonID, request.TierNumber, request.Track)
	if err == nil {
		tier.CreatedAt = existing.CreatedAt
	} else if err != nil && !errors.Is(err, rprepository.ErrNotFound) {
		return dto.TierDefinitionResponse{}, err
	}
	if tier.CreatedAt.IsZero() {
		tier.CreatedAt = now
	}

	if err := s.store.UpsertTierDefinition(ctx, tier); err != nil {
		return dto.TierDefinitionResponse{}, err
	}
	return toTierDefinitionResponse(tier), nil
}

func (s *RoyalPassService) ListTierDefinitions(ctx context.Context, request dto.ListTierDefinitionsRequest) (dto.ListTierDefinitionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListTierDefinitionsResponse{}, err
	}

	tiers, err := s.store.ListTierDefinitions(ctx, request.SeasonID, request.Track, request.ActiveOnly, request.Limit, request.Offset)
	if err != nil {
		return dto.ListTierDefinitionsResponse{}, err
	}

	items := make([]dto.TierDefinitionResponse, 0, len(tiers))
	for _, tier := range tiers {
		items = append(items, toTierDefinitionResponse(tier))
	}
	return dto.ListTierDefinitionsResponse{Items: items, Count: len(items)}, nil
}

func (s *RoyalPassService) ResetRoyalPassProgress(ctx context.Context, request dto.ResetRoyalPassProgressRequest) (dto.ResetRoyalPassProgressResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ResetRoyalPassProgressResponse{}, err
	}

	deletedCount := 0
	var err error
	if request.SeasonID != "" {
		deletedCount, err = s.store.DeleteUserProgressByUserSeason(ctx, request.TargetUserID, request.SeasonID)
	} else {
		deletedCount, err = s.store.DeleteUserProgressByUser(ctx, request.TargetUserID)
	}
	if err != nil {
		return dto.ResetRoyalPassProgressResponse{}, err
	}
	return dto.ResetRoyalPassProgressResponse{Status: "reset", DeletedCount: deletedCount}, nil
}

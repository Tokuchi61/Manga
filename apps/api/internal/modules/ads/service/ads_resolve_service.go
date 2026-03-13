package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
)

func (s *AdsService) ResolvePlacements(ctx context.Context, request dto.ResolvePlacementsRequest) (dto.ResolvePlacementsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ResolvePlacementsResponse{}, err
	}
	if request.NoAds {
		return dto.ResolvePlacementsResponse{Items: []dto.PlacementResolveResponse{}, Count: 0}, nil
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ResolvePlacementsResponse{}, err
	}
	if err := s.requireSurfaceEnabled(cfg.SurfaceEnabled); err != nil {
		return dto.ResolvePlacementsResponse{}, err
	}
	if err := s.requirePlacementEnabled(cfg.PlacementEnabled); err != nil {
		return dto.ResolvePlacementsResponse{}, err
	}
	if err := s.requireCampaignEnabled(cfg.CampaignEnabled); err != nil {
		return dto.ResolvePlacementsResponse{}, err
	}

	now := s.now().UTC()
	placements, err := s.store.ListPlacementDefinitions(ctx, request.Surface, true, 0, 0)
	if err != nil {
		return dto.ResolvePlacementsResponse{}, err
	}

	resolved := make([]dto.PlacementResolveResponse, 0, len(placements))
	for _, placement := range placements {
		if !matchesPlacementTarget(placement, request.TargetType, request.TargetID) {
			continue
		}
		if request.SessionID != "" && placement.FrequencyCap > 0 {
			seenCount, countErr := s.store.CountImpressionsBySessionPlacement(ctx, request.SessionID, placement.PlacementID)
			if countErr != nil {
				return dto.ResolvePlacementsResponse{}, countErr
			}
			if seenCount >= placement.FrequencyCap {
				continue
			}
		}

		campaigns, campaignErr := s.store.ListCampaignDefinitions(ctx, placement.PlacementID, entity.CampaignStateActive, 0, 0)
		if campaignErr != nil {
			return dto.ResolvePlacementsResponse{}, campaignErr
		}
		if len(campaigns) == 0 {
			continue
		}

		previews := make([]dto.CampaignPreviewResponse, 0, len(campaigns))
		for _, campaign := range campaigns {
			if !isCampaignServeable(campaign, now) {
				continue
			}
			previews = append(previews, toCampaignPreview(campaign))
		}
		if len(previews) == 0 {
			continue
		}

		resolved = append(resolved, dto.PlacementResolveResponse{
			PlacementID:  placement.PlacementID,
			Surface:      placement.Surface,
			Priority:     placement.Priority,
			FrequencyCap: placement.FrequencyCap,
			Campaigns:    previews,
		})
	}

	sortByPriority(resolved)
	resolved = applyOffsetLimit(resolved, 0, request.Limit)
	return dto.ResolvePlacementsResponse{Items: resolved, Count: len(resolved)}, nil
}

package service

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
	adsrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/repository"
)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func buildImpressionDedupKey(sessionID string, requestID string) string {
	return normalizeValue(sessionID) + ":" + normalizeValue(requestID)
}

func buildClickDedupKey(sessionID string, requestID string) string {
	return normalizeValue(sessionID) + ":" + normalizeValue(requestID)
}

func matchesPlacementTarget(placement entity.PlacementDefinition, targetType string, targetID string) bool {
	switch normalizeValue(placement.TargetType) {
	case "", "none":
		return true
	}
	if normalizeValue(placement.TargetType) != normalizeValue(targetType) {
		return false
	}
	if normalizeValue(placement.TargetID) == "" {
		return true
	}
	return normalizeValue(placement.TargetID) == normalizeValue(targetID)
}

func isCampaignServeable(campaign entity.CampaignDefinition, now time.Time) bool {
	if normalizeValue(campaign.State) != entity.CampaignStateActive {
		return false
	}
	utcNow := now.UTC()
	if campaign.StartsAt != nil && utcNow.Before(campaign.StartsAt.UTC()) {
		return false
	}
	if campaign.EndsAt != nil && utcNow.After(campaign.EndsAt.UTC()) {
		return false
	}
	return true
}

func toPlacementDefinitionResponse(placement entity.PlacementDefinition) dto.PlacementDefinitionResponse {
	return dto.PlacementDefinitionResponse{
		PlacementID:  placement.PlacementID,
		Surface:      placement.Surface,
		TargetType:   placement.TargetType,
		TargetID:     placement.TargetID,
		Visible:      placement.Visible,
		Priority:     placement.Priority,
		FrequencyCap: placement.FrequencyCap,
		UpdatedAt:    placement.UpdatedAt,
	}
}

func toCampaignDefinitionResponse(campaign entity.CampaignDefinition) dto.CampaignDefinitionResponse {
	return dto.CampaignDefinitionResponse{
		CampaignID:  campaign.CampaignID,
		PlacementID: campaign.PlacementID,
		Name:        campaign.Name,
		State:       campaign.State,
		CreativeURL: campaign.CreativeURL,
		ClickURL:    campaign.ClickURL,
		Weight:      campaign.Weight,
		StartsAt:    campaign.StartsAt,
		EndsAt:      campaign.EndsAt,
		UpdatedAt:   campaign.UpdatedAt,
	}
}

func toCampaignPreview(campaign entity.CampaignDefinition) dto.CampaignPreviewResponse {
	return dto.CampaignPreviewResponse{
		CampaignID:  campaign.CampaignID,
		Name:        campaign.Name,
		State:       campaign.State,
		CreativeURL: campaign.CreativeURL,
		ClickURL:    campaign.ClickURL,
		Weight:      campaign.Weight,
		StartsAt:    campaign.StartsAt,
		EndsAt:      campaign.EndsAt,
	}
}

func applyOffsetLimit[T any](items []T, offset int, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []T{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]T(nil), items[offset:end]...)
}

func sortByPriority(items []dto.PlacementResolveResponse) {
	sort.Slice(items, func(i int, j int) bool {
		return items[i].Priority > items[j].Priority
	})
}

func (s *AdsService) ensurePlacementAndCampaignActive(ctx context.Context, placementID string, campaignID string, now time.Time) (entity.PlacementDefinition, entity.CampaignDefinition, error) {
	placement, err := s.store.GetPlacementDefinition(ctx, placementID)
	if err != nil {
		if errors.Is(err, adsrepository.ErrNotFound) {
			return entity.PlacementDefinition{}, entity.CampaignDefinition{}, ErrNotFound
		}
		return entity.PlacementDefinition{}, entity.CampaignDefinition{}, err
	}
	if !placement.Visible {
		return entity.PlacementDefinition{}, entity.CampaignDefinition{}, ErrNotFound
	}

	campaign, err := s.store.GetCampaignDefinition(ctx, campaignID)
	if err != nil {
		if errors.Is(err, adsrepository.ErrNotFound) {
			return entity.PlacementDefinition{}, entity.CampaignDefinition{}, ErrNotFound
		}
		return entity.PlacementDefinition{}, entity.CampaignDefinition{}, err
	}
	if normalizeValue(campaign.PlacementID) != normalizeValue(placement.PlacementID) {
		return entity.PlacementDefinition{}, entity.CampaignDefinition{}, ErrValidation
	}
	if !isCampaignServeable(campaign, now) {
		return entity.PlacementDefinition{}, entity.CampaignDefinition{}, ErrNotFound
	}

	return placement, campaign, nil
}

func (s *AdsService) incrementImpressionAggregate(ctx context.Context, campaignID string, now time.Time) error {
	aggregate, err := s.store.GetCampaignAggregate(ctx, campaignID)
	if err != nil {
		if !errors.Is(err, adsrepository.ErrNotFound) {
			return err
		}
		aggregate = entity.CampaignAggregate{CampaignID: campaignID}
	}
	aggregate.ImpressionCount++
	aggregate.UpdatedAt = now
	return s.store.UpsertCampaignAggregate(ctx, aggregate)
}

func (s *AdsService) incrementClickAggregate(ctx context.Context, campaignID string, now time.Time) error {
	aggregate, err := s.store.GetCampaignAggregate(ctx, campaignID)
	if err != nil {
		if !errors.Is(err, adsrepository.ErrNotFound) {
			return err
		}
		aggregate = entity.CampaignAggregate{CampaignID: campaignID}
	}
	aggregate.ClickCount++
	aggregate.UpdatedAt = now
	return s.store.UpsertCampaignAggregate(ctx, aggregate)
}

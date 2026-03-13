package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
	adsrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/repository"
	"github.com/google/uuid"
)

func (s *AdsService) IntakeImpression(ctx context.Context, request dto.IntakeImpressionRequest) (dto.IntakeImpressionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.IntakeImpressionResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.IntakeImpressionResponse{}, err
	}
	if err := s.requireSurfaceEnabled(cfg.SurfaceEnabled); err != nil {
		return dto.IntakeImpressionResponse{}, err
	}
	if err := s.requirePlacementEnabled(cfg.PlacementEnabled); err != nil {
		return dto.IntakeImpressionResponse{}, err
	}
	if err := s.requireCampaignEnabled(cfg.CampaignEnabled); err != nil {
		return dto.IntakeImpressionResponse{}, err
	}

	now := s.now().UTC()
	dedupKey := buildImpressionDedupKey(request.SessionID, request.RequestID)
	dedupLog, err := s.store.GetImpressionDedup(ctx, dedupKey)
	if err == nil {
		return dto.IntakeImpressionResponse{
			Status:       "idempotent",
			ImpressionID: dedupLog.ImpressionID,
			CampaignID:   dedupLog.CampaignID,
			CreatedAt:    now,
		}, nil
	}
	if err != nil && !errors.Is(err, adsrepository.ErrNotFound) {
		return dto.IntakeImpressionResponse{}, err
	}

	_, campaign, err := s.ensurePlacementAndCampaignActive(ctx, request.PlacementID, request.CampaignID, now)
	if err != nil {
		return dto.IntakeImpressionResponse{}, err
	}

	impression := entity.ImpressionLog{
		ImpressionID: uuid.NewString(),
		RequestID:    request.RequestID,
		PlacementID:  request.PlacementID,
		CampaignID:   request.CampaignID,
		SessionID:    request.SessionID,
		UserID:       request.UserID,
		Status:       entity.ImpressionStatusAccepted,
		CreatedAt:    now,
	}
	if err := s.store.CreateImpressionLog(ctx, impression); err != nil {
		return dto.IntakeImpressionResponse{}, err
	}
	if err := s.store.PutImpressionDedup(ctx, dedupKey, impression); err != nil {
		return dto.IntakeImpressionResponse{}, err
	}
	if err := s.incrementImpressionAggregate(ctx, campaign.CampaignID, now); err != nil {
		return dto.IntakeImpressionResponse{}, err
	}

	return dto.IntakeImpressionResponse{
		Status:       "accepted",
		ImpressionID: impression.ImpressionID,
		CampaignID:   impression.CampaignID,
		CreatedAt:    now,
	}, nil
}

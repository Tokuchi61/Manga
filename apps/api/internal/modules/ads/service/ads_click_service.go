package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
	adsrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/repository"
	"github.com/google/uuid"
)

func (s *AdsService) IntakeClick(ctx context.Context, request dto.IntakeClickRequest) (dto.IntakeClickResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.IntakeClickResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.IntakeClickResponse{}, err
	}
	if err := s.requireSurfaceEnabled(cfg.SurfaceEnabled); err != nil {
		return dto.IntakeClickResponse{}, err
	}
	if err := s.requirePlacementEnabled(cfg.PlacementEnabled); err != nil {
		return dto.IntakeClickResponse{}, err
	}
	if err := s.requireCampaignEnabled(cfg.CampaignEnabled); err != nil {
		return dto.IntakeClickResponse{}, err
	}
	if err := s.requireClickIntakeEnabled(cfg.ClickIntakeEnabled); err != nil {
		return dto.IntakeClickResponse{}, err
	}

	now := s.now().UTC()
	dedupKey := buildClickDedupKey(request.SessionID, request.RequestID)
	dedupLog, err := s.store.GetClickDedup(ctx, dedupKey)
	if err == nil {
		return dto.IntakeClickResponse{
			Status:     "idempotent",
			ClickID:    dedupLog.ClickID,
			CampaignID: dedupLog.CampaignID,
			CreatedAt:  now,
		}, nil
	}
	if err != nil && !errors.Is(err, adsrepository.ErrNotFound) {
		return dto.IntakeClickResponse{}, err
	}

	_, campaign, err := s.ensurePlacementAndCampaignActive(ctx, request.PlacementID, request.CampaignID, now)
	if err != nil {
		return dto.IntakeClickResponse{}, err
	}

	clickStatus := entity.ClickStatusAccepted
	responseStatus := "accepted"
	if request.InvalidTraffic {
		clickStatus = entity.ClickStatusIgnored
		responseStatus = "ignored_invalid_traffic"
	}

	click := entity.ClickLog{
		ClickID:        uuid.NewString(),
		RequestID:      request.RequestID,
		PlacementID:    request.PlacementID,
		CampaignID:     request.CampaignID,
		SessionID:      request.SessionID,
		UserID:         request.UserID,
		Status:         clickStatus,
		InvalidTraffic: request.InvalidTraffic,
		CreatedAt:      now,
	}
	if err := s.store.CreateClickLog(ctx, click); err != nil {
		return dto.IntakeClickResponse{}, err
	}
	if err := s.store.PutClickDedup(ctx, dedupKey, click); err != nil {
		return dto.IntakeClickResponse{}, err
	}
	if clickStatus == entity.ClickStatusAccepted {
		if err := s.incrementClickAggregate(ctx, campaign.CampaignID, now); err != nil {
			return dto.IntakeClickResponse{}, err
		}
	}

	return dto.IntakeClickResponse{
		Status:     responseStatus,
		ClickID:    click.ClickID,
		CampaignID: click.CampaignID,
		CreatedAt:  now,
	}, nil
}

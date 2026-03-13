package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/dto"
	rprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/repository"
)

func (s *RoyalPassService) GetActorSeasonOverview(ctx context.Context, request dto.GetActorSeasonOverviewRequest) (dto.SeasonOverviewResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.SeasonOverviewResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.SeasonOverviewResponse{}, err
	}
	if err := s.requireSeasonEnabled(cfg.SeasonEnabled); err != nil {
		return dto.SeasonOverviewResponse{}, err
	}

	now := s.now().UTC()
	season, err := s.resolveSeason(ctx, request.SeasonID, now)
	if err != nil {
		return dto.SeasonOverviewResponse{}, err
	}

	tiers, err := s.store.ListTierDefinitions(ctx, season.SeasonID, "", false, 0, 0)
	if err != nil {
		return dto.SeasonOverviewResponse{}, err
	}

	progress, progressErr := s.store.GetUserProgress(ctx, request.ActorUserID, season.SeasonID)
	if progressErr != nil {
		if !errors.Is(progressErr, rprepository.ErrNotFound) {
			return dto.SeasonOverviewResponse{}, progressErr
		}
		progress = fallbackUserProgress(request.ActorUserID, season.SeasonID, now)
	}

	return toSeasonOverviewResponse(season, tiers, progress), nil
}

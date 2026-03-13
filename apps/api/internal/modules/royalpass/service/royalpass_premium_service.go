package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/dto"
	rprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/repository"
)

func (s *RoyalPassService) ActivatePremiumTrack(ctx context.Context, request dto.ActivatePremiumTrackRequest) (dto.ActivatePremiumTrackResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ActivatePremiumTrackResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ActivatePremiumTrackResponse{}, err
	}
	if err := s.requireSeasonEnabled(cfg.SeasonEnabled); err != nil {
		return dto.ActivatePremiumTrackResponse{}, err
	}
	if err := s.requirePremiumEnabled(cfg.PremiumEnabled); err != nil {
		return dto.ActivatePremiumTrackResponse{}, err
	}

	now := s.now().UTC()
	season, err := s.resolveSeason(ctx, request.SeasonID, now)
	if err != nil {
		return dto.ActivatePremiumTrackResponse{}, err
	}
	if !isSeasonWritable(season, now) {
		return dto.ActivatePremiumTrackResponse{}, ErrSeasonUnavailable
	}

	dedupKey := buildPremiumDedupKey(request.ActorUserID, season.SeasonID, request.SourceType, request.ActivationRef)
	dedupProgress, dedupErr := s.store.GetPremiumActivationDedup(ctx, dedupKey)
	if dedupErr == nil {
		return dto.ActivatePremiumTrackResponse{Status: "idempotent", SeasonID: season.SeasonID, PremiumActivated: dedupProgress.PremiumActivated}, nil
	}
	if dedupErr != nil && !errors.Is(dedupErr, rprepository.ErrNotFound) {
		return dto.ActivatePremiumTrackResponse{}, dedupErr
	}

	progress, progressErr := s.store.GetUserProgress(ctx, request.ActorUserID, season.SeasonID)
	if progressErr != nil {
		if !errors.Is(progressErr, rprepository.ErrNotFound) {
			return dto.ActivatePremiumTrackResponse{}, progressErr
		}
		progress = fallbackUserProgress(request.ActorUserID, season.SeasonID, now)
	}

	if progress.PremiumActivated && normalizeValue(progress.PremiumActivationRef) == normalizeValue(request.ActivationRef) {
		return dto.ActivatePremiumTrackResponse{Status: "idempotent", SeasonID: season.SeasonID, PremiumActivated: true}, nil
	}

	progress.PremiumActivated = true
	progress.PremiumActivationSource = request.SourceType
	progress.PremiumActivationRef = request.ActivationRef
	progress.LastRequestID = request.RequestID
	progress.LastCorrelationID = request.CorrelationID
	progress.UpdatedAt = now
	if progress.CreatedAt.IsZero() {
		progress.CreatedAt = now
	}

	if err := s.store.UpsertUserProgress(ctx, progress); err != nil {
		return dto.ActivatePremiumTrackResponse{}, err
	}
	if err := s.store.PutPremiumActivationDedup(ctx, dedupKey, progress); err != nil {
		return dto.ActivatePremiumTrackResponse{}, err
	}

	return dto.ActivatePremiumTrackResponse{Status: "activated", SeasonID: season.SeasonID, PremiumActivated: progress.PremiumActivated}, nil
}

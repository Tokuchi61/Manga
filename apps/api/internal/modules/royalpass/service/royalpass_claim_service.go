package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
	rprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/repository"
)

func (s *RoyalPassService) ClaimTierReward(ctx context.Context, request dto.ClaimTierRewardRequest) (dto.ClaimTierRewardResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ClaimTierRewardResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ClaimTierRewardResponse{}, err
	}
	if err := s.requireSeasonEnabled(cfg.SeasonEnabled); err != nil {
		return dto.ClaimTierRewardResponse{}, err
	}
	if err := s.requireClaimEnabled(cfg.ClaimEnabled); err != nil {
		return dto.ClaimTierRewardResponse{}, err
	}

	now := s.now().UTC()
	season, err := s.resolveSeason(ctx, request.SeasonID, now)
	if err != nil {
		return dto.ClaimTierRewardResponse{}, err
	}
	if !isSeasonWritable(season, now) {
		return dto.ClaimTierRewardResponse{}, ErrSeasonUnavailable
	}

	tier, err := s.store.GetTierDefinition(ctx, season.SeasonID, request.TierNumber, request.Track)
	if err != nil {
		if errors.Is(err, rprepository.ErrNotFound) {
			return dto.ClaimTierRewardResponse{}, ErrNotFound
		}
		return dto.ClaimTierRewardResponse{}, err
	}
	if !tier.Active {
		return dto.ClaimTierRewardResponse{}, ErrNotFound
	}
	if normalizeValue(tier.Track) == entity.TrackPremium {
		if err := s.requirePremiumEnabled(cfg.PremiumEnabled); err != nil {
			return dto.ClaimTierRewardResponse{}, err
		}
	}

	if strings.TrimSpace(request.RequestID) != "" {
		dedupKey := buildClaimDedupKey(request.ActorUserID, season.SeasonID, request.Track, request.TierNumber, request.RequestID)
		_, dedupErr := s.store.GetClaimDedup(ctx, dedupKey)
		if dedupErr == nil {
			return dto.ClaimTierRewardResponse{
				Status:         "idempotent",
				SeasonID:       season.SeasonID,
				TierNumber:     request.TierNumber,
				Track:          request.Track,
				RewardItemID:   tier.RewardItemID,
				RewardQuantity: tier.RewardQuantity,
			}, nil
		}
		if dedupErr != nil && !errors.Is(dedupErr, rprepository.ErrNotFound) {
			return dto.ClaimTierRewardResponse{}, dedupErr
		}
	}

	progress, progressErr := s.store.GetUserProgress(ctx, request.ActorUserID, season.SeasonID)
	if progressErr != nil {
		if !errors.Is(progressErr, rprepository.ErrNotFound) {
			return dto.ClaimTierRewardResponse{}, progressErr
		}
		progress = fallbackUserProgress(request.ActorUserID, season.SeasonID, now)
	}

	if progress.Points < tier.RequiredPoints {
		return dto.ClaimTierRewardResponse{}, ErrTierNotEligible
	}
	if normalizeValue(tier.Track) == entity.TrackPremium && !progress.PremiumActivated {
		return dto.ClaimTierRewardResponse{}, ErrForbiddenAction
	}
	if progress.ClaimedTiers == nil {
		progress.ClaimedTiers = make(map[string]time.Time)
	}

	claimKey := claimTierKey(tier.Track, tier.TierNumber)
	if _, exists := progress.ClaimedTiers[claimKey]; exists {
		return dto.ClaimTierRewardResponse{}, ErrTierAlreadyClaimed
	}
	progress.ClaimedTiers[claimKey] = now
	progress.LastRequestID = request.RequestID
	progress.LastCorrelationID = request.CorrelationID
	progress.UpdatedAt = now
	if progress.CreatedAt.IsZero() {
		progress.CreatedAt = now
	}

	if err := s.store.UpsertUserProgress(ctx, progress); err != nil {
		return dto.ClaimTierRewardResponse{}, err
	}

	if strings.TrimSpace(request.RequestID) != "" {
		dedupKey := buildClaimDedupKey(request.ActorUserID, season.SeasonID, request.Track, request.TierNumber, request.RequestID)
		if err := s.store.PutClaimDedup(ctx, dedupKey, progress); err != nil {
			return dto.ClaimTierRewardResponse{}, err
		}
	}

	return dto.ClaimTierRewardResponse{
		Status:         "claim_requested",
		SeasonID:       season.SeasonID,
		TierNumber:     tier.TierNumber,
		Track:          tier.Track,
		RewardItemID:   tier.RewardItemID,
		RewardQuantity: tier.RewardQuantity,
	}, nil
}

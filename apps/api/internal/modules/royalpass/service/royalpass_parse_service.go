package service

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
	rprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/repository"
)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func isSeasonWithinWindow(season entity.SeasonDefinition, now time.Time) bool {
	utcNow := now.UTC()
	if season.StartsAt != nil && utcNow.Before(season.StartsAt.UTC()) {
		return false
	}
	if season.EndsAt != nil && utcNow.After(season.EndsAt.UTC()) {
		return false
	}
	return true
}

func isSeasonWritable(season entity.SeasonDefinition, now time.Time) bool {
	if normalizeValue(season.State) != entity.SeasonStateActive {
		return false
	}
	return isSeasonWithinWindow(season, now)
}

func (s *RoyalPassService) resolveSeason(ctx context.Context, seasonID string, now time.Time) (entity.SeasonDefinition, error) {
	if strings.TrimSpace(seasonID) != "" {
		season, err := s.store.GetSeasonDefinition(ctx, seasonID)
		if err != nil {
			if errors.Is(err, rprepository.ErrNotFound) {
				return entity.SeasonDefinition{}, ErrNotFound
			}
			return entity.SeasonDefinition{}, err
		}
		return season, nil
	}

	season, err := s.store.ResolveCurrentSeason(ctx, now)
	if err != nil {
		if errors.Is(err, rprepository.ErrNotFound) {
			return entity.SeasonDefinition{}, ErrSeasonUnavailable
		}
		return entity.SeasonDefinition{}, err
	}
	return season, nil
}

func fallbackUserProgress(userID string, seasonID string, now time.Time) entity.UserProgress {
	return entity.UserProgress{
		UserID:           userID,
		SeasonID:         seasonID,
		Points:           0,
		PremiumActivated: false,
		ClaimedTiers:     make(map[string]time.Time),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func claimTierKey(track string, tierNumber int) string {
	return normalizeValue(track) + ":" + normalizeValue(strconv.Itoa(tierNumber))
}

func buildProgressDedupKey(userID string, seasonID string, requestID string) string {
	return normalizeValue(userID) + ":" + normalizeValue(seasonID) + ":" + normalizeValue(requestID)
}

func buildClaimDedupKey(userID string, seasonID string, track string, tierNumber int, requestID string) string {
	return normalizeValue(userID) + ":" + normalizeValue(seasonID) + ":" + normalizeValue(track) + ":" + normalizeValue(strconv.Itoa(tierNumber)) + ":" + normalizeValue(requestID)
}

func buildPremiumDedupKey(userID string, seasonID string, sourceType string, activationRef string) string {
	return normalizeValue(userID) + ":" + normalizeValue(seasonID) + ":" + normalizeValue(sourceType) + ":" + normalizeValue(activationRef)
}

func toTierProgressResponse(tier entity.TierDefinition, progress entity.UserProgress) dto.TierProgressResponse {
	claimedKey := claimTierKey(tier.Track, tier.TierNumber)
	claimedAt, claimed := progress.ClaimedTiers[claimedKey]
	unlocked := progress.Points >= tier.RequiredPoints
	var claimedAtPtr *time.Time
	if claimed {
		claimTime := claimedAt.UTC()
		claimedAtPtr = &claimTime
	}

	return dto.TierProgressResponse{
		TierNumber:     tier.TierNumber,
		Track:          tier.Track,
		RequiredPoints: tier.RequiredPoints,
		RewardItemID:   tier.RewardItemID,
		RewardQuantity: tier.RewardQuantity,
		Active:         tier.Active,
		Unlocked:       unlocked,
		Claimed:        claimed,
		ClaimedAt:      claimedAtPtr,
	}
}

func toSeasonOverviewResponse(season entity.SeasonDefinition, tiers []entity.TierDefinition, progress entity.UserProgress) dto.SeasonOverviewResponse {
	tierItems := make([]dto.TierProgressResponse, 0, len(tiers))
	for _, tier := range tiers {
		tierItems = append(tierItems, toTierProgressResponse(tier, progress))
	}

	sort.Slice(tierItems, func(i int, j int) bool {
		if tierItems[i].TierNumber == tierItems[j].TierNumber {
			return tierItems[i].Track < tierItems[j].Track
		}
		return tierItems[i].TierNumber < tierItems[j].TierNumber
	})

	return dto.SeasonOverviewResponse{
		SeasonID:         season.SeasonID,
		Title:            season.Title,
		State:            season.State,
		StartsAt:         season.StartsAt,
		EndsAt:           season.EndsAt,
		Points:           progress.Points,
		PremiumActivated: progress.PremiumActivated,
		Items:            tierItems,
		Count:            len(tierItems),
	}
}

func toSeasonDefinitionResponse(season entity.SeasonDefinition) dto.SeasonDefinitionResponse {
	return dto.SeasonDefinitionResponse{
		SeasonID:  season.SeasonID,
		Title:     season.Title,
		State:     season.State,
		StartsAt:  season.StartsAt,
		EndsAt:    season.EndsAt,
		UpdatedAt: season.UpdatedAt,
	}
}

func toTierDefinitionResponse(tier entity.TierDefinition) dto.TierDefinitionResponse {
	return dto.TierDefinitionResponse{
		SeasonID:       tier.SeasonID,
		TierNumber:     tier.TierNumber,
		Track:          tier.Track,
		RequiredPoints: tier.RequiredPoints,
		RewardItemID:   tier.RewardItemID,
		RewardQuantity: tier.RewardQuantity,
		Active:         tier.Active,
		UpdatedAt:      tier.UpdatedAt,
	}
}

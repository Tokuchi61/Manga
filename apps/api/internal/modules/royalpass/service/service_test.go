package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/dto"
	rprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRoyalPassServiceOverviewProgressClaimAndPremiumFlow(t *testing.T) {
	store := rprepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	actorID := uuid.NewString()
	setupActiveSeasonAndTierSet(t, svc, now)

	ingestRes, err := svc.IngestRoyalPassProgress(context.Background(), dto.IngestRoyalPassProgressRequest{
		ActorUserID: actorID,
		SeasonID:    "season_2026_03",
		Delta:       60,
		SourceType:  "mission",
		RequestID:   "req-rp-progress-1",
	})
	require.NoError(t, err)
	require.Equal(t, 60, ingestRes.Points)

	freeClaimRes, err := svc.ClaimTierReward(context.Background(), dto.ClaimTierRewardRequest{
		ActorUserID: actorID,
		SeasonID:    "season_2026_03",
		TierNumber:  1,
		Track:       "free",
		RequestID:   "req-rp-claim-free-1",
	})
	require.NoError(t, err)
	require.Equal(t, "claim_requested", freeClaimRes.Status)

	_, err = svc.ClaimTierReward(context.Background(), dto.ClaimTierRewardRequest{
		ActorUserID: actorID,
		SeasonID:    "season_2026_03",
		TierNumber:  1,
		Track:       "premium",
	})
	require.True(t, errors.Is(err, ErrForbiddenAction))

	activateRes, err := svc.ActivatePremiumTrack(context.Background(), dto.ActivatePremiumTrackRequest{
		ActorUserID:   actorID,
		SeasonID:      "season_2026_03",
		SourceType:    "shop",
		ActivationRef: "shop-activation-1",
	})
	require.NoError(t, err)
	require.True(t, activateRes.PremiumActivated)

	premiumClaimRes, err := svc.ClaimTierReward(context.Background(), dto.ClaimTierRewardRequest{
		ActorUserID: actorID,
		SeasonID:    "season_2026_03",
		TierNumber:  1,
		Track:       "premium",
		RequestID:   "req-rp-claim-premium-1",
	})
	require.NoError(t, err)
	require.Equal(t, "claim_requested", premiumClaimRes.Status)

	overviewRes, err := svc.GetActorSeasonOverview(context.Background(), dto.GetActorSeasonOverviewRequest{
		ActorUserID: actorID,
		SeasonID:    "season_2026_03",
	})
	require.NoError(t, err)
	require.Equal(t, 2, overviewRes.Count)
	require.Equal(t, 60, overviewRes.Points)
	require.True(t, overviewRes.PremiumActivated)
}

func TestRoyalPassServiceRuntimeTogglesAffectFlows(t *testing.T) {
	store := rprepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }
	actorID := uuid.NewString()
	setupActiveSeasonAndTierSet(t, svc, now)

	_, err := svc.UpdateSeasonState(context.Background(), dto.UpdateSeasonStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.GetActorSeasonOverview(context.Background(), dto.GetActorSeasonOverviewRequest{ActorUserID: actorID, SeasonID: "season_2026_03"})
	require.True(t, errors.Is(err, ErrSeasonDisabled))

	_, err = svc.UpdateSeasonState(context.Background(), dto.UpdateSeasonStateRequest{Enabled: true})
	require.NoError(t, err)
	_, err = svc.IngestRoyalPassProgress(context.Background(), dto.IngestRoyalPassProgressRequest{ActorUserID: actorID, SeasonID: "season_2026_03", Delta: 60, SourceType: "mission"})
	require.NoError(t, err)

	_, err = svc.UpdateClaimState(context.Background(), dto.UpdateClaimStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.ClaimTierReward(context.Background(), dto.ClaimTierRewardRequest{ActorUserID: actorID, SeasonID: "season_2026_03", TierNumber: 1, Track: "free"})
	require.True(t, errors.Is(err, ErrClaimDisabled))

	_, err = svc.UpdatePremiumState(context.Background(), dto.UpdatePremiumStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.ActivatePremiumTrack(context.Background(), dto.ActivatePremiumTrackRequest{ActorUserID: actorID, SeasonID: "season_2026_03", SourceType: "shop", ActivationRef: "act-2"})
	require.True(t, errors.Is(err, ErrPremiumDisabled))
}

func TestRoyalPassServiceResetProgress(t *testing.T) {
	store := rprepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }
	actorID := uuid.NewString()
	setupActiveSeasonAndTierSet(t, svc, now)

	_, err := svc.IngestRoyalPassProgress(context.Background(), dto.IngestRoyalPassProgressRequest{ActorUserID: actorID, SeasonID: "season_2026_03", Delta: 30, SourceType: "mission"})
	require.NoError(t, err)

	resetRes, err := svc.ResetRoyalPassProgress(context.Background(), dto.ResetRoyalPassProgressRequest{TargetUserID: actorID, SeasonID: "season_2026_03"})
	require.NoError(t, err)
	require.Equal(t, 1, resetRes.DeletedCount)

	overviewRes, err := svc.GetActorSeasonOverview(context.Background(), dto.GetActorSeasonOverviewRequest{ActorUserID: actorID, SeasonID: "season_2026_03"})
	require.NoError(t, err)
	require.Equal(t, 0, overviewRes.Points)
}

func setupActiveSeasonAndTierSet(t *testing.T, svc *RoyalPassService, now time.Time) {
	t.Helper()
	start := now.Add(-24 * time.Hour)
	end := now.Add(24 * time.Hour)

	_, err := svc.UpsertSeasonDefinition(context.Background(), dto.UpsertSeasonDefinitionRequest{
		SeasonID: "season_2026_03",
		Title:    "Spring Launch",
		State:    "active",
		StartsAt: &start,
		EndsAt:   &end,
	})
	require.NoError(t, err)

	_, err = svc.UpsertTierDefinition(context.Background(), dto.UpsertTierDefinitionRequest{
		SeasonID:       "season_2026_03",
		TierNumber:     1,
		Track:          "free",
		RequiredPoints: 50,
		RewardItemID:   "mana_100",
		RewardQuantity: 1,
		Active:         true,
	})
	require.NoError(t, err)

	_, err = svc.UpsertTierDefinition(context.Background(), dto.UpsertTierDefinitionRequest{
		SeasonID:       "season_2026_03",
		TierNumber:     1,
		Track:          "premium",
		RequiredPoints: 50,
		RewardItemID:   "avatar_frame_gold",
		RewardQuantity: 1,
		Active:         true,
	})
	require.NoError(t, err)
}

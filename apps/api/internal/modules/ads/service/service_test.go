package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/dto"
	adsrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/stretchr/testify/require"
)

func TestAdsServiceResolveImpressionClickAndAggregateFlow(t *testing.T) {
	store := adsrepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 21, 13, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	setupPlacementAndCampaign(t, svc, now)

	resolveRes, err := svc.ResolvePlacements(context.Background(), dto.ResolvePlacementsRequest{
		Surface:    "home",
		TargetType: "none",
		SessionID:  "session-1",
	})
	require.NoError(t, err)
	require.Equal(t, 1, resolveRes.Count)

	impressionRes, err := svc.IntakeImpression(context.Background(), dto.IntakeImpressionRequest{
		RequestID:   "req-imp-1",
		PlacementID: "home_top",
		CampaignID:  "campaign_home_1",
		SessionID:   "session-1",
	})
	require.NoError(t, err)
	require.Equal(t, "accepted", impressionRes.Status)

	impressionIdempotentRes, err := svc.IntakeImpression(context.Background(), dto.IntakeImpressionRequest{
		RequestID:   "req-imp-1",
		PlacementID: "home_top",
		CampaignID:  "campaign_home_1",
		SessionID:   "session-1",
	})
	require.NoError(t, err)
	require.Equal(t, "idempotent", impressionIdempotentRes.Status)

	clickRes, err := svc.IntakeClick(context.Background(), dto.IntakeClickRequest{
		RequestID:   "req-click-1",
		PlacementID: "home_top",
		CampaignID:  "campaign_home_1",
		SessionID:   "session-1",
	})
	require.NoError(t, err)
	require.Equal(t, "accepted", clickRes.Status)

	aggRes, err := svc.ListCampaignAggregate(context.Background(), 10, 0)
	require.NoError(t, err)
	require.Equal(t, 1, aggRes.Count)
	require.Equal(t, 1, aggRes.Items[0].ImpressionCount)
	require.Equal(t, 1, aggRes.Items[0].ClickCount)
}

func TestAdsServiceRuntimeTogglesAffectFlows(t *testing.T) {
	store := adsrepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 21, 13, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }
	setupPlacementAndCampaign(t, svc, now)

	_, err := svc.UpdateCampaignState(context.Background(), dto.UpdateCampaignStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.ResolvePlacements(context.Background(), dto.ResolvePlacementsRequest{Surface: "home", TargetType: "none"})
	require.True(t, errors.Is(err, ErrCampaignDisabled))

	_, err = svc.UpdateCampaignState(context.Background(), dto.UpdateCampaignStateRequest{Enabled: true})
	require.NoError(t, err)
	_, err = svc.UpdateClickIntakeState(context.Background(), dto.UpdateClickIntakeStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.IntakeClick(context.Background(), dto.IntakeClickRequest{
		RequestID:   "req-click-disabled",
		PlacementID: "home_top",
		CampaignID:  "campaign_home_1",
		SessionID:   "session-1",
	})
	require.True(t, errors.Is(err, ErrClickIntakeDisabled))
}

func setupPlacementAndCampaign(t *testing.T, svc *AdsService, now time.Time) {
	t.Helper()

	_, err := svc.UpsertPlacementDefinition(context.Background(), dto.UpsertPlacementDefinitionRequest{
		PlacementID:  "home_top",
		Surface:      "home",
		TargetType:   "none",
		Visible:      true,
		Priority:     100,
		FrequencyCap: 3,
	})
	require.NoError(t, err)

	startsAt := now.Add(-2 * time.Hour)
	endsAt := now.Add(24 * time.Hour)
	_, err = svc.UpsertCampaignDefinition(context.Background(), dto.UpsertCampaignDefinitionRequest{
		CampaignID:  "campaign_home_1",
		PlacementID: "home_top",
		Name:        "Launch Campaign",
		State:       "active",
		CreativeURL: "https://cdn.example.com/banner.png",
		ClickURL:    "https://example.com/click",
		Weight:      50,
		StartsAt:    &startsAt,
		EndsAt:      &endsAt,
	})
	require.NoError(t, err)
}

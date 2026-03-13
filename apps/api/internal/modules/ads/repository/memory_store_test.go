package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorePlacementCampaignIntakeAndSnapshot(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 21, 12, 0, 0, 0, time.UTC)

	err := store.UpsertPlacementDefinition(context.Background(), entity.PlacementDefinition{
		PlacementID:  "home_top",
		Surface:      "home",
		TargetType:   "none",
		Visible:      true,
		Priority:     100,
		FrequencyCap: 3,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	require.NoError(t, err)

	err = store.UpsertCampaignDefinition(context.Background(), entity.CampaignDefinition{
		CampaignID:  "campaign_home_1",
		PlacementID: "home_top",
		Name:        "Launch Campaign",
		State:       entity.CampaignStateActive,
		CreativeURL: "https://cdn.example.com/banner.png",
		ClickURL:    "https://example.com/click",
		Weight:      50,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	require.NoError(t, err)

	impression := entity.ImpressionLog{
		ImpressionID: uuid.NewString(),
		RequestID:    "req-impression-1",
		PlacementID:  "home_top",
		CampaignID:   "campaign_home_1",
		SessionID:    "session-1",
		Status:       entity.ImpressionStatusAccepted,
		CreatedAt:    now,
	}
	require.NoError(t, store.CreateImpressionLog(context.Background(), impression))
	require.NoError(t, store.PutImpressionDedup(context.Background(), "imp:session-1:req-impression-1", impression))

	click := entity.ClickLog{
		ClickID:     uuid.NewString(),
		RequestID:   "req-click-1",
		PlacementID: "home_top",
		CampaignID:  "campaign_home_1",
		SessionID:   "session-1",
		Status:      entity.ClickStatusAccepted,
		CreatedAt:   now,
	}
	require.NoError(t, store.CreateClickLog(context.Background(), click))
	require.NoError(t, store.PutClickDedup(context.Background(), "clk:session-1:req-click-1", click))

	require.NoError(t, store.UpsertCampaignAggregate(context.Background(), entity.CampaignAggregate{
		CampaignID:      "campaign_home_1",
		ImpressionCount: 1,
		ClickCount:      1,
		UpdatedAt:       now,
	}))

	payload, err := store.Snapshot()
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	restored := NewMemoryStore()
	require.NoError(t, restored.RestoreSnapshot(payload))

	campaign, err := restored.GetCampaignDefinition(context.Background(), "campaign_home_1")
	require.NoError(t, err)
	require.Equal(t, "Launch Campaign", campaign.Name)

	agg, err := restored.GetCampaignAggregate(context.Background(), "campaign_home_1")
	require.NoError(t, err)
	require.Equal(t, 1, agg.ClickCount)
}

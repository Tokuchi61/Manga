package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
)

func (s *MemoryStore) UpsertCampaignDefinition(_ context.Context, campaign entity.CampaignDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.campaignsByID[normalizeValue(campaign.CampaignID)] = cloneCampaign(campaign)
	return nil
}

func (s *MemoryStore) GetCampaignDefinition(_ context.Context, campaignID string) (entity.CampaignDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	campaign, ok := s.campaignsByID[normalizeValue(campaignID)]
	if !ok {
		return entity.CampaignDefinition{}, ErrNotFound
	}
	return cloneCampaign(campaign), nil
}

func (s *MemoryStore) ListCampaignDefinitions(_ context.Context, placementID string, state string, limit int, offset int) ([]entity.CampaignDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedPlacementID := normalizeValue(placementID)
	normalizedState := normalizeValue(state)
	items := make([]entity.CampaignDefinition, 0, len(s.campaignsByID))
	for _, campaign := range s.campaignsByID {
		if normalizedPlacementID != "" && normalizeValue(campaign.PlacementID) != normalizedPlacementID {
			continue
		}
		if normalizedState != "" && normalizeValue(campaign.State) != normalizedState {
			continue
		}
		items = append(items, cloneCampaign(campaign))
	}

	sortByUpdatedDesc(items, func(i int, j int) bool {
		if items[i].Weight == items[j].Weight {
			return items[i].UpdatedAt.After(items[j].UpdatedAt)
		}
		return items[i].Weight > items[j].Weight
	})

	return applyOffsetLimit(items, offset, limit), nil
}

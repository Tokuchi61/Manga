package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
)

func (s *MemoryStore) GetCampaignAggregate(_ context.Context, campaignID string) (entity.CampaignAggregate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	aggregate, ok := s.aggregateByCampaign[normalizeValue(campaignID)]
	if !ok {
		return entity.CampaignAggregate{}, ErrNotFound
	}
	return cloneAggregate(aggregate), nil
}

func (s *MemoryStore) UpsertCampaignAggregate(_ context.Context, aggregate entity.CampaignAggregate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.aggregateByCampaign[normalizeValue(aggregate.CampaignID)] = cloneAggregate(aggregate)
	return nil
}

func (s *MemoryStore) ListCampaignAggregate(_ context.Context, limit int, offset int) ([]entity.CampaignAggregate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]entity.CampaignAggregate, 0, len(s.aggregateByCampaign))
	for _, aggregate := range s.aggregateByCampaign {
		items = append(items, cloneAggregate(aggregate))
	}

	sortByUpdatedDesc(items, func(i int, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return applyOffsetLimit(items, offset, limit), nil
}

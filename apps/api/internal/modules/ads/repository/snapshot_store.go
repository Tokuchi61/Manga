package repository

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/entity"
)

type adsSnapshotState struct {
	PlacementsByID       map[string]entity.PlacementDefinition
	CampaignsByID        map[string]entity.CampaignDefinition
	ImpressionsByID      map[string]entity.ImpressionLog
	ClicksByID           map[string]entity.ClickLog
	ImpressionDedupByKey map[string]entity.ImpressionLog
	ClickDedupByKey      map[string]entity.ClickLog
	AggregateByCampaign  map[string]entity.CampaignAggregate
	RuntimeConfig        entity.RuntimeConfig
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := adsSnapshotState{
		PlacementsByID:       s.placementsByID,
		CampaignsByID:        s.campaignsByID,
		ImpressionsByID:      s.impressionsByID,
		ClicksByID:           s.clicksByID,
		ImpressionDedupByKey: s.impressionDedupByKey,
		ClickDedupByKey:      s.clickDedupByKey,
		AggregateByCampaign:  s.aggregateByCampaign,
		RuntimeConfig:        s.runtimeConfig,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(state); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *MemoryStore) RestoreSnapshot(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var state adsSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.PlacementsByID == nil {
		state.PlacementsByID = make(map[string]entity.PlacementDefinition)
	}
	if state.CampaignsByID == nil {
		state.CampaignsByID = make(map[string]entity.CampaignDefinition)
	}
	if state.ImpressionsByID == nil {
		state.ImpressionsByID = make(map[string]entity.ImpressionLog)
	}
	if state.ClicksByID == nil {
		state.ClicksByID = make(map[string]entity.ClickLog)
	}
	if state.ImpressionDedupByKey == nil {
		state.ImpressionDedupByKey = make(map[string]entity.ImpressionLog)
	}
	if state.ClickDedupByKey == nil {
		state.ClickDedupByKey = make(map[string]entity.ClickLog)
	}
	if state.AggregateByCampaign == nil {
		state.AggregateByCampaign = make(map[string]entity.CampaignAggregate)
	}
	if state.RuntimeConfig.UpdatedAt.IsZero() {
		state.RuntimeConfig = entity.RuntimeConfig{
			SurfaceEnabled:     true,
			PlacementEnabled:   true,
			CampaignEnabled:    true,
			ClickIntakeEnabled: true,
			UpdatedAt:          time.Now().UTC(),
		}
	}

	s.placementsByID = state.PlacementsByID
	s.campaignsByID = state.CampaignsByID
	s.impressionsByID = state.ImpressionsByID
	s.clicksByID = state.ClicksByID
	s.impressionDedupByKey = state.ImpressionDedupByKey
	s.clickDedupByKey = state.ClickDedupByKey
	s.aggregateByCampaign = state.AggregateByCampaign
	s.runtimeConfig = state.RuntimeConfig

	return nil
}

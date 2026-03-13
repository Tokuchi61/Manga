package repository

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/entity"
)

type royalPassSnapshotState struct {
	SeasonDefinitionsByID       map[string]entity.SeasonDefinition
	TierDefinitionsByKey        map[string]entity.TierDefinition
	UserProgressByKey           map[string]entity.UserProgress
	ProgressDedupByKey          map[string]entity.UserProgress
	ClaimDedupByKey             map[string]entity.UserProgress
	PremiumActivationDedupByKey map[string]entity.UserProgress
	RuntimeConfig               entity.RuntimeConfig
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := royalPassSnapshotState{
		SeasonDefinitionsByID:       s.seasonDefinitionsByID,
		TierDefinitionsByKey:        s.tierDefinitionsByKey,
		UserProgressByKey:           s.userProgressByKey,
		ProgressDedupByKey:          s.progressDedupByKey,
		ClaimDedupByKey:             s.claimDedupByKey,
		PremiumActivationDedupByKey: s.premiumActivationDedupByKey,
		RuntimeConfig:               s.runtimeConfig,
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

	var state royalPassSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.SeasonDefinitionsByID == nil {
		state.SeasonDefinitionsByID = make(map[string]entity.SeasonDefinition)
	}
	if state.TierDefinitionsByKey == nil {
		state.TierDefinitionsByKey = make(map[string]entity.TierDefinition)
	}
	if state.UserProgressByKey == nil {
		state.UserProgressByKey = make(map[string]entity.UserProgress)
	}
	if state.ProgressDedupByKey == nil {
		state.ProgressDedupByKey = make(map[string]entity.UserProgress)
	}
	if state.ClaimDedupByKey == nil {
		state.ClaimDedupByKey = make(map[string]entity.UserProgress)
	}
	if state.PremiumActivationDedupByKey == nil {
		state.PremiumActivationDedupByKey = make(map[string]entity.UserProgress)
	}
	if state.RuntimeConfig.UpdatedAt.IsZero() {
		state.RuntimeConfig = entity.RuntimeConfig{
			SeasonEnabled:  true,
			ClaimEnabled:   true,
			PremiumEnabled: true,
			UpdatedAt:      time.Now().UTC(),
		}
	}

	s.seasonDefinitionsByID = state.SeasonDefinitionsByID
	s.tierDefinitionsByKey = state.TierDefinitionsByKey
	s.userProgressByKey = state.UserProgressByKey
	s.progressDedupByKey = state.ProgressDedupByKey
	s.claimDedupByKey = state.ClaimDedupByKey
	s.premiumActivationDedupByKey = state.PremiumActivationDedupByKey
	s.runtimeConfig = state.RuntimeConfig

	return nil
}

package repository

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
)

type missionSnapshotState struct {
	MissionDefinitionsByID map[string]entity.MissionDefinition
	MissionProgressByKey   map[string]entity.MissionProgress
	ProgressDedupByKey     map[string]entity.MissionProgress
	ClaimDedupByKey        map[string]entity.MissionProgress
	RuntimeConfig          entity.RuntimeConfig
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := missionSnapshotState{
		MissionDefinitionsByID: s.missionDefinitionsByID,
		MissionProgressByKey:   s.missionProgressByKey,
		ProgressDedupByKey:     s.progressDedupByKey,
		ClaimDedupByKey:        s.claimDedupByKey,
		RuntimeConfig:          s.runtimeConfig,
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

	var state missionSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.MissionDefinitionsByID == nil {
		state.MissionDefinitionsByID = make(map[string]entity.MissionDefinition)
	}
	if state.MissionProgressByKey == nil {
		state.MissionProgressByKey = make(map[string]entity.MissionProgress)
	}
	if state.ProgressDedupByKey == nil {
		state.ProgressDedupByKey = make(map[string]entity.MissionProgress)
	}
	if state.ClaimDedupByKey == nil {
		state.ClaimDedupByKey = make(map[string]entity.MissionProgress)
	}
	if state.RuntimeConfig.UpdatedAt.IsZero() {
		state.RuntimeConfig = entity.RuntimeConfig{
			ReadEnabled:           true,
			ClaimEnabled:          true,
			ProgressIngestEnabled: true,
			DailyResetHourUTC:     0,
			UpdatedAt:             time.Now().UTC(),
		}
	}

	s.missionDefinitionsByID = state.MissionDefinitionsByID
	s.missionProgressByKey = state.MissionProgressByKey
	s.progressDedupByKey = state.ProgressDedupByKey
	s.claimDedupByKey = state.ClaimDedupByKey
	s.runtimeConfig = state.RuntimeConfig

	return nil
}

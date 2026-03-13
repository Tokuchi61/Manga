package repository

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
)

type adminSnapshotState struct {
	RuntimeConfig     entity.RuntimeConfig
	ActionsByID       map[string]entity.AdminAction
	ActionDedupByKey  map[string]entity.AdminAction
	OverridesByID     map[string]entity.OverrideRecord
	UserReviewsByID   map[string]entity.UserReviewRecord
	ImpersonationByID map[string]entity.ImpersonationSession
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := adminSnapshotState{
		RuntimeConfig:     s.runtimeConfig,
		ActionsByID:       s.actionsByID,
		ActionDedupByKey:  s.actionDedupByKey,
		OverridesByID:     s.overridesByID,
		UserReviewsByID:   s.userReviewsByID,
		ImpersonationByID: s.impersonationByID,
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

	var state adminSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.ActionsByID == nil {
		state.ActionsByID = make(map[string]entity.AdminAction)
	}
	if state.ActionDedupByKey == nil {
		state.ActionDedupByKey = make(map[string]entity.AdminAction)
	}
	if state.OverridesByID == nil {
		state.OverridesByID = make(map[string]entity.OverrideRecord)
	}
	if state.UserReviewsByID == nil {
		state.UserReviewsByID = make(map[string]entity.UserReviewRecord)
	}
	if state.ImpersonationByID == nil {
		state.ImpersonationByID = make(map[string]entity.ImpersonationSession)
	}
	if state.RuntimeConfig.UpdatedAt.IsZero() {
		state.RuntimeConfig = entity.RuntimeConfig{
			MaintenanceEnabled: false,
			UpdatedAt:          time.Now().UTC(),
		}
	}

	s.runtimeConfig = state.RuntimeConfig
	s.actionsByID = state.ActionsByID
	s.actionDedupByKey = state.ActionDedupByKey
	s.overridesByID = state.OverridesByID
	s.userReviewsByID = state.UserReviewsByID
	s.impersonationByID = state.ImpersonationByID

	return nil
}

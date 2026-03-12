package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
)

type moderationSnapshotState struct {
	CasesByID      map[string]entity.Case
	SourceRefIndex map[string]string
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := moderationSnapshotState{
		CasesByID:      s.casesByID,
		SourceRefIndex: s.sourceRefIndex,
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

	var state moderationSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.CasesByID == nil {
		state.CasesByID = make(map[string]entity.Case)
	}
	if state.SourceRefIndex == nil {
		state.SourceRefIndex = make(map[string]string)
	}

	s.casesByID = state.CasesByID
	s.sourceRefIndex = state.SourceRefIndex
	return nil
}

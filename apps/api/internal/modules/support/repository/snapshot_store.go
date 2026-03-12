package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
)

type supportSnapshotState struct {
	CasesByID      map[string]entity.SupportCase
	RequestIDIndex map[string]map[string]string
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := supportSnapshotState{
		CasesByID:      s.casesByID,
		RequestIDIndex: s.requestIDIndex,
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

	var state supportSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.CasesByID == nil {
		state.CasesByID = make(map[string]entity.SupportCase)
	}
	if state.RequestIDIndex == nil {
		state.RequestIDIndex = make(map[string]map[string]string)
	}

	s.casesByID = state.CasesByID
	s.requestIDIndex = state.RequestIDIndex

	return nil
}

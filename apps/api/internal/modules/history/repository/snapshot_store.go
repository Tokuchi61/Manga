package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
)

type historySnapshotState struct {
	LibraryByKey    map[string]entity.LibraryEntry
	TimelineByUser  map[string][]entity.TimelineEvent
	CheckpointDedup map[string]string
	RuntimeConfig   entity.RuntimeConfig
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := historySnapshotState{
		LibraryByKey:    s.libraryByKey,
		TimelineByUser:  s.timelineByUser,
		CheckpointDedup: s.checkpointDedup,
		RuntimeConfig:   s.runtimeConfig,
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

	var state historySnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.LibraryByKey == nil {
		state.LibraryByKey = make(map[string]entity.LibraryEntry)
	}
	if state.TimelineByUser == nil {
		state.TimelineByUser = make(map[string][]entity.TimelineEvent)
	}
	if state.CheckpointDedup == nil {
		state.CheckpointDedup = make(map[string]string)
	}
	if state.RuntimeConfig.UpdatedAt.IsZero() {
		state.RuntimeConfig.ContinueReadingEnabled = true
		state.RuntimeConfig.LibraryEnabled = true
		state.RuntimeConfig.TimelineEnabled = true
		state.RuntimeConfig.BookmarkWriteEnabled = true
	}

	s.libraryByKey = state.LibraryByKey
	s.timelineByUser = state.TimelineByUser
	s.checkpointDedup = state.CheckpointDedup
	s.runtimeConfig = state.RuntimeConfig

	return nil
}

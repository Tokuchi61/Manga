package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
)

type chapterSnapshotState struct {
	ChaptersByID  map[string]entity.Chapter
	SlugIndex     map[string]map[string]string
	SequenceIndex map[string]map[int]string
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := chapterSnapshotState{
		ChaptersByID:  s.chaptersByID,
		SlugIndex:     s.slugIndex,
		SequenceIndex: s.sequenceIndex,
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

	var state chapterSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.ChaptersByID == nil {
		state.ChaptersByID = make(map[string]entity.Chapter)
	}
	if state.SlugIndex == nil {
		state.SlugIndex = make(map[string]map[string]string)
	}
	if state.SequenceIndex == nil {
		state.SequenceIndex = make(map[string]map[int]string)
	}

	s.chaptersByID = state.ChaptersByID
	s.slugIndex = state.SlugIndex
	s.sequenceIndex = state.SequenceIndex

	return nil
}

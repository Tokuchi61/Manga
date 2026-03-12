package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/entity"
)

type commentSnapshotState struct {
	CommentsByID map[string]entity.Comment
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := commentSnapshotState{CommentsByID: s.commentsByID}

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

	var state commentSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.CommentsByID == nil {
		state.CommentsByID = make(map[string]entity.Comment)
	}

	s.commentsByID = state.CommentsByID

	return nil
}

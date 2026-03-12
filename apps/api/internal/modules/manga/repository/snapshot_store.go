package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
)

type mangaSnapshotState struct {
	MangaByID map[string]entity.Manga
	SlugIndex map[string]string
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := mangaSnapshotState{
		MangaByID: s.mangaByID,
		SlugIndex: s.slugIndex,
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

	var state mangaSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.MangaByID == nil {
		state.MangaByID = make(map[string]entity.Manga)
	}
	if state.SlugIndex == nil {
		state.SlugIndex = make(map[string]string)
	}

	s.mangaByID = state.MangaByID
	s.slugIndex = state.SlugIndex

	return nil
}

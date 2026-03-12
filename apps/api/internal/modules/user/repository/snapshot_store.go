package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
)

type userSnapshotState struct {
	UsersByID           map[string]entity.UserAccount
	CredentialUserIndex map[string]string
	UsernameUserIndex   map[string]string
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := userSnapshotState{
		UsersByID:           s.usersByID,
		CredentialUserIndex: s.credentialUserIndex,
		UsernameUserIndex:   s.usernameUserIndex,
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

	var state userSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.UsersByID == nil {
		state.UsersByID = make(map[string]entity.UserAccount)
	}
	if state.CredentialUserIndex == nil {
		state.CredentialUserIndex = make(map[string]string)
	}
	if state.UsernameUserIndex == nil {
		state.UsernameUserIndex = make(map[string]string)
	}

	s.usersByID = state.UsersByID
	s.credentialUserIndex = state.CredentialUserIndex
	s.usernameUserIndex = state.UsernameUserIndex

	return nil
}

package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/google/uuid"
)

type authSnapshotState struct {
	CredentialsByID      map[uuid.UUID]entity.Credential
	CredentialEmailIdx   map[string]uuid.UUID
	SessionsByID         map[uuid.UUID]entity.Session
	SessionCredentialIdx map[uuid.UUID]map[uuid.UUID]struct{}
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := authSnapshotState{
		CredentialsByID:      s.credentialsByID,
		CredentialEmailIdx:   s.credentialEmailIdx,
		SessionsByID:         s.sessionsByID,
		SessionCredentialIdx: s.sessionCredentialIdx,
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

	var state authSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.CredentialsByID == nil {
		state.CredentialsByID = make(map[uuid.UUID]entity.Credential)
	}
	if state.CredentialEmailIdx == nil {
		state.CredentialEmailIdx = make(map[string]uuid.UUID)
	}
	if state.SessionsByID == nil {
		state.SessionsByID = make(map[uuid.UUID]entity.Session)
	}
	if state.SessionCredentialIdx == nil {
		state.SessionCredentialIdx = make(map[uuid.UUID]map[uuid.UUID]struct{})
	}

	s.credentialsByID = state.CredentialsByID
	s.credentialEmailIdx = state.CredentialEmailIdx
	s.sessionsByID = state.SessionsByID
	s.sessionCredentialIdx = state.SessionCredentialIdx

	// Sensitive token and security-event datasets are intentionally rebuilt in-memory only.
	s.tokensByID = make(map[uuid.UUID]entity.Token)
	s.tokenHashTypeIndex = make(map[string]uuid.UUID)
	s.securityEvents = make([]entity.SecurityEvent, 0)

	return nil
}

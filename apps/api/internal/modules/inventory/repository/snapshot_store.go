package repository

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/entity"
)

type inventorySnapshotState struct {
	ItemDefinitionsByID map[string]entity.ItemDefinition
	InventoryByKey      map[string]entity.InventoryEntry
	GrantDedupByKey     map[string]entity.InventoryEntry
	RuntimeConfig       entity.RuntimeConfig
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := inventorySnapshotState{
		ItemDefinitionsByID: s.itemDefinitionsByID,
		InventoryByKey:      s.inventoryByKey,
		GrantDedupByKey:     s.grantDedupByKey,
		RuntimeConfig:       s.runtimeConfig,
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

	var state inventorySnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.ItemDefinitionsByID == nil {
		state.ItemDefinitionsByID = make(map[string]entity.ItemDefinition)
	}
	if state.InventoryByKey == nil {
		state.InventoryByKey = make(map[string]entity.InventoryEntry)
	}
	if state.GrantDedupByKey == nil {
		state.GrantDedupByKey = make(map[string]entity.InventoryEntry)
	}
	if state.RuntimeConfig.UpdatedAt.IsZero() {
		state.RuntimeConfig = entity.RuntimeConfig{
			ReadEnabled:    true,
			ClaimEnabled:   true,
			ConsumeEnabled: true,
			EquipEnabled:   true,
			UpdatedAt:      time.Now().UTC(),
		}
	}

	s.itemDefinitionsByID = state.ItemDefinitionsByID
	s.inventoryByKey = state.InventoryByKey
	s.grantDedupByKey = state.GrantDedupByKey
	s.runtimeConfig = state.RuntimeConfig

	return nil
}

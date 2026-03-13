package repository

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
)

type shopSnapshotState struct {
	ProductDefinitionsByID map[string]entity.ProductDefinition
	OfferDefinitionsByID   map[string]entity.OfferDefinition
	PurchaseIntentsByID    map[string]entity.PurchaseIntent
	PurchaseDedupByKey     map[string]entity.PurchaseIntent
	RuntimeConfig          entity.RuntimeConfig
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := shopSnapshotState{
		ProductDefinitionsByID: s.productDefinitionsByID,
		OfferDefinitionsByID:   s.offerDefinitionsByID,
		PurchaseIntentsByID:    s.purchaseIntentsByID,
		PurchaseDedupByKey:     s.purchaseDedupByKey,
		RuntimeConfig:          s.runtimeConfig,
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

	var state shopSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.ProductDefinitionsByID == nil {
		state.ProductDefinitionsByID = make(map[string]entity.ProductDefinition)
	}
	if state.OfferDefinitionsByID == nil {
		state.OfferDefinitionsByID = make(map[string]entity.OfferDefinition)
	}
	if state.PurchaseIntentsByID == nil {
		state.PurchaseIntentsByID = make(map[string]entity.PurchaseIntent)
	}
	if state.PurchaseDedupByKey == nil {
		state.PurchaseDedupByKey = make(map[string]entity.PurchaseIntent)
	}
	if state.RuntimeConfig.UpdatedAt.IsZero() {
		state.RuntimeConfig = entity.RuntimeConfig{
			CatalogEnabled:  true,
			PurchaseEnabled: true,
			CampaignEnabled: true,
			UpdatedAt:       time.Now().UTC(),
		}
	}

	s.productDefinitionsByID = state.ProductDefinitionsByID
	s.offerDefinitionsByID = state.OfferDefinitionsByID
	s.purchaseIntentsByID = state.PurchaseIntentsByID
	s.purchaseDedupByKey = state.PurchaseDedupByKey
	s.runtimeConfig = state.RuntimeConfig

	return nil
}

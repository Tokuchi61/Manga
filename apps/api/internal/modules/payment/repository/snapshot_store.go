package repository

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
)

type paymentSnapshotState struct {
	ManaPackagesByID     map[string]entity.ManaPackage
	ProviderSessionsByID map[string]entity.ProviderSession
	TransactionsByID     map[string]entity.Transaction
	CheckoutDedupByKey   map[string]entity.Transaction
	CallbackDedupByEvent map[string]entity.Transaction
	LedgerEntriesByID    map[string]entity.LedgerEntry
	BalanceByUserID      map[string]entity.BalanceSnapshot
	RuntimeConfig        entity.RuntimeConfig
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := paymentSnapshotState{
		ManaPackagesByID:     s.manaPackagesByID,
		ProviderSessionsByID: s.providerSessionsByID,
		TransactionsByID:     s.transactionsByID,
		CheckoutDedupByKey:   s.checkoutDedupByKey,
		CallbackDedupByEvent: s.callbackDedupByEvent,
		LedgerEntriesByID:    s.ledgerEntriesByID,
		BalanceByUserID:      s.balanceByUserID,
		RuntimeConfig:        s.runtimeConfig,
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

	var state paymentSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.ManaPackagesByID == nil {
		state.ManaPackagesByID = make(map[string]entity.ManaPackage)
	}
	if state.ProviderSessionsByID == nil {
		state.ProviderSessionsByID = make(map[string]entity.ProviderSession)
	}
	if state.TransactionsByID == nil {
		state.TransactionsByID = make(map[string]entity.Transaction)
	}
	if state.CheckoutDedupByKey == nil {
		state.CheckoutDedupByKey = make(map[string]entity.Transaction)
	}
	if state.CallbackDedupByEvent == nil {
		state.CallbackDedupByEvent = make(map[string]entity.Transaction)
	}
	if state.LedgerEntriesByID == nil {
		state.LedgerEntriesByID = make(map[string]entity.LedgerEntry)
	}
	if state.BalanceByUserID == nil {
		state.BalanceByUserID = make(map[string]entity.BalanceSnapshot)
	}
	if state.RuntimeConfig.UpdatedAt.IsZero() {
		state.RuntimeConfig = entity.RuntimeConfig{
			ManaPurchaseEnabled:    true,
			CheckoutEnabled:        true,
			TransactionReadEnabled: true,
			CallbackIntakePaused:   false,
			UpdatedAt:              time.Now().UTC(),
		}
	}

	s.manaPackagesByID = state.ManaPackagesByID
	s.providerSessionsByID = state.ProviderSessionsByID
	s.transactionsByID = state.TransactionsByID
	s.checkoutDedupByKey = state.CheckoutDedupByKey
	s.callbackDedupByEvent = state.CallbackDedupByEvent
	s.ledgerEntriesByID = state.LedgerEntriesByID
	s.balanceByUserID = state.BalanceByUserID
	s.runtimeConfig = state.RuntimeConfig

	return nil
}

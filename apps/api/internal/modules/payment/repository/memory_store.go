package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
)

// MemoryStore is stage-19 bootstrap persistence for payment flows.
type MemoryStore struct {
	mu sync.RWMutex

	manaPackagesByID     map[string]entity.ManaPackage
	providerSessionsByID map[string]entity.ProviderSession
	transactionsByID     map[string]entity.Transaction
	checkoutDedupByKey   map[string]entity.Transaction
	callbackDedupByEvent map[string]entity.Transaction
	ledgerEntriesByID    map[string]entity.LedgerEntry
	balanceByUserID      map[string]entity.BalanceSnapshot

	runtimeConfig entity.RuntimeConfig
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	return &MemoryStore{
		manaPackagesByID:     make(map[string]entity.ManaPackage),
		providerSessionsByID: make(map[string]entity.ProviderSession),
		transactionsByID:     make(map[string]entity.Transaction),
		checkoutDedupByKey:   make(map[string]entity.Transaction),
		callbackDedupByEvent: make(map[string]entity.Transaction),
		ledgerEntriesByID:    make(map[string]entity.LedgerEntry),
		balanceByUserID:      make(map[string]entity.BalanceSnapshot),
		runtimeConfig: entity.RuntimeConfig{
			ManaPurchaseEnabled:    true,
			CheckoutEnabled:        true,
			TransactionReadEnabled: true,
			CallbackIntakePaused:   false,
			UpdatedAt:              now,
		},
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func cloneManaPackage(in entity.ManaPackage) entity.ManaPackage {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneProviderSession(in entity.ProviderSession) entity.ProviderSession {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	if in.ExpiresAt != nil {
		expiresAt := in.ExpiresAt.UTC()
		out.ExpiresAt = &expiresAt
	}
	return out
}

func cloneTransaction(in entity.Transaction) entity.Transaction {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneLedgerEntry(in entity.LedgerEntry) entity.LedgerEntry {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	return out
}

func cloneBalanceSnapshot(in entity.BalanceSnapshot) entity.BalanceSnapshot {
	out := in
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneRuntimeConfig(in entity.RuntimeConfig) entity.RuntimeConfig {
	out := in
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func applyOffsetLimit[T any](items []T, offset int, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []T{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]T(nil), items[offset:end]...)
}

func sortByUpdatedDesc[T any](items []T, lessFn func(i int, j int) bool) {
	sort.Slice(items, lessFn)
}

package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
)

var (
	ErrNotFound = errors.New("payment_repository_not_found")
)

// Store defines payment persistence boundary.
type Store interface {
	UpsertManaPackage(ctx context.Context, pkg entity.ManaPackage) error
	GetManaPackage(ctx context.Context, packageID string) (entity.ManaPackage, error)
	ListManaPackages(ctx context.Context, activeOnly bool, limit int, offset int) ([]entity.ManaPackage, error)

	CreateProviderSession(ctx context.Context, session entity.ProviderSession) error
	GetProviderSession(ctx context.Context, sessionID string) (entity.ProviderSession, error)
	UpdateProviderSession(ctx context.Context, session entity.ProviderSession) error

	CreateTransaction(ctx context.Context, tx entity.Transaction) error
	GetTransaction(ctx context.Context, transactionID string) (entity.Transaction, error)
	UpdateTransaction(ctx context.Context, tx entity.Transaction) error
	ListTransactionsByUser(ctx context.Context, userID string, status string, limit int, offset int) ([]entity.Transaction, error)

	GetCheckoutDedup(ctx context.Context, dedupKey string) (entity.Transaction, error)
	PutCheckoutDedup(ctx context.Context, dedupKey string, tx entity.Transaction) error
	GetCallbackDedup(ctx context.Context, providerEventID string) (entity.Transaction, error)
	PutCallbackDedup(ctx context.Context, providerEventID string, tx entity.Transaction) error

	CreateLedgerEntry(ctx context.Context, entry entity.LedgerEntry) error
	ListLedgerEntriesByUser(ctx context.Context, userID string) ([]entity.LedgerEntry, error)
	ListLedgerUsers(ctx context.Context) ([]string, error)

	GetBalanceSnapshot(ctx context.Context, userID string) (entity.BalanceSnapshot, error)
	UpsertBalanceSnapshot(ctx context.Context, snapshot entity.BalanceSnapshot) error

	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error
}

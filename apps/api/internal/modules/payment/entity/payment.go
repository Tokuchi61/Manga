package entity

import "time"

const (
	TransactionStatusPending   = "pending"
	TransactionStatusSuccess   = "success"
	TransactionStatusFailed    = "failed"
	TransactionStatusCancelled = "cancelled"
	TransactionStatusRefunded  = "refunded"
	TransactionStatusReversed  = "reversed"
)

const (
	ProviderSessionStatusStarted   = "started"
	ProviderSessionStatusCompleted = "completed"
	ProviderSessionStatusCancelled = "cancelled"
	ProviderSessionStatusExpired   = "expired"
)

const (
	LedgerEntryTypeCredit = "credit"
	LedgerEntryTypeDebit  = "debit"
)

const (
	CheckoutSourceExternalProvider = "external_provider"
	CheckoutSourceRecoveryReplay   = "recovery_replay"
)

// ManaPackage stores stage-19 payment package metadata.
type ManaPackage struct {
	PackageID     string
	Name          string
	Description   string
	ManaAmount    int
	PriceAmount   int
	PriceCurrency string
	Active        bool
	DisplayOrder  int
	Provider      string
	ProviderSKU   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ProviderSession stores checkout provider session metadata.
type ProviderSession struct {
	SessionID         string
	TransactionID     string
	UserID            string
	PackageID         string
	Provider          string
	ProviderReference string
	Status            string
	ExpiresAt         *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Transaction stores stage-19 payment transaction state.
type Transaction struct {
	TransactionID     string
	UserID            string
	PackageID         string
	AmountMana        int
	MoneyAmount       int
	MoneyCurrency     string
	Source            string
	Status            string
	Provider          string
	ProviderReference string
	RequestID         string
	CorrelationID     string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// LedgerEntry stores wallet impact entries.
type LedgerEntry struct {
	EntryID       string
	UserID        string
	TransactionID string
	EntryType     string
	AmountMana    int
	ReasonCode    string
	CreatedAt     time.Time
}

// BalanceSnapshot stores current wallet projection for a user.
type BalanceSnapshot struct {
	UserID            string
	BalanceMana       int
	LastLedgerEntryID string
	UpdatedAt         time.Time
}

// RuntimeConfig stores stage-19 runtime controls.
type RuntimeConfig struct {
	ManaPurchaseEnabled    bool
	CheckoutEnabled        bool
	TransactionReadEnabled bool
	CallbackIntakePaused   bool
	UpdatedAt              time.Time
}

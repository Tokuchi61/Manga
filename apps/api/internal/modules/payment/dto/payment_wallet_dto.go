package dto

import "time"

// GetOwnWalletRequest resolves actor wallet snapshot.
type GetOwnWalletRequest struct {
	ActorUserID string `json:"-" validate:"required,max=64"`
}

// WalletResponse returns wallet projection payload.
type WalletResponse struct {
	UserID            string    `json:"user_id"`
	BalanceMana       int       `json:"balance_mana"`
	LastLedgerEntryID string    `json:"last_ledger_entry_id,omitempty"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ListOwnTransactionsRequest resolves actor transaction list.
type ListOwnTransactionsRequest struct {
	ActorUserID string `json:"-" validate:"required,max=64"`
	Status      string `json:"-" validate:"omitempty,oneof=pending success failed cancelled refunded reversed"`
	Limit       int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset      int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// TransactionResponse returns transaction payload.
type TransactionResponse struct {
	TransactionID     string    `json:"transaction_id"`
	PackageID         string    `json:"package_id"`
	AmountMana        int       `json:"amount_mana"`
	MoneyAmount       int       `json:"money_amount"`
	MoneyCurrency     string    `json:"money_currency"`
	Source            string    `json:"source"`
	Status            string    `json:"status"`
	Provider          string    `json:"provider"`
	ProviderReference string    `json:"provider_reference,omitempty"`
	RequestID         string    `json:"request_id"`
	CorrelationID     string    `json:"correlation_id,omitempty"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ListOwnTransactionsResponse wraps transaction list payload.
type ListOwnTransactionsResponse struct {
	Items []TransactionResponse `json:"items"`
	Count int                   `json:"count"`
}

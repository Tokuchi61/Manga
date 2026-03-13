package dto

import "time"

// UpdateManaPurchaseStateRequest updates mana-purchase runtime surface.
type UpdateManaPurchaseStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateCheckoutStateRequest updates checkout runtime surface.
type UpdateCheckoutStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateTransactionReadStateRequest updates transaction read runtime surface.
type UpdateTransactionReadStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateCallbackIntakeStateRequest updates callback intake pause state.
type UpdateCallbackIntakeStateRequest struct {
	Paused bool `json:"paused"`
}

// RuntimeConfigResponse is payment runtime control payload.
type RuntimeConfigResponse struct {
	ManaPurchaseEnabled    bool      `json:"mana_purchase_enabled"`
	CheckoutEnabled        bool      `json:"checkout_enabled"`
	TransactionReadEnabled bool      `json:"transaction_read_enabled"`
	CallbackIntakePaused   bool      `json:"callback_intake_paused"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// UpsertManaPackageRequest creates or updates package definition.
type UpsertManaPackageRequest struct {
	PackageID     string `json:"package_id" validate:"required,max=64"`
	Name          string `json:"name" validate:"required,max=160"`
	Description   string `json:"description,omitempty" validate:"omitempty,max=400"`
	ManaAmount    int    `json:"mana_amount" validate:"required,min=1,max=1000000"`
	PriceAmount   int    `json:"price_amount" validate:"required,min=1,max=100000000"`
	PriceCurrency string `json:"price_currency" validate:"required,max=8"`
	Active        bool   `json:"active"`
	DisplayOrder  int    `json:"display_order" validate:"omitempty,min=0,max=100000"`
	Provider      string `json:"provider" validate:"required,max=64"`
	ProviderSKU   string `json:"provider_sku" validate:"required,max=128"`
}

// ListAdminManaPackagesRequest resolves admin package list.
type ListAdminManaPackagesRequest struct {
	ActiveOnly bool `json:"-"`
	Limit      int  `json:"-" validate:"omitempty,min=1,max=100"`
	Offset     int  `json:"-" validate:"omitempty,min=0,max=10000"`
}

// RunReconcileRequest runs reconcile for one actor or all users.
type RunReconcileRequest struct {
	ActorUserID string `json:"actor_user_id,omitempty" validate:"omitempty,max=64"`
}

// RunReconcileResponse returns reconcile summary.
type RunReconcileResponse struct {
	ScannedUsers   int `json:"scanned_users"`
	CorrectedUsers int `json:"corrected_users"`
}

// ProcessRefundRequest triggers refund flow.
type ProcessRefundRequest struct {
	TransactionID string `json:"transaction_id" validate:"required,max=64"`
	ReasonCode    string `json:"reason_code" validate:"required,max=128"`
}

// ProcessReversalRequest triggers reversal flow.
type ProcessReversalRequest struct {
	TransactionID string `json:"transaction_id" validate:"required,max=64"`
	ReasonCode    string `json:"reason_code" validate:"required,max=128"`
}

// AdminTransactionActionResponse returns refund/reversal action result.
type AdminTransactionActionResponse struct {
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
	BalanceMana   int       `json:"balance_mana"`
	UpdatedAt     time.Time `json:"updated_at"`
}

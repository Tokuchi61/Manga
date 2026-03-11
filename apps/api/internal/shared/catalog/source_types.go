package catalog

// PurchaseSourceType defines canonical purchase/checkout initiators.
type PurchaseSourceType string

const (
	PurchaseSourceCatalogPurchase   PurchaseSourceType = "catalog_purchase"
	PurchaseSourcePremiumActivation PurchaseSourceType = "premium_activation"
	PurchaseSourceManaWallet        PurchaseSourceType = "mana_wallet"
	PurchaseSourceExternalProvider  PurchaseSourceType = "external_provider"
	PurchaseSourceRecoveryReplay    PurchaseSourceType = "recovery_replay"
	PurchaseSourceAdminIssue        PurchaseSourceType = "admin_issue"
	PurchaseSourceGiftCode          PurchaseSourceType = "gift_code"
)

var AllPurchaseSourceTypes = []PurchaseSourceType{
	PurchaseSourceCatalogPurchase,
	PurchaseSourcePremiumActivation,
	PurchaseSourceManaWallet,
	PurchaseSourceExternalProvider,
	PurchaseSourceRecoveryReplay,
	PurchaseSourceAdminIssue,
	PurchaseSourceGiftCode,
}

func IsValidPurchaseSourceType(value PurchaseSourceType) bool {
	switch value {
	case PurchaseSourceCatalogPurchase,
		PurchaseSourcePremiumActivation,
		PurchaseSourceManaWallet,
		PurchaseSourceExternalProvider,
		PurchaseSourceRecoveryReplay,
		PurchaseSourceAdminIssue,
		PurchaseSourceGiftCode:
		return true
	default:
		return false
	}
}

// RewardSourceType defines canonical reward/grant initiators.
type RewardSourceType string

const (
	RewardSourceMission              RewardSourceType = "mission"
	RewardSourceRoyalPass            RewardSourceType = "royalpass"
	RewardSourceShop                 RewardSourceType = "shop"
	RewardSourceAdminGrant           RewardSourceType = "admin_grant"
	RewardSourceCompensation         RewardSourceType = "compensation"
	RewardSourceSeasonalEvent        RewardSourceType = "seasonal_event"
	RewardSourceReferral             RewardSourceType = "referral"
	RewardSourceReconciliationRepair RewardSourceType = "reconciliation_repair"
)

var AllRewardSourceTypes = []RewardSourceType{
	RewardSourceMission,
	RewardSourceRoyalPass,
	RewardSourceShop,
	RewardSourceAdminGrant,
	RewardSourceCompensation,
	RewardSourceSeasonalEvent,
	RewardSourceReferral,
	RewardSourceReconciliationRepair,
}

func IsValidRewardSourceType(value RewardSourceType) bool {
	switch value {
	case RewardSourceMission,
		RewardSourceRoyalPass,
		RewardSourceShop,
		RewardSourceAdminGrant,
		RewardSourceCompensation,
		RewardSourceSeasonalEvent,
		RewardSourceReferral,
		RewardSourceReconciliationRepair:
		return true
	default:
		return false
	}
}

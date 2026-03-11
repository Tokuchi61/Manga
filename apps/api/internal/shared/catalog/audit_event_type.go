package catalog

// AuditEventType defines canonical audit action categories.
type AuditEventType string

const (
	AuditEventSecurityAuth     AuditEventType = "security_auth"
	AuditEventAccessPolicy     AuditEventType = "access_policy"
	AuditEventAdminAction      AuditEventType = "admin_action"
	AuditEventModeration       AuditEventType = "moderation_action"
	AuditEventSupportCase      AuditEventType = "support_case"
	AuditEventPaymentFinancial AuditEventType = "payment_financial"
	AuditEventInventoryChange  AuditEventType = "inventory_change"
	AuditEventShopPurchase     AuditEventType = "shop_purchase"
	AuditEventUserState        AuditEventType = "user_state"
	AuditEventNotificationOps  AuditEventType = "notification_ops"
	AuditEventAdsOps           AuditEventType = "ads_ops"
	AuditEventSystemOps        AuditEventType = "system_ops"
)

var AllAuditEventTypes = []AuditEventType{
	AuditEventSecurityAuth,
	AuditEventAccessPolicy,
	AuditEventAdminAction,
	AuditEventModeration,
	AuditEventSupportCase,
	AuditEventPaymentFinancial,
	AuditEventInventoryChange,
	AuditEventShopPurchase,
	AuditEventUserState,
	AuditEventNotificationOps,
	AuditEventAdsOps,
	AuditEventSystemOps,
}

func IsValidAuditEventType(value AuditEventType) bool {
	switch value {
	case AuditEventSecurityAuth,
		AuditEventAccessPolicy,
		AuditEventAdminAction,
		AuditEventModeration,
		AuditEventSupportCase,
		AuditEventPaymentFinancial,
		AuditEventInventoryChange,
		AuditEventShopPurchase,
		AuditEventUserState,
		AuditEventNotificationOps,
		AuditEventAdsOps,
		AuditEventSystemOps:
		return true
	default:
		return false
	}
}

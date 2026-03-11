package catalog

// NotificationCategory defines canonical user-facing notification classes.
type NotificationCategory string

const (
	NotificationAccountSecurity NotificationCategory = "account_security"
	NotificationSocial          NotificationCategory = "social"
	NotificationComment         NotificationCategory = "comment"
	NotificationSupport         NotificationCategory = "support"
	NotificationModeration      NotificationCategory = "moderation"
	NotificationMission         NotificationCategory = "mission"
	NotificationRoyalPass       NotificationCategory = "royalpass"
	NotificationShop            NotificationCategory = "shop"
	NotificationPayment         NotificationCategory = "payment"
	NotificationSystemOps       NotificationCategory = "system_ops"
)

var AllNotificationCategories = []NotificationCategory{
	NotificationAccountSecurity,
	NotificationSocial,
	NotificationComment,
	NotificationSupport,
	NotificationModeration,
	NotificationMission,
	NotificationRoyalPass,
	NotificationShop,
	NotificationPayment,
	NotificationSystemOps,
}

func IsValidNotificationCategory(value NotificationCategory) bool {
	switch value {
	case NotificationAccountSecurity,
		NotificationSocial,
		NotificationComment,
		NotificationSupport,
		NotificationModeration,
		NotificationMission,
		NotificationRoyalPass,
		NotificationShop,
		NotificationPayment,
		NotificationSystemOps:
		return true
	default:
		return false
	}
}

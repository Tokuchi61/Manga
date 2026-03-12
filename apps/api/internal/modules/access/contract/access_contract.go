package contract

import (
	authcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/contract"
	usercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/contract"
)

const (
	PermissionSiteView                  = "site.view"
	PermissionMangaDetailView           = "manga.detail.view"
	PermissionMangaDiscoveryView        = "manga.discovery.view"
	PermissionChapterReadAuthenticated  = "chapter.read.authenticated"
	PermissionChapterEarlyAccessVIP     = "chapter.early_access.vip"
	PermissionCommentWriteAuthenticated = "comment.write.authenticated"
	PermissionHistoryContinueReadingOwn = "history.continue_reading.read.own"
	PermissionHistoryTimelineReadOwn    = "history.timeline.read.own"
	PermissionHistoryLibraryReadOwn     = "history.library.read.own"
	PermissionHistoryLibraryReadPublic  = "history.library.read.public"
	PermissionHistoryBookmarkWriteOwn   = "history.bookmark.write.own"
	PermissionAdsView                   = "ads.view"
	PermissionShopItemPurchase          = "shop.item.purchase"
	PermissionPaymentManaPurchase       = "payment.mana.purchase"
	PermissionPaymentTransactionReadOwn = "payment.transaction.read.own"
	PermissionModerationPanelView       = "moderation.panel.view"
	PermissionModerationActionAny       = "moderation.action.any"
)

var CanonicalPermissions = []string{
	PermissionSiteView,
	PermissionMangaDetailView,
	PermissionMangaDiscoveryView,
	PermissionChapterReadAuthenticated,
	PermissionChapterEarlyAccessVIP,
	PermissionCommentWriteAuthenticated,
	PermissionHistoryContinueReadingOwn,
	PermissionHistoryTimelineReadOwn,
	PermissionHistoryLibraryReadOwn,
	PermissionHistoryLibraryReadPublic,
	PermissionHistoryBookmarkWriteOwn,
	PermissionAdsView,
	PermissionShopItemPurchase,
	PermissionPaymentManaPurchase,
	PermissionPaymentTransactionReadOwn,
	PermissionModerationPanelView,
	PermissionModerationActionAny,
}

// AuthorizationInput combines auth and user owned signals for access decisions.
type AuthorizationInput struct {
	UserID              string
	Permission          string
	FeatureKey          string
	ResourceOwnerUserID string
	Identity            authcontract.VerifiedIdentity
	UserSignal          usercontract.AccessSignal
}

// AuthorizationDecision is the stable access output contract.
type AuthorizationDecision struct {
	Allowed       bool
	Effect        string
	ReasonCode    string
	PolicyVersion int
}

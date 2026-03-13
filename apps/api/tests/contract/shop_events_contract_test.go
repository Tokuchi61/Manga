package contract_test

import (
	"testing"

	shopevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/events"
	"github.com/stretchr/testify/require"
)

func TestShopEventConstants(t *testing.T) {
	require.Equal(t, "shop.purchase.intent.created", shopevents.EventShopPurchaseIntentCreated)
	require.Equal(t, "shop.purchase.completed", shopevents.EventShopPurchaseCompleted)
	require.Equal(t, "shop.purchase.recovery_requested", shopevents.EventShopPurchaseRecoveryRequested)
}

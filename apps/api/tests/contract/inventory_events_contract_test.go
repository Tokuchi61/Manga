package contract_test

import (
	"testing"

	inventoryevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/events"
	"github.com/stretchr/testify/require"
)

func TestInventoryEventConstants(t *testing.T) {
	require.Equal(t, "inventory.granted", inventoryevents.EventInventoryGranted)
	require.Equal(t, "inventory.revoked", inventoryevents.EventInventoryRevoked)
	require.Equal(t, "inventory.consumed", inventoryevents.EventInventoryConsumed)
	require.Equal(t, "inventory.equipped", inventoryevents.EventInventoryEquipped)
}

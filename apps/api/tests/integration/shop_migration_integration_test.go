package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShopMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603130016_shop_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603130016_shop_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS shop_products")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS shop_offers")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS shop_purchase_intents")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS shop_purchase_dedup")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS shop_runtime_controls")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS shop_runtime_controls")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS shop_purchase_dedup")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS shop_purchase_intents")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS shop_products")
}

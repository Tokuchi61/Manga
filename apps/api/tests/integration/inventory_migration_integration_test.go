package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInventoryMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120013_inventory_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120013_inventory_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS inventory_item_definitions")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS inventory_entries")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS inventory_grant_dedup")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS inventory_runtime_controls")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS inventory_runtime_controls")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS inventory_grant_dedup")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS inventory_entries")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS inventory_item_definitions")
}

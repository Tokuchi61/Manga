package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603130017_payment_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603130017_payment_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS payment_mana_packages")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS payment_transactions")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS payment_provider_sessions")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS payment_ledger_entries")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS payment_balance_snapshots")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS payment_runtime_controls")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS payment_runtime_controls")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS payment_provider_sessions")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS payment_transactions")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS payment_mana_packages")
}

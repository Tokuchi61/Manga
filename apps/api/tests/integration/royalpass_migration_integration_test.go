package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoyalPassMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120015_royalpass_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120015_royalpass_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS royalpass_seasons")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS royalpass_tiers")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS royalpass_user_progress")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS royalpass_progress_dedup")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS royalpass_claim_dedup")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS royalpass_runtime_controls")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS royalpass_runtime_controls")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS royalpass_premium_activation_dedup")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS royalpass_user_progress")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS royalpass_seasons")
}

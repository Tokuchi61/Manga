package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdminMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603130019_admin_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603130019_admin_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS admin_runtime_controls")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS admin_actions")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS admin_action_dedup")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS admin_overrides")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS admin_user_reviews")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS admin_impersonation_sessions")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS admin_impersonation_sessions")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS admin_user_reviews")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS admin_overrides")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS admin_action_dedup")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS admin_actions")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS admin_runtime_controls")
}

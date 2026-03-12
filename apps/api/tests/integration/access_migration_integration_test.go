package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccessMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120004_access_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120004_access_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS access_roles")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS access_permissions")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS access_role_permissions")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS access_user_roles")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS access_policy_rules")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS access_temporary_grants")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS access_temporary_grants")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS access_policy_rules")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS access_user_roles")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS access_role_permissions")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS access_permissions")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS access_roles")
}

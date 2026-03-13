package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120002_auth_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120002_auth_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS auth_credentials")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS auth_sessions")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS auth_tokens")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS auth_security_events")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS auth_security_events")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS auth_tokens")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS auth_sessions")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS auth_credentials")
}

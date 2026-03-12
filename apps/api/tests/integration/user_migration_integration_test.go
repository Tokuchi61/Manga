package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserMigrationFilesExistAndContainCoreTable(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120003_user_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120003_user_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS user_accounts")
	require.Contains(t, upSQL, "profile_visibility")
	require.Contains(t, upSQL, "history_visibility_preference")
	require.Contains(t, upSQL, "account_state")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS user_accounts")
}

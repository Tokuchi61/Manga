package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSupportMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120008_support_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120008_support_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS support_entries")
	require.Contains(t, upSQL, "support_kind")
	require.Contains(t, upSQL, "requester_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE")
	require.Contains(t, upSQL, "target_type")
	require.Contains(t, upSQL, "target_id UUID NULL")
	require.Contains(t, upSQL, "duplicate_of_support_id UUID NULL REFERENCES support_entries(id)")
	require.Contains(t, upSQL, "request_id")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS support_replies")
	require.Contains(t, upSQL, "visibility")
	require.Contains(t, upSQL, "public_to_requester")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS support_replies")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS support_entries")
}

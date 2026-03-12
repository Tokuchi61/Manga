package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModerationMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120009_moderation_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120009_moderation_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS moderation_cases")
	require.Contains(t, upSQL, "case_status")
	require.Contains(t, upSQL, "assignment_status")
	require.Contains(t, upSQL, "escalation_status")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS moderation_case_notes")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS moderation_case_actions")
	require.Contains(t, upSQL, "moderation_case_actions_action_type_check")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS moderation_case_actions")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS moderation_case_notes")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS moderation_cases")
}

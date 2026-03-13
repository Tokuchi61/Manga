package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHistoryMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120011_history_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120011_history_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS history_library_entries")
	require.Contains(t, upSQL, "history_library_entries_reading_status_check")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS history_timeline_events")
	require.Contains(t, upSQL, "history_timeline_events_event_check")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS history_runtime_controls")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS history_runtime_controls")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS history_timeline_events")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS history_library_entries")
}

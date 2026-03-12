package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotificationMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120010_notification_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120010_notification_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS notification_entries")
	require.Contains(t, upSQL, "notification_entries_category_check")
	require.Contains(t, upSQL, "notification_entries_channel_check")
	require.Contains(t, upSQL, "notification_entries_state_check")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS notification_delivery_attempts")
	require.Contains(t, upSQL, "notification_delivery_attempts_result_check")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS notification_preferences")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS notification_runtime_controls")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS notification_runtime_controls")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS notification_preferences")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS notification_delivery_attempts")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS notification_entries")
}

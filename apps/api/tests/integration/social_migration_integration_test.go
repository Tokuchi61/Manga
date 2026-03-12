package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSocialMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120012_social_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120012_social_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS social_friendship_requests")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS social_follows")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS social_wall_posts")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS social_message_threads")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS social_messages")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS social_runtime_controls")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS social_runtime_controls")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS social_messages")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS social_wall_posts")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS social_friendship_requests")
}

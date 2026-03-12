package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMissionMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120014_mission_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120014_mission_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS mission_definitions")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS mission_user_progress")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS mission_progress_dedup")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS mission_claim_dedup")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS mission_runtime_controls")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS mission_runtime_controls")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS mission_claim_dedup")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS mission_user_progress")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS mission_definitions")
}

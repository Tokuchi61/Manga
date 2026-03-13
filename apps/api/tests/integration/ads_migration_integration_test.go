package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdsMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603130018_ads_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603130018_ads_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS ads_placements")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS ads_campaigns")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS ads_impressions")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS ads_clicks")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS ads_campaign_aggregates")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS ads_runtime_controls")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS ads_runtime_controls")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS ads_clicks")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS ads_impressions")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS ads_campaigns")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS ads_placements")
}

package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChapterMigrationFilesExistAndContainCoreTables(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120006_chapter_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120006_chapter_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS chapter_entries")
	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS chapter_pages")
	require.Contains(t, upSQL, "manga_id UUID NOT NULL REFERENCES manga_entries(id) ON DELETE CASCADE")
	require.Contains(t, upSQL, "UNIQUE (manga_id, slug)")
	require.Contains(t, upSQL, "UNIQUE (manga_id, sequence_no)")
	require.Contains(t, upSQL, "read_access_level")
	require.Contains(t, upSQL, "early_access_level")
	require.Contains(t, upSQL, "media_health_status")
	require.Contains(t, upSQL, "integrity_status")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS chapter_pages")
	require.Contains(t, downSQL, "DROP TABLE IF EXISTS chapter_entries")
}

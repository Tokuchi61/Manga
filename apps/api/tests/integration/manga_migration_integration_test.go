package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMangaMigrationFilesExistAndContainCoreTable(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120005_manga_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120005_manga_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS manga_entries")
	require.Contains(t, upSQL, "slug TEXT NOT NULL UNIQUE")
	require.Contains(t, upSQL, "publish_state")
	require.Contains(t, upSQL, "visibility")
	require.Contains(t, upSQL, "chapter_count")
	require.Contains(t, upSQL, "comment_count")
	require.Contains(t, upSQL, "view_count")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS manga_entries")
}

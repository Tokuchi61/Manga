package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommentMigrationFilesExistAndContainCoreTable(t *testing.T) {
	upPath := filepath.Join("..", "..", "migrations", "202603120007_comment_create_core_tables.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "202603120007_comment_create_core_tables.down.sql")

	upBytes, err := os.ReadFile(upPath)
	require.NoError(t, err)
	downBytes, err := os.ReadFile(downPath)
	require.NoError(t, err)

	upSQL := string(upBytes)
	downSQL := string(downBytes)

	require.Contains(t, upSQL, "CREATE TABLE IF NOT EXISTS comment_entries")
	require.Contains(t, upSQL, "target_type")
	require.Contains(t, upSQL, "target_id UUID NOT NULL")
	require.Contains(t, upSQL, "author_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE")
	require.Contains(t, upSQL, "parent_comment_id UUID NULL REFERENCES comment_entries(id)")
	require.Contains(t, upSQL, "root_comment_id UUID NULL REFERENCES comment_entries(id)")
	require.Contains(t, upSQL, "moderation_status")
	require.Contains(t, upSQL, "shadowbanned")
	require.Contains(t, upSQL, "deleted_at")

	require.Contains(t, downSQL, "DROP TABLE IF EXISTS comment_entries")
}

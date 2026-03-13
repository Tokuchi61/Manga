package integration_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestMigrationsApplyAndRollbackWithDBTestDSN(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("DB_TEST_DSN"))
	if dsn == "" {
		t.Skip("DB_TEST_DSN is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	schemaName := fmt.Sprintf("migration_test_%d", time.Now().UTC().UnixNano())
	_, err = pool.Exec(ctx, fmt.Sprintf("CREATE SCHEMA %s", schemaName))
	require.NoError(t, err)
	defer func() {
		_, _ = pool.Exec(context.Background(), fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName))
	}()

	_, err = pool.Exec(ctx, fmt.Sprintf("SET search_path TO %s", schemaName))
	require.NoError(t, err)

	upFiles, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.up.sql"))
	require.NoError(t, err)
	require.NotEmpty(t, upFiles)
	sort.Strings(upFiles)

	for _, file := range upFiles {
		payload, readErr := os.ReadFile(file)
		require.NoErrorf(t, readErr, "read migration up file failed: %s", file)
		_, execErr := pool.Exec(ctx, string(payload))
		require.NoErrorf(t, execErr, "apply migration up failed: %s", file)
	}

	var tableCount int
	err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = $1`, schemaName).Scan(&tableCount)
	require.NoError(t, err)
	require.Greater(t, tableCount, 0)

	downFiles, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.down.sql"))
	require.NoError(t, err)
	require.NotEmpty(t, downFiles)
	sort.Strings(downFiles)

	for index := len(downFiles) - 1; index >= 0; index-- {
		file := downFiles[index]
		payload, readErr := os.ReadFile(file)
		require.NoErrorf(t, readErr, "read migration down file failed: %s", file)
		_, execErr := pool.Exec(ctx, string(payload))
		require.NoErrorf(t, execErr, "apply migration down failed: %s", file)
	}

	err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = $1`, schemaName).Scan(&tableCount)
	require.NoError(t, err)
	require.Equal(t, 0, tableCount)
}

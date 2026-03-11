package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadSuccess(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("APP_VERSION", "0.1.0-test")
	t.Setenv("DB_MAIN_DSN", "postgres://main:main@localhost:5432/main?sslmode=disable")
	t.Setenv("DB_TEST_DSN", "postgres://test:test@localhost:5432/test?sslmode=disable")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, "test", cfg.Env)
	require.Equal(t, 9090, cfg.Port)
	require.Equal(t, "0.1.0-test", cfg.AppVersion)
}

func TestLoadFailsWithoutMainDSN(t *testing.T) {
	t.Setenv("APP_VERSION", "0.1.0-test")
	t.Setenv("DB_MAIN_DSN", "")
	t.Setenv("DB_TEST_DSN", "postgres://test:test@localhost:5432/test?sslmode=disable")

	_, err := Load()
	require.Error(t, err)
	require.Contains(t, err.Error(), "DB_MAIN_DSN")
}

func TestLoadFailsWithoutAppVersion(t *testing.T) {
	t.Setenv("APP_VERSION", "")
	t.Setenv("DB_MAIN_DSN", "postgres://main:main@localhost:5432/main?sslmode=disable")
	t.Setenv("DB_TEST_DSN", "postgres://test:test@localhost:5432/test?sslmode=disable")

	_, err := Load()
	require.Error(t, err)
	require.Contains(t, err.Error(), "APP_VERSION")
}

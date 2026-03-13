package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadSuccess(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("APP_VERSION", "0.1.0-test")
	t.Setenv("APP_ALLOW_MEMORY_MODE", "false")
	t.Setenv("DB_MAIN_DSN", "postgres://main:main@localhost:5432/main?sslmode=disable")
	t.Setenv("DB_TEST_DSN", "postgres://test:test@localhost:5432/test?sslmode=disable")
	t.Setenv("STATE_SNAPSHOT_DIR", ".cache/test-state")
	t.Setenv("STATE_SNAPSHOT_INTERVAL", "45s")
	t.Setenv("STATE_SNAPSHOT_WRITE_THROUGH", "true")
	t.Setenv("AUTH_ACCESS_TOKEN_SECRET", "test-secret")
	t.Setenv("AUTH_EXPOSE_SENSITIVE_TOKENS", "false")
	t.Setenv("AUTH_LOGIN_FAILED_ATTEMPT_LIMIT_PER_MINUTE", "7")
	t.Setenv("AUTH_LOGIN_COOLDOWN_SECONDS", "420")
	t.Setenv("AUTH_EMAIL_VERIFICATION_RESEND_COOLDOWN_SECONDS", "90")
	t.Setenv("PAYMENT_CALLBACK_SIGNING_SECRET", "payment-secret")
	t.Setenv("PAYMENT_CALLBACK_TIMESTAMP_SKEW", "6m")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, "test", cfg.Env)
	require.Equal(t, 9090, cfg.Port)
	require.Equal(t, "0.1.0-test", cfg.AppVersion)
	require.Equal(t, ".cache/test-state", cfg.StateSnapshotDir)
	require.Equal(t, 45*time.Second, cfg.StateSnapshotInterval)
	require.True(t, cfg.StateSnapshotWriteThrough)
	require.Equal(t, 7, cfg.AuthLoginFailedAttemptLimitPerMinute)
	require.Equal(t, 420, cfg.AuthLoginCooldownSeconds)
	require.Equal(t, 90, cfg.AuthEmailVerificationResendCooldownSeconds)
	require.Equal(t, "test-secret", cfg.AuthAccessTokenSecret)
	require.Equal(t, 6*time.Minute, cfg.PaymentCallbackTimestampSkew)
	require.False(t, cfg.AuthExposeSensitiveTokens)
}

func TestLoadAllowsMissingMainDSNWhenMemoryModeEnabled(t *testing.T) {
	t.Setenv("APP_VERSION", "0.1.0-test")
	t.Setenv("APP_ALLOW_MEMORY_MODE", "true")
	t.Setenv("DB_MAIN_DSN", "")
	t.Setenv("DB_TEST_DSN", "")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, "", cfg.DBMainDSN)
	require.True(t, cfg.AllowMemoryMode)
}

func TestLoadFailsWhenMissingMainDSNAndMemoryModeDisabled(t *testing.T) {
	t.Setenv("APP_VERSION", "0.1.0-test")
	t.Setenv("APP_ALLOW_MEMORY_MODE", "false")
	t.Setenv("DB_MAIN_DSN", "")

	_, err := Load()
	require.Error(t, err)
	require.Contains(t, err.Error(), "DB_MAIN_DSN")
}

func TestLoadFailsWithoutAppVersion(t *testing.T) {
	t.Setenv("APP_VERSION", "")
	t.Setenv("DB_MAIN_DSN", "postgres://main:main@localhost:5432/main?sslmode=disable")
	t.Setenv("DB_TEST_DSN", "postgres://test:test@localhost:5432/test?sslmode=disable")
	t.Setenv("AUTH_ACCESS_TOKEN_SECRET", "test-secret")
	t.Setenv("AUTH_LOGIN_FAILED_ATTEMPT_LIMIT_PER_MINUTE", "7")
	t.Setenv("AUTH_LOGIN_COOLDOWN_SECONDS", "420")
	t.Setenv("AUTH_EMAIL_VERIFICATION_RESEND_COOLDOWN_SECONDS", "90")
	t.Setenv("PAYMENT_CALLBACK_SIGNING_SECRET", "payment-secret")

	_, err := Load()
	require.Error(t, err)
	require.Contains(t, err.Error(), "APP_VERSION")
}

func TestLoadFailsWithNegativeSnapshotInterval(t *testing.T) {
	t.Setenv("APP_VERSION", "0.1.0-test")
	t.Setenv("APP_ALLOW_MEMORY_MODE", "true")
	t.Setenv("STATE_SNAPSHOT_INTERVAL", "-1s")

	_, err := Load()
	require.Error(t, err)
	require.Contains(t, err.Error(), "STATE_SNAPSHOT_INTERVAL")
}

func TestLoadFailsWhenSensitiveTokenExposureEnabledInProd(t *testing.T) {
	t.Setenv("APP_ENV", "prod")
	t.Setenv("APP_VERSION", "0.1.0-test")
	t.Setenv("DB_MAIN_DSN", "postgres://main:main@localhost:5432/main?sslmode=disable")
	t.Setenv("AUTH_ACCESS_TOKEN_SECRET", "test-secret")
	t.Setenv("PAYMENT_CALLBACK_SIGNING_SECRET", "payment-secret")
	t.Setenv("AUTH_EXPOSE_SENSITIVE_TOKENS", "true")

	_, err := Load()
	require.Error(t, err)
	require.Contains(t, err.Error(), "AUTH_EXPOSE_SENSITIVE_TOKENS")
}

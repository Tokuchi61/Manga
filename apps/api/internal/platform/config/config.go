package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env                                        string        `env:"APP_ENV" envDefault:"local"`
	Port                                       int           `env:"APP_PORT" envDefault:"8080"`
	AppVersion                                 string        `env:"APP_VERSION,required"`
	AllowMemoryMode                            bool          `env:"APP_ALLOW_MEMORY_MODE" envDefault:"false"`
	DBMainDSN                                  string        `env:"DB_MAIN_DSN"`
	DBTestDSN                                  string        `env:"DB_TEST_DSN"`
	DBMaxConns                                 int32         `env:"DB_MAX_CONNS" envDefault:"10"`
	DBConnectTimeout                           time.Duration `env:"DB_CONNECT_TIMEOUT" envDefault:"5s"`
	HTTPShutdownTimeout                        time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
	StateSnapshotDir                           string        `env:"STATE_SNAPSHOT_DIR" envDefault:".cache/state"`
	StateSnapshotInterval                      time.Duration `env:"STATE_SNAPSHOT_INTERVAL" envDefault:"30s"`
	StateSnapshotWriteThrough                  bool          `env:"STATE_SNAPSHOT_WRITE_THROUGH" envDefault:"true"`
	AuthAccessTokenSecret                      string        `env:"AUTH_ACCESS_TOKEN_SECRET"`
	AuthAccessTokenIssuer                      string        `env:"AUTH_ACCESS_TOKEN_ISSUER" envDefault:"novascans-api"`
	AuthExposeSensitiveTokens                  bool          `env:"AUTH_EXPOSE_SENSITIVE_TOKENS" envDefault:"false"`
	AuthLoginFailedAttemptLimitPerMinute       int           `env:"AUTH_LOGIN_FAILED_ATTEMPT_LIMIT_PER_MINUTE" envDefault:"5"`
	AuthLoginCooldownSeconds                   int           `env:"AUTH_LOGIN_COOLDOWN_SECONDS" envDefault:"300"`
	AuthEmailVerificationResendCooldownSeconds int           `env:"AUTH_EMAIL_VERIFICATION_RESEND_COOLDOWN_SECONDS" envDefault:"60"`
	PaymentCallbackSigningSecret               string        `env:"PAYMENT_CALLBACK_SIGNING_SECRET"`
	PaymentCallbackTimestampSkew               time.Duration `env:"PAYMENT_CALLBACK_TIMESTAMP_SKEW" envDefault:"5m"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	cfg.Env = normalize(cfg.Env)
	if strings.TrimSpace(cfg.AppVersion) == "" {
		return Config{}, fmt.Errorf("APP_VERSION cannot be empty")
	}
	if cfg.StateSnapshotInterval < 0 {
		return Config{}, fmt.Errorf("STATE_SNAPSHOT_INTERVAL cannot be negative")
	}
	if cfg.PaymentCallbackTimestampSkew <= 0 {
		return Config{}, fmt.Errorf("PAYMENT_CALLBACK_TIMESTAMP_SKEW must be positive")
	}

	hasMainDB := strings.TrimSpace(cfg.DBMainDSN) != ""
	if !hasMainDB && !cfg.AllowMemoryMode {
		return Config{}, fmt.Errorf("DB_MAIN_DSN is required unless APP_ALLOW_MEMORY_MODE=true")
	}

	if isProdEnv(cfg.Env) {
		if !hasMainDB {
			return Config{}, fmt.Errorf("DB_MAIN_DSN is required in production")
		}
		if strings.TrimSpace(cfg.AuthAccessTokenSecret) == "" {
			return Config{}, fmt.Errorf("AUTH_ACCESS_TOKEN_SECRET is required in production")
		}
		if strings.TrimSpace(cfg.PaymentCallbackSigningSecret) == "" {
			return Config{}, fmt.Errorf("PAYMENT_CALLBACK_SIGNING_SECRET is required in production")
		}
		if cfg.AuthExposeSensitiveTokens {
			return Config{}, fmt.Errorf("AUTH_EXPOSE_SENSITIVE_TOKENS cannot be enabled in production")
		}
	}

	return cfg, nil
}

func (c Config) Addr() string {
	return ":" + strconv.Itoa(c.Port)
}

func isProdEnv(env string) bool {
	env = normalize(env)
	return env == "prod" || env == "production"
}

func normalize(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

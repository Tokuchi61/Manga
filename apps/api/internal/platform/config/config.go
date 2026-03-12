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
	DBMainDSN                                  string        `env:"DB_MAIN_DSN,required"`
	DBTestDSN                                  string        `env:"DB_TEST_DSN,required"`
	DBMaxConns                                 int32         `env:"DB_MAX_CONNS" envDefault:"10"`
	DBConnectTimeout                           time.Duration `env:"DB_CONNECT_TIMEOUT" envDefault:"5s"`
	HTTPShutdownTimeout                        time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
	AuthLoginFailedAttemptLimitPerMinute       int           `env:"AUTH_LOGIN_FAILED_ATTEMPT_LIMIT_PER_MINUTE" envDefault:"5"`
	AuthLoginCooldownSeconds                   int           `env:"AUTH_LOGIN_COOLDOWN_SECONDS" envDefault:"300"`
	AuthEmailVerificationResendCooldownSeconds int           `env:"AUTH_EMAIL_VERIFICATION_RESEND_COOLDOWN_SECONDS" envDefault:"60"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	if strings.TrimSpace(cfg.AppVersion) == "" {
		return Config{}, fmt.Errorf("APP_VERSION cannot be empty")
	}

	if strings.TrimSpace(cfg.DBMainDSN) == "" {
		return Config{}, fmt.Errorf("DB_MAIN_DSN cannot be empty")
	}

	if strings.TrimSpace(cfg.DBTestDSN) == "" {
		return Config{}, fmt.Errorf("DB_TEST_DSN cannot be empty")
	}

	return cfg, nil
}

func (c Config) Addr() string {
	return ":" + strconv.Itoa(c.Port)
}

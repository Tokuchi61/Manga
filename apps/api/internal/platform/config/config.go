package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env                 string        `env:"APP_ENV" envDefault:"local"`
	Port                int           `env:"APP_PORT" envDefault:"8080"`
	AppVersion          string        `env:"APP_VERSION" envDefault:"0.1.0-alpha.0"`
	DBMainDSN           string        `env:"DB_MAIN_DSN,required"`
	DBTestDSN           string        `env:"DB_TEST_DSN,required"`
	DBMaxConns          int32         `env:"DB_MAX_CONNS" envDefault:"10"`
	DBConnectTimeout    time.Duration `env:"DB_CONNECT_TIMEOUT" envDefault:"5s"`
	HTTPShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
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

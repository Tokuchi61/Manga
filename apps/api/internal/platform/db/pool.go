package db

import (
	"context"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DBMainDSN)
	if err != nil {
		return nil, fmt.Errorf("parse db dsn: %w", err)
	}

	if cfg.DBMaxConns > 0 {
		poolCfg.MaxConns = cfg.DBMaxConns
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.DBConnectTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

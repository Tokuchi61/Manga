package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/db"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.Env)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.New(ctx, cfg)
	if err != nil {
		log.Fatal("database init failed", zap.Error(err))
	}
	defer pool.Close()

	auth := authmodule.New(authmodule.RuntimeConfig{
		FailedAttemptLimitPerMinute:       cfg.AuthLoginFailedAttemptLimitPerMinute,
		LoginCooldownSeconds:              cfg.AuthLoginCooldownSeconds,
		VerificationResendCooldownSeconds: cfg.AuthEmailVerificationResendCooldownSeconds,
	})

	registry, err := modules.NewRegistry(auth)
	if err != nil {
		log.Fatal("module registry init failed", zap.Error(err))
	}

	srv := app.New(cfg, log, pool, registry)
	if err := srv.Run(ctx); err != nil {
		log.Fatal("server stopped with error", zap.Error(err))
	}
}

package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	accessmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/access"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	chaptermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter"
	commentmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment"
	mangamodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga"
	supportmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/support"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/db"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/logger"
	"github.com/jackc/pgx/v5/pgxpool"
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

	var pool *pgxpool.Pool
	if strings.TrimSpace(cfg.DBMainDSN) != "" {
		pool, err = db.New(ctx, cfg)
		if err != nil {
			log.Fatal("database init failed", zap.Error(err))
		}
		defer pool.Close()
	} else {
		log.Warn("database disabled; running in memory mode")
	}

	auth := authmodule.New(authmodule.RuntimeConfig{
		FailedAttemptLimitPerMinute:       cfg.AuthLoginFailedAttemptLimitPerMinute,
		LoginCooldownSeconds:              cfg.AuthLoginCooldownSeconds,
		VerificationResendCooldownSeconds: cfg.AuthEmailVerificationResendCooldownSeconds,
	})
	user := usermodule.New()
	manga := mangamodule.New()
	chapter := chaptermodule.New()
	comment := commentmodule.New()
	support := supportmodule.New()
	access := accessmodule.New(accessmodule.RuntimeConfig{})

	user.SetCredentialLookup(auth)
	chapter.SetMangaLookup(manga)
	comment.SetTargetLookups(manga, chapter)
	support.SetTargetLookups(manga, chapter, comment)

	registry, err := modules.NewRegistry(auth, user, access, manga, chapter, comment, support)
	if err != nil {
		log.Fatal("module registry init failed", zap.Error(err))
	}

	srv := app.New(cfg, log, pool, registry)
	if err := srv.Run(ctx); err != nil {
		log.Fatal("server stopped with error", zap.Error(err))
	}
}

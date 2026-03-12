package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/app"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	accessmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/access"
	authmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth"
	chaptermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter"
	commentmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment"
	historymodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/history"
	mangamodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga"
	moderationmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation"
	notificationmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/notification"
	socialmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/social"
	supportmodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/support"
	usermodule "github.com/Tokuchi61/Manga/apps/api/internal/modules/user"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/db"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/logger"
	snapshotplatform "github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type snapshotTarget struct {
	name     string
	snapshot func() ([]byte, error)
	restore  func([]byte) error
}

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
	history := historymodule.New()
	comment := commentmodule.New()
	support := supportmodule.New()
	moderation := moderationmodule.New()
	notification := notificationmodule.New()
	social := socialmodule.New()
	access := accessmodule.New(accessmodule.RuntimeConfig{})

	user.SetCredentialLookup(auth)
	chapter.SetMangaLookup(manga)
	history.SetChapterSignalProvider(chapter)
	comment.SetTargetLookups(manga, chapter)
	support.SetTargetLookups(manga, chapter, comment)
	moderation.SetSupportContracts(support, support)
	notification.SetSupportSignalProvider(support)

	snapshotStore := buildSnapshotStore(ctx, cfg, log, pool)
	targets := []snapshotTarget{
		{name: "auth", snapshot: (&auth).Snapshot, restore: (&auth).RestoreSnapshot},
		{name: "user", snapshot: (&user).Snapshot, restore: (&user).RestoreSnapshot},
		{name: "access", snapshot: (&access).Snapshot, restore: (&access).RestoreSnapshot},
		{name: "manga", snapshot: (&manga).Snapshot, restore: (&manga).RestoreSnapshot},
		{name: "chapter", snapshot: (&chapter).Snapshot, restore: (&chapter).RestoreSnapshot},
		{name: "history", snapshot: (&history).Snapshot, restore: (&history).RestoreSnapshot},
		{name: "comment", snapshot: (&comment).Snapshot, restore: (&comment).RestoreSnapshot},
		{name: "support", snapshot: (&support).Snapshot, restore: (&support).RestoreSnapshot},
		{name: "moderation", snapshot: (&moderation).Snapshot, restore: (&moderation).RestoreSnapshot},
		{name: "notification", snapshot: (&notification).Snapshot, restore: (&notification).RestoreSnapshot},
		{name: "social", snapshot: (&social).Snapshot, restore: (&social).RestoreSnapshot},
	}
	restoreSnapshots(ctx, log, snapshotStore, targets)
	persistSnapshots(context.Background(), log, snapshotStore, targets, cfg.HTTPShutdownTimeout)

	if cfg.StateSnapshotInterval > 0 {
		ticker := time.NewTicker(cfg.StateSnapshotInterval)
		defer ticker.Stop()

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					persistSnapshots(context.Background(), log, snapshotStore, targets, cfg.HTTPShutdownTimeout)
				}
			}
		}()
	}

	registry, err := modules.NewRegistry(auth, user, access, manga, chapter, history, comment, support, moderation, notification, social)
	if err != nil {
		log.Fatal("module registry init failed", zap.Error(err))
	}

	srv := app.New(cfg, log, pool, registry)
	if err := srv.Run(ctx); err != nil {
		log.Fatal("server stopped with error", zap.Error(err))
	}

	persistSnapshots(context.Background(), log, snapshotStore, targets, cfg.HTTPShutdownTimeout)
}

func buildSnapshotStore(ctx context.Context, cfg config.Config, log *zap.Logger, pool *pgxpool.Pool) snapshotplatform.Store {
	if pool != nil {
		postgresStore := snapshotplatform.NewPostgresStore(pool)
		if err := postgresStore.EnsureSchema(ctx); err != nil {
			log.Warn("postgres snapshot store init failed; falling back to file store", zap.Error(err))
		} else {
			log.Info("snapshot persistence enabled", zap.String("backend", "postgres"))
			return postgresStore
		}
	}

	log.Info("snapshot persistence enabled", zap.String("backend", "file"), zap.String("dir", cfg.StateSnapshotDir))
	return snapshotplatform.NewFileStore(cfg.StateSnapshotDir)
}

func restoreSnapshots(ctx context.Context, log *zap.Logger, store snapshotplatform.Store, targets []snapshotTarget) {
	for _, target := range targets {
		payload, err := store.Load(ctx, target.name)
		if err != nil {
			log.Warn("snapshot load failed", zap.String("module", target.name), zap.Error(err))
			continue
		}
		if len(payload) == 0 {
			continue
		}
		if err := target.restore(payload); err != nil {
			log.Warn("snapshot restore failed", zap.String("module", target.name), zap.Error(err))
			continue
		}
		log.Info("snapshot restored", zap.String("module", target.name))
	}
}

func persistSnapshots(baseCtx context.Context, log *zap.Logger, store snapshotplatform.Store, targets []snapshotTarget, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(baseCtx, timeout)
	defer cancel()

	for _, target := range targets {
		payload, err := target.snapshot()
		if err != nil {
			log.Warn("snapshot capture failed", zap.String("module", target.name), zap.Error(err))
			continue
		}
		if len(payload) == 0 {
			continue
		}
		if err := store.Save(ctx, target.name, payload); err != nil {
			log.Warn("snapshot save failed", zap.String("module", target.name), zap.Error(err))
		}
	}
}

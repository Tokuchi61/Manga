package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Server struct {
	cfg        config.Config
	logger     *zap.Logger
	httpServer *http.Server
}

func New(cfg config.Config, logger *zap.Logger, pool *pgxpool.Pool) *Server {
	router := newRouter(cfg, logger, pool)

	return &Server{
		cfg:    cfg,
		logger: logger,
		httpServer: &http.Server{
			Addr:              cfg.Addr(),
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		s.logger.Info(
			"api server starting",
			zap.String("addr", s.httpServer.Addr),
			zap.String("version", s.cfg.AppVersion),
		)

		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.HTTPShutdownTimeout)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

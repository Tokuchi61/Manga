package app

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/config"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type maintenanceStateProvider interface {
	MaintenanceEnabled(ctx context.Context) (bool, error)
}

func NewHTTPHandler(cfg config.Config, logger *zap.Logger, pool *pgxpool.Pool, registry *modules.Registry) http.Handler {
	if logger == nil {
		logger = zap.NewNop()
	}
	if registry == nil {
		registry = modules.EmptyRegistry()
	}

	identity.SetAccessTokenSecret(cfg.AuthAccessTokenSecret)
	maintenanceProvider := resolveMaintenanceProvider(registry)

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(requestLoggingMiddleware(logger))
	r.Use(maintenanceModeMiddleware(maintenanceProvider))
	if cfg.StateSnapshotWriteThrough {
		r.Use(snapshotWriteThroughMiddleware())
	}

	registerCoreRoutes(r, cfg, pool)
	registry.MountRoutes(r)

	return r
}

func registerCoreRoutes(r chi.Router, cfg config.Config, pool *pgxpool.Pool) {
	r.Get("/health", healthHandler())
	r.Get("/ready", readyHandler(pool))
	r.Get("/version", versionHandler(cfg.AppVersion))
}

func requestLoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info(
				"http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("request_id", chiMiddleware.GetReqID(r.Context())),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

func maintenanceModeMiddleware(provider maintenanceStateProvider) func(http.Handler) http.Handler {
	if provider == nil {
		return func(next http.Handler) http.Handler { return next }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isMaintenanceBypassPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			enabled, err := provider.MaintenanceEnabled(r.Context())
			if err != nil {
				_ = writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "maintenance_state_unavailable"})
				return
			}
			if enabled {
				_ = writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "site_maintenance_enabled"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func snapshotWriteThroughMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			if !isMutatingMethod(r.Method) {
				return
			}
			hook := currentPostWriteHook()
			if hook == nil {
				return
			}
			hook(r.Context())
		})
	}
}

func resolveMaintenanceProvider(registry *modules.Registry) maintenanceStateProvider {
	if registry == nil {
		return nil
	}
	for _, module := range registry.Modules() {
		provider, ok := module.(maintenanceStateProvider)
		if !ok {
			continue
		}
		if strings.EqualFold(module.Name(), "admin") {
			return provider
		}
	}
	return nil
}

func isMutatingMethod(method string) bool {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func isMaintenanceBypassPath(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	if path == "/health" || path == "/ready" || path == "/version" {
		return true
	}
	if strings.HasPrefix(path, "/admin") {
		return true
	}
	if path == "/payment/callback" {
		return true
	}
	return false
}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func readyHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pool == nil {
			_ = writeJSON(w, http.StatusOK, map[string]string{"status": "ready", "mode": "memory"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			_ = writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}

		_ = writeJSON(w, http.StatusOK, map[string]string{"status": "ready", "mode": "database"})
	}
}

func versionHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = writeJSON(w, http.StatusOK, map[string]string{"version": version})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

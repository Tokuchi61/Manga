package auth

import (
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// RuntimeConfig maps stage-4 auth runtime keys.
type RuntimeConfig struct {
	FailedAttemptLimitPerMinute       int
	LoginCooldownSeconds              int
	VerificationResendCooldownSeconds int
}

// Module wires auth handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
}

func New(cfg RuntimeConfig) Module {
	store := authrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator, service.Config{
		FailedAttemptLimit:         cfg.FailedAttemptLimitPerMinute,
		LoginCooldown:              time.Duration(cfg.LoginCooldownSeconds) * time.Second,
		VerificationResendCooldown: time.Duration(cfg.VerificationResendCooldownSeconds) * time.Second,
	})

	return Module{httpHandler: handler.New(svc)}
}

func (m Module) Name() string {
	return "auth"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

package auth

import (
	"context"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/handler"
	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// RuntimeConfig maps stage-4 auth runtime keys.
type RuntimeConfig struct {
	FailedAttemptLimitPerMinute       int
	LoginCooldownSeconds              int
	VerificationResendCooldownSeconds int
	AccessTokenSecret                 string
	AccessTokenIssuer                 string
	ExposeSensitiveTokens             bool
}

// Module wires auth handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.AuthService
	snapshotter snapshot.Snapshotter
}

func New(cfg RuntimeConfig) Module {
	store := authrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator, service.Config{
		FailedAttemptLimit:         cfg.FailedAttemptLimitPerMinute,
		LoginCooldown:              time.Duration(cfg.LoginCooldownSeconds) * time.Second,
		VerificationResendCooldown: time.Duration(cfg.VerificationResendCooldownSeconds) * time.Second,
		AccessTokenSecret:          cfg.AccessTokenSecret,
		AccessTokenIssuer:          cfg.AccessTokenIssuer,
		ExposeSensitiveTokens:      cfg.ExposeSensitiveTokens,
	})

	return Module{httpHandler: handler.New(svc), svc: svc, snapshotter: store}
}

func (m Module) Name() string {
	return "auth"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

func (m Module) CredentialExists(ctx context.Context, credentialID string) (bool, error) {
	if m.svc == nil {
		return false, nil
	}
	return m.svc.CredentialExists(ctx, credentialID)
}

func (m *Module) SetUserLookup(lookup service.UserLookup) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetUserLookup(lookup)
}

func (m *Module) Snapshot() ([]byte, error) {
	if m == nil || m.snapshotter == nil {
		return nil, nil
	}
	return m.snapshotter.Snapshot()
}

func (m *Module) RestoreSnapshot(data []byte) error {
	if m == nil || m.snapshotter == nil {
		return nil
	}
	return m.snapshotter.RestoreSnapshot(data)
}

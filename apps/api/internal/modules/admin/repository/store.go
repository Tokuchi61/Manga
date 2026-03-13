package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
)

var (
	ErrNotFound = errors.New("admin_repository_not_found")
)

// Store defines admin persistence boundary.
type Store interface {
	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error

	CreateAdminAction(ctx context.Context, action entity.AdminAction) error
	ListAdminActions(ctx context.Context, actionType string, riskLevel string, limit int, offset int) ([]entity.AdminAction, error)
	GetActionDedup(ctx context.Context, dedupKey string) (entity.AdminAction, error)
	PutActionDedup(ctx context.Context, dedupKey string, action entity.AdminAction) error

	CreateOverrideRecord(ctx context.Context, record entity.OverrideRecord) error
	ListOverrideRecords(ctx context.Context, targetModule string, limit int, offset int) ([]entity.OverrideRecord, error)

	CreateUserReviewRecord(ctx context.Context, record entity.UserReviewRecord) error
	ListUserReviewRecords(ctx context.Context, targetUserID string, limit int, offset int) ([]entity.UserReviewRecord, error)

	CreateImpersonationSession(ctx context.Context, session entity.ImpersonationSession) error
	GetImpersonationSession(ctx context.Context, sessionID string) (entity.ImpersonationSession, error)
	UpdateImpersonationSession(ctx context.Context, session entity.ImpersonationSession) error
	FindActiveImpersonation(ctx context.Context, actorUserID string, targetUserID string) (entity.ImpersonationSession, error)
	ListImpersonationSessions(ctx context.Context, activeOnly bool, limit int, offset int) ([]entity.ImpersonationSession, error)
}

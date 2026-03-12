package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("access_repository_not_found")
	ErrConflict = errors.New("access_repository_conflict")
)

// Store defines access persistence boundary.
type Store interface {
	CreateRole(ctx context.Context, role entity.Role) error
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (entity.Role, error)
	GetRoleByName(ctx context.Context, roleName string) (entity.Role, error)
	ListRoles(ctx context.Context) ([]entity.Role, error)

	CreatePermission(ctx context.Context, permission entity.Permission) error
	GetPermissionByID(ctx context.Context, permissionID uuid.UUID) (entity.Permission, error)
	GetPermissionByName(ctx context.Context, permissionName string) (entity.Permission, error)
	ListPermissions(ctx context.Context) ([]entity.Permission, error)

	AttachPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID, createdAt time.Time) error
	ListPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)

	AssignRoleToUser(ctx context.Context, userRole entity.UserRole) error
	ListUserRoles(ctx context.Context, userID string) ([]entity.UserRole, error)

	CreatePolicyRule(ctx context.Context, rule entity.PolicyRule) error
	ListPolicyRulesByKey(ctx context.Context, key string) ([]entity.PolicyRule, error)

	CreateTemporaryGrant(ctx context.Context, grant entity.TemporaryGrant) error
	ListTemporaryGrantsByUser(ctx context.Context, userID string) ([]entity.TemporaryGrant, error)
}

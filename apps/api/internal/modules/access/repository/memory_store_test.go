package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreRolePermissionAndPolicyLifecycle(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)

	role := entity.Role{ID: uuid.New(), Name: "reader", Priority: 10, CreatedAt: now, UpdatedAt: now}
	require.NoError(t, store.CreateRole(ctx, role))

	permission := entity.Permission{
		ID:           uuid.New(),
		Name:         "history.timeline.read.own",
		Module:       "history",
		Surface:      "timeline",
		Action:       "read",
		AudienceKind: "authenticated",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, store.CreatePermission(ctx, permission))
	require.NoError(t, store.AttachPermissionToRole(ctx, role.ID, permission.ID, now))

	rolePermissions, err := store.ListPermissionsByRole(ctx, role.ID)
	require.NoError(t, err)
	require.Len(t, rolePermissions, 1)
	require.Equal(t, "history.timeline.read.own", rolePermissions[0].Name)

	rule := entity.PolicyRule{
		ID:               uuid.New(),
		Key:              "feature.history.timeline.enabled",
		Effect:           entity.PolicyEffectAllow,
		AudienceKind:     "authenticated",
		AudienceSelector: "-",
		ScopeKind:        "feature",
		ScopeSelector:    "history.timeline",
		Active:           true,
		Version:          1,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	require.NoError(t, store.CreatePolicyRule(ctx, rule))

	loadedRules, err := store.ListPolicyRulesByKey(ctx, "feature.history.timeline.enabled")
	require.NoError(t, err)
	require.Len(t, loadedRules, 1)
	require.Equal(t, entity.PolicyEffectAllow, loadedRules[0].Effect)
}

func TestMemoryStoreRejectsPolicyConflictAndDuplicateAssignments(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)

	role := entity.Role{ID: uuid.New(), Name: "moderator", Priority: 40, CreatedAt: now, UpdatedAt: now}
	require.NoError(t, store.CreateRole(ctx, role))

	permission := entity.Permission{
		ID:           uuid.New(),
		Name:         "moderation.action.any",
		Module:       "moderation",
		Surface:      "action",
		Action:       "write",
		AudienceKind: "authenticated",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, store.CreatePermission(ctx, permission))
	require.NoError(t, store.AttachPermissionToRole(ctx, role.ID, permission.ID, now))
	require.ErrorIs(t, store.AttachPermissionToRole(ctx, role.ID, permission.ID, now), ErrConflict)

	userRole := entity.UserRole{UserID: uuid.NewString(), RoleID: role.ID, CreatedAt: now}
	require.NoError(t, store.AssignRoleToUser(ctx, userRole))
	require.ErrorIs(t, store.AssignRoleToUser(ctx, userRole), ErrConflict)

	policyOne := entity.PolicyRule{
		ID:               uuid.New(),
		Key:              "feature.moderation.action.enabled",
		Effect:           entity.PolicyEffectAllow,
		AudienceKind:     "authenticated",
		AudienceSelector: "-",
		ScopeKind:        "feature",
		ScopeSelector:    "moderation.action",
		Active:           true,
		Version:          1,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	require.NoError(t, store.CreatePolicyRule(ctx, policyOne))

	policyTwo := policyOne
	policyTwo.ID = uuid.New()
	policyTwo.Effect = entity.PolicyEffectDeny
	policyTwo.Version = 2
	require.ErrorIs(t, store.CreatePolicyRule(ctx, policyTwo), ErrConflict)
}

func TestMemoryStoreTemporaryGrantLifecycle(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)

	permission := entity.Permission{
		ID:           uuid.New(),
		Name:         "payment.transaction.read.own",
		Module:       "payment",
		Surface:      "transaction",
		Action:       "read",
		AudienceKind: "authenticated",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, store.CreatePermission(ctx, permission))

	userID := uuid.NewString()
	grant := entity.TemporaryGrant{
		ID:           uuid.New(),
		UserID:       userID,
		PermissionID: permission.ID,
		Reason:       "temporary_support_access",
		ExpiresAt:    now.Add(time.Hour),
		CreatedAt:    now,
	}
	require.NoError(t, store.CreateTemporaryGrant(ctx, grant))

	grants, err := store.ListTemporaryGrantsByUser(ctx, userID)
	require.NoError(t, err)
	require.Len(t, grants, 1)
	require.Equal(t, "temporary_support_access", grants[0].Reason)
}

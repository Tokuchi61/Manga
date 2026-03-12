package service

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
	accessrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func newTestService(nowRef *time.Time) *AccessService {
	svc := New(accessrepository.NewMemoryStore(), validation.New(), Config{})
	svc.now = func() time.Time { return nowRef.UTC() }
	return svc
}

func boolPtr(value bool) *bool {
	return &value
}

func TestEvaluateGuestAndAuthenticatedFlows(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	guestSiteDecision, err := svc.Evaluate(ctx, dto.EvaluateRequest{Permission: "site.view"})
	require.NoError(t, err)
	require.True(t, guestSiteDecision.Allowed)
	require.Equal(t, "guest", guestSiteDecision.SubjectKind)

	guestChapterDecision, err := svc.Evaluate(ctx, dto.EvaluateRequest{Permission: "chapter.read.authenticated"})
	require.NoError(t, err)
	require.False(t, guestChapterDecision.Allowed)
	require.Equal(t, "chapter_requires_authenticated", guestChapterDecision.ReasonCode)

	userID := uuid.NewString()
	authenticatedDecision, err := svc.Evaluate(ctx, dto.EvaluateRequest{
		UserID:     userID,
		Permission: "chapter.read.authenticated",
		Identity:   &dto.EvaluateIdentity{CredentialID: uuid.NewString(), EmailVerified: true},
		UserSignal: &dto.EvaluateUserSignal{AccountState: "active"},
	})
	require.NoError(t, err)
	require.True(t, authenticatedDecision.Allowed)
	require.Equal(t, "authenticated", authenticatedDecision.SubjectKind)
}

func TestEvaluatePolicyPrecedenceAndConflict(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	_, err := svc.CreatePolicyRule(ctx, dto.CreatePolicyRuleRequest{
		Key:              "feature.chapter.read.enabled",
		Effect:           "allow",
		AudienceKind:     "authenticated_non_vip",
		AudienceSelector: "-",
		ScopeKind:        "feature",
		ScopeSelector:    "chapter.read",
		Active:           boolPtr(true),
	})
	require.NoError(t, err)

	_, err = svc.CreatePolicyRule(ctx, dto.CreatePolicyRuleRequest{
		Key:              "feature.chapter.read.enabled",
		Effect:           "emergency_deny",
		AudienceKind:     "all",
		AudienceSelector: "-",
		ScopeKind:        "feature",
		ScopeSelector:    "chapter.read",
		Active:           boolPtr(true),
	})
	require.NoError(t, err)

	_, err = svc.CreatePolicyRule(ctx, dto.CreatePolicyRuleRequest{
		Key:              "feature.chapter.read.enabled",
		Effect:           "deny",
		AudienceKind:     "all",
		AudienceSelector: "-",
		ScopeKind:        "feature",
		ScopeSelector:    "chapter.read",
		Active:           boolPtr(true),
	})
	require.ErrorIs(t, err, ErrPolicyConflict)

	decision, err := svc.Evaluate(ctx, dto.EvaluateRequest{
		UserID:        uuid.NewString(),
		Permission:    "chapter.read.authenticated",
		ScopeKind:     "feature",
		ScopeSelector: "chapter.read",
		Identity:      &dto.EvaluateIdentity{CredentialID: uuid.NewString(), EmailVerified: true},
		UserSignal:    &dto.EvaluateUserSignal{AccountState: "active"},
	})
	require.NoError(t, err)
	require.False(t, decision.Allowed)
	require.Equal(t, "emergency_deny", decision.ReasonCode)
	require.GreaterOrEqual(t, decision.PolicyVersion, 1)
}

func TestEvaluateOwnScopeAndTemporaryGrant(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	userID := uuid.NewString()
	otherUserID := uuid.NewString()
	identity := &dto.EvaluateIdentity{CredentialID: uuid.NewString(), EmailVerified: true}
	userSignal := &dto.EvaluateUserSignal{AccountState: "active"}

	ownDenied, err := svc.Evaluate(ctx, dto.EvaluateRequest{
		UserID:              userID,
		Permission:          "history.timeline.read.own",
		ResourceOwnerUserID: otherUserID,
		Identity:            identity,
		UserSignal:          userSignal,
	})
	require.NoError(t, err)
	require.False(t, ownDenied.Allowed)
	require.Equal(t, "ownership_mismatch", ownDenied.ReasonCode)

	_, err = svc.CreatePermission(ctx, dto.CreatePermissionRequest{
		Name:         "support.internal_note.write.any",
		Module:       "support",
		Surface:      "internal_note",
		Action:       "write",
		AudienceKind: "authenticated",
	})
	require.NoError(t, err)

	noGrantDecision, err := svc.Evaluate(ctx, dto.EvaluateRequest{
		UserID:     userID,
		Permission: "support.internal_note.write.any",
		Identity:   identity,
		UserSignal: userSignal,
	})
	require.NoError(t, err)
	require.False(t, noGrantDecision.Allowed)
	require.Equal(t, "permission_not_granted", noGrantDecision.ReasonCode)

	_, err = svc.CreateTemporaryGrant(ctx, userID, dto.CreateTemporaryGrantRequest{
		PermissionName: "support.internal_note.write.any",
		Reason:         "temporary_support_scope",
		ExpiresAt:      now.Add(time.Hour),
	})
	require.NoError(t, err)

	withGrantDecision, err := svc.Evaluate(ctx, dto.EvaluateRequest{
		UserID:     userID,
		Permission: "support.internal_note.write.any",
		Identity:   identity,
		UserSignal: userSignal,
	})
	require.NoError(t, err)
	require.True(t, withGrantDecision.Allowed)
}

func TestEvaluateSuperAdminBypassAndAdminOverride(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	_, err := svc.CreatePermission(ctx, dto.CreatePermissionRequest{
		Name:         "admin.secret.panel.any",
		Module:       "admin",
		Surface:      "panel",
		Action:       "read",
		AudienceKind: "authenticated",
	})
	require.NoError(t, err)

	superAdminUser := uuid.NewString()
	_, err = svc.AssignRoleToUser(ctx, superAdminUser, dto.AssignRoleToUserRequest{RoleName: "super_admin"})
	require.NoError(t, err)

	superAdminDecision, err := svc.Evaluate(ctx, dto.EvaluateRequest{
		UserID:     superAdminUser,
		Permission: "admin.secret.panel.any",
		Identity:   &dto.EvaluateIdentity{CredentialID: uuid.NewString(), EmailVerified: true},
		UserSignal: &dto.EvaluateUserSignal{AccountState: "active"},
	})
	require.NoError(t, err)
	require.True(t, superAdminDecision.Allowed)
	require.Equal(t, "super_admin_bypass", superAdminDecision.ReasonCode)

	adminUser := uuid.NewString()
	_, err = svc.AssignRoleToUser(ctx, adminUser, dto.AssignRoleToUserRequest{RoleName: "admin"})
	require.NoError(t, err)

	adminOverrideDecision, err := svc.Evaluate(ctx, dto.EvaluateRequest{
		UserID:        adminUser,
		Permission:    "admin.secret.panel.any",
		AllowOverride: true,
		Identity:      &dto.EvaluateIdentity{CredentialID: uuid.NewString(), EmailVerified: true},
		UserSignal:    &dto.EvaluateUserSignal{AccountState: "active"},
	})
	require.NoError(t, err)
	require.True(t, adminOverrideDecision.Allowed)
	require.Equal(t, "admin_override", adminOverrideDecision.ReasonCode)
}

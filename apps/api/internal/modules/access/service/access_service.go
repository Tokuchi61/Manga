package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	accesscontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	accessrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/repository"
	accessvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
)

var (
	ErrValidation              = errors.New("access_validation_failed")
	ErrRoleAlreadyExists       = errors.New("access_role_already_exists")
	ErrRoleNotFound            = errors.New("access_role_not_found")
	ErrPermissionAlreadyExists = errors.New("access_permission_already_exists")
	ErrPermissionNotFound      = errors.New("access_permission_not_found")
	ErrPolicyConflict          = errors.New("access_policy_conflict")
	ErrAuthorizationDenied     = errors.New("access_authorization_denied")
)

// Config defines access runtime settings.
type Config struct{}

// AccessService owns stage-6 authorization and policy flows.
type AccessService struct {
	store     accessrepository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store accessrepository.Store, validator *validation.Validator, cfg Config) *AccessService {
	if store == nil {
		store = accessrepository.NewMemoryStore()
	}

	_ = cfg
	svc := &AccessService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
	if err := svc.bootstrapDefaults(); err != nil {
		panic(fmt.Sprintf("access bootstrap failed: %v", err))
	}
	return svc
}
func (s *AccessService) validateInput(payload any) error {
	if err := accessvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *AccessService) bootstrapDefaults() error {
	ctx := context.Background()
	now := s.now().UTC()

	if err := s.ensureRole(ctx, now, "guest", 10, true, false); err != nil {
		return err
	}
	if err := s.ensureRole(ctx, now, "authenticated", 20, true, false); err != nil {
		return err
	}
	if err := s.ensureRole(ctx, now, "vip", 30, false, false); err != nil {
		return err
	}
	if err := s.ensureRole(ctx, now, "moderator", 40, false, false); err != nil {
		return err
	}
	if err := s.ensureRole(ctx, now, "admin", 50, false, false); err != nil {
		return err
	}
	if err := s.ensureRole(ctx, now, "super_admin", 100, false, true); err != nil {
		return err
	}

	seedPermissions := []struct {
		Name         string
		Module       string
		Surface      string
		Action       string
		AudienceKind string
		Roles        []string
	}{
		{Name: accesscontract.PermissionSiteView, Module: "site", Surface: "view", Action: "read", AudienceKind: "all", Roles: []string{"guest", "authenticated", "vip"}},
		{Name: accesscontract.PermissionMangaListView, Module: "manga", Surface: "list", Action: "read", AudienceKind: "all", Roles: []string{"guest", "authenticated", "vip"}},
		{Name: accesscontract.PermissionMangaDetailView, Module: "manga", Surface: "detail", Action: "read", AudienceKind: "all", Roles: []string{"guest", "authenticated", "vip"}},
		{Name: accesscontract.PermissionMangaDiscoveryView, Module: "manga", Surface: "discovery", Action: "read", AudienceKind: "all", Roles: []string{"guest", "authenticated", "vip"}},
		{Name: accesscontract.PermissionChapterReadAuthenticated, Module: "chapter", Surface: "read", Action: "read", AudienceKind: "authenticated", Roles: []string{"authenticated", "vip"}},
		{Name: accesscontract.PermissionChapterEarlyAccessVIP, Module: "chapter", Surface: "early_access", Action: "read", AudienceKind: "vip", Roles: []string{"vip"}},
		{Name: accesscontract.PermissionCommentWriteAuthenticated, Module: "comment", Surface: "write", Action: "write", AudienceKind: "authenticated", Roles: []string{"authenticated", "vip"}},
		{Name: accesscontract.PermissionHistoryContinueReadingOwn, Module: "history", Surface: "continue_reading", Action: "read", AudienceKind: "authenticated", Roles: []string{"authenticated", "vip"}},
		{Name: accesscontract.PermissionHistoryTimelineReadOwn, Module: "history", Surface: "timeline", Action: "read", AudienceKind: "authenticated", Roles: []string{"authenticated", "vip"}},
		{Name: accesscontract.PermissionHistoryLibraryReadOwn, Module: "history", Surface: "library", Action: "read", AudienceKind: "authenticated", Roles: []string{"authenticated", "vip"}},
		{Name: accesscontract.PermissionHistoryLibraryReadPublic, Module: "history", Surface: "library", Action: "read", AudienceKind: "all", Roles: []string{"guest", "authenticated", "vip"}},
		{Name: accesscontract.PermissionHistoryBookmarkWriteOwn, Module: "history", Surface: "bookmark", Action: "write", AudienceKind: "authenticated", Roles: []string{"authenticated", "vip"}},
		{Name: accesscontract.PermissionAdsView, Module: "ads", Surface: "view", Action: "read", AudienceKind: "all", Roles: []string{"guest", "authenticated", "vip"}},
		{Name: accesscontract.PermissionShopItemPurchase, Module: "shop", Surface: "purchase", Action: "write", AudienceKind: "authenticated", Roles: []string{"authenticated", "vip"}},
		{Name: accesscontract.PermissionPaymentManaPurchase, Module: "payment", Surface: "mana_purchase", Action: "write", AudienceKind: "authenticated", Roles: []string{"authenticated", "vip"}},
		{Name: accesscontract.PermissionPaymentTransactionReadOwn, Module: "payment", Surface: "transaction", Action: "read", AudienceKind: "authenticated", Roles: []string{"authenticated", "vip"}},
		{Name: accesscontract.PermissionModerationPanelView, Module: "moderation", Surface: "panel", Action: "read", AudienceKind: "authenticated", Roles: []string{"moderator", "admin", "super_admin"}},
		{Name: accesscontract.PermissionModerationActionAny, Module: "moderation", Surface: "action", Action: "write", AudienceKind: "authenticated", Roles: []string{"moderator", "admin", "super_admin"}},
	}

	for _, permission := range seedPermissions {
		permissionID, permErr := s.ensurePermission(ctx, now, permission.Name, permission.Module, permission.Surface, permission.Action, permission.AudienceKind)
		if permErr != nil {
			return permErr
		}
		for _, roleName := range permission.Roles {
			if err := s.ensureRolePermission(ctx, now, roleName, permissionID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *AccessService) ensureRole(ctx context.Context, now time.Time, roleName string, priority int, isDefault bool, isSuperAdmin bool) error {
	_, err := s.store.GetRoleByName(ctx, roleName)
	if err == nil {
		return nil
	}
	if !errors.Is(err, accessrepository.ErrNotFound) {
		return err
	}
	return s.store.CreateRole(ctx, entity.Role{
		ID:           uuid.New(),
		Name:         roleName,
		Priority:     priority,
		IsDefault:    isDefault,
		IsSuperAdmin: isSuperAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

func (s *AccessService) ensurePermission(ctx context.Context, now time.Time, name string, module string, surface string, action string, audienceKind string) (uuid.UUID, error) {
	permission, err := s.store.GetPermissionByName(ctx, name)
	if err == nil {
		return permission.ID, nil
	}
	if !errors.Is(err, accessrepository.ErrNotFound) {
		return uuid.Nil, err
	}

	permission = entity.Permission{
		ID:           uuid.New(),
		Name:         name,
		Module:       module,
		Surface:      surface,
		Action:       action,
		AudienceKind: audienceKind,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if createErr := s.store.CreatePermission(ctx, permission); createErr != nil {
		if errors.Is(createErr, accessrepository.ErrConflict) {
			persisted, getErr := s.store.GetPermissionByName(ctx, name)
			if getErr != nil {
				return uuid.Nil, getErr
			}
			return persisted.ID, nil
		}
		return uuid.Nil, createErr
	}
	return permission.ID, nil
}

func (s *AccessService) ensureRolePermission(ctx context.Context, now time.Time, roleName string, permissionID uuid.UUID) error {
	role, err := s.store.GetRoleByName(ctx, roleName)
	if err != nil {
		return err
	}
	err = s.store.AttachPermissionToRole(ctx, role.ID, permissionID, now)
	if errors.Is(err, accessrepository.ErrConflict) {
		return nil
	}
	return err
}

package service

import (
	"context"
	"errors"
	"sort"
	"strings"

	accesscontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	accessrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/repository"
)

func (s *AccessService) Evaluate(ctx context.Context, request dto.EvaluateRequest) (dto.EvaluateResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.EvaluateResponse{}, err
	}

	permissionName := normalizeName(request.Permission)
	userID := normalizeName(request.UserID)
	if userID != "" {
		if _, err := parseUUID(userID, "user_id"); err != nil {
			return dto.EvaluateResponse{}, err
		}
	}

	scopeKind := normalizeOptional(request.ScopeKind, "feature")
	scopeSelector := normalizeOptional(request.ScopeSelector, "-")
	audienceSelector := normalizeOptional(request.AudienceSelector, "-")
	subjectKind := deriveSubjectKind(request.Identity, request.UserSignal)

	if subjectKind == "blocked" {
		return deniedDecision(permissionName, "deny", "subject_blocked", "subject blocked by auth/user state", 0, subjectKind), nil
	}

	permission, err := s.store.GetPermissionByName(ctx, permissionName)
	if err != nil {
		if errors.Is(err, accessrepository.ErrNotFound) {
			return deniedDecision(permissionName, "deny", "permission_not_found", "unknown permission", 0, subjectKind), nil
		}
		return dto.EvaluateResponse{}, err
	}

	policyKey := normalizeName(request.FeatureKey)
	if policyKey == "" {
		policyKey = mapPermissionToFeatureKey(permissionName)
	}

	policyVersion := 0
	if policyKey != "" {
		rules, ruleErr := s.store.ListPolicyRulesByKey(ctx, policyKey)
		if ruleErr != nil {
			return dto.EvaluateResponse{}, ruleErr
		}

		effect, version := evaluatePolicyRules(rules, subjectKind, userID, scopeKind, scopeSelector, audienceSelector)
		if version > policyVersion {
			policyVersion = version
		}
		switch effect {
		case string(entity.PolicyEffectEmergencyDeny):
			return deniedDecision(permissionName, effect, "emergency_deny", "emergency deny policy applied", policyVersion, subjectKind), nil
		case string(entity.PolicyEffectDeny):
			return deniedDecision(permissionName, effect, "policy_deny", "deny policy applied", policyVersion, subjectKind), nil
		}
	}

	roleNames, grantedPermissions, roleErr := s.resolveRoleAndPermissionSet(ctx, userID, subjectKind)
	if roleErr != nil {
		return dto.EvaluateResponse{}, roleErr
	}

	if hasRole(roleNames, "super_admin") {
		return allowedDecision(permissionName, "allow", "super_admin_bypass", "super admin bypass", policyVersion, subjectKind), nil
	}

	if request.AllowOverride && (hasRole(roleNames, "admin") || hasRole(roleNames, "moderator")) {
		return allowedDecision(permissionName, "allow", "admin_override", "admin/moderator override", policyVersion, subjectKind), nil
	}

	if strings.HasSuffix(permissionName, ".own") {
		if userID == "" {
			return deniedDecision(permissionName, "require_auth", "auth_required", "own scope needs authenticated subject", policyVersion, subjectKind), nil
		}

		ownerID := normalizeName(request.ResourceOwnerUserID)
		if ownerID == "" {
			return deniedDecision(permissionName, "deny", "owner_missing", "resource owner is required for own scope", policyVersion, subjectKind), nil
		}
		if _, err := parseUUID(ownerID, "resource_owner_user_id"); err != nil {
			return dto.EvaluateResponse{}, err
		}
		if ownerID != userID && !(hasRole(roleNames, "admin") || hasRole(roleNames, "moderator")) {
			return deniedDecision(permissionName, "deny", "ownership_mismatch", "subject is not owner", policyVersion, subjectKind), nil
		}
	}

	if permissionName == accesscontract.PermissionChapterReadAuthenticated && subjectKind == "guest" {
		return deniedDecision(permissionName, "require_auth", "chapter_requires_authenticated", "chapter read needs authenticated subject", policyVersion, subjectKind), nil
	}
	if permissionName == accesscontract.PermissionChapterEarlyAccessVIP && subjectKind != "vip" {
		return deniedDecision(permissionName, "require_entitlement", "vip_required", "early access requires vip", policyVersion, subjectKind), nil
	}
	if permissionName == accesscontract.PermissionAdsView && subjectKind == "vip" {
		return deniedDecision(permissionName, "deny", "vip_no_ads", "vip no-ads precedence", policyVersion, subjectKind), nil
	}

	if permission.AudienceKind == "authenticated" && subjectKind == "guest" {
		return deniedDecision(permissionName, "require_auth", "auth_required", "permission requires authenticated subject", policyVersion, subjectKind), nil
	}
	if permission.AudienceKind == "vip" && subjectKind != "vip" {
		return deniedDecision(permissionName, "require_entitlement", "vip_required", "permission requires vip subject", policyVersion, subjectKind), nil
	}

	if _, granted := grantedPermissions[permissionName]; !granted {
		return deniedDecision(permissionName, "deny", "permission_not_granted", "subject has no matching permission", policyVersion, subjectKind), nil
	}

	return allowedDecision(permissionName, "allow", "granted", "permission granted", policyVersion, subjectKind), nil
}

func (s *AccessService) resolveRoleAndPermissionSet(ctx context.Context, userID string, subjectKind string) (map[string]struct{}, map[string]struct{}, error) {
	roleNames := make(map[string]struct{})
	roleIDs := make(map[string]struct{})

	addRoleByName := func(roleName string) error {
		role, err := s.store.GetRoleByName(ctx, roleName)
		if err != nil {
			if errors.Is(err, accessrepository.ErrNotFound) {
				return nil
			}
			return err
		}
		roleNames[role.Name] = struct{}{}
		roleIDs[role.ID.String()] = struct{}{}
		return nil
	}

	switch subjectKind {
	case "guest":
		if err := addRoleByName("guest"); err != nil {
			return nil, nil, err
		}
	case "authenticated":
		if err := addRoleByName("authenticated"); err != nil {
			return nil, nil, err
		}
	case "vip":
		if err := addRoleByName("authenticated"); err != nil {
			return nil, nil, err
		}
		if err := addRoleByName("vip"); err != nil {
			return nil, nil, err
		}
	}

	if userID != "" {
		assignedRoles, err := s.store.ListUserRoles(ctx, userID)
		if err != nil {
			return nil, nil, err
		}
		now := s.now().UTC()
		for _, assignedRole := range assignedRoles {
			if !assignedRole.IsActive(now) {
				continue
			}
			role, err := s.store.GetRoleByID(ctx, assignedRole.RoleID)
			if err != nil {
				if errors.Is(err, accessrepository.ErrNotFound) {
					continue
				}
				return nil, nil, err
			}
			roleNames[role.Name] = struct{}{}
			roleIDs[role.ID.String()] = struct{}{}
		}
	}

	grantedPermissions := make(map[string]struct{})
	for roleID := range roleIDs {
		parsedRoleID, err := parseUUID(roleID, "role_id")
		if err != nil {
			return nil, nil, err
		}
		permissions, err := s.store.ListPermissionsByRole(ctx, parsedRoleID)
		if err != nil {
			if errors.Is(err, accessrepository.ErrNotFound) {
				continue
			}
			return nil, nil, err
		}
		for _, permission := range permissions {
			grantedPermissions[permission.Name] = struct{}{}
		}
	}

	if userID != "" {
		tempGrants, err := s.store.ListTemporaryGrantsByUser(ctx, userID)
		if err != nil {
			return nil, nil, err
		}
		now := s.now().UTC()
		for _, grant := range tempGrants {
			if !grant.IsActive(now) {
				continue
			}
			permission, err := s.store.GetPermissionByID(ctx, grant.PermissionID)
			if err != nil {
				if errors.Is(err, accessrepository.ErrNotFound) {
					continue
				}
				return nil, nil, err
			}
			grantedPermissions[permission.Name] = struct{}{}
		}
	}

	return roleNames, grantedPermissions, nil
}

func evaluatePolicyRules(rules []entity.PolicyRule, subjectKind string, userID string, scopeKind string, scopeSelector string, audienceSelector string) (string, int) {
	filtered := make([]entity.PolicyRule, 0)
	for _, rule := range rules {
		if !rule.Active {
			continue
		}
		if normalizeName(rule.ScopeKind) != normalizeName(scopeKind) {
			continue
		}
		if normalizeName(rule.ScopeSelector) != "-" && normalizeName(rule.ScopeSelector) != normalizeName(scopeSelector) {
			continue
		}
		if !matchesAudience(rule, subjectKind) {
			continue
		}
		if !matchesAudienceSelector(rule, userID, audienceSelector) {
			continue
		}
		filtered = append(filtered, rule)
	}

	if len(filtered) == 0 {
		return "", 0
	}

	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Version == filtered[j].Version {
			return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
		}
		return filtered[i].Version > filtered[j].Version
	})

	maxVersion := filtered[0].Version
	hasAllow := false
	for _, rule := range filtered {
		switch rule.Effect {
		case entity.PolicyEffectEmergencyDeny:
			return string(entity.PolicyEffectEmergencyDeny), maxVersion
		case entity.PolicyEffectDeny:
			return string(entity.PolicyEffectDeny), maxVersion
		case entity.PolicyEffectAllow:
			hasAllow = true
		}
	}
	if hasAllow {
		return string(entity.PolicyEffectAllow), maxVersion
	}
	return "", maxVersion
}

func matchesAudience(rule entity.PolicyRule, subjectKind string) bool {
	audienceKind := normalizeName(rule.AudienceKind)
	switch audienceKind {
	case "all":
		return true
	case "guest":
		return subjectKind == "guest"
	case "authenticated":
		return subjectKind == "authenticated" || subjectKind == "vip"
	case "authenticated_non_vip":
		return subjectKind == "authenticated"
	case "vip":
		return subjectKind == "vip"
	default:
		return false
	}
}

func matchesAudienceSelector(rule entity.PolicyRule, userID string, requestAudienceSelector string) bool {
	selector := normalizeName(rule.AudienceSelector)
	if selector == "-" {
		return true
	}
	if strings.HasPrefix(selector, "user:") {
		targetUserID := normalizeName(strings.TrimPrefix(selector, "user:"))
		return targetUserID != "" && targetUserID == userID
	}
	return selector == normalizeName(requestAudienceSelector)
}

func hasRole(roleNames map[string]struct{}, roleName string) bool {
	_, ok := roleNames[normalizeName(roleName)]
	return ok
}

func deriveSubjectKind(identity *dto.EvaluateIdentity, userSignal *dto.EvaluateUserSignal) string {
	if identity == nil || normalizeName(identity.CredentialID) == "" {
		return "guest"
	}
	if identity.Banned || identity.Suspended {
		return "blocked"
	}
	if userSignal == nil {
		return "authenticated"
	}
	state := normalizeName(userSignal.AccountState)
	if state == "banned" || state == "deactivated" {
		return "blocked"
	}
	if userSignal.VIPActive && !userSignal.VIPFrozen {
		return "vip"
	}
	return "authenticated"
}

func allowedDecision(permission string, effect string, reasonCode string, reason string, policyVersion int, subjectKind string) dto.EvaluateResponse {
	return dto.EvaluateResponse{
		Allowed:       true,
		Effect:        effect,
		ReasonCode:    reasonCode,
		Reason:        reason,
		Permission:    permission,
		PolicyVersion: policyVersion,
		SubjectKind:   subjectKind,
	}
}

func deniedDecision(permission string, effect string, reasonCode string, reason string, policyVersion int, subjectKind string) dto.EvaluateResponse {
	return dto.EvaluateResponse{
		Allowed:       false,
		Effect:        effect,
		ReasonCode:    reasonCode,
		Reason:        reason,
		Permission:    permission,
		PolicyVersion: policyVersion,
		SubjectKind:   subjectKind,
	}
}

func mapPermissionToFeatureKey(permission string) string {
	switch permission {
	case accesscontract.PermissionMangaDiscoveryView:
		return "feature.manga.discovery.enabled"
	case accesscontract.PermissionChapterReadAuthenticated:
		return "feature.chapter.read.enabled"
	case accesscontract.PermissionChapterEarlyAccessVIP:
		return "feature.chapter.early_access.enabled"
	case accesscontract.PermissionCommentWriteAuthenticated:
		return "feature.comment.write.enabled"
	case accesscontract.PermissionHistoryContinueReadingOwn:
		return "feature.history.continue_reading.enabled"
	case accesscontract.PermissionHistoryTimelineReadOwn:
		return "feature.history.timeline.enabled"
	case accesscontract.PermissionHistoryLibraryReadOwn, accesscontract.PermissionHistoryLibraryReadPublic:
		return "feature.history.library.enabled"
	case accesscontract.PermissionHistoryBookmarkWriteOwn:
		return "feature.history.bookmark_write.enabled"
	case accesscontract.PermissionAdsView:
		return "feature.ads.surface.enabled"
	case accesscontract.PermissionShopItemPurchase:
		return "feature.shop.purchase.enabled"
	case accesscontract.PermissionPaymentManaPurchase:
		return "feature.payment.mana_purchase.enabled"
	case accesscontract.PermissionPaymentTransactionReadOwn:
		return "feature.payment.transaction_read.enabled"
	case accesscontract.PermissionModerationPanelView:
		return "feature.moderation.panel.enabled"
	case accesscontract.PermissionModerationActionAny:
		return "feature.moderation.action.enabled"
	default:
		return ""
	}
}

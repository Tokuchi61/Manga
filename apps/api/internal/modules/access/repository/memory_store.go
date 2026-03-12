package repository

import (
	"strings"
	"sync"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	"github.com/google/uuid"
)

// MemoryStore is a stage-6 bootstrap persistence for access flows.
type MemoryStore struct {
	mu sync.RWMutex

	rolesByID     map[uuid.UUID]entity.Role
	roleNameIndex map[string]uuid.UUID

	permissionsByID     map[uuid.UUID]entity.Permission
	permissionNameIndex map[string]uuid.UUID

	rolePermissionIndex map[uuid.UUID]map[uuid.UUID]struct{}

	userRoleIndex map[string]map[uuid.UUID]entity.UserRole

	policyRulesByID      map[uuid.UUID]entity.PolicyRule
	policyByKeyIndex     map[string]map[uuid.UUID]struct{}
	activePolicyConflict map[string]uuid.UUID

	temporaryGrantsByID map[uuid.UUID]entity.TemporaryGrant
	userGrantIndex      map[string]map[uuid.UUID]struct{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rolesByID:            make(map[uuid.UUID]entity.Role),
		roleNameIndex:        make(map[string]uuid.UUID),
		permissionsByID:      make(map[uuid.UUID]entity.Permission),
		permissionNameIndex:  make(map[string]uuid.UUID),
		rolePermissionIndex:  make(map[uuid.UUID]map[uuid.UUID]struct{}),
		userRoleIndex:        make(map[string]map[uuid.UUID]entity.UserRole),
		policyRulesByID:      make(map[uuid.UUID]entity.PolicyRule),
		policyByKeyIndex:     make(map[string]map[uuid.UUID]struct{}),
		activePolicyConflict: make(map[string]uuid.UUID),
		temporaryGrantsByID:  make(map[uuid.UUID]entity.TemporaryGrant),
		userGrantIndex:       make(map[string]map[uuid.UUID]struct{}),
	}
}

func normalizeRoleName(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func normalizePermissionName(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func normalizeUserID(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func normalizePolicyKey(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

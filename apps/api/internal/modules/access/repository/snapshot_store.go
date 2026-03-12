package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	"github.com/google/uuid"
)

type accessSnapshotState struct {
	RolesByID            map[uuid.UUID]entity.Role
	RoleNameIndex        map[string]uuid.UUID
	PermissionsByID      map[uuid.UUID]entity.Permission
	PermissionNameIndex  map[string]uuid.UUID
	RolePermissionIndex  map[uuid.UUID]map[uuid.UUID]struct{}
	UserRoleIndex        map[string]map[uuid.UUID]entity.UserRole
	PolicyRulesByID      map[uuid.UUID]entity.PolicyRule
	PolicyByKeyIndex     map[string]map[uuid.UUID]struct{}
	ActivePolicyConflict map[string]uuid.UUID
	TemporaryGrantsByID  map[uuid.UUID]entity.TemporaryGrant
	UserGrantIndex       map[string]map[uuid.UUID]struct{}
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := accessSnapshotState{
		RolesByID:            s.rolesByID,
		RoleNameIndex:        s.roleNameIndex,
		PermissionsByID:      s.permissionsByID,
		PermissionNameIndex:  s.permissionNameIndex,
		RolePermissionIndex:  s.rolePermissionIndex,
		UserRoleIndex:        s.userRoleIndex,
		PolicyRulesByID:      s.policyRulesByID,
		PolicyByKeyIndex:     s.policyByKeyIndex,
		ActivePolicyConflict: s.activePolicyConflict,
		TemporaryGrantsByID:  s.temporaryGrantsByID,
		UserGrantIndex:       s.userGrantIndex,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(state); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *MemoryStore) RestoreSnapshot(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var state accessSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.RolesByID == nil {
		state.RolesByID = make(map[uuid.UUID]entity.Role)
	}
	if state.RoleNameIndex == nil {
		state.RoleNameIndex = make(map[string]uuid.UUID)
	}
	if state.PermissionsByID == nil {
		state.PermissionsByID = make(map[uuid.UUID]entity.Permission)
	}
	if state.PermissionNameIndex == nil {
		state.PermissionNameIndex = make(map[string]uuid.UUID)
	}
	if state.RolePermissionIndex == nil {
		state.RolePermissionIndex = make(map[uuid.UUID]map[uuid.UUID]struct{})
	}
	if state.UserRoleIndex == nil {
		state.UserRoleIndex = make(map[string]map[uuid.UUID]entity.UserRole)
	}
	if state.PolicyRulesByID == nil {
		state.PolicyRulesByID = make(map[uuid.UUID]entity.PolicyRule)
	}
	if state.PolicyByKeyIndex == nil {
		state.PolicyByKeyIndex = make(map[string]map[uuid.UUID]struct{})
	}
	if state.ActivePolicyConflict == nil {
		state.ActivePolicyConflict = make(map[string]uuid.UUID)
	}
	if state.TemporaryGrantsByID == nil {
		state.TemporaryGrantsByID = make(map[uuid.UUID]entity.TemporaryGrant)
	}
	if state.UserGrantIndex == nil {
		state.UserGrantIndex = make(map[string]map[uuid.UUID]struct{})
	}

	s.rolesByID = state.RolesByID
	s.roleNameIndex = state.RoleNameIndex
	s.permissionsByID = state.PermissionsByID
	s.permissionNameIndex = state.PermissionNameIndex
	s.rolePermissionIndex = state.RolePermissionIndex
	s.userRoleIndex = state.UserRoleIndex
	s.policyRulesByID = state.PolicyRulesByID
	s.policyByKeyIndex = state.PolicyByKeyIndex
	s.activePolicyConflict = state.ActivePolicyConflict
	s.temporaryGrantsByID = state.TemporaryGrantsByID
	s.userGrantIndex = state.UserGrantIndex

	return nil
}

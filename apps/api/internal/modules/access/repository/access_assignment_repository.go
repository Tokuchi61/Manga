package repository

import (
	"context"
	"sort"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	"github.com/google/uuid"
)

func (s *MemoryStore) AssignRoleToUser(_ context.Context, userRole entity.UserRole) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.rolesByID[userRole.RoleID]; !ok {
		return ErrNotFound
	}

	userID := normalizeUserID(userRole.UserID)
	if _, ok := s.userRoleIndex[userID]; !ok {
		s.userRoleIndex[userID] = make(map[uuid.UUID]entity.UserRole)
	}
	if _, exists := s.userRoleIndex[userID][userRole.RoleID]; exists {
		return ErrConflict
	}

	userRole.UserID = userID
	s.userRoleIndex[userID][userRole.RoleID] = userRole
	return nil
}

func (s *MemoryStore) ListUserRoles(_ context.Context, userID string) ([]entity.UserRole, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rolesByID := s.userRoleIndex[normalizeUserID(userID)]
	result := make([]entity.UserRole, 0, len(rolesByID))
	for _, role := range rolesByID {
		result = append(result, role)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

func (s *MemoryStore) CreateTemporaryGrant(_ context.Context, grant entity.TemporaryGrant) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.permissionsByID[grant.PermissionID]; !ok {
		return ErrNotFound
	}

	userID := normalizeUserID(grant.UserID)
	grant.UserID = userID
	s.temporaryGrantsByID[grant.ID] = grant
	if _, ok := s.userGrantIndex[userID]; !ok {
		s.userGrantIndex[userID] = make(map[uuid.UUID]struct{})
	}
	s.userGrantIndex[userID][grant.ID] = struct{}{}
	return nil
}

func (s *MemoryStore) ListTemporaryGrantsByUser(_ context.Context, userID string) ([]entity.TemporaryGrant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	grantIDs := s.userGrantIndex[normalizeUserID(userID)]
	result := make([]entity.TemporaryGrant, 0, len(grantIDs))
	for grantID := range grantIDs {
		grant, ok := s.temporaryGrantsByID[grantID]
		if ok {
			result = append(result, grant)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

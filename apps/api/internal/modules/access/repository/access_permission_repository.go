package repository

import (
	"context"
	"sort"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	"github.com/google/uuid"
)

func (s *MemoryStore) CreatePermission(_ context.Context, permission entity.Permission) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nameKey := normalizePermissionName(permission.Name)
	if _, exists := s.permissionNameIndex[nameKey]; exists {
		return ErrConflict
	}

	permission.Name = nameKey
	s.permissionsByID[permission.ID] = permission
	s.permissionNameIndex[nameKey] = permission.ID
	return nil
}

func (s *MemoryStore) GetPermissionByID(_ context.Context, permissionID uuid.UUID) (entity.Permission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	permission, ok := s.permissionsByID[permissionID]
	if !ok {
		return entity.Permission{}, ErrNotFound
	}
	return permission, nil
}

func (s *MemoryStore) GetPermissionByName(_ context.Context, permissionName string) (entity.Permission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	permissionID, ok := s.permissionNameIndex[normalizePermissionName(permissionName)]
	if !ok {
		return entity.Permission{}, ErrNotFound
	}
	permission, ok := s.permissionsByID[permissionID]
	if !ok {
		return entity.Permission{}, ErrNotFound
	}
	return permission, nil
}

func (s *MemoryStore) ListPermissions(_ context.Context) ([]entity.Permission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]entity.Permission, 0, len(s.permissionsByID))
	for _, permission := range s.permissionsByID {
		result = append(result, permission)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result, nil
}

func (s *MemoryStore) AttachPermissionToRole(_ context.Context, roleID uuid.UUID, permissionID uuid.UUID, _ time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.rolesByID[roleID]; !ok {
		return ErrNotFound
	}
	if _, ok := s.permissionsByID[permissionID]; !ok {
		return ErrNotFound
	}

	if _, ok := s.rolePermissionIndex[roleID]; !ok {
		s.rolePermissionIndex[roleID] = make(map[uuid.UUID]struct{})
	}
	if _, exists := s.rolePermissionIndex[roleID][permissionID]; exists {
		return ErrConflict
	}
	s.rolePermissionIndex[roleID][permissionID] = struct{}{}
	return nil
}

func (s *MemoryStore) ListPermissionsByRole(_ context.Context, roleID uuid.UUID) ([]entity.Permission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.rolesByID[roleID]; !ok {
		return nil, ErrNotFound
	}

	permissionIDs := s.rolePermissionIndex[roleID]
	result := make([]entity.Permission, 0, len(permissionIDs))
	for permissionID := range permissionIDs {
		permission, ok := s.permissionsByID[permissionID]
		if ok {
			result = append(result, permission)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result, nil
}

package repository

import (
	"context"
	"sort"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	"github.com/google/uuid"
)

func (s *MemoryStore) CreateRole(_ context.Context, role entity.Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nameKey := normalizeRoleName(role.Name)
	if _, exists := s.roleNameIndex[nameKey]; exists {
		return ErrConflict
	}

	role.Name = nameKey
	s.rolesByID[role.ID] = role
	s.roleNameIndex[nameKey] = role.ID
	return nil
}

func (s *MemoryStore) GetRoleByID(_ context.Context, roleID uuid.UUID) (entity.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	role, ok := s.rolesByID[roleID]
	if !ok {
		return entity.Role{}, ErrNotFound
	}
	return role, nil
}

func (s *MemoryStore) GetRoleByName(_ context.Context, roleName string) (entity.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roleID, ok := s.roleNameIndex[normalizeRoleName(roleName)]
	if !ok {
		return entity.Role{}, ErrNotFound
	}
	role, ok := s.rolesByID[roleID]
	if !ok {
		return entity.Role{}, ErrNotFound
	}
	return role, nil
}

func (s *MemoryStore) ListRoles(_ context.Context) ([]entity.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]entity.Role, 0, len(s.rolesByID))
	for _, role := range s.rolesByID {
		result = append(result, role)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Priority == result[j].Priority {
			return result[i].Name < result[j].Name
		}
		return result[i].Priority > result[j].Priority
	})

	return result, nil
}

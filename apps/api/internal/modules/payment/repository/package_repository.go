package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
)

func (s *MemoryStore) UpsertManaPackage(_ context.Context, pkg entity.ManaPackage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.manaPackagesByID[normalizeValue(pkg.PackageID)] = cloneManaPackage(pkg)
	return nil
}

func (s *MemoryStore) GetManaPackage(_ context.Context, packageID string) (entity.ManaPackage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pkg, ok := s.manaPackagesByID[normalizeValue(packageID)]
	if !ok {
		return entity.ManaPackage{}, ErrNotFound
	}
	return cloneManaPackage(pkg), nil
}

func (s *MemoryStore) ListManaPackages(_ context.Context, activeOnly bool, limit int, offset int) ([]entity.ManaPackage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]entity.ManaPackage, 0, len(s.manaPackagesByID))
	for _, pkg := range s.manaPackagesByID {
		if activeOnly && !pkg.Active {
			continue
		}
		items = append(items, cloneManaPackage(pkg))
	}

	sortByUpdatedDesc(items, func(i int, j int) bool {
		if items[i].DisplayOrder == items[j].DisplayOrder {
			return items[i].UpdatedAt.After(items[j].UpdatedAt)
		}
		return items[i].DisplayOrder < items[j].DisplayOrder
	})

	return applyOffsetLimit(items, offset, limit), nil
}

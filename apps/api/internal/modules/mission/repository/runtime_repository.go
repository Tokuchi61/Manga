package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
)

func (s *MemoryStore) GetRuntimeConfig(_ context.Context) (entity.RuntimeConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return cloneRuntimeConfig(s.runtimeConfig), nil
}

func (s *MemoryStore) UpdateRuntimeConfig(_ context.Context, cfg entity.RuntimeConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.runtimeConfig = cloneRuntimeConfig(cfg)
	return nil
}

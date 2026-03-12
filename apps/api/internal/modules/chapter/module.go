package chapter

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/handler"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires chapter handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.ChapterService
	snapshotter snapshot.Snapshotter
}

func New() Module {
	store := chapterrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc, snapshotter: store}
}

func (m *Module) SetMangaLookup(lookup service.MangaLookup) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetMangaLookup(lookup)
}

func (m Module) Name() string {
	return "chapter"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

func (m Module) TargetExists(ctx context.Context, chapterID string) (bool, error) {
	if m.svc == nil {
		return false, nil
	}
	return m.svc.TargetExists(ctx, chapterID)
}

func (m *Module) Snapshot() ([]byte, error) {
	if m == nil || m.snapshotter == nil {
		return nil, nil
	}
	return m.snapshotter.Snapshot()
}

func (m *Module) RestoreSnapshot(data []byte) error {
	if m == nil || m.snapshotter == nil {
		return nil
	}
	return m.snapshotter.RestoreSnapshot(data)
}

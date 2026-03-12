package snapshot

import "context"

// Snapshotter captures and restores in-memory module state.
type Snapshotter interface {
	Snapshot() ([]byte, error)
	RestoreSnapshot(data []byte) error
}

// Store persists module snapshots.
type Store interface {
	Load(ctx context.Context, module string) ([]byte, error)
	Save(ctx context.Context, module string, payload []byte) error
}

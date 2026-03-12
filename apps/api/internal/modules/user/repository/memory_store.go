package repository

import (
	"strings"
	"sync"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
)

// MemoryStore is a stage-5 bootstrap persistence for user flows.
type MemoryStore struct {
	mu sync.RWMutex

	usersByID           map[string]entity.UserAccount
	credentialUserIndex map[string]string
	usernameUserIndex   map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		usersByID:           make(map[string]entity.UserAccount),
		credentialUserIndex: make(map[string]string),
		usernameUserIndex:   make(map[string]string),
	}
}

func normalizeID(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

func normalizeUsername(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

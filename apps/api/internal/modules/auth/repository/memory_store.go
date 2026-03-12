package repository

import (
	"strings"
	"sync"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/google/uuid"
)

// MemoryStore is a stage-4 bootstrap persistence for auth flows.
type MemoryStore struct {
	mu sync.RWMutex

	credentialsByID    map[uuid.UUID]entity.Credential
	credentialEmailIdx map[string]uuid.UUID

	sessionsByID         map[uuid.UUID]entity.Session
	sessionCredentialIdx map[uuid.UUID]map[uuid.UUID]struct{}

	tokensByID         map[uuid.UUID]entity.Token
	tokenHashTypeIndex map[string]uuid.UUID

	securityEvents []entity.SecurityEvent
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		credentialsByID:      make(map[uuid.UUID]entity.Credential),
		credentialEmailIdx:   make(map[string]uuid.UUID),
		sessionsByID:         make(map[uuid.UUID]entity.Session),
		sessionCredentialIdx: make(map[uuid.UUID]map[uuid.UUID]struct{}),
		tokensByID:           make(map[uuid.UUID]entity.Token),
		tokenHashTypeIndex:   make(map[string]uuid.UUID),
		securityEvents:       make([]entity.SecurityEvent, 0),
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func tokenIndexKey(hash string, tokenType entity.TokenType) string {
	return string(tokenType) + ":" + hash
}

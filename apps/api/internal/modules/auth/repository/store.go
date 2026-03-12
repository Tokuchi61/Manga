package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("auth_repository_not_found")
	ErrConflict = errors.New("auth_repository_conflict")
)

// Store defines auth persistence boundary.
type Store interface {
	CreateCredential(ctx context.Context, credential entity.Credential) error
	GetCredentialByEmail(ctx context.Context, email string) (entity.Credential, error)
	GetCredentialByID(ctx context.Context, credentialID uuid.UUID) (entity.Credential, error)
	UpdateCredential(ctx context.Context, credential entity.Credential) error

	CreateSession(ctx context.Context, session entity.Session) error
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (entity.Session, error)
	ListSessionsByCredential(ctx context.Context, credentialID uuid.UUID) ([]entity.Session, error)
	UpdateSession(ctx context.Context, session entity.Session) error
	RevokeSession(ctx context.Context, sessionID uuid.UUID, revokedAtUnix int64) error
	RevokeOtherSessions(ctx context.Context, credentialID uuid.UUID, currentSessionID uuid.UUID, revokedAtUnix int64) error
	RevokeAllSessions(ctx context.Context, credentialID uuid.UUID, revokedAtUnix int64) error

	CreateToken(ctx context.Context, token entity.Token) error
	GetTokenByHashAndType(ctx context.Context, hash string, tokenType entity.TokenType) (entity.Token, error)
	UpdateToken(ctx context.Context, token entity.Token) error

	AppendSecurityEvent(ctx context.Context, event entity.SecurityEvent) error
	ListSecurityEvents(ctx context.Context, credentialID uuid.UUID) ([]entity.SecurityEvent, error)
}

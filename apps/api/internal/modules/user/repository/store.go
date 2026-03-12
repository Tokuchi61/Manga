package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
)

var (
	ErrNotFound = errors.New("user_repository_not_found")
	ErrConflict = errors.New("user_repository_conflict")
)

// Store defines user persistence boundary.
type Store interface {
	CreateUser(ctx context.Context, user entity.UserAccount) error
	GetUserByID(ctx context.Context, userID string) (entity.UserAccount, error)
	GetUserByCredentialID(ctx context.Context, credentialID string) (entity.UserAccount, error)
	UpdateUser(ctx context.Context, user entity.UserAccount) error
}

package service

import (
	"context"
	"errors"

	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	"github.com/google/uuid"
)

// CredentialExists exposes auth-owned credential existence checks for other modules.
func (s *AuthService) CredentialExists(ctx context.Context, credentialID string) (bool, error) {
	parsedID, err := uuid.Parse(credentialID)
	if err != nil {
		return false, nil
	}

	_, err = s.store.GetCredentialByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, authrepository.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

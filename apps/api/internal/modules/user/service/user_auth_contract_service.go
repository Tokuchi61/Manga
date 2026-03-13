package service

import (
	"context"
	"errors"

	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
)

func (s *UserService) ResolveUserIDByCredentialID(ctx context.Context, credentialID string) (string, error) {
	parsedCredentialID, err := parseID(credentialID, "credential_id")
	if err != nil {
		return "", err
	}

	user, err := s.store.GetUserByCredentialID(ctx, parsedCredentialID)
	if err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return "", nil
		}
		return "", err
	}
	return user.ID, nil
}

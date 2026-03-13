package service

import (
	"context"
	"strings"
)

// UserLookup resolves user identity from credential scope.
type UserLookup interface {
	ResolveUserIDByCredentialID(ctx context.Context, credentialID string) (string, error)
}

func (s *AuthService) SetUserLookup(lookup UserLookup) {
	if s == nil {
		return
	}
	s.userLookup = lookup
}

func (s *AuthService) resolveUserID(ctx context.Context, credentialID string) (string, error) {
	if s == nil || s.userLookup == nil || strings.TrimSpace(credentialID) == "" {
		return "", nil
	}
	return s.userLookup.ResolveUserIDByCredentialID(ctx, credentialID)
}

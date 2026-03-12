package service

import (
	"fmt"

	"github.com/google/uuid"
)

func parseCredentialAndSession(credentialIDRaw string, sessionIDRaw string) (uuid.UUID, uuid.UUID, error) {
	credentialID, err := uuid.Parse(credentialIDRaw)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("%w: invalid credential id", ErrValidation)
	}
	sessionID, err := uuid.Parse(sessionIDRaw)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("%w: invalid session id", ErrValidation)
	}
	return credentialID, sessionID, nil
}

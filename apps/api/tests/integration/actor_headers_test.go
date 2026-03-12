package integration_test

import (
	"net/http"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func setActorHeaders(req *http.Request, userID string, credentialID string, roles string) {
	if req == nil {
		return
	}

	if strings.TrimSpace(userID) != "" {
		req.Header.Set(identity.HeaderActorUserID, userID)
	} else {
		req.Header.Del(identity.HeaderActorUserID)
	}

	if strings.TrimSpace(credentialID) != "" {
		req.Header.Set(identity.HeaderActorCredentialID, credentialID)
	} else {
		req.Header.Del(identity.HeaderActorCredentialID)
	}

	if strings.TrimSpace(roles) != "" {
		req.Header.Set(identity.HeaderActorRoles, roles)
	} else {
		req.Header.Del(identity.HeaderActorRoles)
	}
}

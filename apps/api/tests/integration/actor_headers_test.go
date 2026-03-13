package integration_test

import (
	"net/http"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func setActorHeaders(req *http.Request, userID string, credentialID string, roles string) {
	if req == nil {
		return
	}

	req.Header.Del(identity.HeaderAuthorization)
	userID = strings.TrimSpace(userID)
	credentialID = strings.TrimSpace(credentialID)
	roles = strings.TrimSpace(roles)
	if userID == "" && credentialID == "" {
		return
	}

	roleItems := make([]string, 0)
	if roles != "" {
		for _, role := range strings.Split(roles, ",") {
			role = strings.TrimSpace(role)
			if role == "" {
				continue
			}
			roleItems = append(roleItems, role)
		}
	}

	token, err := identity.IssueAccessToken(identity.TokenClaims{
		UserID:       userID,
		CredentialID: credentialID,
		Roles:        roleItems,
		ExpiresAt:    time.Now().UTC().Add(time.Hour),
	})
	if err != nil {
		panic(err)
	}
	req.Header.Set(identity.HeaderAuthorization, "Bearer "+token)
}

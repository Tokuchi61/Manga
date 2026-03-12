package handler

import (
	"net/http"
	"strings"
	"time"

	accesscontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/contract"
	authcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/contract"
	usercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

func (h *HTTPHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Permission          string `json:"permission"`
		FeatureKey          string `json:"feature_key,omitempty"`
		ResourceOwnerUserID string `json:"resource_owner_user_id,omitempty"`
		ScopeKind           string `json:"scope_kind,omitempty"`
		ScopeSelector       string `json:"scope_selector,omitempty"`
		AudienceSelector    string `json:"audience_selector,omitempty"`
		AllowOverride       bool   `json:"allow_override,omitempty"`
	}
	if err := decodeJSON(r, &requestBody); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := identity.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_actor_user_id")
		return
	}

	contractInput := accesscontract.AuthorizationInput{
		UserID:              userID,
		Permission:          strings.TrimSpace(requestBody.Permission),
		FeatureKey:          strings.TrimSpace(requestBody.FeatureKey),
		ResourceOwnerUserID: strings.TrimSpace(requestBody.ResourceOwnerUserID),
		ScopeKind:           strings.TrimSpace(requestBody.ScopeKind),
		ScopeSelector:       strings.TrimSpace(requestBody.ScopeSelector),
		AudienceSelector:    strings.TrimSpace(requestBody.AudienceSelector),
		AllowOverride:       requestBody.AllowOverride,
		UserSignal: usercontract.AccessSignal{
			UserID:    userID,
			UpdatedAt: time.Now().UTC(),
		},
	}
	if credentialID, ok := identity.CredentialID(r.Context()); ok {
		contractInput.Identity = authcontract.VerifiedIdentity{
			CredentialID:    credentialID,
			EmailVerified:   true,
			AuthenticatedAt: time.Now().UTC(),
		}
	}

	res, err := h.service.EvaluateInput(r.Context(), contractInput)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

package service

import (
	"context"
	"strings"

	accesscontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
)

// EvaluateInput executes access evaluation using the cross-module authorization contract.
func (s *AccessService) EvaluateInput(ctx context.Context, input accesscontract.AuthorizationInput) (dto.EvaluateResponse, error) {
	request := dto.EvaluateRequest{
		UserID:              strings.TrimSpace(input.UserID),
		Permission:          strings.TrimSpace(input.Permission),
		FeatureKey:          strings.TrimSpace(input.FeatureKey),
		ResourceOwnerUserID: strings.TrimSpace(input.ResourceOwnerUserID),
		ScopeKind:           strings.TrimSpace(input.ScopeKind),
		ScopeSelector:       strings.TrimSpace(input.ScopeSelector),
		AudienceSelector:    strings.TrimSpace(input.AudienceSelector),
		AllowOverride:       input.AllowOverride,
	}

	if strings.TrimSpace(input.Identity.CredentialID) != "" {
		request.Identity = &dto.EvaluateIdentity{
			CredentialID:  strings.TrimSpace(input.Identity.CredentialID),
			SessionID:     strings.TrimSpace(input.Identity.SessionID),
			EmailVerified: input.Identity.EmailVerified,
			Suspended:     input.Identity.Suspended,
			Banned:        input.Identity.Banned,
		}
	}

	if strings.TrimSpace(input.UserSignal.UserID) != "" {
		request.UserSignal = &dto.EvaluateUserSignal{
			AccountState:                strings.TrimSpace(input.UserSignal.AccountState),
			ProfileVisibility:           strings.TrimSpace(input.UserSignal.ProfileVisibility),
			HistoryVisibilityPreference: strings.TrimSpace(input.UserSignal.HistoryVisibilityPreference),
			VIPActive:                   input.UserSignal.VIPActive,
			VIPFrozen:                   input.UserSignal.VIPFrozen,
		}
	}

	return s.Evaluate(ctx, request)
}

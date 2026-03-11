package catalog

// PolicyEffect defines canonical access and availability outcomes.
type PolicyEffect string

const (
	PolicyEffectAllow              PolicyEffect = "allow"
	PolicyEffectDeny               PolicyEffect = "deny"
	PolicyEffectDenySoft           PolicyEffect = "deny_soft"
	PolicyEffectRequireAuth        PolicyEffect = "require_auth"
	PolicyEffectRequireRole        PolicyEffect = "require_role"
	PolicyEffectRequireEntitlement PolicyEffect = "require_entitlement"
	PolicyEffectReadOnly           PolicyEffect = "read_only"
	PolicyEffectWriteOff           PolicyEffect = "write_off"
	PolicyEffectMask               PolicyEffect = "mask"
	PolicyEffectNeedsReview        PolicyEffect = "needs_review"
)

var AllPolicyEffects = []PolicyEffect{
	PolicyEffectAllow,
	PolicyEffectDeny,
	PolicyEffectDenySoft,
	PolicyEffectRequireAuth,
	PolicyEffectRequireRole,
	PolicyEffectRequireEntitlement,
	PolicyEffectReadOnly,
	PolicyEffectWriteOff,
	PolicyEffectMask,
	PolicyEffectNeedsReview,
}

func IsValidPolicyEffect(value PolicyEffect) bool {
	switch value {
	case PolicyEffectAllow,
		PolicyEffectDeny,
		PolicyEffectDenySoft,
		PolicyEffectRequireAuth,
		PolicyEffectRequireRole,
		PolicyEffectRequireEntitlement,
		PolicyEffectReadOnly,
		PolicyEffectWriteOff,
		PolicyEffectMask,
		PolicyEffectNeedsReview:
		return true
	default:
		return false
	}
}

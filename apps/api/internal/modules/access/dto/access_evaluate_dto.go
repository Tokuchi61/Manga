package dto

type EvaluateIdentity struct {
	CredentialID  string `json:"credential_id"`
	SessionID     string `json:"session_id"`
	EmailVerified bool   `json:"email_verified"`
	Suspended     bool   `json:"suspended"`
	Banned        bool   `json:"banned"`
}

type EvaluateUserSignal struct {
	AccountState                string `json:"account_state"`
	ProfileVisibility           string `json:"profile_visibility"`
	HistoryVisibilityPreference string `json:"history_visibility_preference"`
	VIPActive                   bool   `json:"vip_active"`
	VIPFrozen                   bool   `json:"vip_frozen"`
}

type EvaluateRequest struct {
	UserID              string              `json:"user_id,omitempty"`
	Permission          string              `json:"permission" validate:"required,min=3,max=128"`
	FeatureKey          string              `json:"feature_key,omitempty"`
	ResourceOwnerUserID string              `json:"resource_owner_user_id,omitempty"`
	ScopeKind           string              `json:"scope_kind,omitempty"`
	ScopeSelector       string              `json:"scope_selector,omitempty"`
	AudienceSelector    string              `json:"audience_selector,omitempty"`
	Identity            *EvaluateIdentity   `json:"identity,omitempty"`
	UserSignal          *EvaluateUserSignal `json:"user_signal,omitempty"`
	AllowOverride       bool                `json:"allow_override"`
}

type EvaluateResponse struct {
	Allowed       bool   `json:"allowed"`
	Effect        string `json:"effect"`
	ReasonCode    string `json:"reason_code"`
	Reason        string `json:"reason"`
	Permission    string `json:"permission"`
	PolicyVersion int    `json:"policy_version"`
	SubjectKind   string `json:"subject_kind"`
}

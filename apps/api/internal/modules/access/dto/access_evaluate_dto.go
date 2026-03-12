package dto

type EvaluateIdentity struct {
	CredentialID  string `json:"-"`
	SessionID     string `json:"-"`
	EmailVerified bool   `json:"-"`
	Suspended     bool   `json:"-"`
	Banned        bool   `json:"-"`
}

type EvaluateUserSignal struct {
	AccountState                string `json:"-"`
	ProfileVisibility           string `json:"-"`
	HistoryVisibilityPreference string `json:"-"`
	VIPActive                   bool   `json:"-"`
	VIPFrozen                   bool   `json:"-"`
}

type EvaluateRequest struct {
	UserID              string              `json:"-"`
	Permission          string              `json:"permission" validate:"required,min=3,max=128"`
	FeatureKey          string              `json:"feature_key,omitempty"`
	ResourceOwnerUserID string              `json:"resource_owner_user_id,omitempty"`
	ScopeKind           string              `json:"scope_kind,omitempty"`
	ScopeSelector       string              `json:"scope_selector,omitempty"`
	AudienceSelector    string              `json:"audience_selector,omitempty"`
	Identity            *EvaluateIdentity   `json:"-"`
	UserSignal          *EvaluateUserSignal `json:"-"`
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

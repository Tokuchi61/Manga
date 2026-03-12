package dto

import "time"

type CreatePolicyRuleRequest struct {
	Key              string `json:"key" validate:"required,min=3,max=256"`
	Effect           string `json:"effect" validate:"required"`
	AudienceKind     string `json:"audience_kind" validate:"required"`
	AudienceSelector string `json:"audience_selector" validate:"required"`
	ScopeKind        string `json:"scope_kind" validate:"required"`
	ScopeSelector    string `json:"scope_selector" validate:"required"`
	Active           *bool  `json:"active,omitempty"`
}

type CreatePolicyRuleResponse struct {
	PolicyRuleID     string    `json:"policy_rule_id"`
	Key              string    `json:"key"`
	Effect           string    `json:"effect"`
	AudienceKind     string    `json:"audience_kind"`
	AudienceSelector string    `json:"audience_selector"`
	ScopeKind        string    `json:"scope_kind"`
	ScopeSelector    string    `json:"scope_selector"`
	Active           bool      `json:"active"`
	Version          int       `json:"version"`
	CreatedAt        time.Time `json:"created_at"`
}

package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type PolicyEffect string

const (
	PolicyEffectAllow         PolicyEffect = "allow"
	PolicyEffectDeny          PolicyEffect = "deny"
	PolicyEffectEmergencyDeny PolicyEffect = "emergency_deny"
)

func IsValidPolicyEffect(value PolicyEffect) bool {
	switch value {
	case PolicyEffectAllow, PolicyEffectDeny, PolicyEffectEmergencyDeny:
		return true
	default:
		return false
	}
}

type PolicyRule struct {
	ID               uuid.UUID
	Key              string
	Effect           PolicyEffect
	AudienceKind     string
	AudienceSelector string
	ScopeKind        string
	ScopeSelector    string
	Active           bool
	Version          int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (p PolicyRule) ConflictKey() string {
	return strings.Join([]string{
		strings.ToLower(strings.TrimSpace(p.Key)),
		strings.ToLower(strings.TrimSpace(p.AudienceKind)),
		strings.ToLower(strings.TrimSpace(p.AudienceSelector)),
		strings.ToLower(strings.TrimSpace(p.ScopeKind)),
		strings.ToLower(strings.TrimSpace(p.ScopeSelector)),
	}, "|")
}

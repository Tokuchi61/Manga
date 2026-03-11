package settings

import (
	"regexp"
	"strings"
)

// Audience values used by runtime settings.
type Audience string

const (
	AudienceAll                 Audience = "all"
	AudienceGuest               Audience = "guest"
	AudienceAuthenticated       Audience = "authenticated"
	AudienceAuthenticatedNonVIP Audience = "authenticated_non_vip"
	AudienceVIP                 Audience = "vip"
)

var AllAudiences = []Audience{
	AudienceAll,
	AudienceGuest,
	AudienceAuthenticated,
	AudienceAuthenticatedNonVIP,
	AudienceVIP,
}

// ScopeKind values used by runtime settings.
type ScopeKind string

const (
	ScopeSite            ScopeKind = "site"
	ScopeModule          ScopeKind = "module"
	ScopeFeature         ScopeKind = "feature"
	ScopeResourceContext ScopeKind = "resource/context"
)

var AllScopeKinds = []ScopeKind{
	ScopeSite,
	ScopeModule,
	ScopeFeature,
	ScopeResourceContext,
}

// DisabledBehavior captures canonical shutdown/degrade behaviors.
type DisabledBehavior string

const (
	DisabledBehaviorVisibilityOff  DisabledBehavior = "visibility_off"
	DisabledBehaviorReadOnly       DisabledBehavior = "read_only"
	DisabledBehaviorWriteOff       DisabledBehavior = "write_off"
	DisabledBehaviorIntakePause    DisabledBehavior = "intake_pause"
	DisabledBehaviorReadOnlyIntake DisabledBehavior = "read_only_intake"
	DisabledBehaviorAttachmentOff  DisabledBehavior = "attachment_off"
	DisabledBehaviorPreviewOff     DisabledBehavior = "preview_off"
	DisabledBehaviorBenefitPause   DisabledBehavior = "benefit_pause"
)

var AllDisabledBehaviors = []DisabledBehavior{
	DisabledBehaviorVisibilityOff,
	DisabledBehaviorReadOnly,
	DisabledBehaviorWriteOff,
	DisabledBehaviorIntakePause,
	DisabledBehaviorReadOnlyIntake,
	DisabledBehaviorAttachmentOff,
	DisabledBehaviorPreviewOff,
	DisabledBehaviorBenefitPause,
}

// ErrorResponsePolicy captures canonical external response modes.
type ErrorResponsePolicy string

const (
	ErrorResponseNotFound               ErrorResponsePolicy = "not_found"
	ErrorResponseForbidden              ErrorResponsePolicy = "forbidden"
	ErrorResponseRateLimited            ErrorResponsePolicy = "rate_limited"
	ErrorResponseValidationError        ErrorResponsePolicy = "validation_error"
	ErrorResponseTemporarilyUnavailable ErrorResponsePolicy = "temporarily_unavailable"
)

var AllErrorResponsePolicies = []ErrorResponsePolicy{
	ErrorResponseNotFound,
	ErrorResponseForbidden,
	ErrorResponseRateLimited,
	ErrorResponseValidationError,
	ErrorResponseTemporarilyUnavailable,
}

// EntitlementImpactPolicy captures canonical entitlement side effects.
type EntitlementImpactPolicy string

const (
	EntitlementImpactNone                  EntitlementImpactPolicy = "none"
	EntitlementImpactFreezeOnSystemDisable EntitlementImpactPolicy = "freeze_on_system_disable"
)

var AllEntitlementImpactPolicies = []EntitlementImpactPolicy{
	EntitlementImpactNone,
	EntitlementImpactFreezeOnSystemDisable,
}

// ApplyMode captures settings application timing.
type ApplyMode string

const (
	ApplyModeImmediate    ApplyMode = "immediate"
	ApplyModeCacheRefresh ApplyMode = "cache_refresh"
	ApplyModeScheduled    ApplyMode = "scheduled"
)

// CacheStrategy captures settings cache behavior.
type CacheStrategy string

const (
	CacheStrategyNone             CacheStrategy = "none"
	CacheStrategyTTL              CacheStrategy = "ttl"
	CacheStrategyManualInvalidate CacheStrategy = "manual_invalidate"
)

// ScheduleSupport captures scheduling capabilities per setting.
type ScheduleSupport string

const (
	ScheduleSupportNone       ScheduleSupport = "none"
	ScheduleSupportStartAt    ScheduleSupport = "start_at"
	ScheduleSupportTimeWindow ScheduleSupport = "time_window"
)

var AccessInterpretationOrder = []string{
	"global",
	"module",
	"surface",
	"audience",
	"entitlement",
	"action",
	"rate_limit",
}

// KillSwitchLevel captures the shutdown granularity.
type KillSwitchLevel string

const (
	KillSwitchGlobal                  KillSwitchLevel = "global"
	KillSwitchModuleLevel             KillSwitchLevel = "module-level"
	KillSwitchSurfaceLevel            KillSwitchLevel = "surface-level"
	KillSwitchActionLevel             KillSwitchLevel = "action-level"
	KillSwitchIntakeOnly              KillSwitchLevel = "intake-only"
	KillSwitchWriteOnly               KillSwitchLevel = "write-only"
	KillSwitchExternalIntegrationOnly KillSwitchLevel = "external-integration-only"
)

var AllKillSwitchLevels = []KillSwitchLevel{
	KillSwitchGlobal,
	KillSwitchModuleLevel,
	KillSwitchSurfaceLevel,
	KillSwitchActionLevel,
	KillSwitchIntakeOnly,
	KillSwitchWriteOnly,
	KillSwitchExternalIntegrationOnly,
}

var SettingsRecordFields = []string{
	"key",
	"description",
	"category",
	"owner_module",
	"consumer_layer",
	"value_type",
	"default_value",
	"allowed_range_or_enum",
	"scope_kind",
	"scope_selector",
	"audience_kind",
	"audience_selector",
	"sensitive",
	"apply_mode",
	"cache_strategy",
	"schedule_support",
	"audit_required",
	"affected_surfaces",
	"disabled_behavior",
	"error_response_policy",
	"entitlement_impact_policy",
	"status",
	"notes",
}

var FeatureToggleKeyPattern = regexp.MustCompile(`^feature\.[a-z][a-z0-9_]*\.[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)*\.enabled$`)
var ModuleSurfaceMetricPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*\.[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)*$`)
var SiteMetricPattern = regexp.MustCompile(`^site\.[a-z][a-z0-9_]*\.[a-z][a-z0-9_]*$`)

func IsFeatureToggleKey(key string) bool {
	return FeatureToggleKeyPattern.MatchString(key)
}

func IsModuleSurfaceMetricKey(key string) bool {
	if strings.HasPrefix(key, "site.") {
		return false
	}
	return ModuleSurfaceMetricPattern.MatchString(key)
}

func IsSiteMetricKey(key string) bool {
	return SiteMetricPattern.MatchString(key)
}

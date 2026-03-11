package settings

// SettingRecord defines canonical runtime settings inventory fields.
type SettingRecord struct {
	Key                     string
	Description             string
	Category                string
	OwnerModule             string
	ConsumerLayer           string
	ValueType               string
	DefaultValue            string
	AllowedRangeOrEnum      string
	ScopeKind               ScopeKind
	ScopeSelector           string
	AudienceKind            string
	AudienceSelector        string
	Sensitive               bool
	ApplyMode               ApplyMode
	CacheStrategy           CacheStrategy
	ScheduleSupport         ScheduleSupport
	AuditRequired           bool
	AffectedSurfaces        []string
	DisabledBehavior        DisabledBehavior
	ErrorResponsePolicy     ErrorResponsePolicy
	EntitlementImpactPolicy EntitlementImpactPolicy
	Status                  string
	Notes                   string
}

func (r SettingRecord) HasRequiredIdentity() bool {
	return r.Key != "" && r.OwnerModule != "" && r.ConsumerLayer != ""
}

func (r SettingRecord) IsStatusKnown() bool {
	switch r.Status {
	case "active", "planned", "deprecated":
		return true
	default:
		return false
	}
}

package policy

import "strings"

// ModuleStatus defines canonical module lifecycle statuses.
type ModuleStatus string

const (
	ModuleStatusPlanned    ModuleStatus = "planned"
	ModuleStatusActive     ModuleStatus = "active"
	ModuleStatusDeprecated ModuleStatus = "deprecated"
	ModuleStatusArchived   ModuleStatus = "archived"
)

var AllModuleStatuses = []ModuleStatus{
	ModuleStatusPlanned,
	ModuleStatusActive,
	ModuleStatusDeprecated,
	ModuleStatusArchived,
}

// SuggestedDomainGroups are optional grouping hints for scaling.
var SuggestedDomainGroups = []string{
	"identity",
	"content",
	"community",
	"operations",
	"engagement",
	"commerce",
	"gameplay",
}

var ModuleInventoryRequiredFields = []string{
	"name",
	"domain_group_optional",
	"description",
	"status",
	"doc_path",
}

// ModuleInventoryRecord captures canonical module inventory fields.
type ModuleInventoryRecord struct {
	Name        string
	DomainGroup string
	Description string
	Status      ModuleStatus
	DocPath     string
}

func IsValidModuleStatus(value ModuleStatus) bool {
	switch value {
	case ModuleStatusPlanned,
		ModuleStatusActive,
		ModuleStatusDeprecated,
		ModuleStatusArchived:
		return true
	default:
		return false
	}
}

func IsSuggestedDomainGroup(value string) bool {
	normalized := strings.TrimSpace(strings.ToLower(value))
	for _, item := range SuggestedDomainGroups {
		if item == normalized {
			return true
		}
	}
	return false
}

func BuildModuleRootPath(moduleName string, domainGroup string) string {
	name := strings.TrimSpace(moduleName)
	group := strings.TrimSpace(domainGroup)
	if group == "" {
		return "apps/api/internal/modules/" + name
	}
	return "apps/api/internal/modules/" + group + "/" + name
}

func (r ModuleInventoryRecord) HasRequiredFields() bool {
	if strings.TrimSpace(r.Name) == "" {
		return false
	}
	if strings.TrimSpace(r.Description) == "" {
		return false
	}
	if strings.TrimSpace(r.DocPath) == "" {
		return false
	}
	return IsValidModuleStatus(r.Status)
}

func (r ModuleInventoryRecord) IsDocPathUnderDocsModules() bool {
	normalized := strings.TrimSpace(strings.ToLower(strings.ReplaceAll(r.DocPath, `\\`, "/")))
	if normalized == "docs/modules.md" {
		return true
	}
	return strings.HasPrefix(normalized, "docs/modules/")
}

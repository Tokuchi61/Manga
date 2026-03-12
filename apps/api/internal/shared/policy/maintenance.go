package policy

// MaintenanceRefactorRule captures stage-3 expansion guardrails.
type MaintenanceRefactorRule string

const (
	MaintenanceRuleNoUnnecessaryRefactor         MaintenanceRefactorRule = "no_unnecessary_refactor"
	MaintenanceRuleKeepChangesLogicalAndSmall    MaintenanceRefactorRule = "changes_must_be_logical_small_and_revertible"
	MaintenanceRuleUpdateDocsWithBoundaryChange  MaintenanceRefactorRule = "update_docs_when_architecture_or_boundaries_change"
	MaintenanceRuleAvoidCrossModuleTightCoupling MaintenanceRefactorRule = "avoid_cross_module_tight_coupling"
	MaintenanceRuleKeepOwnerBoundariesClear      MaintenanceRefactorRule = "keep_owner_boundaries_clear"
)

var MaintenanceRefactorRules = []MaintenanceRefactorRule{
	MaintenanceRuleNoUnnecessaryRefactor,
	MaintenanceRuleKeepChangesLogicalAndSmall,
	MaintenanceRuleUpdateDocsWithBoundaryChange,
	MaintenanceRuleAvoidCrossModuleTightCoupling,
	MaintenanceRuleKeepOwnerBoundariesClear,
}

var StageThreeOperationalChecklist = []string{
	"stage_tests_must_pass",
	"docker_build_and_run_must_pass",
	"versioning_must_be_applied_before_push",
	"changelog_and_upgrade_docs_must_be_updated",
}

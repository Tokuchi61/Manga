package catalog

// VisibilityState defines canonical visibility states across modules.
type VisibilityState string

const (
	VisibilityStatePublic   VisibilityState = "public"
	VisibilityStateLimited  VisibilityState = "limited"
	VisibilityStatePrivate  VisibilityState = "private"
	VisibilityStateHidden   VisibilityState = "hidden"
	VisibilityStateRemoved  VisibilityState = "removed"
	VisibilityStateArchived VisibilityState = "archived"
)

var AllVisibilityStates = []VisibilityState{
	VisibilityStatePublic,
	VisibilityStateLimited,
	VisibilityStatePrivate,
	VisibilityStateHidden,
	VisibilityStateRemoved,
	VisibilityStateArchived,
}

func IsValidVisibilityState(value VisibilityState) bool {
	switch value {
	case VisibilityStatePublic,
		VisibilityStateLimited,
		VisibilityStatePrivate,
		VisibilityStateHidden,
		VisibilityStateRemoved,
		VisibilityStateArchived:
		return true
	default:
		return false
	}
}

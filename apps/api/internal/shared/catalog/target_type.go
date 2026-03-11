package catalog

// TargetType defines canonical cross-module target references.
type TargetType string

const (
	TargetTypeManga   TargetType = "manga"
	TargetTypeChapter TargetType = "chapter"
	TargetTypeComment TargetType = "comment"
	TargetTypeSocial  TargetType = "social"
)

var AllTargetTypes = []TargetType{
	TargetTypeManga,
	TargetTypeChapter,
	TargetTypeComment,
	TargetTypeSocial,
}

func IsValidTargetType(value TargetType) bool {
	switch value {
	case TargetTypeManga, TargetTypeChapter, TargetTypeComment, TargetTypeSocial:
		return true
	default:
		return false
	}
}

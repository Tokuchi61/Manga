package contract

import "time"

// ChapterAccessDefaults exposes manga-owned chapter default access metadata.
type ChapterAccessDefaults struct {
	MangaID            string
	ReadAccessLevel    string
	EarlyAccessEnabled bool
	EarlyAccessLevel   string
	UpdatedAt          time.Time
}

// TargetReference represents manga target relation for comment/support surfaces.
type TargetReference struct {
	TargetType string
	TargetID   string
}

const TargetTypeManga = "manga"

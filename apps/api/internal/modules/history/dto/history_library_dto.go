package dto

// ListLibraryRequest resolves own library surface.
type ListLibraryRequest struct {
	UserID       string `json:"-" validate:"required,uuid4"`
	Status       string `json:"-" validate:"omitempty,oneof=in_progress completed dropped"`
	Bookmarked   *bool  `json:"-"`
	Favorited    *bool  `json:"-"`
	SharedOnly   bool   `json:"-"`
	Limit        int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset       int    `json:"-" validate:"omitempty,min=0"`
	SortBy       string `json:"-" validate:"omitempty,oneof=newest oldest"`
}

// ListPublicLibraryRequest resolves public library surface.
type ListPublicLibraryRequest struct {
	OwnerUserID string `json:"-" validate:"required,uuid4"`
	Limit       int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset      int    `json:"-" validate:"omitempty,min=0"`
	SortBy      string `json:"-" validate:"omitempty,oneof=newest oldest"`
}

// ListLibraryResponse wraps library payload.
type ListLibraryResponse struct {
	Items []HistoryEntryResponse `json:"items"`
	Count int                    `json:"count"`
}

// UpdateBookmarkRequest updates bookmark/favorite flags.
type UpdateBookmarkRequest struct {
	UserID     string `json:"-" validate:"required,uuid4"`
	MangaID    string `json:"-" validate:"required,uuid4"`
	Bookmarked bool   `json:"bookmarked"`
	Favorited  bool   `json:"favorited"`
}

// UpdateShareRequest updates entry-level public share metadata.
type UpdateShareRequest struct {
	UserID      string `json:"-" validate:"required,uuid4"`
	MangaID     string `json:"-" validate:"required,uuid4"`
	SharePublic bool   `json:"share_public"`
}

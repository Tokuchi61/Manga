package chapter

import "errors"

var (
	// ErrUnavailable can be used by future runtime gate checks.
	ErrUnavailable = errors.New("chapter_module_unavailable")
)

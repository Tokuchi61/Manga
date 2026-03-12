package manga

import "errors"

var (
	// ErrUnavailable can be used by future runtime gate checks.
	ErrUnavailable = errors.New("manga_module_unavailable")
)

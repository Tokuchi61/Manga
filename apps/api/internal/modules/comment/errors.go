package comment

import "errors"

var (
	// ErrUnavailable can be used by future runtime gate checks.
	ErrUnavailable = errors.New("comment_module_unavailable")
)

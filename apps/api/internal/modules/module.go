package modules

import "github.com/go-chi/chi/v5"

// Module defines the minimal runtime contract for leaf modules.
type Module interface {
	Name() string
	RegisterRoutes(router chi.Router)
}

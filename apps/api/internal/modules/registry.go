package modules

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
)

var moduleNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Registry keeps module wiring centralized in the app bootstrap layer.
type Registry struct {
	items []Module
}

func NewRegistry(modules ...Module) (*Registry, error) {
	seen := make(map[string]struct{}, len(modules))
	items := make([]Module, 0, len(modules))

	for _, module := range modules {
		if module == nil {
			return nil, fmt.Errorf("module cannot be nil")
		}

		name := strings.TrimSpace(module.Name())
		if name == "" {
			return nil, fmt.Errorf("module name cannot be empty")
		}
		if !moduleNamePattern.MatchString(name) {
			return nil, fmt.Errorf("invalid module name %q", name)
		}
		if _, exists := seen[name]; exists {
			return nil, fmt.Errorf("duplicate module name %q", name)
		}

		seen[name] = struct{}{}
		items = append(items, module)
	}

	return &Registry{items: items}, nil
}

func EmptyRegistry() *Registry {
	return &Registry{items: []Module{}}
}

func (r *Registry) Names() []string {
	if r == nil {
		return nil
	}

	names := make([]string, 0, len(r.items))
	for _, module := range r.items {
		names = append(names, strings.TrimSpace(module.Name()))
	}
	sort.Strings(names)

	return names
}

func (r *Registry) Modules() []Module {
	if r == nil {
		return nil
	}
	items := make([]Module, 0, len(r.items))
	items = append(items, r.items...)
	return items
}

func (r *Registry) MountRoutes(router chi.Router) {
	if r == nil {
		return
	}
	for _, module := range r.items {
		module.RegisterRoutes(router)
	}
}

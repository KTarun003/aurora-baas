package codegen

import (
	"fmt"
	"sync"
)

// Generator is the interface that all code generators must implement
type Generator interface {
	// Generate produces a complete project from the given context
	Generate(ctx *GeneratorContext) (*GeneratedProject, error)

	// Name returns the name of this generator (e.g., "typescript", "python")
	Name() string
}

// Registry manages available generators
type Registry struct {
	mu         sync.RWMutex
	generators map[string]Generator
}

// NewRegistry creates a new generator registry
func NewRegistry() *Registry {
	return &Registry{
		generators: make(map[string]Generator),
	}
}

// Register adds a generator to the registry
func (r *Registry) Register(gen Generator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.generators[gen.Name()] = gen
}

// Get retrieves a generator by name
func (r *Registry) Get(name string) (Generator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	gen, ok := r.generators[name]
	return gen, ok
}

// List returns all registered generator names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.generators))
	for name := range r.generators {
		names = append(names, name)
	}
	return names
}

// Generate is a convenience method to get a generator and generate in one call
func (r *Registry) Generate(generatorName string, ctx *GeneratorContext) (*GeneratedProject, error) {
	gen, ok := r.Get(generatorName)
	if !ok {
		return nil, fmt.Errorf("generator not found: %s", generatorName)
	}
	return gen.Generate(ctx)
}

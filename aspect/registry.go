// Package aspect - registry manages function registration and advice attachment
package aspect

import (
	"fmt"
	"sync"
)

// -------------------------------------------- Global Variables --------------------------------------------

var (
	// defaultRegistry is the global default registry used by the fluent API
	defaultRegistry *Registry
	defaultRegOnce  sync.Once
)

// -------------------------------------------- Types --------------------------------------------

// Registry stores function references and their associated advice chains.
type Registry struct {
	mu      sync.RWMutex
	entries map[FuncKey]*AdviceChain
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[FuncKey]*AdviceChain),
	}
}

// DefaultRegistry returns the global default registry.
func DefaultRegistry() *Registry {
	defaultRegOnce.Do(func() {
		defaultRegistry = NewRegistry()
	})
	return defaultRegistry
}

// -------------------------------------------- Public Functions --------------------------------------------

// Register registers a function with the given name.
// Returns error if the function is already registered.
func (registry *Registry) Register(name FuncKey) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if name == "" {
		return fmt.Errorf("function name cannot be empty")
	}

	if _, exists := registry.entries[name]; exists {
		return fmt.Errorf("function '%s' is already registered", name)
	}

	registry.entries[name] = NewAdviceChain()
	return nil
}

// RegisterOrGet registers a function if not already registered, otherwise returns existing chain.
// Always returns the advice chain and never errors.
func (registry *Registry) RegisterOrGet(name FuncKey) *AdviceChain {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if name == "" {
		panic("function name cannot be empty")
	}

	if chain, exists := registry.entries[name]; exists {
		return chain
	}

	chain := NewAdviceChain()
	registry.entries[name] = chain
	return chain
}

// MustRegister registers a function and panics on error.
// Useful for initialization code where registration must succeed.
func (registry *Registry) MustRegister(name FuncKey) {
	if err := registry.Register(name); err != nil {
		panic(err)
	}
}

// AddAdvice adds an advice to the specified function.
// Returns error if the function is not registered.
func (registry *Registry) AddAdvice(funcKey FuncKey, advice Advice) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if funcKey == "" {
		return fmt.Errorf("function name cannot be empty")
	}

	chain, exists := registry.entries[funcKey]
	if !exists {
		return fmt.Errorf("function '%s' is not registered", funcKey)
	}

	chain.Add(advice)
	return nil
}

// MustAddAdvice adds advice and panics on error.
// Useful for initialization code where advice addition must succeed.
func (registry *Registry) MustAddAdvice(funcKey FuncKey, advice Advice) {
	if err := registry.AddAdvice(funcKey, advice); err != nil {
		panic(err)
	}
}

// GetAdviceChain retrieves the advice chain for a function.
// Returns error if the function is not registered.
func (registry *Registry) GetAdviceChain(funcKey FuncKey) (*AdviceChain, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	if funcKey == "" {
		return nil, fmt.Errorf("function name cannot be empty")
	}

	chain, exists := registry.entries[funcKey]
	if !exists {
		return nil, fmt.Errorf("function '%s' is not registered", funcKey)
	}

	return chain, nil
}

// IsRegistered checks if a function is registered.
func (registry *Registry) IsRegistered(name FuncKey) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	_, exists := registry.entries[name]
	return exists
}

// Unregister removes a function from the registry.
// Does nothing if the function is not registered.
func (registry *Registry) Unregister(name FuncKey) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	delete(registry.entries, name)
}

// ListRegistered returns all registered function names.
func (registry *Registry) ListRegistered() []FuncKey {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	names := make([]FuncKey, 0, len(registry.entries))
	for name := range registry.entries {
		names = append(names, name)
	}
	return names
}

// Clear removes all registered functions from the registry.
func (registry *Registry) Clear() {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.entries = make(map[FuncKey]*AdviceChain)
}

// Count returns the number of registered functions.
func (registry *Registry) Count() int {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	return len(registry.entries)
}

// GetAdviceCount returns the total number of advice for a function.
// Returns 0 if the function is not registered.
func (registry *Registry) GetAdviceCount(funcKey FuncKey) int {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	chain, exists := registry.entries[funcKey]
	if !exists {
		return 0
	}

	return chain.Count()
}

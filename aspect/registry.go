// Package aspect - registry manages function registration and advice attachment
package aspect

import (
	"fmt"
	"sync"
)

// -------------------------------------------- Types --------------------------------------------

// Registry stores function references and their associated advice chains.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]*AdviceChain
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]*AdviceChain),
	}
}

// -------------------------------------------- Public Functions --------------------------------------------

// Register registers a function with the given name.
// Returns error if the function is already registered.
func (registry *Registry) Register(name string) error {
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
func (registry *Registry) RegisterOrGet(name string) *AdviceChain {
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
func (registry *Registry) MustRegister(name string) {
	if err := registry.Register(name); err != nil {
		panic(err)
	}
}

// AddAdvice adds an advice to the specified function.
// Returns error if the function is not registered.
func (registry *Registry) AddAdvice(functionName string, advice Advice) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if functionName == "" {
		return fmt.Errorf("function name cannot be empty")
	}

	chain, exists := registry.entries[functionName]
	if !exists {
		return fmt.Errorf("function '%s' is not registered", functionName)
	}

	chain.Add(advice)
	return nil
}

// MustAddAdvice adds advice and panics on error.
// Useful for initialization code where advice addition must succeed.
func (registry *Registry) MustAddAdvice(functionName string, advice Advice) {
	if err := registry.AddAdvice(functionName, advice); err != nil {
		panic(err)
	}
}

// GetAdviceChain retrieves the advice chain for a function.
// Returns error if the function is not registered.
func (registry *Registry) GetAdviceChain(functionName string) (*AdviceChain, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	if functionName == "" {
		return nil, fmt.Errorf("function name cannot be empty")
	}

	chain, exists := registry.entries[functionName]
	if !exists {
		return nil, fmt.Errorf("function '%s' is not registered", functionName)
	}

	return chain, nil
}

// IsRegistered checks if a function is registered.
func (registry *Registry) IsRegistered(name string) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	_, exists := registry.entries[name]
	return exists
}

// Unregister removes a function from the registry.
// Does nothing if the function is not registered.
func (registry *Registry) Unregister(name string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	delete(registry.entries, name)
}

// ListRegistered returns all registered function names.
func (registry *Registry) ListRegistered() []string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	names := make([]string, 0, len(registry.entries))
	for name := range registry.entries {
		names = append(names, name)
	}
	return names
}

// Clear removes all registered functions from the registry.
func (registry *Registry) Clear() {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.entries = make(map[string]*AdviceChain)
}

// Count returns the number of registered functions.
func (registry *Registry) Count() int {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	return len(registry.entries)
}

// GetAdviceCount returns the total number of advice for a function.
// Returns 0 if the function is not registered.
func (registry *Registry) GetAdviceCount(functionName string) int {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	chain, exists := registry.entries[functionName]
	if !exists {
		return 0
	}

	return chain.Count()
}

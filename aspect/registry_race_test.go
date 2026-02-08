// Package aspect - registry_race_test targets concurrent access to the registry.
package aspect

import (
	"sync"
	"testing"
	"time"
)

func TestRegistryRace_RegisterAndLookupConcurrent(t *testing.T) {
	registry := NewRegistry()
	registry.Clear() // Ensure clean state
	defer registry.Clear()

	const numGoroutines = 10
	var wg sync.WaitGroup

	// Launch goroutines to register and then immediately try to wrap/check existence
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			funcName := FuncKey("ConcurrentFunc_" + string(rune('A'+id)))

			_ = registry.Register(funcName)
			// Brief pause to allow other goroutines to potentially interfere
			time.Sleep(time.Millisecond * 10)
			if !registry.IsRegistered(funcName) {
				t.Errorf("Function %s should be registered after calling Register", funcName)
			}

			// Attempt to add an advice (this requires the function to be registered first)
			err := registry.AddAdvice(funcName, Advice{
				Type:     Before,
				Priority: id * 10, // Unique priority
				Handler:  func(ctx *Context) error { return nil },
			})
			if err != nil {
				t.Errorf("Failed to add advice to %s: %v", funcName, err)
			}

			// Attempt to wrap the function concurrently
			_ = Wrap0(registry, funcName, func() {})
		}(i)
	}

	wg.Wait()
}

func TestRegistryRace_RegisterAndUnregisterConcurrent(t *testing.T) {
	registry := NewRegistry()
	registry.Clear() // Ensure clean state
	defer registry.Clear()

	const numGoroutines = 5
	var wg sync.WaitGroup

	// Pre-register a set of functions - twice as many as goroutines
	funcNames := make([]string, numGoroutines*2)
	for i := 0; i < len(funcNames); i++ {
		name := "TestFunc_" + string(rune('A'+i))
		funcNames[i] = name
		_ = registry.Register(FuncKey(name))
	}

	// Half the goroutines register new functions, half unregister existing ones
	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)

		// Goroutine 1: Register a new function
		go func(id int) {
			defer wg.Done()

			newName := FuncKey("NewlyAddedFunc_" + string(rune('Z'-id)))
			_ = registry.Register(newName)
			_ = registry.AddAdvice(newName, Advice{Type: Before, Handler: func(*Context) error { return nil }, Priority: 100})
		}(i)

		// Goroutine 2: Unregister an existing function
		go func(id int) {
			defer wg.Done()

			nameToUnregister := FuncKey(funcNames[id])
			registry.Unregister(nameToUnregister)
		}(i)
	}

	wg.Wait()

	// Verify state consistency after operations
	// We should only check the functions that were supposed to be unregistered (indices 0-4)
	for i := 0; i < numGoroutines; i++ {
		name := FuncKey(funcNames[i])
		if registry.IsRegistered(name) {
			t.Errorf("Function %s should have been unregistered", name)
		}
	}

	// The remaining functions (indices 5-9) should still be registered
	for i := numGoroutines; i < len(funcNames); i++ {
		name := FuncKey(funcNames[i])
		if !registry.IsRegistered(name) {
			t.Errorf("Function %s should still be registered", name)
		}
	}
	// Note: Checking newly added functions might be racy itself, so we focus on ensuring no panic/crash.
}

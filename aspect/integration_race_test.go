// Package aspect - integration_race_test runs the complete AOP workflow under concurrent load.
package aspect

import (
	"errors"
	"sync"
	"testing"
)

func TestIntegrationRace_FullWorkflow(t *testing.T) {
	registry := NewRegistry()
	registry.Clear() // Ensure clean state
	defer registry.Clear()

	funcName := FuncKey("RaceTestTargetFunc")
	registry.Register(funcName)

	// Add various types of advice to trigger the full execution path
	registry.MustAddAdvice(funcName, Advice{
		Type:     Before,
		Priority: 10,
		Handler: func(ctx *Context) error {
			ctx.Metadata["before_executed"] = true
			return nil
		},
	})

	registry.MustAddAdvice(funcName, Advice{
		Type:     Around,
		Priority: 50,
		Handler: func(ctx *Context) error {
			// Simulate some work in Around advice
			ctx.Metadata["around_executed"] = true
			return nil // Do not set ctx.Skipped, let target run
		},
	})

	registry.MustAddAdvice(funcName, Advice{
		Type:     AfterReturning,
		Priority: 100,
		Handler: func(ctx *Context) error {
			ctx.Metadata["after_returning_executed"] = true
			return nil
		},
	})

	registry.MustAddAdvice(funcName, Advice{
		Type:     After,
		Priority: 200,
		Handler: func(ctx *Context) error {
			ctx.Metadata["after_executed"] = true
			return nil
		},
	})

	// Target function to be wrapped
	targetFunc := func(input int) (int, error) {
		// Simulate some work
		result := input * 2
		if input < 0 {
			return 0, nil // Or maybe return an error, adjust test expectation
		}
		return result, nil
	}

	wrappedFunc := Wrap1RE(registry, funcName, targetFunc)

	var wg sync.WaitGroup
	const numGoroutines = 20
	const callsPerGoroutine = 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				inputVal := goroutineID*100 + j
				result, err := wrappedFunc(inputVal)

				expectedResult := inputVal * 2
				if result != expectedResult || err != nil {
					t.Errorf("Goroutine %d, Call %d: Expected (%d, nil), got (%d, %v)", goroutineID, j, expectedResult, result, err)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestIntegrationRace_FullWorkflowWithError(t *testing.T) {
	registry := NewRegistry()
	registry.Clear() // Ensure clean state
	defer registry.Clear()

	funcName := FuncKey("RaceTestTargetFuncWithError")
	registry.Register(funcName)

	// Add advice that interacts with the error field
	registry.MustAddAdvice(funcName, Advice{
		Type:     Before,
		Priority: 10,
		Handler: func(ctx *Context) error {
			ctx.Metadata["before_executed"] = true
			return nil
		},
	})

	registry.MustAddAdvice(funcName, Advice{
		Type:     AfterReturning, // This should NOT run if target returns an error
		Priority: 100,
		Handler: func(ctx *Context) error {
			if ctx.Error == nil {
				ctx.Metadata["after_returning_executed"] = true
			}
			return nil
		},
	})

	registry.MustAddAdvice(funcName, Advice{
		Type:     After, // This SHOULD run regardless of error
		Priority: 200,
		Handler: func(ctx *Context) error {
			ctx.Metadata["after_executed"] = true
			return nil
		},
	})

	// Target function that can return an error
	targetFunc := func(input int) (int, error) {
		if input%2 == 0 {
			return 0, nil // Success case
		}
		return 0, errors.New("simulated error") // Error case
	}

	wrappedFunc := Wrap1RE(registry, funcName, targetFunc)

	var wg sync.WaitGroup
	const numGoroutines = 10
	const callsPerGoroutine = 20

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				inputVal := goroutineID*100 + j
				_, err := wrappedFunc(inputVal)

				// Just check for successful completion without panicking due to race conditions
				_ = err
			}
		}(i)
	}

	wg.Wait()
}

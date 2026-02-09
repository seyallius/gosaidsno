// Package aspect - context_test validates context propagation functionality
package aspect

import (
	"context"
	"testing"
	"time"
)

// TestContextCancellation verifies that context cancellation propagates through advice
func TestContextCancellation(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()
	defer registry.Clear()

	var executionOrder []string

	// Register function
	registry.MustRegister("TestContextCancellation")

	// Add multiple before advice to test cancellation propagation
	registry.MustAddAdvice("TestContextCancellation", Advice{
		Type:     Before,
		Priority: 100,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, "before-1")
			// Check if context is cancelled
			select {
			case <-c.Context().Done():
				return c.Context().Err()
			default:
				// Context not cancelled, continue
			}
			return nil
		},
	})

	registry.MustAddAdvice("TestContextCancellation", Advice{
		Type:     Before,
		Priority: 50,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, "before-2")
			// Check if context is cancelled
			select {
			case <-c.Context().Done():
				return c.Context().Err()
			default:
				// Context not cancelled, continue
			}
			return nil
		},
	})

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the context
	cancel()

	// Create and wrap function
	testFn := func(ctx context.Context) string {
		executionOrder = append(executionOrder, "target")
		return "result"
	}
	wrappedFn := Wrap0RCtx[string](registry, "TestContextCancellation", testFn)

	// Execute the wrapped function with cancelled context
	defer func() {
		if r := recover(); r != nil {
			// Expected: context cancellation should cause panic
		}
	}()

	result := wrappedFn(ctx)

	// The function should not reach the target due to context cancellation
	// So result should be the default zero value
	if result != "" {
		t.Errorf("expected empty result due to context cancellation, got '%s'", result)
	}

	// Check execution order - should only execute first advice before cancellation is detected
	if len(executionOrder) == 0 {
		t.Fatal("no advice executed")
	}

	// Should have executed at least the first advice
	if executionOrder[0] != "before-1" {
		t.Errorf("expected first advice to execute, got %v", executionOrder)
	}
}

// TestContextDeadline verifies that context deadline propagates through advice
func TestContextDeadline(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()
	defer registry.Clear()

	var executionOrder []string

	// Register function
	registry.MustRegister("TestContextDeadline")

	// Add advice that checks context
	registry.MustAddAdvice("TestContextDeadline", Advice{
		Type:     Before,
		Priority: 100,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, "before")
			// Check if context has deadline
			deadline, ok := c.Context().Deadline()
			if !ok {
				t.Error("expected context to have deadline")
			} else {
				// Verify deadline is in the future (or very close to now)
				if deadline.Before(time.Now().Add(-1 * time.Second)) {
					t.Error("deadline appears to be in the past")
				}
			}
			return nil
		},
	})

	// Create a context with deadline
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Create and wrap function
	testFn := func(ctx context.Context) {
		executionOrder = append(executionOrder, "target")
	}
	wrappedFn := Wrap0Ctx(registry, "TestContextDeadline", testFn)

	// Execute the wrapped function
	wrappedFn(ctx)

	// Check execution order
	expectedOrder := []string{"before", "target"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

// TestContextValues verifies that context values propagate through advice
func TestContextValues(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()
	defer registry.Clear()

	var capturedValue string

	// Register function
	registry.MustRegister("TestContextValues")

	// Add advice that accesses context values
	registry.MustAddAdvice("TestContextValues", Advice{
		Type:     Before,
		Priority: 100,
		Handler: func(c *Context) error {
			// Get value from context
			value := c.Context().Value("test_key")
			if value != nil {
				if strVal, ok := value.(string); ok {
					capturedValue = strVal
				}
			}
			return nil
		},
	})

	// Create a context with values
	ctx := context.WithValue(context.Background(), "test_key", "test_value")

	// Create and wrap function
	testFn := func(ctx context.Context) {
		// Access context value in target function too
		value := ctx.Value("test_key")
		if value != "test_value" {
			t.Errorf("target function expected 'test_value', got %v", value)
		}
	}
	wrappedFn := Wrap0Ctx(registry, "TestContextValues", testFn)

	// Execute the wrapped function
	wrappedFn(ctx)

	// Check that value was captured in advice
	if capturedValue != "test_value" {
		t.Errorf("expected 'test_value' in advice, got '%s'", capturedValue)
	}
}

// TestContextPropagationThroughAllAdviceTypes verifies context propagation through all advice types
func TestContextPropagationThroughAllAdviceTypes(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()
	defer registry.Clear()

	var executionOrder []string

	// Register function
	registry.MustRegister("TestAllAdviceTypes")

	// Add all types of advice
	registry.MustAddAdvice("TestAllAdviceTypes", Advice{
		Type:     Before,
		Priority: 100,
		Handler: func(c *Context) error {
			if c.Context() == nil {
				t.Error("Before advice: context should not be nil")
			}
			executionOrder = append(executionOrder, "before")
			return nil
		},
	})

	registry.MustAddAdvice("TestAllAdviceTypes", Advice{
		Type:     Around,
		Priority: 100,
		Handler: func(c *Context) error {
			if c.Context() == nil {
				t.Error("Around advice: context should not be nil")
			}
			executionOrder = append(executionOrder, "around")
			return nil
		},
	})

	registry.MustAddAdvice("TestAllAdviceTypes", Advice{
		Type:     AfterReturning,
		Priority: 100,
		Handler: func(c *Context) error {
			if c.Context() == nil {
				t.Error("AfterReturning advice: context should not be nil")
			}
			executionOrder = append(executionOrder, "after-returning")
			return nil
		},
	})

	registry.MustAddAdvice("TestAllAdviceTypes", Advice{
		Type:     After,
		Priority: 100,
		Handler: func(c *Context) error {
			if c.Context() == nil {
				t.Error("After advice: context should not be nil")
			}
			executionOrder = append(executionOrder, "after")
			return nil
		},
	})

	// Create context
	ctx := context.Background()

	// Create and wrap function
	testFn := func(ctx context.Context) string {
		if ctx == nil {
			t.Error("Target function: context should not be nil")
		}
		executionOrder = append(executionOrder, "target")
		return "success"
	}
	wrappedFn := Wrap0RCtx[string](registry, "TestAllAdviceTypes", testFn)

	// Execute the wrapped function
	result := wrappedFn(ctx)

	// Check result
	if result != "success" {
		t.Errorf("expected 'success', got '%s'", result)
	}

	// Check execution order
	expectedOrder := []string{"before", "around", "target", "after-returning", "after"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

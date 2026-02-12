// Package aspect - fluent_test validates the fluent API for advice registration
package aspect

import (
	"errors"
	"testing"
)

// TestFluentAPI_BasicUsage tests the basic usage of the fluent API
func TestFluentAPI_BasicUsage(t *testing.T) {
	// Reset the default registry for clean test
	DefaultRegistry().Clear()
	defer DefaultRegistry().Clear()

	var executionOrder []string

	// Use the fluent API to configure advice
	For("TestFunction").
		WithBefore(func(c *Context) error {
			executionOrder = append(executionOrder, "before")
			return nil
		}).
		WithAfter(func(c *Context) error {
			executionOrder = append(executionOrder, "after")
			return nil
		}).
		WithAround(func(c *Context) error {
			executionOrder = append(executionOrder, "around-start")
			// Don't skip execution
			return nil
		})

	// Create and wrap a test function
	testFn := func() {
		executionOrder = append(executionOrder, "target")
	}
	builder := For("TestFunction")
	wrappedFn := Wrap0(builder.GetRegistry(), builder.GetFuncKey(), testFn)

	// Execute the wrapped function
	wrappedFn()

	// Check execution order
	expectedOrder := []string{"before", "around-start", "target", "after"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

// TestFluentAPI_WithPriorities tests the fluent API with priorities
func TestFluentAPI_WithPriorities(t *testing.T) {
	// Reset the default registry for clean test
	DefaultRegistry().Clear()
	defer DefaultRegistry().Clear()

	var executionOrder []string

	// Use the fluent API with priorities
	For("TestFunctionWithPriorities").
		WithBeforeP(func(c *Context) error {
			executionOrder = append(executionOrder, "before-low")
			return nil
		}, 10).
		WithBeforeP(func(c *Context) error {
			executionOrder = append(executionOrder, "before-high")
			return nil
		}, 100).
		WithAfterP(func(c *Context) error {
			executionOrder = append(executionOrder, "after-low")
			return nil
		}, 10).
		WithAfterP(func(c *Context) error {
			executionOrder = append(executionOrder, "after-high")
			return nil
		}, 100)

	// Create and wrap a test function
	testFn := func() {}
	builder := For("TestFunctionWithPriorities")
	wrappedFn := Wrap0(builder.GetRegistry(), builder.GetFuncKey(), testFn)

	// Execute the wrapped function
	wrappedFn()

	// Check execution order - higher priority should execute first
	// So before-high (priority 100) should run before before-low (priority 10)
	// And after-high (priority 100) should run before after-low (priority 10)
	expectedOrder := []string{"before-high", "before-low", "after-high", "after-low"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

// TestFluentAPI_WithRegistry tests the fluent API with a custom registry
func TestFluentAPI_WithRegistry(t *testing.T) {
	registry := NewRegistry()
	defer registry.Clear()

	var executionOrder []string

	// Use the fluent API with a custom registry
	ForWithRegistry(registry, "CustomRegistryFunction").
		WithBefore(func(c *Context) error {
			executionOrder = append(executionOrder, "custom-before")
			return nil
		}).
		WithAfter(func(c *Context) error {
			executionOrder = append(executionOrder, "custom-after")
			return nil
		})

	// Create and wrap a test function
	testFn := func() {}
	builder := ForWithRegistry(registry, "CustomRegistryFunction")
	wrappedFn := Wrap0(builder.GetRegistry(), builder.GetFuncKey(), testFn)

	// Execute the wrapped function
	wrappedFn()

	// Check execution order
	expectedOrder := []string{"custom-before", "custom-after"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

// TestFluentAPI_AroundAdviceSkip tests Around advice that skips execution
func TestFluentAPI_AroundAdviceSkip(t *testing.T) {
	// Reset the default registry for clean test
	DefaultRegistry().Clear()
	defer DefaultRegistry().Clear()

	var executionOrder []string

	// Use the fluent API with Around advice that skips execution
	For("SkippedFunction").
		WithBefore(func(c *Context) error {
			executionOrder = append(executionOrder, "before")
			return nil
		}).
		WithAround(func(c *Context) error {
			executionOrder = append(executionOrder, "around-skip")
			c.Skipped = true // Skip target execution
			c.SetResult(0, "skipped-result")
			return nil
		}).
		WithAfter(func(c *Context) error {
			executionOrder = append(executionOrder, "after")
			return nil
		})

	// Create and wrap a test function that would normally execute
	testFn := func() string {
		executionOrder = append(executionOrder, "target")
		return "normal-result"
	}
	builder := For("SkippedFunction")
	wrappedFn := Wrap0R(builder.GetRegistry(), builder.GetFuncKey(), testFn)

	// Execute the wrapped function
	result := wrappedFn()

	// Check execution order - target should be skipped
	expectedOrder := []string{"before", "around-skip", "after"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}

	// Check that the result comes from the Around advice
	if result != "skipped-result" {
		t.Errorf("expected result 'skipped-result', got '%s'", result)
	}
}

// TestFluentAPI_AfterReturning tests AfterReturning advice
func TestFluentAPI_AfterReturning(t *testing.T) {
	// Reset the default registry for clean test
	DefaultRegistry().Clear()
	defer DefaultRegistry().Clear()

	var executionOrder []string

	// Use the fluent API with AfterReturning advice
	For("AfterReturningFunction").
		WithAfterReturning(func(c *Context) error {
			executionOrder = append(executionOrder, "after-returning")
			return nil
		}).
		WithAfter(func(c *Context) error {
			executionOrder = append(executionOrder, "after")
			return nil
		})

	// Create and wrap a test function
	testFn := func() {}
	builder := For("AfterReturningFunction")
	wrappedFn := Wrap0(builder.GetRegistry(), builder.GetFuncKey(), testFn)

	// Execute the wrapped function
	wrappedFn()

	// Check execution order - AfterReturning should execute when no error occurs
	expectedOrder := []string{"after-returning", "after"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

// TestFluentAPI_AfterThrowing tests AfterThrowing advice
func TestFluentAPI_AfterThrowing(t *testing.T) {
	// Reset the default registry for clean test
	DefaultRegistry().Clear()
	defer DefaultRegistry().Clear()

	var executionOrder []string

	// Use the fluent API with AfterThrowing advice
	For("AfterThrowingFunction").
		WithAfterThrowing(func(c *Context) error {
			executionOrder = append(executionOrder, "after-throwing")
			return nil
		}).
		WithAfter(func(c *Context) error {
			executionOrder = append(executionOrder, "after")
			return nil
		})

	// Create and wrap a test function that panics
	testFn := func() {
		panic("test panic")
	}
	builder := For("AfterThrowingFunction")
	wrappedFn := Wrap0(builder.GetRegistry(), builder.GetFuncKey(), testFn)

	// Execute the wrapped function and catch the panic
	defer func() {
		if r := recover(); r != nil && r != "test panic" {
			t.Errorf("expected no panic! %v", r)
		}
	}()

	wrappedFn()

	// Check execution order - AfterThrowing should execute when panic occurs
	expectedOrder := []string{"after-throwing", "after"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

// TestFluentAPI_MultipleTypes tests combining multiple advice types
func TestFluentAPI_MultipleTypes(t *testing.T) {
	// Reset the default registry for clean test
	DefaultRegistry().Clear()
	defer DefaultRegistry().Clear()

	var executionOrder []string

	// Use the fluent API with multiple advice types
	For("MultiTypeFunction").
		WithBefore(func(c *Context) error {
			executionOrder = append(executionOrder, "before")
			return nil
		}).
		WithAround(func(c *Context) error {
			executionOrder = append(executionOrder, "around")
			return nil
		}).
		WithAfterReturning(func(c *Context) error {
			executionOrder = append(executionOrder, "after-returning")
			return nil
		}).
		WithAfter(func(c *Context) error {
			executionOrder = append(executionOrder, "after")
			return nil
		})

	// Create and wrap a test function
	testFn := func() {}
	builder := For("MultiTypeFunction")
	wrappedFn := Wrap0(builder.GetRegistry(), builder.GetFuncKey(), testFn)

	// Execute the wrapped function
	wrappedFn()

	// Check execution order
	expectedOrder := []string{"before", "around", "after-returning", "after"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

// TestFluentAPI_WithError tests error propagation through advice
func TestFluentAPI_WithError(t *testing.T) {
	// Reset the default registry for clean test
	DefaultRegistry().Clear()
	defer DefaultRegistry().Clear()

	var executionOrder []string

	// Use the fluent API with error handling
	For("ErrorFunction").
		WithBefore(func(c *Context) error {
			executionOrder = append(executionOrder, "before")
			return nil
		}).
		WithAfter(func(c *Context) error {
			executionOrder = append(executionOrder, "after")
			return nil
		})

	// Create and wrap a test function that returns an error
	testFn := func() error {
		executionOrder = append(executionOrder, "target-error")
		return errors.New("test error")
	}
	builder := For("ErrorFunction")
	wrappedFn := Wrap0E(builder.GetRegistry(), builder.GetFuncKey(), testFn)

	// Execute the wrapped function
	err := wrappedFn()

	// Check execution order
	expectedOrder := []string{"before", "target-error", "after"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("step %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}

	// Check that the error is properly propagated
	if err == nil || err.Error() != "test error" {
		t.Errorf("expected error 'test error', got %v", err)
	}
}

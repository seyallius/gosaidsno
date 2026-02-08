// Package aspect - integration_test validates end-to-end AOP functionality
package aspect

import (
	"errors"
	"testing"
	"time"
)

// -------------------------------------------- Integration Tests --------------------------------------------

func TestIntegration_CompleteWorkflow(t *testing.T) {
	// Clean up
	registry := NewRegistry()
	registry.Clear()

	// Register function
	err := registry.Register("CompleteWorkflow")
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	// Track execution order
	var executionOrder []string

	// Add Before advice
	_ = registry.AddAdvice("CompleteWorkflow", Advice{
		Type:     Before,
		Priority: 100,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, "before")
			return nil
		},
	})

	// Add Around advice
	_ = registry.AddAdvice("CompleteWorkflow", Advice{
		Type:     Around,
		Priority: 100,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, "around")
			return nil
		},
	})

	// Add AfterReturning advice
	_ = registry.AddAdvice("CompleteWorkflow", Advice{
		Type:     AfterReturning,
		Priority: 100,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, "afterReturning")
			return nil
		},
	})

	// Add After advice
	_ = registry.AddAdvice("CompleteWorkflow", Advice{
		Type:     After,
		Priority: 100,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, "after")
			return nil
		},
	})

	// Target function
	targetFunc := func(x int) int {
		executionOrder = append(executionOrder, "target")
		return x * 2
	}

	// Wrap and execute
	wrapped := Wrap1R(registry, "CompleteWorkflow", targetFunc)
	result := wrapped(5)

	// Verify result
	if result != 10 {
		t.Errorf("expected 10, got %d", result)
	}

	// Verify execution order
	expectedOrder := []string{"before", "around", "target", "afterReturning", "after"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d steps, got %d: %v", len(expectedOrder), len(executionOrder), executionOrder)
	}

	for i, step := range expectedOrder {
		if executionOrder[i] != step {
			t.Errorf("step %d: expected %s, got %s", i, step, executionOrder[i])
		}
	}

	// Clean up
	registry.Clear()
}

func TestIntegration_TimingPattern(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()

	_ = registry.Register("TimedOperation")

	// Before: record start time
	_ = registry.AddAdvice("TimedOperation", Advice{
		Type:     Before,
		Priority: 100,
		Handler: func(c *Context) error {
			c.SetMetadataVal("startTime", time.Now())
			c.SetMetadataVal("startTime", time.Now())
			return nil
		},
	})

	// After: calculate duration
	var duration time.Duration
	_ = registry.AddAdvice("TimedOperation", Advice{
		Type:     After,
		Priority: 100,
		Handler: func(c *Context) error {
			val, _ := c.GetMetadataVal("startTime")
			startTime := val.(time.Time)
			duration = time.Since(startTime)
			return nil
		},
	})

	// Target function (sleeps for a bit)
	targetFunc := func(ms int) {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	wrapped := Wrap1(registry, "TimedOperation", targetFunc)
	wrapped(50)

	if duration < 50*time.Millisecond {
		t.Errorf("expected duration >= 50ms, got %v", duration)
	}

	registry.Clear()
}

func TestIntegration_CachingPattern(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()

	_ = registry.Register("CachedFetch")

	cache := make(map[string]string)
	cache["key1"] = "cached_value"

	var targetExecuted bool

	// Around: check cache
	_ = registry.AddAdvice("CachedFetch", Advice{
		Type:     Around,
		Priority: 100,
		Handler: func(c *Context) error {
			key := c.Args[0].(string)
			if cached, ok := cache[key]; ok {
				c.SetResult(0, cached)
				c.Skipped = true
			}
			return nil
		},
	})

	targetFunc := func(key string) string {
		targetExecuted = true
		return "fresh_value"
	}

	wrapped := Wrap1R(registry, "CachedFetch", targetFunc)

	// Cache hit
	targetExecuted = false
	result1 := wrapped("key1")
	if result1 != "cached_value" {
		t.Errorf("expected cached_value, got %s", result1)
	}
	if targetExecuted {
		t.Error("target should not execute on cache hit")
	}

	// Cache miss
	targetExecuted = false
	result2 := wrapped("key2")
	if result2 != "fresh_value" {
		t.Errorf("expected fresh_value, got %s", result2)
	}
	if !targetExecuted {
		t.Error("target should execute on cache miss")
	}

	registry.Clear()
}

func TestIntegration_PanicRecoveryPattern(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()

	_ = registry.Register("RiskyOperation")

	var panicCaught bool
	var panicValue interface{}

	// AfterThrowing: catch panic
	_ = registry.AddAdvice("RiskyOperation", Advice{
		Type:     AfterThrowing,
		Priority: 100,
		Handler: func(c *Context) error {
			panicCaught = true
			panicValue = c.PanicValue
			return nil
		},
	})

	targetFunc := func(x int) {
		if x == 0 {
			panic("division by zero")
		}
	}

	wrapped := Wrap1(registry, "RiskyOperation", targetFunc)

	// Normal execution
	wrapped(10) // Should not panic

	// Panicking execution
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic to be re-thrown")
		}
		if r != "division by zero" {
			t.Errorf("unexpected panic value: %v", r)
		}
	}()

	wrapped(0) // Should panic

	if !panicCaught {
		t.Error("panic should have been caught by AfterThrowing")
	}
	if panicValue != "division by zero" {
		t.Errorf("unexpected panic value: %v", panicValue)
	}

	registry.Clear()
}

func TestIntegration_ErrorHandlingPattern(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()

	_ = registry.Register("ErrorOperation")

	var capturedError error

	// After: capture error
	_ = registry.AddAdvice("ErrorOperation", Advice{
		Type:     After,
		Priority: 100,
		Handler: func(c *Context) error {
			capturedError = c.Error
			return nil
		},
	})

	targetFunc := func(x int) error {
		if x < 0 {
			return errors.New("negative value")
		}
		return nil
	}

	wrapped := Wrap1E(registry, "ErrorOperation", targetFunc)

	// Success case
	err := wrapped(10)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if capturedError != nil {
		t.Error("should not capture error on success")
	}

	// Error case
	err = wrapped(-5)
	if err == nil {
		t.Error("expected error")
	}
	if capturedError == nil {
		t.Error("should capture error")
	}
	if capturedError.Error() != "negative value" {
		t.Errorf("unexpected error message: %s", capturedError.Error())
	}

	registry.Clear()
}

func TestIntegration_MultipleAdvicePriority(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()

	_ = registry.Register("PriorityTest")

	var executionOrder []int

	// Add advice with different priorities
	_ = registry.AddAdvice("PriorityTest", Advice{
		Type:     Before,
		Priority: 10,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, 10)
			return nil
		},
	})

	_ = registry.AddAdvice("PriorityTest", Advice{
		Type:     Before,
		Priority: 50,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, 50)
			return nil
		},
	})

	_ = registry.AddAdvice("PriorityTest", Advice{
		Type:     Before,
		Priority: 30,
		Handler: func(c *Context) error {
			executionOrder = append(executionOrder, 30)
			return nil
		},
	})

	targetFunc := func() {}

	wrapped := Wrap0(registry, "PriorityTest", targetFunc)
	wrapped()

	// Should execute in order: 50, 30, 10 (highest priority first)
	expected := []int{50, 30, 10}
	if len(executionOrder) != len(expected) {
		t.Fatalf("expected %d executions, got %d", len(expected), len(executionOrder))
	}

	for i, priority := range expected {
		if executionOrder[i] != priority {
			t.Errorf("step %d: expected priority %d, got %d", i, priority, executionOrder[i])
		}
	}

	registry.Clear()
}

func TestIntegration_AfterReturningOnlyOnSuccess(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()

	_ = registry.Register("ConditionalSuccess")

	var afterReturningCalled bool

	_ = registry.AddAdvice("ConditionalSuccess", Advice{
		Type:     AfterReturning,
		Priority: 100,
		Handler: func(c *Context) error {
			afterReturningCalled = true
			return nil
		},
	})

	// Success case
	successFunc := func(x int) error {
		return nil
	}

	wrapped := Wrap1E(registry, "ConditionalSuccess", successFunc)
	_ = wrapped(10)

	if !afterReturningCalled {
		t.Error("AfterReturning should be called on success")
	}

	// Error case
	afterReturningCalled = false
	errorFunc := func(x int) error {
		return errors.New("error")
	}

	wrapped2 := Wrap1E(registry, "ConditionalSuccess", errorFunc)
	_ = wrapped2(10)

	if afterReturningCalled {
		t.Error("AfterReturning should NOT be called on error")
	}

	registry.Clear()
}

func TestIntegration_MetadataPassingBetweenAdvice(t *testing.T) {
	registry := NewRegistry()
	registry.Clear()

	_ = registry.Register("MetadataTest")

	// Before: set metadata
	_ = registry.AddAdvice("MetadataTest", Advice{
		Type:     Before,
		Priority: 100,
		Handler: func(c *Context) error {
			c.SetMetadataVal("userId", "user_123")
			c.SetMetadataVal("requestId", "req_456")
			return nil
		},
	})

	// After: read metadata
	var userId, requestId string
	_ = registry.AddAdvice("MetadataTest", Advice{
		Type:     After,
		Priority: 100,
		Handler: func(c *Context) error {
			userId = c.Metadata["userId"].(string)
			requestId = c.Metadata["requestId"].(string)
			return nil
		},
	})

	targetFunc := func() {}

	wrapped := Wrap0(registry, "MetadataTest", targetFunc)
	wrapped()

	if userId != "user_123" {
		t.Errorf("expected user_123, got %s", userId)
	}
	if requestId != "req_456" {
		t.Errorf("expected req_456, got %s", requestId)
	}

	registry.Clear()
}

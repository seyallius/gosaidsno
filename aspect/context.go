// Package aspect - context provides execution context for aspect-oriented advice
package aspect

import "fmt"

// -------------------------------------------- Types --------------------------------------------

// Context holds the execution state for a single function invocation.
// It captures arguments, return values, errors, and panic information.
type Context struct {
	FunctionName string         // FunctionName is the registered name of the wrapped function.
	Args         []any          // Args contains the function arguments (caller must cast to correct types).
	Results      []any          // Results contains the function return values (populated after execution).
	Error        error          // Error holds any error returned by the function.
	PanicValue   any            // PanicValue holds the recovered panic value if a panic occurred.
	Metadata     map[string]any // Metadata allows storing custom key-value pairs for advice communication.
	Skipped      bool           // Skipped indicates if the target function execution should be skipped (set by Around advice).
}

// NewContext creates a new execution context for the given function.
func NewContext(functionName string, args ...any) *Context {
	return &Context{
		FunctionName: functionName,
		Args:         args,
		Metadata:     make(map[string]any),
		Results:      make([]any, 0),
	}
}

// -------------------------------------------- Public Functions --------------------------------------------

// SetResult sets a return value at the specified index.
func (c *Context) SetResult(index int, value any) {
	if index < 0 {
		return // Invalid index
	}

	// Extend results slice if needed
	for len(c.Results) <= index {
		c.Results = append(c.Results, nil)
	}
	c.Results[index] = value
}

// GetResult retrieves a return value at the specified index.
func (c *Context) GetResult(index int) any {
	if index < 0 || index >= len(c.Results) {
		return nil
	}
	return c.Results[index]
}

// HasPanic returns true if a panic was recovered during execution.
func (c *Context) HasPanic() bool {
	return c.PanicValue != nil
}

// String returns a formatted string representation of the context implementing fmt.Stringer interface.
func (c *Context) String() string {
	return fmt.Sprintf("Context{Function: %s, Args: %v, Results: %v, Error: %v, Panic: %v}",
		c.FunctionName, c.Args, c.Results, c.Error, c.PanicValue)
}

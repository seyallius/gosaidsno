// Package aspect - wrap provides function wrapping utilities with AOP advice execution
package aspect

import (
	"context"
	"fmt"
)

// -------------------------------------------- Public Functions --------------------------------------------

// Wrap0 wraps a function with no arguments and no return values.
func Wrap0(registry *Registry, funcKey FuncKey, fn func()) func() {
	return func() {
		executeWithAdvice(registry, funcKey, func(c *Context) {
			fn()
		})
	}
}

// Wrap0R wraps a function with no arguments and one return value.
func Wrap0R[R any](registry *Registry, funcKey FuncKey, fn func() R) func() R {
	return func() R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(c *Context) {
			result = fn()
			c.SetResult(0, result)
		})

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}

		return result
	}
}

// Wrap0E wraps a function with no arguments and returns error.
func Wrap0E(registry *Registry, funcKey FuncKey, fn func() error) func() error {
	return func() error {
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *Context) {
			err = fn()
			c.Error = err
		})

		if c != nil && c.Error != nil {
			err = c.Error
		}

		return err
	}
}

// Wrap0RE wraps a function with no arguments and returns (result, error).
func Wrap0RE[R any](registry *Registry, funcKey FuncKey, fn func() (R, error)) func() (R, error) {
	return func() (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *Context) {
			result, err = fn()
			c.SetResult(0, result)
			c.Error = err
		})

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}
		if c != nil && c.Error != nil {
			err = c.Error
		}

		return result, err
	}
}

// Wrap1 wraps a function with one argument and no return values.
func Wrap1[A any](registry *Registry, funcKey FuncKey, fn func(A)) func(A) {
	return func(a A) {
		executeWithAdvice(registry, funcKey, func(c *Context) {
			fn(a)
		}, a)
	}
}

// Wrap1R wraps a function with one argument and one return value.
func Wrap1R[A, R any](registry *Registry, funcKey FuncKey, fn func(A) R) func(A) R {
	return func(a A) R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(c *Context) {
			result = fn(a)
			c.SetResult(0, result)
		}, a)

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}

		return result
	}
}

// Wrap1RE wraps a function with one argument and returns (result, error).
func Wrap1RE[A, R any](registry *Registry, funcKey FuncKey, fn func(A) (R, error)) func(A) (R, error) {
	return func(a A) (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *Context) {
			result, err = fn(a)
			c.SetResult(0, result)
			c.Error = err
		}, a)

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}
		if c != nil && c.Error != nil {
			err = c.Error
		}

		return result, err
	}
}

// Wrap1E wraps a function with one argument and returns error.
func Wrap1E[A any](registry *Registry, funcKey FuncKey, fn func(A) error) func(A) error {
	return func(a A) error {
		var err error
		executeWithAdvice(registry, funcKey, func(c *Context) {
			err = fn(a)
			c.Error = err
		}, a)
		return err
	}
}

// Wrap2 wraps a function with two arguments and no return values.
func Wrap2[A, B any](registry *Registry, funcKey FuncKey, fn func(A, B)) func(A, B) {
	return func(a A, b B) {
		executeWithAdvice(registry, funcKey, func(c *Context) {
			fn(a, b)
		}, a, b)
	}
}

// Wrap2R wraps a function with two arguments and one return value.
func Wrap2R[A, B, R any](registry *Registry, funcKey FuncKey, fn func(A, B) R) func(A, B) R {
	return func(a A, b B) R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(c *Context) {
			result = fn(a, b)
			c.SetResult(0, result)
		}, a, b)

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}

		return result
	}
}

// Wrap2RE wraps a function with two arguments and returns (result, error).
func Wrap2RE[A, B, R any](registry *Registry, funcKey FuncKey, fn func(A, B) (R, error)) func(A, B) (R, error) {
	return func(a A, b B) (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *Context) {
			result, err = fn(a, b)
			c.SetResult(0, result)
			c.Error = err
		}, a, b)

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}
		if c != nil && c.Error != nil {
			err = c.Error
		}

		return result, err
	}
}

// Wrap2E wraps a function with two arguments and returns error.
func Wrap2E[A, B any](registry *Registry, funcKey FuncKey, fn func(A, B) error) func(A, B) error {
	return func(a A, b B) error {
		var err error
		executeWithAdvice(registry, funcKey, func(c *Context) {
			err = fn(a, b)
			c.Error = err
		}, a, b)
		return err
	}
}

// Wrap3 wraps a function with three arguments and no return values.
func Wrap3[A, B, C any](registry *Registry, funcKey FuncKey, fn func(A, B, C)) func(A, B, C) {
	return func(a A, b B, c C) {
		executeWithAdvice(registry, funcKey, func(ct *Context) {
			fn(a, b, c)
		}, a, b, c)
	}
}

// Wrap3R wraps a function with three arguments and one return value.
func Wrap3R[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(A, B, C) R) func(A, B, C) R {
	return func(a A, b B, paramC C) R {
		var result R
		ctx := executeWithAdvice(registry, funcKey, func(ct *Context) {
			result = fn(a, b, paramC)
			ct.SetResult(0, result)
		}, a, b, paramC)

		// If Around advice set a result and skipped execution, use that result
		if ctx != nil && ctx.Skipped && len(ctx.Results) > 0 && ctx.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := ctx.Results[0].(R); ok {
				result = res
			}
		}

		return result
	}
}

// Wrap3RE wraps a function with three arguments and returns (result, error).
func Wrap3RE[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(A, B, C) (R, error)) func(A, B, C) (R, error) {
	return func(a A, b B, paramC C) (R, error) {
		var result R
		var err error
		ctx := executeWithAdvice(registry, funcKey, func(ct *Context) {
			result, err = fn(a, b, paramC)
			ct.SetResult(0, result)
			ct.Error = err
		}, a, b, paramC)

		// If Around advice set a result and skipped execution, use that result
		if ctx != nil && ctx.Skipped && len(ctx.Results) > 0 && ctx.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := ctx.Results[0].(R); ok {
				result = res
			}
		}
		if ctx != nil && ctx.Error != nil {
			err = ctx.Error
		}

		return result, err
	}
}

// Wrap3E wraps a function with three arguments and returns error.
func Wrap3E[A, B, C any](registry *Registry, funcKey FuncKey, fn func(A, B, C) error) func(A, B, C) error {
	return func(a A, b B, c C) error {
		var err error
		executeWithAdvice(registry, funcKey, func(ct *Context) {
			err = fn(a, b, c)
			ct.Error = err
		}, a, b, c)
		return err
	}
}

// Context-aware wrapper functions that accept context.Context

// Wrap0Ctx wraps a function with no arguments and no return values with context.
func Wrap0Ctx(registry *Registry, funcKey FuncKey, fn func(context.Context)) func(context.Context) {
	return func(ctx context.Context) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			fn(c.Context())
		})
	}
}

// Wrap0RCtx wraps a function with no arguments and one return value with context.
func Wrap0RCtx[R any](registry *Registry, funcKey FuncKey, fn func(context.Context) R) func(context.Context) R {
	return func(ctx context.Context) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result = fn(c.Context())
			c.SetResult(0, result)
		})

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}

		return result
	}
}

// Wrap0ECtx wraps a function with no arguments and returns error with context.
func Wrap0ECtx(registry *Registry, funcKey FuncKey, fn func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			err = fn(c.Context())
			c.Error = err
		})

		if c != nil && c.Error != nil {
			err = c.Error
		}

		return err
	}
}

// Wrap0RECtx wraps a function with no arguments and returns (result, error) with context.
func Wrap0RECtx[R any](registry *Registry, funcKey FuncKey, fn func(context.Context) (R, error)) func(context.Context) (R, error) {
	return func(ctx context.Context) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result, err = fn(c.Context())
			c.SetResult(0, result)
			c.Error = err
		})

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}
		if c != nil && c.Error != nil {
			err = c.Error
		}

		return result, err
	}
}

// Wrap1Ctx wraps a function with one argument and no return values with context.
func Wrap1Ctx[A any](registry *Registry, funcKey FuncKey, fn func(context.Context, A)) func(context.Context, A) {
	return func(ctx context.Context, a A) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			fn(c.Context(), a)
		}, a)
	}
}

// Wrap1RCtx wraps a function with one argument and one return value with context.
func Wrap1RCtx[A, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A) R) func(context.Context, A) R {
	return func(ctx context.Context, a A) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result = fn(c.Context(), a)
			c.SetResult(0, result)
		}, a)

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}

		return result
	}
}

// Wrap1RECtx wraps a function with one argument and returns (result, error) with context.
func Wrap1RECtx[A, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A) (R, error)) func(context.Context, A) (R, error) {
	return func(ctx context.Context, a A) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result, err = fn(c.Context(), a)
			c.SetResult(0, result)
			c.Error = err
		}, a)

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}
		if c != nil && c.Error != nil {
			err = c.Error
		}

		return result, err
	}
}

// Wrap1ECtx wraps a function with one argument and returns error with context.
func Wrap1ECtx[A any](registry *Registry, funcKey FuncKey, fn func(context.Context, A) error) func(context.Context, A) error {
	return func(ctx context.Context, a A) error {
		var err error
		executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			err = fn(c.Context(), a)
			c.Error = err
		}, a)
		return err
	}
}

// Wrap2Ctx wraps a function with two arguments and no return values with context.
func Wrap2Ctx[A, B any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B)) func(context.Context, A, B) {
	return func(ctx context.Context, a A, b B) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			fn(c.Context(), a, b)
		}, a, b)
	}
}

// Wrap2RCtx wraps a function with two arguments and one return value with context.
func Wrap2RCtx[A, B, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B) R) func(context.Context, A, B) R {
	return func(ctx context.Context, a A, b B) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result = fn(c.Context(), a, b)
			c.SetResult(0, result)
		}, a, b)

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}

		return result
	}
}

// Wrap2RECtx wraps a function with two arguments and returns (result, error) with context.
func Wrap2RECtx[A, B, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B) (R, error)) func(context.Context, A, B) (R, error) {
	return func(ctx context.Context, a A, b B) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result, err = fn(c.Context(), a, b)
			c.SetResult(0, result)
			c.Error = err
		}, a, b)

		// If Around advice set a result and skipped execution, use that result
		if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := c.Results[0].(R); ok {
				result = res
			}
		}
		if c != nil && c.Error != nil {
			err = c.Error
		}

		return result, err
	}
}

// Wrap2ECtx wraps a function with two arguments and returns error with context.
func Wrap2ECtx[A, B any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B) error) func(context.Context, A, B) error {
	return func(ctx context.Context, a A, b B) error {
		var err error
		executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			err = fn(c.Context(), a, b)
			c.Error = err
		}, a, b)
		return err
	}
}

// Wrap3Ctx wraps a function with three arguments and no return values with context.
func Wrap3Ctx[A, B, C any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B, C)) func(context.Context, A, B, C) {
	return func(ctx context.Context, a A, b B, c C) {
		executeWithAdviceContext(registry, funcKey, ctx, func(ct *Context) {
			fn(ct.Context(), a, b, c)
		}, a, b, c)
	}
}

// Wrap3RCtx wraps a function with three arguments and one return value with context.
func Wrap3RCtx[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B, C) R) func(context.Context, A, B, C) R {
	return func(ctx context.Context, a A, b B, paramC C) R {
		var result R
		cxt := executeWithAdviceContext(registry, funcKey, ctx, func(ct *Context) {
			result = fn(ct.Context(), a, b, paramC)
			ct.SetResult(0, result)
		}, a, b, paramC)

		// If Around advice set a result and skipped execution, use that result
		if cxt != nil && cxt.Skipped && len(cxt.Results) > 0 && cxt.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := cxt.Results[0].(R); ok {
				result = res
			}
		}

		return result
	}
}

// Wrap3RECtx wraps a function with three arguments and returns (result, error) with context.
func Wrap3RECtx[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B, C) (R, error)) func(context.Context, A, B, C) (R, error) {
	return func(ctx context.Context, a A, b B, paramC C) (R, error) {
		var result R
		var err error
		cxt := executeWithAdviceContext(registry, funcKey, ctx, func(ct *Context) {
			result, err = fn(ct.Context(), a, b, paramC)
			ct.SetResult(0, result)
			ct.Error = err
		}, a, b, paramC)

		// If Around advice set a result and skipped execution, use that result
		if cxt != nil && cxt.Skipped && len(cxt.Results) > 0 && cxt.Results[0] != nil {
			// Safe type assertion with proper handling
			if res, ok := cxt.Results[0].(R); ok {
				result = res
			}
		}
		if cxt != nil && cxt.Error != nil {
			err = cxt.Error
		}

		return result, err
	}
}

// Wrap3ECtx wraps a function with three arguments and returns error with context.
func Wrap3ECtx[A, B, C any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B, C) error) func(context.Context, A, B, C) error {
	return func(ctx context.Context, a A, b B, c C) error {
		var err error
		executeWithAdviceContext(registry, funcKey, ctx, func(ct *Context) {
			err = fn(ct.Context(), a, b, c)
			ct.Error = err
		}, a, b, c)
		return err
	}
}

// -------------------------------------------- Private Helper Functions --------------------------------------------

// executeWithAdvice executes a function with full advice chain support and returns the context.
func executeWithAdvice(registry *Registry, functionName FuncKey, targetFn func(*Context), args ...any) *Context {
	return executeWithAdviceContext(registry, functionName, context.Background(), targetFn, args...)
}

// executeWithAdviceContext executes a function with full advice chain support using a specific context.Context.
func executeWithAdviceContext(registry *Registry, functionName FuncKey, ctx context.Context, targetFn func(*Context), args ...any) *Context {
	// Get advice chain from registry
	chain, err := registry.GetAdviceChain(functionName)
	if err != nil {
		// No advice registered, just execute target function
		c := NewContextWithContext(ctx, functionName, args...)
		targetFn(c)
		return c
	}

	// Create execution context
	c := NewContextWithContext(ctx, functionName, args...)

	// Defer After advice (always runs)
	defer func() {
		_ = chain.ExecuteAfter(c)
	}()

	// Defer panic recovery and AfterThrowing advice
	defer func() {
		if r := recover(); r != nil {
			c.PanicValue = r
			_ = chain.ExecuteAfterThrowing(c)

			// Re-panic to maintain panic semantics
			panic(r)
		}
	}()

	// Execute Before advice
	if err = chain.ExecuteBefore(c); err != nil {
		panic(fmt.Errorf("before advice failed: %w", err))
	}

	// Execute Around advice (if any)
	if chain.HasAround() {
		if err := chain.ExecuteAround(c); err != nil {
			panic(fmt.Errorf("around advice failed: %w", err))
		}
		// If Around advice sets Skipped, don't execute target function
		if c.Skipped {
			// Execute AfterReturning if no error (Around advice might have set result)
			if c.Error == nil && !c.HasPanic() {
				_ = chain.ExecuteAfterReturning(c)
			}
			return c
		}
	}

	// Execute target function
	targetFn(c)

	// Execute AfterReturning advice (only if no error and no panic)
	if c.Error == nil && !c.HasPanic() {
		_ = chain.ExecuteAfterReturning(c)
	}

	return c
}

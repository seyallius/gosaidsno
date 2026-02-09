// Package aspect - wrap provides function wrapping utilities with AOP advice execution
package aspect

import (
	"context"
	"fmt"
)

// -------------------------------------------- Public Functions --------------------------------------------

// -- 0 Arguments --

// Wrap0 wraps a function with no arguments and no return values.
func Wrap0(registry *Registry, funcKey FuncKey, fn func()) func() {
	return func() {
		executeWithAdvice(registry, funcKey, func(c *Context) {
			fn()
		})
	}
}

// Wrap0Ctx wraps a function with context, no arguments, no returns.
func Wrap0Ctx(registry *Registry, funcKey FuncKey, fn func(context.Context)) func(context.Context) {
	return func(ctx context.Context) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			fn(c.Context())
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
		return resolveResult(c, result)
	}
}

// Wrap0RCtx wraps a function with context, no arguments, one return.
func Wrap0RCtx[R any](registry *Registry, funcKey FuncKey, fn func(context.Context) R) func(context.Context) R {
	return func(ctx context.Context) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result = fn(c.Context())
			c.SetResult(0, result)
		})
		return resolveResult(c, result)
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
		return resolveError(c, err)
	}
}

// Wrap0ECtx wraps a function with context, no arguments, returns error.
func Wrap0ECtx(registry *Registry, funcKey FuncKey, fn func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			err = fn(c.Context())
			c.Error = err
		})
		return resolveError(c, err)
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
		return resolveResultError(c, result, err)
	}
}

// Wrap0RECtx wraps a function with context, no arguments, returns (result, error).
func Wrap0RECtx[R any](registry *Registry, funcKey FuncKey, fn func(context.Context) (R, error)) func(context.Context) (R, error) {
	return func(ctx context.Context) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result, err = fn(c.Context())
			c.SetResult(0, result)
			c.Error = err
		})
		return resolveResultError(c, result, err)
	}
}

// -- 1 Argument --

// Wrap1 wraps a function with one argument and no return values.
func Wrap1[A any](registry *Registry, funcKey FuncKey, fn func(A)) func(A) {
	return func(a A) {
		executeWithAdvice(registry, funcKey, func(c *Context) {
			fn(a)
		}, a)
	}
}

// Wrap1Ctx wraps a function with context, 1 arg, no returns.
func Wrap1Ctx[A any](registry *Registry, funcKey FuncKey, fn func(context.Context, A)) func(context.Context, A) {
	return func(ctx context.Context, a A) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			fn(c.Context(), a)
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
		return resolveResult(c, result)
	}
}

// Wrap1RCtx wraps a function with context, 1 arg, one return.
func Wrap1RCtx[A, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A) R) func(context.Context, A) R {
	return func(ctx context.Context, a A) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result = fn(c.Context(), a)
			c.SetResult(0, result)
		}, a)
		return resolveResult(c, result)
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
		return err // executeWithAdvice returns *Context, but simplistic wrapper can just return captured err if no complex mutation needed
	}
}

// Wrap1ECtx wraps a function with context, 1 arg, returns error.
func Wrap1ECtx[A any](registry *Registry, funcKey FuncKey, fn func(context.Context, A) error) func(context.Context, A) error {
	return func(ctx context.Context, a A) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			err = fn(c.Context(), a)
			c.Error = err
		}, a)
		return resolveError(c, err)
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
		return resolveResultError(c, result, err)
	}
}

// Wrap1RECtx wraps a function with context, 1 arg, returns (result, error).
func Wrap1RECtx[A, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A) (R, error)) func(context.Context, A) (R, error) {
	return func(ctx context.Context, a A) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result, err = fn(c.Context(), a)
			c.SetResult(0, result)
			c.Error = err
		}, a)
		return resolveResultError(c, result, err)
	}
}

// -- 2 Arguments --

// Wrap2 wraps a function with two arguments and no return values.
func Wrap2[A, B any](registry *Registry, funcKey FuncKey, fn func(A, B)) func(A, B) {
	return func(a A, b B) {
		executeWithAdvice(registry, funcKey, func(c *Context) {
			fn(a, b)
		}, a, b)
	}
}

// Wrap2Ctx wraps a function with context, 2 args, no returns.
func Wrap2Ctx[A, B any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B)) func(context.Context, A, B) {
	return func(ctx context.Context, a A, b B) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			fn(c.Context(), a, b)
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
		return resolveResult(c, result)
	}
}

// Wrap2RCtx wraps a function with context, 2 args, one return.
func Wrap2RCtx[A, B, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B) R) func(context.Context, A, B) R {
	return func(ctx context.Context, a A, b B) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result = fn(c.Context(), a, b)
			c.SetResult(0, result)
		}, a, b)
		return resolveResult(c, result)
	}
}

// Wrap2E wraps a function with two arguments and returns error.
func Wrap2E[A, B any](registry *Registry, funcKey FuncKey, fn func(A, B) error) func(A, B) error {
	return func(a A, b B) error {
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *Context) {
			err = fn(a, b)
			c.Error = err
		}, a, b)
		return resolveError(c, err)
	}
}

// Wrap2ECtx wraps a function with context, 2 args, returns error.
func Wrap2ECtx[A, B any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B) error) func(context.Context, A, B) error {
	return func(ctx context.Context, a A, b B) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			err = fn(c.Context(), a, b)
			c.Error = err
		}, a, b)
		return resolveError(c, err)
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
		return resolveResultError(c, result, err)
	}
}

// Wrap2RECtx wraps a function with context, 2 args, returns (result, error).
func Wrap2RECtx[A, B, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B) (R, error)) func(context.Context, A, B) (R, error) {
	return func(ctx context.Context, a A, b B) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *Context) {
			result, err = fn(c.Context(), a, b)
			c.SetResult(0, result)
			c.Error = err
		}, a, b)
		return resolveResultError(c, result, err)
	}
}

// -- 3 Arguments --

// Wrap3 wraps a function with three arguments and no return values.
func Wrap3[A, B, C any](registry *Registry, funcKey FuncKey, fn func(A, B, C)) func(A, B, C) {
	return func(a A, b B, c C) {
		executeWithAdvice(registry, funcKey, func(ct *Context) {
			fn(a, b, c)
		}, a, b, c)
	}
}

// Wrap3Ctx wraps a function with context, 3 args, no returns.
func Wrap3Ctx[A, B, C any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B, C)) func(context.Context, A, B, C) {
	return func(ctx context.Context, a A, b B, c C) {
		executeWithAdviceContext(registry, funcKey, ctx, func(ct *Context) {
			fn(ct.Context(), a, b, c)
		}, a, b, c)
	}
}

// Wrap3R wraps a function with three arguments and one return value.
func Wrap3R[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(A, B, C) R) func(A, B, C) R {
	return func(a A, b B, paramC C) R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(ct *Context) {
			result = fn(a, b, paramC)
			ct.SetResult(0, result)
		}, a, b, paramC)
		return resolveResult(c, result)
	}
}

// Wrap3RCtx wraps a function with context, 3 args, one return.
func Wrap3RCtx[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B, C) R) func(context.Context, A, B, C) R {
	return func(ctx context.Context, a A, b B, paramC C) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(ct *Context) {
			result = fn(ct.Context(), a, b, paramC)
			ct.SetResult(0, result)
		}, a, b, paramC)
		return resolveResult(c, result)
	}
}

// Wrap3E wraps a function with three arguments and returns error.
func Wrap3E[A, B, C any](registry *Registry, funcKey FuncKey, fn func(A, B, C) error) func(A, B, C) error {
	return func(a A, b B, c C) error {
		var err error
		ctx := executeWithAdvice(registry, funcKey, func(ct *Context) {
			err = fn(a, b, c)
			ct.Error = err
		}, a, b, c)
		return resolveError(ctx, err)
	}
}

// Wrap3ECtx wraps a function with context, 3 args, returns error.
func Wrap3ECtx[A, B, C any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B, C) error) func(context.Context, A, B, C) error {
	return func(ctx context.Context, a A, b B, c C) error {
		var err error
		ct := executeWithAdviceContext(registry, funcKey, ctx, func(ct *Context) {
			err = fn(ct.Context(), a, b, c)
			ct.Error = err
		}, a, b, c)
		return resolveError(ct, err)
	}
}

// Wrap3RE wraps a function with three arguments and returns (result, error).
func Wrap3RE[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(A, B, C) (R, error)) func(A, B, C) (R, error) {
	return func(a A, b B, paramC C) (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(ct *Context) {
			result, err = fn(a, b, paramC)
			ct.SetResult(0, result)
			ct.Error = err
		}, a, b, paramC)
		return resolveResultError(c, result, err)
	}
}

// Wrap3RECtx wraps a function with context, 3 args, returns (result, error).
func Wrap3RECtx[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(context.Context, A, B, C) (R, error)) func(context.Context, A, B, C) (R, error) {
	return func(ctx context.Context, a A, b B, paramC C) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(ct *Context) {
			result, err = fn(ct.Context(), a, b, paramC)
			ct.SetResult(0, result)
			ct.Error = err
		}, a, b, paramC)
		return resolveResultError(c, result, err)
	}
}

// -------------------------------------------- Private Helper Functions --------------------------------------------

// resolveResult handles the logic for extracting a generic result from the context,
// checking for advice skips, and performing safe type assertions.
func resolveResult[R any](c *Context, original R) R {
	// If Around advice skipped execution and set a result, try to use it
	if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
		if res, ok := c.Results[0].(R); ok {
			return res
		}
	}
	return original
}

// resolveError handles the logic for extracting an error from the context,
// allowing advice chains to replace the original error.
func resolveError(c *Context, original error) error {
	if c != nil && c.Error != nil {
		return c.Error
	}
	return original
}

// resolveResultError combines result and error resolution for functions returning (R, error).
func resolveResultError[R any](c *Context, origRes R, origErr error) (R, error) {
	finalRes := resolveResult(c, origRes)
	finalErr := resolveError(c, origErr)
	return finalRes, finalErr
}

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

	// Ensure After advice always runs
	defer func() {
		_ = chain.ExecuteAfter(c)
	}()

	// Handle Panic Recovery and Throwing advice
	defer func() {
		if r := recover(); r != nil {
			c.PanicValue = r
			_ = chain.ExecuteAfterThrowing(c)
			// Re-panic to maintain panic semantics for the caller
			panic(r)
		}
	}()

	// Execute Before advice
	if err = chain.ExecuteBefore(c); err != nil {
		// If Before advice fails, we typically panic or stop (design choice from original code)
		panic(fmt.Errorf("before advice failed: %w", err))
	}

	// Execute Around advice
	if chain.HasAround() {
		if err := chain.ExecuteAround(c); err != nil {
			panic(fmt.Errorf("around advice failed: %w", err))
		}
		// If Around advice sets Skipped, we skip the target function
		if c.Skipped {
			// Execute AfterReturning if no error (Around advice might have set result)
			if c.Error == nil && !c.HasPanic() {
				_ = chain.ExecuteAfterReturning(c)
			}
			return c
		}
	}

	// Execute Target Function
	targetFn(c)

	// Execute AfterReturning advice (only if no error and no panic occurred)
	if c.Error == nil && !c.HasPanic() {
		_ = chain.ExecuteAfterReturning(c)
	}

	return c
}

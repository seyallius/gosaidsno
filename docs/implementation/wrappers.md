# Wrapper Functions

Wrapper functions serve as the bridge between your code and the AOP system. They're the entry point that transforms regular function calls into AOP-enabled calls with all the associated advice execution.

## Wrapper Structure

```go
func Wrap1RE[A, R any](name string, fn func(A) (R, error)) func(A) (R, error) {
    return func(a A) (R, error) {
        var result R
        var err error
        
        executeWithAdvice(name, func(ctx *Context) {
            result, err = fn(a)           // Execute target function
            ctx.SetResult(0, result)      // Store result in context
            ctx.Error = err               // Store error in context
        }, a)
        
        return result, err
    }
}
```

## Why This Two-Step Process?

The wrapper performs several critical functions:

### 1. Closure Creation
The wrapper creates a closure that captures the original function, preserving it for later execution while adding AOP capabilities.

### 2. Type Safety Maintenance
Through generics, the wrapper maintains type safety for the function signature while adding AOP functionality.

### 3. Context Management
The wrapper handles the creation and management of the execution context, ensuring proper data flow between the original function and its advice.

### 4. Result Handling
The wrapper manages the return of results from the execution engine back to the caller.

## Generic Implementation Benefits

### Type Safety
- Compile-time checking of function signatures
- No runtime type errors for correct usage
- Full IDE support with autocompletion

### Performance
- No reflection overhead
- Direct function calls without type inspection
- Optimized by the compiler for specific types

### Maintainability
- Clear function signatures
- Easy to understand type relationships
- Good error messages from the compiler

## Limitations of Generic Wrappers

### Limited Arity Support
Currently, only functions with up to 3 arguments are supported:
- `Wrap0`, `Wrap0R`, `Wrap0RE` - No arguments
- `Wrap1`, `Wrap1R`, `Wrap1RE`, `Wrap1E` - One argument
- `Wrap2`, `Wrap2R`, `Wrap2RE`, `Wrap2E` - Two arguments
- `Wrap3RE` - Three arguments

For functions with more arguments, you can:
- Create custom wrappers
- Refactor to use a single struct parameter
- Use variadic functions with manual handling

### Complex Signature Challenges
Functions with complex return types or multiple return values beyond (result, error) patterns require custom handling.

### Code Duplication
Each argument count requires its own set of wrapper functions, leading to some code duplication. This is a trade-off for type safety and performance.

## Wrapper Generation Pattern

Different wrapper types handle different function signatures:

### No Return Values
```go
func Wrap1[A any](name string, fn func(A)) func(A) {
    return func(a A) {
        executeWithAdvice(name, func(ctx *Context) {
            fn(a)
        }, a)
    }
}
```

### Single Return Value
```go
func Wrap1R[A, R any](name string, fn func(A) R) func(A) R {
    return func(a A) R {
        var result R
        ctx := executeWithAdvice(name, func(ctx *Context) {
            result = fn(a)
            ctx.SetResult(0, result)
        }, a)
        
        // Handle result from Around advice if target was skipped
        if ctx != nil && ctx.Skipped && len(ctx.Results) > 0 && ctx.Results[0] != nil {
            if res, ok := ctx.Results[0].(R); ok {
                result = res
            }
        }
        
        return result
    }
}
```

### Result and Error
```go
func Wrap1RE[A, R any](name string, fn func(A) (R, error)) func(A) (R, error) {
    return func(a A) (R, error) {
        var result R
        var err error
        executeWithAdvice(name, func(ctx *Context) {
            result, err = fn(a)
            ctx.SetResult(0, result)
            ctx.Error = err
        }, a)
        return result, err
    }
}
```

## Performance Considerations

- **Function Call Overhead**: One additional function call per wrapped function
- **Closure Creation**: Minimal overhead during wrapper creation
- **Type Assertion**: When Around advice provides results, type assertions occur
- **Memory Allocation**: Context creation per call

## Memory Management

Each wrapper creates a closure that captures the original function, but this is typically a one-time cost during application initialization. The actual function call involves creating a context and executing the advice chain.

The wrapper functions are the user-facing API of the AOP system, providing a clean, type-safe interface that hides the complexity of advice execution while maintaining performance and usability.
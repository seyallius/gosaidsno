# Context Object

The Context is the most important object in the system - it carries state between all components and serves as the communication hub for the entire AOP execution. It's the only object that's passed to every advice function.

## Context Structure

```go
type Context struct {
    FunctionName string         // Identity of the function being executed
    Args         []any          // Input parameters (can be modified)
    Results      []any          // Output results (can be modified)
    Error        error          // Function error (can be modified)
    PanicValue   any            // Panic value if function panicked
    Metadata     map[string]any // Shared data between advice functions
    Skipped      bool           // Signal from Around advice to skip target
}
```

## Why Interface{} for Args and Results?

Using `[]any` provides flexibility but sacrifices type safety. This is a deliberate trade-off:

### Benefits:
- Supports any function signature
- Allows modification of arguments and results
- Keeps the API simple and unified
- Enables Around advice to set results without knowing types

### Drawbacks:
- Loss of compile-time type safety
- Requires careful type assertions
- Runtime panics possible on incorrect type assertions

### Best Practice:
Always type-assert carefully when accessing Args/Results:
```go
// Safe type assertion
if arg, ok := c.Args[0].(string); ok {
    // Use arg safely
} else {
    // Handle type mismatch
}
```

## Metadata Design: Communication Between Advice

The Metadata field enables communication between different advice functions:

```go
// Authentication advice sets user info
c.Metadata["user"] = authenticatedUser

// Authorization advice reads user info
user := c.Metadata["user"].(*User)
```

### Considerations:
- Metadata is untyped and unchecked - mistakes lead to runtime panics
- No validation of key names or value types
- Potential for key collisions between different advice functions
- Memory leak possible if metadata is not cleaned up

### Best Practices:
- Use consistent naming conventions for keys
- Document metadata keys and expected types
- Consider using typed accessors to encapsulate type assertions

## Function Identity

The FunctionName field provides identity for the executing function:
- Used for logging and debugging
- Helps advice functions know which function they're operating on
- Enables conditional behavior based on function name
- Useful for metrics and monitoring

## Execution Control

The Skipped field allows Around advice to control execution flow:
- When set to true, the target function is skipped
- AfterReturning advice may still execute depending on other conditions
- Enables caching and other optimization patterns

## Error and Panic Handling

The Context handles both errors and panics:
- Error field captures function errors
- PanicValue field captures panic values
- Both can be modified by advice functions
- Enables error transformation and recovery patterns

## Memory Management

The Context is created per function call and destroyed when the call completes:
- Args and Results slices grow as needed
- Metadata map is initialized lazily
- No shared state between function calls
- Garbage collected after function execution

## Thread Safety

Context objects are not shared between goroutines:
- Each function call gets its own context
- No synchronization needed within a context
- Safe for concurrent execution of the same function from different goroutines

## Limitations

- No type safety for Args/Results/Metadata
- Potential for runtime panics on type mismatches
- Memory overhead for each function call
- No validation of metadata usage

The Context object is the backbone of the AOP system, enabling all the communication and state management needed for effective cross-cutting concern implementation.
# Execution Engine

The `executeWithAdvice` function is the core orchestration logic that brings together all components of the AOP system. It manages the complex flow of executing advice in the correct order around the target function.

## Core Execution Flow

```go
func executeWithAdvice(functionName string, targetFn func(*Context), args ...any) *Context {
    // Get advice chain from registry
    chain, err := GetAdviceChain(functionName)
    if err != nil {
        // No advice registered, just execute target function
        c := NewContext(functionName, args...)
        targetFn(c)
        return c
    }

    // Create execution context
    c := NewContext(functionName, args...)

    // Defer After advice (always runs)
    defer func() {
        _ = chain.ExecuteAfter(c)
    }()

    // Defer panic recovery and AfterThrowing advice
    defer func() {
        if r := recover(); r != nil {
            c.PanicValue = r
            _ = chain.ExecuteAfterThrowing(c)
            panic(r)  // Re-panic to maintain original behavior
        }
    }()

    // Execute Before advice
    if err := chain.ExecuteBefore(c); err != nil {
        panic(fmt.Errorf("before advice failed: %w", err))
    }

    // Execute Around advice (if any)
    if chain.HasAround() {
        if err := chain.ExecuteAround(c); err != nil {
            panic(fmt.Errorf("around advice failed: %w", err))
        }
        // If Around advice sets Skipped, don't execute target function
        if c.Skipped {
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
```

## Execution Flow Analysis

### 1. Setup Phase
- Get advice chain from registry
- Create execution context with arguments
- If no advice is registered, execute target function directly

### 2. Defer Setup
- **After advice**: Deferred to ensure it always runs
- **Panic recovery**: Captures panics and executes AfterThrowing advice

### 3. Before Execution
- Execute all Before advice in priority order
- If any Before advice fails, execution stops with panic

### 4. Around Decision
- Execute Around advice in priority order
- Check if execution should be skipped based on `c.Skipped`
- If skipped, may execute AfterReturning and return early

### 5. Target Execution
- Execute the original target function
- Function results and errors are captured in context

### 6. Post-Execution
- Execute AfterReturning if no error/panic
- After advice runs via defer (always)
- AfterThrowing runs via defer if panic occurred

## Error Handling Philosophy

### Fail-Fast for Before/Around
- Before and Around advice errors cause immediate panic
- Prevents execution of target function with invalid state
- Ensures validation and preparation steps succeed

### Graceful for After
- After advice errors are generally ignored
- Cleanup should not fail and shouldn't affect main function
- May be logged but doesn't change execution flow

### Panic Preservation
- Panics from target functions are captured and re-thrown
- AfterThrowing advice runs before re-panicking
- Maintains original panic behavior for calling code

## Deferral Strategy

The execution engine uses defer strategically:
- **After advice**: Deferred to guarantee execution
- **Panic handling**: Deferred to catch any panics from target or advice
- **Order matters**: Defer statements execute in LIFO order

## Around Advice Special Handling

Around advice has unique behavior:
- Can skip target function execution entirely
- Can modify results without calling target
- Multiple Around advice functions execute in sequence
- Last Around advice to run "wraps" the others

## Performance Considerations

### Time Complexity
- **Registry lookup**: O(1) average case
- **Context creation**: O(1) for basic fields, O(m) for m arguments
- **Advice execution**: O(p log p + a) where p is advice with same priority, a is total advice count
- **Deferred execution**: Happens after main execution

### Memory Allocation
- One Context per function call
- Temporary slices for sorted advice
- Stack frames for deferred functions

## Thread Safety

The execution engine is thread-safe because:
- Context objects are not shared between goroutines
- Registry access is synchronized
- Advice chains are immutable after creation
- No shared mutable state between concurrent executions

## Exception Handling

The execution engine handles multiple types of exceptional conditions:
- **Advice errors**: Handled according to error philosophy
- **Target function panics**: Recovered and re-thrown
- **System panics**: Preserved to maintain program behavior

The execution engine is the heart of the AOP system, coordinating the complex dance of advice execution around the target function while maintaining proper error handling and guarantees.
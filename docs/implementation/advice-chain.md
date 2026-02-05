# Advice Chain

The AdviceChain is where the execution orchestration happens - it determines the order and conditions under which advice executes. It's the core component that manages the actual AOP behavior for each registered function.

## Chain Structure

```go
type AdviceChain struct {
    before         []Advice  // Before advice stored in priority order
    after          []Advice  // After advice stored in priority order
    around         []Advice  // Around advice stored in priority order
    afterReturning []Advice  // AfterReturning advice stored in priority order
    afterThrowing  []Advice  // AfterThrowing advice stored in priority order
}
```

## Why Separate Storage for Each Advice Type?

Each advice type has different execution semantics:

- **Before**: Always executes first (unless Around skips)
- **Around**: Can conditionally execute the target function
- **After**: Always executes (cleanup guarantee)
- **AfterReturning**: Only on success
- **AfterThrowing**: Only on panic

Separating them allows for efficient execution without type checking during runtime.

## Priority System Implementation

```go
// Sorting happens during execution, not during addition
sort.Slice(sortedAdviceList, func(i, j int) bool {
    return sortedAdviceList[i].Priority > sortedAdviceList[j].Priority
})
```

### Design Choice
Sorting during execution allows dynamic priority changes but adds O(n log n) complexity per execution. This trade-off was chosen because:

- Advice addition happens infrequently (during setup)
- Function execution happens frequently (runtime)
- Dynamic priority changes might be needed in some scenarios

## Execution Flow

The advice chain controls the execution flow for each function call:

1. **Before advice**: Executed in priority order (highest first)
2. **Around advice**: Executed in priority order, can skip target function
3. **Target function**: Executes unless skipped by Around advice
4. **AfterReturning**: Executes only if no error/panic
5. **AfterThrowing**: Executes only if panic occurred
6. **After advice**: Always executes (cleanup)

## Limitations of the Current Design

### Fixed Advice Types
Only the 5 predefined types are supported. This limits flexibility but keeps the implementation simple and predictable.

### Priority Conflicts
Multiple teams or modules might use overlapping priority ranges, leading to unexpected execution orders. Coordination is needed for large applications.

### No Conditional Advice
Advice always executes based on type, not conditions. This simplifies the implementation but reduces flexibility.

## Adding Advice

Advice is added to the appropriate slice based on its type:

```go
func (ac *AdviceChain) Add(advice Advice) {
    switch advice.Type {
    case Before:
        ac.before = append(ac.before, advice)
    case After:
        ac.after = append(ac.after, advice)
    // ... other cases
    }
}
```

This approach is simple and efficient, but requires type checking during addition.

## Performance Considerations

- **Addition**: O(1) - simple slice append
- **Execution**: O(n log n) where n is the number of advice of the same type (due to sorting)
- **Memory**: O(n) where n is the total number of advice for the function

## Thread Safety

The advice chain is immutable after creation (aside from the initial additions), so it's safe for concurrent access. No locking is needed during execution.

## Extensibility

The current design makes it relatively easy to add new advice types by:
1. Adding a new type to the AdviceType enum
2. Adding a new slice to the AdviceChain struct
3. Adding a new case to the Add method
4. Adding a new execution method

The advice chain is the heart of the AOP system, orchestrating the complex interactions between different cross-cutting concerns while maintaining predictable execution order.
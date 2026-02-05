# Registry System

The registry is the central hub that maintains associations between function names and their advice chains. It's responsible for managing the global state of registered functions.

## Registry Structure

```go
type Registry struct {
    mu      sync.RWMutex      // Thread safety for concurrent access
    entries map[string]*AdviceChain  // Function name → AdviceChain mapping
}
```

## Design Decisions

### Thread Safety
The registry uses `sync.RWMutex` because we expect more read operations (function calls) than write operations (registrations). This allows multiple goroutines to access different advice chains concurrently while preventing race conditions during registration.

### Global Singleton
A default global registry is provided for convenience, but the system supports multiple registries for different use cases.

### String Keys
Human-readable function names are used as keys, making debugging and logging more intuitive.

## Implementation Details

### Registration Process
When a function is registered:
1. Acquire write lock
2. Check for duplicates
3. Create new advice chain
4. Store in the entries map
5. Release lock

### Lookup Process
During function execution:
1. Acquire read lock
2. Find the advice chain by name
3. Return the chain (or error if not found)
4. Release lock

## Limitations and Considerations

### Memory Growth
The registry grows with each registered function and never shrinks. This means:
- Memory usage increases monotonically
- Applications with dynamic function registration should monitor memory
- Consider using local registries for temporary or test scenarios

### Name Collisions
Duplicate function names cause registration errors. This prevents accidental overwrites but requires careful naming conventions.

### Startup Cost
All registration should happen during application initialization to avoid contention during runtime.

## When to Use Local Registries

Local registries are beneficial for:

- **Testing isolation**: Each test can have its own registry
- **Microservices**: Different services may have different AOP requirements
- **Dynamic scenarios**: When function registration happens at runtime
- **Multi-tenant applications**: Isolating AOP configurations

Example:
```go
// Use local registries for:
// - Testing isolation
// - Microservices with different AOP requirements
// - Dynamic function registration scenarios

localRegistry := aspect.NewRegistry()
localRegistry.Register("MyFunc", myFunc)
```

## Performance Characteristics

- **Registration**: O(1) average case with hash map lookup
- **Lookup**: O(1) average case during function calls
- **Concurrency**: Optimized for read-heavy workloads
- **Memory**: O(n) where n is the number of registered functions

## Thread Safety Guarantees

The registry provides thread-safe operations:
- Multiple goroutines can safely call registered functions concurrently
- Registration operations are serialized
- No race conditions between registration and function calls

Understanding the registry system helps you plan your application's AOP architecture and avoid potential issues with memory usage and naming conflicts.
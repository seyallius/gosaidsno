# Examples

This section provides real-world examples of how to use gosaidsno for various common patterns and use cases.

## Available Examples

The examples directory contains several complete, runnable examples demonstrating different aspects of gosaidsno:

### 1. Basic Usage (`examples/01_basic_usage`)

Demonstrates core AOP features with real-world scenarios:

- Request/response logging
- Method execution timing
- Input validation (Before advice)
- Success-only actions (AfterReturning)

### 2. Caching Pattern (`examples/02_caching_pattern`)

Shows how to implement caching with Around advice:

- Database query caching
- API response caching
- TTL-based cache expiration
- Cache hit/miss metrics

### 3. Authentication (`examples/03_authentication`)

Demonstrates authentication and authorization patterns:

- Token validation
- Session management
- Role-based access control (RBAC)
- Audit logging

### 4. Circuit Breaker (`examples/04_circuit_breaker`)

Shows fault tolerance with circuit breaker pattern:

- External API fault tolerance
- Database connection protection
- Service degradation handling

### 5. Retry Pattern (`examples/05_retry_pattern`)

Implements automatic retries with exponential backoff:

- Transient failure handling
- Network request retries
- Exponential backoff calculation

## Running Examples

You can run any example with:

```bash
# Basic usage - logging, timing, validation
go run examples/01_basic_usage/main.go

# Caching with Around advice
go run examples/02_caching_pattern/main.go

# Authentication and authorization
go run examples/03_authentication/main.go

# Circuit breaker for fault tolerance
go run examples/04_circuit_breaker/main.go

# Retry with exponential backoff
go run examples/05_retry_pattern/main.go
```

## Example Structure

All examples follow the same pattern:

```go
// 1. Setup AOP once at startup
func setupAOP() {
aspect.MustRegister("FunctionName")
aspect.MustAddAdvice("FunctionName", /* ... */)
}

// 2. Implement business logic
func businessLogicImpl(args...) result {
// Pure business logic
}

// 3. Wrap functions (once, during initialization)
var BusinessLogic = aspect.Wrap*("FunctionName", businessLogicImpl)

// 4. Use normally throughout application
func main() {
setupAOP() // Once

// Use wrapped functions
result := BusinessLogic(args...)
}
```

## Advice Execution Order

Understanding the execution order is crucial for complex examples:

1. **Before** (high priority → low priority)
2. **Around** (high priority → low priority) - can skip step 3
3. Target function (only if not skipped by Around advice)
4. **AfterReturning** (only if success, high priority → low)
5. **AfterThrowing** (only if panic, high priority → low)
6. **After** (always runs, high priority → low)

## Common Metadata Keys

Examples use these metadata conventions for communication between advice:

- `startTime` - For timing (type: `time.Time`)
- `userID` - For auth/audit (type: `string`)
- `role` - For authorization (type: `string`)
- `attempt` - For retry logic (type: `int`)
- `maxRetries` - Retry configuration (type: `int`)

## Performance Considerations

Each example includes:

- Timing measurements
- Success/failure counts
- Performance impact of advice
- Cache hit rates (where applicable)

For more details on any specific example, see the [full examples documentation](./examples/README.md).
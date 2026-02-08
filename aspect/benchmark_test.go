package aspect

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// Benchmark Results:
//
// with pool
//Benchmark_NoAdvice-6                   	 2842268	     397.9 ns/op	     288 B/op	       8 allocs/op
//Benchmark_BeforeAdvice-6               	 2963404	     402.0 ns/op	     312 B/op	       9 allocs/op
//Benchmark_AroundAdvice-6               	 2891487	     413.7 ns/op	     312 B/op	       9 allocs/op
//Benchmark_HotPath_MicroserviceCall-6   	 83270	     	 31331 ns/op	     752 B/op	      15 allocs/op
//Benchmark_HotPath_CacheHeavy-6         	 125724	     	 9776 ns/op	     	 418 B/op	      13 allocs/op
//Benchmark_HighConcurrency-6            	 3340183	     365.6 ns/op	     750 B/op	      16 allocs/op
//Benchmark_MixedWorkload-6              	 36288	     	 35478 ns/op	     638 B/op	      11 allocs/op
//Benchmark_MetadataHeavy-6              	 2199160	     556.1 ns/op	     1416 B/op	      21 allocs/op
//Benchmark_PoolEffectiveness/WithPool-6 	 3562045	     328.3 ns/op	     744 B/op	      16 allocs/op
//Benchmark_PoolEffectiveness/WithoutPool-6  3926685	     383.7 ns/op	     744 B/op	      15 allocs/op
//Benchmark_RealWorldExample-6               36604	     	 34777 ns/op	     906 B/op	      20 allocs/op
//
// without pool
//Benchmark_NoAdvice-6                   	 3239707	     423.1 ns/op	     304 B/op	       9 allocs/op
//Benchmark_BeforeAdvice-6               	 3296284	     372.0 ns/op	     311 B/op	       8 allocs/op
//Benchmark_AroundAdvice-6               	 2812153	     457.4 ns/op	     312 B/op	       8 allocs/op
//Benchmark_HotPath_MicroserviceCall-6   	  192961	     30095 ns/op	     752 B/op	      15 allocs/op
//Benchmark_HotPath_CacheHeavy-6         	  112795	     11670 ns/op	     418 B/op	      13 allocs/op
//Benchmark_HighConcurrency-6            	 2865145	     588.7 ns/op	     749 B/op	      16 allocs/op
//Benchmark_MixedWorkload-6              	   34382	     34703 ns/op	     638 B/op	      11 allocs/op
//Benchmark_MetadataHeavy-6              	 1574728	     795.4 ns/op	    1416 B/op	      20 allocs/op
//Benchmark_PoolEffectiveness/WithPool-6 	 2073848	     558.6 ns/op	     744 B/op	      16 allocs/op
//Benchmark_PoolEffectiveness/WithoutPool-6  2684281	     548.3 ns/op	     744 B/op	      15 allocs/op
//Benchmark_RealWorldExample-6                 28915	     37826 ns/op	     912 B/op	      20 allocs/op

// -------------------------------------------- Setup Functions --------------------------------------------

// createRegistryWithLoggingAdvice creates a registry with realistic logging advice
func createRegistryWithLoggingAdvice() *Registry {
	reg := NewRegistry()
	reg.MustRegister("businessLogic")

	// Common pattern: Logging + Metrics + Validation
	reg.MustAddAdvice("businessLogic", Advice{
		Type:     Before,
		Priority: 100,
		Handler: func(c *Context) error {
			c.SetMetadataVal("startTime", time.Now())
			c.SetMetadataVal("traceId", generateTraceID())
			return nil
		},
	})

	reg.MustAddAdvice("businessLogic", Advice{
		Type:     After,
		Priority: 100,
		Handler: func(c *Context) error {
			// Log execution time
			if start, ok := c.GetMetadataVal("startTime"); ok {
				duration := time.Since(start.(time.Time))
				_ = duration // In real code, log this
			}
			return nil
		},
	})

	return reg
}

func generateTraceID() string {
	return "trace-12345"
}

// createRegistryWithCaching creates a registry with caching pattern
func createRegistryWithCaching() *Registry {
	reg := NewRegistry()
	reg.MustRegister("dataFetch")

	cache := sync.Map{}

	reg.MustAddAdvice("dataFetch", Advice{
		Type:     Around,
		Priority: 200, // High priority to check cache first
		Handler: func(c *Context) error {
			key := c.Args[0].(string)
			if val, ok := cache.Load(key); ok {
				c.SetResult(0, val)
				c.Skipped = true
			}
			return nil
		},
	})

	reg.MustAddAdvice("dataFetch", Advice{
		Type:     AfterReturning,
		Priority: 100,
		Handler: func(c *Context) error {
			if !c.Skipped {
				key := c.Args[0].(string)
				cache.Store(key, c.Results[0])
			}
			return nil
		},
	})

	return reg
}

// createRegistryWithErrorHandling creates a registry with error handling pattern
func createRegistryWithErrorHandling() *Registry {
	reg := NewRegistry()
	reg.MustRegister("apiCall")

	reg.MustAddAdvice("apiCall", Advice{
		Type:     AfterThrowing,
		Priority: 100,
		Handler: func(c *Context) error {
			// Log panic for monitoring
			_ = c.PanicValue
			return nil
		},
	})

	reg.MustAddAdvice("apiCall", Advice{
		Type:     After,
		Priority: 50,
		Handler: func(c *Context) error {
			if c.Error != nil {
				// Convert errors to consistent format
				// In real code: c.Error = formatError(c.Error)
			}
			return nil
		},
	})

	return reg
}

// -------------------------------------------- Benchmark Tests --------------------------------------------

// Benchmark_NoAdvice
//
//	NoPool	3118063	       351.7 ns/op	     368 B/op	      11 allocs/op
//	Pool	4296924	       265.2 ns/op	      96 B/op	       6 allocs/op
func Benchmark_NoAdvice(b *testing.B) {
	reg := NewRegistry()

	fn := func(a int) int {
		return a + 1
	}

	wrapped := Wrap1R(reg, "fn", fn)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = wrapped(i)
	}
}

// Benchmark_BeforeAdvice
//
//	NoPool	3458373	       330.1 ns/op	     344 B/op	      10 allocs/op
//	Pool	5014700	       244.2 ns/op	     104 B/op	       6 allocs/op
func Benchmark_BeforeAdvice(b *testing.B) {
	reg := NewRegistry()

	_ = reg.Register("fn")
	if err := reg.AddAdvice("fn", Advice{
		Type: Before,
		Handler: func(c *Context) error {
			return nil
		},
	}); err != nil {
		b.Fatalf("failed to add advice: %v", err)
	}

	fn := func(a int) int {
		return a + 1
	}

	wrapped := Wrap1R(reg, "fn", fn)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = wrapped(i)
	}
}

// Benchmark_AroundAdvice
//
//	NoPool	3348843	       341.4 ns/op	     344 B/op	      10 allocs/op
//	Pool	4946163	       244.3 ns/op	     104 B/op	       6 allocs/op
func Benchmark_AroundAdvice(b *testing.B) {
	reg := NewRegistry()

	_ = reg.Register("fn")
	if err := reg.AddAdvice("fn", Advice{
		Type: Around,
		Handler: func(c *Context) error {
			return nil
		},
	}); err != nil {
		b.Fatalf("failed to add advice: %v", err)
	}

	fn := func(a int) int {
		return a + 1
	}

	wrapped := Wrap1R(reg, "fn", fn)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = wrapped(i)
	}
}

// Benchmark_HotPath_MicroserviceCall simulates a hot path like an API endpoint
func Benchmark_HotPath_MicroserviceCall(b *testing.B) {
	reg := createRegistryWithLoggingAdvice()

	// Simulate business logic that might be called millions of times
	businessLogic := func(userID int, requestData string) (string, error) {
		// Simulate some processing
		time.Sleep(10 * time.Microsecond)
		if userID <= 0 {
			return "", errors.New("invalid user ID")
		}
		return "processed:" + requestData, nil
	}

	wrapped := Wrap2RE(reg, "businessLogic", businessLogic)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			userID := i%100 + 1 // Mix of valid and invalid IDs
			_, _ = wrapped(userID, "request-data")
		}
	})
}

// Benchmark_HotPath_CacheHeavy simulates cache-heavy workload (common in web servers)
func Benchmark_HotPath_CacheHeavy(b *testing.B) {
	reg := createRegistryWithCaching()

	// Simulate expensive database/API call
	expensiveFetch := func(key string) (string, error) {
		time.Sleep(50 * time.Microsecond) // Simulate network/db latency
		return "data-for-" + key, nil
	}

	wrapped := Wrap1RE(reg, "dataFetch", expensiveFetch)

	// Pre-warm cache with some keys
	for i := 0; i < 100; i++ {
		wrapped("key-" + string(rune(i)))
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			// 80% cache hits, 20% cache misses (realistic ratio)
			key := "key-"
			if i%5 == 0 { // 20% miss
				key += "miss-" + string(rune(i))
			} else { // 80% hit
				key += string(rune(i % 100))
			}
			_, _ = wrapped(key)
		}
	})
}

// Benchmark_HighConcurrency simulates high concurrency scenario (e.g., web server)
func Benchmark_HighConcurrency(b *testing.B) {
	reg := NewRegistry()
	reg.MustRegister("concurrentOperation")

	// Add multiple advice types (realistic for production)
	reg.MustAddAdvice("concurrentOperation", Advice{
		Type: Before,
		Handler: func(c *Context) error {
			c.SetMetadataVal("reqId", time.Now().UnixNano())
			return nil
		},
	})

	reg.MustAddAdvice("concurrentOperation", Advice{
		Type: Around,
		Handler: func(c *Context) error {
			// Rate limiting/validation logic
			return nil
		},
	})

	reg.MustAddAdvice("concurrentOperation", Advice{
		Type: After,
		Handler: func(c *Context) error {
			// Cleanup/metrics
			return nil
		},
	})

	operation := func(x int) int {
		// Simulate CPU work
		result := x
		for i := 0; i < 100; i++ {
			result = (result * 31) % 1000
		}
		return result
	}

	wrapped := Wrap1R(reg, "concurrentOperation", operation)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			_ = wrapped(i)
		}
	})
}

// Benchmark_MixedWorkload simulates mixed short and long operations
func Benchmark_MixedWorkload(b *testing.B) {
	reg := NewRegistry()

	// Register multiple functions with different advice configurations
	funcs := []string{"fastPath", "slowPath", "errorPath"}
	for _, name := range funcs {
		reg.MustRegister(FuncKey(name))

		// Common advice for all
		reg.MustAddAdvice(FuncKey(name), Advice{
			Type: Before,
			Handler: func(c *Context) error {
				c.SetMetadataVal("timestamp", time.Now())
				return nil
			},
		})

		if name == "slowPath" {
			reg.MustAddAdvice(FuncKey(name), Advice{
				Type: Around,
				Handler: func(c *Context) error {
					// Timeout logic for slow operations
					return nil
				},
			})
		}

		if name == "errorPath" {
			reg.MustAddAdvice(FuncKey(name), Advice{
				Type: AfterThrowing,
				Handler: func(c *Context) error {
					// Error recovery
					return nil
				},
			})
		}
	}

	// Different types of operations
	fastOp := func(x int) int { return x * 2 }
	slowOp := func(x int) int {
		time.Sleep(100 * time.Microsecond)
		return x * 3
	}
	errorOp := func(x int) int {
		if x%10 == 0 {
			panic("simulated error")
		}
		return x * 4
	}

	fastWrapped := Wrap1R(reg, "fastPath", fastOp)
	slowWrapped := Wrap1R(reg, "slowPath", slowOp)
	errorWrapped := Wrap1R(reg, "errorPath", errorOp)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			// Mix of operations: 70% fast, 20% slow, 10% error
			switch i % 10 {
			case 0, 1, 2, 3, 4, 5, 6: // 70% fast
				_ = fastWrapped(i)
			case 7, 8: // 20% slow
				_ = slowWrapped(i)
			case 9: // 10% error (with panic recovery)
				func() {
					defer func() { _ = recover() }()
					_ = errorWrapped(i)
				}()
			}
		}
	})
}

// Benchmark_MetadataHeavy simulates advice that heavily uses metadata (common for tracing)
func Benchmark_MetadataHeavy(b *testing.B) {
	reg := NewRegistry()
	reg.MustRegister("tracingHeavy")

	// Simulate distributed tracing with lots of metadata
	reg.MustAddAdvice("tracingHeavy", Advice{
		Type: Before,
		Handler: func(c *Context) error {
			// Set multiple metadata values (common in tracing)
			c.SetMetadataVal("traceId", "trace-123")
			c.SetMetadataVal("spanId", "span-456")
			c.SetMetadataVal("parentId", "parent-789")
			c.SetMetadataVal("service", "user-service")
			c.SetMetadataVal("endpoint", "/api/users")
			c.SetMetadataVal("startTime", time.Now())
			c.SetMetadataVal("correlationId", "corr-abc")
			c.SetMetadataVal("userId", c.Args[0])
			return nil
		},
	})

	reg.MustAddAdvice("tracingHeavy", Advice{
		Type: Around,
		Handler: func(c *Context) error {
			// Add more metadata
			c.SetMetadataVal("aroundStart", time.Now())
			return nil
		},
	})

	reg.MustAddAdvice("tracingHeavy", Advice{
		Type: After,
		Handler: func(c *Context) error {
			// Read all metadata for logging
			_, _ = c.GetMetadataVal("traceId")
			_, _ = c.GetMetadataVal("spanId")
			_, _ = c.GetMetadataVal("parentId")
			_, _ = c.GetMetadataVal("service")
			_, _ = c.GetMetadataVal("endpoint")
			_, _ = c.GetMetadataVal("startTime")
			_, _ = c.GetMetadataVal("correlationId")
			_, _ = c.GetMetadataVal("userId")
			_, _ = c.GetMetadataVal("aroundStart")
			return nil
		},
	})

	operation := func(userID int) string {
		return "user-data"
	}

	wrapped := Wrap1R(reg, "tracingHeavy", operation)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			_ = wrapped(i)
		}
	})
}

// Benchmark_PoolEffectiveness isolates the pool effect
func Benchmark_PoolEffectiveness(b *testing.B) {
	// Test with and without pool
	b.Run("WithPool", func(b *testing.B) {
		// Enable pooling
		enableContextPooling = true
		defer func() { enableContextPooling = false }()

		benchmarkPoolScenario(b)
	})

	b.Run("WithoutPool", func(b *testing.B) {
		// Disable pooling
		enableContextPooling = false
		defer func() { enableContextPooling = true }()

		benchmarkPoolScenario(b)
	})
}

func benchmarkPoolScenario(b *testing.B) {
	reg := NewRegistry()
	reg.MustRegister("poolTest")

	// Add moderate amount of advice
	reg.MustAddAdvice("poolTest", Advice{
		Type: Before,
		Handler: func(c *Context) error {
			c.SetMetadataVal("val1", 1)
			c.SetMetadataVal("val2", 2)
			return nil
		},
	})

	reg.MustAddAdvice("poolTest", Advice{
		Type: Around,
		Handler: func(c *Context) error {
			c.SetMetadataVal("val3", 3)
			return nil
		},
	})

	reg.MustAddAdvice("poolTest", Advice{
		Type: After,
		Handler: func(c *Context) error {
			_, _ = c.GetMetadataVal("val1")
			_, _ = c.GetMetadataVal("val2")
			_, _ = c.GetMetadataVal("val3")
			return nil
		},
	})

	operation := func(x int) int {
		// Quick operation
		return x * 7
	}

	wrapped := Wrap1R(reg, "poolTest", operation)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			_ = wrapped(i)
		}
	})
}

// Benchmark_RealWorldExample e-commerce checkout flow
func Benchmark_RealWorldExample(b *testing.B) {
	reg := NewRegistry()
	reg.MustRegister("checkout")

	// Real-world advice patterns
	reg.MustAddAdvice("checkout", Advice{
		Type:     Before,
		Priority: 300,
		Handler: func(c *Context) error {
			// Authentication & authorization
			c.SetMetadataVal("userId", c.Args[0])
			c.SetMetadataVal("startTime", time.Now())
			return nil
		},
	})

	reg.MustAddAdvice("checkout", Advice{
		Type:     Before,
		Priority: 200,
		Handler: func(c *Context) error {
			// Request validation
			amount := c.Args[1].(float64)
			if amount <= 0 {
				return errors.New("invalid amount")
			}
			return nil
		},
	})

	reg.MustAddAdvice("checkout", Advice{
		Type:     Around,
		Priority: 100,
		Handler: func(c *Context) error {
			// Transaction management
			c.SetMetadataVal("transactionId", "txn-"+string(rune(time.Now().UnixNano())))
			return nil
		},
	})

	reg.MustAddAdvice("checkout", Advice{
		Type:     AfterReturning,
		Priority: 100,
		Handler: func(c *Context) error {
			// Audit logging
			userId, _ := c.GetMetadataVal("userId")
			txnId, _ := c.GetMetadataVal("transactionId")
			startTime, _ := c.GetMetadataVal("startTime")
			duration := time.Since(startTime.(time.Time))
			_ = userId
			_ = txnId
			_ = duration
			return nil
		},
	})

	reg.MustAddAdvice("checkout", Advice{
		Type:     AfterThrowing,
		Priority: 100,
		Handler: func(c *Context) error {
			// Error recovery/notification
			_ = c.PanicValue
			return nil
		},
	})

	// Simulate checkout logic
	checkout := func(userID int, amount float64) (string, error) {
		// Complex business logic
		time.Sleep(20 * time.Microsecond)

		if amount > 10000 {
			return "", errors.New("amount too high")
		}

		return "order-created", nil
	}

	wrapped := Wrap2RE(reg, "checkout", checkout)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			// Mix of valid and invalid requests
			userID := i%1000 + 1
			amount := float64((i%2000)*10 + 1)
			_, _ = wrapped(userID, amount)
		}
	})
}

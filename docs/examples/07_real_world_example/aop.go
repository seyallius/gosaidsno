package main

import (
	"errors"
	"log"
	"time"

	"github.com/seyallius/gosaidno/aspect"
	"github.com/seyallius/gosaidno/docs/examples/utils"
)

// -------------------------------------------- AOP Setup (Centralized) --------------------------------------------

// setupAOP configures all cross-cutting concerns using the fluent API
func setupAOP() {
	log.Println("=== Setting up AOP Cross-Cutting Concerns ===")

	// 1. Logging setup using fluent API
	setupLogging()

	// 2. Timing setup using fluent API
	setupTiming()

	// 3. Validation setup using fluent API
	setupValidation()

	// 4. Caching setup using fluent API
	setupCaching()

	// 5. Error handling setup using fluent API
	setupErrorHandling()

	log.Println("=== AOP Setup Complete ===")
	log.Println()
}

func setupLogging() {
	log.Println("   📝 Setting up logging advice...")

	// Apply logging to all service methods
	aspect.For("UserService.GetUser").
		WithBefore(func(c *aspect.Context) error {
			utils.LogBefore(c, 100, "LOGGING")
			username := c.Args[0].(string)
			log.Printf("   📝 [LOG] Starting GetUser for username: %s", username)
			return nil
		}).
		WithAfter(func(c *aspect.Context) error {
			utils.LogAfter(c, 100, "LOGGING")
			username := c.Args[0].(string)
			status := "SUCCESS"
			if c.Error != nil {
				status = "FAILED"
			}
			log.Printf("   📝 [LOG] Completed GetUser for %s - Status: %s", username, status)
			if c.Error != nil {
				log.Printf("   ❌ Error: %v", c.Error)
			}
			return nil
		})

	aspect.For("UserService.CreateUser").
		WithBefore(func(c *aspect.Context) error {
			utils.LogBefore(c, 100, "LOGGING")
			user := c.Args[0].(*User)
			log.Printf("   📝 [LOG] Starting CreateUser for user: %s", user.Username)
			return nil
		}).
		WithAfter(func(c *aspect.Context) error {
			utils.LogAfter(c, 100, "LOGGING")
			user := c.Args[0].(*User)
			status := "SUCCESS"
			if c.Error != nil {
				status = "FAILED"
			}
			log.Printf("   📝 [LOG] Completed CreateUser for %s - Status: %s", user.Username, status)
			return nil
		})

	aspect.For("OrderService.CreateOrder").
		WithBefore(func(c *aspect.Context) error {
			utils.LogBefore(c, 100, "LOGGING")
			userID := c.Args[0].(string)
			amount := c.Args[1].(float64)
			log.Printf("   📝 [LOG] Starting CreateOrder for user: %s, amount: %.2f", userID, amount)
			return nil
		}).
		WithAfter(func(c *aspect.Context) error {
			utils.LogAfter(c, 100, "LOGGING")
			userID := c.Args[0].(string)
			amount := c.Args[1].(float64)
			status := "SUCCESS"
			if c.Error != nil {
				status = "FAILED"
			}
			log.Printf("   📝 [LOG] Completed CreateOrder for %s, amount: %.2f - Status: %s", userID, amount, status)
			return nil
		})
}

func setupTiming() {
	log.Println("   ⏱️  Setting up timing advice...")

	// Apply timing to all service methods
	for _, funcName := range []aspect.FuncKey{
		"UserService.GetUser", "UserService.CreateUser",
		"OrderService.CreateOrder", "OrderService.GetOrder"} {

		aspect.For(funcName).
			WithBefore(func(c *aspect.Context) error {
				utils.LogBefore(c, 90, "TIMING")
				c.Metadata["start_time"] = time.Now()
				log.Printf("   ⏱️  [TIMING] Started timer for %s", c.FunctionName)
				return nil
			}).
			WithAfter(func(c *aspect.Context) error {
				utils.LogAfter(c, 90, "TIMING")
				startTime, ok := c.Metadata["start_time"].(time.Time)
				if !ok {
					return nil
				}
				duration := time.Since(startTime)
				log.Printf("   ⏱️  [PERF] %s took %v", c.FunctionName, duration)
				return nil
			})
	}
}

func setupValidation() {
	log.Println("   ✅ Setting up validation advice...")

	// Validation for CreateUser
	aspect.For("UserService.CreateUser").
		WithBefore(func(c *aspect.Context) error {
			utils.LogBefore(c, 110, "VALIDATION")
			user := c.Args[0].(*User)

			if user.Username == "" {
				log.Printf("   ❌ [VALIDATE] Username cannot be empty")
				return errors.New("username cannot be empty")
			}
			if user.Email == "" {
				log.Printf("   ❌ [VALIDATE] Email cannot be empty")
				return errors.New("email cannot be empty")
			}
			if user.Username == "admin" {
				log.Printf("   ❌ [VALIDATE] Username 'admin' is reserved")
				return errors.New("username 'admin' is reserved")
			}

			log.Printf("   ✅ [VALIDATE] User validation passed for: %s", user.Username)
			return nil
		})

	// Validation for CreateOrder
	aspect.For("OrderService.CreateOrder").
		WithBefore(func(c *aspect.Context) error {
			utils.LogBefore(c, 110, "VALIDATION")
			userID := c.Args[0].(string)
			amount := c.Args[1].(float64)

			if userID == "" {
				log.Printf("   ❌ [VALIDATE] UserID cannot be empty")
				return errors.New("userID cannot be empty")
			}
			if amount <= 0 {
				log.Printf("   ❌ [VALIDATE] Amount must be positive")
				return errors.New("amount must be positive")
			}
			if amount > 10000 {
				log.Printf("   ❌ [VALIDATE] Amount exceeds maximum allowed ($10,000)")
				return errors.New("amount exceeds maximum allowed ($10,000)")
			}

			log.Printf("   ✅ [VALIDATE] Order validation passed for user: %s, amount: %.2f", userID, amount)
			return nil
		})
}

func setupCaching() {
	log.Println("   💾 Setting up caching advice...")

	// Simple in-memory cache
	userCache := make(map[string]*User)

	// Around advice for caching GetUser
	aspect.For("UserService.GetUser").
		WithAround(func(c *aspect.Context) error {
			username := c.Args[0].(string)

			// Check cache first
			if cachedUser, exists := userCache[username]; exists {
				log.Printf("   💾 [CACHE] Cache HIT for user: %s", username)
				c.SetResult(0, cachedUser)
				c.Skipped = true // Skip target execution
				return nil
			}

			log.Printf("   🔍 [CACHE] Cache MISS for user: %s", username)
			return nil // Let target execute
		}).
		WithAfterReturning(func(c *aspect.Context) error {
			// Populate cache after successful execution
			username := c.Args[0].(string)
			user := c.Results[0].(*User)
			userCache[username] = user
			log.Printf("   💾 [CACHE] Cached user: %s", username)
			return nil
		})
}

func setupErrorHandling() {
	log.Println("   🚨 Setting up error handling advice...")

	// Apply error handling to all service methods
	for _, funcName := range []aspect.FuncKey{
		"UserService.GetUser", "UserService.CreateUser",
		"OrderService.CreateOrder", "OrderService.GetOrder"} {

		aspect.For(funcName).
			WithAfterThrowing(func(c *aspect.Context) error {
				utils.LogAfterThrowing(c, 100, "ERROR HANDLING")
				log.Printf("   🚨 [ERROR] Function %s panicked: %v", c.FunctionName, c.PanicValue)
				log.Printf("   🔧 [RECOVERY] Recovered from panic in %s", c.FunctionName)
				return nil
			})
	}
}

// -------------------------------------------- Wrapped Service Functions (Global Variables) --------------------------------------------

// These are the wrapped versions of our service methods
// Named following the convention: [Service][Method]Wrapped
var (
	// Create service instances to hold data
	userServiceInstance  = NewUserService()
	orderServiceInstance = NewOrderService()

	// UserService wrapped functions
	UserServiceGetUser = func(username string) (*User, error) {
		builder := aspect.For("UserService.GetUser")
		return aspect.Wrap1RE[string, *User](
			builder.GetRegistry(),
			builder.GetFuncKey(),
			userServiceInstance.GetUser,
		)(username)
	}

	UserServiceCreateUser = func(user *User) error {
		builder := aspect.For("UserService.CreateUser")
		return aspect.Wrap1E[*User](
			builder.GetRegistry(),
			builder.GetFuncKey(),
			userServiceInstance.CreateUser,
		)(user)
	}

	// OrderService wrapped functions
	OrderServiceCreateOrder = func(userID string, amount float64) (*Order, error) {
		builder := aspect.For("OrderService.CreateOrder")
		return aspect.Wrap2RE[string, float64, *Order](
			builder.GetRegistry(),
			builder.GetFuncKey(),
			orderServiceInstance.CreateOrder,
		)(userID, amount)
	}

	OrderServiceGetOrder = func(orderID string) (*Order, error) {
		builder := aspect.For("OrderService.GetOrder")
		return aspect.Wrap1RE[string, *Order](
			builder.GetRegistry(),
			builder.GetFuncKey(),
			orderServiceInstance.GetOrder,
		)(orderID)
	}
)

// Alternative approach: Create a wrapper struct to group related functions
type WrappedServices struct {
	UserService  *WrappedUserService
	OrderService *WrappedOrderService
}

type WrappedUserService struct {
	GetUser    func(string) (*User, error)
	CreateUser func(*User) error
}

type WrappedOrderService struct {
	CreateOrder func(string, float64) (*Order, error)
	GetOrder    func(string) (*Order, error)
}

// Create wrapped services instance
func NewWrappedServices() *WrappedServices {
	// Use the same service instances as the global variables to maintain data consistency
	return &WrappedServices{
		UserService: &WrappedUserService{
			GetUser: func(username string) (*User, error) {
				builder := aspect.For("UserService.GetUser")
				return aspect.Wrap1RE[string, *User](
					builder.GetRegistry(),
					builder.GetFuncKey(),
					userServiceInstance.GetUser,
				)(username)
			},
			CreateUser: func(user *User) error {
				builder := aspect.For("UserService.CreateUser")
				return aspect.Wrap1E[*User](
					builder.GetRegistry(),
					builder.GetFuncKey(),
					userServiceInstance.CreateUser,
				)(user)
			},
		},
		OrderService: &WrappedOrderService{
			CreateOrder: func(userID string, amount float64) (*Order, error) {
				builder := aspect.For("OrderService.CreateOrder")
				return aspect.Wrap2RE[string, float64, *Order](
					builder.GetRegistry(),
					builder.GetFuncKey(),
					orderServiceInstance.CreateOrder,
				)(userID, amount)
			},
			GetOrder: func(orderID string) (*Order, error) {
				builder := aspect.For("OrderService.GetOrder")
				return aspect.Wrap1RE[string, *Order](
					builder.GetRegistry(),
					builder.GetFuncKey(),
					orderServiceInstance.GetOrder,
				)(orderID)
			},
		},
	}
}

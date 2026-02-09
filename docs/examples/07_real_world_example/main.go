// Package main - real_world_example demonstrates a complete real-world usage of gosaidsno
package main

import (
	"fmt"
	"log"
	"time"
)

// -------------------------------------------- Real-World Usage Examples --------------------------------------------

func example1_UserOperations() {
	fmt.Println("\n========== Example 1: User Operations with AOP ==========")

	// Create a user first
	newUser := &User{
		ID:       "user_123",
		Username: "john_doe",
		Email:    "john@example.com",
		Created:  time.Now(),
	}

	log.Println("\n--- Creating user (first call) ---")
	err := UserServiceCreateUser(newUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}

	log.Println("\n--- Retrieving user (cache miss) ---")
	user, err := UserServiceGetUser("john_doe")
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return
	}
	fmt.Printf("\n🎯 Retrieved user: %s (%s)\n", user.Username, user.Email)

	log.Println("\n--- Retrieving same user again (cache hit) ---")
	user2, err := UserServiceGetUser("john_doe")
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return
	}
	fmt.Printf("\n🎯 Retrieved user from cache: %s (%s)\n", user2.Username, user2.Email)
}

func example2_OrderOperations() {
	fmt.Println("\n========== Example 2: Order Operations with AOP ==========")

	log.Println("\n--- Creating order ---")
	order, err := OrderServiceCreateOrder("user_123", 99.99)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return
	}
	fmt.Printf("\n✅ Created order: %s for user: %s, amount: $%.2f\n", order.ID, order.UserID, order.Amount)

	log.Println("\n--- Retrieving order ---")
	retrievedOrder, err := OrderServiceGetOrder(order.ID)
	if err != nil {
		log.Printf("Error getting order: %v", err)
		return
	}
	fmt.Printf("\n🎯 Retrieved order: %s, amount: $%.2f\n", retrievedOrder.ID, retrievedOrder.Amount)
}

func example3_ValidationErrors() {
	fmt.Println("\n========== Example 3: Validation Errors ==========")

	log.Println("\n--- Attempting to create user with invalid data ---")
	invalidUser := &User{
		ID:       "user_456",
		Username: "", // Empty username - should fail validation
		Email:    "invalid@example.com",
		Created:  time.Now(),
	}

	// Use defer/recover to handle the expected panic from validation failure
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("\n❌ Expected validation error caught: %v\n", r)
			}
		}()
		_ = UserServiceCreateUser(invalidUser)
	}()

	log.Println("\n--- Attempting to create order with invalid amount ---")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("\n❌ Expected validation error caught: %v\n", r)
			}
		}()
		_, _ = OrderServiceCreateOrder("user_123", -50.00) // Negative amount - should fail validation
	}()
}

func example4_UsingWrappedServicesStruct() {
	fmt.Println("\n========== Example 4: Using Wrapped Services Struct ==========")

	// Create wrapped services using the struct approach
	wrappedServices := NewWrappedServices()

	// Use the wrapped services
	log.Println("\n--- Creating user via wrapped services struct ---")
	newUser := &User{
		ID:       "user_789",
		Username: "jane_doe",
		Email:    "jane@example.com",
		Created:  time.Now(),
	}

	err := wrappedServices.UserService.CreateUser(newUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}

	log.Println("\n--- Getting user via wrapped services struct ---")
	user, err := wrappedServices.UserService.GetUser("jane_doe")
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return
	}
	fmt.Printf("\n🎯 Retrieved via struct: %s (%s)\n", user.Username, user.Email)
}

// -------------------------------------------- Main Function --------------------------------------------

func main() {
	// Setup AOP once at application startup
	// This is typically done in an init() function or application bootstrap
	setupAOP()

	// Run examples
	example1_UserOperations()
	example2_OrderOperations()
	example3_ValidationErrors()
	example4_UsingWrappedServicesStruct()

	fmt.Println("\n========== All Real-World Examples Complete ==========")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- AOP setup is centralized and done once at startup")
	fmt.Println("- Original business logic remains unchanged")
	fmt.Println("- Cross-cutting concerns (logging, timing, validation, caching) are cleanly separated")
	fmt.Println("- Wrapped functions maintain the same signatures as original functions")
	fmt.Println("- Multiple approaches for organizing wrapped functions (globals vs structs)")
}

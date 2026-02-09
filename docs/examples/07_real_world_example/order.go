package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// -------------------------------------------- Types --------------------------------------------

// OrderService handles order-related operations
type OrderService struct {
	orders map[string]*Order
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string]*Order),
	}
}

// -------------------------------------------- Public Functions --------------------------------------------

// CreateOrder creates a new order
func (os *OrderService) CreateOrder(userID string, amount float64) (*Order, error) {
	log.Printf("   🛒 [BUSINESS] Creating order for user: %s, amount: %.2f", userID, amount)

	// Simulate database insert
	time.Sleep(100 * time.Millisecond) // Simulate DB delay

	order := &Order{
		ID:     fmt.Sprintf("order_%d", time.Now().UnixNano()),
		UserID: userID,
		Amount: amount,
		Date:   time.Now(),
	}
	os.orders[order.ID] = order

	log.Printf("   ✅ [BUSINESS] Order created successfully: %s", order.ID)
	return order, nil
}

// GetOrder retrieves an order by ID
func (os *OrderService) GetOrder(orderID string) (*Order, error) {
	log.Printf("   🛒 [BUSINESS] Retrieving order: %s", orderID)

	// Simulate database lookup
	time.Sleep(60 * time.Millisecond) // Simulate DB delay

	order, exists := os.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}

	log.Printf("   ✅ [BUSINESS] Order retrieved successfully")
	return order, nil
}

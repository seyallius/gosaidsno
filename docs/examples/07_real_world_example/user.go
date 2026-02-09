package main

import (
	"errors"
	"log"
	"time"
)

// -------------------------------------------- Types --------------------------------------------

// UserService handles user-related operations
type UserService struct {
	users map[string]*User
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*User),
	}
}

// -------------------------------------------- Public Functions --------------------------------------------

// GetUser retrieves a user by username
// Original function in the service package
func (us *UserService) GetUser(username string) (*User, error) {
	log.Printf("   👨‍💼 [BUSINESS] Retrieving user: %s", username)

	// Simulate database lookup
	time.Sleep(50 * time.Millisecond) // Simulate DB delay

	user, exists := us.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	log.Printf("   ✅ [BUSINESS] User retrieved successfully")
	return user, nil
}

// CreateUser creates a new user
func (us *UserService) CreateUser(user *User) error {
	log.Printf("   👨‍💼 [BUSINESS] Creating user: %s", user.Username)

	// Simulate database insert
	time.Sleep(75 * time.Millisecond) // Simulate DB delay

	us.users[user.Username] = user
	log.Printf("   ✅ [BUSINESS] User created successfully")
	return nil
}

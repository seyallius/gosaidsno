package main

import "time"

type User struct {
	ID       string
	Username string
	Email    string
	Created  time.Time
}

type Order struct {
	ID     string
	UserID string
	Amount float64
	Date   time.Time
}

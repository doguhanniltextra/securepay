package models

import "time"

// PaymentStatus represents the status of a payment transaction
type PaymentStatus string

const (
	StatusPending   PaymentStatus = "PENDING"
	StatusCompleted PaymentStatus = "COMPLETED"
	StatusFailed    PaymentStatus = "FAILED"
)

// Payment represents a transaction record in the database
type Payment struct {
	ID             string
	FromAccount    string
	ToAccount      string
	Amount         float64
	Currency       string
	Status         string
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Version        int
}

// PaymentInitiatedEvent represents the event structure published to Kafka
type PaymentInitiatedEvent struct {
	PaymentID   string  `json:"payment_id"`
	FromAccount string  `json:"from_account"`
	ToAccount   string  `json:"to_account"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Timestamp   string  `json:"timestamp"`
}

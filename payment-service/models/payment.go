package models

import (
	"fmt"
	"time"
)

// PaymentStatus represents the status of a payment transaction
type PaymentStatus string

const (
	StatusPending   PaymentStatus = "PENDING"
	StatusCompleted PaymentStatus = "COMPLETED"
	StatusFailed    PaymentStatus = "FAILED"
)

// IsValidTransition ensures valid state transitions
// PENDING -> COMPLETED
// PENDING -> FAILED
func IsValidTransition(current, next PaymentStatus) error {
	if current == next {
		return nil
	}
	if current == StatusCompleted || current == StatusFailed {
		return fmt.Errorf("cannot transition from terminal state %s to %s", current, next)
	}
	if current == StatusPending && (next == StatusCompleted || next == StatusFailed) {
		return nil
	}
	return fmt.Errorf("invalid transition from %s to %s", current, next)
}

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

// PaymentResultEvent represents the result produced by account-service.
// Consumed by payment-service to update transaction status.
type PaymentResultEvent struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"` // "COMPLETED" or "FAILED"
	Reason    string `json:"reason,omitempty"`
	Timestamp string `json:"timestamp"`
}


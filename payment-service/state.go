package main

import (
	"fmt"
)

// PaymentStatus defines the status of a payment
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

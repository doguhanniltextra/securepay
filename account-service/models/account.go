package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Account represents the account entity in the database
type Account struct {
	ID        string    `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

// PaymentInitiatedEvent represents the Kafka event payload
type PaymentInitiatedEvent struct {
	PaymentID   string  `json:"payment_id"`
	FromAccount string  `json:"from_account"`
	ToAccount   string  `json:"to_account"`
	Amount      decimal.Decimal `json:"amount"`
	Currency    string  `json:"currency"`
	Timestamp   string  `json:"timestamp"`
}

// PaymentResultEvent represents the result of payment processing.
// Produced by account-service after attempting to process a payment.
type PaymentResultEvent struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"` // "COMPLETED" or "FAILED"
	Reason    string `json:"reason,omitempty"`
	Timestamp string `json:"timestamp"`
}

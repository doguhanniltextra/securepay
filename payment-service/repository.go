package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	// Driver registration
	_ "github.com/lib/pq" 

	pb "securepay/proto/gen/go/payment/v1"
)

// Repository defines the interface for database operations
type Repository interface {
	SavePayment(ctx context.Context, req *pb.InitiatePaymentRequest) error
	GetPayment(ctx context.Context, paymentId string) (*Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentId string, status PaymentStatus) error
}

// Payment represents the payments.transactions table structure
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

// PostgresRepository implements Repository
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// SavePayment saves a new payment to the database with PENDING status and version 1
func (r *PostgresRepository) SavePayment(ctx context.Context, req *pb.InitiatePaymentRequest) error {
	query := `
		INSERT INTO payments.transactions (
			id, from_account, to_account, amount, currency, status, idempotency_key, created_at, updated_at, version
		) VALUES (
			$1, $2, $3, $4, $5, 'PENDING', $6, NOW(), NOW(), 1
		)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		req.PaymentId,
		req.FromAccount,
		req.ToAccount,
		req.Amount,
		req.Currency,
		req.IdempotencyKey,
	)

	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	return nil
}

// GetPayment fetches a payment by ID
func (r *PostgresRepository) GetPayment(ctx context.Context, paymentId string) (*Payment, error) {
	query := `
		SELECT id, from_account, to_account, amount, currency, status, idempotency_key, created_at, updated_at, version
		FROM payments.transactions
		WHERE id = $1
	`
	
	var p Payment
	err := r.db.QueryRowContext(ctx, query, paymentId).Scan(
		&p.ID,
		&p.FromAccount,
		&p.ToAccount,
		&p.Amount,
		&p.Currency,
		&p.Status,
		&p.IdempotencyKey,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &p, nil
}

// UpdatePaymentStatus updates the status of a payment
func (r *PostgresRepository) UpdatePaymentStatus(ctx context.Context, paymentID string, status PaymentStatus) error {
	query := `
		UPDATE payments.transactions 
		SET status = $1, updated_at = NOW(), version = version + 1 
		WHERE id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query, status, paymentID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("payment not found")
	}
	
	return nil
}

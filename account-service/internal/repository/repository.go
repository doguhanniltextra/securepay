package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"securepay/account-service/models"

	_ "github.com/lib/pq"
)

// Repository defines the interface for database operations
type Repository interface {
	GetAccount(ctx context.Context, accountID string) (*models.Account, error)
	UpsertAccount(ctx context.Context, account *models.Account) error
	ProcessPayment(ctx context.Context, fromAccountID, toAccountID string, amount float64) error
}

// PostgresRepository implements Repository
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetAccount fetches account details by ID
func (r *PostgresRepository) GetAccount(ctx context.Context, accountID string) (*models.Account, error) {
	query := `
		SELECT account_id, balance, currency, created_at, updated_at, version
		FROM accounts.balances
		WHERE account_id = $1
	`

	var acc models.Account
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(
		&acc.ID,
		&acc.Balance,
		&acc.Currency,
		&acc.CreatedAt,
		&acc.UpdatedAt,
		&acc.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &acc, nil
}

// UpsertAccount inserts or updates an account (used for seeding)
func (r *PostgresRepository) UpsertAccount(ctx context.Context, account *models.Account) error {
	query := `
		INSERT INTO accounts.balances (account_id, balance, currency, created_at, updated_at, version)
		VALUES ($1, $2, $3, NOW(), NOW(), 1)
		ON CONFLICT (account_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, account.ID, account.Balance, account.Currency)
	if err != nil {
		return fmt.Errorf("failed to upsert account: %w", err)
	}
	return nil
}


// ProcessPayment handles the transactional balance update
func (r *PostgresRepository) ProcessPayment(ctx context.Context, fromAccountID, toAccountID string, amount float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 1. Lock From Account and Check Balance
	var fromBalance float64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM accounts.balances WHERE account_id = $1 FOR UPDATE", fromAccountID).Scan(&fromBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("from_account not found")
		}
		return fmt.Errorf("failed to lock from_account: %w", err)
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// 2. Deduct from From Account
	_, err = tx.ExecContext(ctx, "UPDATE accounts.balances SET balance = balance - $1, version = version + 1, updated_at = NOW() WHERE account_id = $2", amount, fromAccountID)
	if err != nil {
		return fmt.Errorf("failed to deduct balance: %w", err)
	}

	// 3. Add to To Account
	// Also check if To Account exists first or rely on UPDATE returning 0?
	// It's safer to check existence or use UPSERT if allowed, but strict rule is exists.
	// But let's assume it exists as per payment-service validation.
	res, err := tx.ExecContext(ctx, "UPDATE accounts.balances SET balance = balance + $1, version = version + 1, updated_at = NOW() WHERE account_id = $2", amount, toAccountID)
	if err != nil {
		return fmt.Errorf("failed to credit balance: %w", err)
	}
	
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("to_account not found")
	}

	slog.Info("Successfully processed payment", "from", fromAccountID, "to", toAccountID, "amount", amount)
	return nil
}
